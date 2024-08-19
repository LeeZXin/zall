package notifymd

import "github.com/LeeZXin/zall/pkg/notify/notify"

type InsertTplReqDTO struct {
	Name      string
	ApiKey    string
	NotifyCfg notify.Cfg
	TeamId    int64
}

type UpdateTplReqDTO struct {
	Id        int64
	Name      string
	NotifyCfg notify.Cfg
}

type ListTplReqDTO struct {
	Name     string
	TeamId   int64
	PageNum  int
	PageSize int
}
