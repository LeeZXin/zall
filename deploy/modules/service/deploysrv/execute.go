package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	runner         *executor.Executor
	initRunnerOnce = sync.Once{}
)

func initRunner() {
	initRunnerOnce.Do(func() {
		runner, _ = executor.NewExecutor(10, 1024, time.Minute, executor.AbortStrategy)
	})
}

func executeDeployOnStartPlan(plan deploymd.Plan, dp deploy.Pipeline, defaultEnvs map[string]string, stageMap map[int][]deploymd.Stage) error {
	initRunner()
	return runner.Execute(func() {
		for index, stage := range dp.Deploy {
			if stage.Confirm != nil && stage.Confirm.NeedInteract ||
				// 执行阶段
				!runStage(dp, index, stageMap[index], defaultEnvs) {
				break
			}
			// 自动完成发布计划
			if index == len(dp.Deploy)-1 {
				setPlanStatusSuccessful(plan.Id)
			}
		}
	})
}

func executeDeployOnConfirmStage(plan deploymd.Plan, dp deploy.Pipeline, defaultEnv map[string]string, startIndex int) error {
	initRunner()
	return runner.Execute(func() {
		for index := startIndex; index < len(dp.Deploy); index++ {
			stage := dp.Deploy[index]
			//如果未来需要交互则停下
			if index > startIndex &&
				stage.Confirm != nil && stage.Confirm.NeedInteract {
				break
			}
			stages := getStagesByPlanIdAndIndex(plan.Id, index)
			if len(stages) == 0 ||
				!runStage(dp, index, stages, defaultEnv) {
				break
			}
			// 自动完成发布计划
			if index == len(dp.Deploy)-1 {
				setPlanStatusSuccessful(plan.Id)
			}
		}
	})
}

func redoAgentStagesInRunner(dp deploy.Pipeline, stages []deploymd.Stage, defaultEnv map[string]string) error {
	initRunner()
	return runner.Execute(func() {
		redoAgentStages(dp, stages, defaultEnv)
	})
}

func redoAgentStages(dp deploy.Pipeline, stages []deploymd.Stage, defaultEnv map[string]string) {
	for _, stage := range stages {
		if updateStageStatusWithOldStatusById(stage.Id, deploymd.PendingStageStatus, stage.StageStatus) {
			redoAgentStage(dp, stage, defaultEnv)
		}
	}
}

func redoAgentStage(dp deploy.Pipeline, stage deploymd.Stage, defaultEnv map[string]string) {
	if runAgentScript(stage, defaultEnv) {
		// 判断是否自动执行下一个节点或自动完成
		if isStageAllDone(stage.PlanId, stage.StageIndex) {
			// 自动完成发布计划
			if stage.StageIndex == len(dp.Deploy)-1 {
				setPlanStatusSuccessful(stage.PlanId)
			} else {
				// 自动执行下一个节点
				for index := stage.StageIndex + 1; index < len(dp.Deploy); index++ {
					if dp.Deploy[index].Confirm != nil && dp.Deploy[index].Confirm.NeedInteract {
						break
					}
					stages := getStagesByPlanIdAndIndex(stage.PlanId, index)
					if len(stages) == 0 ||
						!runStage(dp, index, stages, defaultEnv) {
						break
					}
					// 自动完成发布计划
					if index == len(dp.Deploy)-1 {
						setPlanStatusSuccessful(stage.PlanId)
					}
				}
			}
		}
	}
}

func redoAgentStageInRunner(dp deploy.Pipeline, stage deploymd.Stage, defaultEnv map[string]string) error {
	initRunner()
	return runner.Execute(func() {
		redoAgentStage(dp, stage, defaultEnv)
	})
}

func isStageAllDone(planId int64, index int) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.ExistNotSuccessfulStagesByPlanIdAndIndex(ctx, planId, index)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	return !b
}

func getStagesByPlanIdAndIndex(planId int64, index int) []deploymd.Stage {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	ret, err := deploymd.GetStagesByPlanIdAndIndex(ctx, planId, index)
	if err != nil {
		logger.Logger.Error(err)
	}
	return ret
}

func runStage(dp deploy.Pipeline, index int, stages []deploymd.Stage, args map[string]string) bool {
	if len(stages) == 0 {
		return true
	}
	finalHasErr := false
	parallel := dp.Deploy[index].Parallel
	if parallel <= 1 {
		for _, stage := range stages {
			if !runAgentScript(stage, util.MergeMap(args, dp.Agents[stage.Agent].With)) {
				finalHasErr = true
			}
		}
	} else {
		// 并行执行
		extr, _ := executor.NewExecutor(parallel, 0, time.Minute, executor.StillQueuedStrategy)
		var (
			wg     sync.WaitGroup
			hasErr atomic.Bool
		)
		wg.Add(len(stages))
		for i := range stages {
			stage := stages[i]
			extr.Execute(func() {
				defer wg.Done()
				if !runAgentScript(stage, util.MergeMap(args, dp.Agents[stage.Agent].With)) {
					hasErr.Store(true)
				}
			})
		}
		wg.Wait()
		extr.Shutdown()
		finalHasErr = hasErr.Load()
	}
	return !finalHasErr
}

func runAgentScript(stage deploymd.Stage, args map[string]string) bool {
	if !updateStageStatusAndInput(stage, args, deploymd.RunningStageStatus, deploymd.PendingStageStatus) {
		return false
	}
	log, err := sshagent.NewServiceCommand(stage.AgentHost, stage.AgentToken, stage.AppId).
		Execute(strings.NewReader(stage.Script), args, stage.TaskId)
	if err != nil {
		updateStageStatusAndLog(stage, err.Error(), deploymd.FailedStageStatus, deploymd.RunningStageStatus)
		return false
	}
	updateStageStatusAndLog(stage, log, deploymd.SuccessfulStageStatus, deploymd.RunningStageStatus)
	return true
}

func setPlanStatusSuccessful(planId int64) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, err := deploymd.UpdatePlanStatusWithOldStatus(ctx, planId, deploymd.SuccessfulPlanStatus, deploymd.RunningPlanStatus)
	if err != nil {
		logger.Logger.Error(err)
	}
}

func updateStageStatusAndInput(stage deploymd.Stage, args map[string]string, newStatus, oldStatus deploymd.StageStatus) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStageStatusAndInputArgsWithOldStatus(
		ctx,
		stage.PlanId,
		stage.StageIndex,
		stage.Agent,
		args,
		newStatus,
		oldStatus,
	)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	return b
}

func updateStageStatusAndLog(stage deploymd.Stage, log string, newStatus, oldStatus deploymd.StageStatus) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStageStatusAndExecuteLogWithOldStatus(ctx, stage.PlanId, stage.StageIndex, stage.Agent, log, newStatus, oldStatus)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	return b
}
