package productmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertProduct(ctx context.Context, reqDTO InsertProductReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Product{
			Env:     reqDTO.Env,
			AppId:   reqDTO.AppId,
			Name:    reqDTO.Name,
			Creator: reqDTO.Creator,
		})
	return err
}

func GetProductByAppIdAndNameAndEnv(ctx context.Context, reqDTO GetProductReqDTO) (Product, bool, error) {
	var ret Product
	b, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("name = ?", reqDTO.Name).
		And("env = ?", reqDTO.Env).
		Get(&ret)
	return ret, b, err
}

func GetProductById(ctx context.Context, id int64) (Product, bool, error) {
	var ret Product
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func DeleteProductById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(Product))
	return rows == 1, err
}

func ListProduct(ctx context.Context, appId, env string) ([]Product, error) {
	ret := make([]Product, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		OrderBy("id desc").
		Limit(20).
		Find(&ret)
	return ret, err
}
