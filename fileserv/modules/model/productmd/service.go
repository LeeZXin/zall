package productmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertProduct(ctx context.Context, reqDTO InsertProductReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zfile_product_" + reqDTO.Env).
		Insert(&Product{
			AppId:   reqDTO.AppId,
			Name:    reqDTO.Name,
			Creator: reqDTO.Creator,
		})
	return err
}

func GetProduct(ctx context.Context, reqDTO GetProductReqDTO) (Product, bool, error) {
	ret := Product{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("name = ?", reqDTO.Name).
		Get(&ret)
	return ret, b, err
}

func ListProduct(ctx context.Context, appId, env string) ([]Product, error) {
	ret := make([]Product, 0)
	err := xormutil.MustGetXormSession(ctx).
		Table("zfile_product_"+env).
		Where("app_id = ?", appId).
		OrderBy("id desc").
		Limit(20).
		Find(&ret)
	return ret, err
}
