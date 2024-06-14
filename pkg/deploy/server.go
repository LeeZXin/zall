package deploy

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/pkg/process"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/gliderlabs/ssh"
	"github.com/kballard/go-shellquote"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	OutputLimit = 1024 * 1024 * 2
)

type handler func(ssh.Session, map[string]string, string, string)

type AgentServer struct {
	*zssh.Server
	token           string
	servicesDir     string
	serviceExecutor *executor.Executor
	cmdMap          *cmdMap
	handlerMap      map[string]handler
}

func (s *AgentServer) CancelAllCmd() {
	cmds := s.cmdMap.GetAll()
	for _, cmd := range cmds {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	time.Sleep(time.Second)
}

func NewAgentServer() zsf.LifeCycle {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("os.Getwd err: %v", err)
	}
	pwdDir := pwd
	agent := new(AgentServer)
	poolSize := static.GetInt("deploy.agent.poolSize")
	if poolSize <= 0 {
		poolSize = 10
	}
	queueSize := static.GetInt("deploy.agent.queueSize")
	if queueSize <= 0 {
		queueSize = 1024
	}
	agent.serviceExecutor, _ = executor.NewExecutor(poolSize, queueSize, time.Minute, executor.AbortStrategy)
	agent.token = static.GetString("deploy.agent.token")
	agent.servicesDir = filepath.Join(pwdDir, "services")
	agent.cmdMap = newCmdMap()
	agent.handlerMap = map[string]handler{
		"execute": func(session ssh.Session, args map[string]string, workdir, tempDir string) {
			id := idutil.RandomUuid()
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
			cmd, err := newCommand(session.Context(), "bash -c "+cmdPath, session, session, workdir)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			if !agent.cmdMap.PutIfAbsent(id, cmd) {
				util.ExitWithErrMsg(session, "duplicated id")
				return
			}
			defer agent.cmdMap.Remove(id)
			err = cmd.Start()
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			err = cmd.Wait()
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			session.Exit(0)
		},
	}
	agentPort := static.GetInt("deploy.agent.port")
	if agentPort <= 0 {
		agentPort = 7777
	}
	serv, err := zssh.NewServer(&zssh.ServerOpts{
		Port:    agentPort,
		HostKey: filepath.Join(pwdDir, "data", "ssh", "deploy.rsa"),
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			if ctx.User() != "deploy" {
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
			//
			service := cmd.Args["s"]
			if service == "" {
				util.ExitWithErrMsg(session, "invalid service")
				return
			}
			workdir := filepath.Join(pwdDir, service)
			// 创建临时目录
			tempDir := filepath.Join(workdir, "temp")
			err = util.Mkdir(tempDir)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			fn(session, cmd.Args, workdir, tempDir)
		},
	})
	if err != nil {
		logger.Logger.Fatal(err)
	}
	agent.Server = serv
	quit.AddShutdownHook(agent.CancelAllCmd)
	return agent
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
	if output.Len() > OutputLimit {
		output.Truncate(OutputLimit)
	}
	return output.String(), nil
}

var (
	clientCfg  *gossh.ClientConfig
	clientOnce sync.Once
)

func initClientCfg() {
	clientOnce.Do(func() {
		pwd, err := os.Getwd()
		if err != nil {
			logger.Logger.Fatal(err)
		}
		hostKey, err := util.ReadOrGenRsaKey(filepath.Join(pwd, "data", "ssh", "deploy.rsa"))
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
		clientCfg = zssh.NewCommonClientConfig("deploy", keySigner)
	})
}

type AgentCommand struct {
	agentHost  string
	agentToken string
	service    string
}

func NewAgentCommand(agentHost, agentToken, service string) *AgentCommand {
	initClientCfg()
	return &AgentCommand{
		agentHost:  agentHost,
		agentToken: agentToken,
		service:    service,
	}
}

func (c *AgentCommand) Execute(cmd io.Reader, envs map[string]string) (string, error) {
	return execute(
		c.agentHost,
		fmt.Sprintf("execute -t %s -s %s", c.agentToken, c.service),
		cmd,
		envs,
	)
}
