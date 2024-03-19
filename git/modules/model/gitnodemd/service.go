package gitnodemd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func GetById(ctx context.Context, id int64) (Node, bool, error) {
	var ret Node
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func InsertNode(ctx context.Context, reqDTO InsertNodeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(Node{
		Name:     reqDTO.Name,
		HttpHost: reqDTO.HttpHost,
		SshHost:  reqDTO.SshHost,
	})
	return err
}

func DeleteNode(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(Node))
	return rows == 1, err
}

func UpdateNode(ctx context.Context, reqDTO UpdateNodeReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "http_host", "ssh_host").
		Limit(1).
		Update(Node{
			HttpHost: reqDTO.HttpHost,
			SshHost:  reqDTO.SshHost,
			Name:     reqDTO.Name,
		})
	return rows == 1, err
}

func GetAllNodes(ctx context.Context) ([]Node, error) {
	ret := make([]Node, 0)
	session := xormutil.MustGetXormSession(ctx)
	err := session.OrderBy("id asc").Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, err
}
