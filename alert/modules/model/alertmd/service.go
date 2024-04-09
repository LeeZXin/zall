package alertmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func IsConfigNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertConfig(ctx context.Context, reqDTO InsertConfigReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Config{
			Name:        reqDTO.Name,
			Content:     &reqDTO.Alert,
			AppId:       reqDTO.AppId,
			IntervalSec: reqDTO.IntervalSec,
			SilenceSec:  reqDTO.SilenceSec,
			Enabled:     reqDTO.Enabled,
			NextTime:    reqDTO.NextTime,
		})
	return err
}

func UpdateConfig(ctx context.Context, reqDTO UpdateConfigReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "content", "interval_sec", "silence_sec", "enabled", "next_time").
		Update(&Config{
			Name:        reqDTO.Name,
			Content:     &reqDTO.Alert,
			IntervalSec: reqDTO.IntervalSec,
			SilenceSec:  reqDTO.SilenceSec,
			Enabled:     reqDTO.Enabled,
			NextTime:    reqDTO.NextTime,
		})
	return rows == 1, err
}

func ListConfig(ctx context.Context, reqDTO ListConfigReqDTO) ([]Config, error) {
	session := xormutil.MustGetXormSession(ctx).Where("app_id = ?", reqDTO.AppId)
	if reqDTO.Cursor > 0 {
		session.And("id > ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	ret := make([]Config, 0)
	err := session.OrderBy("id asc").Find(&ret)
	return ret, err
}

func GetConfigById(ctx context.Context, id int64) (Config, bool, error) {
	var ret Config
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func DeleteConfig(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(Config))
	return rows == 1, err
}

func UpdateNextTime(ctx context.Context, id int64, nextTime int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("next_time").
		Update(&Config{
			NextTime: nextTime,
		})
	return rows == 1, err
}

func InsertInstance(ctx context.Context, instanceId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Instance{
			InstanceId:    instanceId,
			HeartbeatTime: time.Now().UnixMilli(),
		})
	return err
}

func UpdateHeartbeatTime(ctx context.Context, instanceId string, heartbeatTime int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		Cols("heartbeat_time").
		Update(&Instance{
			HeartbeatTime: heartbeatTime,
		})
	return rows == 1, err
}

func DeleteInstance(ctx context.Context, instanceId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		Limit(1).
		Delete(new(Instance))
	return err
}

func GetValidInstances(ctx context.Context, after int64) ([]Instance, error) {
	ret := make([]Instance, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("heartbeat_time>= ?", after).
		OrderBy("id asc").
		Find(&ret)
	return ret, err
}

func IterateConfig(ctx context.Context, nextTime int64, enabled bool, fn func(*Config) error) error {
	enabledInt := 0
	if enabled {
		enabledInt = 1
	}
	return xormutil.MustGetXormSession(ctx).
		Where("next_time <= ?", nextTime).
		And("enabled = ?", enabledInt).
		Iterate(new(Config), func(_ int, bean any) error {
			return fn(bean.(*Config))
		})
}
