package action

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/git/process"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/executor/completable"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	EmptyArgs   = errors.New("empty args")
	ThereHasBug = errors.New("there has bug")
	// 脚本存放路径
	scriptDir = filepath.Join(common.ResourcesDir, "actions")
)

type GraphCfg struct {
	Name string            `json:"name" yaml:"name"`
	Jobs map[string]JobCfg `json:"jobs" yaml:"jobs"`
}

func (c *GraphCfg) String() string {
	return fmt.Sprintf("name: %s, jobs: %v", c.Name, c.Jobs)
}

func (c *GraphCfg) IsValid() error {
	if c.Name == "" || len(c.Jobs) == 0 {
		return EmptyArgs
	}
	allJobNames := hashset.NewHashSet[string]()
	// 检查是否有重复的jobName
	for k, cfg := range c.Jobs {
		if err := cfg.IsValid(); err != nil {
			return err
		}
		// 有重复的名字
		if allJobNames.Contains(k) {
			return fmt.Errorf("job has duplicated name: %v", k)
		}
		allJobNames.Add(k)
	}
	// 检查job needs
	for k, cfg := range c.Jobs {
		for _, n := range cfg.Needs {
			b := allJobNames.Contains(n)
			// 检查jobNeeds 是否存在
			if !b {
				return fmt.Errorf("job node does not exist: %v", n)
			}
			// 检查jobNeeds是否指向自己
			if n == k {
				return fmt.Errorf("job needs point to itself: %v", n)
			}
		}
	}
	// 检查job是否有环
	return c.checkRoundJob()
}

type jobTemp struct {
	Name  string
	Needs *hashset.HashSet[string]
	Next  *hashset.HashSet[string]
}

func newJobTemp(name string) *jobTemp {
	return &jobTemp{
		Name:  name,
		Needs: hashset.NewHashSet[string](),
		Next:  hashset.NewHashSet[string](),
	}
}

func (c *GraphCfg) checkRoundJob() error {
	tmap := make(map[string]*jobTemp, len(c.Jobs))
	for k, cfg := range c.Jobs {
		t := newJobTemp(k)
		if len(cfg.Needs) > 0 {
			t.Needs.Add(cfg.Needs...)
		}
		tmap[k] = t
	}
	for k, cfg := range c.Jobs {
		for _, need := range cfg.Needs {
			tmap[need].Next.Add(k)
		}
	}
	// 寻找深度优先遍历开始节点
	starts := make([]string, 0)
	for k, temp := range tmap {
		if temp.Next.Size() == 0 {
			starts = append(starts, k)
		}
	}
	// 深度优先遍历
	for _, start := range starts {
		if err := c.dfs([]string{}, tmap[start], tmap); err != nil {
			return err
		}
	}
	return nil
}

func (c *GraphCfg) dfs(path []string, t *jobTemp, all map[string]*jobTemp) error {
	if util.FindInSlice(path, t.Name) {
		return fmt.Errorf("round job: %v %v", path, t.Name)
	}
	p := append(path[:], t.Name)
	for _, key := range t.Needs.AllKeys() {
		if err := c.dfs(p, all[key], all); err != nil {
			return err
		}
	}
	return nil
}

func (c *GraphCfg) ConvertToGraph() (*Graph, error) {
	if err := c.IsValid(); err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	// 转换jobs
	jobs := make([]*Job, 0, len(c.Jobs))
	for k, j := range c.Jobs {
		jobs = append(jobs, j.convertToJob(k))
	}
	graphJobs(jobs, c.Jobs)
	return &Graph{
		Name:    c.Name,
		allJobs: jobs,
	}, nil
}

type StepCfg struct {
	Name   string            `json:"name" yaml:"name"`
	Uses   string            `json:"uses" yaml:"uses"`
	With   map[string]string `json:"with" yaml:"with"`
	Script string            `json:"script" yaml:"script"`
}

