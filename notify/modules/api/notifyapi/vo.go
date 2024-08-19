package notifyapi

import "github.com/LeeZXin/zall/pkg/notify/notify"

type CreateNotifyTplReqVO struct {
	Name   string     `json:"name"`
	TeamId int64      `json:"teamId"`
	Cfg    notify.Cfg `json:"cfg"`
}

type UpdateNotifyTplReqVO struct {
	Id   int64      `json:"id"`
	Name string     `json:"name"`
	Cfg  notify.Cfg `json:"cfg"`
}

type ListNotifyTplReqVO struct {
	Name    string `json:"name"`
	PageNum int    `json:"pageNum"`
	TeamId  int64  `json:"teamId"`
}

type TplVO struct {
	Id        int64      `json:"id"`
	TeamId    int64      `json:"teamId"`
	Name      string     `json:"name"`
	ApiKey    string     `json:"apiKey"`
	NotifyCfg notify.Cfg `json:"notifyCfg"`
}

type SimpleTplVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
