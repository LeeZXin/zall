package promapi

import (
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
)

type CreateScrapeReqVO struct {
	Endpoint   string            `json:"endpoint"`
	AppId      string            `json:"appId"`
	Target     string            `json:"target"`
	TargetType prommd.TargetType `json:"targetType"`
	Env        string            `json:"env"`
}

type UpdateScrapeReqVO struct {
	ScrapeId   int64             `json:"scrapeId"`
	Endpoint   string            `json:"endpoint"`
	Target     string            `json:"target"`
	TargetType prommd.TargetType `json:"targetType"`
}

type ScrapeVO struct {
	Id         int64             `json:"id"`
	Endpoint   string            `json:"endpoint"`
	AppId      string            `json:"appId"`
	Target     string            `json:"target"`
	TargetType prommd.TargetType `json:"targetType"`
	Created    string            `json:"created"`
	Env        string            `json:"env"`
}
