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

func executeDeployOnStartPlanService(serviceId int64, appId string, dp deploy.Deploy, env map[string]string) error {
	return runner.Execute(func() {
		for index, stage := range dp.Deploy {
			if stage.Confirm.NeedInteract {
				break
			}
			runStage(dp, serviceId, appId, index, env)
		}
	})
}

func executeDeployByIndex(serviceId int64, appId string, dp deploy.Deploy, index int, env map[string]string) error {
	return runner.Execute(func() {
		runStage(dp, serviceId, appId, index, env)
	})
}

func rollbackStage(serviceId int64, appId string, dp deploy.Deploy, index int, agent deploy.Agent, env map[string]string) error {
	return runner.Execute(func() {
		stage := dp.Deploy[index]
		if updateStepStatusAndInputArgs(serviceId, index, agent.Id, deploymd.SuccessStepStatus, deploymd.RollbackStepStatus, env) {
			log, err := agent.RunScript(stage.Confirm.Script, appId, env)
			if err == nil {
				updateRollbackLog(serviceId, index, agent.Id, log)
			} else {
				updateRollbackLog(serviceId, index, agent.Id, err.Error())
			}
		}
	})
}

func runStage(dp deploy.Deploy, serviceId int64, appId string, index int, input map[string]string) {
	stage := dp.Deploy[index]
	if len(stage.Agents) > 0 {
		agentMap := dp.GetAgentMap()
		result := true
		for _, id := range stage.Agents {
			agent := agentMap[id]
			if updateStepStatusAndInputArgs(serviceId, index, id, deploymd.RunningStepStatus, deploymd.PendingStepStatus, input) {
				log, err := agent.RunScript(stage.Confirm.Script, appId, input)
				result = result && err == nil
				if err == nil {
					updateStepStatusAndLog(serviceId, index, id, log, deploymd.SuccessStepStatus, deploymd.RunningStepStatus)
				} else {
					updateStepStatusAndLog(serviceId, index, id, err.Error(), deploymd.FailStepStatus, deploymd.RunningStepStatus)
				}
			}
		}
	} else {
		result := true
		for _, agent := range dp.Agents {
			if updateStepStatusAndInputArgs(serviceId, index, agent.Id, deploymd.RunningStepStatus, deploymd.PendingStepStatus, input) {
				log, err := agent.RunScript(stage.Confirm.Script, appId, input)
				result = result && err == nil
				if err == nil {
					updateStepStatusAndLog(serviceId, index, agent.Id, log, deploymd.SuccessStepStatus, deploymd.RunningStepStatus)
				} else {
					updateStepStatusAndLog(serviceId, index, agent.Id, err.Error(), deploymd.FailStepStatus, deploymd.RunningStepStatus)
				}
			}
		}
	}
}

func updateStepStatusAndInputArgs(serviceId int64, index int, agent string, newStatus, oldStatus deploymd.StepStatus, env map[string]string) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStepStatusAndInputArgsWithOldStatus(ctx, serviceId, index, agent, env, newStatus, oldStatus)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	return b
}

func updateStepStatusAndLog(serviceId int64, index int, agent, log string, newStatus, oldStatus deploymd.StepStatus) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStepStatusAndExecuteLogWithOldStatus(ctx, serviceId, index, agent, log, newStatus, oldStatus)
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
