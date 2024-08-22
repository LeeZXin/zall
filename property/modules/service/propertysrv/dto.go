package propertysrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/property/modules/model/propertymd"
	"github.com/LeeZXin/zall/util"
	"time"
)

type PropertySourceDTO struct {
	Id        int64    `json:"id"`
	Name      string   `json:"name"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Env       string   `json:"env"`
}

type SimplePropertySourceDTO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ListPropertySourceReqDTO struct {
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListPropertySourceReqDTO) IsValid() error {
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAllPropertySourceReqDTO struct {
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAllPropertySourceReqDTO) IsValid() error {
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListBindPropertySourceReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListBindPropertySourceReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type BindAppAndPropertySourceReqDTO struct {
	AppId        string              `json:"appId"`
	SourceIdList []int64             `json:"sourceIdList"`
	Env          string              `json:"env"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *BindAppAndPropertySourceReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	for _, i := range r.SourceIdList {
		if i <= 0 {
			return util.InvalidArgsError()
		}
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	return nil
}

type ListPropertySourceByFileIdReqDTO struct {
	FileId   int64               `json:"fileId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListPropertySourceByFileIdReqDTO) IsValid() error {
	if r.FileId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CreatePropertySourceReqDTO struct {
	Name      string              `json:"name"`
	Endpoints []string            `json:"endpoints"`
	Username  string              `json:"username"`
	Password  string              `json:"password"`
	Env       string              `json:"env"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *CreatePropertySourceReqDTO) IsValid() error {
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if len(r.Endpoints) == 0 {
		return util.InvalidArgsError()
	}
	for _, endpoint := range r.Endpoints {
		if !util.GenIpPortPattern().MatchString(endpoint) {
			return util.InvalidArgsError()
		}
	}
	if !propertymd.IsPropertySourceNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeletePropertySourceReqDTO struct {
	SourceId int64               `json:"sourceId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeletePropertySourceReqDTO) IsValid() error {
	if r.SourceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdatePropertySourceReqDTO struct {
	SourceId  int64               `json:"sourceId"`
	Name      string              `json:"name"`
	Endpoints []string            `json:"endpoints"`
	Username  string              `json:"username"`
	Password  string              `json:"password"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *UpdatePropertySourceReqDTO) IsValid() error {
	if r.SourceId <= 0 {
		return util.InvalidArgsError()
	}
	if !propertymd.IsPropertySourceNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if len(r.Endpoints) == 0 {
		return util.InvalidArgsError()
	}
	for _, endpoint := range r.Endpoints {
		if !util.GenIpPortPattern().MatchString(endpoint) {
			return util.InvalidArgsError()
		}
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
	if !cfgsrv.ContainsEnv(r.Env) {
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

type DeleteFileReqDTO struct {
	FileId   int64               `json:"fileId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteFileReqDTO) IsValid() error {
	if r.FileId <= 0 {
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

type FileDTO struct {
	Id    int64  `json:"id"`
	AppId string `json:"appId"`
	Name  string `json:"name"`
	Env   string `json:"env"`
}

type DeployHistoryReqDTO struct {
	HistoryId    int64               `json:"historyId"`
	SourceIdList []int64             `json:"sourceIdList"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *DeployHistoryReqDTO) IsValid() error {
	if r.HistoryId <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.SourceIdList) == 0 {
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
	HistoryId int64               `json:"historyId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *ListDeployReqDTO) IsValid() error {
	if r.HistoryId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type HistoryDTO struct {
	Id          int64     `json:"id"`
	FileName    string    `json:"fileName"`
	FileId      int64     `json:"fileId"`
	Content     string    `json:"content"`
	Version     string    `json:"version"`
	LastVersion string    `json:"lastVersion"`
	Created     time.Time `json:"created"`
	Creator     string    `json:"creator"`
	Env         string    `json:"env"`
}

type DeployDTO struct {
	NodeName  string    `json:"nodeName"`
	Endpoints string    `json:"endpoints"`
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
