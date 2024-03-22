package detectmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func IsNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertInstance(ctx context.Context, instanceId string) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Instance{
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
		Where("heartbeat_time >= ?", after).
		OrderBy("id asc").
		Find(&ret)
	return ret, err
}

func DeleteInValidInstances(ctx context.Context, before int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("heartbeat_time < ?", before).
		Delete(new(Instance))
	return err
}

func IterateDetect(ctx context.Context, fn func(detect *TcpDetect) error) error {
	return xormutil.MustGetXormSession(ctx).
		Where("enabled = 1").
		Iterate(new(TcpDetect), func(_ int, bean interface{}) error {
			return fn(bean.(*TcpDetect))
		})
}

func UpdateDetectHeartbeatTime(ctx context.Context, id int64, heartbeatTime int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("heartbeat_time").
		Update(&TcpDetect{
			HeartbeatTime: heartbeatTime,
		})
	return err
}

func InsertLog(ctx context.Context, reqDTO InsertLogReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&DetectLog{
		DetectId: reqDTO.DetectId,
		Ip:       reqDTO.Ip,
		Port:     reqDTO.Port,
		Valid:    reqDTO.Valid,
	})
	return err
}

func DeleteLogByTime(ctx context.Context, before time.Time) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("created < ?", before.Format(time.DateTime)).
		Delete(new(DetectLog))
	return err
}

func InsertDetect(ctx context.Context, reqDTO InsertDetectReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&TcpDetect{
		Ip:            reqDTO.Ip,
		Port:          reqDTO.Port,
		Name:          reqDTO.Name,
		Enabled:       reqDTO.Enabled,
		HeartbeatTime: reqDTO.HeartbeatTime,
	})
	return err
}

func UpdateDetect(ctx context.Context, reqDTO UpdateDetectReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("ip", "port", "name").
		Update(&TcpDetect{
			Ip:   reqDTO.Ip,
			Port: reqDTO.Port,
			Name: reqDTO.Name,
		})
	return rows == 1, err
}

func DeleteDetect(ctx context.Context, id int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(TcpDetect))
	return err
}

func DeleteLog(ctx context.Context, detectId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("detect_id = ?", detectId).
		Delete(new(DetectLog))
	return err
}

func ListDetect(ctx context.Context, reqDTO ListDetectReqDTO) ([]TcpDetect, error) {
	session := xormutil.MustGetXormSession(ctx)
	if reqDTO.Name != "" {
		session.And("name like ?", reqDTO.Name+"%")
	}
	if reqDTO.Cursor > 0 {
		session.And("id > ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	ret := make([]TcpDetect, 0)
	err := session.OrderBy("id asc").Find(&ret)
	return ret, err
}

func ListLog(ctx context.Context, reqDTO ListLogReqDTO) ([]DetectLog, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("detect_id = ?", reqDTO.Id)
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	ret := make([]DetectLog, 0)
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}

func SetDetectEnabled(ctx context.Context, id int64, enabled bool) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("enabled").
		Update(&TcpDetect{
			Enabled: enabled,
		})
	return rows == 1, err
}

func GetDetectByIpPort(ctx context.Context, ip string, port int) (TcpDetect, bool, error) {
	var ret TcpDetect
	b, err := xormutil.MustGetXormSession(ctx).Where("ip = ?", ip).And("port = ?", port).Get(&ret)
	return ret, b, err
}
