package ssh

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/pkg/git/process"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
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
	"time"
)

type handler func(ssh.Session, map[string]string, string, string)

var (
	agentToken string
	serviceDir string
	pwdDir     string
	handlerMap map[string]handler
)

var (
	clientCfg  *gossh.ClientConfig
	clientOnce sync.Once
)

func NewAgentServer() zsf.LifeCycle {
	agentToken = static.GetString("ssh.agent.token")
	handlerMap = map[string]handler{
		"execute": func(session ssh.Session, args map[string]string, workdir, tempDir string) {
			cmdPath := filepath.Join(tempDir, idutil.RandomUuid())
			file, err := os.OpenFile(cmdPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			defer func() {
				file.Close()
				util.RemoveAll(cmdPath)
			}()
			_, err = io.Copy(file, session)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			err = executeCommand("chmod +x "+cmdPath, session, workdir)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			err = executeCommand("bash -c "+cmdPath, session, workdir)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			session.Exit(0)
		},
	}
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("os.Getwd err: %v", err)
	}
	pwdDir = pwd
	serviceDir = filepath.Join(pwd, "data", "services")
	util.MkdirAll(serviceDir)
	agentPort := static.GetInt("ssh.agent.port")
	if agentPort <= 0 {
		agentPort = 6666
	}
	serv, err := NewServer(&ServerOpts{
		Port:    agentPort,
		HostKey: filepath.Join(pwd, "data", "ssh", "agent.rsa"),
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
			fn, b := handlerMap[cmd.Operation]
			if !b {
				util.ExitWithErrMsg(session, "unrecognized command")
				return
			}
			// token校验
			if cmd.Args["t"] != agentToken {
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
	return serv
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
	ctx, cancelFunc := context.WithTimeout(session.Context(), time.Hour)
	defer cancelFunc()
	fields, err := shellquote.Split(line)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if len(fields) > 1 {
		cmd = exec.CommandContext(ctx, fields[0], fields[1:]...)
	} else if len(fields) == 1 {
		cmd = exec.CommandContext(ctx, fields[0])
	} else {
		return nil
	}
	process.SetSysProcAttribute(cmd)
	cmd.Env = append(os.Environ(), session.Environ()...)
	cmd.Dir = workdir
	cmd.Stdout = session
	cmd.Stderr = session
	return cmd.Run()
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

func executeAsync(sshHost, command string, cmd io.Reader, envs map[string]string) (io.Reader, error) {
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
	appId string
	cfg   AgentCfg
}

func (c *ServiceCommand) Execute(cmd io.Reader, envs map[string]string) (string, error) {
	return execute(c.cfg.Host, fmt.Sprintf("execute -s %s -t %s", c.appId, c.cfg.Token), cmd, envs)
}

func (c *ServiceCommand) ExecuteAsync(cmd io.Reader, envs map[string]string) (io.Reader, error) {
	return executeAsync(c.cfg.Host, fmt.Sprintf("execute -s %s -t %s", c.appId, c.cfg.Token), cmd, envs)
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

func NewServiceCommand(cfg AgentCfg, appId string) *ServiceCommand {
	initClientCfg()
	return &ServiceCommand{
		appId: appId,
		cfg:   cfg,
	}
}

type AgentCommand struct {
	workdir string
	cfg     AgentCfg
}

func NewAgentCommand(cfg AgentCfg, workdir string) *AgentCommand {
	initClientCfg()
	return &AgentCommand{
		workdir: workdir,
		cfg:     cfg,
	}
}

func (c *AgentCommand) Execute(cmd io.Reader, envs map[string]string) (string, error) {
	return execute(c.cfg.Host, fmt.Sprintf("execute -w %s -t %s", c.workdir, c.cfg.Token), cmd, envs)
}

func (c *AgentCommand) ExecuteAsync(cmd io.Reader, envs map[string]string) (io.Reader, error) {
	return executeAsync(c.cfg.Host, fmt.Sprintf("execute -w %s -t %s", c.workdir, c.cfg.Token), cmd, envs)
}

type AgentCfg struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

func (c *AgentCfg) IsValid() bool {
	return util.IpPortPattern.MatchString(c.Host) && len(c.Token) <= 1024
}

func (c *AgentCfg) FromDB(content []byte) error {
	if c == nil {
		*c = AgentCfg{}
	}
	return json.Unmarshal(content, c)
}

func (c *AgentCfg) ToDB() ([]byte, error) {
	return json.Marshal(c)
}
