package discoverysrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/discovery/modules/model/discoverymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/services/lb"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
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
						case event.AppManageDiscoverySourceAction:
							return cfg.AppSource.ManageDiscoverySource
						default:
							return false
						}
					}
					return false
				})
			}
		})
		psub.Subscribe(event.AppDiscoveryTopic, func(data any) {
			req, ok := data.(event.AppDiscoveryEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Source.Env]
					if ok {
						switch req.Action {
						case event.AppDiscoveryMarkAsDownAction:
							return cfg.AppDiscovery.MarkAsDown
						case event.AppDiscoveryMarkAsUpAction:
							return cfg.AppDiscovery.MarkAsUp
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
	app, team, err := checkManageDiscoverySourcePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
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
	nodes, err := discoverymd.BatchGetEtcdNodeByIdList(ctx, reqDTO.SourceIdList, []string{"id", "name", "env"})
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
	insertReqs := listutil.MapNe(nodes, func(t discoverymd.EtcdNode) discoverymd.InsertAppEtcdNodeBindReqDTO {
		return discoverymd.InsertAppEtcdNodeBindReqDTO{
			NodeId: t.Id,
			AppId:  reqDTO.AppId,
			Env:    reqDTO.Env,
		}
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
	notifyDiscoverySourceEvent(
		reqDTO.Operator,
		team,
		app,
		nodes,
		reqDTO.Env,
	)
	return nil
}

// ListBindDiscoverySource 获取绑定注册中心来源
func ListBindDiscoverySource(ctx context.Context, reqDTO ListBindDiscoverySourceReqDTO) ([]SimpleBindDiscoverySourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
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
	nodeIdList := listutil.MapNe(binds, func(t discoverymd.AppEtcdNodeBind) int64 {
		bindMap[t.NodeId] = t
		return t.NodeId
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
	app, _, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
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
	prefix := common.ServicePrefix + app.AppId + "-http/"
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
				InstanceId: strings.TrimPrefix(string(kv.Key), prefix),
			})
		}
	}
	return ret, nil
}

// MarkAsDownService 下线服务
func MarkAsDownService(ctx context.Context, reqDTO MarkAsDownServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	app, team, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return err
	}
	client, err := newEtcdClient(source)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	defer client.Close()
	kv := clientv3.NewKV(client)
	response, err := kv.Get(ctx, common.ServicePrefix+app.AppId+"-http/"+reqDTO.InstanceId)
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
	// 已经下线
	if s.IsDown {
		return nil
	}
	resp, err := http.Get(fmt.Sprintf("%s://%s:%d/actuator/v1/markAsDownServer", s.Protocol, s.Host, s.Port))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.OperationFailedError()
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return util.OperationFailedError()
	}
	notifyDiscoveryEvent(
		reqDTO.Operator,
		team,
		app,
		event.AppDiscoveryMarkAsDownAction,
		source,
	)
	time.Sleep(time.Second)
	return nil
}

// MarkAsUpService 上线服务
func MarkAsUpService(ctx context.Context, reqDTO MarkAsUpServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	app, team, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return err
	}
	client, err := newEtcdClient(source)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	defer client.Close()
	kv := clientv3.NewKV(client)
	response, err := kv.Get(ctx, common.ServicePrefix+app.AppId+"-http/"+reqDTO.InstanceId)
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
	// 已经上线
	if !s.IsDown {
		return nil
	}
	resp, err := http.Get(fmt.Sprintf("%s://%s:%d/actuator/v1/markAsUpServer", s.Protocol, s.Host, s.Port))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.OperationFailedError()
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return util.OperationFailedError()
	}
	notifyDiscoveryEvent(
		reqDTO.Operator,
		team,
		app,
		event.AppDiscoveryMarkAsUpAction,
		source,
	)
	time.Sleep(time.Second)
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

func checkAppDevelopPermByBindId(ctx context.Context, operator apisession.UserInfo, bindId int64) (appmd.App, teammd.Team, discoverymd.EtcdNode, error) {
	bind, b, err := discoverymd.GetAppEtcdNodeBindById(ctx, bindId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, discoverymd.EtcdNode{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, discoverymd.EtcdNode{}, util.InvalidArgsError()
	}
	node, b, err := discoverymd.GetEtcdNodeById(ctx, bind.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, discoverymd.EtcdNode{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, discoverymd.EtcdNode{}, util.InvalidArgsError()
	}
	app, team, err := checkAppDevelopPermByAppId(ctx, operator, bind.AppId)
	return app, team, node, err
}

func checkManageDiscoverySourcePermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (appmd.App, teammd.Team, error) {
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
	if p.IsAdmin {
		return app, team, nil
	}
	if p.PermDetail.GetAppPerm(appId).CanManageDiscoverySource {
		return app, team, nil
	}
	return app, team, util.UnauthorizedError()
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

func notifyDiscoveryEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, action event.AppDiscoveryEventAction, source discoverymd.EtcdNode) {
	initPsub()
	psub.Publish(event.AppDiscoveryTopic, event.AppDiscoveryEvent{
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
		Action: action,
		Source: event.AppSource{
			Id:   source.Id,
			Name: source.Name,
			Env:  source.Env,
		},
	})
}

func notifyDiscoverySourceEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, sources []discoverymd.EtcdNode, env string) {
	initPsub()
	srcs := listutil.MapNe(sources, func(t discoverymd.EtcdNode) event.AppSource {
		return event.AppSource{
			Id:   t.Id,
			Name: t.Name,
			Env:  t.Env,
		}
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
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, event.AppManageDiscoverySourceAction.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, event.AppManageDiscoverySourceAction.GetI18nValue()),
		},
		Env:     env,
		Action:  event.AppManageDiscoverySourceAction,
		Sources: srcs,
	})
}
