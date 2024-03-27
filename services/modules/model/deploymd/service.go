package deploymd

import (
	"context"
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
			CurrProductVersion: reqDTO.CurrProductVersion,
		})
	return err
}

func UpdateService(ctx context.Context, reqDTO UpdateServiceReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_service_"+reqDTO.Env).
		Where("config_id = ?", reqDTO.ConfigId).
		Cols("service_config", "curr_product_version", "last_product_version").
		Update(&Service{
			ServiceConfig:      reqDTO.ServiceConfig,
			CurrProductVersion: reqDTO.CurrProductVersion,
			LastProductVersion: reqDTO.LastProductVersion,
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

func InsertLog(ctx context.Context, reqDTO InsertLogReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_log_" + reqDTO.Env).
		Insert(&Log{
			ConfigId:       reqDTO.ConfigId,
			AppId:          reqDTO.AppId,
			ServiceType:    reqDTO.ServiceType,
			ServiceConfig:  reqDTO.ServiceConfig,
			ProductVersion: reqDTO.ProductVersion,
			DeployOutput:   reqDTO.DeployOutput,
			Operator:       reqDTO.Operator,
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
