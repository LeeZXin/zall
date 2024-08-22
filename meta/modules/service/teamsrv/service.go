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
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/model/teamhookmd"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/timer/modules/model/timermd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"strings"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.TeamTopic, func(data any) {
			req, ok := data.(event.TeamEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.TeamCreateAction:
						return events.Team.Create
					case event.TeamUpdateAction:
						return events.Team.Update
					case event.TeamDeleteAction:
						return events.Team.Delete
					default:
						return false
					}
				})
			}
		})
		psub.Subscribe(event.TeamRoleTopic, func(data any) {
			req, ok := data.(event.TeamRoleEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.TeamRoleCreateAction:
						return events.TeamRole.Create
					case event.TeamRoleUpdateAction:
						return events.TeamRole.Update
					case event.TeamRoleDeleteAction:
						return events.TeamRole.Delete
					default:
						return false
					}
				})
			}
		})
		psub.Subscribe(event.TeamUserTopic, func(data any) {
			req, ok := data.(event.TeamUserEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.TeamUserCreateAction:
						return events.TeamUser.Create
					case event.TeamUserChangeRoleAction:
						return events.TeamUser.ChangeRole
					case event.TeamUserDeleteAction:
						return events.TeamUser.Delete
					default:
						return false
					}
				})
			}
		})
	})
}

func GetUserPermDetail(ctx context.Context, teamId int64, account string) (teammd.UserPermDetailDTO, bool) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	r, b, err := teammd.GetUserPermDetail(ctx, teamId, account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	return r, b
}

