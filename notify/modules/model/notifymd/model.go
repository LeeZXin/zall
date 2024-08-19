package notifymd

import (
	"github.com/LeeZXin/zall/pkg/notify/notify"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	TplTableName = "znotify_tpl"
)

type Tpl struct {
	Id        int64                            `json:"id" xorm:"pk autoincr"`
	Name      string                           `json:"name"`
	ApiKey    string                           `json:"apiKey"`
	NotifyCfg *xormutil.Conversion[notify.Cfg] `json:"notifyCfg"`
	TeamId    int64                            `json:"teamId"`
	Created   time.Time                        `json:"created" xorm:"created"`
	Updated   time.Time                        `json:"updated" xorm:"updated"`
}

func (*Tpl) TableName() string {
	return TplTableName
}
