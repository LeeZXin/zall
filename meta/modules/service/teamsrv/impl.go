package teamsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
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
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/patrickmn/go-cache"
	"strconv"
	"time"
)

type innerImpl struct {
	permCache *cache.Cache
}

func (s *innerImpl) GetUserPermDetail(ctx context.Context, teamId int64, account string) (teammd.UserPermDetailDTO, bool) {
	key := strconv.FormatInt(teamId, 10) + "_" + account
	v, b := s.permCache.Get(key)
	if b {
		r := v.(teammd.UserPermDetailDTO)
		return r, r.RoleId != 0
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	r, b, err := teammd.GetUserPermDetail(ctx, teamId, account)
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

// CreateTeam 创建团队
func (*outerImpl) CreateTeam(ctx context.Context, reqDTO CreateTeamReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TeamSrvKeysVO.CreateTeam),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		// 判断是否允许用户自行创建项目
		sysCfg, b := cfgsrv.Inner.GetSysCfg(ctx)
		if !b || !sysCfg.AllowUserCreateTeam {
			err = util.UnauthorizedError()
			return
		}
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 创建项目
		pu, err := teammd.InsertTeam(ctx, teammd.InsertTeamReqDTO{
			Name: reqDTO.Name,
		})
		if err != nil {
			return err
		}
		// 创建管理员组
		group, err := teammd.InsertRole(ctx, teammd.InsertRoleReqDTO{
			Name:       i18n.GetByKey(i18n.TeamAdminUserGroupName),
			TeamId:     pu.Id,
			PermDetail: perm.DefaultPermDetail,
			IsAdmin:    true,
		})
		if err != nil {
			return err
		}
		// 创建关联关系
		return teammd.InsertUser(ctx, teammd.InsertUserReqDTO{
			TeamId:  pu.Id,
			Account: reqDTO.Operator.Account,
			RoleId:  group.Id,
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
	ctx, closer := xormstore.Context(ctx)
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

// IsAdmin 是否是团队管理员
func (*outerImpl) IsAdmin(ctx context.Context, reqDTO IsAdminReqDTO) (bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	detail, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return false, err
	}
	if !b {
		return false, nil
	}
	return detail.IsAdmin, nil
}

// GetTeamPerm 获取团队权限
func (*outerImpl) GetTeamPerm(ctx context.Context, reqDTO GetTeamPermReqDTO) (perm.TeamPerm, error) {
	if err := reqDTO.IsValid(); err != nil {
		return perm.TeamPerm{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	detail, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return perm.TeamPerm{}, err
	}
	if !b {
		return perm.TeamPerm{}, nil
	}
	return detail.PermDetail.TeamPerm, nil
}

// GetTeam 获取团队信息
func (*outerImpl) GetTeam(ctx context.Context, reqDTO GetTeamReqDTO) (TeamDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TeamDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TeamDTO{}, err
	}
	if !b {
		return TeamDTO{}, nil
	}
	team, b, err := teammd.GetByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TeamDTO{}, err
	}
	if !b {
		return TeamDTO{}, nil
	}
	return TeamDTO{
		TeamId: team.Id,
		Name:   team.Name,
	}, nil
}

func (*outerImpl) ListUser(ctx context.Context, reqDTO ListUserReqDTO) ([]UserDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	users, err := teammd.ListUser(ctx, teammd.ListUserReqDTO{
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
	groupIdList, _ := listutil.Map(users, func(t teammd.User) (int64, error) {
		return t.RoleId, nil
	})
	groupIdList = listutil.Distinct(groupIdList...)
	groups, err := teammd.GetByRoleIdList(ctx, groupIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	groupIdNameMap, _ := listutil.CollectToMap(groups, func(t teammd.Role) (int64, error) {
		return t.Id, nil
	}, func(t teammd.Role) (string, error) {
		return t.Name, nil
	})
	ret, _ := listutil.Map(users, func(t teammd.User) (UserDTO, error) {
		return UserDTO{
			TeamId:   t.TeamId,
			Account:  t.Account,
			RoleId:   t.RoleId,
			RoleName: groupIdNameMap[t.RoleId],
			Created:  t.Created,
		}, nil
	})
	return ret, next, nil
}

func (*outerImpl) DeleteUser(ctx context.Context, reqDTO DeleteUserReqDTO) (err error) {
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
	ctx, closer := xormstore.Context(ctx)
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
	_, err = teammd.DeleteUser(ctx, reqDTO.TeamId, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) UpsertUser(ctx context.Context, reqDTO UpsertUserReqDTO) (err error) {
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
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return
	}
	// 校验groupId是否正确
	group, b, err := teammd.GetByRoleId(ctx, reqDTO.RoleId)
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
		err = teammd.InsertUser(ctx, teammd.InsertUserReqDTO{
			TeamId:  reqDTO.TeamId,
			Account: reqDTO.Account,
			RoleId:  reqDTO.RoleId,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
	} else {
		_, err = teammd.UpdateUser(ctx, teammd.UpdateUserReqDTO{
			TeamId:  reqDTO.TeamId,
			Account: reqDTO.Account,
			RoleId:  reqDTO.RoleId,
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
	_, b, err := teammd.GetByTeamId(ctx, teamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 如果是企业管理员
	if operator.IsAdmin {
		return nil
	}
	// 判断权限
	pu, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
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

func checkTeamUserPermByGroupId(ctx context.Context, operator apisession.UserInfo, groupId int64) (teammd.Role, error) {
	// 检查权限
	group, b, err := teammd.GetByRoleId(ctx, groupId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Role{}, util.InternalError(err)
	}
	if !b {
		return teammd.Role{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return group, nil
	}
	// 检查项目管理员权限
	pu, b, err := teammd.GetUserPermDetail(ctx, group.TeamId, operator.Account)
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

func (*outerImpl) InsertRole(ctx context.Context, reqDTO InsertRoleReqDTO) (err error) {
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
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err = checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return
	}
	if _, err = teammd.InsertRole(ctx, teammd.InsertRoleReqDTO{
		Name:       reqDTO.Name,
		TeamId:     reqDTO.TeamId,
		PermDetail: reqDTO.Perm,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) UpdateRoleName(ctx context.Context, reqDTO UpdateRoleNameReqDTO) (err error) {
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
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	_, err = checkTeamUserPermByGroupId(ctx, reqDTO.Operator, reqDTO.RoleId)
	if err != nil {
		return
	}
	_, err = teammd.UpdateRoleName(ctx, reqDTO.RoleId, reqDTO.Name)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) UpdateRolePerm(ctx context.Context, reqDTO UpdateRolePermReqDTO) (err error) {
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
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	group, err := checkTeamUserPermByGroupId(ctx, reqDTO.Operator, reqDTO.RoleId)
	if err != nil {
		return
	}
	// 管理员项目组无法编辑权限
	if group.IsAdmin {
		err = util.NewBizErr(apicode.CannotUpdateTeamUserAdminGroupCode, i18n.TeamUserGroupUpdateAdminNotAllow)
		return
	}
	// 检查仓库Id配置
	if len(reqDTO.Perm.RepoPermList) > 0 {
		idList, _ := listutil.Map(reqDTO.Perm.RepoPermList, func(t perm.RepoPermWithId) (int64, error) {
			return t.RepoId, nil
		})
		var repos []repomd.Repo
		repos, err = repomd.GetRepoByIdList(ctx, idList)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		for _, repo := range repos {
			// 仓库不属于项目组里面的
			if repo.TeamId != group.TeamId {
				return util.InvalidArgsError()
			}
		}
	}
	if _, err = teammd.UpdateRolePerm(ctx, reqDTO.RoleId, reqDTO.Perm); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return nil
}

func (*outerImpl) DeleteRole(ctx context.Context, reqDTO DeleteRoleReqDTO) (err error) {
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
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	group, err := checkTeamUserPermByGroupId(ctx, reqDTO.Operator, reqDTO.RoleId)
	if err != nil {
		return
	}
	b, err := teammd.ExistRole(ctx, group.TeamId, reqDTO.RoleId)
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
	if _, err = teammd.DeleteRole(ctx, reqDTO.RoleId); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return nil
}

func (*outerImpl) ListRole(ctx context.Context, reqDTO ListRoleReqDTO) ([]RoleDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return nil, err
	}
	groups, err := teammd.ListRole(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(groups, func(t teammd.Role) (RoleDTO, error) {
		perm := t.Perm
		return RoleDTO{
			RoleId: t.Id,
			TeamId: t.TeamId,
			Name:   t.Name,
			Perm:   *perm,
		}, nil
	})
}

// ListTeam 展示用户所在团队列表
func (*outerImpl) ListTeam(ctx context.Context, reqDTO ListTeamReqDTO) ([]TeamDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	puList, err := teammd.ListUserByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	teamIdList, _ := listutil.Map(puList, func(t teammd.User) (int64, error) {
		return t.TeamId, nil
	})
	teamList, err := teammd.GetByTeamIdList(ctx, teamIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(teamList, func(t teammd.Team) (TeamDTO, error) {
		return TeamDTO{
			TeamId: t.Id,
			Name:   t.Name,
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
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err = checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return err
	}
	// 检查是否还有挂在该项目组的仓库
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
	// 检查是否有挂在该项目组的app
	appCount, err := appmd.CountByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 存在 不允许删除
	if appCount > 0 {
		err = util.NewBizErr(apicode.DataAlreadyExistsCode, i18n.AppRemainCountGreaterThanZero)
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 项目表
		_, err := teammd.DeleteTeam(ctx, reqDTO.TeamId)
		if err != nil {
			return err
		}
		// 项目用户组
		_, err = teammd.DeleteAllRoleByTeamId(ctx, reqDTO.TeamId)
		if err != nil {
			return err
		}
		// 项目用户
		_, err = teammd.DeleteAllUserByTeamId(ctx, reqDTO.TeamId)
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
