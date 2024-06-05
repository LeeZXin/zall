package propertysrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/property/modules/model/propertymd"
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
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListEtcdNodeReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
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
	Env       string              `json:"env"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *InsertEtcdNodeReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !propertymd.IsNodeIdValid(r.NodeId) {
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
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteEtcdNodeReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !propertymd.IsNodeIdValid(r.NodeId) {
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
	Env       string              `json:"env"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *UpdateEtcdNodeReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !propertymd.IsNodeIdValid(r.NodeId) {
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
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GrantAuthReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
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

type GetAuthReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetAuthReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
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

type CreateFileReqDTO struct {
	AppId    string              `json:"appId"`
	Name     string              `json:"name"`
	Content  string              `json:"content"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateFileReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !propertymd.IsFileNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type NewVersionReqDTO struct {
	FileId      int64               `json:"fileId"`
	Content     string              `json:"content"`
	LastVersion string              `json:"lastVersion"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *NewVersionReqDTO) IsValid() error {
	if r.LastVersion == "" {
		return util.InvalidArgsError()
	}
	if r.FileId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeletePropContentReqDTO struct {
	Id       int64               `json:"id"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeletePropContentReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListFileReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListFileReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
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

type FileDTO struct {
	Id    int64  `json:"id"`
	AppId string `json:"appId"`
	Name  string `json:"name"`
	Env   string `json:"env"`
}

type DeployPropContentReqDTO struct {
	Id           int64               `json:"id"`
	Version      string              `json:"version"`
	EtcdNodeList []string            `json:"etcdNodeList"`
	Env          string              `json:"env"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *DeployPropContentReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
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

type PageHistoryReqDTO struct {
	FileId   int64               `json:"fileId"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *PageHistoryReqDTO) IsValid() error {
	if r.FileId <= 0 {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
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
	Env       string              `json:"env"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *ListDeployReqDTO) IsValid() error {
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if len(r.Version) > 32 {
		return util.InvalidArgsError()
	}
	if len(r.NodeId) > 0 && !propertymd.IsNodeIdValid(r.NodeId) {
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
	Id          int64     `json:"id"`
	FileId      int64     `json:"fileId"`
	Content     string    `json:"content"`
	Version     string    `json:"version"`
	LastVersion string    `json:"lastVersion"`
	Created     time.Time `json:"created"`
	Creator     string    `json:"creator"`
}

type DeployDTO struct {
	ContentId int64     `json:"contentId"`
	Content   string    `json:"content"`
	Version   string    `json:"version"`
	NodeId    string    `json:"nodeId"`
	Created   time.Time `json:"created"`
	Creator   string    `json:"creator"`
}

type GetHistoryByVersionReqDTO struct {
	FileId   int64               `json:"fileId"`
	Version  string              `json:"version"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetHistoryByVersionReqDTO) IsValid() error {
	if r.FileId <= 0 {
		return util.InvalidArgsError()
	}
	if r.Version == "" || len(r.Version) > 64 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
