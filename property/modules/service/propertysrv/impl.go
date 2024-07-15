package propertysrv

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/property/modules/model/propertymd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type outerImpl struct {
}

func (*outerImpl) ListPropertySource(ctx context.Context, reqDTO ListPropertySourceReqDTO) ([]PropertySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkManagePropertySourcePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return nil, err
	}
	nodes, err := propertymd.ListEtcdNode(ctx, propertymd.ListEtcdNodeReqDTO{
		AppId: reqDTO.AppId,
		Env:   reqDTO.Env,
		Cols:  []string{"id", "name", "endpoints", "username", "password", "env"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(nodes, func(t propertymd.EtcdNode) (PropertySourceDTO, error) {
		return PropertySourceDTO{
			Id:        t.Id,
			Name:      t.Name,
			Endpoints: strings.Split(t.Endpoints, ";"),
			Username:  t.Username,
			Password:  t.Password,
			Env:       t.Env,
		}, nil
	})
}

// ListPropertySourceByFileId 配置来源
func (*outerImpl) ListPropertySourceByFileId(ctx context.Context, reqDTO ListPropertySourceByFileIdReqDTO) ([]SimplePropertySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	file, err := checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId)
	if err != nil {
		return nil, err
	}
	nodes, err := propertymd.ListEtcdNode(ctx, propertymd.ListEtcdNodeReqDTO{
		AppId: file.AppId,
		Env:   file.Env,
		Cols:  []string{"id", "name"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(nodes, func(t propertymd.EtcdNode) (SimplePropertySourceDTO, error) {
		return SimplePropertySourceDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
}

// CreatePropertySource 新增配置来源
func (*outerImpl) CreatePropertySource(ctx context.Context, reqDTO CreatePropertySourceReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err = checkManagePropertySourcePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return err
	}
	err = propertymd.InsertEtcdNode(ctx, propertymd.InsertEtcdNodeReqDTO{
		AppId:     reqDTO.AppId,
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
		Env:       reqDTO.Env,
		Name:      reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeletePropertySource(ctx context.Context, reqDTO DeletePropertySourceReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkManagePropertySourcePermBySourceId(ctx, reqDTO.Operator, reqDTO.SourceId); err != nil {
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除节点
		err2 := propertymd.DeleteEtcdNode(ctx, reqDTO.SourceId)
		if err2 != nil {
			return err2
		}
		// todo 删除历史
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// UpdatePropertySource 编辑配置来源
func (*outerImpl) UpdatePropertySource(ctx context.Context, reqDTO UpdatePropertySourceReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkManagePropertySourcePermBySourceId(ctx, reqDTO.Operator, reqDTO.SourceId); err != nil {
		return
	}
	_, err = propertymd.UpdateEtcdNode(ctx, propertymd.UpdateEtcdNodeReqDTO{
		Id:        reqDTO.SourceId,
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
		Name:      reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) CreateFile(ctx context.Context, reqDTO CreateFileReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return err
	}
	b, err := propertymd.ExistFile(ctx, reqDTO.AppId, reqDTO.Name, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		file, err2 := propertymd.InsertFile(ctx, propertymd.InsertFileReqDTO{
			AppId: reqDTO.AppId,
			Name:  reqDTO.Name,
			Env:   reqDTO.Env,
		})
		if err2 != nil {
			return err2
		}
		return propertymd.InsertHistory(ctx, propertymd.InsertHistoryReqDTO{
			FileId:  file.Id,
			Content: reqDTO.Content,
			Version: genVersion(),
			Creator: reqDTO.Operator.Account,
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) NewVersion(ctx context.Context, reqDTO NewVersionReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId)
	if err != nil {
		return err
	}
	// 检查lastVersion是否正确
	exist, err := propertymd.ExistHistoryByVersion(ctx, reqDTO.FileId, reqDTO.LastVersion)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !exist {
		return util.InvalidArgsError()
	}
	err = propertymd.InsertHistory(ctx, propertymd.InsertHistoryReqDTO{
		FileId:      reqDTO.FileId,
		Content:     reqDTO.Content,
		Version:     genVersion(),
		LastVersion: reqDTO.LastVersion,
		Creator:     reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) DeleteFile(ctx context.Context, reqDTO DeleteFileReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var file propertymd.File
	if file, err = checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId); err != nil {
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除配置文件
		_, err := propertymd.DeleteFileById(ctx, reqDTO.FileId)
		if err != nil {
			return err
		}
		// 删除配置历史
		err = propertymd.DeleteHistoryByFileId(ctx, reqDTO.FileId)
		if err != nil {
			return err
		}
		// 删除部署记录
		return propertymd.DeleteDeployByFileId(ctx, reqDTO.FileId)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	go deleteFromEtcd(file)
	return
}

// ListFile 配置文件列表
func (*outerImpl) ListFile(ctx context.Context, reqDTO ListFileReqDTO) ([]FileDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return nil, err
	}
	contents, err := propertymd.ListFile(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(contents, func(t propertymd.File) (FileDTO, error) {
		return FileDTO{
			Id:    t.Id,
			AppId: t.AppId,
			Name:  t.Name,
			Env:   t.Env,
		}, nil
	})
}

// GetHistoryByVersion 获取最新版本的配置
func (*outerImpl) GetHistoryByVersion(ctx context.Context, reqDTO GetHistoryByVersionReqDTO) (HistoryDTO, bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return HistoryDTO{}, false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	file, err := checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId)
	if err != nil {
		return HistoryDTO{}, false, err
	}
	history, b, err := propertymd.GetHistoryByVersion(ctx, reqDTO.FileId, reqDTO.Version)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return HistoryDTO{}, false, util.InternalError(err)
	}
	if !b {
		return HistoryDTO{}, false, nil
	}
	return HistoryDTO{
		Id:          history.Id,
		FileName:    file.Name,
		FileId:      history.FileId,
		Content:     history.Content,
		Version:     history.Version,
		LastVersion: history.LastVersion,
		Created:     history.Created,
		Creator:     history.Creator,
		Env:         file.Env,
	}, true, nil
}

// DeployHistory 部署配置
func (*outerImpl) DeployHistory(ctx context.Context, reqDTO DeployHistoryReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		file    propertymd.File
		history propertymd.History
	)
	history, file, err = checkAppDevelopPermByHistoryId(ctx, reqDTO.Operator, reqDTO.HistoryId)
	if err != nil {
		return
	}
	sources, err := propertymd.BatchGetEtcdNodes(ctx, reqDTO.SourceIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if len(sources) == 0 {
		return util.InvalidArgsError()
	}
	for _, source := range sources {
		if source.AppId != file.AppId || source.Env != file.Env {
			return util.InvalidArgsError()
		}
	}
	go func() {
		for _, source := range sources {
			deployToEtcd(file, history, source, reqDTO.Operator)
		}
	}()
	return nil
}

func (*outerImpl) PageHistory(ctx context.Context, reqDTO PageHistoryReqDTO) ([]HistoryDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	file, err := checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId)
	if err != nil {
		return nil, 0, err
	}
	const pageSize = 10
	histories, total, err := propertymd.PageHistory(ctx, propertymd.PageHistoryReqDTO{
		FileId:   reqDTO.FileId,
		PageNum:  reqDTO.PageNum,
		PageSize: pageSize,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret, _ := listutil.Map(histories, func(t propertymd.History) (HistoryDTO, error) {
		return HistoryDTO{
			Id:          t.Id,
			FileName:    file.Name,
			FileId:      t.FileId,
			Content:     t.Content,
			Version:     t.Version,
			LastVersion: t.LastVersion,
			Created:     t.Created,
			Creator:     t.Creator,
			Env:         file.Env,
		}, nil
	})
	return ret, total, nil
}

func (*outerImpl) ListDeploy(ctx context.Context, reqDTO ListDeployReqDTO) ([]DeployDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if _, _, err := checkAppDevelopPermByHistoryId(ctx, reqDTO.Operator, reqDTO.HistoryId); err != nil {
		return nil, err
	}
	deploys, err := propertymd.ListDeployByHistoryId(ctx, reqDTO.HistoryId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret, _ := listutil.Map(deploys, func(t propertymd.Deploy) (DeployDTO, error) {
		return DeployDTO{
			NodeName:  t.NodeName,
			Endpoints: t.Endpoints,
			Created:   t.Created,
			Creator:   t.Creator,
		}, nil
	})
	return ret, nil
}

func deleteFromEtcd(file propertymd.File) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	nodes, err := propertymd.ListEtcdNode(ctx, propertymd.ListEtcdNodeReqDTO{
		AppId: file.AppId,
		Env:   file.Env,
	})
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	for _, node := range nodes {
		etcdClient, err := newEtcdClient(node)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			continue
		}
		kv := clientv3.NewKV(etcdClient)
		_, err = kv.Delete(context.Background(), common.PropertyPrefix+file.AppId+"/"+file.Name)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		etcdClient.Close()
	}
}

type contentVal struct {
	Version string `json:"version"`
	Content string `json:"content"`
}

func deployToEtcd(file propertymd.File, history propertymd.History, node propertymd.EtcdNode, operator apisession.UserInfo) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	etcdClient, err := newEtcdClient(node)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return
	}
	defer etcdClient.Close()
	kv := clientv3.NewKV(etcdClient)
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		err2 := propertymd.InsertDeploy(ctx, propertymd.InsertDeployReqDTO{
			HistoryId: history.Id,
			NodeName:  node.Name,
			FileId:    file.Id,
			Endpoints: node.Endpoints,
			Username:  node.Username,
			Password:  node.Password,
			Creator:   operator.Account,
		})
		if err2 != nil {
			return err
		}
		jsonBytes, _ := json.Marshal(contentVal{
			Version: history.Version,
			Content: history.Content,
		})
		_, err2 = kv.Put(ctx,
			common.PropertyPrefix+file.AppId+"/"+file.Name,
			string(jsonBytes),
		)
		return err2
	})
	if err != nil {
		logger.Logger.Errorf("deploy history: %v app: %v failed with err: %v", history.Id, file.AppId, err)
	}
}

func newEtcdClient(node propertymd.EtcdNode) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(node.Endpoints, ";"),
		DialTimeout: 5 * time.Second,
		Username:    node.Username,
		Password:    node.Password,
		Logger:      zap.NewNop(),
	})
}

func genVersion() string {
	now := time.Now()
	rint := strconv.Itoa(rand.Intn(1000000))
	if len(rint) < 6 {
		rint = "000000" + rint
		rint = rint[len(rint)-6:]
	} else if len(rint) > 6 {
		rint = rint[len(rint)-6:]
	}
	return now.Format("20060102150405") + rint
}

func checkAppDevelopPermByFileId(ctx context.Context, operator apisession.UserInfo, fileId int64) (propertymd.File, error) {
	file, b, err := propertymd.GetFileById(ctx, fileId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return propertymd.File{}, util.InternalError(err)
	}
	if !b {
		return propertymd.File{}, util.InvalidArgsError()
	}
	return file, checkAppDevelopPermByAppId(ctx, operator, file.AppId)
}

func checkAppDevelopPermByHistoryId(ctx context.Context, operator apisession.UserInfo, historyId int64) (propertymd.History, propertymd.File, error) {
	history, b, err := propertymd.GetHistoryById(ctx, historyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return propertymd.History{}, propertymd.File{}, util.InternalError(err)
	}
	if !b {
		return propertymd.History{}, propertymd.File{}, util.InvalidArgsError()
	}
	file, err := checkAppDevelopPermByFileId(ctx, operator, history.FileId)
	return history, file, err
}

func checkAppDevelopPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
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
	if p.PermDetail.DevelopAppList.Contains(appId) {
		return nil
	}
	return util.UnauthorizedError()
}

func checkManagePropertySourcePermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
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
	if p.PermDetail.TeamPerm.CanManagePropertySource {
		return nil
	}
	return util.UnauthorizedError()
}

func checkManagePropertySourcePermBySourceId(ctx context.Context, operator apisession.UserInfo, sourceId int64) error {
	ps, b, err := propertymd.GetEtcdNodeById(ctx, sourceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	app, b, err := appmd.GetByAppId(ctx, ps.AppId)
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
	if p.PermDetail.TeamPerm.CanManagePropertySource {
		return nil
	}
	return util.UnauthorizedError()
}
