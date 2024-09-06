package deploymd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsPlanNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsArtifactVersionValid(version string) bool {
	return len(version) > 0 && len(version) <= 128
}

func IsServiceSourceNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsPipelineVarsNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertPlan(ctx context.Context, reqDTO InsertPlanReqDTO) (Plan, error) {
	ret := Plan{
		AppId:           reqDTO.AppId,
		PipelineId:      reqDTO.PipelineId,
		PipelineName:    reqDTO.PipelineName,
		Name:            reqDTO.Name,
		ArtifactVersion: reqDTO.ArtifactVersion,
		PlanStatus:      reqDTO.PlanStatus,
		Env:             reqDTO.Env,
		Creator:         reqDTO.Creator,
		PipelineConfig:  reqDTO.PipelineConfig,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func DeletePlanByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("app_id = ?", appId).Delete(new(Plan))
	return err
}

func DeleteStageByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("app_id = ?", appId).Delete(new(Stage))
	return err
}

func ExistPendingOrRunningPlanByPipelineId(ctx context.Context, pipelineId int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("pipeline_id = ?", pipelineId).
		In("plan_status", PendingPlanStatus, RunningPlanStatus).
		Exist(new(Plan))
}

func ExistRunningStatusByPlanId(ctx context.Context, planId int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_status = ?", RunningStageStatus).
		Exist(new(Stage))
}

func ExistNotSuccessfulStagesByPlanIdAndIndex(ctx context.Context, planId int64, index int) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index = ?", index).
		And("stage_status != ?", SuccessfulStageStatus).
		Exist(new(Stage))
}

