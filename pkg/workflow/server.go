package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/process"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/gliderlabs/ssh"
	"github.com/kballard/go-shellquote"
	gossh "golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	validTaskIdRegexp = regexp.MustCompile(`^\d{10}\S+$`)
)

type handler func(ssh.Session, map[string]string, string, string)

var (
	pwdDir string
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("os.Getwd err: %v", err)
	}
	pwdDir = pwd
}

var (
	clientCfg  *gossh.ClientConfig
	clientOnce sync.Once
)

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

type AgentServer struct {
	*zssh.Server
	token            string
	graphMap         *graphMap
	handlerMap       map[string]handler
	workflowDir      string
	workflowExecutor *executor.Executor
}

func (a *AgentServer) GetWorkflowBaseDir(taskId string) string {
	dateStr := taskId[:8]
	hourStr := taskId[8:10]
	id := taskId[10:]
	return filepath.Join(a.workflowDir, "action", dateStr, hourStr, id)
}

func (a *AgentServer) CancelAllGraph() {
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

type TaskStatusCallbackReq struct {
	Status   string `json:"status"`
	Duration int64  `json:"duration"`
}

func getBaseStatus(store Store) BaseStatus {
	var (
		ret BaseStatus
		err error
	)
	ret.Status, ret.Duration, err = store.ReadStatus()
	if err != nil {
		ret.Status = UnknownStatus
	} else {
		if ret.Status == FailStatus {
			content, _ := store.ReadErrLog()
			if len(content) > 0 {
				ret.ErrLog = content
			}
		}
	}
	beginTime, err := store.ReadBeginTime()
	if err == nil {
		ret.BeginTime = beginTime.UnixMilli()
	}
	return ret
}

func getJobStatus(baseDir string, jobName string, jobCfg action.JobCfg) JobStatus {
	jobDir := filepath.Join(baseDir, jobName)
	store := newFileStore(jobDir)
	if !store.IsExists() {
		return JobStatus{
			JobName: jobName,
			BaseStatus: BaseStatus{
				Status: UnExecuted,
			},
		}
	}
	var ret JobStatus
	ret.JobName = jobName
	ret.BaseStatus = getBaseStatus(store)
	ret.Steps = make([]StepStatus, 0, len(jobCfg.Steps))
	for i, step := range jobCfg.Steps {
		ret.Steps = append(ret.Steps, getStepStatus(baseDir, jobName, i, step.Name))
	}
	return ret
}

func getStepStatus(baseDir string, jobName string, index int, stepName string) StepStatus {
	stepDir := filepath.Join(baseDir, jobName, strconv.Itoa(index))
	store := newFileStore(stepDir)
	if !store.IsExists() {
		return StepStatus{
			StepName: stepName,
		}
	}
	var ret StepStatus
	ret.StepName = stepName
	ret.BaseStatus = getBaseStatus(store)
	return ret
}

func getTaskStatus(baseDir string) TaskStatus {
	var ret TaskStatus
	store := newFileStore(baseDir)
	origin, err := store.ReadOrigin()
	if err == nil {
		var p action.GraphCfg
		// 解析yaml
		err = yaml.Unmarshal(origin, &p)
		if err != nil || p.IsValid() != nil {
			return TaskStatus{}
		}
		jobsMap := make(map[string]action.JobCfg, len(p.Jobs))
		for name, cfg := range p.Jobs {
			jobsMap[name] = cfg
		}
		jobsInDfsOrder := dfsOrder(p)
		ret.JobStatus = make([]JobStatus, 0, len(p.Jobs))
		for _, jobName := range jobsInDfsOrder {
			ret.JobStatus = append(ret.JobStatus, getJobStatus(baseDir, jobName, jobsMap[jobName]))
		}
		//
	}
	ret.BaseStatus = getBaseStatus(store)
	return ret
}

// 返回深度优先遍历顺序
func dfsOrder(p action.GraphCfg) []string {
	type jobNode struct {
		name  string
		next  *hashset.HashSet[string]
		needs *hashset.HashSet[string]
	}
	nodesMap := make(map[string]*jobNode, len(p.Jobs))
	// 初始化nodesMap
	for name, cfg := range p.Jobs {
		nodesMap[name] = &jobNode{
			name:  name,
			next:  hashset.NewHashSet[string](),
			needs: hashset.NewHashSet[string](cfg.Needs...),
		}
	}
	// 补充next
	for name, cfg := range p.Jobs {
		for _, need := range cfg.Needs {
			n, b := nodesMap[need]
			if b {
				n.next.Add(name)
			}
		}
	}
	var dfs func(...*jobNode)
	ret := make([]string, 0, len(p.Jobs))
	visited := make(map[string]bool, len(p.Jobs))
	noNeedsLayers := make([]*jobNode, 0)
	for _, node := range nodesMap {
		if node.needs.Size() == 0 {
			noNeedsLayers = append(noNeedsLayers, node)
		}
	}
	dfs = func(nodes ...*jobNode) {
		for _, node := range nodes {
			if visited[node.name] {
				continue
			}
			ret = append(ret, node.name)
			visited[node.name] = true
			if node.next.Size() > 0 {
				nextNodes := make([]*jobNode, 0, node.next.Size())
				node.next.Range(func(n string) {
					nextNodes = append(nextNodes, nodesMap[n])
				})
				nextNodes, _ = listutil.Filter(nextNodes, func(n *jobNode) (bool, error) {
					return n != nil, nil
				})
				if len(nextNodes) > 0 {
					dfs(nextNodes...)
				}
			}
		}
	}
	dfs(noNeedsLayers...)
	return ret
}

func mkdir(dir string) bool {
	return util.Mkdir(dir) == nil
}

func notifyCallback(callbackUrl, token, taskId string, req any) {
	// 通知回调
	httputil.Post(context.Background(),
		http.DefaultClient,
		callbackUrl+"?taskId="+taskId,
		map[string]string{
			"Authorization": token,
		},
		req,
		nil,
	)
}

func NewAgentServer() zsf.LifeCycle {
	agent := new(AgentServer)
	poolSize := static.GetInt("workflow.agent.poolSize")
	if poolSize <= 0 {
		poolSize = 10
	}
	queueSize := static.GetInt("workflow.agent.queueSize")
	if queueSize <= 0 {
		queueSize = 1024
	}
	agent.workflowExecutor, _ = executor.NewExecutor(poolSize, queueSize, time.Minute, executor.AbortStrategy)
	agent.token = static.GetString("workflow.agent.token")
	agent.graphMap = newGraphMap()
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
			readCloser, err := newFileStore(stepDir).ReadLog()
			if err == nil {
				defer readCloser.Close()
				io.Copy(session, readCloser)
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
			origin, err := newFileStore(baseDir).ReadOrigin()
			if err != nil {
				util.ExitWithErrMsg(session, "unknown id")
				return
			}
			session.Write(origin)
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
			callbackUrl := envs[action.EnvCallBackUrl]
			token := envs[action.EnvCallBackToken]
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
			if !agent.graphMap.PutIfAbsent(taskId, graph) {
				// 不太可能会发生
				graph.Cancel(action.TaskCancelErr)
				util.ExitWithErrMsg(session, "duplicated biz id")
				return
			}
			taskStore := newFileStore(logDir)
			// 首先置为排队状态
			taskStore.StoreStatus(QueueStatus, 0)
			if rErr := agent.workflowExecutor.Execute(func() {
				defer agent.graphMap.Remove(taskId)
				// 写入开始时间
				taskStore.StoreBeginTime(now)
				// 写入原始内容
				taskStore.StoreOrigin(input)
				// 初始状态 执行状态
				taskStore.StoreStatus(RunningStatus, 0)
				// 通知回调
				notifyCallback(callbackUrl, token, taskId, TaskStatusCallbackReq{
					Status: RunningStatus,
				})
				err := graph.Run(action.RunOpts{
					Workdir: filepath.Join(agent.workflowDir, "temp", taskId),
					StepOutputFunc: func(stat action.StepOutputStat) {
						defer stat.Output.Close()
						stepDir := filepath.Join(logDir, stat.JobName, strconv.Itoa(stat.Index))
						if mkdir(stepDir) {
							newFileStore(stepDir).StoreLog(stat.Output)
						}
					},
					JobBeforeFunc: func(stat action.JobBeforeStat) error {
						jobDir := filepath.Join(logDir, stat.JobName)
						err := util.Mkdir(jobDir)
						if err == nil {
							jobStore := newFileStore(jobDir)
							// 记录job开始时间
							jobStore.StoreBeginTime(stat.BeginTime)
							// 设置初始状态
							jobStore.StoreStatus(RunningStatus, 0)
						}
						return err
					},
					JobAfterFunc: func(err error, stat action.JobRunStat) {
						jobDir := filepath.Join(logDir, stat.JobName)
						jobStore := newFileStore(jobDir)
						if err == nil {
							jobStore.StoreStatus(SuccessStatus, stat.Duration)
						} else {
							if err == context.DeadlineExceeded {
								jobStore.StoreStatus(TimeoutStatus, stat.Duration)
							} else {
								jobStore.StoreStatus(FailStatus, stat.Duration)
							}
							jobStore.StoreErrLog(err)
						}
					},
					StepAfterFunc: func(err error, stat action.StepRunStat) {
						stepDir := filepath.Join(logDir, stat.JobName, strconv.Itoa(stat.Index))
						if mkdir(stepDir) {
							stepStore := newFileStore(stepDir)
							// 记录step开始时间
							stepStore.StoreBeginTime(stat.BeginTime)
							if err == nil {
								stepStore.StoreStatus(SuccessStatus, stat.Duration)
							} else {
								stepStore.StoreStatus(FailStatus, stat.Duration)
								stepStore.StoreErrLog(err)
							}
						}
					},
					Args: envs,
				})
				if err != nil {
					graph.Cancel(action.TaskCancelErr)
				}
				var status string
				if err == nil {
					status = SuccessStatus
				} else {
					switch err {
					case context.DeadlineExceeded:
						status = TimeoutStatus
					case action.TaskCancelErr:
						status = CancelStatus
					default:
						status = FailStatus
					}
					taskStore.StoreErrLog(err)
				}
				duration := graph.SinceBeginTime()
				taskStore.StoreStatus(status, duration)
				if callbackUrl != "" {
					if callbackUrl != "" {
						// 通知回调
						notifyCallback(callbackUrl, token, taskId, TaskStatusCallbackReq{
							Status:   status,
							Duration: duration.Milliseconds(),
						})
					}
				}
			}); rErr != nil {
				agent.graphMap.Remove(taskId)
				util.ExitWithErrMsg(session, "out of capacity")
				return
			}
			session.Exit(0)
		},
		"killWorkflow": func(session ssh.Session, args map[string]string, _, _ string) {
			taskId := args["i"]
			graph := agent.graphMap.GetById(taskId)
			if graph == nil {
				util.ExitWithErrMsg(session, "unknown taskId: "+args["i"])
				return
			}
			taskStore := newFileStore(agent.GetWorkflowBaseDir(taskId))
			taskStore.StoreStatus(CancelStatus, graph.SinceBeginTime())
			logger.Logger.Infof("cancel task: %s", taskId)
			graph.Cancel(action.TaskCancelErr)
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
			workdir := pwdDir
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
