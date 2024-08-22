package propertysrv

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/property/modules/model/propertymd"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.AppPropertyFileTopic, func(data any) {
			req, ok := data.(event.AppPropertyFileEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppPropertyFileCreateFileAction:
							return cfg.AppPropertyFile.Create
						case event.AppPropertyFileDeleteFileAction:
							return cfg.AppPropertyFile.Delete
						default:
							return false
						}
					}
					return false
				})
			}
		})
		psub.Subscribe(event.AppPropertyVersionTopic, func(data any) {
			req, ok := data.(event.AppPropertyVersionEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppPropertyVersionNewAction:
							return cfg.AppPropertyVersion.New
						case event.AppPropertyVersionDeployAction:
							return cfg.AppPropertyVersion.Deploy
						default:
							return false
						}
					}
					return false
				})
			}
		})
		psub.Subscribe(event.AppSourceTopic, func(data any) {
			req, ok := data.(event.AppSourceEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppManagePropertySourceAction:
							return cfg.AppSource.ManagePropertySource
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

func ListPropertySource(ctx context.Context, reqDTO ListPropertySourceReqDTO) ([]PropertySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	nodes, err := propertymd.ListEtcdNode(ctx, propertymd.ListEtcdNodeReqDTO{
		Env:  reqDTO.Env,
		Cols: []string{"id", "name", "endpoints", "username", "password", "env"},
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

// ListAllPropertySource 所有配置来源
func ListAllPropertySource(ctx context.Context, reqDTO ListAllPropertySourceReqDTO) ([]SimplePropertySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	nodes, err := propertymd.ListEtcdNode(ctx, propertymd.ListEtcdNodeReqDTO{
		Env:  reqDTO.Env,
		Cols: []string{"id", "name"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return etcdNodeMd2SimpleDto(nodes)
}

func etcdNodeMd2SimpleDto(nodes []propertymd.EtcdNode) ([]SimplePropertySourceDTO, error) {
	return listutil.Map(nodes, func(t propertymd.EtcdNode) (SimplePropertySourceDTO, error) {
		return SimplePropertySourceDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
}

// ListBindPropertySource 获取绑定的配置来源
func ListBindPropertySource(ctx context.Context, reqDTO ListBindPropertySourceReqDTO) ([]SimplePropertySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	binds, err := propertymd.ListAppEtcdNodeBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if len(binds) == 0 {
		return []SimplePropertySourceDTO{}, nil
	}
	nodeIdList, _ := listutil.Map(binds, func(t propertymd.AppEtcdNodeBind) (int64, error) {
		return t.NodeId, nil
	})
	nodes, err := propertymd.BatchGetEtcdNodesById(ctx, nodeIdList, []string{"id", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return etcdNodeMd2SimpleDto(nodes)
}

// ListPropertySourceByFileId 配置来源
func ListPropertySourceByFileId(ctx context.Context, reqDTO ListPropertySourceByFileIdReqDTO) ([]SimplePropertySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	file, _, _, err := checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId)
	if err != nil {
		return nil, err
	}
	binds, err := propertymd.ListAppEtcdNodeBindByAppIdAndEnv(ctx, file.AppId, file.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if len(binds) == 0 {
		return []SimplePropertySourceDTO{}, nil
	}
	nodeIdList, _ := listutil.Map(binds, func(t propertymd.AppEtcdNodeBind) (int64, error) {
		return t.NodeId, nil
	})
	nodes, err := propertymd.BatchGetEtcdNodesById(ctx, nodeIdList, nil)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return etcdNodeMd2SimpleDto(nodes)
}

// BindAppAndPropertySource 绑定应用服务和配置来源
func BindAppAndPropertySource(ctx context.Context, reqDTO BindAppAndPropertySourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	app, team, err := checkManagePropertySourcePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return err
	}
	if len(reqDTO.SourceIdList) == 0 {
		err = propertymd.DeleteAppEtcdNodeBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		// 通知
		notifyPropertySourceEvent(
			reqDTO.Operator,
			team,
			app,
			nil,
			reqDTO.Env,
		)
		return nil
	}
	nodes, err := propertymd.BatchGetEtcdNodesById(ctx, reqDTO.SourceIdList, []string{"id", "name", "env"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if len(nodes) == 0 {
		return util.InvalidArgsError()
	}
	for _, node := range nodes {
		if node.Env != reqDTO.Env {
			return util.InvalidArgsError()
		}
	}
	insertReqs, _ := listutil.Map(nodes, func(t propertymd.EtcdNode) (propertymd.InsertAppEtcdNodeBindReqDTO, error) {
		return propertymd.InsertAppEtcdNodeBindReqDTO{
			NodeId: t.Id,
			AppId:  reqDTO.AppId,
			Env:    reqDTO.Env,
		}, nil
	})
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 先删除
		err2 := propertymd.DeleteAppEtcdNodeBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
		if err2 != nil {
			return err2
		}
		// 批量插入
		return propertymd.BatchInsertAppEtcdNodeBind(ctx, insertReqs)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 通知
	notifyPropertySourceEvent(
		reqDTO.Operator,
		team,
		app,
		nodes,
		reqDTO.Env,
	)
	return nil
}

// CreatePropertySource 新增配置来源
func CreatePropertySource(ctx context.Context, reqDTO CreatePropertySourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	err := propertymd.InsertEtcdNode(ctx, propertymd.InsertEtcdNodeReqDTO{
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
		Env:       reqDTO.Env,
		Name:      reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeletePropertySource 删除配置来源
func DeletePropertySource(ctx context.Context, reqDTO DeletePropertySourceReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除节点
		err2 := propertymd.DeleteEtcdNodeById(ctx, reqDTO.SourceId)
		if err2 != nil {
			return err2
		}
		// 删除绑定
		return propertymd.DeleteAppEtcdNodeBindByNodeId(ctx, reqDTO.SourceId)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// UpdatePropertySource 编辑配置来源
func UpdatePropertySource(ctx context.Context, reqDTO UpdatePropertySourceReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
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

// CreateFile 创建配置文件
func CreateFile(ctx context.Context, reqDTO CreateFileReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	app, team, err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
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
	var file propertymd.File
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		var err2 error
		file, err2 = propertymd.InsertFile(ctx, propertymd.InsertFileReqDTO{
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
	notifyPropertyFileEvent(
		reqDTO.Operator,
		team,
		app,
		file,
		event.AppPropertyFileCreateFileAction,
	)
	return nil
}

// NewVersion 创建配置文件新版本
func NewVersion(ctx context.Context, reqDTO NewVersionReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	file, app, team, err := checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId)
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
	version := genVersion()
	err = propertymd.InsertHistory(ctx, propertymd.InsertHistoryReqDTO{
		FileId:      reqDTO.FileId,
		Content:     reqDTO.Content,
		Version:     version,
		LastVersion: reqDTO.LastVersion,
		Creator:     reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyPropertyVersionEvent(
		reqDTO.Operator,
		team,
		app,
		file,
		event.AppPropertyVersionNewAction,
		version,
		reqDTO.LastVersion,
	)
	return nil
}

// DeleteFile 删除配置文件
func DeleteFile(ctx context.Context, reqDTO DeleteFileReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	file, app, team, err := checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除配置文件
		_, err2 := propertymd.DeleteFileById(ctx, reqDTO.FileId)
		if err2 != nil {
			return err2
		}
		// 删除配置历史
		err2 = propertymd.DeleteHistoryByFileId(ctx, reqDTO.FileId)
		if err2 != nil {
			return err2
		}
		// 删除部署记录
		return propertymd.DeleteDeployByFileId(ctx, reqDTO.FileId)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	go deleteFromEtcd(file)
	notifyPropertyFileEvent(
		reqDTO.Operator,
		team,
		app,
		file,
		event.AppPropertyFileDeleteFileAction,
	)
	return nil
}

// ListFile 配置文件列表
func ListFile(ctx context.Context, reqDTO ListFileReqDTO) ([]FileDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
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
func GetHistoryByVersion(ctx context.Context, reqDTO GetHistoryByVersionReqDTO) (HistoryDTO, bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return HistoryDTO{}, false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	file, _, _, err := checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId)
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
func DeployHistory(ctx context.Context, reqDTO DeployHistoryReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	history, file, app, team, err := checkPropertyDeployPermByHistoryId(ctx, reqDTO.Operator, reqDTO.HistoryId)
	if err != nil {
		return err
	}
	// 校验发布节点
	binds, err := propertymd.BatchGetAppEtcdNodeBindByNodeIdListAndAppId(ctx, reqDTO.SourceIdList, file.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if len(binds) == 0 {
		return util.InvalidArgsError()
	}
	nodeIdList, _ := listutil.Map(binds, func(t propertymd.AppEtcdNodeBind) (int64, error) {
		return t.NodeId, nil
	})
	nodes, err := propertymd.BatchGetEtcdNodesById(ctx, nodeIdList, nil)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if len(nodes) == 0 {
		return util.InvalidArgsError()
	}
	for _, node := range nodes {
		if node.Env != file.Env {
			return util.InvalidArgsError()
		}
	}
	go func() {
		for _, source := range nodes {
			deployToEtcd(file, history, source, reqDTO.Operator.Account)
		}
	}()
	notifyPropertyVersionEvent(
		reqDTO.Operator,
		team,
		app,
		file,
		event.AppPropertyVersionDeployAction,
		history.Version,
		history.LastVersion,
	)
	return nil
}

func ListHistory(ctx context.Context, reqDTO ListHistoryReqDTO) ([]HistoryDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	file, _, _, err := checkAppDevelopPermByFileId(ctx, reqDTO.Operator, reqDTO.FileId)
	if err != nil {
		return nil, 0, err
	}
	histories, total, err := propertymd.ListHistory(ctx, propertymd.ListHistoryReqDTO{
		FileId:   reqDTO.FileId,
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
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

func ListDeploy(ctx context.Context, reqDTO ListDeployReqDTO) ([]DeployDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByHistoryId(ctx, reqDTO.Operator, reqDTO.HistoryId)
	if err != nil {
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
	binds, err := propertymd.ListAppEtcdNodeBindByAppIdAndEnv(ctx, file.AppId, file.Env)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	nodeIdList, _ := listutil.Map(binds, func(t propertymd.AppEtcdNodeBind) (int64, error) {
		return t.NodeId, nil
	})
	nodes, err := propertymd.BatchGetEtcdNodesById(ctx, nodeIdList, nil)
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

func deployToEtcd(file propertymd.File, history propertymd.History, node propertymd.EtcdNode, creator string) {
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
		jsonBytes, _ := json.Marshal(contentVal{
			Version: history.Version,
			Content: history.Content,
		})
		_, err2 := kv.Put(ctx,
			common.PropertyPrefix+file.AppId+"/"+file.Name,
			string(jsonBytes),
		)
		if err2 != nil {
			return err
		}
		return propertymd.InsertDeploy(ctx, propertymd.InsertDeployReqDTO{
			HistoryId: history.Id,
			NodeName:  node.Name,
			FileId:    file.Id,
			AppId:     file.AppId,
			Endpoints: node.Endpoints,
			Username:  node.Username,
			Password:  node.Password,
			Creator:   creator,
		})
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

func checkAppDevelopPermByFileId(ctx context.Context, operator apisession.UserInfo, fileId int64) (propertymd.File, appmd.App, teammd.Team, error) {
	file, b, err := propertymd.GetFileById(ctx, fileId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return propertymd.File{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return propertymd.File{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	app, team, err := checkAppDevelopPermByAppId(ctx, operator, file.AppId)
	return file, app, team, err
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
	file, _, _, err := checkAppDevelopPermByFileId(ctx, operator, history.FileId)
	return history, file, err
}

func checkAppDevelopPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (appmd.App, teammd.Team, error) {
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
		return appmd.App{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return app, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
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

func checkManagePropertySourcePermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (appmd.App, teammd.Team, error) {
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
		logger.Logger.WithContext(ctx).Error(err)
		return app, team, util.InternalError(err)
	}
	if !b {
		return app, team, util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.GetAppPerm(appId).CanManagePropertySource {
		return app, team, nil
	}
	return app, team, util.UnauthorizedError()
}

func checkPropertyDeployPermByHistoryId(ctx context.Context, operator apisession.UserInfo, historyId int64) (propertymd.History, propertymd.File, appmd.App, teammd.Team, error) {
	history, b, err := propertymd.GetHistoryById(ctx, historyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return propertymd.History{}, propertymd.File{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return propertymd.History{}, propertymd.File{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	file, b, err := propertymd.GetFileById(ctx, history.FileId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return propertymd.History{}, propertymd.File{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return propertymd.History{}, propertymd.File{}, appmd.App{}, teammd.Team{}, util.ThereHasBugErr()
	}
	app, b, err := appmd.GetByAppId(ctx, file.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return propertymd.History{}, propertymd.File{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return propertymd.History{}, propertymd.File{}, appmd.App{}, teammd.Team{}, util.ThereHasBugErr()
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return propertymd.History{}, propertymd.File{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return propertymd.History{}, propertymd.File{}, appmd.App{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return history, file, app, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return history, file, app, team, util.InternalError(err)
	}
	if !b {
		return history, file, app, team, util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.GetAppPerm(file.AppId).CanDeployProperty {
		return history, file, app, team, nil
	}
	return history, file, app, team, util.UnauthorizedError()
}

func notifyPropertyFileEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, file propertymd.File, action event.AppPropertyFileEventAction) {
	initPsub()
	psub.Publish(event.AppPropertyFileTopic, event.AppPropertyFileEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Name,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		BasePropertyFile: event.BasePropertyFile{
			FileId:   file.Id,
			FileName: file.Name,
			Env:      file.Env,
		},
		Action: action,
	})
}

func notifyPropertySourceEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, sources []propertymd.EtcdNode, env string) {
	initPsub()
	srcs, _ := listutil.Map(sources, func(t propertymd.EtcdNode) (event.AppSource, error) {
		return event.AppSource{
			Id:   t.Id,
			Name: t.Name,
			Env:  t.Env,
		}, nil
	})
	psub.Publish(event.AppSourceTopic, event.AppSourceEvent{
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
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, event.AppManagePropertySourceAction.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, event.AppManagePropertySourceAction.GetI18nValue()),
		},
		Action:  event.AppManagePropertySourceAction,
		Env:     env,
		Sources: srcs,
	})
}

func notifyPropertyVersionEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, file propertymd.File, action event.AppPropertyVersionEventAction, version, oldVersion string) {
	initPsub()
	psub.Publish(event.AppPropertyVersionTopic, event.AppPropertyVersionEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Name,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		BasePropertyFile: event.BasePropertyFile{
			FileId:   file.Id,
			FileName: file.Name,
			Env:      file.Env,
		},
		Version: version,
		Action:  action,
	})
}
