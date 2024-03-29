package branchmd

import (
	"context"
	"encoding/json"
	"github.com/IGLOU-EU/go-wildcard/v2"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
)

var (
	validBranchPattern = regexp.MustCompile(`^\S{1,32}$`)
)

func IsWildcardBranchValid(branch string) bool {
	return validBranchPattern.MatchString(branch)
}

func InsertProtectedBranch(ctx context.Context, reqDTO InsertProtectedBranchReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ProtectedBranch{
		Branch: reqDTO.Branch,
		RepoId: reqDTO.RepoId,
		Cfg:    reqDTO.Cfg.ToString(),
	})
	return err
}

func DeleteById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		And("id = ?", id).
		Limit(1).
		Delete(new(ProtectedBranch))
	return rows == 1, err
}

func GetProtectedBranch(ctx context.Context, repoId int64, branch string) (ProtectedBranchDTO, bool, error) {
	ret := ProtectedBranch{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		And("branch = ?", branch).
		Limit(1).
		Get(&ret)
	return protectedBranch2DTO(ret), b, err
}

func GetById(ctx context.Context, id int64) (ProtectedBranchDTO, bool, error) {
	ret := ProtectedBranch{}
	b, err := xormutil.MustGetXormSession(ctx).
		And("id = ?", id).
		Get(&ret)
	return protectedBranch2DTO(ret), b, err
}

func ListProtectedBranch(ctx context.Context, repoId int64) ([]ProtectedBranchDTO, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId)
	ret := make([]ProtectedBranch, 0)
	if err := session.Find(&ret); err != nil {
		return nil, err
	}
	return listutil.Map(ret, func(t ProtectedBranch) (ProtectedBranchDTO, error) {
		return protectedBranch2DTO(t), nil
	})
}

func IsProtectedBranch(ctx context.Context, repoId int64, branch string) (ProtectedBranchCfg, bool, error) {
	pbList, err := ListProtectedBranch(ctx, repoId)
	if err != nil {
		return ProtectedBranchCfg{}, false, err
	}
	for _, pb := range pbList {
		if wildcard.Match(pb.Branch, branch) {
			return pb.Cfg, true, nil
		}
	}
	return ProtectedBranchCfg{}, false, nil
}

func protectedBranch2DTO(b ProtectedBranch) ProtectedBranchDTO {
	var cfg ProtectedBranchCfg
	json.Unmarshal([]byte(b.Cfg), &cfg)
	return ProtectedBranchDTO{
		Id:     b.Id,
		RepoId: b.RepoId,
		Branch: b.Branch,
		Cfg:    cfg,
	}
}
