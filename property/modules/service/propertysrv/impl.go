package propertysrv

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/property/modules/model/propertymd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/strutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"go.etcd.io/etcd/api/v3/authpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type outerImpl struct {
}

func (*outerImpl) ListSimpleEtcdNode(ctx context.Context, env string) ([]string, error) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	nodes, err := propertymd.ListEtcdNode(ctx, env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(nodes, func(t propertymd.EtcdNode) (string, error) {
		return t.NodeId, nil
	})
}

func (*outerImpl) ListEtcdNode(ctx context.Context, reqDTO ListEtcdNodeReqDTO) ([]EtcdNodeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	// 系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	nodes, err := propertymd.ListEtcdNode(ctx, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(nodes, func(t propertymd.EtcdNode) (EtcdNodeDTO, error) {
		return EtcdNodeDTO{
			NodeId:    t.NodeId,
			Endpoints: strings.Split(t.Endpoints, ";"),
			Username:  t.Username,
			Password:  t.Password,
		}, nil
	})
}

func (*outerImpl) InsertEtcdNode(ctx context.Context, reqDTO InsertEtcdNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.InsertEtcdNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := propertymd.GetEtcdNodeByNodeId(ctx, reqDTO.NodeId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.AlreadyExistsError()
		return
	}
	err = propertymd.InsertEtcdNode(ctx, propertymd.InsertEtcdNodeReqDTO{
		NodeId:    reqDTO.NodeId,
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
		Env:       reqDTO.Env,
	})
	if err != nil {
		if xormutil.IsDuplicatedEntryError(err) {
			err = util.AlreadyExistsError()
			return
		}
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteEtcdNode(ctx context.Context, reqDTO DeleteEtcdNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.DeleteEtcdNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err = propertymd.DeleteEtcdNode(ctx, reqDTO.NodeId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) UpdateEtcdNode(ctx context.Context, reqDTO UpdateEtcdNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.UpdateEtcdNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = propertymd.UpdateEtcdNode(ctx, propertymd.UpdateEtcdNodeReqDTO{
		NodeId:    reqDTO.NodeId,
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
		Env:       reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) GrantAuth(ctx context.Context, reqDTO GrantAuthReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.GrantAuth),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return
	}
	if err = grantAuth(ctx, reqDTO.AppId, reqDTO.Env); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) GetAuth(ctx context.Context, reqDTO GetAuthReqDTO) (string, string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", "", err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return "", "", err
	}
	auth, b, err := propertymd.GetAuthByAppId(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", "", util.InternalError(err)
	}
	if !b {
		return "", "", util.InvalidArgsError()
	}
	return auth.Username, auth.Password, nil
}

func (*outerImpl) CreateFile(ctx context.Context, reqDTO CreateFileReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
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
			FileId:      file.Id,
			Content:     reqDTO.Content,
			Version:     genVersion(),
			LastVersion: "",
			Env:         reqDTO.Env,
			Creator:     reqDTO.Operator.Account,
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
	file, err := checkPerm(ctx, reqDTO.Operator, reqDTO.FileId)
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
		Env:         file.Env,
		Creator:     reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) DeletePropContent(ctx context.Context, reqDTO DeletePropContentReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.DeletePropContent),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var content propertymd.File
	if content, err = checkPerm(ctx, reqDTO.Operator, reqDTO.Id); err != nil {
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除配置文件
		_, err := propertymd.DeletePropContent(ctx, reqDTO.Id, reqDTO.Env)
		if err != nil {
			return err
		}
		// 删除配置历史
		err = propertymd.DeleteHistory(ctx, reqDTO.Id, reqDTO.Env)
		if err != nil {
			return err
		}
		// 删除部署记录
		return propertymd.DeleteDeploy(ctx, reqDTO.Id, reqDTO.Env)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	go deleteFromEtcd(content, reqDTO.Env)
	return
}

// ListFile 配置文件列表
func (*outerImpl) ListFile(ctx context.Context, reqDTO ListFileReqDTO) ([]FileDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
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
	if _, err := checkPerm(ctx, reqDTO.Operator, reqDTO.FileId); err != nil {
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
		FileId:  history.FileId,
		Content: history.Content,
		Version: history.Version,
		Created: history.Created,
		Creator: history.Creator,
	}, true, nil
}

func (*outerImpl) DeployPropContent(ctx context.Context, reqDTO DeployPropContentReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.DeployPropContent),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var content propertymd.File
	if content, err = checkPerm(ctx, reqDTO.Operator, reqDTO.Id); err != nil {
		return
	}
	nodes, err := propertymd.BatchGetEtcdNodes(ctx, reqDTO.EtcdNodeList, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if len(nodes) == 0 {
		return util.InvalidArgsError()
	}
	// 获取历史版本记录
	var (
		history propertymd.History
		b       bool
	)
	history, b, err = propertymd.GetHistoryByVersion(ctx, reqDTO.Id, reqDTO.Version)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	go func() {
		for _, node := range nodes {
			deployToEtcd(reqDTO.Id, content.AppId, content.Name, history.Content, history.Version, node, reqDTO.Env, reqDTO.Operator.Account)
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
	if _, err := checkPerm(ctx, reqDTO.Operator, reqDTO.FileId); err != nil {
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
			FileId:      t.FileId,
			Content:     t.Content,
			Version:     t.Version,
			Created:     t.Created,
			Creator:     t.Creator,
			LastVersion: t.LastVersion,
		}, nil
	})
	return ret, total, nil
}

func (*outerImpl) ListDeploy(ctx context.Context, reqDTO ListDeployReqDTO) ([]DeployDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if _, err := checkPerm(ctx, reqDTO.Operator, reqDTO.ContentId); err != nil {
		return nil, 0, err
	}
	deploys, err := propertymd.ListDeploy(ctx, propertymd.ListDeployReqDTO{
		ContentId: reqDTO.ContentId,
		Version:   reqDTO.Version,
		Cursor:    reqDTO.Cursor,
		Limit:     reqDTO.Limit,
		NodeId:    reqDTO.NodeId,
		Env:       reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var cursor int64 = 0
	if len(deploys) == reqDTO.Limit {
		cursor = deploys[len(deploys)-1].Id
	}
	ret, _ := listutil.Map(deploys, func(t propertymd.PropDeploy) (DeployDTO, error) {
		return DeployDTO{
			ContentId: t.ContentId,
			Content:   t.Content,
			Version:   t.Version,
			NodeId:    t.NodeId,
			Created:   t.Created,
			Creator:   t.Creator,
		}, nil
	})
	return ret, cursor, nil
}

func deleteFromEtcd(content propertymd.File, env string) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	nodes, err := propertymd.ListEtcdNode(ctx, env)
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
		_, err = kv.Delete(context.Background(), common.PropertyPrefix+content.AppId+"/"+content.Name)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		etcdClient.Close()
	}
}

func grantAuthToEtcd(auth propertymd.EtcdAuth, env string) {
	ctx := context.Background()
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	nodes, err := propertymd.ListEtcdNode(ctx, env)
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
		authClient := clientv3.NewAuth(etcdClient)
		_, err = authClient.UserAdd(
			ctx,
			auth.Username,
			auth.Password,
		)
		if err != nil && !strings.Contains(err.Error(), "user name already exists") {
			logger.Logger.WithContext(ctx).Error(err)
		} else {
			roleName := "prop_" + auth.AppId + "_read"
			roleGet, err := authClient.RoleGet(ctx, roleName)
			if err != nil {
				if strings.Contains(err.Error(), "role name not found") {
					_, err = authClient.RoleAdd(ctx, roleName)
					if err != nil {
						logger.Logger.WithContext(ctx).Error(err)
						continue
					}
				} else {
					logger.Logger.WithContext(ctx).Error(err)
					continue
				}
			}
			key := common.PropertyPrefix + auth.AppId + "/"
			rangeEnd := common.PropertyPrefix + auth.AppId + "0"
			if roleGet == nil || !permHasTargetKey(roleGet.Perm, key, rangeEnd) {
				if _, err = authClient.RoleGrantPermission(
					ctx,
					roleName,
					key,
					rangeEnd,
					clientv3.PermissionType(clientv3.PermRead),
				); err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					continue
				}
			}
			if _, err = authClient.UserGrantRole(ctx, auth.Username, roleName); err != nil {
				logger.Logger.WithContext(ctx).Error(err)
			}
		}
		etcdClient.Close()
	}
}

