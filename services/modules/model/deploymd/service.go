package deploymd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertDeploy(ctx context.Context, reqDTO InsertDeployReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_" + reqDTO.Env).
		Insert(Deploy{
			AppId:  reqDTO.AppId,
			Config: &reqDTO.Config,
		})
	return err
}

func UpdateDeploy(ctx context.Context, reqDTO UpdateDeployReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_"+reqDTO.Env).
		Where("app_id = ?", reqDTO.AppId).
		Cols("config").
		Update(Deploy{
			Config: &reqDTO.Config,
		})
	return rows == 1, err
}

func GetByAppId(ctx context.Context, appId, env string) (Deploy, bool, error) {
	var ret Deploy
	b, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_"+env).
		Where("app_id = ?", appId).
		Get(&ret)
	return ret, b, err
}

func DeleteByAppId(ctx context.Context, appId, env string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zservice_deploy_"+env).
		Where("app_id = ?", appId).
		Delete(new(Deploy))
	return rows == 1, err
}
