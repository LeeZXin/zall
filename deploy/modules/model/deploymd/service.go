package deploymd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsConfigNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsPlanNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertConfig(ctx context.Context, reqDTO InsertConfigReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(Config{
			AppId:   reqDTO.AppId,
			Name:    reqDTO.Name,
			Content: reqDTO.Content,
			Env:     reqDTO.Env,
		})
	return err
}

func UpdateConfig(ctx context.Context, reqDTO UpdateConfigReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.ConfigId).
		Cols("content", "name").
		Limit(1).
		Update(Config{
			Name:    reqDTO.Name,
			Content: reqDTO.Content,
		})
	return rows == 1, err
}

func DeleteConfigById(ctx context.Context, configId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", configId).
		Delete(new(Config))
	return rows == 1, err
}

func GetConfigById(ctx context.Context, configId int64) (Config, bool, error) {
	var ret Config
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", configId).
		Get(&ret)
	return ret, b, err
}

func ListConfigByAppId(ctx context.Context, appId, env string) ([]Config, error) {
	ret := make([]Config, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		Desc("id").
		Find(&ret)
	return ret, err
}

func InsertPlanService(ctx context.Context, reqDTO InsertPlanServiceReqDTO) (PlanService, error) {
	ret := PlanService{
		ConfigId:           reqDTO.ConfigId,
		CurrProductVersion: reqDTO.CurrProductVersion,
		LastProductVersion: reqDTO.LastProductVersion,
		DeployConfig:       reqDTO.DeployConfig,
		PlanId:             reqDTO.PlanId,
		ServiceStatus:      reqDTO.Status,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetServiceByConfigId(ctx context.Context, configId int64) (PlanService, bool, error) {
	var ret PlanService
	b, err := xormutil.MustGetXormSession(ctx).
		Where("config_id = ?", configId).
		Get(&ret)
	return ret, b, err
}

func GetServiceById(ctx context.Context, serviceId int64) (PlanService, bool, error) {
	var ret PlanService
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", serviceId).
		Get(&ret)
	return ret, b, err
}

func InsertPlan(ctx context.Context, reqDTO InsertPlanReqDTO) (Plan, error) {
	ret := Plan{
		Name:     reqDTO.Name,
		IsClosed: reqDTO.IsClosed,
		TeamId:   reqDTO.TeamId,
		Env:      reqDTO.Env,
		Creator:  reqDTO.Creator,
		Expired:  reqDTO.Expired,
	}
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&ret)
	return ret, err
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

func IterateService(ctx context.Context, env string, statusList []ServiceStatus, fn func(*PlanService) error) error {
	return xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+env).
		In("active_status", statusList).
		Iterate(new(PlanService), func(_ int, bean interface{}) error {
			return fn(bean.(*PlanService))
		})
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

func BatchGetConfigById(ctx context.Context, idList []int64) ([]Config, error) {
	ret := make([]Config, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", idList).
		Find(&ret)
	return ret, err
}

func BatchGetSimpleConfigById(ctx context.Context, idList []int64) ([]Config, error) {
	ret := make([]Config, 0)
	err := xormutil.MustGetXormSession(ctx).
		Cols("id", "app_id", "name").
		In("id", idList).
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

func ClosePlan(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("is_closed = 0").
		Cols("is_closed").
		Update(&Plan{
			IsClosed: true,
		})
	return rows == 1, err
}

func ListPlan(ctx context.Context, reqDTO ListPlanReqDTO) ([]Plan, int64, error) {
	ret := make([]Plan, 0)
	total, err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId).
		And("env = ?", reqDTO.Env).
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		Desc("id").
		FindAndCount(&ret)
	return ret, total, err
}

func BatchInsertDeployStep(ctx context.Context, reqDTOList ...InsertDeployStepReqDTO) error {
	steps, _ := listutil.Map(reqDTOList, func(t InsertDeployStepReqDTO) (Step, error) {
		return Step{
			ServiceId:  t.ServiceId,
			StepIndex:  t.StepIndex,
			Agent:      t.Agent,
			StepStatus: t.StepStatus,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(steps)
	return err
}

func ListPlanServiceByPlanId(ctx context.Context, planId int64) ([]PlanService, error) {
	ret := make([]PlanService, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("plan_id = ?", planId).
		Find(&ret)
	return ret, err
}

func ExistPlanServiceByConfigId(ctx context.Context, configId int64, statusList []ServiceStatus) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("config_id = ?", configId).
		In("service_status", statusList).
		Exist(new(PlanService))
}

func UpdateServiceStatusByIdWithOldStatus(ctx context.Context, serviceId int64, newStatus, oldStatus ServiceStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", serviceId).
		And("service_status = ?", oldStatus).
		Cols("service_status").
		Update(&PlanService{
			ServiceStatus: newStatus,
		})
	return rows == 1, err
}

func DeletePendingPlanServiceById(ctx context.Context, serviceId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", serviceId).
		And("service_status = ?", PendingServiceStatus).
		Delete(new(PlanService))
	return rows == 1, err
}

func UpdateStepStatusAndInputArgsWithOldStatus(ctx context.Context, serviceId int64, index int, agent string, inputArgs map[string]string, newStatus, oldStatus StepStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("service_id = ?", serviceId).
		And("step_index = ?", index).
		And("agent = ?", agent).
		And("step_status = ?", oldStatus).
		Cols("step_status", "input_args").
		Update(&Step{
			StepStatus: newStatus,
			InputArgs: &xormutil.Conversion[map[string]string]{
				Data: inputArgs,
			},
		})
	return rows == 1, err
}

func UpdateStepStatusAndExecuteLogWithOldStatus(ctx context.Context, serviceId int64, index int, agent, log string, newStatus, oldStatus StepStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("service_id = ?", serviceId).
		And("step_index = ?", index).
		And("agent = ?", agent).
		And("step_status = ?", oldStatus).
		Cols("step_status", "execute_log").
		Update(&Step{
			ExecuteLog: log,
			StepStatus: newStatus,
		})
	return rows == 1, err
}

func UpdateRollbackLogWithOldStatus(ctx context.Context, serviceId int64, index int, agent, log string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("service_id = ?", serviceId).
		And("step_index = ?", index).
		And("agent = ?", agent).
		And("step_status = ?", RollbackStepStatus).
		Cols("rollback_log").
		Update(&Step{
			RollbackLog: log,
		})
	return rows == 1, err
}

func ListStepByServiceIdAndLessThanIndex(ctx context.Context, serviceId int64, index int) ([]Step, error) {
	ret := make([]Step, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("service_id = ?", serviceId).
		And("step_index < ?", index).
		Asc("id").
		Find(&ret)
	return ret, err
}

func GetStepByStepId(ctx context.Context, stepId int64) (Step, bool, error) {
	var ret Step
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", stepId).Get(&ret)
	return ret, b, err
}
