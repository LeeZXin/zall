package zalletmd

import (
	"context"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
	"xorm.io/builder"
)

func IsZalletNodeIdValid(nodeId string) bool {
	return regexp.MustCompile("[\\w-]{1,32}").MatchString(nodeId)
}

func IsZalletNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsZalletAgentHostValid(agentHost string) bool {
	return util.GenIpPortPattern().MatchString(agentHost)
}

func IsZalletAgentTokenValid(agentToken string) bool {
	return len(agentToken) <= 1024
}

func InsertZalletNode(ctx context.Context, reqDTO InsertZalletNodeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&ZalletNode{
			NodeId:     reqDTO.NodeId,
			Name:       reqDTO.Name,
			AgentHost:  reqDTO.AgentHost,
			AgentToken: reqDTO.AgentToken,
		})
	return err
}

func GetZalletNodeByNodeId(ctx context.Context, nodeId string) (ZalletNode, bool, error) {
	var ret ZalletNode
	b, err := xormutil.MustGetXormSession(ctx).Where("node_id = ?", nodeId).Get(&ret)
	return ret, b, err
}

func ExistZalletNodeByNodeId(ctx context.Context, nodeId string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).Where("node_id = ?", nodeId).Exist(new(ZalletNode))
}

func ListZalletNode(ctx context.Context, reqDTO ListZalletNodeReqDTO) ([]ZalletNode, int64, error) {
	session := xormutil.MustGetXormSession(ctx)
	if reqDTO.Name != "" {
		session.And(
			builder.Or(
				builder.Expr("name like ?", reqDTO.Name+"%"),
				builder.Expr("node_id like ?", reqDTO.Name+"%"),
			))
	}
	if len(reqDTO.Cols) > 0 {
		session.Cols(reqDTO.Cols...)
	}
	ret := make([]ZalletNode, 0)
	total, err := session.Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).FindAndCount(&ret)
	return ret, total, err
}

func ListAllZalletNode(ctx context.Context, cols []string) ([]ZalletNode, error) {
	session := xormutil.MustGetXormSession(ctx)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	ret := make([]ZalletNode, 0)
	err := session.Find(&ret)
	return ret, err
}

func DeleteZalletNodeById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(ZalletNode))
	return rows == 1, err
}

func UpdateZalletNode(ctx context.Context, reqDTO UpdateZalletNodeReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "agent_host", "agent_token").
		Update(&ZalletNode{
			Name:       reqDTO.Name,
			AgentHost:  reqDTO.AgentHost,
			AgentToken: reqDTO.AgentToken,
		})
	return rows == 1, err
}
