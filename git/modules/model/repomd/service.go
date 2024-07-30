package repomd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
	"time"
)

var (
	validRepoNamePattern = regexp.MustCompile(`^[\w-]{1,32}$`)
	validBranchPattern   = regexp.MustCompile(`^[\w.-]{0,128}$`)
)

func IsRepoDescValid(desc string) bool {
	return len(desc) < 255
}

func IsRepoNameValid(name string) bool {
	return validRepoNamePattern.MatchString(name)
}

func IsBranchValid(branch string) bool {
	return validBranchPattern.MatchString(branch)
}

func GetByPath(ctx context.Context, path string) (Repo, bool, error) {
	var ret Repo
	b, err := xormutil.MustGetXormSession(ctx).
		Where("path = ?", path).
		And("is_deleted = 0").
		Get(&ret)
	return ret, b, err
}

func GetByPathWithoutJudgingDeleted(ctx context.Context, path string) (Repo, bool, error) {
	var ret Repo
	b, err := xormutil.MustGetXormSession(ctx).
		Where("path = ?", path).
		Get(&ret)
	return ret, b, err
}

func GetByRepoId(ctx context.Context, repoId int64) (Repo, bool, error) {
	var ret Repo
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		And("is_deleted = 0").
		Get(&ret)
	return ret, b, err
}

func GetByRepoIdWithoutJudgingDeleted(ctx context.Context, repoId int64) (Repo, bool, error) {
	var ret Repo
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Get(&ret)
	return ret, b, err
}

func UpdateGitSizeAndLfsSize(ctx context.Context, repoId int64, gitSize, lfsSize int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Cols("git_size", "lfs_size").
		Update(&Repo{
			GitSize: gitSize,
			LfsSize: lfsSize,
		})
	return err
}

func UpdateGitSize(ctx context.Context, repoId int64, gitSize int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		And("is_deleted = 0").
		Cols("git_size").
		Update(&Repo{
			GitSize: gitSize,
		})
	return err
}

func UpdateLastOperated(ctx context.Context, repoId int64, lastOperated time.Time) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Cols("last_operated").
		Update(&Repo{
			LastOperated: lastOperated,
		})
	return rows == 1, err
}

func ListRepoByTeamId(ctx context.Context, teamId int64) ([]Repo, error) {
	ret := make([]Repo, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		And("is_deleted = 0").
		Find(&ret)
	return ret, err
}

func GetDeletedRepoListByTeamId(ctx context.Context, teamId int64) ([]Repo, error) {
	ret := make([]Repo, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		And("is_deleted = 1").
		Desc("deleted").
		Find(&ret)
	return ret, err
}

func GetRepoByIdList(ctx context.Context, repoIdList []int64) ([]Repo, error) {
	ret := make([]Repo, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", repoIdList).
		And("is_deleted = 0").
		Find(&ret)
	return ret, err
}

func InsertRepo(ctx context.Context, reqDTO InsertRepoReqDTO) (Repo, error) {
	r := Repo{
		Name:          reqDTO.Name,
		Path:          reqDTO.Path,
		TeamId:        reqDTO.TeamId,
		RepoDesc:      reqDTO.RepoDesc,
		DefaultBranch: reqDTO.DefaultBranch,
		GitSize:       reqDTO.GitSize,
		LfsSize:       reqDTO.LfsSize,
		Cfg:           &reqDTO.Cfg,
		LastOperated:  reqDTO.LastOperated,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&r)
	return r, err
}

func DeleteRepo(ctx context.Context, repoId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Delete(new(Repo))
	return rows == 1, err
}

func SetRepoDeleted(ctx context.Context, repoId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Cols("is_deleted", "deleted").
		Update(&Repo{
			IsDeleted: true,
			Deleted:   time.Now(),
		})
	return rows == 1, err
}

func SetRepoUnDeleted(ctx context.Context, repoId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Cols("is_deleted").
		Update(&Repo{
			IsDeleted: false,
		})
	return rows == 1, err
}

func CountRepoByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Count(new(Repo))
}

func TransferTeam(ctx context.Context, repoId, teamId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Cols("team_id").
		Limit(1).
		Update(&Repo{
			TeamId: teamId,
		})
	return rows == 1, err
}

func UpdateRepoIsArchived(ctx context.Context, repoId int64, isArchived bool) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", repoId).
		Cols("is_archived").
		Update(&Repo{
			IsArchived: isArchived,
		})
	return rows == 1, err
}

func UpdateRepo(ctx context.Context, reqDTO UpdateRepoReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", reqDTO.Id).
		Cols("cfg", "repo_desc").
		Update(&Repo{
			RepoDesc: reqDTO.RepoDesc,
			Cfg:      &reqDTO.Cfg,
		})
	return rows == 1, err
}
