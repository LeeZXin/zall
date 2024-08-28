package alertmd

import (
	"context"
	"github.com/LeeZXin/zall/pkg/alert"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsConfigNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertConfig(ctx context.Context, reqDTO InsertConfigReqDTO) (Config, error) {
	ret := Config{
		Name: reqDTO.Name,
		Content: &xormutil.Conversion[alert.Alert]{
			Data: reqDTO.Alert,
		},
		AppId:       reqDTO.AppId,
		IntervalSec: reqDTO.IntervalSec,
		IsEnabled:   reqDTO.IsEnabled,
		Env:         reqDTO.Env,
		Creator:     reqDTO.Creator,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateConfig(ctx context.Context, reqDTO UpdateConfigReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "content", "interval_sec", "next_time").
		Update(&Config{
			Name: reqDTO.Name,
			Content: &xormutil.Conversion[alert.Alert]{
				Data: reqDTO.Alert,
			},
			IntervalSec: reqDTO.IntervalSec,
		})
	return rows == 1, err
}

func ListConfig(ctx context.Context, reqDTO ListConfigReqDTO) ([]Config, int64, error) {
	ret := make([]Config, 0)
	total, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("env = ?", reqDTO.Env).
		OrderBy("id desc").
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		FindAndCount(&ret)
	return ret, total, err
}

func GetConfigById(ctx context.Context, id int64) (Config, bool, error) {
	var ret Config
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func DeleteConfigById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(Config))
	return rows == 1, err
}

func UpdateExecuteNextTime(ctx context.Context, configId, nextTime, runVersion int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("config_id = ?", configId).
		And("run_version = ?", runVersion).
		Cols("next_time", "run_version").
		Update(&Execute{
			NextTime:   nextTime,
			RunVersion: runVersion + 1,
		})
	return rows == 1, err
}

func EnableExecute(ctx context.Context, configId int64, nextTime int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("config_id = ?", configId).
		And("is_enabled = 0").
		Cols("is_enabled", "next_time").
		Update(&Execute{
			IsEnabled: true,
			NextTime:  nextTime,
		})
	return rows == 1, err
}

func DisableExecute(ctx context.Context, configId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("config_id = ?", configId).
		And("is_enabled = 1").
		Cols("is_enabled", "next_time").
		Update(&Execute{
			IsEnabled: false,
			NextTime:  0,
		})
	return rows == 1, err
}

func IterateExecute(ctx context.Context, nextTime int64, env string, fn func(*Execute) error) error {
	return xormutil.MustGetXormSession(ctx).
		Where("next_time <= ?", nextTime).
		And("env = ?", env).
		And("is_enabled = 1").
		Iterate(new(Execute), func(idx int, obj interface{}) error {
			if err := fn(obj.(*Execute)); err != nil {
				return err
			}
			return nil
		})
}

func InsertExecute(ctx context.Context, reqDTO InsertExecuteReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Execute{
			ConfigId:  reqDTO.ConfigId,
			IsEnabled: reqDTO.IsEnabled,
			NextTime:  reqDTO.NextTime,
			Env:       reqDTO.Env,
		})
	return err
}

func DeleteExecuteByConfigId(ctx context.Context, configId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("config_id = ?", configId).
		Delete(new(Execute))
	return rows == 1, err
}

func EnableConfig(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("is_enabled = 0").
		Cols("is_enabled").
		Update(&Config{
			IsEnabled: true,
		})
	return rows == 1, err
}

func DisableConfig(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("is_enabled = 1").
		Cols("is_enabled").
		Update(&Config{
			IsEnabled: false,
		})
	return rows == 1, err
}
