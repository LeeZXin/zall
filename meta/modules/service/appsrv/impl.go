package appsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/prop/modules/service/propsrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
)

type outerImpl struct{}

func (*outerImpl) InsertApp(ctx context.Context, reqDTO InsertAppReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.AppSrvKeysVO.InsertApp),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验teamId
	if err = checkPermAdmin(ctx, reqDTO.Operator, reqDTO.TeamId); err != nil {
		return
	}
	_, b, err := appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.AlreadyExistsError()
		return
	}
	err = appmd.InsertApp(ctx, appmd.InsertAppReqDTO{
		AppId:  reqDTO.AppId,
		Name:   reqDTO.Name,
		TeamId: reqDTO.TeamId,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	propsrv.Inner.GrantAuth(ctx, reqDTO.AppId)
	return
}

func (*outerImpl) DeleteApp(ctx context.Context, reqDTO DeleteAppReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.AppSrvKeysVO.DeleteApp),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验teamId
	if err = checkPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return
	}
	_, err = appmd.DeleteByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListApp(ctx context.Context, reqDTO ListAppReqDTO) ([]AppDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	appIdList, err := checkPerm(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return nil, 0, err
	}
	var (
		apps []appmd.App
	)
	if len(appIdList) > 0 {
		apps, err = appmd.GetByAppIdList(ctx, appIdList)
	} else {
		apps, err = appmd.ListApp(ctx, appmd.ListAppReqDTO{
			AppId:  reqDTO.AppId,
			TeamId: reqDTO.TeamId,
			Cursor: reqDTO.Cursor,
			Limit:  reqDTO.Limit,
		})
	}
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret, _ := listutil.Map(apps, func(t appmd.App) (AppDTO, error) {
		return AppDTO{
			AppId: t.AppId,
			Name:  t.Name,
		}, nil
	})
	if reqDTO.Limit > 0 && len(apps) == reqDTO.Limit {
		return ret, apps[len(apps)-1].Id, nil
	}
	return ret, 0, nil
}

func (*outerImpl) UpdateApp(ctx context.Context, reqDTO UpdateAppReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.AppSrvKeysVO.TransferTeam),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验teamId
	if err = checkPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return
	}
	_, err = appmd.UpdateApp(ctx, appmd.UpdateAppReqDTO{
		AppId: reqDTO.AppId,
		Name:  reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// TransferTeam 迁移team
func (*outerImpl) TransferTeam(ctx context.Context, reqDTO TransferTeamReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.AppSrvKeysVO.TransferTeam),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才能迁移team
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, b, err := teammd.GetByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	app, b, err := appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 查询team角色是否包含该Id
	groups, err := teammd.ListTeamUserGroup(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 项目组仍有该仓库的特殊配置
	for _, group := range groups {
		for _, repoPerm := range group.GetPermDetail().AppPermList {
			if repoPerm.AppId == reqDTO.AppId {
				err = util.NewBizErr(apicode.OperationFailedErrCode, i18n.AppPermsContainerTargetAppId)
				return
			}
		}
	}
	_, err = appmd.TransferTeam(ctx, reqDTO.AppId, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func checkPerm(ctx context.Context, operator apisession.UserInfo, teamId int64) ([]string, error) {
	if operator.IsAdmin {
		return nil, nil
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, teamId, operator.Account)
	if !b {
		return nil, util.UnauthorizedError()
	}
	return listutil.Map(p.PermDetail.AppPermList, func(t perm.AppPermWithAppId) (string, error) {
		return t.AppId, nil
	})
}

func checkPermAdmin(ctx context.Context, operator apisession.UserInfo, teamId int64) error {
	if operator.IsAdmin {
		_, b, err := teammd.GetByTeamId(ctx, teamId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !b {
			return util.InvalidArgsError()
		}
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, teamId, operator.Account)
	if !b || !p.IsAdmin {
		return util.UnauthorizedError()
	}
	return nil
}

func checkPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, app.TeamId, operator.Account)
	if !b || !p.IsAdmin {
		return util.UnauthorizedError()
	}
	return nil
}
