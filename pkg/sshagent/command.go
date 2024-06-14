package sshagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/pkg/process"
	"github.com/gliderlabs/ssh"
	"github.com/kballard/go-shellquote"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"os"
	"os/exec"
	"strings"
)

type AgentCommand struct {
	agentHost  string
	agentToken string
}

func NewAgentCommand(agentHost, agentToken string) *AgentCommand {
	initClientCfg()
	return &AgentCommand{
		agentHost:  agentHost,
		agentToken: agentToken,
	}
}

type ServiceCommand struct {
	agentHost  string
	agentToken string
	service    string
}

func NewServiceCommand(agentHost, agentToken, service string) *ServiceCommand {
	initClientCfg()
	return &ServiceCommand{
		agentHost:  agentHost,
		agentToken: agentToken,
		service:    service,
	}
}

func (c *AgentCommand) ExecuteWorkflow(content, bizId string, envs map[string]string) error {
	_, err := execute(c.agentHost,
		fmt.Sprintf("executeWorkflow -t %s -i %s", c.agentToken, bizId),
		strings.NewReader(content),
		envs,
	)
	return err
}

func (c *AgentCommand) GetWorkflowTaskStatus(taskId string) (TaskStatus, error) {
	result, err := execute(c.agentHost,
		fmt.Sprintf("getWorkflowTaskStatus -t %s -i %s", c.agentToken, taskId),
		nil,
		nil,
	)
	if err != nil {
		return TaskStatus{}, err
	}
	var ret TaskStatus
	json.Unmarshal([]byte(result), &ret)
	return ret, nil
}

func (c *AgentCommand) GetLogContent(taskId string, jobName string, stepIndex int) (string, error) {
	return execute(c.agentHost,
		fmt.Sprintf("getWorkflowStepLog -t %s -i %s -j %s -n %d", c.agentToken, taskId, jobName, stepIndex),
		nil,
		nil,
	)
}

func (c *AgentCommand) KillWorkflow(taskId string) error {
	_, err := execute(c.agentHost,
		fmt.Sprintf("killWorkflow -t %s -i %s", c.agentToken, taskId),
		nil,
		nil,
	)
	return err
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

func (c *ServiceCommand) Execute(cmd io.Reader, envs map[string]string) (string, error) {
	return execute(
		c.agentHost,
		fmt.Sprintf("execute -t %s -s %s", c.agentToken, c.service),
		cmd,
		envs,
	)
}
