package deploymd

import "github.com/LeeZXin/zall/pkg/deploy"

type InsertDeployReqDTO struct {
	AppId  string
	Config deploy.Config
	Env    string
}

type UpdateDeployReqDTO struct {
	AppId  string
	Config deploy.Config
	Env    string
}
