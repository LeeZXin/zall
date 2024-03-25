package propsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/prop/modules/model/propmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
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
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
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
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
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
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteEtcdNodeReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
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
	Env       string              `json:"env"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *UpdateEtcdNodeReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
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
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GrantAuthReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
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
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
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

type InsertPropContentReqDTO struct {
	AppId    string              `json:"appId"`
	Name     string              `json:"name"`
	Content  string              `json:"content"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertPropContentReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
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
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdatePropContentReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
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

type DeletePropContentReqDTO struct {
	Id       int64               `json:"id"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeletePropContentReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
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

type ListPropContentReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListPropContentReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
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

type PropContentDTO struct {
	Id    int64  `json:"id"`
	AppId string `json:"appId"`
	Name  string `json:"name"`
}

type DeployPropContentReqDTO struct {
	Id           int64               `json:"id"`
	Version      string              `json:"version"`
	EtcdNodeList []string            `json:"etcdNodeList"`
	Env          string              `json:"env"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *DeployPropContentReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
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

type ListHistoryReqDTO struct {
	ContentId int64               `json:"contentId"`
	Version   string              `json:"version"`
	Cursor    int64               `json:"cursor"`
	Limit     int                 `json:"limit"`
	Env       string              `json:"env"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *ListHistoryReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
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
	Env       string              `json:"env"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *ListDeployReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
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
	Creator   string    `json:"creator"`
}

type DeployDTO struct {
	ContentId int64     `json:"contentId"`
	Content   string    `json:"content"`
	Version   string    `json:"version"`
	NodeId    string    `json:"nodeId"`
	Created   time.Time `json:"created"`
	Creator   string    `json:"creator"`
}
