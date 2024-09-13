package artifactmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertArtifact(ctx context.Context, reqDTO InsertArtifactReqDTO) (Artifact, error) {
	ret := Artifact{
		Env:     reqDTO.Env,
		AppId:   reqDTO.AppId,
		Name:    reqDTO.Name,
		Creator: reqDTO.Creator,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetArtifactByAppIdAndNameAndEnv(ctx context.Context, reqDTO GetArtifactReqDTO) (Artifact, bool, error) {
	var ret Artifact
	b, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("name = ?", reqDTO.Name).
		And("env = ?", reqDTO.Env).
		Get(&ret)
	return ret, b, err
}

func GetArtifactById(ctx context.Context, id int64) (Artifact, bool, error) {
	var ret Artifact
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func DeleteArtifactById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(Artifact))
	return rows == 1, err
}

func ListArtifact(ctx context.Context, reqDTO ListArtifactReqDTO) ([]Artifact, int64, error) {
	ret := make([]Artifact, 0)
	total, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("env = ?", reqDTO.Env).
		OrderBy("id desc").
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		FindAndCount(&ret)
	return ret, total, err
}

func ListByAppIdOffsetN(ctx context.Context, appId, env string, offset int) ([]Artifact, error) {
	ret := make([]Artifact, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		Desc("id").
		Limit(99999, offset).
		Find(&ret)
	return ret, err
}

func ListLatestArtifact(ctx context.Context, appId, env string, size int) ([]Artifact, error) {
	ret := make([]Artifact, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		OrderBy("id desc").
		Limit(size).
		Find(&ret)
	return ret, err
}
