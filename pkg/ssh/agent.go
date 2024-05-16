package ssh

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/git/process"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/gliderlabs/ssh"
	"github.com/kballard/go-shellquote"
	"github.com/spf13/cast"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
)

var (
	validCommandIdRegexp = regexp.MustCompile(`^\S+$`)
)

type handler func(ssh.Session, map[string]string, string, string)

var (
	serviceDir string
	pwdDir     string
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("os.Getwd err: %v", err)
	}
	pwdDir = pwd
	serviceDir = filepath.Join(pwd, "data", "services")
	util.MkdirAll(serviceDir)
}

var (
	clientCfg  *gossh.ClientConfig
	clientOnce sync.Once
)

type cmdMap struct {
	sync.Mutex
	container map[string]*exec.Cmd
}

func newCmdMap() *cmdMap {
	return &cmdMap{
		container: make(map[string]*exec.Cmd),
	}
}

func (m *cmdMap) PutIfAbsent(id string, cmd *exec.Cmd) bool {
	m.Lock()
	defer m.Unlock()
	_, b := m.container[id]
	if b {
		return false
	}
	m.container[id] = cmd
	return true
}

func (m *cmdMap) GetById(id string) *exec.Cmd {
	m.Lock()
	defer m.Unlock()
	return m.container[id]
}

func (m *cmdMap) Remove(id string) {
	m.Lock()
	defer m.Unlock()
	delete(m.container, id)
}

type graphMap struct {
	sync.Mutex
	container map[int64]*action.Graph
}

func newGraphMap() *graphMap {
	return &graphMap{
		container: make(map[int64]*action.Graph),
	}
}

func (m *graphMap) PutIfAbsent(id int64, graph *action.Graph) bool {
	m.Lock()
	defer m.Unlock()
	_, b := m.container[id]
	if b {
		return false
	}
	m.container[id] = graph
	return true
}

func (m *graphMap) GetById(id int64) *action.Graph {
	m.Lock()
	defer m.Unlock()
	return m.container[id]
}

func (m *graphMap) Remove(id int64) {
	m.Lock()
	defer m.Unlock()
	delete(m.container, id)
}

type Agent struct {
	*Server
	token      string
	graphMap   *graphMap
	cmdMap     *cmdMap
	handlerMap map[string]handler
}

