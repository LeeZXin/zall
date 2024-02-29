package repomd

import (
	"context"
	"errors"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsAccessTokenAccountValid(account string) bool {
	return len(account) == 16
}

func IsAccessTokenTokenValid(token string) bool {
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

func ListAllRepo(ctx context.Context, teamId int64) ([]Repo, error) {
	ret := make([]Repo, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Find(&ret)
	return ret, err
}

func ListRepoByIdList(ctx context.Context, repoIdList []int64) ([]Repo, error) {
	ret := make([]Repo, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", repoIdList).
		Find(&ret)
	return ret, err
}

func UpdateIsEmpty(ctx context.Context, repoId string, isEmpty bool) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", repoId).
		Cols("is_empty").
		Limit(1).
		Update(&Repo{
			IsEmpty: isEmpty,
		})
	return err
}

func InsertRepo(ctx context.Context, reqDTO InsertRepoReqDTO) (Repo, error) {
	r := Repo{
		Name:          reqDTO.Name,
		Path:          reqDTO.Path,
		Author:        reqDTO.Author,
		TeamId:        reqDTO.TeamId,
		RepoDesc:      reqDTO.RepoDesc,
		DefaultBranch: reqDTO.DefaultBranch,
		IsEmpty:       reqDTO.IsEmpty,
		GitSize:       reqDTO.GitSize,
		LfsSize:       reqDTO.LfsSize,
		NodeId:        reqDTO.NodeId,
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

func GetAccessToken(ctx context.Context, reqDTO GetAccessTokenReqDTO) (AccessToken, bool, error) {
	var ret AccessToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", reqDTO.RepoId).
		And("account = ?", reqDTO.Account).
		Get(&ret)
	return ret, b, err
}

func DeleteAccessToken(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(AccessToken))
	return rows == 1, err
}

func InsertAccessToken(ctx context.Context, reqDTO InsertAccessTokenReqDTO) (AccessToken, error) {
	ret := AccessToken{
		RepoId:  reqDTO.RepoId,
		Account: reqDTO.Account,
		Token:   reqDTO.Token,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetByTid(ctx context.Context, tid int64) (AccessToken, bool, error) {
	var ret AccessToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", tid).
		Get(&ret)
	return ret, b, err
}

func ListAccessToken(ctx context.Context, repoId int64) ([]AccessToken, error) {
	ret := make([]AccessToken, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		Find(&ret)
	return ret, err
}

func InsertAction(ctx context.Context, reqDTO InsertActionReqDTO) error {
	ret := Action{
		RepoId:         reqDTO.RepoId,
		Content:        reqDTO.Content,
		AssignInstance: reqDTO.AssignInstance,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func UpdateAction(ctx context.Context, reqDTO UpdateActionReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.ActionId).
		Cols("content", "assign_instance").
		Update(&Action{
			Content:        reqDTO.Content,
			AssignInstance: reqDTO.AssignInstance,
		})
	return rows == 1, err
}

func DeleteAction(ctx context.Context, actionId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", actionId).
		Delete(new(Action))
	return rows == 1, err
}

func ListAction(ctx context.Context, repoId int64) ([]Action, error) {
	ret := make([]Action, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		Find(&ret)
	return ret, err
}

func GetByActionId(ctx context.Context, actionId int64) (Action, bool, error) {
	var ret Action
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", actionId).
		Get(&ret)
	return ret, b, err
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
