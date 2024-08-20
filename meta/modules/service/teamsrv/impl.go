package teamsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/notify/modules/model/notifymd"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/teamhook/modules/model/teamhookmd"
	"github.com/LeeZXin/zall/timer/modules/model/taskmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type innerImpl struct {
}

func (s *innerImpl) GetUserPermDetail(ctx context.Context, teamId int64, account string) (teammd.UserPermDetailDTO, bool) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	r, b, err := teammd.GetUserPermDetail(ctx, teamId, account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	return r, b
}

type outerImpl struct{}

// CreateTeam 创建团队
func (*outerImpl) CreateTeam(ctx context.Context, reqDTO CreateTeamReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		var sysCfg cfgsrv.SysCfg
		// 判断是否允许用户自行创建团队
		sysCfg, err = cfgsrv.Inner.GetSysCfg(ctx)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		if !sysCfg.AllowUserCreateTeam {
			err = util.UnauthorizedError()
			return
		}
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 创建项目
		pu, err2 := teammd.InsertTeam(ctx, teammd.InsertTeamReqDTO{
			Name: reqDTO.Name,
		})
		if err2 != nil {
			return err2
		}
		// 创建管理员组
		role, err2 := teammd.InsertRole(ctx, teammd.InsertRoleReqDTO{
			Name:       i18n.GetByKey(i18n.TeamAdminUserGroupName),
			TeamId:     pu.Id,
			PermDetail: perm.DefaultPermDetail,
			IsAdmin:    true,
		})
		if err2 != nil {
			return err2
		}
		// 创建关联关系
		return teammd.InsertUser(ctx, teammd.InsertUserReqDTO{
			TeamId:  pu.Id,
			Account: reqDTO.Operator.Account,
			RoleId:  role.Id,
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

// UpdateTeam 编辑团队
func (*outerImpl) UpdateTeam(ctx context.Context, reqDTO UpdateTeamReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkTeamAdmin(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
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

// GetTeam 获取团队信息
func (*outerImpl) GetTeam(ctx context.Context, reqDTO GetTeamReqDTO) (TeamWithPermDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TeamWithPermDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	team, b, err := teammd.GetByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TeamWithPermDTO{}, util.InternalError(err)
	}
	if !b {
		return TeamWithPermDTO{}, util.InvalidArgsError()
	}
	if reqDTO.Operator.IsAdmin {
		return TeamWithPermDTO{
			Id:      team.Id,
			Name:    team.Name,
			IsAdmin: true,
			Perm:    perm.DefaultTeamPerm,
		}, nil
	}
	userPerm, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TeamWithPermDTO{}, err
	}
	if !b {
		return TeamWithPermDTO{}, util.UnauthorizedError()
	}
	return TeamWithPermDTO{
		Id:      team.Id,
		Name:    team.Name,
		IsAdmin: userPerm.IsAdmin,
		Perm:    userPerm.PermDetail.TeamPerm,
	}, nil
}

// ListUserByTeamId 获取成员账号
func (*outerImpl) ListUserByTeamId(ctx context.Context, reqDTO ListUserByTeamIdReqDTO) ([]UserDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		p, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, err
		}
		if !b && !p.IsAdmin {
			return nil, nil
		}
	}
	accounts, err := teammd.ListUserAccountByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	users, err := usermd.ListUserByAccounts(ctx, accounts, []string{"account", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(users, func(t usermd.User) (UserDTO, error) {
		return UserDTO{
			Account: t.Account,
			Name:    t.Name,
		}, nil
	})
}

// ListRoleUser 展示角色成员
func (*outerImpl) ListRoleUser(ctx context.Context, reqDTO ListRoleUserReqDTO) ([]RoleUserDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkTeamAdmin(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return nil, err
	}
	teamUsers, err := teammd.ListUserByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	// 角色信息
	roleList, err := teammd.ListRole(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	roleMap := make(map[int64]teammd.Role, len(roleList))
	for _, role := range roleList {
		roleMap[role.Id] = role
	}
	teamUserMap := make(map[string]teammd.User, len(teamUsers))
	for _, teamUser := range teamUsers {
		teamUserMap[teamUser.Account] = teamUser
	}
	// 用户姓名信息
	accounts, _ := listutil.Map(teamUsers, func(t teammd.User) (string, error) {
		return t.Account, nil
	})
	users, err := usermd.ListUserByAccounts(ctx, accounts, []string{"account", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	userMap := make(map[string]string, len(users))
	for _, user := range users {
		userMap[user.Account] = user.Name
	}
	// 保持按角色id排序
	return listutil.Map(teamUsers, func(t teammd.User) (RoleUserDTO, error) {
		roleId := teamUserMap[t.Account].RoleId
		return RoleUserDTO{
			Id:       t.Id,
			Account:  t.Account,
			Name:     userMap[t.Account],
			RoleId:   roleId,
			RoleName: roleMap[roleId].Name,
		}, nil
	})
}

// DeleteUser 删除团队和成员的绑定关系
func (*outerImpl) DeleteUser(ctx context.Context, reqDTO DeleteUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	teamUser, b, err := teammd.GetTeamUserById(ctx, reqDTO.RelationId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if err = checkTeamAdmin(ctx, teamUser.TeamId, reqDTO.Operator); err != nil {
		return err
	}
	_, err = teammd.DeleteUserById(ctx, reqDTO.RelationId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) CreateUser(ctx context.Context, reqDTO CreateUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	role, err := checkTeamUserPermByRoleId(ctx, reqDTO.Operator, reqDTO.RoleId)
	if err != nil {
		return err
	}
	// 去重
	accounts := hashset.NewHashSet(reqDTO.Accounts...).AllKeys()
	usersCount, err := usermd.CountUserByAccounts(ctx, accounts)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if int(usersCount) != len(accounts) {
		return util.InvalidArgsError()
	}
	// 校验重复的数据
	b, err := teammd.ExistUserByTeamIdAndAccounts(ctx, role.TeamId, accounts)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.InvalidArgsError()
	}
	reqList, _ := listutil.Map(accounts, func(t string) (teammd.InsertUserReqDTO, error) {
		return teammd.InsertUserReqDTO{
			TeamId:  role.TeamId,
			Account: t,
			RoleId:  reqDTO.RoleId,
		}, nil
	})
	err = teammd.BatchInsertUser(ctx, reqList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func checkTeamAdmin(ctx context.Context, teamId int64, operator apisession.UserInfo) error {
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
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 不存在或不是管理员角色
	if !b || !p.IsAdmin {
		return util.UnauthorizedError()
	}
	return nil
}

func checkTeamUserPermByRoleId(ctx context.Context, operator apisession.UserInfo, roleId int64) (teammd.Role, error) {
	// 检查权限
	role, b, err := teammd.GetRoleById(ctx, roleId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Role{}, util.InternalError(err)
	}
	if !b {
		return teammd.Role{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return role, nil
	}
	// 检查项目管理员权限
	p, b, err := teammd.GetUserPermDetail(ctx, role.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return role, util.InternalError(err)
	}
	// 不存在或不是管理员角色
	if !b || !p.IsAdmin {
		return role, util.UnauthorizedError()
	}
	return role, nil
}

func (*outerImpl) CreateRole(ctx context.Context, reqDTO CreateRoleReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err := checkTeamAdmin(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return err
	}
	if err := checkReqPerm(ctx, reqDTO.Perm, reqDTO.TeamId); err != nil {
		return err
	}
	if _, err := teammd.InsertRole(ctx, teammd.InsertRoleReqDTO{
		Name:       reqDTO.Name,
		TeamId:     reqDTO.TeamId,
		PermDetail: reqDTO.Perm,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) UpdateRole(ctx context.Context, reqDTO UpdateRoleReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	role, err := checkTeamUserPermByRoleId(ctx, reqDTO.Operator, reqDTO.RoleId)
	if err != nil {
		return err
	}
	// 管理员项目组无法编辑权限
	if role.IsAdmin {
		return util.InvalidArgsError()
	}
	if err = checkReqPerm(ctx, reqDTO.Perm, role.TeamId); err != nil {
		return err
	}
	if _, err = teammd.UpdateRoleById(ctx, teammd.UpdateRoleReqDTO{
		RoleId: reqDTO.RoleId,
		Name:   reqDTO.Name,
		Perm:   reqDTO.Perm,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func checkReqPerm(ctx context.Context, permDetail perm.Detail, teamId int64) error {
	// 检查仓库id
	if len(permDetail.RepoPermList) > 0 {
		repoIdList, _ := listutil.Map(permDetail.RepoPermList, func(t perm.RepoPermWithId) (int64, error) {
			return t.RepoId, nil
		})
		repoList, err := repomd.GetRepoByIdList(ctx, repoIdList)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if len(repoList) != len(repoIdList) {
			return util.InvalidArgsError()
		}
		for _, repo := range repoList {
			if repo.TeamId != teamId {
				return util.InvalidArgsError()
			}
		}
	}
	// 检查应用服务列表
	if len(permDetail.AppPermList) > 0 {
		appIdList, _ := listutil.Map(permDetail.AppPermList, func(t perm.AppPermWithId) (string, error) {
			return t.AppId, nil
		})
		appList, err := appmd.GetByAppIdList(ctx, appIdList)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if len(appList) != len(appIdList) {
			return util.InvalidArgsError()
		}
		for _, app := range appList {
			if app.TeamId != teamId {
				return util.InvalidArgsError()
			}
		}
	}
	return nil
}

func (*outerImpl) DeleteRole(ctx context.Context, reqDTO DeleteRoleReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	_, err := checkTeamUserPermByRoleId(ctx, reqDTO.Operator, reqDTO.RoleId)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := teammd.DeleteRoleById(ctx, reqDTO.RoleId)
		if err2 != nil {
			return err2
		}
		_, err2 = teammd.DeleteAllUserByRoleId(ctx, reqDTO.RoleId)
		return err2
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListRole 角色列表
func (*outerImpl) ListRole(ctx context.Context, reqDTO ListRoleReqDTO) ([]RoleDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkTeamAdmin(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return nil, err
	}
	roles, err := teammd.ListRole(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(roles, func(t teammd.Role) (RoleDTO, error) {
		return RoleDTO{
			RoleId:  t.Id,
			TeamId:  t.TeamId,
			Name:    t.Name,
			Perm:    *t.Perm,
			IsAdmin: t.IsAdmin,
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
	var (
		teamList []teammd.Team
	)
	if reqDTO.Operator.IsAdmin {
		var err error
		teamList, err = teammd.ListAllTeam(ctx)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
	} else {
		puList, err := teammd.ListUserByAccount(ctx, reqDTO.Operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		teamIdList, _ := listutil.Map(puList, func(t teammd.User) (int64, error) {
			return t.TeamId, nil
		})
		teamList, err = teammd.GetTeamsByTeamIdList(ctx, teamIdList)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
	}
	return listutil.Map(teamList, func(t teammd.Team) (TeamDTO, error) {
		return TeamDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
}

// DeleteTeam 删除项目
func (*outerImpl) DeleteTeam(ctx context.Context, reqDTO DeleteTeamReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err = checkTeamAdmin(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return err
	}
	// 检查是否还有挂在该项目组的仓库
	repoCount, err := repomd.CountRepoByTeamId(ctx, reqDTO.TeamId)
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
	appCount, err := appmd.CountAppByTeamId(ctx, reqDTO.TeamId)
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
	// 检查是否挂在团队的定时任务
	taskCount, err := taskmd.CountTaskByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 存在 不允许删除
	if taskCount > 0 {
		err = util.NewBizErr(apicode.DataAlreadyExistsCode, i18n.TimerTaskRemainCountGreaterThanZero)
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 项目表
		_, err2 := teammd.DeleteTeam(ctx, reqDTO.TeamId)
		if err2 != nil {
			return err2
		}
		// 项目用户组
		_, err2 = teammd.DeleteAllRoleByTeamId(ctx, reqDTO.TeamId)
		if err2 != nil {
			return err2
		}
		// 项目用户
		_, err2 = teammd.DeleteAllUserByTeamId(ctx, reqDTO.TeamId)
		if err2 != nil {
			return err2
		}
		// 外部通知模板
		err2 = notifymd.DeleteTplByTeamId(ctx, reqDTO.TeamId)
		if err2 != nil {
			return err2
		}
		// 删除team hook
		return teamhookmd.DeleteTeamHookByTeamId(ctx, reqDTO.TeamId)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return nil
}

// ChangeRole 更换角色
func (*outerImpl) ChangeRole(ctx context.Context, reqDTO ChangeRoleReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	teamUser, b, err := teammd.GetTeamUserById(ctx, reqDTO.RelationId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	b, err = teammd.ExistRoleById(ctx, reqDTO.RoleId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 检查权限
	if err = checkTeamAdmin(ctx, teamUser.TeamId, reqDTO.Operator); err != nil {
		return err
	}
	_, err = teammd.ChangeRoleById(ctx, reqDTO.RelationId, reqDTO.RoleId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListAllByAdmin 管理员查看所有团队
func (*outerImpl) ListAllByAdmin(ctx context.Context, reqDTO ListAllByAdminReqDTO) ([]TeamDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	teams, err := teammd.ListAllTeam(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(teams, func(t teammd.Team) (TeamDTO, error) {
		return TeamDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
}
