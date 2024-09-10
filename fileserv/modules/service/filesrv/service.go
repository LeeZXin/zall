package filesrv

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/fileserv/modules/model/artifactmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
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
	avatarStorage   files.Storage
	artifactStorage files.Storage
	domain          string
	initPsubOnce    = sync.Once{}
)

func InitStorage() {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("os.Getwd err: %v", err)
	}
	dataDir := filepath.Join(pwd, "data")
	avatarDir := filepath.Join(dataDir, "avatar")
	artifactDir := filepath.Join(dataDir, "artifact")
	avatarTempDir := filepath.Join(avatarDir, "temp")
	artifactTempDir := filepath.Join(artifactDir, "temp")
	util.MkdirAll(
		avatarDir, artifactDir,
		avatarTempDir, artifactTempDir,
	)
	avatarStorage, _ = files.NewLocalStorage(avatarDir, avatarTempDir)
	artifactStorage, _ = files.NewLocalStorage(artifactDir, artifactTempDir)
	domain = static.GetString("files.domain")
}

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.AppArtifactTopic, func(data any) {
			req, ok := data.(event.AppArtifactEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppArtifactDeleteAction:
							return cfg.AppArtifact.Delete
						case event.AppArtifactUploadAction:
							return cfg.AppArtifact.Upload
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

func UploadArtifact(ctx context.Context, reqDTO UploadArtifactReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := artifactmd.GetArtifactByAppIdAndNameAndEnv(ctx, artifactmd.GetArtifactReqDTO{
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
	operator, b, err := usermd.GetByAccount(ctx, reqDTO.Creator)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", util.InvalidArgsError()
	}
	app, b, err := appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", util.ThereHasBugErr()
	}
	artifact, err := artifactmd.InsertArtifact(ctx, artifactmd.InsertArtifactReqDTO{
		AppId:   reqDTO.AppId,
		Name:    reqDTO.Name,
		Creator: reqDTO.Creator,
		Env:     reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	_, err = artifactStorage.Save(ctx, filepath.Join(reqDTO.Env, reqDTO.AppId, reqDTO.Name), reqDTO.Body)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	notifyArtifactEvent(
		apisession.UserInfo{
			Account:      operator.Account,
			Name:         operator.Name,
			Email:        operator.Email,
			IsProhibited: operator.IsProhibited,
			AvatarUrl:    operator.AvatarUrl,
			IsAdmin:      operator.IsAdmin,
			IsDba:        operator.IsDba,
		},
		team,
		app,
		artifact,
		event.AppArtifactUploadAction,
	)
	return domain + fmt.Sprintf("/api/files/artifact/get/%s/%s/%s", reqDTO.AppId, url.QueryEscape(reqDTO.Name), reqDTO.Env), nil
}

func GetArtifact(ctx context.Context, reqDTO GetArtifactReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	artifactPath := filepath.Join(reqDTO.Env, reqDTO.AppId, reqDTO.Name)
	b, err := artifactStorage.Exists(ctx, artifactPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", nil
	}
	return filepath.Join(artifactStorage.StoreDir(), artifactPath), nil
}

// ListArtifact 制品库列表
func ListArtifact(ctx context.Context, reqDTO ListArtifactReqDTO) ([]ArtifactDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return nil, 0, err
	}
	artifacts, total, err := artifactmd.ListArtifact(ctx, artifactmd.ListArtifactReqDTO{
		AppId:    reqDTO.AppId,
		Env:      reqDTO.Env,
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret := listutil.MapNe(artifacts, func(t artifactmd.Artifact) ArtifactDTO {
		return ArtifactDTO{
			Id:      t.Id,
			Name:    t.Name,
			Creator: t.Creator,
			Created: t.Created,
		}
	})
	return ret, total, nil
}

// ListLatestArtifact 最新制品库列表
func ListLatestArtifact(ctx context.Context, reqDTO ListLatestArtifactReqDTO) ([]ArtifactDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	artifacts, err := artifactmd.ListLatestArtifact(ctx, reqDTO.AppId, reqDTO.Env, 10)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret := listutil.MapNe(artifacts, func(t artifactmd.Artifact) ArtifactDTO {
		return ArtifactDTO{
			Id:      t.Id,
			Name:    t.Name,
			Creator: t.Creator,
			Created: t.Created,
		}
	})
	return ret, nil
}

// DeleteArtifact 删除制品
func DeleteArtifact(ctx context.Context, reqDTO DeleteArtifactReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	artifact, app, team, err := checkAppDevelopPermByArtifactId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := artifactmd.DeleteArtifactById(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		return artifactStorage.Delete(ctx, filepath.Join(artifact.Env, artifact.AppId, artifact.Name))
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyArtifactEvent(
		reqDTO.Operator,
		team,
		app,
		artifact,
		event.AppArtifactDeleteAction,
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

func checkAppDevelopPermByArtifactId(ctx context.Context, artifactId int64, operator apisession.UserInfo) (artifactmd.Artifact, appmd.App, teammd.Team, error) {
	artifact, b, err := artifactmd.GetArtifactById(ctx, artifactId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return artifactmd.Artifact{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return artifactmd.Artifact{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	app, team, err := checkAppDevelopPermByAppId(ctx, artifact.AppId, operator)
	return artifact, app, team, err
}

func notifyArtifactEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, artifact artifactmd.Artifact, action event.AppArtifactEventAction) {
	initPsub()
	psub.Publish(event.AppArtifactTopic, event.AppArtifactEvent{
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
		ArtifactId:   artifact.Id,
		ArtifactName: artifact.Name,
		Env:          artifact.Env,
		Action:       action,
	})
}
