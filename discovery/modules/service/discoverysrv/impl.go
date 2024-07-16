package discoverysrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/discovery/modules/model/discoverymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
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

type outImpl struct{}

// ListDiscoverySource 注册中心来源
func (*outImpl) ListDiscoverySource(ctx context.Context, reqDTO ListDiscoverySourceReqDTO) ([]DiscoverySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkManageDiscoverySourcePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return nil, err
	}
	nodes, err := discoverymd.ListEtcdNode(ctx, discoverymd.ListEtcdNodeReqDTO{
		AppId: reqDTO.AppId,
		Env:   reqDTO.Env,
		Cols:  []string{"id", "name", "endpoints", "username", "password", "env"},
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

// ListSimpleDiscoverySource 注册中心来源
func (*outImpl) ListSimpleDiscoverySource(ctx context.Context, reqDTO ListDiscoverySourceReqDTO) ([]SimpleDiscoverySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return nil, err
	}
	nodes, err := discoverymd.ListEtcdNode(ctx, discoverymd.ListEtcdNodeReqDTO{
		AppId: reqDTO.AppId,
		Env:   reqDTO.Env,
		Cols:  []string{"id", "name"},
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

// CreateDiscoverySource 创建注册中心来源
func (*outImpl) CreateDiscoverySource(ctx context.Context, reqDTO CreateDiscoverySourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkManageDiscoverySourcePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return err
	}
	err = discoverymd.InsertEtcdNode(ctx, discoverymd.InsertEtcdNodeReqDTO{
		AppId:     reqDTO.AppId,
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
func (*outImpl) DeleteDiscoverySource(ctx context.Context, reqDTO DeleteDiscoverySourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkManageDiscoverySourcePermBySourceId(ctx, reqDTO.Operator, reqDTO.SourceId)
	if err != nil {
		return err
	}
	err = discoverymd.DeleteEtcdNodeById(ctx, reqDTO.SourceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

// UpdateDiscoverySource 编辑注册中心来源
func (*outImpl) UpdateDiscoverySource(ctx context.Context, reqDTO UpdateDiscoverySourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkManageDiscoverySourcePermBySourceId(ctx, reqDTO.Operator, reqDTO.SourceId)
	if err != nil {
		return err
	}
	_, err = discoverymd.UpdateEtcdNode(ctx, discoverymd.UpdateEtcdNodeReqDTO{
		Id:        reqDTO.SourceId,
		Name:      reqDTO.Name,
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

// ListDiscoveryService 服务列表
func (*outImpl) ListDiscoveryService(ctx context.Context, reqDTO ListDiscoveryServiceReqDTO) ([]ServiceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	source, err := checkAppDevelopPermBySourceId(ctx, reqDTO.Operator, reqDTO.SourceId)
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
	prefix := common.ServicePrefix + source.AppId + "-http" + "/"
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
	services, err := discoverymd.ListDownServiceBySourceId(ctx, reqDTO.SourceId)
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
func (*outImpl) DeregisterService(ctx context.Context, reqDTO DeregisterServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	source, err := checkAppDevelopPermBySourceId(ctx, reqDTO.Operator, reqDTO.SourceId)
	if err != nil {
		return err
	}
	_, b, err := discoverymd.GetDownServiceBySourceIdAndInstanceId(ctx, reqDTO.SourceId, reqDTO.InstanceId)
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
	response, err := kv.Get(ctx, common.ServicePrefix+source.AppId+"-http"+"/"+reqDTO.InstanceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if len(response.Kvs) == 0 || len(response.Kvs) > 1 {
		return util.InvalidArgsError()
	}
	var s lb.Server
	err = json.Unmarshal(response.Kvs[0].Value, &s)
	if err != nil {
		return util.ThereHasBugErr()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		err2 := discoverymd.InsertDownService(ctx, discoverymd.InsertDownServiceDTO{
			AppId:      source.AppId,
			SourceId:   source.Id,
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
func (*outImpl) ReRegisterService(ctx context.Context, reqDTO ReRegisterServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkAppDevelopPermBySourceId(ctx, reqDTO.Operator, reqDTO.SourceId)
	if err != nil {
		return err
	}
	service, b, err := discoverymd.GetDownServiceBySourceIdAndInstanceId(ctx, reqDTO.SourceId, reqDTO.InstanceId)
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
func (*outImpl) DeleteDownService(ctx context.Context, reqDTO DeleteDownServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkAppDevelopPermBySourceId(ctx, reqDTO.Operator, reqDTO.SourceId)
	if err != nil {
		return err
	}
	service, b, err := discoverymd.GetDownServiceBySourceIdAndInstanceId(ctx, reqDTO.SourceId, reqDTO.InstanceId)
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

func checkManageDiscoverySourcePermBySourceId(ctx context.Context, operator apisession.UserInfo, sourceId int64) (discoverymd.EtcdNode, error) {
	node, b, err := discoverymd.GetEtcdNodeById(ctx, sourceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return discoverymd.EtcdNode{}, util.InternalError(err)
	}
	if !b {
		return discoverymd.EtcdNode{}, util.InvalidArgsError()
	}
	return node, checkManageDiscoverySourcePermByAppId(ctx, operator, node.AppId)
}

func checkAppDevelopPermBySourceId(ctx context.Context, operator apisession.UserInfo, sourceId int64) (discoverymd.EtcdNode, error) {
	node, b, err := discoverymd.GetEtcdNodeById(ctx, sourceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return discoverymd.EtcdNode{}, util.InternalError(err)
	}
	if !b {
		return discoverymd.EtcdNode{}, util.InvalidArgsError()
	}
	return node, checkAppDevelopPermByAppId(ctx, operator, node.AppId)
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
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	if p.PermDetail.TeamPerm.CanManageDiscoverySource {
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
