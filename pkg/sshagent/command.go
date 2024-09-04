package sshagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	gossh "golang.org/x/crypto/ssh"
	"io"
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

func (c *ServiceCommand) Execute(cmd io.Reader, envs map[string]string, taskId string) (string, error) {
	return execute(
		c.agentHost,
		fmt.Sprintf("execute -t %s -s %s -i %s", c.agentToken, c.service, taskId),
		cmd,
		envs,
	)
}

func (c *ServiceCommand) Kill(taskId string) error {
	_, err := execute(
		c.agentHost,
		fmt.Sprintf("kill -t %s -i %s", c.agentToken, taskId),
		nil,
		nil,
	)
	return err
}
