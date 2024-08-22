package filesrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"io"
	"time"
)

type UploadAvatarReqDTO struct {
	Body     io.Reader           `json:"-"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UploadAvatarReqDTO) IsValid() error {
	if r.Body == nil {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetAvatarReqDTO struct {
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetAvatarReqDTO) IsValid() error {
	if r.Name == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UploadProductReqDTO struct {
	AppId   string    `json:"appId"`
	Name    string    `json:"name"`
	Creator string    `json:"creator"`
	Env     string    `json:"env"`
	Body    io.Reader `json:"-"`
}

func (r *UploadProductReqDTO) IsValid() error {
	if r.Name == "" || !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if r.Env == "" || r.Creator == "" {
		return util.InvalidArgsError()
	}
	if r.Body == nil {
		return util.InvalidArgsError()
	}
	return nil
}

type GetProductReqDTO struct {
	AppId string `json:"appId"`
	Name  string `json:"name"`
	Env   string `json:"env"`
}

func (r *GetProductReqDTO) IsValid() error {
	if r.Name == "" || !appmd.IsAppIdValid(r.AppId) || r.Env == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type ListProductReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListProductReqDTO) IsValid() error {
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ProductDTO struct {
	Id      int64
	Name    string
	Creator string
	Created time.Time
}

type DeleteProductReqDTO struct {
	ProductId int64               `json:"productId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *DeleteProductReqDTO) IsValid() error {
	if r.ProductId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