func (c *StepCfg) String() string {
	return fmt.Sprintf("NodeId: %s uses: %s with: %v", c.Name, c.Uses, c.With)
}

func (c *StepCfg) convertToStep() *Step {
	cpyMap := make(map[string]string, len(c.With))
	for k, v := range c.With {
		cpyMap[k] = v
	}
	return &Step{
		name:   c.Name,
		uses:   c.Uses,
		with:   cpyMap,
		script: c.Script,
	}
}

func (c *StepCfg) IsValid() error {
	if c.Name == "" || (c.Uses == "" && c.Script == "") {
		return EmptyArgs
	}
	return nil
}

type JobCfg struct {
	Needs   []string  `json:"needs" yaml:"needs"`
	Steps   []StepCfg `json:"steps" yaml:"steps"`
	Timeout int64     `json:"timeout" yaml:"timeout"`
}

func (c *JobCfg) String() string {
	return fmt.Sprintf("needs: %v, steps: %v", c.Needs, c.Steps)
}

func (c *JobCfg) IsValid() error {
	if len(c.Steps) == 0 {
		return EmptyArgs
	}
	for _, cfg := range c.Steps {
		if err := cfg.IsValid(); err != nil {
			return err
		}
	}
	return nil
}

func (c *JobCfg) convertToJob(jobName string) *Job {
	steps := make([]*Step, 0, len(c.Steps))
	for _, s := range c.Steps {
		steps = append(steps, s.convertToStep())
	}
	return &Job{
		name:    jobName,
		timeout: time.Duration(c.Timeout) * time.Second,
		steps:   steps,
		needs:   hashset.NewHashSet[*Job](),
		next:    hashset.NewHashSet[*Job](),
	}
}

type RunOpts struct {
	tempDir string
	// 执行前触发
	// err == nil 就不会触发
	BeforeStartFunc func(GraphRunStat) error
	StepOutputFunc  func(StepOutputStat)
	StepAfterFunc   func(error, StepRunStat)
	Args            map[string]string
}

type GraphRunStat struct {
	Name string
}

type StepRunStat struct {
	JobName   string
	Index     int
	Duration  time.Duration
	EventTime time.Time
}

type StepOutputStat struct {
	JobName   string
	Index     int
	EventTime time.Time
	Output    io.ReadCloser
}

type Graph struct {
	Name    string
	allJobs []*Job
}

func (g *Graph) ListJobInfo() []JobInfo {
	ret, _ := listutil.Map(g.allJobs, func(t *Job) (JobInfo, error) {
		steps := make([]StepInfo, 0, len(t.steps))
		for i, step := range t.steps {
			steps = append(steps, StepInfo{
				Index: i,
				Name:  step.name,
			})
		}
		return JobInfo{
			Name:  t.name,
			Steps: steps,
		}, nil
	})
	return ret
}

func (g *Graph) Run(opts RunOpts) error {
	tempDir := filepath.Join(git.ActionDir(), "action-"+util.RandomIdWithTime())
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create tempDir err:%v", err)
	}
	defer util.RemoveAll(tempDir)
	if opts.BeforeStartFunc != nil {
		if err = opts.BeforeStartFunc(GraphRunStat{
			Name: g.Name,
		}); err != nil {
			return err
		}
	}
	opts.tempDir = tempDir
	futures := make(map[string]completable.Future[any])
	// 找到最后一层节点
	layers, _ := listutil.Filter(g.allJobs, func(j *Job) (bool, error) {
		return j.next.Size() == 0, nil
	})
	finalFutures, _ := listutil.Map(layers, func(t *Job) (completable.IBase, error) {
		return loadJob(futures, t, opts), nil
	})
	if len(finalFutures) > 0 {
		// 最后一层的节点就可以不用异步
		future := completable.ThenAllOf(finalFutures...)
		_, err = future.Get()
		return err
	}
	// finalLayers必须大于0 不应该会走到这 否则就是bug
	return ThereHasBug
}

