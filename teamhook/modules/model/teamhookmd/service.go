package teamhookmd

import (
	"context"
	"github.com/LeeZXin/zall/pkg/commonhook"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsTeamHookNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertTeamHook(ctx context.Context, reqDTO InsertTeamHookReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&TeamHook{
			Name:   reqDTO.Name,
			TeamId: reqDTO.TeamId,
			Events: &xormutil.Conversion[teamhook.Events]{
				Data: reqDTO.Events,
			},
			HookType: reqDTO.HookType,
			HookCfg: &xormutil.Conversion[commonhook.Cfg]{
				Data: reqDTO.HookCfg,
			},
		})
	return err
}

func UpdateTeamHook(ctx context.Context, reqDTO UpdateTeamHookReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "events", "hook_type", "hook_cfg").
		Update(&TeamHook{
			Name: reqDTO.Name,
			Events: &xormutil.Conversion[teamhook.Events]{
				Data: reqDTO.Events,
			},
			HookType: reqDTO.HookType,
			HookCfg: &xormutil.Conversion[commonhook.Cfg]{
				Data: reqDTO.HookCfg,
			},
		})
	return rows == 1, err
}

func DeleteTeamHookById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(TeamHook))
	return rows == 1, err
}

func GetTeamHookById(ctx context.Context, id int64) (TeamHook, bool, error) {
	var ret TeamHook
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func ListTeamHookByTeamId(ctx context.Context, teamId int64, cols []string) ([]TeamHook, error) {
	ret := make([]TeamHook, 0)
	session := xormutil.MustGetXormSession(ctx).Where("team_id = ?", teamId)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func DeleteTeamHookByTeamId(ctx context.Context, teamId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Delete(new(TeamHook))
	return err
}
