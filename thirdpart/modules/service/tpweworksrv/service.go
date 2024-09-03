package tpweworksrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/thirdpart/modules/model/tpweworkmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

// ListAccessToken token列表
func ListAccessToken(ctx context.Context, reqDTO ListAccessTokenReqDTO) ([]AccessTokenDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkManageAccessTokenPermByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return nil, 0, err
	}
	tokens, total, err := tpweworkmd.ListAccessToken(ctx, tpweworkmd.ListAccessTokenReqDTO{
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
		Key:      reqDTO.Key,
		TeamId:   reqDTO.TeamId,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret := listutil.MapNe(tokens, func(t tpweworkmd.AccessToken) AccessTokenDTO {
		return AccessTokenDTO{
			Id:      t.Id,
			TeamId:  t.TeamId,
			Name:    t.Name,
			CorpId:  t.CorpId,
			Creator: t.Creator,
			Secret:  t.Secret,
			Token:   t.Token,
			ApiKey:  t.ApiKey,
			Expired: time.UnixMilli(t.ExpireTime),
		}
	})
	return ret, total, nil
}

// CreateAccessToken 创建token
func CreateAccessToken(ctx context.Context, reqDTO CreateAccessTokenReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkManageAccessTokenPermByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	b, err := tpweworkmd.ExistAccessTokenByCorpIdAndSecret(ctx, reqDTO.CorpId, reqDTO.Secret)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	at, err := tpweworkmd.InsertAccessToken(ctx, tpweworkmd.InsertAccessTokenReqDTO{
		Name:    reqDTO.Name,
		TeamId:  reqDTO.TeamId,
		CorpId:  reqDTO.CorpId,
		Secret:  reqDTO.Secret,
		Creator: reqDTO.Operator.Account,
		ApiKey:  idutil.RandomUuid(),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	refreshAccessToken(ctx, &at)
	return nil
}

// UpdateAccessToken 编辑token
func UpdateAccessToken(ctx context.Context, reqDTO UpdateAccessTokenReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	at, _, err := checkManageAccessTokenPermByTokenId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	b, err := tpweworkmd.UpdateAccessToken(ctx, tpweworkmd.UpdateAccessTokenReqDTO{
		Id:     reqDTO.Id,
		Name:   reqDTO.Name,
		CorpId: reqDTO.CorpId,
		Secret: reqDTO.Secret,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		// 如果不一致需要重刷token
		if at.CorpId != reqDTO.CorpId || at.Secret != reqDTO.Secret {
			at.CorpId = reqDTO.CorpId
			at.Secret = reqDTO.Secret
			at.Name = reqDTO.Name
			refreshAccessToken(ctx, &at)
		}
	}
	return nil
}

func DeleteAccessToken(ctx context.Context, reqDTO DeleteAccessTokenReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, _, err := checkManageAccessTokenPermByTokenId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = tpweworkmd.DeleteAccessTokenById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func RefreshAccessToken(ctx context.Context, reqDTO RefreshAccessTokenReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	at, _, err := checkManageAccessTokenPermByTokenId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = refreshAccessToken(ctx, &at)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func ChangeAccessTokenApiKey(ctx context.Context, reqDTO ChangeAccessTokenApiKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, _, err := checkManageAccessTokenPermByTokenId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = tpweworkmd.UpdateAccessTokenApiKeyById(ctx, reqDTO.Id, idutil.RandomUuid())
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func GetAccessTokenByApiKey(ctx context.Context, apiKey string) (string, error) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	at, b, err := tpweworkmd.GetAccessTokenByApiKey(ctx, apiKey)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", util.InvalidArgsError()
	}
	if !at.IsNotExpired() {
		return "", util.OperationFailedError()
	}
	return at.Token, nil
}

func checkManageAccessTokenPermByTokenId(ctx context.Context, tokenId int64, operator apisession.UserInfo) (tpweworkmd.AccessToken, teammd.Team, error) {
	at, b, err := tpweworkmd.GetAccessTokenById(ctx, tokenId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return tpweworkmd.AccessToken{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return tpweworkmd.AccessToken{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, err := checkManageAccessTokenPermByTeamId(ctx, at.TeamId, operator)
	return at, team, err
}

func checkManageAccessTokenPermByTeamId(ctx context.Context, teamId int64, operator apisession.UserInfo) (teammd.Team, error) {
	team, b, err := teammd.GetByTeamId(ctx, teamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return teammd.Team{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return team, util.InternalError(err)
	}
	if !b {
		return team, util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanManageWeworkAccessToken {
		return team, nil
	}
	return team, util.UnauthorizedError()
}
