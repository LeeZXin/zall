package tpfeishusrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/thirdpart/modules/model/tpfeishumd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.FeishuAccessTokenTopic, func(data any) {
			req, ok := data.(event.FeishuAccessTokenEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.FeishuAccessTokenCreateAction:
						return events.FeishuAccessToken.Create
					case event.FeishuAccessTokenUpdateAction:
						return events.FeishuAccessToken.Update
					case event.FeishuAccessTokenDeleteAction:
						return events.FeishuAccessToken.Delete
					case event.FeishuAccessTokenChangeApiKeyAction:
						return events.FeishuAccessToken.ChangeApiKey
					case event.FeishuAccessTokenRefreshAction:
						return events.FeishuAccessToken.Refresh
					default:
						return false
					}
				})
			}
		})
	})
}

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
	tokens, total, err := tpfeishumd.ListAccessToken(ctx, tpfeishumd.ListAccessTokenReqDTO{
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
		Key:      reqDTO.Key,
		TeamId:   reqDTO.TeamId,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret := listutil.MapNe(tokens, func(t tpfeishumd.AccessToken) AccessTokenDTO {
		return AccessTokenDTO{
			Id:          t.Id,
			TeamId:      t.TeamId,
			Name:        t.Name,
			AppId:       t.AppId,
			Creator:     t.Creator,
			Secret:      t.Secret,
			Token:       t.Token,
			TenantToken: t.TenantToken,
			ApiKey:      t.ApiKey,
			Expired:     time.UnixMilli(t.ExpireTime),
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
	team, err := checkManageAccessTokenPermByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	b, err := tpfeishumd.ExistAccessTokenByAppIdAndSecret(ctx, reqDTO.AppId, reqDTO.Secret)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	at, err := tpfeishumd.InsertAccessToken(ctx, tpfeishumd.InsertAccessTokenReqDTO{
		Name:    reqDTO.Name,
		TeamId:  reqDTO.TeamId,
		AppId:   reqDTO.AppId,
		Secret:  reqDTO.Secret,
		Creator: reqDTO.Operator.Account,
		ApiKey:  idutil.RandomUuid(),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	refreshAccessToken(ctx, &at)
	notifyEvent(team, at, reqDTO.Operator, event.FeishuAccessTokenCreateAction)
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
	at, team, err := checkManageAccessTokenPermByTokenId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	b, err := tpfeishumd.UpdateAccessToken(ctx, tpfeishumd.UpdateAccessTokenReqDTO{
		Id:     reqDTO.Id,
		Name:   reqDTO.Name,
		AppId:  reqDTO.AppId,
		Secret: reqDTO.Secret,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		// 如果不一致需要重刷token
		if at.AppId != reqDTO.AppId || at.Secret != reqDTO.Secret {
			at.AppId = reqDTO.AppId
			at.Secret = reqDTO.Secret
			at.Name = reqDTO.Name
			refreshAccessToken(ctx, &at)
		}
	}
	notifyEvent(team, at, reqDTO.Operator, event.FeishuAccessTokenUpdateAction)
	return nil
}

func DeleteAccessToken(ctx context.Context, reqDTO DeleteAccessTokenReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	at, team, err := checkManageAccessTokenPermByTokenId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = tpfeishumd.DeleteAccessTokenById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(team, at, reqDTO.Operator, event.FeishuAccessTokenDeleteAction)
	return nil
}

func RefreshAccessToken(ctx context.Context, reqDTO RefreshAccessTokenReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	at, team, err := checkManageAccessTokenPermByTokenId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = refreshAccessToken(ctx, &at)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(team, at, reqDTO.Operator, event.FeishuAccessTokenRefreshAction)
	return nil
}

func ChangeAccessTokenApiKey(ctx context.Context, reqDTO ChangeAccessTokenApiKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	at, team, err := checkManageAccessTokenPermByTokenId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = tpfeishumd.UpdateAccessTokenApiKeyById(ctx, reqDTO.Id, idutil.RandomUuid())
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(team, at, reqDTO.Operator, event.FeishuAccessTokenChangeApiKeyAction)
	return nil
}

func GetAccessTokenByApiKey(ctx context.Context, apiKey string) (string, string, error) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	at, b, err := tpfeishumd.GetAccessTokenByApiKey(ctx, apiKey)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", "", util.InternalError(err)
	}
	if !b {
		return "", "", util.InvalidArgsError()
	}
	if !at.IsNotExpired() {
		return "", "", util.OperationFailedError()
	}
	return at.Token, at.TenantToken, nil
}

func checkManageAccessTokenPermByTokenId(ctx context.Context, tokenId int64, operator apisession.UserInfo) (tpfeishumd.AccessToken, teammd.Team, error) {
	at, b, err := tpfeishumd.GetAccessTokenById(ctx, tokenId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return tpfeishumd.AccessToken{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return tpfeishumd.AccessToken{}, teammd.Team{}, util.InvalidArgsError()
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
	if p.IsAdmin || p.PermDetail.TeamPerm.CanManageFeishuAccessToken {
		return team, nil
	}
	return team, util.UnauthorizedError()
}

func notifyEvent(team teammd.Team, at tpfeishumd.AccessToken, operator apisession.UserInfo, action event.FeishuAccessTokenEventAction) {
	initPsub()
	psub.Publish(event.FeishuAccessTokenTopic, event.FeishuAccessTokenEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action: action,
		Name:   at.Name,
		AppId:  at.AppId,
	})
}
