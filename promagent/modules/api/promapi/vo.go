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

type ScrapeBySaVO struct {
	Id         int64             `json:"id"`
	Endpoint   string            `json:"endpoint"`
	AppId      string            `json:"appId"`
	AppName    string            `json:"appName"`
	TeamId     int64             `json:"teamId"`
	TeamName   string            `json:"teamName"`
	Target     string            `json:"target"`
	TargetType prommd.TargetType `json:"targetType"`
	Env        string            `json:"env"`
}

type ScrapeByTeamVO struct {
	Id         int64             `json:"id"`
	Endpoint   string            `json:"endpoint"`
	Target     string            `json:"target"`
	TargetType prommd.TargetType `json:"targetType"`
	Env        string            `json:"env"`
}

type ListScrapeReqVO struct {
	Endpoint string `json:"endpoint"`
	AppId    string `json:"appId"`
	Env      string `json:"env"`
	PageNum  int    `json:"pageNum"`
}
