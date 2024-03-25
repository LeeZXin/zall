package appmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
)

var (
	validAppIdPattern = regexp.MustCompile("[\\w-]{1,32}")
)

func IsAppIdValid(app string) bool {
	return validAppIdPattern.MatchString(app)
}

func IsAppNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertApp(ctx context.Context, reqDTO InsertAppReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&App{
		AppId:  reqDTO.AppId,
		TeamId: reqDTO.TeamId,
		Name:   reqDTO.Name,
	})
	return err
}

func UpdateApp(ctx context.Context, reqDTO UpdateAppReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		Cols("name").
		Update(&App{
			Name: reqDTO.Name,
		})
	return rows == 1, err
}

func DeleteByAppId(ctx context.Context, appId string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("app_id = ?", appId).Limit(1).Delete(new(App))
	return rows == 1, err
}

func ListApp(ctx context.Context, reqDTO ListAppReqDTO) ([]App, error) {
	session := xormutil.MustGetXormSession(ctx).Where("team_id = ?", reqDTO.TeamId)
	if reqDTO.AppId != "" {
		session.And("app_id like ?", reqDTO.AppId+"%")
	}
	ret := make([]App, 0)
	return ret, session.OrderBy("id asc").Find(&ret)
}

func GetByAppIdList(ctx context.Context, appIdList []string) ([]App, error) {
	ret := make([]App, 0)
	err := xormutil.MustGetXormSession(ctx).In("app_id", appIdList).Find(&ret)
	return ret, err
}

func GetByAppId(ctx context.Context, appId string) (App, bool, error) {
	var ret = App{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Get(&ret)
	return ret, b, err
}

func CountByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Count(new(App))
}

func TransferTeam(ctx context.Context, appId string, teamId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Cols("team_id").
		Update(&App{
			TeamId: teamId,
		})
	return rows == 1, err
}
