package productsrv

import (
	"context"
	"github.com/LeeZXin/zall/fileserv/modules/model/productmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

func (*outerImpl) ListProduct(ctx context.Context, reqDTO ListProductReqDTO) ([]ProductDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkAccessProductPerm(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	products, err := productmd.ListProduct(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(products, func(t productmd.Product) (ProductDTO, error) {
		return ProductDTO{
			Name:    t.Name,
			Creator: t.Creator,
			Created: t.Created,
		}, nil
	})
}

func checkAccessProductPerm(ctx context.Context, appId string, operator apisession.UserInfo) error {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	contains, _ := listutil.Contains(p.PermDetail.DevelopAppList, func(s string) (bool, error) {
		return s == appId, nil
	})
	if contains {
		return nil
	}
	return nil
}
