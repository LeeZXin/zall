package alertapi

import (
	"github.com/LeeZXin/zall/pkg/alert"
)

type InsertConfigReqVO struct {
	Name        string      `json:"name"`
	Alert       alert.Alert `json:"alert"`
	AppId       string      `json:"appId"`
	IntervalSec int         `json:"intervalSec"`
	SilenceSec  int         `json:"silenceSec"`
	Enabled     bool        `json:"enabled"`
}

type UpdateConfigReqVO struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	Alert       alert.Alert `json:"alert"`
	IntervalSec int         `json:"intervalSec"`
	SilenceSec  int         `json:"silenceSec"`
	Enabled     bool        `json:"enabled"`
}

type DeleteConfigReqVO struct {
	Id int64 `json:"id"`
}

type ListConfigReqVO struct {
	Cursor int64  `json:"cursor"`
	Limit  int    `json:"limit"`
	AppId  string `json:"appId"`
}

type ConfigVO struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	AppId       string       `json:"appId"`
	Content     *alert.Alert `json:"content"`
	IntervalSec int          `json:"intervalSec"`
	SilenceSec  int          `json:"silenceSec"`
	Enabled     bool         `json:"enabled"`
	Created     string       `json:"created"`
}
