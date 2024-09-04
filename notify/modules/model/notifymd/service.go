package notifymd

import (
	"context"
	"github.com/LeeZXin/zall/pkg/notify/notify"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsTplNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertTpl(ctx context.Context, reqDTO InsertTplReqDTO) (Tpl, error) {
	ret := Tpl{
		Name:   reqDTO.Name,
		ApiKey: reqDTO.ApiKey,
		NotifyCfg: &xormutil.Conversion[notify.Cfg]{
			Data: reqDTO.NotifyCfg,
		},
		TeamId: reqDTO.TeamId,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateTpl(ctx context.Context, reqDTO UpdateTplReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "notify_cfg").Update(&Tpl{
		Name: reqDTO.Name,
		NotifyCfg: &xormutil.Conversion[notify.Cfg]{
			Data: reqDTO.NotifyCfg,
		},
	})
	return rows == 1, err
}

func DeleteTpl(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(Tpl))
	return rows == 1, err
}

func GetTplById(ctx context.Context, id int64, cols []string) (Tpl, bool, error) {
	var ret Tpl
	session := xormutil.MustGetXormSession(ctx).Where("id = ?", id)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	b, err := session.Get(&ret)
	return ret, b, err
}

func ExistTplById(ctx context.Context, id int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).Where("id = ?", id).Exist(new(Tpl))
}

func UpdateTplApiKeyById(ctx context.Context, id int64, apiKey string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("api_key").
		Update(&Tpl{
			ApiKey: apiKey,
		})
	return rows == 1, err
}

func GetTplByApiKey(ctx context.Context, apiKey string) (Tpl, bool, error) {
	var ret Tpl
	b, err := xormutil.MustGetXormSession(ctx).Where("api_key = ?", apiKey).Get(&ret)
	return ret, b, err
}

func ListTpl(ctx context.Context, reqDTO ListTplReqDTO) ([]Tpl, int64, error) {
	session := xormutil.MustGetXormSession(ctx).Where("team_id = ?", reqDTO.TeamId)
	if reqDTO.Name != "" {
		session.And("name like ?", reqDTO.Name+"%")
	}
	ret := make([]Tpl, 0)
	total, err := session.
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		Desc("id").
		FindAndCount(&ret)
	return ret, total, err
}

func ListAllTplByTeamId(ctx context.Context, teamId int64, cols []string) ([]Tpl, error) {
	session := xormutil.MustGetXormSession(ctx).Where("team_id = ?", teamId)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	ret := make([]Tpl, 0)
	err := session.
		Find(&ret)
	return ret, err
}

func DeleteTplByTeamId(ctx context.Context, teamId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("team_id = ?", teamId).Delete(new(Tpl))
	return err
}
