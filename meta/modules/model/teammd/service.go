package teammd

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsGroupNameValid(name string) bool {
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

func GetTeamUserPermDetail(ctx context.Context, teamId int64, account string) (TeamUserPermDetailDTO, bool, error) {
	pu, b, err := GetTeamUser(ctx, teamId, account)
	if err != nil || !b {
		return TeamUserPermDetailDTO{}, b, err
	}
	if pu.GroupId == 0 {
		return TeamUserPermDetailDTO{}, true, nil
	}
	group, b, err := GetByGroupId(ctx, pu.GroupId)
	if err != nil || !b {
		return TeamUserPermDetailDTO{}, b, err
	}
	return TeamUserPermDetailDTO{
		GroupId:    group.Id,
		IsAdmin:    group.IsAdmin,
		PermDetail: group.GetPermDetail(),
	}, true, nil
}

func InsertTeamUser(ctx context.Context, reqDTO InsertTeamUserReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&TeamUser{
		TeamId:  reqDTO.TeamId,
		Account: reqDTO.Account,
		GroupId: reqDTO.GroupId,
	})
	return err
}

func UpdateTeamUser(ctx context.Context, reqDTO UpdateTeamUserReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId).
		And("account = ?", reqDTO.Account).
		Cols("group_id").
		Limit(1).
		Update(&TeamUser{
			GroupId: reqDTO.GroupId,
		})
	return rows == 1, err
}

func DeleteTeamUser(ctx context.Context, teamId int64, account string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		And("account = ?", account).
		Limit(1).
		Delete(new(TeamUser))
	return rows == 1, err
}

func ListTeamUser(ctx context.Context, reqDTO ListTeamUserReqDTO) ([]TeamUser, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId)
	if reqDTO.Account != "" {
		session.And("account like ?", "%"+reqDTO.Account+"%")
	}
	if reqDTO.Cursor > 0 {
		session.And("id > ?", reqDTO.Cursor)
	}
	ret := make([]TeamUser, 0)
	err := session.Limit(reqDTO.Limit).OrderBy("id asc").Find(&ret)
	return ret, err
}

func DeleteAllTeamUserByAccount(ctx context.Context, account string) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		And("account = ?", account).
		Delete(new(TeamUser))
}

func GetTeamUser(ctx context.Context, teamId int64, account string) (TeamUser, bool, error) {
	ret := TeamUser{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		And("account = ?", account).
		Get(&ret)
	return ret, b, err
}

func InsertTeamUserGroup(ctx context.Context, reqDTO InsertTeamUserGroupReqDTO) (TeamUserGroup, error) {
	m, _ := json.Marshal(reqDTO.PermDetail)
	ret := TeamUserGroup{
		TeamId:  reqDTO.TeamId,
		Name:    reqDTO.Name,
		Perm:    string(m),
		IsAdmin: reqDTO.IsAdmin,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetByGroupId(ctx context.Context, id int64) (TeamUserGroup, bool, error) {
	ret := TeamUserGroup{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func UpdateTeamUserGroupName(ctx context.Context, id int64, name string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("name").
		Limit(1).
		Update(&TeamUserGroup{
			Name: name,
		})
	return rows == 1, err
}

func UpdateTeamUserGroupPerm(ctx context.Context, id int64, detail perm.Detail) (bool, error) {
	m, _ := json.Marshal(detail)
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("perm").
		Limit(1).
		Update(&TeamUserGroup{
			Perm: string(m),
		})
	return rows == 1, err
}

func DeleteTeamUserGroup(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(TeamUserGroup))
	return rows == 1, err
}

func ExistTeamUser(ctx context.Context, teamId, groupId int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		And("group_id = ?", groupId).
		Exist(new(TeamUser))
}

func ListTeamUserGroup(ctx context.Context, id int64) ([]TeamUserGroup, error) {
	ret := make([]TeamUserGroup, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", id).
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

func DeleteAllTeamUserGroupByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Delete(new(TeamUserGroup))
}

func DeleteAllTeamUserByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Delete(new(TeamUser))
}

func ListTeamUserByAccount(ctx context.Context, account string) ([]TeamUser, error) {
	ret := make([]TeamUser, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		Find(&ret)
	return ret, err
}

func GetTeamByTeamIdList(ctx context.Context, teamIdList []int64) ([]Team, error) {
	ret := make([]Team, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", teamIdList).
		Find(&ret)
	return ret, err
}

func GetByGroupIdList(ctx context.Context, teamId int64, groupIdList []int64) ([]TeamUserGroup, error) {
	ret := make([]TeamUserGroup, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		In("group_id", groupIdList).
		Find(&ret)
	return ret, err
}
