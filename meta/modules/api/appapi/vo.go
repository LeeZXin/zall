package appapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type InsertAppReqVO struct {
	AppId  string `json:"appId"`
	TeamId int64  `json:"teamId"`
	Name   string `json:"name"`
}

type DeleteAppReqVO struct {
	AppId string `json:"appId"`
}

type UpdateAppReqVO struct {
	AppId string `json:"appId"`
	Name  string `json:"name"`
}

type ListAppReqVO struct {
	AppId  string `json:"appId"`
	TeamId int64  `json:"teamId"`
	Cursor int64  `json:"cursor"`
	Limit  int    `json:"limit"`
}

type TransferTeamReqVO struct {
	AppId  string `json:"appId"`
	TeamId int64  `json:"teamId"`
}

type AppVO struct {
	AppId string `json:"appId"`
	Name  string `json:"name"`
}

type ListAppRespVO struct {
	ginutil.BaseResp
	Data   []AppVO `json:"data"`
	Cursor int64   `json:"cursor"`
}
