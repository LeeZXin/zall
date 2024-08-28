package teamhookmd

import (
	"github.com/LeeZXin/zall/pkg/commonhook"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	TeamHookTableName = "zall_team_hook"
)

type TeamHook struct {
	Id       int64                                 `json:"id" xorm:"pk autoincr"`
	Name     string                                `json:"name"`
	TeamId   int64                                 `json:"teamId"`
	Events   *xormutil.Conversion[teamhook.Events] `json:"events"`
	HookType commonhook.HookType                   `json:"hookType"`
	HookCfg  *xormutil.Conversion[commonhook.Cfg]  `json:"hookCfg"`
	Created  time.Time                             `json:"created" xorm:"created"`
	Updated  time.Time                             `json:"updated" xorm:"updated"`
}

func (*TeamHook) TableName() string {
	return TeamHookTableName
}

func (w *TeamHook) GetEvents() teamhook.Events {
	if w.Events == nil {
		return teamhook.Events{}
	}
	return w.Events.Data
}

func (w *TeamHook) GetHookCfg() commonhook.Cfg {
	if w.HookCfg == nil {
		return commonhook.Cfg{}
	}
	return w.HookCfg.Data
}
