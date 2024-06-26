package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

var (
	runner *executor.Executor
)

func initRunner() {
	runner, _ = executor.NewExecutor(10, 1024, time.Minute, executor.AbortStrategy)
}

func executeDeployOnStartPlanService(planId int64, appId string, dp deploy.Service, env map[string]string) error {
	return runner.Execute(func() {
		for index, stage := range dp.Deploy {
			if stage.Confirm.NeedInteract {
				break
			}
			err := runStage(dp, planId, appId, index, env)
			if err != nil {
				break
			}
		}
	})
}

func executeDeployByIndex(planId int64, appId string, dp deploy.Service, index int, env map[string]string) error {
	return runner.Execute(func() {
		runStage(dp, planId, appId, index, env)
	})
}

func runStage(dp deploy.Service, planId int64, appId string, index int, input map[string]string) error {
	stage := dp.Deploy[index]
	if len(stage.Agents) > 0 {
		for _, id := range stage.Agents {
			agent := dp.Agents[id]
			if !updateStageStatusAndInputArgs(planId, index, id, deploymd.RunningStageStatus, deploymd.PendingStageStatus, input) {
				break
			}
			log, err := agent.RunScript(dp.Actions[stage.Confirm.Action].Script, appId, input)
			if err != nil {
				updateStageStatusAndLog(planId, index, id, err.Error(), deploymd.FailStageStatus, deploymd.RunningStageStatus)
				return err
			}
			updateStageStatusAndLog(planId, index, id, log, deploymd.SuccessStageStatus, deploymd.RunningStageStatus)
		}
	} else {
		for id, agent := range dp.Agents {
			if !updateStageStatusAndInputArgs(planId, index, id, deploymd.RunningStageStatus, deploymd.PendingStageStatus, input) {
				break
			}
			log, err := agent.RunScript(dp.Actions[stage.Confirm.Action].Script, appId, input)
			if err != nil {
				updateStageStatusAndLog(planId, index, id, err.Error(), deploymd.FailStageStatus, deploymd.RunningStageStatus)
				return err
			}
			updateStageStatusAndLog(planId, index, id, log, deploymd.SuccessStageStatus, deploymd.RunningStageStatus)
		}
	}
	// 自动完成发布计划
	if index == len(dp.Deploy)-1 {
		closePlan(planId)
	}
	return nil
}

func closePlan(planId int64) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	err := xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := deploymd.ClosePlan(ctx, planId, deploymd.RunningPlanStatus)
		if err2 != nil {
			return err2
		}
		return deploymd.UpdateDeployServiceIsPlanDoneTrue(ctx, planId)
	})
	if err != nil {
		logger.Logger.Error(err)
	}

}

func updateStageStatusAndInputArgs(serviceId int64, index int, agent string, newStatus, oldStatus deploymd.StageStatus, env map[string]string) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStageStatusAndInputArgsWithOldStatus(ctx, serviceId, index, agent, env, newStatus, oldStatus)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	return b
}

func updateStageStatusAndLog(serviceId int64, index int, agent, log string, newStatus, oldStatus deploymd.StageStatus) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStageStatusAndExecuteLogWithOldStatus(ctx, serviceId, index, agent, log, newStatus, oldStatus)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	return b
}

func updateRollbackLog(serviceId int64, index int, agent, log string) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateRollbackLogWithOldStatus(ctx, serviceId, index, agent, log)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	return b
}
