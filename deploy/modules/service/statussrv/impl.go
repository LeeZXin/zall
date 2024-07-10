package statussrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/deploy/modules/model/statusmd"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zall/pkg/status"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

var (
	zalletCmd string
	sockFile  string
)

func NewOuterService() OuterService {
	zalletCmd = static.GetString("zallet.cmd")
	if zalletCmd == "" {
		zalletCmd = "zallet"
	}
	sockFile = static.GetString("zallet.sock")
	return new(impl)
}

type impl struct{}

// GetActions 获取基本信息
func (*impl) GetActions(context.Context) []status.Action {
	return []status.Action{
		{
			Label: "重启",
			Api: status.Api{
				Url:    "/api/service/v1/status/restart",
				Method: http.MethodPut,
			},
		},
		{
			Label: "关闭",
			Api: status.Api{
				Url:    "/api/service/v1/status/kill",
				Method: http.MethodPut,
			},
		},
	}
}

// ListService 服务列表
func (*impl) ListService(ctx context.Context, reqDTO status.ListServiceReq) ([]status.Service, error) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	services, err := statusmd.ListService(ctx, statusmd.ListServiceReqDTO{
		AppId: reqDTO.App,
		Env:   reqDTO.Env,
	})
	if err != nil {
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(services, func(t statusmd.Service) (status.Service, error) {
		addr, _ := netip.ParseAddrPort(t.AgentHost)
		return status.Service{
			Id:         t.ServiceId,
			App:        t.App,
			Status:     t.ServiceStatus,
			Host:       addr.Addr().String(),
			Env:        t.Env,
			CpuPercent: t.CpuPercent,
			MemPercent: t.MemPercent,
			Created:    t.Created.Format(time.DateTime),
		}, nil
	})
	return data, nil
}

// KillService 关闭服务
func (*impl) KillService(ctx context.Context, serviceId string) error {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	service, b, err := statusmd.GetServiceByServiceId(ctx, serviceId)
	if err != nil {
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	var cmd string
	if sockFile == "" {
		cmd = fmt.Sprintf("%s kill -service %s", zalletCmd, serviceId)
	} else {
		cmd = fmt.Sprintf("%s kill -service %s -sock %s", zalletCmd, serviceId, sockFile)
	}
	logger.Logger.WithContext(ctx).Infof("serviceId: %s cmd: %s", serviceId, cmd)
	_, err = sshagent.NewServiceCommand(service.AgentHost, service.AgentToken, service.App).
		Execute(strings.NewReader(cmd), nil, idutil.RandomUuid())
	if err != nil {
		return bizerr.NewBizErr(apicode.OperationFailedErrCode.Int(), err.Error())
	}
	return nil
}

// RestartService 重启服务
func (*impl) RestartService(ctx context.Context, serviceId string) error {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	service, b, err := statusmd.GetServiceByServiceId(ctx, serviceId)
	if err != nil {
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	var cmd string
	if sockFile == "" {
		cmd = fmt.Sprintf("%s restart -service %s", zalletCmd, serviceId)
	} else {
		cmd = fmt.Sprintf("%s restart -service %s -sock %s", zalletCmd, serviceId, sockFile)
	}
	logger.Logger.WithContext(ctx).Infof("serviceId: %s cmd: %s", serviceId, cmd)
	_, err = sshagent.NewServiceCommand(service.AgentHost, service.AgentToken, service.App).
		Execute(strings.NewReader(cmd), nil, idutil.RandomUuid())
	if err != nil {
		return bizerr.NewBizErr(apicode.OperationFailedErrCode.Int(), err.Error())
	}
	return nil
}
