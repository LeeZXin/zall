package branchmd

import (
	"context"
	"github.com/LeeZXin/zall/pkg/branch"
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
		Pattern: reqDTO.Pattern,
		RepoId:  reqDTO.RepoId,
		Cfg:     &reqDTO.Cfg,
	})
	return err
}

func UpdateProtectedBranch(ctx context.Context, reqDTO UpdateProtectedBranchReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("pattern", "cfg").
		Update(&ProtectedBranch{
			Pattern: reqDTO.Pattern,
			Cfg:     &reqDTO.Cfg,
		})
	return rows == 1, err
}

func DeleteProtectedBranchById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		And("id = ?", id).
		Limit(1).
		Delete(new(ProtectedBranch))
	return rows == 1, err
}

func DeleteProtectedBranchByRepoId(ctx context.Context, repoId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		And("repo_id = ?", repoId).
		Limit(1).
		Delete(new(ProtectedBranch))
	return err
}

func GetProtectedBranchById(ctx context.Context, id int64) (ProtectedBranch, bool, error) {
	ret := ProtectedBranch{}
	b, err := xormutil.MustGetXormSession(ctx).
		And("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ListProtectedBranch(ctx context.Context, repoId int64) (ProtectedBranchList, error) {
	ret := make([]ProtectedBranch, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		Find(&ret)
	return ret, err
}

func IsProtectedBranch(ctx context.Context, repoId int64, pattern string) (branch.ProtectedBranchCfg, bool, error) {
	pbList, err := ListProtectedBranch(ctx, repoId)
	if err != nil {
		return branch.ProtectedBranchCfg{}, false, err
	}
	isProtectedBranch, protectedBranch := pbList.IsProtectedBranch(pattern)
	return protectedBranch.GetCfg(), isProtectedBranch, nil
}
