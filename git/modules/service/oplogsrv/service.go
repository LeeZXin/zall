package oplogsrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/oplogmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

func InsertOpLog(t OpLog) {
	reqBody, _ := json.Marshal(t)
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	err := oplogmd.BatchInsertLog(ctx, []oplogmd.InsertOpLogReqDTO{
		{
			RepoId:   t.RepoId,
			Operator: t.Operator,
			Content:  t.Log,
			Created:  t.EventTime,
			ReqBody:  string(reqBody),
		},
	})
	if err != nil {
		logger.Logger.Error(err)
	}
}

func PageOpLog(ctx context.Context, reqDTO PageOpLogReqDTO) ([]OpLogDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkTeamAdmin(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, 0, err
	}
	const pageSize = 10
	d := reqDTO.DateTime
	logs, total, err := oplogmd.PageLog(ctx, oplogmd.PageOpLogReqDTO{
		RepoId:    reqDTO.RepoId,
		PageNum:   reqDTO.PageNum,
		PageSize:  pageSize,
		Account:   reqDTO.Account,
		BeginTime: time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()),
		EndTime:   time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location()),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, err
	}
	data, _ := listutil.Map(logs, func(t oplogmd.OpLog) (OpLogDTO, error) {
		return OpLogDTO{
			Id:      t.Id,
			Account: t.Operator,
			Created: t.Created,
			Content: t.Content,
			ReqBody: t.ReqBody,
		}, nil
	})
	return data, total, nil
}

func FormatI18n(key i18n.KeyItem, str ...any) string {
	return fmt.Sprintf(i18n.GetByKey(key), str...)
}

func checkTeamAdmin(ctx context.Context, repoId int64, operator apisession.UserInfo) error {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	return util.UnauthorizedError()
}
