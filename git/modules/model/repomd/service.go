package repomd

import (
	"context"
	"errors"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsRepoTokenAccountValid(account string) bool {
	return len(account) == 16
}

func IsRepoTokenTokenValid(token string) bool {
	return len(token) == 16
}

func GetByPath(ctx context.Context, path string) (RepoInfo, bool, error) {
	var ret Repo
	b, err := xormutil.MustGetXormSession(ctx).
		Where("path = ?", path).
		Get(&ret)
	return ret.ToRepoInfo(), b, err
}

func GetByRepoId(ctx context.Context, repoId int64) (RepoInfo, bool, error) {
	var ret Repo
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Get(&ret)
	return ret.ToRepoInfo(), b, err
}

func UpdateGitSize(ctx context.Context, repoId int64, gitSize int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Cols("git_size").
		Update(&Repo{
			GitSize: gitSize,
		})
	return err
}

func UpdateLfsSize(ctx context.Context, repoId int64, lfsSize int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Cols("lfs_size").
		Incr("lfs_size", lfsSize).
		Update(new(Repo))
	return err
}

func GetRepoListByTeamId(ctx context.Context, teamId int64) ([]Repo, error) {
	ret := make([]Repo, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Find(&ret)
	return ret, err
}

func GetRepoByIdList(ctx context.Context, repoIdList []int64) ([]Repo, error) {
	ret := make([]Repo, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", repoIdList).
		Find(&ret)
	return ret, err
}

func InsertRepo(ctx context.Context, reqDTO InsertRepoReqDTO) (Repo, error) {
	r := Repo{
		Name:          reqDTO.Name,
		Path:          reqDTO.Path,
		Author:        reqDTO.Author,
		TeamId:        reqDTO.TeamId,
		RepoDesc:      reqDTO.RepoDesc,
		DefaultBranch: reqDTO.DefaultBranch,
		GitSize:       reqDTO.GitSize,
		LfsSize:       reqDTO.LfsSize,
		Cfg:           reqDTO.Cfg.ToString(),
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

func CountByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Count(new(Repo))
}

func GetRepoToken(ctx context.Context, reqDTO GetRepoTokenReqDTO) (RepoToken, bool, error) {
	var ret RepoToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", reqDTO.RepoId).
		And("account = ?", reqDTO.Account).
		Get(&ret)
	return ret, b, err
}

func DeleteRepoToken(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(RepoToken))
	return rows == 1, err
}

func InsertRepoToken(ctx context.Context, reqDTO InsertRepoTokenReqDTO) (RepoToken, error) {
	ret := RepoToken{
		RepoId:  reqDTO.RepoId,
		Account: reqDTO.Account,
		Token:   reqDTO.Token,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetByTokenId(ctx context.Context, id int64) (RepoToken, bool, error) {
	var ret RepoToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ListRepoToken(ctx context.Context, repoId int64) ([]RepoToken, error) {
	ret := make([]RepoToken, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		Find(&ret)
	return ret, err
}

func IterateRepo(ctx context.Context, fn func(*Repo) error) error {
	if fn == nil {
		return errors.New("nil iterate fn")
	}
	return xormutil.MustGetXormSession(ctx).
		Iterate(
			new(Repo),
			func(_ int, bean any) error {
				return fn(bean.(*Repo))
			})
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
