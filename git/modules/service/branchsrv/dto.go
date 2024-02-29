package branchsrv

import (
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
)

type InsertProtectedBranchReqDTO struct {
	RepoId   int64                       `json:"repoId"`
	Branch   string                      `json:"branch"`
	Cfg      branchmd.ProtectedBranchCfg `json:"cfg"`
	Operator apisession.UserInfo         `json:"operator"`
}

func (r *InsertProtectedBranchReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !branchmd.IsWildcardBranchValid(r.Branch) {
		return util.InvalidArgsError()
	}
	if len(r.Cfg.ReviewerList) > 50 {
		return util.InvalidArgsError()
	}
	if r.Cfg.ReviewCountWhenCreatePr < len(r.Cfg.ReviewerList) {
		return util.NewBizErr(apicode.InvalidReviewCountWhenCreatePrCode, i18n.ProtectedBranchInvalidReviewCountWhenCreatePr)
	}
	return nil
}

type DeleteProtectedBranchReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteProtectedBranchReqDTO) IsValid() error {
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
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ProtectedBranchDTO struct {
	Id     int64                       `json:"id"`
	RepoId int64                       `json:"repoId"`
	Branch string                      `json:"branch"`
	Cfg    branchmd.ProtectedBranchCfg `json:"cfg"`
}