// loadJob 递归调用 从后置节点往前置节点递归整个graph
func loadJob(all map[string]completable.Future[any], j *Job, opts RunOpts) completable.Future[any] {
	// 防止重复执行
	f, b := all[j.name]
	if b {
		return f
	}
	if j.needs.Size() == 0 {
		all[j.name] = completable.CallAsync(func() (any, error) {
			return nil, j.Run(opts)
		})
	} else {
		needs := make([]completable.IBase, 0, j.needs.Size())
		j.needs.Range(func(j *Job) {
			needs = append(needs, loadJob(all, j, opts))
		})
		all[j.name] = completable.CallAsync(func() (any, error) {
			// 等待前置节点执行完，还得执行自己
			allOfAsync := completable.ThenAllOf(needs...)
			_, err := allOfAsync.Get()
			if err != nil {
				return nil, err
			}
			return nil, j.Run(opts)
		})
	}
	return all[j.name]
}

type Step struct {
	name   string
	uses   string
	with   map[string]string
	script string
}

func (s *Step) Run(opts RunOpts, ctx context.Context, j *Job, index int) error {
	var cmd *exec.Cmd
	if s.script != "" {
		cmd = exec.CommandContext(ctx, "bash", "-c", s.script)
	} else {
		cmd = exec.CommandContext(ctx, "bash", filepath.Join(scriptDir, s.uses))
	}
	env := make([]string, 0, len(s.with)+1+len(opts.Args))
	for k, v := range s.with {
		env = append(env, k+"="+v)
	}
	for k, v := range opts.Args {
		env = append(env, k+"="+v)
	}
	env = append(env, "ACTION_STEP_NAME="+s.name)
	if len(s.with) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Dir = opts.tempDir
	var (
		stdoutReader *io.PipeReader
		stdoutWriter *io.PipeWriter
	)
	if opts.StepOutputFunc != nil {
		stdoutReader, stdoutWriter = io.Pipe()
		defer func() {
			stdoutWriter.Close()
		}()
		cmd.Stderr = stdoutWriter
		cmd.Stdout = stdoutWriter
	}
	process.SetSysProcAttribute(cmd)
	beginTime := time.Now()
	err := cmd.Start()
	if err != nil {
		return err
	}
	if opts.StepOutputFunc != nil {
		go opts.StepOutputFunc(StepOutputStat{
			JobName:   j.name,
			Index:     index,
			EventTime: beginTime,
			Output:    stdoutReader,
		})
	}
	err = cmd.Wait()
	endTime := time.Now()
	if opts.StepAfterFunc != nil {
		opts.StepAfterFunc(err, StepRunStat{
			JobName:   j.name,
			Index:     index,
			Duration:  endTime.Sub(beginTime),
			EventTime: beginTime,
		})
	}
	return err
}

type Job struct {
	name    string
	steps   []*Step
	needs   *hashset.HashSet[*Job]
	next    *hashset.HashSet[*Job]
	timeout time.Duration
}

func (j *Job) Run(opts RunOpts) error {
	ctx := context.Background()
	if j.timeout > 0 {
		var cancelFn context.CancelFunc
		ctx, cancelFn = context.WithTimeout(ctx, j.timeout)
		defer cancelFn()
	}
	for i, s := range j.steps {
		if err := s.Run(opts, ctx, j, i); err != nil {
			return err
		}
	}
	return nil
}

func graphJobs(jobs []*Job, c map[string]JobCfg) {
	m := make(map[string]*Job)
	for _, j := range jobs {
		m[j.name] = j
	}
	for k, cfg := range c {
		for _, need := range cfg.Needs {
			m[k].needs.Add(m[need])
			m[need].next.Add(m[k])
		}
	}
}

type JobInfo struct {
	Name  string
	Steps []StepInfo
}

type StepInfo struct {
	Index int
	Name  string
}
