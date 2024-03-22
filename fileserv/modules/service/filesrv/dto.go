package filesrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"io"
)

type UploadIconReqDTO struct {
	Name     string              `json:"name"`
	Body     io.Reader           `json:"-"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UploadIconReqDTO) IsValid() error {
	if r.Name == "" {
		return util.InvalidArgsError()
	}
	if r.Body == nil {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetIconReqDTO struct {
	Id       string              `json:"id"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetIconReqDTO) IsValid() error {
	if r.Name == "" || len(r.Id) != 32 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UploadAvatarReqDTO struct {
	Name     string              `json:"name"`
	Body     io.Reader           `json:"-"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UploadAvatarReqDTO) IsValid() error {
	if r.Name == "" {
		return util.InvalidArgsError()
	}
	if r.Body == nil {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetAvatarReqDTO struct {
	Id       string              `json:"id"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetAvatarReqDTO) IsValid() error {
	if r.Name == "" || len(r.Id) != 32 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UploadNormalReqDTO struct {
	Name string    `json:"name"`
	Body io.Reader `json:"-"`
}

func (r *UploadNormalReqDTO) IsValid() error {
	if r.Name == "" {
		return util.InvalidArgsError()
	}
	if r.Body == nil {
		return util.InvalidArgsError()
	}
	return nil
}

type GetNormalReqDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (r *GetNormalReqDTO) IsValid() error {
	if r.Name == "" || len(r.Id) != 32 {
		return util.InvalidArgsError()
	}
	return nil
}

type UploadProductReqDTO struct {
	App  string    `json:"app"`
	Name string    `json:"name"`
	Body io.Reader `json:"-"`
}

func (r *UploadProductReqDTO) IsValid() error {
	if r.Name == "" || !appmd.IsAppIdValid(r.App) {
		return util.InvalidArgsError()
	}
	if r.Body == nil {
		return util.InvalidArgsError()
	}
	return nil
}

type GetProductReqDTO struct {
	App  string `json:"app"`
	Name string `json:"name"`
}

func (r *GetProductReqDTO) IsValid() error {
	if r.Name == "" || !appmd.IsAppIdValid(r.App) {
		return util.InvalidArgsError()
	}
	return nil
}

type ListProductReqDTO struct {
	App string `json:"app"`
}

func (r *ListProductReqDTO) IsValid() error {
	if r.App == "" {
		return util.InvalidArgsError()
	}
	return nil
}
