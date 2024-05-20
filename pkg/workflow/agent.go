package workflow

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/process"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/gliderlabs/ssh"
	"github.com/kballard/go-shellquote"
	"github.com/spf13/cast"
	gossh "golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	RunningStatus = "running"
	SuccessStatus = "success"
	FailStatus    = "fail"
	TimeoutStatus = "timeout"
	CancelStatus  = "cancel"
	UnknownStatus = "unknown"
)

const (
	originFileName = "origin"
	statusFileName = "status"
	beginFileName  = "begin"
	errLogFileName = "error.log"
	logFileName    = "log"
)

var (
	validCommandIdRegexp = regexp.MustCompile(`^\S+$`)
	validTaskIdRegexp    = regexp.MustCompile(`^\d{10}\S+$`)
	validAppIdRegexp     = regexp.MustCompile("[\\w-]{1,32}")
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
	container map[string]*action.Graph
}

func newGraphMap() *graphMap {
	return &graphMap{
		container: make(map[string]*action.Graph),
	}
}

func (m *graphMap) GetAll() map[string]*action.Graph {
	m.Lock()
	defer m.Unlock()
	ret := make(map[string]*action.Graph, len(m.container))
	for k, v := range m.container {
		ret[k] = v
	}
	return ret
}

func (m *graphMap) PutIfAbsent(id string, graph *action.Graph) bool {
	m.Lock()
	defer m.Unlock()
	_, b := m.container[id]
	if b {
		return false
	}
	m.container[id] = graph
	return true
}

func (m *graphMap) GetById(id string) *action.Graph {
	m.Lock()
	defer m.Unlock()
	return m.container[id]
}

func (m *graphMap) Remove(id string) {
	m.Lock()
	defer m.Unlock()
	delete(m.container, id)
}

type Agent struct {
	*zssh.Server
	token       string
	graphMap    *graphMap
	cmdMap      *cmdMap
	handlerMap  map[string]handler
	workflowDir string
}

func (a *Agent) GetWorkflowBaseDir(taskId string) string {
	dateStr := taskId[:8]
	hourStr := taskId[8:10]
	id := taskId[10:]
	return filepath.Join(a.workflowDir, "action", dateStr, hourStr, id)
}

func (a *Agent) CancelAllGraph() {
	graphs := a.graphMap.GetAll()
	for _, graph := range graphs {
		graph.Cancel(action.TaskCancelErr)
	}
	time.Sleep(time.Second)
}

type BaseStatus struct {
	Status    string `json:"status"`
	Duration  int64  `json:"duration"`
	ErrLog    string `json:"errLog"`
	BeginTime int64  `json:"beginTime"`
}

type TaskStatus struct {
	BaseStatus
	JobStatus []JobStatus `json:"jobStatus"`
}

type JobStatus struct {
	JobName string `json:"jobName"`
	BaseStatus
	Steps []StepStatus `json:"steps"`
}

type StepStatus struct {
	StepName string `json:"stepName"`
	BaseStatus
}

func getBaseStatus(dir string) BaseStatus {
	var ret BaseStatus
	content, err := os.ReadFile(filepath.Join(dir, statusFileName))
	if err != nil {
		ret.Status = UnknownStatus
	} else {
		ret.Status, ret.Duration = convertStatusFileContent(content)
		if ret.Status == FailStatus {
			content, _ = os.ReadFile(filepath.Join(dir, errLogFileName))
			if len(content) > 0 {
				ret.ErrLog = string(content)
			}
		}
	}
	content, _ = os.ReadFile(filepath.Join(dir, beginFileName))
	if len(content) > 0 {
		ret.BeginTime = cast.ToInt64(string(content))
	}
	return ret
}

func getJobStatus(baseDir string, jobName string, jobCfg action.JobCfg) JobStatus {
	jobDir := filepath.Join(baseDir, jobName)
	exist, _ := util.IsExist(jobDir)
	if !exist {
		return JobStatus{
			JobName: jobName,
		}
	}
	var ret JobStatus
	ret.JobName = jobName
	ret.BaseStatus = getBaseStatus(jobDir)
	ret.Steps = make([]StepStatus, 0, len(jobCfg.Steps))
	for i, step := range jobCfg.Steps {
		ret.Steps = append(ret.Steps, getStepStatus(baseDir, jobName, i, step.Name))
	}
	return ret
}

