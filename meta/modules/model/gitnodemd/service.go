package gitnodemd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"strings"
)

func IsNodeIdValid(nodeId string) bool {
	return len(nodeId) > 0 && len(nodeId) <= 32
}

func InsertGitNode(ctx context.Context, reqDTO InsertNodeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&GitNode{
		NodeId:   reqDTO.NodeId,
		HttpHost: strings.Join(reqDTO.HttpHosts, ","),
		SshHost:  strings.Join(reqDTO.SshHosts, ","),
		Version:  0,
	})
	return err
}

func GetByNodeId(ctx context.Context, nodeId string) (GitNodeDTO, bool, error) {
	var ret GitNode
	b, err := xormutil.MustGetXormSession(ctx).Where("node_id = ?", nodeId).Get(&ret)
	return md2dto(ret), b, err
}

func DeleteNode(ctx context.Context, nodeId string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("node_id = ?", nodeId).
		Limit(1).
		Delete(new(GitNode))
	return rows == 1, err
}

func UpdateNode(ctx context.Context, reqDTO UpdateNodeReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("node_id = ?", reqDTO.NodeId).
		Cols("version", "http_host", "ssh_host").
		Incr("version").
		Limit(1).
		Update(GitNode{
			HttpHost: strings.Join(reqDTO.HttpHosts, ","),
			SshHost:  strings.Join(reqDTO.SshHosts, ","),
		})
	return rows == 1, err
}

func GetAll(ctx context.Context) ([]GitNodeDTO, error) {
	ret := make([]GitNode, 0)
	session := xormutil.MustGetXormSession(ctx)
	err := session.OrderBy("id asc").Find(&ret)
	if err != nil {
		return nil, err
	}
	return listutil.Map(ret, func(t GitNode) (GitNodeDTO, error) {
		return md2dto(t), nil
	})
}

func md2dto(node GitNode) GitNodeDTO {
	return GitNodeDTO{
		NodeId:    node.NodeId,
		HttpHosts: strings.Split(node.HttpHost, ","),
		SshHosts:  strings.Split(node.SshHost, ","),
		Version:   node.Version,
	}
}
