package discoverysrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/discovery/modules/model/discoverymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/actuator"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/services/lb"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// ListDiscoverySource 注册中心来源
func ListDiscoverySource(ctx context.Context, reqDTO ListDiscoverySourceReqDTO) ([]DiscoverySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	nodes, err := discoverymd.ListEtcdNode(ctx, discoverymd.ListEtcdNodeReqDTO{
		Env:  reqDTO.Env,
		Cols: []string{"id", "name", "endpoints", "username", "password", "env"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(nodes, func(t discoverymd.EtcdNode) (DiscoverySourceDTO, error) {
		return DiscoverySourceDTO{
			Id:        t.Id,
			Name:      t.Name,
			Endpoints: strings.Split(t.Endpoints, ";"),
			Username:  t.Username,
			Password:  t.Password,
			Env:       t.Env,
		}, nil
	})
}

// ListAllDiscoverySource 获取所有注册中心来源
func ListAllDiscoverySource(ctx context.Context, reqDTO ListAllDiscoverySourceReqDTO) ([]SimpleDiscoverySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	nodes, err := discoverymd.ListEtcdNode(ctx, discoverymd.ListEtcdNodeReqDTO{
		Env:  reqDTO.Env,
		Cols: []string{"id", "name"},
	})
	if err != nil {
		return nil, err
	}
	return listutil.Map(nodes, func(t discoverymd.EtcdNode) (SimpleDiscoverySourceDTO, error) {
		return SimpleDiscoverySourceDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
}

// BindAppAndDiscoverySource 绑定注册中心来源
func BindAppAndDiscoverySource(ctx context.Context, reqDTO BindAppAndDiscoverySourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkManageDiscoverySourcePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return err
	}
	if len(reqDTO.SourceIdList) == 0 {
		err = discoverymd.DeleteAppEtcdNodeBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		return nil
	}
	nodes, err := discoverymd.BatchGetEtcdNodeByIdList(ctx, reqDTO.SourceIdList, []string{"id", "env"})
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
	insertReqs, _ := listutil.Map(nodes, func(t discoverymd.EtcdNode) (discoverymd.InsertAppEtcdNodeBindReqDTO, error) {
		return discoverymd.InsertAppEtcdNodeBindReqDTO{
			NodeId: t.Id,
			AppId:  reqDTO.AppId,
			Env:    reqDTO.Env,
		}, nil
	})
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 先删除
		err2 := discoverymd.DeleteAppEtcdNodeBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
		if err2 != nil {
			return err2
		}
		// 批量插入
		return discoverymd.BatchInsertAppEtcdNodeBind(ctx, insertReqs)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListBindDiscoverySource 获取绑定注册中心来源
func ListBindDiscoverySource(ctx context.Context, reqDTO ListBindDiscoverySourceReqDTO) ([]SimpleBindDiscoverySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return nil, err
	}
	binds, err := discoverymd.ListAppEtcdNodeBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		return nil, err
	}
	if len(binds) == 0 {
		return []SimpleBindDiscoverySourceDTO{}, nil
	}
	bindMap := make(map[int64]discoverymd.AppEtcdNodeBind, len(binds))
	nodeIdList, _ := listutil.Map(binds, func(t discoverymd.AppEtcdNodeBind) (int64, error) {
		bindMap[t.NodeId] = t
		return t.NodeId, nil
	})
	nodes, err := discoverymd.BatchGetEtcdNodeByIdList(ctx, nodeIdList, []string{"id", "name"})
	return listutil.Map(nodes, func(t discoverymd.EtcdNode) (SimpleBindDiscoverySourceDTO, error) {
		bind := bindMap[t.Id]
		return SimpleBindDiscoverySourceDTO{
			Id:     t.Id,
			Name:   t.Name,
			BindId: bind.Id,
			Env:    bind.Env,
		}, nil
	})
}

// CreateDiscoverySource 创建注册中心来源
func CreateDiscoverySource(ctx context.Context, reqDTO CreateDiscoverySourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	err := discoverymd.InsertEtcdNode(ctx, discoverymd.InsertEtcdNodeReqDTO{
		Name:      reqDTO.Name,
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
		Env:       reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

// DeleteDiscoverySource 删除注册中心来源
func DeleteDiscoverySource(ctx context.Context, reqDTO DeleteDiscoverySourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.InvalidArgsError()
	}
	err := discoverymd.DeleteEtcdNodeById(ctx, reqDTO.SourceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateDiscoverySource 编辑注册中心来源
func UpdateDiscoverySource(ctx context.Context, reqDTO UpdateDiscoverySourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.InvalidArgsError()
	}
	_, err := discoverymd.UpdateEtcdNode(ctx, discoverymd.UpdateEtcdNodeReqDTO{
		Id:        reqDTO.SourceId,
		Name:      reqDTO.Name,
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListDiscoveryService 服务列表
func ListDiscoveryService(ctx context.Context, reqDTO ListDiscoveryServiceReqDTO) ([]ServiceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	appId, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return nil, err
	}
	client, err := newEtcdClient(source)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	defer client.Close()
	kv := clientv3.NewKV(client)
	prefix := common.ServicePrefix + appId + "-http/"
	response, err := kv.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret := make([]ServiceDTO, 0, len(response.Kvs))
	for _, kv := range response.Kvs {
		var s lb.Server
		err = json.Unmarshal(kv.Value, &s)
		if err == nil {
			ret = append(ret, ServiceDTO{
				Server:     s,
				Up:         true,
				InstanceId: strings.TrimPrefix(string(kv.Key), prefix),
			})
		}
	}
	services, err := discoverymd.ListDownServiceBySourceIdAndAppId(ctx, source.Id, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	for _, service := range services {
		if service.DownService == nil {
			continue
		}
		ret = append(ret, ServiceDTO{
			Server:     service.DownService.Data,
			Up:         false,
			InstanceId: service.InstanceId,
		})
	}
	return ret, nil
}

// DeregisterService 下线服务
func DeregisterService(ctx context.Context, reqDTO DeregisterServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	appId, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return err
	}
	_, b, err := discoverymd.GetDownServiceBySourceIdAndInstanceId(ctx, source.Id, reqDTO.InstanceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	client, err := newEtcdClient(source)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	defer client.Close()
	kv := clientv3.NewKV(client)
	response, err := kv.Get(ctx, common.ServicePrefix+appId+"-http/"+reqDTO.InstanceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if len(response.Kvs) == 0 || len(response.Kvs) > 1 {
		return util.OperationFailedError()
	}
	var s lb.Server
	err = json.Unmarshal(response.Kvs[0].Value, &s)
	if err != nil {
		return util.ThereHasBugErr()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		err2 := discoverymd.InsertDownService(ctx, discoverymd.InsertDownServiceDTO{
			SourceId:   source.Id,
			AppId:      appId,
			Service:    s,
			InstanceId: reqDTO.InstanceId,
		})
		if err2 != nil {
			return err2
		}
		resp, err2 := http.Get(fmt.Sprintf("http://%s:%d/actuator/v1/deregisterServer", s.Host, actuator.DefaultServerPort))
		if err2 != nil {
			return err2
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return errors.New("failed")
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.OperationFailedError()
	}
	time.Sleep(time.Second)
	return nil
}

// ReRegisterService 上线服务
func ReRegisterService(ctx context.Context, reqDTO ReRegisterServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return err
	}
	service, b, err := discoverymd.GetDownServiceBySourceIdAndInstanceId(ctx, source.Id, reqDTO.InstanceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || service.DownService == nil {
		return util.InvalidArgsError()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := discoverymd.DeleteDownServiceById(ctx, service.Id)
		if err2 != nil {
			return err2
		}
		resp, err2 := http.Get(fmt.Sprintf("http://%s:%d/actuator/v1/registerServer", service.DownService.Data.Host, actuator.DefaultServerPort))
		if err2 != nil {
			return err2
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return errors.New("failed")
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.OperationFailedError()
	}
	time.Sleep(time.Second)
	return nil
}

// DeleteDownService 删除下线服务
func DeleteDownService(ctx context.Context, reqDTO DeleteDownServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return err
	}
	service, b, err := discoverymd.GetDownServiceBySourceIdAndInstanceId(ctx, source.Id, reqDTO.InstanceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || service.DownService == nil {
		return util.InvalidArgsError()
	}
	_, err = discoverymd.DeleteDownServiceById(ctx, service.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func newEtcdClient(node discoverymd.EtcdNode) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(node.Endpoints, ";"),
		DialTimeout: 5 * time.Second,
		Username:    node.Username,
		Password:    node.Password,
		Logger:      zap.NewNop(),
	})
}

func checkAppDevelopPermByBindId(ctx context.Context, operator apisession.UserInfo, bindId int64) (string, discoverymd.EtcdNode, error) {
	bind, b, err := discoverymd.GetAppEtcdNodeBindById(ctx, bindId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", discoverymd.EtcdNode{}, util.InternalError(err)
	}
	if !b {
		return "", discoverymd.EtcdNode{}, util.InvalidArgsError()
	}
	node, b, err := discoverymd.GetEtcdNodeById(ctx, bind.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", discoverymd.EtcdNode{}, util.InternalError(err)
	}
	if !b {
		return "", discoverymd.EtcdNode{}, util.InvalidArgsError()
	}
	return bind.AppId, node, checkAppDevelopPermByAppId(ctx, operator, bind.AppId)
}

func checkManageDiscoverySourcePermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
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
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	if p.PermDetail.GetAppPerm(appId).CanManageDiscoverySource {
		return nil
	}
	return util.UnauthorizedError()
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
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	if p.PermDetail.GetAppPerm(appId).CanDevelop {
		return nil
	}
	return util.UnauthorizedError()
}