func getStepStatus(baseDir string, jobName string, index int, stepName string) StepStatus {
	stepDir := filepath.Join(baseDir, jobName, strconv.Itoa(index))
	exist, _ := util.IsExist(stepDir)
	if !exist {
		return StepStatus{
			StepName: stepName,
		}
	}
	var ret StepStatus
	ret.StepName = stepName
	ret.BaseStatus = getBaseStatus(stepDir)
	return ret
}

func getTaskStatus(baseDir string) TaskStatus {
	origin, err := os.ReadFile(filepath.Join(baseDir, originFileName))
	if err != nil {
		return TaskStatus{}
	}
	var (
		p   action.GraphCfg
		ret TaskStatus
	)
	// 解析yaml
	err = yaml.Unmarshal(origin, &p)
	if err != nil || p.IsValid() != nil {
		return TaskStatus{}
	}
	ret.JobStatus = make([]JobStatus, 0, len(p.Jobs))
	for jobName, jobCfg := range p.Jobs {
		ret.JobStatus = append(ret.JobStatus, getJobStatus(baseDir, jobName, jobCfg))
	}
	sort.SliceStable(ret.JobStatus, func(i, j int) bool {
		return ret.JobStatus[i].JobName < ret.JobStatus[j].JobName
	})
	ret.BaseStatus = getBaseStatus(baseDir)
	return ret
}

func mkdir(dir string) bool {
	return util.Mkdir(dir) == nil
}

func toStatusMsg(status string, duration time.Duration) string {
	return fmt.Sprintf("%s %d", status, duration.Milliseconds())
}

func toStatusMsgBytes(status string, duration time.Duration) []byte {
	return []byte(toStatusMsg(status, duration))
}

func convertStatusFileContent(content []byte) (string, int64) {
	fields := strings.Fields(strings.TrimSpace(string(content)))
	if len(fields) != 2 {
		return UnknownStatus, 0
	}
	return fields[0], cast.ToInt64(fields[1])
}

