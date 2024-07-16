package appapi

type CreateAppReqVO struct {
	AppId  string `json:"appId"`
	TeamId int64  `json:"teamId"`
	Name   string `json:"name"`
}

type UpdateAppReqVO struct {
	AppId string `json:"appId"`
	Name  string `json:"name"`
}

type TransferTeamReqVO struct {
	AppId  string `json:"appId"`
	TeamId int64  `json:"teamId"`
}

type AppVO struct {
	AppId string `json:"appId"`
	Name  string `json:"name"`
}