func NewAgentServer() zsf.LifeCycle {
	agent := new(Agent)
	agent.token = static.GetString("ssh.agent.token")
	agent.graphMap = newGraphMap()
	agent.cmdMap = newCmdMap()
	agent.handlerMap = map[string]handler{
		"executeWorkflow": func(session ssh.Session, args map[string]string, workdir, tempDir string) {
			input, err := io.ReadAll(session)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			envs := util.CutEnv(session.Environ())
			req := executeWorkflowReq{
				PrId:        cast.ToInt64(envs[action.ActionPrId]),
				WfId:        cast.ToInt64(envs[action.ActionWfId]),
				Operator:    envs[action.ActionOperator],
				TriggerType: workflowmd.TriggerType(cast.ToInt(envs[action.ActionTriggerType])),
				Branch:      envs[action.ActionBranch],
				YamlContent: string(input),
			}
			if !req.IsValid() {
				util.ExitWithErrMsg(session, "invalid request")
				return
			}
			err = executeWorkflow(req, agent.graphMap)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
		},
		"killWorkflow": func(session ssh.Session, args map[string]string, workdir, tempDir string) {
			id := cast.ToInt64(args["i"])
			graph := agent.graphMap.GetById(id)
			if graph != nil {
				util.ExitWithErrMsg(session, "unknown id:"+args["i"])
				return
			}
			graph.Cancel()
			agent.graphMap.Remove(id)
			session.Exit(0)
		},
		"kill": func(session ssh.Session, args map[string]string, workdir, tempDir string) {
			id := args["i"]
			if !validCommandIdRegexp.MatchString(id) {
				util.ExitWithErrMsg(session, "invalid id")
				return
			}
			cmd := agent.cmdMap.GetById(id)
			if cmd == nil {
				util.ExitWithErrMsg(session, "unknown id:"+id)
				return
			}
			// 杀死子进程 带负数的pid
			err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			agent.cmdMap.Remove(id)
			session.Exit(0)
		},
		"execute": func(session ssh.Session, args map[string]string, workdir, tempDir string) {
			id := args["i"]
			if !validCommandIdRegexp.MatchString(id) {
				util.ExitWithErrMsg(session, "invalid id")
				return
			}
			cmd := agent.cmdMap.GetById(id)
			if cmd != nil {
				util.ExitWithErrMsg(session, "duplicated id:"+id)
				return
			}
			cmdPath := filepath.Join(tempDir, id)
			file, err := os.OpenFile(cmdPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			defer util.RemoveAll(cmdPath)
			_, err = io.Copy(file, session)
			file.Close()
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			err = executeCommand("chmod +x "+cmdPath, session, workdir)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			cmd, err = newCommand(session.Context(), "bash -c "+cmdPath, session, session, workdir)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			err = cmd.Start()
			if err != nil {
				logger.Logger.Info()
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			if !agent.cmdMap.PutIfAbsent(id, cmd) {
				util.ExitWithErrMsg(session, "duplicated id")
				return
			}
			err = cmd.Wait()
			agent.cmdMap.Remove(id)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			session.Exit(0)
		},
	}
	agentPort := static.GetInt("ssh.agent.port")
	if agentPort <= 0 {
		agentPort = 6666
	}
	serv, err := NewServer(&ServerOpts{
		Port:    agentPort,
		HostKey: filepath.Join(pwdDir, "data", "ssh", "agent.rsa"),
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			if ctx.User() != "workflow" {
				return false
			}
			return true
		},
		SessionHandler: func(session ssh.Session) {
			cmd, err := splitCommand(session.RawCommand())
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			fn, b := agent.handlerMap[cmd.Operation]
			if !b {
				util.ExitWithErrMsg(session, "unrecognized command")
				return
			}
			// token校验
			if cmd.Args["t"] != agent.token {
				util.ExitWithErrMsg(session, "invalid token")
				return
			}
			appId := cmd.Args["s"]
			dir := cmd.Args["w"]
			if appId != "" {
				// appId校验
				if !appmd.IsAppIdValid(appId) {
					util.ExitWithErrMsg(session, "invalid app id")
					return
				}
				// 创建工作目录
				workdir := filepath.Join(serviceDir, appId)
				err = mkdir(workdir)
				if err != nil {
					util.ExitWithErrMsg(session, err.Error())
					return
				}
				// 创建临时目录
				tempDir := filepath.Join(workdir, "temp")
				err = mkdir(tempDir)
				if err != nil {
					util.ExitWithErrMsg(session, err.Error())
					return
				}
				fn(session, cmd.Args, workdir, tempDir)
			} else if dir != "" {
				if !filepath.IsAbs(dir) {
					util.ExitWithErrMsg(session, "invalid -w arg")
					return
				}
				// 创建工作目录
				workdir := dir
				err = mkdir(workdir)
				if err != nil {
					util.ExitWithErrMsg(session, err.Error())
					return
				}
				// 创建临时目录
				tempDir := filepath.Join(workdir, "temp")
				err = mkdir(tempDir)
				if err != nil {
					util.ExitWithErrMsg(session, err.Error())
					return
				}
				fn(session, cmd.Args, workdir, tempDir)
			} else {
				// 没有配置工作目录就用当前agent目录作为工作目录
				workdir := pwdDir
				// 创建临时目录
				tempDir := filepath.Join(workdir, "temp")
				err = mkdir(tempDir)
				if err != nil {
					util.ExitWithErrMsg(session, err.Error())
					return
				}
				fn(session, cmd.Args, workdir, tempDir)
			}
		},
	})
	if err != nil {
		logger.Logger.Fatal(err)
	}
	agent.Server = serv
	return agent
}

func mkdir(dir string) error {
	exist, err := util.IsExist(dir)
	if err != nil {
		return err
	}
	if !exist {
		err = util.Mkdir(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeCommand(line string, session ssh.Session, workdir string) error {
	cmd, err := newCommand(session.Context(), line, session, session, workdir)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func newCommand(ctx context.Context, line string, stdout, stderr io.Writer, workdir string) (*exec.Cmd, error) {
	fields, err := shellquote.Split(line)
	if err != nil {
		return nil, err
	}
	var cmd *exec.Cmd
	if len(fields) > 1 {
		cmd = exec.CommandContext(ctx, fields[0], fields[1:]...)
	} else if len(fields) == 1 {
		cmd = exec.CommandContext(ctx, fields[0])
	} else {
		return nil, fmt.Errorf("empty command")
	}
	process.SetSysProcAttribute(cmd)
	cmd.Env = os.Environ()
	cmd.Dir = workdir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd, nil
}

type command struct {
	Args      map[string]string
	Operation string
}

func splitCommand(cmd string) (command, error) {
	fields, err := shellquote.Split(cmd)
	if err != nil {
		return command{}, err
	}
	flen := len(fields)
	if flen == 0 {
		return command{}, fmt.Errorf("unregonized command: %v", cmd)
	}
	args := make(map[string]string)
	ret := command{
		Operation: fields[0],
		Args:      args,
	}
	i := 1
	for i < flen {
		if fields[i][0] == '-' && i < flen-1 && fields[i+1][0] != '-' {
			args[fields[i][1:]] = fields[i+1]
			i++
		}
		i++
	}
	return ret, nil
}

type Command interface {
	Execute(io.Reader, map[string]string) (string, error)
	ExecuteAsync(io.Reader, map[string]string) (io.Reader, error)
}

func execute(sshHost, command string, cmd io.Reader, envs map[string]string) (string, error) {
	client, err := gossh.Dial("tcp", sshHost, clientCfg)
	if err != nil {
		return "", err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	for k, v := range envs {
		err = session.Setenv(k, v)
		if err != nil {
			return "", err
		}
	}
	stderr := new(bytes.Buffer)
	output := new(bytes.Buffer)
	session.Stdin = cmd
	session.Stdout = output
	session.Stderr = stderr
	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("%w - %s", err, stderr.String())
	}
	return output.String(), nil
}

func executeAsync(sshHost, command string, cmd io.Reader, envs map[string]string) (io.ReadCloser, error) {
	client, err := gossh.Dial("tcp", sshHost, clientCfg)
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	for k, v := range envs {
		err = session.Setenv(k, v)
		if err != nil {
			return nil, err
		}
	}
	stderr := new(bytes.Buffer)
	pipeReader, pipeWriter := io.Pipe()
	session.Stdin = cmd
	session.Stdout = pipeWriter
	session.Stderr = stderr
	go func() {
		err := session.Run(command)
		if err != nil {
			pipeWriter.CloseWithError(fmt.Errorf("%w - %s", err, stderr.String()))
		} else {
			pipeWriter.Close()
		}
		client.Close()
	}()
	return pipeReader, nil
}

type ServiceCommand struct {
	appId      string
	agentHost  string
	agentToken string
}

func (c *ServiceCommand) Execute(cmd io.Reader, envs map[string]string) (string, error) {
	return execute(c.agentHost, fmt.Sprintf("execute -s %s -t %s", c.appId, c.agentToken), cmd, envs)
}

func (c *ServiceCommand) ExecuteAsync(cmd io.Reader, envs map[string]string) (io.Reader, error) {
	return executeAsync(c.agentHost, fmt.Sprintf("execute -s %s -t %s", c.appId, c.agentToken), cmd, envs)
}

func initClientCfg() {
	clientOnce.Do(func() {
		pwd, err := os.Getwd()
		if err != nil {
			logger.Logger.Fatal(err)
		}
		hostKey, err := util.ReadOrGenRsaKey(filepath.Join(pwd, "data", "ssh", "agent.rsa"))
		if err != nil {
			logger.Logger.Fatal(err)
		}
		privateKey, err := os.ReadFile(hostKey)
		if err != nil {
			logger.Logger.Fatal(err)
		}
		keySigner, err := gossh.ParsePrivateKey(privateKey)
		if err != nil {
			logger.Logger.Fatal(err)
		}
		clientCfg = NewCommonClientConfig("workflow", keySigner)
	})
}

func NewServiceCommand(agentHost, agentToken, appId string) *ServiceCommand {
	initClientCfg()
	return &ServiceCommand{
		appId:      appId,
		agentHost:  agentHost,
		agentToken: agentToken,
	}
}

type AgentCommand struct {
	workdir    string
	agentHost  string
	agentToken string
}

func NewAgentCommand(agentHost, agentToken, workdir string) *AgentCommand {
	initClientCfg()
	return &AgentCommand{
		workdir:    workdir,
		agentHost:  agentHost,
		agentToken: agentToken,
	}
}

func (c *AgentCommand) Execute(cmd io.Reader, envs map[string]string) (string, error) {
	return execute(c.agentHost, fmt.Sprintf("execute -w %s -t %s", c.workdir, c.agentToken), cmd, envs)
}

func (c *AgentCommand) ExecuteAsync(cmd io.Reader, envs map[string]string) (io.Reader, error) {
	return executeAsync(c.agentHost, fmt.Sprintf("execute -w %s -t %s", c.workdir, c.agentToken), cmd, envs)
}

func (c *AgentCommand) ExecuteWorkflowAsync(content string, envs map[string]string) error {
	reader, err := executeAsync(c.agentHost,
		fmt.Sprintf("executeWorkflow -w %s -t %s", c.workdir, c.agentToken),
		strings.NewReader(content),
		envs,
	)
	if err != nil {
		return err
	}
	// read nothing
	reader.Close()
	return nil
}
