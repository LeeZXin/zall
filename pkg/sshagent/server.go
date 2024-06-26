package sshagent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/pkg/action"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/gliderlabs/ssh"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

var (
	validTaskIdRegexp *regexp.Regexp
)

type handler func(ssh.Session, map[string]string, string, string)

type AgentServer struct {
	*zssh.Server
	token            string
	graphMap         *graphMap
	handlerMap       map[string]handler
	workflowDir      string
	servicesDir      string
	cmdMap           *cmdMap
	workflowExecutor *executor.Executor
	serviceExecutor  *executor.Executor
}

func (a *AgentServer) GetWorkflowBaseDir(taskId string) string {
	dateStr := taskId[:8]
	hourStr := taskId[8:10]
	id := taskId[10:]
	return filepath.Join(a.workflowDir, "action", dateStr, hourStr, id)
}

func (a *AgentServer) CancelAll() {
	graphs := a.graphMap.GetAll()
	for _, graph := range graphs {
		graph.Cancel(action.TaskCancelErr)
	}
	cmds := a.cmdMap.GetAll()
	for _, cmd := range cmds {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
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
	validTaskIdRegexp = regexp.MustCompile(`^\d{10}\S+$`)
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("os.Getwd err: %v", err)
	}
	pwdDir = pwd
	agent := new(AgentServer)
	poolSize := static.GetInt("ssh.agent.workflow.poolSize")
	if poolSize <= 0 {
		poolSize = 10
	}
	queueSize := static.GetInt("ssh.agent.workflow.queueSize")
	if queueSize <= 0 {
		queueSize = 1024
	}
	agent.workflowExecutor, _ = executor.NewExecutor(poolSize, queueSize, time.Minute, executor.AbortStrategy)
	poolSize = static.GetInt("ssh.agent.service.poolSize")
	if poolSize <= 0 {
		poolSize = 10
	}
	queueSize = static.GetInt("ssh.agent.service.queueSize")
	if queueSize <= 0 {
		queueSize = 1024
	}
	agent.serviceExecutor, _ = executor.NewExecutor(poolSize, queueSize, time.Minute, executor.AbortStrategy)
	agent.token = static.GetString("ssh.agent.token")
	agent.graphMap = newGraphMap()
	agent.cmdMap = newCmdMap()
	agent.workflowDir = filepath.Join(pwdDir, "workflow")
	agent.servicesDir = filepath.Join(pwdDir, "services")
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
		"executeWorkflow": func(session ssh.Session, args map[string]string, workDir string, tempDir string) {
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
					Workdir: filepath.Join(workDir, "temp", taskId),
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
			graph.Cancel(action.TaskCancelErr)
			session.Exit(0)
		},
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
			cmd, err := newCommand(session.Context(), "bash -c "+cmdPath, session, session, workdir, session.Environ())
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
	agentPort := static.GetInt("ssh.agent.port")
	if agentPort <= 0 {
		agentPort = 6666
	}
	serv, err := zssh.NewServer(&zssh.ServerOpts{
		Port:    agentPort,
		HostKey: filepath.Join(pwdDir, "data", "ssh", "sshAgent.rsa"),
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			if ctx.User() != "zall" {
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
			var workdir string
			if cmd.Operation == "execute" {
				service := cmd.Args["s"]
				if service == "" {
					util.ExitWithErrMsg(session, "invalid service")
					return
				}
				workdir = filepath.Join(agent.servicesDir, service)
			} else {
				workdir = agent.workflowDir
			}
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
	quit.AddShutdownHook(agent.CancelAll, true)
	return agent
}
