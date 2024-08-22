package filesrv

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/fileserv/modules/model/productmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/files"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf-utils/typesniffer"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	avatarStorage  files.Storage
	productStorage files.Storage
	domain         string
	initPsubOnce   = sync.Once{}
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

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.AppProductTopic, func(data any) {
			req, ok := data.(event.AppProductEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppProductDeleteAction:
							return cfg.AppProduct.Delete
						default:
							return false
						}
					}
					return false
				})
			}
		})
	})
}

// UploadAvatar 上传头像
func UploadAvatar(ctx context.Context, reqDTO UploadAvatarReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	file, err := io.ReadAll(reqDTO.Body)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	sniffedType := typesniffer.DetectContentType(file)
	if !sniffedType.IsImage() {
		return "", util.NewBizErr(apicode.InvalidArgsCode, i18n.AvatarNotImageError)
	}
	id := idutil.RandomUuid()
	ext := strings.TrimPrefix(sniffedType.GetMimeType(), "image/")
	_, err = avatarStorage.Save(ctx,
		convertPath(id)+"."+ext,
		bytes.NewReader(file))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	return domain + fmt.Sprintf("/api/files/avatar/get/%s", id+"."+ext), nil
}

func convertPath(id string) string {
	return filepath.Join(id[:8], id[8:16], id[16:24], id[24:])
}

// GetAvatar 获取头像路径
func GetAvatar(ctx context.Context, reqDTO GetAvatarReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	ext := path.Ext(reqDTO.Name)
	id := strings.TrimSuffix(reqDTO.Name, ext)
	if len(id) != 32 {
		return "", util.InvalidArgsError()
	}
	avatarPath := convertPath(id) + ext
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

func UploadProduct(ctx context.Context, reqDTO UploadProductReqDTO) (string, error) {
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

func GetProduct(ctx context.Context, reqDTO GetProductReqDTO) (string, error) {
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
func ListProduct(ctx context.Context, reqDTO ListProductReqDTO) ([]ProductDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return nil, 0, err
	}
	products, total, err := productmd.ListProduct(ctx, productmd.ListProductReqDTO{
		AppId:    reqDTO.AppId,
		Env:      reqDTO.Env,
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret, _ := listutil.Map(products, func(t productmd.Product) (ProductDTO, error) {
		return ProductDTO{
			Id:      t.Id,
			Name:    t.Name,
			Creator: t.Creator,
			Created: t.Created,
		}, nil
	})
	return ret, total, nil
}

// DeleteProduct 删除制品
func DeleteProduct(ctx context.Context, reqDTO DeleteProductReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	product, app, team, err := checkAppDevelopPermByProductId(ctx, reqDTO.ProductId, reqDTO.Operator)
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
	notifyProductEvent(
		reqDTO.Operator,
		team,
		app,
		product,
		event.AppProductDeleteAction,
	)
	return nil
}

func checkAppDevelopPermByAppId(ctx context.Context, appId string, operator apisession.UserInfo) (appmd.App, teammd.Team, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return app, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		return app, team, util.InternalError(err)
	}
	if !b {
		return app, team, util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.GetAppPerm(appId).CanDevelop {
		return app, team, nil
	}
	return app, team, util.UnauthorizedError()
}

func checkAppDevelopPermByProductId(ctx context.Context, productId int64, operator apisession.UserInfo) (productmd.Product, appmd.App, teammd.Team, error) {
	product, b, err := productmd.GetProductById(ctx, productId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return productmd.Product{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return productmd.Product{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	app, team, err := checkAppDevelopPermByAppId(ctx, product.AppId, operator)
	return product, app, team, err
}

func notifyProductEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, product productmd.Product, action event.AppProductEventAction) {
	initPsub()
	psub.Publish(event.AppProductTopic, event.AppProductEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		ProductId:   product.Id,
		ProductName: product.Name,
		Env:         product.Env,
		Action:      action,
	})
}
