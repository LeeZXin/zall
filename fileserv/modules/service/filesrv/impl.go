package filesrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/pkg/files"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"os"
	"path/filepath"
)

var (
	iconStorage   files.Storage
	normalStorage files.Storage
	avatarStorage files.Storage
	domain        string
)

func InitStorage() {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("os.Getwd err: %v", err)
	}
	dataDir := filepath.Join(pwd, "data")
	iconDir := filepath.Join(dataDir, "icon")
	avatarDir := filepath.Join(dataDir, "avatar")
	normalDir := filepath.Join(dataDir, "normal")
	iconTempDir := filepath.Join(iconDir, "temp")
	avatarTempDir := filepath.Join(avatarDir, "temp")
	normalTempDir := filepath.Join(normalDir, "temp")
	util.MkdirAll(iconDir, normalDir)
	iconStorage, _ = files.NewLocalStorage(iconDir, iconTempDir)
	normalStorage, _ = files.NewLocalStorage(normalDir, normalTempDir)
	avatarStorage, _ = files.NewLocalStorage(avatarDir, avatarTempDir)
	domain = static.GetString("files.domain")
}

type outerImpl struct{}

func (*outerImpl) UploadIcon(ctx context.Context, reqDTO UploadIconReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", nil
	}
	if !reqDTO.Operator.IsAdmin {
		return "", util.UnauthorizedError()
	}
	id := idutil.RandomUuid()
	_, err := iconStorage.Save(ctx, filepath.Join(convertPointerPath(id), reqDTO.Name), reqDTO.Body)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	return domain + fmt.Sprintf("/api/files/icon/get/%s/%s", id, reqDTO.Name), nil
}

func (*outerImpl) GetIcon(ctx context.Context, reqDTO GetIconReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", nil
	}
	iconPath := filepath.Join(convertPointerPath(reqDTO.Id), reqDTO.Name)
	b, err := iconStorage.Exists(ctx, iconPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", nil
	}
	return filepath.Join(iconStorage.StoreDir(), iconPath), nil
}

func (*outerImpl) UploadAvatar(ctx context.Context, reqDTO UploadAvatarReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", nil
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
		return "", nil
	}
	iconPath := filepath.Join(convertPointerPath(reqDTO.Id), reqDTO.Name)
	b, err := avatarStorage.Exists(ctx, iconPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", nil
	}
	return filepath.Join(iconStorage.StoreDir(), iconPath), nil
}

func (*outerImpl) UploadNormal(ctx context.Context, reqDTO UploadNormalReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", nil
	}
	id := idutil.RandomUuid()
	_, err := normalStorage.Save(ctx, filepath.Join(convertPointerPath(id), reqDTO.Name), reqDTO.Body)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	return domain + fmt.Sprintf("/api/files/normal/get/%s/%s", id, reqDTO.Name), nil
}

func (*outerImpl) GetNormal(ctx context.Context, reqDTO GetNormalReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", nil
	}
	iconPath := filepath.Join(convertPointerPath(reqDTO.Id), reqDTO.Name)
	b, err := normalStorage.Exists(ctx, iconPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", nil
	}
	return filepath.Join(normalStorage.StoreDir(), iconPath), nil
}

func convertPointerPath(id string) string {
	return filepath.Join(id[0:8], id[8:16], id[16:24], id[24:])
}
