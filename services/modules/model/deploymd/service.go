package deploymd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
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
			AppId:     reqDTO.AppId,
			Name:      reqDTO.Name,
			Content:   reqDTO.Content,
			Env:       reqDTO.Env,
			IsEnabled: reqDTO.IsEnabled,
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

func GetConfigById(ctx context.Context, configId int64) (Config, bool, error) {
	var ret Config
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", configId).
		Get(&ret)
	return ret, b, err
}

func BatchSetConfigIsEnabled(ctx context.Context, configIdList []int64, isEnabled bool) error {
	_, err := xormutil.MustGetXormSession(ctx).
		In("id", configIdList).
		Cols("is_enabled").
		Update(&Config{
			IsEnabled: isEnabled,
		})
	return err
}

func ListConfigForUpdate(ctx context.Context, appId, env string) ([]Config, error) {
	ret := make([]Config, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		ForUpdate().
		Find(&ret)
	return ret, err
}

func DeleteConfigById(ctx context.Context, configId int64, env string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_config_"+env).
		Where("id = ?", configId).
		Limit(1).
		Delete(new(Config))
	return rows == 1, err
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

func GetEnabledConfigByAppId(ctx context.Context, appId, env string) (Config, bool, error) {
	var ret Config
	b, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		And("is_enabled = 1").
		Get(&ret)
	return ret, b, err
}

func InsertService(ctx context.Context, reqDTO InsertServiceReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_" + reqDTO.Env).
		Insert(&Service{
			ConfigId:           reqDTO.ConfigId,
			ServiceConfig:      reqDTO.ServiceConfig,
			ActiveStatus:       reqDTO.ActiveStatus,
			StartTime:          reqDTO.StartTime,
			CurrProductVersion: reqDTO.CurrProductVersion,
		})
	return err
}

func UpdateServiceWithOldStatus(ctx context.Context, oldStatus ActiveStatus, reqDTO UpdateServiceReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+reqDTO.Env).
		Where("config_id = ?", reqDTO.ConfigId).
		And("active_status = ?", oldStatus).
		Cols("service_config", "curr_product_version", "last_product_version", "active_status", "probe_time", "start_time").
		Update(&Service{
			ServiceConfig:      reqDTO.ServiceConfig,
			CurrProductVersion: reqDTO.CurrProductVersion,
			LastProductVersion: reqDTO.LastProductVersion,
			ActiveStatus:       reqDTO.ActiveStatus,
			StartTime:          reqDTO.StartTime,
			ProbeTime:          reqDTO.ProbeTime,
		})
	return rows == 1, err
}

func UpdateServiceActiveStatus(ctx context.Context, configId int64, env string, status ActiveStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+env).
		Where("config_id = ?", configId).
		Cols("active_status").
		Update(&Service{
			ActiveStatus: status,
		})
	return rows == 1, err
}

func UpdateServiceActiveStatusAndProbeTimeWithOldStatus(ctx context.Context, configId int64, env string, newStatus, oldStatus ActiveStatus, probeTime int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+env).
		Where("config_id = ?", configId).
		And("active_status = ?", oldStatus).
		Cols("active_status", "probe_time").
		Update(&Service{
			ActiveStatus: newStatus,
			ProbeTime:    probeTime,
		})
	return rows == 1, err
}

func UpdateServiceActiveStatusWithOldStatus(ctx context.Context, configId int64, env string, newStatus, oldStatus ActiveStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+env).
		Where("config_id = ?", configId).
		And("active_status = ?", oldStatus).
		Cols("active_status").
		Update(&Service{
			ActiveStatus: newStatus,
		})
	return rows == 1, err
}

func DeleteServiceByConfigId(ctx context.Context, configId int64, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+env).
		Where("config_id = ?", configId).
		Limit(1).
		Delete(new(Service))
	return err
}

func ListServiceByConfigIdList(ctx context.Context, configIdList []int64, env string) ([]Service, error) {
	ret := make([]Service, 0)
	err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+env).
		In("config_id", configIdList).
		Find(&ret)
	return ret, err
}

func GetServiceByConfigId(ctx context.Context, configId int64, env string) (Service, bool, error) {
	var ret Service
	b, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+env).
		Where("config_id = ?", configId).
		Get(&ret)
	return ret, b, err
}

