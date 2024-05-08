package branchsrv

import (
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/branch"
	"github.com/LeeZXin/zall/util"
)

type CreateProtectedBranchReqDTO struct {
	RepoId   int64                     `json:"repoId"`
	Pattern  string                    `json:"pattern"`
	Cfg      branch.ProtectedBranchCfg `json:"cfg"`
	Operator apisession.UserInfo       `json:"operator"`
}

func (r *CreateProtectedBranchReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !branchmd.IsWildcardBranchValid(r.Pattern) {
		return util.InvalidArgsError()
	}
	if !r.Cfg.IsValid() {
		return util.InvalidArgsError()
	}
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteProtectedBranchReqDTO struct {
	ProtectedBranchId int64               `json:"protectedBranchId"`
	Operator          apisession.UserInfo `json:"operator"`
}

func (r *DeleteProtectedBranchReqDTO) IsValid() error {
	if r.ProtectedBranchId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListProtectedBranchReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListProtectedBranchReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateProtectedBranchReqDTO struct {
	ProtectedBranchId int64
	Pattern           string
	Cfg               branch.ProtectedBranchCfg
	Operator          apisession.UserInfo `json:"operator"`
}

func (r *UpdateProtectedBranchReqDTO) IsValid() error {
	if r.ProtectedBranchId <= 0 {
		return util.InvalidArgsError()
	}
	if !branchmd.IsWildcardBranchValid(r.Pattern) {
		return util.InvalidArgsError()
	}
	if !r.Cfg.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ProtectedBranchDTO struct {
	Id      int64
	Pattern string
	Cfg     branch.ProtectedBranchCfg
}
