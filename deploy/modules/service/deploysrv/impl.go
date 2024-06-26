package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"gopkg.in/yaml.v3"
)

type outerImpl struct{}

func NewOuterService() OuterService {
	initRunner()
	return new(outerImpl)
}

// CreatePlan 创建发布计划
func (*outerImpl) CreatePlan(ctx context.Context, reqDTO CreatePlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	service, err := checkDeployPlanPermByServiceId(ctx, reqDTO.Operator, reqDTO.ServiceId)
	if err != nil {
		return err
	}
	// 检查是否有其他发布计划在
	b, err := deploymd.ExistPendingOrRunningPlanByServiceId(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	_, err = deploymd.InsertPlan(ctx, deploymd.InsertPlanReqDTO{
		Name:           reqDTO.Name,
		PlanStatus:     deploymd.PendingPlanStatus,
		AppId:          service.AppId,
		ServiceId:      reqDTO.ServiceId,
		ProductVersion: reqDTO.ProductVersion,
		Creator:        reqDTO.Operator.Account,
		Env:            service.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// StartPlan 开始发布计划
func (*outerImpl) StartPlan(ctx context.Context, reqDTO StartPlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, b, err := deploymd.GetPlanById(ctx, reqDTO.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || plan.PlanStatus != deploymd.PendingPlanStatus {
		return util.InvalidArgsError()
	}
	service, err := checkAppDevelopPermByServiceId(ctx, reqDTO.Operator, plan.ServiceId)
	if err != nil {
		return err
	}
	var serviceConfig deploy.Service
	err = yaml.Unmarshal([]byte(service.Config), &serviceConfig)
	if err != nil || !serviceConfig.IsValid() {
		return util.ThereHasBugErr()
	}
	insertDeployList := serviceConfig2DeployServiceReq(reqDTO.PlanId, service.Id, &serviceConfig)
	insertStageList := serviceConfig2DeployStageReq(reqDTO.PlanId, &serviceConfig)
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b2, err2 := deploymd.StartPlan(ctx, reqDTO.PlanId, serviceConfig)
		if err2 != nil {
			return err2
		}
		if b2 {
			// 删除之前的部署列表
			err2 = deploymd.DeleteDeployServiceByServiceId(ctx, service.Id)
			if err2 != nil {
				return err2
			}
			// 新增部署列表
			err2 = deploymd.BatchInsertDeployService(ctx, insertDeployList...)
			if err2 != nil {
				return err2
			}
			// 新增部署步骤
			return deploymd.BatchInsertDeployStage(ctx, insertStageList...)
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	executeDeployOnStartPlanService(reqDTO.PlanId, service.AppId, serviceConfig, map[string]string{
		deploy.CurrentProductVersionKey: plan.ProductVersion,
		deploy.OperatorAccountKey:       reqDTO.Operator.Account,
		deploy.AppKey:                   plan.AppId,
	})
	return nil
}

func serviceConfig2DeployStageReq(planId int64, s *deploy.Service) []deploymd.InsertDeployStageReqDTO {
	ret := make([]deploymd.InsertDeployStageReqDTO, 0)
	for i, stage := range s.Deploy {
		if len(stage.Agents) > 0 {
			for _, agent := range stage.Agents {
				ret = append(ret, deploymd.InsertDeployStageReqDTO{
					PlanId:      planId,
					StageIndex:  i,
					Agent:       agent,
					StageStatus: deploymd.PendingStageStatus,
				})
			}
		} else {
			for agent := range s.Agents {
				ret = append(ret, deploymd.InsertDeployStageReqDTO{
					PlanId:      planId,
					StageIndex:  i,
					Agent:       agent,
					StageStatus: deploymd.PendingStageStatus,
				})
			}
		}
	}
	return ret
}

func serviceConfig2DeployServiceReq(planId, serviceId int64, s *deploy.Service) []deploymd.InsertDeployServiceReqDTO {
	ret := make([]deploymd.InsertDeployServiceReqDTO, 0)
	for _, process := range s.Process {
		c := &deploymd.ProcessConfig{
			Host:       process.Name,
			AgentHost:  s.Agents[process.Agent].Host,
			AgentToken: s.Agents[process.Agent].Token,
		}
		ret = append(ret, deploymd.InsertDeployServiceReqDTO{
			PlanId:    planId,
			ServiceId: serviceId,
			Config: deploymd.DeployServiceConfig{
				Type:    s.Type,
				Process: c,
			},
			Probed: 0,
		})
	}
	if s.K8s != nil {
		ret = append(ret, deploymd.InsertDeployServiceReqDTO{
			PlanId:    planId,
			ServiceId: serviceId,
			Config: deploymd.DeployServiceConfig{
				Type: s.Type,
				K8s: &deploymd.K8sConfig{
					AgentHost:       s.Agents[s.K8s.Agent].Host,
					AgentToken:      s.Agents[s.K8s.Agent].Token,
					GetStatusScript: s.K8s.GetStatusScript,
				},
			},
			Probed: 0,
		})
	}
	return ret
}

// ClosePlan 关闭发布计划
func (*outerImpl) ClosePlan(ctx context.Context, reqDTO ClosePlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, b, err := deploymd.GetPlanById(ctx, reqDTO.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || plan.PlanStatus == deploymd.ClosedPlanStatus {
		return util.InvalidArgsError()
	}
	if _, err = checkAppDevelopPermByServiceId(ctx, reqDTO.Operator, plan.ServiceId); err != nil {
		return err
	}
	b, err = deploymd.ExistUnsuccessfulDeployStage(ctx, reqDTO.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		_, err = deploymd.ClosePlan(ctx, reqDTO.PlanId, plan.PlanStatus)
	} else {
		err = xormstore.WithTx(ctx, func(ctx context.Context) error {
			_, err2 := deploymd.ClosePlan(ctx, reqDTO.PlanId, plan.PlanStatus)
			if err2 != nil {
				return err2
			}
			return deploymd.UpdateDeployServiceIsPlanDoneTrue(ctx, reqDTO.PlanId)
		})
	}
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListPlan 发布计划列表
func (*outerImpl) ListPlan(ctx context.Context, reqDTO ListPlanReqDTO) ([]PlanDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return nil, 0, err
	}
	const pageSize = 10
	plans, total, err := deploymd.ListPlan(ctx, deploymd.ListPlanReqDTO{
		AppId:    reqDTO.AppId,
		PageNum:  reqDTO.PageNum,
		PageSize: pageSize,
		Env:      reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	serviceMap := make(map[int64]string, len(plans))
	serviceIdList, _ := listutil.Map(plans, func(t deploymd.Plan) (int64, error) {
		return t.ServiceId, nil
	})
	serviceIdList = listutil.Distinct(serviceIdList...)
	if len(serviceIdList) > 0 {
		services, err := deploymd.BatchGetServiceById(ctx, serviceIdList, []string{"id", "name"})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, 0, util.InternalError(err)
		}
		for _, srv := range services {
			serviceMap[srv.Id] = srv.Name
		}
	}
	data, _ := listutil.Map(plans, func(t deploymd.Plan) (PlanDTO, error) {
		return PlanDTO{
			Id:             t.Id,
			ServiceId:      t.ServiceId,
			ServiceName:    serviceMap[t.ServiceId],
			Name:           t.Name,
			ProductVersion: t.ProductVersion,
			PlanStatus:     t.PlanStatus,
			Env:            t.Env,
			Creator:        t.Creator,
			Created:        t.Created,
		}, nil
	})
	return data, total, nil
}

// ConfirmPlanServiceStep 执行流水线其中一个环节
func (*outerImpl) ConfirmPlanServiceStep(ctx context.Context, reqDTO ConfirmPlanServiceStepReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	return nil
}

// RollbackPlanServiceStep 回滚执行流水线其中一个环节
func (*outerImpl) RollbackPlanServiceStep(ctx context.Context, reqDTO RollbackPlanServiceStepReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	return nil
}

// CreateService 创建服务
func (*outerImpl) CreateService(ctx context.Context, reqDTO CreateServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePerm(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return err
	}
	err := deploymd.InsertService(ctx, deploymd.InsertServiceReqDTO{
		AppId:       reqDTO.AppId,
		Config:      reqDTO.Config,
		Env:         reqDTO.Env,
		Name:        reqDTO.Name,
		ServiceType: reqDTO.service.Type,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateService 编辑服务
func (*outerImpl) UpdateService(ctx context.Context, reqDTO UpdateServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePermByServiceId(ctx, reqDTO.ServiceId, reqDTO.Operator); err != nil {
		return err
	}
	_, err := deploymd.UpdateService(ctx, deploymd.UpdateServiceReqDTO{
		ServiceId:   reqDTO.ServiceId,
		Name:        reqDTO.Name,
		Config:      reqDTO.Config,
		ServiceType: reqDTO.service.Type,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteService 删除服务
func (*outerImpl) DeleteService(ctx context.Context, reqDTO DeleteServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePermByServiceId(ctx, reqDTO.ServiceId, reqDTO.Operator); err != nil {
		return err
	}
	// 存在正在执行的发布计划
	b, err := deploymd.ExistPendingOrRunningPlanByServiceId(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.AlreadyExistsError()
	}
	_, err = deploymd.DeleteService(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListService 服务列表
func (*outerImpl) ListService(ctx context.Context, reqDTO ListServiceReqDTO) ([]ServiceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePerm(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return nil, err
	}
	services, err := deploymd.ListService(ctx, deploymd.ListServiceReqDTO{
		AppId: reqDTO.AppId,
		Env:   reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(services, func(t deploymd.Service) (ServiceDTO, error) {
		return ServiceDTO{
			Id:          t.Id,
			AppId:       t.AppId,
			Config:      t.Config,
			Env:         t.Env,
			Name:        t.Name,
			ServiceType: t.ServiceType,
		}, nil
	})
	return data, nil
}

// ListServiceWhenCreatePlan 创建发布计划时展示的服务列表
func (*outerImpl) ListServiceWhenCreatePlan(ctx context.Context, reqDTO ListServiceWhenCreatePlanReqDTO) ([]SimpleServiceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkDeployPlanPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return nil, err
	}
	services, err := deploymd.ListService(ctx, deploymd.ListServiceReqDTO{
		AppId: reqDTO.AppId,
		Env:   reqDTO.Env,
		Cols:  []string{"id", "env", "name", "service_type"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(services, func(t deploymd.Service) (SimpleServiceDTO, error) {
		return SimpleServiceDTO{
			Id:          t.Id,
			Env:         t.Env,
			Name:        t.Name,
			ServiceType: t.ServiceType,
		}, nil
	})
	return data, nil
}

func checkServicePerm(ctx context.Context, appId string, operator apisession.UserInfo) error {
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
	if !b || (!p.IsAdmin && !p.PermDetail.TeamPerm.CanManageService) {
		return util.UnauthorizedError()
	}
	return nil
}

func checkServicePermByServiceId(ctx context.Context, probeId int64, operator apisession.UserInfo) error {
	probe, b, err := deploymd.GetServiceById(ctx, probeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkServicePerm(ctx, probe.AppId, operator)
}
