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
	"github.com/LeeZXin/zsf-utils/collections/hashset"
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
		// 判断是否允许用户自行创建团队
		sysCfg, b := cfgsrv.Inner.GetSysCfg()
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
		role, err := teammd.InsertRole(ctx, teammd.InsertRoleReqDTO{
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
			RoleId:  role.Id,
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

// IsAdmin 是否是团队管理员
func (*outerImpl) IsAdmin(ctx context.Context, reqDTO IsAdminReqDTO) (bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return false, err
	}
	if reqDTO.Operator.IsAdmin {
		return true, nil
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
func (*outerImpl) GetTeam(ctx context.Context, reqDTO GetTeamReqDTO) (teammd.Team, error) {
	if err := reqDTO.IsValid(); err != nil {
		return teammd.Team{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Team{}, err
	}
	if !b {
		return teammd.Team{}, nil
	}
	team, _, err := teammd.GetByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Team{}, util.InternalError(err)
	}
	return team, nil
}

// ListUserByTeamId 获取成员账号
func (*outerImpl) ListUserByTeamId(ctx context.Context, reqDTO ListUserByTeamIdReqDTO) ([]UserDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, err
	}
	if !b {
		return nil, nil
	}
	accounts, err := teammd.ListUserAccountByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	users, err := usermd.ListUserByAccounts(ctx, accounts)
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
	users, err := usermd.ListUserByAccounts(ctx, accounts)
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
	if len(permDetail.DevelopAppList) > 0 {
		appList, err := appmd.GetByAppIdList(ctx, permDetail.DevelopAppList)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if len(appList) != len(permDetail.DevelopAppList) {
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
func (*outerImpl) ListTeam(ctx context.Context, reqDTO ListTeamReqDTO) ([]teammd.Team, error) {
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
	teamList, err := teammd.GetTeamsByTeamIdList(ctx, teamIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return teamList, nil
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
	if err = checkTeamAdmin(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
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
