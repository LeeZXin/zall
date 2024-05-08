package teammd

import (
	"context"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsRoleNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 64
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
	pu, b, err := GetTeamUser(ctx, teamId, account)
	if err != nil || !b {
		return UserPermDetailDTO{}, b, err
	}
	if pu.RoleId == 0 {
		return UserPermDetailDTO{}, true, nil
	}
	group, b, err := GetByRoleId(ctx, pu.RoleId)
	if err != nil || !b {
		return UserPermDetailDTO{}, b, err
	}
	return UserPermDetailDTO{
		RoleId:     group.Id,
		IsAdmin:    group.IsAdmin,
		PermDetail: *group.Perm,
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

func UpdateUser(ctx context.Context, reqDTO UpdateUserReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId).
		And("account = ?", reqDTO.Account).
		Cols("group_id").
		Limit(1).
		Update(&User{
			RoleId: reqDTO.RoleId,
		})
	return rows == 1, err
}

func DeleteUser(ctx context.Context, teamId int64, account string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		And("account = ?", account).
		Limit(1).
		Delete(new(User))
	return rows == 1, err
}

func ListUser(ctx context.Context, reqDTO ListUserReqDTO) ([]User, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId)
	if reqDTO.Account != "" {
		session.And("account like ?", "%"+reqDTO.Account+"%")
	}
	if reqDTO.Cursor > 0 {
		session.And("id > ?", reqDTO.Cursor)
	}
	ret := make([]User, 0)
	err := session.Limit(reqDTO.Limit).OrderBy("id asc").Find(&ret)
	return ret, err
}

func DeleteAllTeamUserByAccount(ctx context.Context, account string) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		And("account = ?", account).
		Delete(new(User))
}

func GetTeamUser(ctx context.Context, teamId int64, account string) (User, bool, error) {
	ret := User{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		And("account = ?", account).
		Get(&ret)
	return ret, b, err
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

func GetByRoleId(ctx context.Context, id int64) (Role, bool, error) {
	ret := Role{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func UpdateRoleName(ctx context.Context, id int64, name string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("name").
		Limit(1).
		Update(&Role{
			Name: name,
		})
	return rows == 1, err
}

func UpdateRolePerm(ctx context.Context, id int64, detail perm.Detail) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("perm").
		Limit(1).
		Update(&Role{
			Perm: &detail,
		})
	return rows == 1, err
}

func DeleteRole(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(Role))
	return rows == 1, err
}

func ExistRole(ctx context.Context, teamId, roleId int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		And("role_id = ?", roleId).
		Exist(new(User))
}

func ListRole(ctx context.Context, teamId int64) ([]Role, error) {
	ret := make([]Role, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Find(&ret)
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

func ListUserByAccount(ctx context.Context, account string) ([]User, error) {
	ret := make([]User, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		Find(&ret)
	return ret, err
}

func GetByTeamIdList(ctx context.Context, teamIdList []int64) ([]Team, error) {
	ret := make([]Team, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", teamIdList).
		Find(&ret)
	return ret, err
}

func GetByRoleIdList(ctx context.Context, roleIdList []int64) ([]Role, error) {
	ret := make([]Role, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("role_id", roleIdList).
		Find(&ret)
	return ret, err
}

func ListAccountByTeamId(ctx context.Context, teamId int64) ([]string, error) {
	ret := make([]string, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Cols("account").
		Table(TeamUserTableName).
		Find(&ret)
	return ret, err
}
