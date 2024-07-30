package filesrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/fileserv/modules/model/productmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/files"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"net/url"
	"os"
	"path/filepath"
)

var (
	avatarStorage  files.Storage
	productStorage files.Storage
	domain         string
)

func InitStorage() {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("os.Getwd err: %v", err)
	}
	dataDir := filepath.Join(pwd, "data")
	avatarDir := filepath.Join(dataDir, "avatar")
	productDir := filepath.Join(dataDir, "product")
	avatarTempDir := filepath.Join(avatarDir, "temp")
	productTempDir := filepath.Join(productDir, "temp")
	util.MkdirAll(
		avatarDir, productDir,
		avatarTempDir, productTempDir,
	)
	avatarStorage, _ = files.NewLocalStorage(avatarDir, avatarTempDir)
	productStorage, _ = files.NewLocalStorage(productDir, productTempDir)
	domain = static.GetString("files.domain")
}

type outerImpl struct{}

func (*outerImpl) UploadAvatar(ctx context.Context, reqDTO UploadAvatarReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	id := idutil.RandomUuid()
	_, err := avatarStorage.Save(ctx, filepath.Join(convertPointerPath(id), reqDTO.Name), reqDTO.Body)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	return domain + fmt.Sprintf("/api/files/avatar/get/%s/%s", id, reqDTO.Name), nil
}

func (*outerImpl) GetAvatar(ctx context.Context, reqDTO GetAvatarReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	avatarPath := filepath.Join(convertPointerPath(reqDTO.Id), reqDTO.Name)
	b, err := avatarStorage.Exists(ctx, avatarPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", nil
	}
	return filepath.Join(avatarStorage.StoreDir(), avatarPath), nil
}

func (*outerImpl) UploadProduct(ctx context.Context, reqDTO UploadProductReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := productmd.GetProductByAppIdAndNameAndEnv(ctx, productmd.GetProductReqDTO{
		AppId: reqDTO.AppId,
		Name:  reqDTO.Name,
		Env:   reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if b {
		return "", util.AlreadyExistsError()
	}
	_, b, err = appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", util.InvalidArgsError()
	}
	err = productmd.InsertProduct(ctx, productmd.InsertProductReqDTO{
		AppId:   reqDTO.AppId,
		Name:    reqDTO.Name,
		Creator: reqDTO.Creator,
		Env:     reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	_, err = productStorage.Save(ctx, filepath.Join(reqDTO.Env, reqDTO.AppId, reqDTO.Name), reqDTO.Body)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	return domain + fmt.Sprintf("/api/files/product/get/%s/%s/%s", reqDTO.AppId, url.QueryEscape(reqDTO.Name), reqDTO.Env), nil
}

func (*outerImpl) GetProduct(ctx context.Context, reqDTO GetProductReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	productPath := filepath.Join(reqDTO.Env, reqDTO.AppId, reqDTO.Name)
	b, err := productStorage.Exists(ctx, productPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", nil
	}
	return filepath.Join(productStorage.StoreDir(), productPath), nil
}

// ListProduct 制品库列表
func (*outerImpl) ListProduct(ctx context.Context, reqDTO ListProductReqDTO) ([]ProductDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkAppDevelopPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
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
			Id:      t.Id,
			Name:    t.Name,
			Creator: t.Creator,
			Created: t.Created,
		}, nil
	})
}

// DeleteProduct 删除制品
func (*outerImpl) DeleteProduct(ctx context.Context, reqDTO DeleteProductReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	product, err := checkAppDevelopPermByProductId(ctx, reqDTO.ProductId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := productmd.DeleteProductById(ctx, reqDTO.ProductId)
		if err2 != nil {
			return err2
		}
		return productStorage.Delete(ctx, filepath.Join(product.Env, product.AppId, product.Name))
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func checkAppDevelopPermByAppId(ctx context.Context, appId string, operator apisession.UserInfo) error {
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
	if p.PermDetail.GetAppPerm(appId).CanDevelop {
		return nil
	}
	return nil
}

func checkAppDevelopPermByProductId(ctx context.Context, productId int64, operator apisession.UserInfo) (productmd.Product, error) {
	product, b, err := productmd.GetProductById(ctx, productId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return productmd.Product{}, util.InternalError(err)
	}
	if !b {
		return productmd.Product{}, util.InvalidArgsError()
	}
	return product, checkAppDevelopPermByAppId(ctx, product.AppId, operator)
}

func convertPointerPath(id string) string {
	return filepath.Join(id[0:8], id[8:16], id[16:24], id[24:])
}