func permHasTargetKey(perms []*authpb.Permission, key, rangeEnd string) bool {
	for _, perm := range perms {
		if string(perm.Key) == key && string(perm.RangeEnd) == rangeEnd {
			return true
		}
	}
	return false
}

type contentVal struct {
	Version string `json:"version"`
	Content string `json:"content"`
}

func deployToEtcd(id int64, appId, name, content, version string, node propertymd.EtcdNode, env string, creator string) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	etcdClient, err := newEtcdClient(node)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return
	}
	defer etcdClient.Close()
	kv := clientv3.NewKV(etcdClient)
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		err := propertymd.InsertDeploy(ctx, propertymd.InsertDeployReqDTO{
			ContentId:    id,
			Content:      content,
			Version:      version,
			NodeId:       node.NodeId,
			ContentAppId: appId,
			ContentName:  name,
			Endpoints:    node.Endpoints,
			Username:     node.Username,
			Password:     node.Password,
			Creator:      creator,
			Env:          env,
		})
		if err != nil {
			return err
		}
		jsonBytes, _ := json.Marshal(contentVal{
			Version: version,
			Content: content,
		})
		_, err = kv.Put(ctx,
			common.PropertyPrefix+appId+"/"+name,
			string(jsonBytes),
		)
		return err
	})
	if err != nil {
		logger.Logger.Error(err)
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

func checkPerm(ctx context.Context, operator apisession.UserInfo, id int64) (propertymd.File, error) {
	file, b, err := propertymd.GetFileById(ctx, id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return propertymd.File{}, util.InternalError(err)
	}
	if !b {
		return propertymd.File{}, util.InvalidArgsError()
	}
	return file, checkPermByAppId(ctx, operator, file.AppId)
}

func checkPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
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

type innerImpl struct{}

func (*innerImpl) GrantAuth(ctx context.Context, appId, env string) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := grantAuth(ctx, appId, env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
}

func (*innerImpl) CheckConsistent(env string) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	// 检查数据库->etcd
	nodes, err := propertymd.ListEtcdNode(ctx, env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return
	}
	clientMap := make(map[string]*clientv3.Client, 8)
	defer func() {
		for _, client := range clientMap {
			client.Close()
		}
	}()
	if err = propertymd.IteratePropContent(ctx, env, func(content *propertymd.File) error {
		return checkDb2EtcdConsistent(content, nodes, clientMap, env)
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	if err = checkEtcd2DbConsistent(nodes, clientMap, env); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
}

// checkEtcd2DbConsistent etcd -> db数据一致性
func checkEtcd2DbConsistent(nodes []propertymd.EtcdNode, clientMap map[string]*clientv3.Client, env string) error {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	for _, node := range nodes {
		if err := propertymd.IterateDeletedDeployByNodeId(ctx, node.NodeId, env, func(deploy *propertymd.PropDeploy) error {
			client := clientMap[node.NodeId]
			var err error
			if client == nil {
				client, err = newEtcdClient(node)
				if err != nil {
					return err
				}
				clientMap[node.NodeId] = client
			}
			kv := clientv3.NewKV(client)
			key := common.PropertyPrefix + deploy.ContentAppId + "/" + deploy.ContentName
			response, err := kv.Get(
				ctx,
				key,
			)
			if err != nil {
				logger.Logger.Error(err)
			} else {
				// 如果版本号相同 则删除
				if response.Count > 0 && checkConsistentVersion(response.Kvs[0].Value, deploy.ContentAppId, deploy.ContentName, deploy.Version) {
					logger.Logger.Infof("find db not exists but etcd exists, delete key: %s", key)
					_, err = kv.Delete(ctx, key)
					if err != nil {
						logger.Logger.Error(err)
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

// checkDb2EtcdConsistent db -> etcd数据一致性
func checkDb2EtcdConsistent(content *propertymd.File, nodes []propertymd.EtcdNode, clientMap map[string]*clientv3.Client, env string) error {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	for _, node := range nodes {
		deploy, b, err := propertymd.GetLatestDeployByNodeId(ctx, content.Id, node.NodeId, env)
		if err != nil {
			return err
		}
		// 10秒内的部署忽略
		if time.Since(deploy.Created) < 10*time.Second {
			continue
		}
		if b {
			client := clientMap[node.NodeId]
			if client == nil {
				client, err = newEtcdClient(node)
				if err != nil {
					return err
				}
				clientMap[node.NodeId] = client
			}
			kv := clientv3.NewKV(client)
			key := common.PropertyPrefix + content.AppId + "/" + content.Name
			response, err := kv.Get(
				ctx,
				key,
			)
			if err != nil {
				logger.Logger.Error(err)
			} else {
				if response.Count == 0 || !checkConsistentVersion(response.Kvs[0].Value, content.AppId, content.Name, deploy.Version) {
					logger.Logger.Infof("find db exists but etcd not exists: %s version: %s", key, deploy.Version)
					jsonBytes, _ := json.Marshal(contentVal{
						Version: deploy.Version,
						Content: deploy.Content,
					})
					_, err = kv.Put(ctx, key, string(jsonBytes))
					if err != nil {
						logger.Logger.Error(err)
					}
				}
			}
		}
	}
	return nil
}

func checkConsistentVersion(val []byte, appId, name, version string) bool {
	var ret contentVal
	err := json.Unmarshal(val, &ret)
	if err != nil {
		logger.Logger.Errorf("read value from etcd is not json format: %s %s", appId, name)
		return false
	}
	return ret.Version == version
}

func grantAuth(ctx context.Context, appId, env string) error {
	auth, b, err := propertymd.GetAuthByAppId(ctx, appId, env)
	if err != nil {
		return err
	}
	if !b {
		insertReq := propertymd.InsertAuthReqDTO{
			AppId:    appId,
			Username: "prop_" + appId,
			Password: strutil.RandomStr(16),
			Env:      env,
		}
		err = propertymd.InsertAuth(ctx, insertReq)
		if err != nil {
			return err
		}
		auth = propertymd.EtcdAuth{
			AppId:    insertReq.AppId,
			Username: insertReq.Username,
			Password: insertReq.Password,
		}
	}
	go grantAuthToEtcd(auth, env)
	return nil
}
