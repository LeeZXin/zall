package deploymd

import (
	"context"
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
		Table("zservice_deploy_config_" + reqDTO.Env).
		Insert(Config{
			AppId:       reqDTO.AppId,
			Name:        reqDTO.Name,
			Content:     reqDTO.Content,
			ServiceType: reqDTO.ServiceType,
		})
	return err
}

func UpdateConfig(ctx context.Context, reqDTO UpdateConfigReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_config_"+reqDTO.Env).
		Where("id = ?", reqDTO.ConfigId).
		Cols("content", "name").
		Limit(1).
		Update(Config{
			Name:    reqDTO.Name,
			Content: reqDTO.Content,
		})
	return rows == 1, err
}

func GetConfigById(ctx context.Context, configId int64, env string) (Config, bool, error) {
	var ret Config
	b, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_config_"+env).
		Where("id = ?", configId).
		Get(&ret)
	return ret, b, err
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
		Table("zservice_deploy_config_"+env).
		Where("app_id = ?", appId).
		Find(&ret)
	return ret, err
}

func InsertService(ctx context.Context, reqDTO InsertServiceReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_" + reqDTO.Env).
		Insert(&Service{
			ConfigId:           reqDTO.ConfigId,
			ServiceType:        reqDTO.ServiceType,
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

func InsertPlan(ctx context.Context, reqDTO InsertPlanReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_plan_" + reqDTO.Env).
		Insert(&Plan{
			Name:       reqDTO.Name,
			PlanStatus: reqDTO.PlanStatus,
			TeamId:     reqDTO.TeamId,
			Creator:    reqDTO.Creator,
		})
	return err
}

func InsertDeployLog(ctx context.Context, reqDTO InsertDeployLogReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_log_" + reqDTO.Env).
		Insert(&DeployLog{
			ConfigId:       reqDTO.ConfigId,
			AppId:          reqDTO.AppId,
			ServiceType:    reqDTO.ServiceType,
			ServiceConfig:  reqDTO.ServiceConfig,
			ProductVersion: reqDTO.ProductVersion,
			DeployOutput:   reqDTO.DeployOutput,
			Operator:       reqDTO.Operator,
			PlanId:         reqDTO.PlanId,
		})
	return err
}

func GetTeamConfig(ctx context.Context, teamId int64, env string) (TeamConfig, bool, error) {
	var ret TeamConfig
	b, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_team_config_"+env).
		Where("team_id = ?", teamId).
		Get(&ret)
	return ret, b, err
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