// CreateTeam 创建团队
func CreateTeam(ctx context.Context, reqDTO CreateTeamReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		var sysCfg cfgsrv.SysCfg
		// 判断是否允许用户自行创建团队
		sysCfg, err := cfgsrv.GetSysCfg(ctx)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !sysCfg.AllowUserCreateTeam {
			return util.UnauthorizedError()
		}
	}
	var team teammd.Team
	err := xormstore.WithTx(ctx, func(ctx context.Context) error {
		var err2 error
		// 创建项目
		team, err2 = teammd.InsertTeam(ctx, teammd.InsertTeamReqDTO{
			Name: reqDTO.Name,
		})
		if err2 != nil {
			return err2
		}
		// 创建管理员组
		role, err2 := teammd.InsertRole(ctx, teammd.InsertRoleReqDTO{
			Name:       i18n.GetByKey(i18n.TeamAdminUserGroupName),
			TeamId:     team.Id,
			PermDetail: perm.DefaultPermDetail,
			IsAdmin:    true,
		})
		if err2 != nil {
			return err2
		}
		// 创建关联关系
		return teammd.InsertUser(ctx, teammd.InsertUserReqDTO{
			TeamId:  team.Id,
			Account: reqDTO.Operator.Account,
			RoleId:  role.Id,
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTeamEvent(
		reqDTO.Operator,
		team,
		event.TeamCreateAction,
	)
	return nil
}

// UpdateTeam 编辑团队
func UpdateTeam(ctx context.Context, reqDTO UpdateTeamReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	team, err := checkTeamAdminByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = teammd.UpdateTeam(ctx, teammd.UpdateTeamReqDTO{
		TeamId: reqDTO.TeamId,
		Name:   reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTeamEvent(
		reqDTO.Operator,
		team,
		event.TeamUpdateAction,
	)
	return nil
}

// GetTeam 获取团队信息
func GetTeam(ctx context.Context, reqDTO GetTeamReqDTO) (TeamWithPermDTO, error) {
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
func ListUserByTeamId(ctx context.Context, reqDTO ListUserByTeamIdReqDTO) ([]UserDTO, error) {
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
func ListRoleUser(ctx context.Context, reqDTO ListRoleUserReqDTO) ([]RoleUserDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if _, err := checkTeamAdminByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
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
func DeleteUser(ctx context.Context, reqDTO DeleteUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	teamUser, team, err := checkTeamAdminPermByRelationId(ctx, reqDTO.Operator, reqDTO.RelationId)
	if err != nil {
		return err
	}
	role, b, err := teammd.GetRoleById(ctx, teamUser.RoleId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	user, b, err := usermd.GetByAccount(ctx, teamUser.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	_, err = teammd.DeleteUserById(ctx, reqDTO.RelationId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTeamUserEvent(
		reqDTO.Operator,
		team,
		role,
		formatUserAccountName(user),
		event.TeamUserDeleteAction,
	)
	return nil
}

func CreateUser(ctx context.Context, reqDTO CreateUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	role, team, err := checkTeamAdminPermByRoleId(ctx, reqDTO.Operator, reqDTO.RoleId)
	if err != nil {
		return err
	}
	// 去重
	accounts := hashset.NewHashSet(reqDTO.Accounts...).AllKeys()
	users, err := usermd.ListUserByAccounts(ctx, accounts, []string{"account", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if len(users) != len(accounts) {
		return util.InvalidArgsError()
	}
	// 校验重复的数据
	b, err := teammd.ExistUserByTeamIdAndAccounts(ctx, role.TeamId, accounts)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
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
	userAccountNames, _ := listutil.Map(users, func(t usermd.User) (string, error) {
		return formatUserAccountName(t), nil
	})
	notifyTeamUserEvent(
		reqDTO.Operator,
		team,
		role,
		strings.Join(userAccountNames, ","),
		event.TeamUserCreateAction,
	)
	return nil
}

func formatUserAccountName(t usermd.User) string {
	return t.Account + "(" + t.Name + ")"
}

func checkTeamAdminByTeamId(ctx context.Context, teamId int64, operator apisession.UserInfo) (teammd.Team, error) {
	team, b, err := teammd.GetByTeamId(ctx, teamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return teammd.Team{}, util.InvalidArgsError()
	}
	// 如果是企业管理员
	if operator.IsAdmin {
		return team, nil
	}
	// 判断权限
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return team, util.InternalError(err)
	}
	// 不存在或不是管理员角色
	if !b || !p.IsAdmin {
		return team, util.UnauthorizedError()
	}
	return team, nil
}

func checkTeamAdminPermByRoleId(ctx context.Context, operator apisession.UserInfo, roleId int64) (teammd.Role, teammd.Team, error) {
	// 检查权限
	role, b, err := teammd.GetRoleById(ctx, roleId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Role{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return teammd.Role{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, err := checkTeamAdminByTeamId(ctx, role.TeamId, operator)
	return role, team, err
}

func checkTeamAdminPermByRelationId(ctx context.Context, operator apisession.UserInfo, relationId int64) (teammd.User, teammd.Team, error) {
	user, b, err := teammd.GetTeamUserById(ctx, relationId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.User{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return teammd.User{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, err := checkTeamAdminByTeamId(ctx, user.TeamId, operator)
	return user, team, err
}

func CreateRole(ctx context.Context, reqDTO CreateRoleReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	team, err := checkTeamAdminByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = checkReqPerm(ctx, reqDTO.Perm, reqDTO.TeamId)
	if err != nil {
		return err
	}
	role, err := teammd.InsertRole(ctx, teammd.InsertRoleReqDTO{
		Name:       reqDTO.Name,
		TeamId:     reqDTO.TeamId,
		PermDetail: reqDTO.Perm,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTeamRoleEvent(
		reqDTO.Operator,
		team,
		role,
		event.TeamRoleCreateAction,
	)
	return nil
}

func UpdateRole(ctx context.Context, reqDTO UpdateRoleReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	role, team, err := checkTeamAdminPermByRoleId(ctx, reqDTO.Operator, reqDTO.RoleId)
	if err != nil {
		return err
	}
	// 管理员项目组无法编辑权限
	if role.IsAdmin {
		return util.InvalidArgsError()
	}
	err = checkReqPerm(ctx, reqDTO.Perm, role.TeamId)
	if err != nil {
		return err
	}
	_, err = teammd.UpdateRoleById(ctx, teammd.UpdateRoleReqDTO{
		RoleId: reqDTO.RoleId,
		Name:   reqDTO.Name,
		Perm:   reqDTO.Perm,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTeamRoleEvent(
		reqDTO.Operator,
		team,
		role,
		event.TeamRoleUpdateAction,
	)
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

func DeleteRole(ctx context.Context, reqDTO DeleteRoleReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	role, team, err := checkTeamAdminPermByRoleId(ctx, reqDTO.Operator, reqDTO.RoleId)
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
	notifyTeamRoleEvent(
		reqDTO.Operator,
		team,
		role,
		event.TeamRoleDeleteAction,
	)
	return nil
}

// ListRole 角色列表
func ListRole(ctx context.Context, reqDTO ListRoleReqDTO) ([]RoleDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if _, err := checkTeamAdminByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
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
func ListTeam(ctx context.Context, reqDTO ListTeamReqDTO) ([]TeamDTO, error) {
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
func DeleteTeam(ctx context.Context, reqDTO DeleteTeamReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	team, err := checkTeamAdminByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 检查是否还有挂在该项目组的仓库
	repoCount, err := repomd.CountRepoByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 存在 不允许删除
	if repoCount > 0 {
		return util.NewBizErr(apicode.DataAlreadyExistsCode, i18n.RepoRemainCountGreaterThanZero)
	}
	// 检查是否有挂在该项目组的app
	appCount, err := appmd.CountAppByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 存在 不允许删除
	if appCount > 0 {
		return util.NewBizErr(apicode.DataAlreadyExistsCode, i18n.AppRemainCountGreaterThanZero)
	}
	// 检查是否挂在团队的定时任务
	taskCount, err := timermd.CountTimerByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 存在 不允许删除
	if taskCount > 0 {
		return util.NewBizErr(apicode.DataAlreadyExistsCode, i18n.TimerTaskRemainCountGreaterThanZero)
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
		return util.InternalError(err)
	}
	notifyTeamEvent(
		reqDTO.Operator,
		team,
		event.TeamDeleteAction,
	)
	return nil
}

// ChangeRole 更换角色
func ChangeRole(ctx context.Context, reqDTO ChangeRoleReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	teamUser, team, err := checkTeamAdminPermByRelationId(ctx, reqDTO.Operator, reqDTO.RelationId)
	if err != nil {
		return err
	}
	if teamUser.RoleId == reqDTO.RoleId {
		return nil
	}
	// 检查角色字段
	role, b, err := teammd.GetRoleById(ctx, reqDTO.RoleId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	user, b, err := usermd.GetByAccount(ctx, teamUser.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	_, err = teammd.ChangeRoleById(ctx, reqDTO.RelationId, reqDTO.RoleId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTeamUserEvent(
		reqDTO.Operator,
		team,
		role,
		formatUserAccountName(user),
		event.TeamUserChangeRoleAction,
	)
	return nil
}

// ListAllByAdmin 管理员查看所有团队
func ListAllByAdmin(ctx context.Context, reqDTO ListAllByAdminReqDTO) ([]TeamDTO, error) {
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

func notifyTeamUserEvent(operator apisession.UserInfo, team teammd.Team, role teammd.Role, user string, action event.TeamUserEventAction) {
	initPsub()
	psub.Publish(event.TeamUserTopic, event.TeamUserEvent{
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
		BaseRole: event.BaseRole{
			RoleId:   role.Id,
			RoleName: role.Name,
		},
		User:   user,
		Action: action,
	})
}

func notifyTeamRoleEvent(operator apisession.UserInfo, team teammd.Team, role teammd.Role, action event.TeamRoleEventAction) {
	initPsub()
	psub.Publish(event.TeamRoleTopic, event.TeamRoleEvent{
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
		BaseRole: event.BaseRole{
			RoleId:   role.Id,
			RoleName: role.Name,
		},
		Action: action,
	},
	)
}

func notifyTeamEvent(operator apisession.UserInfo, team teammd.Team, action event.TeamEventAction) {
	initPsub()
	psub.Publish(event.TeamTopic, event.TeamEvent{
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
	})
}
