package teammd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsRoleNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsTeamNameValid(name string) bool {
	return len(name) > 0 && len(name) < 32
}

func GetByTeamId(ctx context.Context, id int64) (Team, bool, error) {
	var ret Team
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ExistByTeamId(ctx context.Context, id int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Exist(new(Team))
}

func InsertTeam(ctx context.Context, reqDTO InsertTeamReqDTO) (Team, error) {
	ret := Team{
		Name: reqDTO.Name,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateTeam(ctx context.Context, reqDTO UpdateTeamReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.TeamId).
		Cols("name").
		Limit(1).
		Update(&Team{
			Name: reqDTO.Name,
		})
	return rows == 1, err
}

func GetUserPermDetail(ctx context.Context, teamId int64, account string) (UserPermDetailDTO, bool, error) {
	role, b, err := GetRoleByTeamIdAndAccount(ctx, teamId, account)
	if err != nil || !b {
		return UserPermDetailDTO{}, b, err
	}
	return UserPermDetailDTO{
		RoleId:     role.Id,
		IsAdmin:    role.IsAdmin,
		PermDetail: *role.Perm,
	}, true, nil
}

func InsertUser(ctx context.Context, reqDTO InsertUserReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&User{
		TeamId:  reqDTO.TeamId,
		Account: reqDTO.Account,
		RoleId:  reqDTO.RoleId,
	})
	return err
}

func BatchInsertUser(ctx context.Context, reqDTOList []InsertUserReqDTO) error {
	userList, _ := listutil.Map(reqDTOList, func(t InsertUserReqDTO) (*User, error) {
		return &User{
			TeamId:  t.TeamId,
			Account: t.Account,
			RoleId:  t.RoleId,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(userList)
	return err
}

func DeleteUserById(ctx context.Context, relationId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", relationId).
		Delete(new(User))
	return rows == 1, err
}

func DeleteAllTeamUserByAccount(ctx context.Context, account string) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		And("account = ?", account).
		Delete(new(User))
}

func GetTeamUserById(ctx context.Context, relationId int64) (User, bool, error) {
	ret := User{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", relationId).
		Get(&ret)
	return ret, b, err
}

func ChangeRoleById(ctx context.Context, relationId, roleId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", relationId).
		Cols("role_id").
		Update(&User{
			RoleId: roleId,
		})
	return rows == 1, err
}

func ExistUserByTeamIdAndAccounts(ctx context.Context, teamId int64, accounts []string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		In("account", accounts).
		Exist(new(User))
}

func InsertRole(ctx context.Context, reqDTO InsertRoleReqDTO) (Role, error) {
	ret := Role{
		TeamId:  reqDTO.TeamId,
		Name:    reqDTO.Name,
		Perm:    &reqDTO.PermDetail,
		IsAdmin: reqDTO.IsAdmin,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetRoleById(ctx context.Context, id int64) (Role, bool, error) {
	ret := Role{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ExistRoleById(ctx context.Context, id int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Exist(new(Role))
}

func GetRoleByTeamIdAndAccount(ctx context.Context, teamId int64, account string) (Role, bool, error) {
	var ret Role
	b, err := xormutil.MustGetXormSession(ctx).
		SQL("select * from zall_team_role where id = (select role_id from zall_team_user where team_id = ? and account = ? limit 1)",
			teamId, account).
		Get(&ret)
	return ret, b, err
}

func UpdateRoleById(ctx context.Context, reqDTO UpdateRoleReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.RoleId).
		Cols("perm", "name").
		Limit(1).
		Update(&Role{
			Name: reqDTO.Name,
			Perm: &reqDTO.Perm,
		})
	return rows == 1, err
}

func DeleteRoleById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(Role))
	return rows == 1, err
}

func ListRole(ctx context.Context, teamId int64) ([]Role, error) {
	ret := make([]Role, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Find(&ret)
	return ret, err
}

func ListAllTeam(ctx context.Context) ([]Team, error) {
	ret := make([]Team, 0)
	err := xormutil.MustGetXormSession(ctx).Find(&ret)
	return ret, err
}

func DeleteTeam(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(Team))
	return rows == 1, err
}

func DeleteAllRoleByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Delete(new(Role))
}

func DeleteAllUserByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Delete(new(User))
}

func DeleteAllUserByRoleId(ctx context.Context, roleId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("role_id = ?", roleId).
		Delete(new(User))
}

func ListUserByAccount(ctx context.Context, account string) ([]User, error) {
	ret := make([]User, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		Find(&ret)
	return ret, err
}

func GetTeamsByTeamIdList(ctx context.Context, teamIdList []int64) ([]Team, error) {
	ret := make([]Team, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", teamIdList).
		Find(&ret)
	return ret, err
}

func GetRolesByRoleIdList(ctx context.Context, roleIdList []int64) ([]Role, error) {
	ret := make([]Role, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("role_id", roleIdList).
		Find(&ret)
	return ret, err
}

func ListUserAccountByTeamId(ctx context.Context, teamId int64) ([]string, error) {
	ret := make([]string, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Cols("account").
		Table(TeamUserTableName).
		Find(&ret)
	return ret, err
}

func ListUserByTeamId(ctx context.Context, teamId int64) ([]User, error) {
	ret := make([]User, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Asc("role_id", "id").
		Find(&ret)
	return ret, err
}

func IsUserAnyTeamAdmin(ctx context.Context, account string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		And("is_admin = ?", 1).
		Exist(new(User))
}