func NewAgentServer() zsf.LifeCycle {
	agent := new(Agent)
	agent.token = static.GetString("workflow.agent.token")
	agent.graphMap = newGraphMap()
	agent.cmdMap = newCmdMap()
	agent.workflowDir = filepath.Join(pwdDir, "workflow")
	agent.handlerMap = map[string]handler{
		"getWorkflowStepLog": func(session ssh.Session, args map[string]string, _ string, _ string) {
			taskId := args["i"]
			if !validTaskIdRegexp.MatchString(taskId) {
				util.ExitWithErrMsg(session, "invalid id")
				return
			}
			jobName := args["j"]
			if !action.ValidJobNameRegexp.MatchString(jobName) {
				util.ExitWithErrMsg(session, "invalid job name")
				return
			}
			index := args["n"]
			if index == "" {
				util.ExitWithErrMsg(session, "invalid index")
				return
			}
			stepDir := filepath.Join(agent.GetWorkflowBaseDir(taskId), jobName, index)
			exist, _ := util.IsExist(stepDir)
			if !exist {
				util.ExitWithErrMsg(session, "unknown step")
				return
			}
			content, _ := os.ReadFile(filepath.Join(stepDir, logFileName))
			if len(content) > 0 {
				session.Write(content)
			}
			session.Exit(0)
		},
		"getWorkflowTaskOrigin": func(session ssh.Session, args map[string]string, _ string, _ string) {
			taskId := args["i"]
			if !validTaskIdRegexp.MatchString(taskId) {
				util.ExitWithErrMsg(session, "invalid id")
				return
			}
			baseDir := agent.GetWorkflowBaseDir(taskId)
			file, err := os.ReadFile(filepath.Join(baseDir, originFileName))
			if err != nil {
				util.ExitWithErrMsg(session, "unknown id")
				return
			}
			session.Write(file)
			session.Exit(0)
		},
		"getWorkflowTaskStatus": func(session ssh.Session, args map[string]string, _ string, _ string) {
			taskId := args["i"]
			if !validTaskIdRegexp.MatchString(taskId) {
				util.ExitWithErrMsg(session, "invalid id")
				return
			}
			taskStatus := getTaskStatus(agent.GetWorkflowBaseDir(taskId))
			m, _ := json.Marshal(taskStatus)
			fmt.Fprint(session, string(m)+"\n")
			session.Exit(0)
		},
		"executeWorkflow": func(session ssh.Session, args map[string]string, _ string, _ string) {
			taskId := args["i"]
			if !validTaskIdRegexp.MatchString(taskId) {
				util.ExitWithErrMsg(session, "invalid task id")
				return
			}
			input, err := io.ReadAll(session)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			var p action.GraphCfg
			// 解析yaml
			err = yaml.Unmarshal(input, &p)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			graph, err := p.ConvertToGraph()
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			// 环境变量
			envs := util.CutEnv(session.Environ())
			now := time.Now()
			logDir := agent.GetWorkflowBaseDir(taskId)
			exist, err := util.IsExist(logDir)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			if exist {
				util.ExitWithErrMsg(session, "duplicated biz id")
				return
			}
			err = util.Mkdir(logDir)
			if err != nil {
				util.ExitWithErrMsg(session, err.Error())
				return
			}
			go func() {
				if !agent.graphMap.PutIfAbsent(taskId, graph) {
					graph.Cancel(action.TaskCancelErr)
					return
				}
				defer agent.graphMap.Remove(taskId)
				// 写入开始时间
				if err := util.WriteFile(filepath.Join(logDir, beginFileName), []byte(strconv.FormatInt(now.UnixMilli(), 10))); err != nil {
					return
				}
				// 写入原始内容
				if err := util.WriteFile(filepath.Join(logDir, originFileName), input); err != nil {
					return
				}
				// 初始状态
				if err := util.WriteFile(filepath.Join(logDir, statusFileName), toStatusMsgBytes(RunningStatus, 0)); err != nil {
					return
				}
				err := graph.Run(action.RunOpts{
					Workdir: filepath.Join(agent.workflowDir, "temp", taskId),
					StepOutputFunc: func(stat action.StepOutputStat) {
						defer stat.Output.Close()
						stepDir := filepath.Join(logDir, stat.JobName, strconv.Itoa(stat.Index))
						if mkdir(stepDir) {
							var logFile *os.File
							// 记录日志
							logFile, err = os.OpenFile(filepath.Join(stepDir, logFileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
							if err == nil {
								defer logFile.Close()
								// 增加缓存
								writer := bufio.NewWriter(logFile)
								defer writer.Flush()
								io.Copy(writer, stat.Output)
							}
						}
					},
					JobAfterFunc: func(err error, stat action.JobRunStat) {
						jobDir := filepath.Join(logDir, stat.JobName)
						if mkdir(jobDir) {
							// 记录job开始时间
							util.WriteFile(filepath.Join(jobDir, beginFileName), []byte(strconv.FormatInt(stat.BeginTime.UnixMilli(), 10)))
							var (
								taskContent []byte
							)
							if err == nil {
								taskContent = toStatusMsgBytes(SuccessStatus, stat.Duration)
							} else {
								if err == context.DeadlineExceeded {
									taskContent = toStatusMsgBytes(TimeoutStatus, stat.Duration)
								} else {
									taskContent = toStatusMsgBytes(FailStatus, stat.Duration)
								}
								util.WriteFile(filepath.Join(jobDir, errLogFileName), []byte(err.Error()))
							}
							util.WriteFile(filepath.Join(jobDir, statusFileName), taskContent)
						}
					},
					StepAfterFunc: func(err error, stat action.StepRunStat) {
						stepDir := filepath.Join(logDir, stat.JobName, strconv.Itoa(stat.Index))
						if mkdir(stepDir) {
							// 记录step开始时间
							util.WriteFile(filepath.Join(stepDir, beginFileName), []byte(strconv.FormatInt(stat.BeginTime.UnixMilli(), 10)))
							var stepContent []byte
							if err == nil {
								stepContent = toStatusMsgBytes(SuccessStatus, stat.Duration)
							} else {
								stepContent = toStatusMsgBytes(FailStatus, stat.Duration)
								util.WriteFile(filepath.Join(stepDir, errLogFileName), []byte(err.Error()))
							}
							util.WriteFile(filepath.Join(stepDir, statusFileName), stepContent)
						}
					},
					Args: args,
				})
				var content []byte
				if err == nil {
					content = toStatusMsgBytes(SuccessStatus, graph.SinceBeginTime())
				} else {
					switch err {
					case context.DeadlineExceeded:
						content = toStatusMsgBytes(TimeoutStatus, graph.SinceBeginTime())
					case action.TaskCancelErr:
						content = toStatusMsgBytes(CancelStatus, graph.SinceBeginTime())
					default:
						content = toStatusMsgBytes(FailStatus, graph.SinceBeginTime())
					}
					util.WriteFile(filepath.Join(logDir, errLogFileName), []byte(err.Error()))
				}
				util.WriteFile(filepath.Join(logDir, statusFileName), content)
				callbackUrl := envs[action.EnvCallBackUrl]
				token := envs[action.EnvCallBackToken]
				if callbackUrl != "" {
					// 通知回调
					httputil.Post(context.Background(),
						http.DefaultClient,
						callbackUrl+"?taskId="+taskId,
						map[string]string{
							"Authorization": token,
						},
						getTaskStatus(logDir),
						nil,
					)
				}
			}()
			session.Exit(0)
		},
		"killWorkflow": func(session ssh.Session, args map[string]string, workdir, tempDir string) {
			id := args["i"]
			graph := agent.graphMap.GetById(id)
			if graph == nil {
				util.ExitWithErrMsg(session, "unknown id: "+args["i"])
				return
			}
			logger.Logger.Infof("cancel task: %s", id)
			graph.Cancel(action.TaskCancelErr)
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
				util.ExitWithErrMsg(session, "unknown id: "+id)
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
	agentPort := static.GetInt("workflow.agent.port")
	if agentPort <= 0 {
		agentPort = 6666
	}
	serv, err := zssh.NewServer(&zssh.ServerOpts{
		Port:    agentPort,
		HostKey: filepath.Join(pwdDir, "data", "ssh", "workflow.rsa"),
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
				if !validAppIdRegexp.MatchString(appId) {
					util.ExitWithErrMsg(session, "invalid app id")
					return
				}
				// 创建工作目录
				workdir := filepath.Join(serviceDir, appId)
				err = util.Mkdir(workdir)
				if err != nil {
					util.ExitWithErrMsg(session, err.Error())
					return
				}
				// 创建临时目录
				tempDir := filepath.Join(workdir, "temp")
				err = util.Mkdir(tempDir)
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
				err = util.Mkdir(workdir)
				if err != nil {
					util.ExitWithErrMsg(session, err.Error())
					return
				}
				// 创建临时目录
				tempDir := filepath.Join(workdir, "temp")
				err = util.Mkdir(tempDir)
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
				err = util.Mkdir(tempDir)
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
	quit.AddShutdownHook(agent.CancelAllGraph, true)
	return agent
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
		hostKey, err := util.ReadOrGenRsaKey(filepath.Join(pwd, "data", "ssh", "workflow.rsa"))
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
		clientCfg = zssh.NewCommonClientConfig("workflow", keySigner)
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

func (c *AgentCommand) ExecuteWorkflow(content, bizId string, envs map[string]string) error {
	_, err := execute(c.agentHost,
		fmt.Sprintf("executeWorkflow -t %s -i %s", c.agentToken, bizId),
		strings.NewReader(content),
		envs,
	)
	return err
}

func (c *AgentCommand) KillWorkflow(taskId string) error {
	_, err := execute(c.agentHost,
		fmt.Sprintf("killWorkflow -t %s -i %s", c.agentToken, taskId),
		nil,
		nil,
	)
	return err
}
