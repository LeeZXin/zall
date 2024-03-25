package deployapi

import "github.com/LeeZXin/zall/pkg/deploy"

type GetDeployReqVO struct {
	AppId string `json:"appId"`
	Env   string `json:"env"`
}

type UpdateDeployReqVO struct {
	AppId  string        `json:"appId"`
	Env    string        `json:"env"`
	Config deploy.Config `json:"config"`
}
