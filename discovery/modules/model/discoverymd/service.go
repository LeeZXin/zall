package discoverymd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsDiscoverySourceNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func ListEtcdNode(ctx context.Context, reqDTO ListEtcdNodeReqDTO) ([]EtcdNode, error) {
	ret := make([]EtcdNode, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("env = ?", reqDTO.Env)
	if len(reqDTO.Cols) > 0 {
		session.Cols(reqDTO.Cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func InsertEtcdNode(ctx context.Context, reqDTO InsertEtcdNodeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&EtcdNode{
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

func BatchGetEtcdNodeByIdList(ctx context.Context, idList []int64, cols []string) ([]EtcdNode, error) {
	ret := make([]EtcdNode, 0)
	session := xormutil.MustGetXormSession(ctx).
		In("id", idList)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func BatchInsertAppEtcdNodeBind(ctx context.Context, reqDTOs []InsertAppEtcdNodeBindReqDTO) error {
	binds := listutil.MapNe(reqDTOs, func(t InsertAppEtcdNodeBindReqDTO) AppEtcdNodeBind {
		return AppEtcdNodeBind{
			NodeId: t.NodeId,
			AppId:  t.AppId,
			Env:    t.Env,
		}
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(binds)
	return err
}

func DeleteAppEtcdNodeBindByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(AppEtcdNodeBind))
	return err
}

func DeleteAppEtcdNodeBindByAppIdAndEnv(ctx context.Context, appId, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		Delete(new(AppEtcdNodeBind))
	return err
}

func BatchGetAppEtcdNodeBindByNodeIdListAndAppId(ctx context.Context, nodeIdList []int64, appId string) ([]AppEtcdNodeBind, error) {
	ret := make([]AppEtcdNodeBind, 0)
	err := xormutil.MustGetXormSession(ctx).
		And("app_id = ?", appId).
		In("node_id", nodeIdList).
		Find(&ret)
	return ret, err
}

func ListAppEtcdNodeBindByAppIdAndEnv(ctx context.Context, appId, env string) ([]AppEtcdNodeBind, error) {
	ret := make([]AppEtcdNodeBind, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		Find(&ret)
	return ret, err
}

func GetAppEtcdNodeBindById(ctx context.Context, id int64) (AppEtcdNodeBind, bool, error) {
	var ret AppEtcdNodeBind
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}
