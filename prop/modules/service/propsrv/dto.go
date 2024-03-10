package propsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/prop/modules/model/propmd"
	"github.com/LeeZXin/zall/util"
	"time"
)

type EtcdNodeDTO struct {
	NodeId    string   `json:"nodeId"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}

type ListEtcdNodeReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListEtcdNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertEtcdNodeReqDTO struct {
	NodeId    string              `json:"nodeId"`
	Endpoints []string            `json:"endpoints"`
	Username  string              `json:"username"`
	Password  string              `json:"password"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *InsertEtcdNodeReqDTO) IsValid() error {
	if !propmd.IsNodeIdValid(r.NodeId) {
		return util.InvalidArgsError()
	}
	if len(r.Endpoints) == 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteEtcdNodeReqDTO struct {
	NodeId   string              `json:"nodeId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteEtcdNodeReqDTO) IsValid() error {
	if !propmd.IsNodeIdValid(r.NodeId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateEtcdNodeReqDTO struct {
	NodeId    string              `json:"nodeId"`
	Endpoints []string            `json:"endpoints"`
	Username  string              `json:"username"`
	Password  string              `json:"password"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *UpdateEtcdNodeReqDTO) IsValid() error {
	if !propmd.IsNodeIdValid(r.NodeId) {
		return util.InvalidArgsError()
	}
	if len(r.Endpoints) == 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GrantAuthReqDTO struct {
	AppId    string              `json:"appId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GrantAuthReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetAuthReqDTO struct {
	AppId    string              `json:"appId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetAuthReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertPropContentReqDTO struct {
	AppId    string              `json:"appId"`
	Name     string              `json:"name"`
	Content  string              `json:"content"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertPropContentReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !propmd.IsPropContentNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdatePropContentReqDTO struct {
	Id       int64               `json:"id"`
	Content  string              `json:"content"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdatePropContentReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeletePropContentReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeletePropContentReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListPropContentReqDTO struct {
	AppId    string              `json:"appId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListPropContentReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type PropContentDTO struct {
	Id    int64  `json:"id"`
	AppId string `json:"appId"`
	Name  string `json:"name"`
}

type DeployPropContentReqDTO struct {
	Id           int64               `json:"id"`
	Version      string              `json:"version"`
	EtcdNodeList []string            `json:"etcdNodeList"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *DeployPropContentReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.Version) == 0 || len(r.Version) > 32 {
		return util.InvalidArgsError()
	}
	if len(r.EtcdNodeList) == 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListHistoryReqDTO struct {
	ContentId int64               `json:"contentId"`
	Version   string              `json:"version"`
	Cursor    int64               `json:"cursor"`
	Limit     int                 `json:"limit"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *ListHistoryReqDTO) IsValid() error {
	if r.ContentId <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.Version) > 32 {
		return util.InvalidArgsError()
	}
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDeployReqDTO struct {
	ContentId int64               `json:"contentId"`
	NodeId    string              `json:"nodeId"`
	Version   string              `json:"version"`
	Cursor    int64               `json:"cursor"`
	Limit     int                 `json:"limit"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *ListDeployReqDTO) IsValid() error {
	if len(r.Version) > 32 {
		return util.InvalidArgsError()
	}
	if len(r.NodeId) > 0 && !propmd.IsNodeIdValid(r.NodeId) {
		return util.InvalidArgsError()
	}
	if r.ContentId <= 0 {
		return util.InvalidArgsError()
	}
	if r.Cursor < 0 || r.Limit < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type HistoryDTO struct {
	ContentId int64     `json:"contentId"`
	Content   string    `json:"content"`
	Version   string    `json:"version"`
	Created   time.Time `json:"created"`
}

type DeployDTO struct {
	ContentId int64     `json:"contentId"`
	Content   string    `json:"content"`
	Version   string    `json:"version"`
	NodeId    string    `json:"nodeId"`
	Created   time.Time `json:"created"`
}
