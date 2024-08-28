package alertapi

import (
	"github.com/LeeZXin/zall/pkg/alert"
)

type CreateConfigReqVO struct {
	Name        string      `json:"name"`
	Alert       alert.Alert `json:"alert"`
	AppId       string      `json:"appId"`
	IntervalSec int         `json:"intervalSec"`
	Env         string      `json:"env"`
}

type UpdateConfigReqVO struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	Alert       alert.Alert `json:"alert"`
	IntervalSec int         `json:"intervalSec"`
}

type ListConfigReqVO struct {
	PageNum int    `json:"pageNum"`
	AppId   string `json:"appId"`
	Env     string `json:"env"`
}

type ConfigVO struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	AppId       string      `json:"appId"`
	Content     alert.Alert `json:"content"`
	IntervalSec int         `json:"intervalSec"`
	IsEnabled   bool        `json:"isEnabled"`
	Creator     string      `json:"creator"`
	Env         string      `json:"env"`
}
