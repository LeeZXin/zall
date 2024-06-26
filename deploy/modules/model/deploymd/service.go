package deploymd

import (
	"context"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func IsPlanNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsProductVersionValid(productVersion string) bool {
	return len(productVersion) > 0 && len(productVersion) <= 128
}

func InsertPlan(ctx context.Context, reqDTO InsertPlanReqDTO) (Plan, error) {
	ret := Plan{
		AppId:          reqDTO.AppId,
		ServiceId:      reqDTO.ServiceId,
		Name:           reqDTO.Name,
		ProductVersion: reqDTO.ProductVersion,
		PlanStatus:     reqDTO.PlanStatus,
		Env:            reqDTO.Env,
		Creator:        reqDTO.Creator,
		ServiceConfig:  &reqDTO.ServiceConfig,
	}
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&ret)
	return ret, err
}

func ExistPendingOrRunningPlanByServiceId(ctx context.Context, serviceId int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("service_id = ?", serviceId).
		In("plan_status", PendingPlanStatus, RunningPlanStatus).
		Exist(new(Plan))
}

func InsertDeployLog(ctx context.Context, reqDTO InsertDeployLogReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_log_" + reqDTO.Env).
		Insert(&DeployLog{
			ConfigId:       reqDTO.ConfigId,
			AppId:          reqDTO.AppId,
			ServiceConfig:  reqDTO.ServiceConfig,
			ProductVersion: reqDTO.ProductVersion,
			DeployOutput:   reqDTO.DeployOutput,
			Operator:       reqDTO.Operator,
			PlanId:         reqDTO.PlanId,
		})
	return err
}

func InsertOpLog(ctx context.Context, reqDTO InsertOpLogReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_op_log_" + reqDTO.Env).
		Insert(&OpLog{
			Op:             reqDTO.Op,
			ConfigId:       reqDTO.ConfigId,
			Operator:       reqDTO.Operator,
			ScriptOutput:   reqDTO.ScriptOutput,
			ProductVersion: reqDTO.ProductVersion,
		})
	return err
}

func ListDeployLog(ctx context.Context, reqDTO ListDeployLogReqDTO) ([]DeployLog, error) {
	session := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_log_"+reqDTO.Env).
		Where("config_id = ?", reqDTO.ConfigId)
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	ret := make([]DeployLog, 0)
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}

func ListOpLog(ctx context.Context, reqDTO ListOpLogReqDTO) ([]OpLog, error) {
	session := xormutil.MustGetXormSession(ctx).
		Table("zservice_op_log_"+reqDTO.Env).
		Where("config_id = ?", reqDTO.ConfigId)
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	ret := make([]OpLog, 0)
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}

func GetPlanById(ctx context.Context, id int64) (Plan, bool, error) {
	var ret Plan
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ClosePlan(ctx context.Context, id int64, oldStatus PlanStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("plan_status = ?", oldStatus).
		Cols("plan_status").
		Update(&Plan{
			PlanStatus: ClosedPlanStatus,
		})
	return rows == 1, err
}

func UpdateDeployServiceIsPlanDoneTrue(ctx context.Context, planId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("is_plan_done = 0").
		Cols("is_plan_done").
		Update(&DeployService{
			IsPlanDone: true,
		})
	return err
}

func ExistUnsuccessfulDeployStage(ctx context.Context, planId int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		In("stage_status", PendingStageStatus, RunningStageStatus, FailStageStatus, RollbackStageStatus).
		Exist(new(Stage))
}

func StartPlan(ctx context.Context, id int64, serviceConfig deploy.Service) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("plan_status = ?", PendingPlanStatus).
		Cols("plan_status", "service_config").
		Update(&Plan{
			PlanStatus:    RunningPlanStatus,
			ServiceConfig: &serviceConfig,
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

func BatchInsertDeployStage(ctx context.Context, reqDTOList ...InsertDeployStageReqDTO) error {
	stages, _ := listutil.Map(reqDTOList, func(t InsertDeployStageReqDTO) (Stage, error) {
		return Stage{
			PlanId:      t.PlanId,
			StageIndex:  t.StageIndex,
			Agent:       t.Agent,
			StageStatus: t.StageStatus,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(stages)
	return err
}

func BatchInsertDeployService(ctx context.Context, reqDTOList ...InsertDeployServiceReqDTO) error {
	processes, _ := listutil.Map(reqDTOList, func(t InsertDeployServiceReqDTO) (DeployService, error) {
		return DeployService{
			PlanId:    t.PlanId,
			ServiceId: t.ServiceId,
			Config:    &t.Config,
			Probed:    t.Probed,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(processes)
	return err
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

func UpdateRollbackLogWithOldStatus(ctx context.Context, planId int64, index int, agent, log string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		And("stage_index = ?", index).
		And("agent = ?", agent).
		And("stage_status = ?", RollbackStageStatus).
		Cols("rollback_log").
		Update(&Stage{
			RollbackLog: log,
		})
	return rows == 1, err
}

func ListStageByServiceIdAndLessThanIndex(ctx context.Context, serviceId int64, index int) ([]Stage, error) {
	ret := make([]Stage, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("service_id = ?", serviceId).
		And("stage_index < ?", index).
		Asc("id").
		Find(&ret)
	return ret, err
}

func GetStageByStageId(ctx context.Context, stageId int64) (Stage, bool, error) {
	var ret Stage
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", stageId).Get(&ret)
	return ret, b, err
}

func IsServiceNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertService(ctx context.Context, reqDTO InsertServiceReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Service{
			AppId:       reqDTO.AppId,
			Name:        reqDTO.Name,
			Config:      reqDTO.Config,
			Env:         reqDTO.Env,
			ServiceType: reqDTO.ServiceType,
		})
	return err
}

func ListService(ctx context.Context, reqDTO ListServiceReqDTO) ([]Service, error) {
	ret := make([]Service, 0)
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

func DeleteService(ctx context.Context, serviceId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", serviceId).
		Delete(new(Service))
	return rows == 1, err
}

func UpdateService(ctx context.Context, reqDTO UpdateServiceReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.ServiceId).
		Cols("config", "name", "service_type").
		Update(&Service{
			ServiceType: reqDTO.ServiceType,
			Name:        reqDTO.Name,
			Config:      reqDTO.Config,
		})
	return rows == 1, err
}

func GetServiceById(ctx context.Context, serviceId int64) (Service, bool, error) {
	var ret Service
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", serviceId).Get(&ret)
	return ret, b, err
}

func BatchGetServiceById(ctx context.Context, serviceIdList []int64, cols []string) ([]Service, error) {
	ret := make([]Service, 0)
	session := xormutil.MustGetXormSession(ctx).
		In("id", serviceIdList)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func DeleteDeployServiceByServiceId(ctx context.Context, serviceId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("service_id = ?", serviceId).
		Delete(new(DeployService))
	return err
}

func IterateDeployService(ctx context.Context, fn func(*DeployService) error) error {
	return xormutil.MustGetXormSession(ctx).
		Iterate(new(DeployService), func(_ int, bean interface{}) error {
			return fn(bean.(*DeployService))
		})
}

func UpdateDeployServiceProbed(ctx context.Context, serviceId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", serviceId).
		Cols("probed", "fail_times").
		Update(&DeployService{
			Probed:         time.Now().UnixMilli(),
			ProbeFailTimes: 0,
		})
	return rows == 1, err
}

func IncrDeployServiceProbeFailTimes(ctx context.Context, serviceId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", serviceId).
		Cols("probe_fail_times").
		Incr("probe_fail_times", 1).
		Update(new(DeployService))
	return rows == 1, err
}
