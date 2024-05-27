package oplogmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func BatchInsertLog(ctx context.Context, reqDTOList []InsertOpLogReqDTO) error {
	opList, _ := listutil.Map(reqDTOList, func(t InsertOpLogReqDTO) (*OpLog, error) {
		return &OpLog{
			RepoId:   t.RepoId,
			Operator: t.Operator,
			Content:  t.Content,
			ReqBody:  t.ReqBody,
			Created:  t.Created,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(opList)
	return err
}

func PageLog(ctx context.Context, reqDTO PageOpLogReqDTO) ([]OpLog, int64, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", reqDTO.RepoId).
		And("created between ? and ?", reqDTO.BeginTime.Format(time.DateTime), reqDTO.EndTime.Format(time.DateTime))
	if reqDTO.Account != "" {
		session.And("operator like ?", reqDTO.Account+"%")
	}
	ret := make([]OpLog, 0)
	total, err := session.
		Desc("id").
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		FindAndCount(&ret)
	return ret, total, err
}