func ListStagesByPlanIdAndLETIndex(ctx context.Context, planId int64, index int) ([]Stage, error) {
	ret := make([]Stage, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index <= ?", index).
		OrderBy("id asc").
		Find(&ret)
	return ret, err
}

func GetPlanById(ctx context.Context, id int64) (Plan, bool, error) {
	var ret Plan
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func UpdatePlanStatusWithOldStatus(ctx context.Context, id int64, newStatus, oldStatus PlanStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("plan_status = ?", oldStatus).
		Cols("plan_status").
		Update(&Plan{
			PlanStatus: newStatus,
		})
	return rows == 1, err
}

func ClosePlanAndUpdateConfig(ctx context.Context, id int64, oldStatus PlanStatus, pipelineConfig string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("plan_status = ?", oldStatus).
		Cols("plan_status", "pipeline_config").
		Update(&Plan{
			PlanStatus:     ClosedPlanStatus,
			PipelineConfig: pipelineConfig,
		})
	return rows == 1, err
}

func StartPlan(ctx context.Context, id int64, pipelineConfig string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("plan_status = ?", PendingPlanStatus).
		Cols("plan_status", "pipeline_config").
		Update(&Plan{
			PlanStatus:     RunningPlanStatus,
			PipelineConfig: pipelineConfig,
		})
	return rows == 1, err
}

func ListPlan(ctx context.Context, reqDTO ListPlanReqDTO) ([]Plan, int64, error) {
	ret := make([]Plan, 0)
	total, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("env = ?", reqDTO.Env).
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		Desc("id").
		FindAndCount(&ret)
	return ret, total, err
}

func ListStageByPlanIdAndStageIndexAndStatus(ctx context.Context, planId int64, stageIndex int, status StageStatus) ([]Stage, error) {
	ret := make([]Stage, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index = ?", stageIndex).
		And("stage_status = ?", status).
		Find(&ret)
	return ret, err
}

func GetStagesByPlanIdAndIndex(ctx context.Context, planId int64, stageIndex int) ([]Stage, error) {
	ret := make([]Stage, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index = ?", stageIndex).
		Find(&ret)
	return ret, err
}

func MergeInputArgsByPlanIdAndLTIndex(ctx context.Context, planId int64, stageIndex int) (map[string]string, error) {
	ret := make([]Stage, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index < ?", stageIndex).
		Cols("input_args").
		Asc("id").
		Find(&ret)
	if err != nil {
		return nil, err
	}
	args := make(map[string]string)
	for _, stage := range ret {
		if stage.InputArgs != nil && stage.InputArgs.Data != nil {
			for k, v := range stage.InputArgs.Data {
				args[k] = v
			}
		}
	}
	return args, err
}

func MergeInputArgsByPlanIdAndLETIndex(ctx context.Context, planId int64, stageIndex int) (map[string]string, error) {
	ret := make([]Stage, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index <= ?", stageIndex).
		Cols("input_args").
		Asc("id").
		Find(&ret)
	if err != nil {
		return nil, err
	}
	args := make(map[string]string)
	for _, stage := range ret {
		if stage.InputArgs != nil && stage.InputArgs.Data != nil {
			for k, v := range stage.InputArgs.Data {
				args[k] = v
			}
		}
	}
	return args, err
}

func BatchInsertDeployStage(ctx context.Context, reqDTOList ...InsertDeployStageReqDTO) ([]Stage, error) {
	stages := listutil.MapNe(reqDTOList, func(t InsertDeployStageReqDTO) Stage {
		return Stage{
			PlanId:      t.PlanId,
			AppId:       t.AppId,
			StageIndex:  t.StageIndex,
			Agent:       t.Agent,
			AgentHost:   t.AgentHost,
			AgentToken:  t.AgentToken,
			StageStatus: t.StageStatus,
			TaskId:      t.TaskId,
			Script:      t.Script,
		}
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(stages)
	return stages, err
}

func UpdateStageStatusAndInputArgsWithOldStatus(ctx context.Context, planId int64, index int, agent string, inputArgs map[string]string, newStatus, oldStatus StageStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index = ?", index).
		And("agent = ?", agent).
		And("stage_status = ?", oldStatus).
		Cols("stage_status", "input_args").
		Update(&Stage{
			StageStatus: newStatus,
			InputArgs: &xormutil.Conversion[map[string]string]{
				Data: inputArgs,
			},
		})
	return rows == 1, err
}

func UpdateStageStatusWithOldStatusById(ctx context.Context, id int64, newStatus, oldStatus StageStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("stage_status = ?", oldStatus).
		Cols("stage_status").
		Update(&Stage{
			StageStatus: newStatus,
		})
	return rows == 1, err
}

func KillStages(ctx context.Context, planId int64, index int, errLog string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index = ?", index).
		In("stage_status", RunningStageStatus, PendingStageStatus).
		Cols("stage_status", "execute_log").
		Update(&Stage{
			StageStatus: FailedStageStatus,
			ExecuteLog:  errLog,
		})
	return err
}

func UpdateStageStatusAndExecuteLogWithOldStatus(ctx context.Context, planId int64, index int, agent, log string, newStatus, oldStatus StageStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index = ?", index).
		And("agent = ?", agent).
		And("stage_status = ?", oldStatus).
		Cols("stage_status", "execute_log").
		Update(&Stage{
			ExecuteLog:  log,
			StageStatus: newStatus,
		})
	return rows == 1, err
}

func GetStageByStageId(ctx context.Context, stageId int64) (Stage, bool, error) {
	var ret Stage
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", stageId).
		Get(&ret)
	return ret, b, err
}

func GetStageListByPlanId(ctx context.Context, planId int64) ([]Stage, error) {
	ret := make([]Stage, 0)
	err := xormutil.MustGetXormSession(ctx).Where("plan_id = ?", planId).Find(&ret)
	return ret, err
}

func IsPipelineNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertPipeline(ctx context.Context, reqDTO InsertPipelineReqDTO) (Pipeline, error) {
	ret := Pipeline{
		AppId:  reqDTO.AppId,
		Name:   reqDTO.Name,
		Config: reqDTO.Config,
		Env:    reqDTO.Env,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func ListPipeline(ctx context.Context, reqDTO ListPipelineReqDTO) ([]Pipeline, error) {
	ret := make([]Pipeline, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("env = ?", reqDTO.Env).
		Desc("id")
	if len(reqDTO.Cols) > 0 {
		session.Cols(reqDTO.Cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func DeletePipelineById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(Pipeline))
	return rows == 1, err
}

func DeletePipelineByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(Pipeline))
	return err
}

func UpdatePipeline(ctx context.Context, reqDTO UpdatePipelineReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("config", "name").
		Update(&Pipeline{
			Name:   reqDTO.Name,
			Config: reqDTO.Config,
		})
	return rows == 1, err
}

func GetPipelineById(ctx context.Context, pipelineId int64) (Pipeline, bool, error) {
	var ret Pipeline
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", pipelineId).Get(&ret)
	return ret, b, err
}

func ListServiceSource(ctx context.Context, reqDTO ListServiceSourceReqDTO) ([]ServiceSource, error) {
	ret := make([]ServiceSource, 0)
	session := xormutil.MustGetXormSession(ctx).
		And("env = ?", reqDTO.Env)
	if len(reqDTO.Cols) > 0 {
		session.Cols(reqDTO.Cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func InsertServiceSource(ctx context.Context, reqDTO InsertServiceSourceReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ServiceSource{
		Name:   reqDTO.Name,
		Env:    reqDTO.Env,
		Host:   reqDTO.Host,
		ApiKey: reqDTO.ApiKey,
	})
	return err
}

func UpdateServiceSource(ctx context.Context, reqDTO UpdateServiceSourceReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "host", "api_key").
		Update(&ServiceSource{
			Name:   reqDTO.Name,
			Host:   reqDTO.Host,
			ApiKey: reqDTO.ApiKey,
		})
	return rows == 1, err
}

func DeleteServiceSourceById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(ServiceSource))
	return rows == 1, err
}

func GetServiceSourceById(ctx context.Context, id int64) (ServiceSource, bool, error) {
	var ret ServiceSource
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func GetAppServiceSourceBindById(ctx context.Context, id int64) (AppServiceSourceBind, bool, error) {
	var ret AppServiceSourceBind
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func BatchGetServiceSourceByIdList(ctx context.Context, idList []int64, cols []string) ([]ServiceSource, error) {
	ret := make([]ServiceSource, 0)
	session := xormutil.MustGetXormSession(ctx).In("id", idList)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func ListPipelineVars(ctx context.Context, appId, env string, cols []string) ([]PipelineVars, error) {
	ret := make([]PipelineVars, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func ListPipelineVarsMap(ctx context.Context, appId, env string) (map[string]string, error) {
	ret, err := ListPipelineVars(ctx, appId, env, []string{"name", "content"})
	if err != nil {
		return nil, err
	}
	varsMap := make(map[string]string, len(ret))
	for _, vars := range ret {
		varsMap[vars.Name] = vars.Content
	}
	return varsMap, nil
}

func InsertPipelineVars(ctx context.Context, reqDTO InsertPipelineVarsReqDTO) (PipelineVars, error) {
	ret := PipelineVars{
		AppId:   reqDTO.AppId,
		Env:     reqDTO.Env,
		Name:    reqDTO.Name,
		Content: reqDTO.Content,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdatePipelineVars(ctx context.Context, reqDTO UpdatePipelineVarsReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("content").
		Update(&PipelineVars{
			Content: reqDTO.Content,
		})
	return rows == 1, err
}

func DeletePipelineVarsById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(PipelineVars))
	return rows == 1, err
}

func DeletePipelineVarsByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(PipelineVars))
	return err
}

func GetPipelineVarsById(ctx context.Context, id int64) (PipelineVars, bool, error) {
	var ret PipelineVars
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ExistPipelineVarsByAppIdAndEnvAndName(ctx context.Context, appId, env, name string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		And("name = ?", name).
		Exist(new(PipelineVars))
}

func BatchInsertAppServiceSourceBind(ctx context.Context, reqDTOs []InsertAppServiceSourceBindReqDTO) error {
	binds := listutil.MapNe(reqDTOs, func(t InsertAppServiceSourceBindReqDTO) AppServiceSourceBind {
		return AppServiceSourceBind{
			SourceId: t.SourceId,
			AppId:    t.AppId,
			Env:      t.Env,
		}
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(binds)
	return err
}

func DeleteAppServiceSourceBindBySourceId(ctx context.Context, sourceId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("source_id = ?", sourceId).
		Delete(new(AppServiceSourceBind))
	return err
}

func DeleteAppServiceSourceBindByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(AppServiceSourceBind))
	return err
}

func DeleteAppServiceSourceBindByAppIdAndEnv(ctx context.Context, appId, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		Delete(new(AppServiceSourceBind))
	return err
}

func ListAppServiceSourceBindByAppIdAndEnv(ctx context.Context, appId, env string) ([]AppServiceSourceBind, error) {
	ret := make([]AppServiceSourceBind, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		Find(&ret)
	return ret, err
}