func InsertPlan(ctx context.Context, reqDTO InsertPlanReqDTO) (Plan, error) {
	ret := Plan{
		Name:       reqDTO.Name,
		PlanStatus: reqDTO.PlanStatus,
		PlanType:   reqDTO.PlanType,
		TeamId:     reqDTO.TeamId,
		Creator:    reqDTO.Creator,
		Expired:    reqDTO.Expired,
	}
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_plan_" + reqDTO.Env).
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

func InsertProbeInstance(ctx context.Context, instanceId, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_probe_instance_" + env).
		Insert(&ProbeInstance{
			InstanceId:    instanceId,
			HeartbeatTime: time.Now().UnixMilli(),
		})
	return err
}

func GetValidProbeInstances(ctx context.Context, env string, after int64) ([]ProbeInstance, error) {
	ret := make([]ProbeInstance, 0)
	err := xormutil.MustGetXormSession(ctx).
		Table("zservice_probe_instance_"+env).
		Where("heartbeat_time > ?", after).
		OrderBy("id asc").
		Find(&ret)
	return ret, err
}

func DeleteProbeInstance(ctx context.Context, instanceId, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_probe_instance_"+env).
		Where("instance_id = ?", instanceId).
		Delete(new(ProbeInstance))
	return err
}

func UpdateProbeInstanceHeartbeatTime(ctx context.Context, instanceId, env string, heartbeatTime int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_probe_instance_"+env).
		Where("instance_id = ?", instanceId).
		Cols("heartbeat_time").
		Update(&ProbeInstance{
			HeartbeatTime: heartbeatTime,
		})
	return rows == 1, err
}

func IterateService(ctx context.Context, env string, statusList []ActiveStatus, fn func(*Service) error) error {
	return xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+env).
		In("active_status", statusList).
		Iterate(new(Service), func(_ int, bean interface{}) error {
			return fn(bean.(*Service))
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

func BatchGetConfigById(ctx context.Context, idList []int64, env string) ([]Config, error) {
	ret := make([]Config, 0)
	err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_config_"+env).
		In("id", idList).
		Find(&ret)
	return ret, err
}

func BatchGetSimpleConfigById(ctx context.Context, idList []int64, env string) ([]Config, error) {
	ret := make([]Config, 0)
	err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_config_"+env).
		Cols("id", "app_id", "name").
		In("id", idList).
		Find(&ret)
	return ret, err
}

func InsertPlanApproval(ctx context.Context, reqDTO InsertPlanApprovalReqDTO) (PlanApproval, error) {
	ret := PlanApproval{
		Name:        reqDTO.Name,
		TeamId:      reqDTO.TeamId,
		DeployItems: reqDTO.DeployItems,
		Creator:     reqDTO.Creator,
	}
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_approval_" + reqDTO.Env).
		Insert(&ret)
	return ret, err
}

func BatchInsertApprovalNotify(ctx context.Context, env string, reqDTOs ...InsertApprovalNotifyReqDTO) error {
	notifyList, _ := listutil.Map(reqDTOs, func(t InsertApprovalNotifyReqDTO) (ApprovalNotify, error) {
		return ApprovalNotify{
			Aid:          t.Aid,
			Account:      t.Account,
			NotifyStatus: t.NotifyStatus,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_approval_notify_" + env).
		Insert(notifyList)
	return err
}

func BatchInsertPlanItem(ctx context.Context, env string, reqDTOs ...InsertPlanItemReqDTO) error {
	itemList, _ := listutil.Map(reqDTOs, func(t InsertPlanItemReqDTO) (PlanItem, error) {
		return PlanItem{
			PlanId:             t.PlanId,
			ConfigId:           t.ConfigId,
			ProductVersion:     t.ProductVersion,
			LastProductVersion: t.LastProductVersion,
			ItemStatus:         t.ItemStatus,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_plan_item_" + env).
		Insert(itemList)
	return err
}

func ListPlanItemByPlanId(ctx context.Context, planId int64, env string) ([]PlanItem, error) {
	ret := make([]PlanItem, 0)
	err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_plan_item_"+env).
		Where("plan_id = ?", planId).
		Find(&ret)
	return ret, err
}

func GetPlanById(ctx context.Context, id int64, env string) (Plan, bool, error) {
	var ret Plan
	b, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_plan_"+env).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func GetPlanItemById(ctx context.Context, id int64, env string) (PlanItem, bool, error) {
	var ret PlanItem
	b, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_plan_item_"+env).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func UpdateItemStatusWithOldStatus(ctx context.Context, id int64, env string, newStatus, oldStatus PlanItemStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Table("zservice_deploy_plan_item_"+env).
		Where("id = ?", id).
		And("item_status = ?", oldStatus).
		Cols("item_status").
		Update(&PlanItem{
			ItemStatus: newStatus,
		})
	return rows == 1, err
}

func UpdatePlanStatusById(ctx context.Context, id int64, env string, newStatus PlanStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Table("zservice_deploy_plan_"+env).
		Where("id = ?", id).
		Cols("plan_status").
		Update(&Plan{
			PlanStatus: newStatus,
		})
	return rows == 1, err
}

func ListPlan(ctx context.Context, reqDTO ListPlanReqDTO) ([]Plan, error) {
	session := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_plan_"+reqDTO.Env).
		Where("team_id = ?", reqDTO.TeamId)
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	ret := make([]Plan, 0)
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}
