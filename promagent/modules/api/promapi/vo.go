package promapi

import (
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
)

type InsertScrapeReqVO struct {
	ServerUrl  string            `json:"serverUrl"`
	AppId      string            `json:"appId"`
	Target     string            `json:"target"`
	TargetType prommd.TargetType `json:"targetType"`
	Env        string            `json:"env"`
}

type UpdateScrapeReqVO struct {
	Id         int64             `json:"id"`
	ServerUrl  string            `json:"serverUrl"`
	Target     string            `json:"target"`
	TargetType prommd.TargetType `json:"targetType"`
	Env        string            `json:"env"`
}

type ListScrapeReqVO struct {
	AppId string `json:"appId"`
	Env   string `json:"env"`
}

type DeleteScrapeReqVO struct {
	Id  int64  `json:"id"`
	Env string `json:"env"`
}

type ScrapeVO struct {
	Id         int64  `json:"id"`
	ServerUrl  string `json:"serverUrl"`
	AppId      string `json:"appId"`
	Target     string `json:"target"`
	TargetType string `json:"targetType"`
	Created    string `json:"created"`
}
