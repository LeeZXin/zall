package branchmd

import (
	"context"
	"github.com/IGLOU-EU/go-wildcard/v2"
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
		Cfg:    &reqDTO.Cfg,
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

func GetProtectedBranch(ctx context.Context, repoId int64, branch string) (ProtectedBranch, bool, error) {
	ret := ProtectedBranch{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		And("branch = ?", branch).
		Get(&ret)
	return ret, b, err
}

func GetById(ctx context.Context, id int64) (ProtectedBranch, bool, error) {
	ret := ProtectedBranch{}
	b, err := xormutil.MustGetXormSession(ctx).
		And("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ListProtectedBranch(ctx context.Context, repoId int64) ([]ProtectedBranch, error) {
	ret := make([]ProtectedBranch, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		Find(&ret)
	return ret, err
}

func IsProtectedBranch(ctx context.Context, repoId int64, branch string) (ProtectedBranchCfg, bool, error) {
	pbList, err := ListProtectedBranch(ctx, repoId)
	if err != nil {
		return ProtectedBranchCfg{}, false, err
	}
	for _, pb := range pbList {
		if wildcard.Match(pb.Branch, branch) && pb.Cfg != nil {
			return *pb.Cfg, true, nil
		}
	}
	return ProtectedBranchCfg{}, false, nil
}
