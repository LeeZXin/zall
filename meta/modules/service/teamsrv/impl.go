package teamsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"github.com/patrickmn/go-cache"
	"strconv"
	"time"
)

type innerImpl struct {
	permCache *cache.Cache
}

func (s *innerImpl) GetTeamUserPermDetail(ctx context.Context, teamId int64, account string) (teammd.TeamUserPermDetailDTO, bool) {
	key := strconv.FormatInt(teamId, 10) + "_" + account
	v, b := s.permCache.Get(key)
	if b {
		r := v.(teammd.TeamUserPermDetailDTO)
		return r, r.GroupId != 0
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	r, b, err := teammd.GetTeamUserPermDetail(ctx, teamId, account)
	if err != nil || !b {
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		s.permCache.Set(key, r, time.Second)
	} else {
		s.permCache.Set(key, r, time.Minute)
	}
	return r, b
}

type outerImpl struct{}

// InsertTeam 创建项目
func (*outerImpl) InsertTeam(ctx context.Context, reqDTO InsertTeamReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.InsertTeam),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		// 判断是否允许用户自行创建项目
		sysCfg, b := cfgsrv.Inner.GetSysCfg(ctx)
		if !b || !sysCfg.AllowUserCreateTeam {
			err = util.UnauthorizedError()
			return
		}
	}
	err = mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		// 创建项目
		pu, err := teammd.InsertTeam(ctx, teammd.InsertTeamReqDTO{
			Name: reqDTO.Name,
		})
		if err != nil {
			return err
		}
		// 创建管理员组
		group, err := teammd.InsertTeamUserGroup(ctx, teammd.InsertTeamUserGroupReqDTO{
			Name:       i18n.GetByKey(i18n.TeamAdminUserGroupName),
			TeamId:     pu.Id,
			PermDetail: perm.DefaultPermDetail,
			IsAdmin:    true,
		})
		if err != nil {
			return err
		}
		// 创建关联关系
		return teammd.InsertTeamUser(ctx, teammd.InsertTeamUserReqDTO{
			TeamId:  pu.Id,
			Account: reqDTO.Operator.Account,
			GroupId: group.Id,
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) UpdateTeam(ctx context.Context, reqDTO UpdateTeamReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.UpdateTeam),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err = checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return err
	}
	_, err = teammd.UpdateTeam(ctx, teammd.UpdateTeamReqDTO{
		TeamId: reqDTO.TeamId,
		Name:   reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) ListTeamUser(ctx context.Context, reqDTO ListTeamUserReqDTO) ([]TeamUserDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err := checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	users, err := teammd.ListTeamUser(ctx, teammd.ListTeamUserReqDTO{
		TeamId:  reqDTO.TeamId,
		Account: reqDTO.Account,
		Cursor:  reqDTO.Cursor,
		Limit:   reqDTO.Limit,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(users) == reqDTO.Limit {
		next = users[len(users)-1].Id
	}
	groupIdList, _ := listutil.Map(users, func(t teammd.TeamUser) (int64, error) {
		return t.GroupId, nil
	})
	groupIdList = listutil.Distinct(groupIdList...)
	groups, err := teammd.GetByGroupIdList(ctx, reqDTO.TeamId, groupIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	groupIdNameMap, _ := listutil.CollectToMap(groups, func(t teammd.TeamUserGroup) (int64, error) {
		return t.Id, nil
	}, func(t teammd.TeamUserGroup) (string, error) {
		return t.Name, nil
	})
	ret, _ := listutil.Map(users, func(t teammd.TeamUser) (TeamUserDTO, error) {
		return TeamUserDTO{
			TeamId:    t.TeamId,
			Account:   t.Account,
			GroupId:   t.GroupId,
			GroupName: groupIdNameMap[t.GroupId],
			Created:   t.Created,
		}, nil
	})
	return ret, next, nil
}

func (*outerImpl) DeleteTeamUser(ctx context.Context, reqDTO DeleteTeamUserReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.DeleteTeamUser),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err = checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return
	}
	_, b, err := teammd.GetTeamUser(ctx, reqDTO.TeamId, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		return util.InvalidArgsError()
	}
	_, err = teammd.DeleteTeamUser(ctx, reqDTO.TeamId, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) UpsertTeamUser(ctx context.Context, reqDTO UpsertTeamUserReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.UpsertTeamUser),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err = checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return
	}
	// 校验groupId是否正确
	group, b, err := teammd.GetByGroupId(ctx, reqDTO.GroupId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b || group.TeamId != reqDTO.TeamId {
		err = util.InvalidArgsError()
		return
	}
	// 校验账号是否存在
	_, b, err = usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	_, b, err = teammd.GetTeamUser(ctx, reqDTO.TeamId, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		// 不存在则插入
		err = teammd.InsertTeamUser(ctx, teammd.InsertTeamUserReqDTO{
			TeamId:  reqDTO.TeamId,
			Account: reqDTO.Account,
			GroupId: reqDTO.GroupId,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
	} else {
		_, err = teammd.UpdateTeamUser(ctx, teammd.UpdateTeamUserReqDTO{
			TeamId:  reqDTO.TeamId,
			Account: reqDTO.Account,
			GroupId: reqDTO.GroupId,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
	}
	return
}

func checkPerm(ctx context.Context, teamId int64, operator apisession.UserInfo) error {
	// 如果是企业管理员
	if operator.IsAdmin {
		return nil
	}
	_, b, err := teammd.GetByTeamId(ctx, teamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 判断权限
	pu, b, err := teammd.GetTeamUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 不存在或不是管理员角色
	if !b || !pu.IsAdmin {
		return util.UnauthorizedError()
	}
	return nil
}

func checkTeamUserPermByGroupId(ctx context.Context, operator apisession.UserInfo, groupId int64) (teammd.TeamUserGroup, error) {
	// 检查权限
	group, b, err := teammd.GetByGroupId(ctx, groupId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.TeamUserGroup{}, util.InternalError(err)
	}
	if !b {
		return teammd.TeamUserGroup{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return group, nil
	}
	// 检查项目管理员权限
	pu, b, err := teammd.GetTeamUserPermDetail(ctx, group.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return group, util.InternalError(err)
	}
	// 不存在或不是管理员角色
	if !b || !pu.IsAdmin {
		return group, util.UnauthorizedError()
	}
	return group, nil
}

func (*outerImpl) InsertTeamUserGroup(ctx context.Context, reqDTO InsertTeamUserGroupReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.InsertTeamUserGroup),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err = checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return
	}
	if _, err = teammd.InsertTeamUserGroup(ctx, teammd.InsertTeamUserGroupReqDTO{
		Name:       reqDTO.Name,
		TeamId:     reqDTO.TeamId,
		PermDetail: reqDTO.Perm,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) UpdateTeamUserGroupName(ctx context.Context, reqDTO UpdateTeamUserGroupNameReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.UpdateTeamUserGroupName),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	_, err = checkTeamUserPermByGroupId(ctx, reqDTO.Operator, reqDTO.GroupId)
	if err != nil {
		return
	}
	_, err = teammd.UpdateTeamUserGroupName(ctx, reqDTO.GroupId, reqDTO.Name)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) UpdateTeamUserGroupPerm(ctx context.Context, reqDTO UpdateTeamUserGroupPermReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.UpdateTeamUserGroupPerm),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	group, err := checkTeamUserPermByGroupId(ctx, reqDTO.Operator, reqDTO.GroupId)
	if err != nil {
		return
	}
	// 管理员项目组无法编辑权限
	if group.IsAdmin {
		err = util.NewBizErr(apicode.CannotUpdateTeamUserAdminGroupCode, i18n.TeamUserGroupUpdateAdminNotAllow)
		return
	}
	if _, err = teammd.UpdateTeamUserGroupPerm(ctx, reqDTO.GroupId, reqDTO.Perm); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return nil
}

func (*outerImpl) DeleteTeamUserGroup(ctx context.Context, reqDTO DeleteTeamUserGroupReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.DeleteTeamUserGroup),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	group, err := checkTeamUserPermByGroupId(ctx, reqDTO.Operator, reqDTO.GroupId)
	if err != nil {
		return
	}
	b, err := teammd.ExistTeamUser(ctx, group.TeamId, reqDTO.GroupId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 存在属于该groupId的用户
	if b {
		err = util.NewBizErr(apicode.TeamUserGroupHasUserWhenDelCode, i18n.TeamUserGroupHasUserWhenDel)
		return
	}
	if _, err = teammd.DeleteTeamUserGroup(ctx, reqDTO.GroupId); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return nil
}

func (*outerImpl) ListTeamUserGroup(ctx context.Context, reqDTO ListTeamUserGroupReqDTO) ([]TeamUserGroupDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err := checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return nil, err
	}
	groups, err := teammd.ListTeamUserGroup(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(groups, func(t teammd.TeamUserGroup) (TeamUserGroupDTO, error) {
		return TeamUserGroupDTO{
			GroupId: t.Id,
			TeamId:  t.TeamId,
			Name:    t.Name,
			Perm:    t.GetPermDetail(),
		}, nil
	})
}

func (*outerImpl) ListTeam(ctx context.Context, reqDTO ListTeamReqDTO) ([]TeamDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	puList, err := teammd.ListTeamUserByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	teamIdList, _ := listutil.Map(puList, func(t teammd.TeamUser) (int64, error) {
		return t.TeamId, nil
	})
	teamList, err := teammd.GetTeamByTeamIdList(ctx, teamIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(teamList, func(t teammd.Team) (TeamDTO, error) {
		return TeamDTO{
			TeamId:  t.Id,
			Name:    t.Name,
			Created: t.Created,
		}, nil
	})
}

// DeleteTeam 删除项目
func (*outerImpl) DeleteTeam(ctx context.Context, reqDTO DeleteTeamReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.DeleteTeam),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err = checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return err
	}
	// 检查是否还有挂在该项目的仓库
	repoCount, err := repomd.CountByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 存在 不允许删除
	if repoCount > 0 {
		err = util.NewBizErr(apicode.DataAlreadyExistsCode, i18n.RepoRemainCountGreaterThanZero)
		return
	}
	err = mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		// 项目表
		_, err := teammd.DeleteTeam(ctx, reqDTO.TeamId)
		if err != nil {
			return err
		}
		// 项目用户组
		_, err = teammd.DeleteAllTeamUserGroupByTeamId(ctx, reqDTO.TeamId)
		if err != nil {
			return err
		}
		// 项目用户
		_, err = teammd.DeleteAllTeamUserByTeamId(ctx, reqDTO.TeamId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return nil
}
