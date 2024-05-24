package repomd

import (
	"context"
	"errors"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func GetByPath(ctx context.Context, path string) (Repo, bool, error) {
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
		Update(&Repo{
			LfsSize: lfsSize,
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
		Cfg:           &reqDTO.Cfg,
		RepoStatus:    reqDTO.RepoStatus,
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

func CountByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Count(new(Repo))
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
