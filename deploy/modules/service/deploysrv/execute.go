package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"sync"
	"sync/atomic"
	"time"
)

var (
	runner *executor.Executor
)

func initRunner() {
	runner, _ = executor.NewExecutor(10, 1024, time.Minute, executor.AbortStrategy)
}

func executeDeployOnStartPlan(planId int64, appId string, dp deploy.Pipeline, env map[string]string, taskIdMapList []map[string]string, varsMap map[string]string) error {
	return runner.Execute(func() {
		for index, stage := range dp.Deploy {
			if stage.Confirm != nil && stage.Confirm.NeedInteract ||
				!runStage(dp, planId, appId, index, env, taskIdMapList[index], varsMap) {
				break
			}
			// 自动完成发布计划
			if index == len(dp.Deploy)-1 {
				setPlanStatusSuccessful(planId)
			}
		}
	})
}

func executeDeployOnConfirmStage(planId int64, appId string, dp deploy.Pipeline, env, varsMap map[string]string, startIndex int) error {
	return runner.Execute(func() {
		for index := startIndex; index < len(dp.Deploy); index++ {
			stage := dp.Deploy[index]
			if index > startIndex && stage.Confirm != nil && stage.Confirm.NeedInteract {
				break
			}
			taskIdMap := getTaskIdMap(planId, index)
			if taskIdMap == nil ||
				!runStage(dp, planId, appId, index, env, taskIdMap, varsMap) {
				break
			}
			// 自动完成发布计划
			if index == len(dp.Deploy)-1 {
				setPlanStatusSuccessful(planId)
			}
		}
	})
}

func redoAgentStage(planId int64, appId string, dp deploy.Pipeline, index int, agentId string, env, varsMap map[string]string, taskId string) bool {
	if runAgentScript(agentId, dp.Agents[agentId], dp.Deploy[index], appId, planId, index, env, varsMap, dp, taskId) {
		// 判断是否自动执行下一个节点或自动完成
		if isStageAllDone(planId, index) {
			// 自动完成发布计划
			if index == len(dp.Deploy)-1 {
				setPlanStatusSuccessful(planId)
			} else {
				// 自动执行下一个节点
				for i := index + 1; i < len(dp.Deploy); i++ {
					stage := dp.Deploy[i]
					if stage.Confirm != nil && stage.Confirm.NeedInteract {
						break
					}
					taskIdMap := getTaskIdMap(planId, i)
					if taskIdMap == nil ||
						!runStage(dp, planId, appId, i, env, taskIdMap, varsMap) {
						break
					}
					// 自动完成发布计划
					if i == len(dp.Deploy)-1 {
						setPlanStatusSuccessful(planId)
					}
				}
			}
		}
		return true
	}
	return false
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

func getTaskIdMap(planId int64, index int) map[string]string {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	ret, err := deploymd.GetStageTaskIdMap(ctx, planId, index)
	if err != nil {
		logger.Logger.Error(err)
		return nil
	}
	return ret
}

type idAgent struct {
	Id    string
	Agent deploy.Agent
}

func runStage(dp deploy.Pipeline, planId int64, appId string, index int, args, taskIdMap, varsMap map[string]string) bool {
	stage := dp.Deploy[index]
	finalHasErr := false
	agentMapList := make([]idAgent, 0)
	if len(stage.Agents) > 0 {
		for _, id := range stage.Agents {
			agentMapList = append(agentMapList, idAgent{
				Id:    id,
				Agent: dp.Agents[id],
			})
		}
	} else {
		for id, agent := range dp.Agents {
			agentMapList = append(agentMapList, idAgent{
				Id:    id,
				Agent: agent,
			})
		}
	}
	if stage.Parallel <= 1 {
		var hasErr bool
		for _, ia := range agentMapList {
			if !runAgentScript(ia.Id, ia.Agent, stage, appId, planId, index, args, varsMap, dp, taskIdMap[ia.Id]) {
				hasErr = true
			}
		}
		finalHasErr = hasErr
	} else {
		// 并行执行
		parallel, _ := executor.NewExecutor(stage.Parallel, 0, time.Minute, executor.StillQueuedStrategy)
		var (
			wg     sync.WaitGroup
			hasErr atomic.Bool
		)
		wg.Add(len(agentMapList))
		for _, ia := range agentMapList {
			fid := ia.Id
			fagent := ia.Agent
			parallel.Execute(func() {
				defer wg.Done()
				if !runAgentScript(fid, fagent, stage, appId, planId, index, args, varsMap, dp, taskIdMap[fid]) {
					hasErr.Store(true)
				}
			})
		}
		wg.Wait()
		parallel.Shutdown()
		finalHasErr = hasErr.Load()
	}
	return !finalHasErr
}

func runAgentScript(id string, agent deploy.Agent, stage deploy.Stage, appId string, planId int64, index int, args, varsMap map[string]string, dp deploy.Pipeline, taskId string) bool {
	if !updateStageStatusAndInputArgs(planId, index, id, args, deploymd.RunningStageStatus, deploymd.PendingStageStatus) {
		return false
	}
	log, err := agent.RunScript(dp.Actions[stage.Action].Script, appId, util.MergeMap(args, varsMap), taskId)
	if err != nil {
		updateStageStatusAndLog(planId, index, id, err.Error(), deploymd.FailedStageStatus, deploymd.RunningStageStatus)
		return false
	}
	updateStageStatusAndLog(planId, index, id, log, deploymd.SuccessfulStageStatus, deploymd.RunningStageStatus)
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

func updateStageStatusAndInputArgs(planId int64, index int, agent string, args map[string]string, newStatus, oldStatus deploymd.StageStatus) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStageStatusAndInputArgsWithOldStatus(ctx, planId, index, agent, args, newStatus, oldStatus)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	return b
}

func updateStageStatusAndLog(planId int64, index int, agent, log string, newStatus, oldStatus deploymd.StageStatus) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStageStatusAndExecuteLogWithOldStatus(ctx, planId, index, agent, log, newStatus, oldStatus)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	return b
}
