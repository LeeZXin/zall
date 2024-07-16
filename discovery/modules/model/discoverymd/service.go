package discoverymd

import (
	"context"
	"github.com/LeeZXin/zsf/services/lb"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsDiscoverySourceNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func ListEtcdNode(ctx context.Context, reqDTO ListEtcdNodeReqDTO) ([]EtcdNode, error) {
	ret := make([]EtcdNode, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("env = ?", reqDTO.Env)
	if len(reqDTO.Cols) > 0 {
		session.Cols(reqDTO.Cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func InsertEtcdNode(ctx context.Context, reqDTO InsertEtcdNodeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&EtcdNode{
			AppId:     reqDTO.AppId,
			Name:      reqDTO.Name,
			Env:       reqDTO.Env,
			Endpoints: reqDTO.Endpoints,
			Username:  reqDTO.Username,
			Password:  reqDTO.Password,
		})
	return err
}

func DeleteEtcdNodeById(ctx context.Context, id int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(EtcdNode))
	return err
}

func DeleteEtcdNodeByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(EtcdNode))
	return err
}

func UpdateEtcdNode(ctx context.Context, reqDTO UpdateEtcdNodeReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("endpoints", "username", "password", "name").
		Update(&EtcdNode{
			Endpoints: reqDTO.Endpoints,
			Username:  reqDTO.Username,
			Password:  reqDTO.Password,
			Name:      reqDTO.Name,
		})
	return rows == 1, err
}

func GetEtcdNodeById(ctx context.Context, id int64) (EtcdNode, bool, error) {
	var ret EtcdNode
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func InsertDownService(ctx context.Context, reqDTO InsertDownServiceDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&DownService{
			AppId:    reqDTO.AppId,
			SourceId: reqDTO.SourceId,
			DownService: &xormutil.Conversion[lb.Server]{
				Data: reqDTO.Service,
			},
			InstanceId: reqDTO.InstanceId,
		})
	return err
}

func DeleteDownServiceById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(DownService))
	return rows == 1, err
}

func DeleteDownServiceByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(DownService))
	return err
}

func ListDownServiceBySourceId(ctx context.Context, sourceId int64) ([]DownService, error) {
	ret := make([]DownService, 0)
	err := xormutil.MustGetXormSession(ctx).Where("source_id = ?", sourceId).Find(&ret)
	return ret, err
}

func GetDownServiceBySourceIdAndInstanceId(ctx context.Context, sourceId int64, instanceId string) (DownService, bool, error) {
	var ret DownService
	b, err := xormutil.MustGetXormSession(ctx).
		Where("source_id = ?", sourceId).
		And("instance_id = ?", instanceId).
		Get(&ret)
	return ret, b, err
}
