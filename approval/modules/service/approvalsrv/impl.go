package approvalsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/approval/modules/model/approvalmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/approval"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/hashicorp/go-bexpr"
	"sort"
	"time"
)

type innerImpl struct{}

func (*innerImpl) InsertFlow(ctx context.Context, reqDTO InsertFlowReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	process, b, err := approvalmd.GetProcessByPid(ctx, reqDTO.Pid)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("process not found: %s", reqDTO.Pid)
	}
	flow, err := approvalmd.InsertFlow(ctx, approvalmd.InsertFlowReqDTO{
		ProcessId:   process.Id,
		ProcessName: process.Name,
		Process:     *process.Content,
		CurrIndex:   1,
		FlowStatus:  approvalmd.PendingFlowStatus,
		Creator:     reqDTO.Account,
		BizId:       reqDTO.BizId,
	})
	if err != nil {
		return err
	}
	executeFlow(flow)
	return nil
}

func (*innerImpl) InsertAttachedProcess(ctx context.Context, reqDTO InsertAttachedProcessReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := approvalmd.GetProcessByPid(ctx, reqDTO.Pid)
	if err != nil {
		return err
	}
	if b {
		return fmt.Errorf("pid: %s already exists", reqDTO.Pid)
	}
	return approvalmd.InsertProcess(ctx, approvalmd.InsertProcessReqDTO{
		Pid:        reqDTO.Pid,
		Name:       reqDTO.Name,
		Process:    reqDTO.Process.Convert(),
		SourceType: approvalmd.SystemSourceType,
	})
}

func (*innerImpl) UpdateAttachedProcess(ctx context.Context, reqDTO UpdateAttachedProcessReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := approvalmd.UpdateProcessByPid(ctx, approvalmd.UpdateProcessByPidReqDTO{
		Pid:     reqDTO.Pid,
		Name:    reqDTO.Name,
		Process: reqDTO.Process.Convert(),
	})
	return err
}

// DeleteAttachedProcess 删除系统审批流
func (*innerImpl) DeleteAttachedProcess(ctx context.Context, pid string) error {
	if pid == "" {
		return util.InvalidArgsError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := approvalmd.DeleteProcessByPid(ctx, pid)
	return err
}

func executeFlow(flow approvalmd.Flow) {
	go func() {
		flowCtx := &approval.FlowContext{
			Context: context.Background(),
			FlowId:  flow.Id,
			BizId:   flow.BizId,
			Kvs:     flow.Kvs,
			Process: flow.ProcessContent,
		}
		runFlow(flowCtx, &flow, flow.ProcessContent.Node)
	}()
}

func runFlow(flowCtx *approval.FlowContext, flow *approvalmd.Flow, node *approval.Node) {
	if node == nil {
		return
	}
	node.FindAndDo(flow.CurrIndex, map[approval.NodeType]func(*approval.Node){
		approval.PeopleNode: func(node *approval.Node) {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			err := xormstore.WithTx(ctx, func(ctx context.Context) error {
				return approvalmd.InsertNotify(ctx, approvalmd.InsertNotifyReqDTO{
					FlowId:    flow.Id,
					Accounts:  node.Accounts,
					Done:      false,
					FlowIndex: flow.CurrIndex,
				})
			})
			if err != nil {
				logger.Logger.Error(err)
			}
		},
		approval.ApiNode: func(node *approval.Node) {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			if node.Api == nil {
				logger.Logger.Errorf("flowId: %d currIndex: %d its api node has no config", flow.Id, flow.CurrIndex)
				_, err := approvalmd.UpdateFlowStatusAndErrMsgWithOldStatus(
					ctx,
					flow.Id,
					i18n.GetByKey(i18n.SystemInternalError),
					approvalmd.ErrFlowStatus,
					approvalmd.PendingFlowStatus,
				)
				if err != nil {
					logger.Logger.Error(err)
				}
				return
			}
			response, err := node.Api.DoRequest(flowCtx)
			if err != nil {
				logger.Logger.Errorf("flowId: %d currIndex: %d its api request err: %v", flow.Id, flow.CurrIndex, err)
				_, err = approvalmd.UpdateFlowStatusAndErrMsgWithOldStatus(
					ctx,
					flow.Id,
					i18n.GetByKey(i18n.SystemInternalError),
					approvalmd.ErrFlowStatus,
					approvalmd.PendingFlowStatus,
				)
				if err != nil {
					logger.Logger.Error(err)
				}
				return
			}
			runNext(flowCtx, flow, node, response)
		},
		approval.MethodNode: func(node *approval.Node) {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			if node.Method == nil {
				logger.Logger.Errorf("flowId: %d currIndex: %d its method node has no config", flow.Id, flow.CurrIndex)
				_, err := approvalmd.UpdateFlowStatusAndErrMsgWithOldStatus(
					ctx,
					flow.Id,
					i18n.GetByKey(i18n.SystemInternalError),
					approvalmd.ErrFlowStatus,
					approvalmd.PendingFlowStatus,
				)
				if err != nil {
					logger.Logger.Error(err)
				}
				return
			}
			response, err := node.Method.DoMethod(flowCtx)
			if err != nil {
				logger.Logger.Errorf("flowId: %d currIndex: %d its method: %s err: %v", flow.Id, flow.CurrIndex, node.Method.Name, err)
				_, err = approvalmd.UpdateFlowStatusAndErrMsgWithOldStatus(
					ctx,
					flow.Id,
					i18n.GetByKey(i18n.SystemInternalError),
					approvalmd.ErrFlowStatus,
					approvalmd.PendingFlowStatus,
				)
				if err != nil {
					logger.Logger.Error(err)
				}
				return
			}
			runNext(flowCtx, flow, node, response)
		},
		approval.DisagreeNode: func(node *approval.Node) {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			_, err := approvalmd.UpdateFlowStatusWithOldStatus(
				ctx,
				flow.Id,
				approvalmd.DisagreeFlowStatus,
				approvalmd.PendingFlowStatus,
			)
			if err != nil {
				logger.Logger.Error(err)
			}
		},
		approval.AgreeNode: func(node *approval.Node) {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			_, err := approvalmd.UpdateFlowStatusWithOldStatus(
				ctx,
				flow.Id,
				approvalmd.AgreeFlowStatus,
				approvalmd.PendingFlowStatus,
			)
			if err != nil {
				logger.Logger.Error(err)
			}
		},
	})
}

func runNext(flowCtx *approval.FlowContext, flow *approvalmd.Flow, node *approval.Node, response map[string]any) {
	ctx, closer := xormstore.Context(flowCtx)
	defer closer.Close()
	if len(node.Next) == 0 {
		logger.Logger.Errorf("flowId: %d currIndex: %d has no next", flow.Id, flow.CurrIndex)
		_, err := approvalmd.UpdateFlowStatusAndErrMsgWithOldStatus(
			ctx,
			flow.Id,
			i18n.GetByKey(i18n.SystemInternalError),
			approvalmd.ErrFlowStatus,
			approvalmd.PendingFlowStatus,
		)
		if err != nil {
			logger.Logger.Error(err)
		}
		return
	}
	for _, next := range node.Next {
		if doCondition(next.Condition, response) {
			// 更新currIndex
			b, err := approvalmd.UpdateFlowCurrIndexWithOldCurrIndex(ctx, flow.Id, next.Node.NodeId, flow.CurrIndex)
			if err != nil {
				logger.Logger.Error(err)
				return
			}
			if b {
				flow.CurrIndex = next.Node.NodeId
				runFlow(flowCtx, flow, next.Node)
			}
			return
		}
	}
	logger.Logger.Errorf("flowId: %d currIndex: %d no next condition matches", flow.Id, flow.CurrIndex)
	_, err := approvalmd.UpdateFlowStatusAndErrMsgWithOldStatus(
		ctx,
		flow.Id,
		i18n.GetByKey(i18n.SystemInternalError),
		approvalmd.ErrFlowStatus,
		approvalmd.PendingFlowStatus,
	)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
}

func doCondition(condition string, vars map[string]any) bool {
	if condition == "" {
		return true
	}
	eval, err := bexpr.CreateEvaluator(condition)
	if err != nil {
		logger.Logger.Errorf("condition: %s copmile with err: %v", condition, err)
		return false
	}
	if vars == nil {
		vars = make(map[string]any)
	}
	ret, err := eval.Evaluate(vars)
	if err != nil {
		logger.Logger.Errorf("condition: %s with vars: %v evaluate with err: %v", condition, vars, err)
		return false
	}
	return ret
}

type outerImpl struct{}

func (*outerImpl) Agree(ctx context.Context, reqDTO AgreeFlowReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	notify, b, err := approvalmd.GetNotifyById(ctx, reqDTO.NotifyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || notify.Account != reqDTO.Operator.Account {
		return util.InvalidArgsError()
	}
	if notify.Done {
		return util.AlreadyExistsError()
	}
	flow, b, err := approvalmd.GetFlowById(ctx, notify.FlowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || flow.FlowStatus != approvalmd.PendingFlowStatus {
		return util.InvalidArgsError()
	}
	if flow.CurrIndex != notify.FlowIndex {
		return util.InvalidArgsError()
	}
	process := flow.ProcessContent
	node := process.Find(notify.FlowIndex)
	if node == nil {
		return util.InvalidArgsError()
	}
	if node.NodeType != approval.PeopleNode {
		return util.InvalidArgsError()
	}
	findAccount := false
	for _, account := range node.Accounts {
		if account == reqDTO.Operator.Account {
			findAccount = true
			break
		}
	}
	if !findAccount {
		return util.UnauthorizedError()
	}
	b, err = approvalmd.UpdateNotifyDoneWithOldDone(ctx, reqDTO.NotifyId, true, false, approvalmd.AgreeOp)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		count, err := approvalmd.CountOperate(ctx, flow.Id, flow.CurrIndex, approvalmd.AgreeOp)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if int(count) >= node.AtLeastNum {
			// 其他人被迫同意
			_ = approvalmd.UpdateNotifyDoneWithOldDoneByFlowIdAndIndex(ctx, true, false, approvalmd.AutoAgreeOp, flow.Id, flow.CurrIndex)
			kvs := flow.Kvs
			flowCtx := &approval.FlowContext{
				Context: context.Background(),
				FlowId:  flow.Id,
				BizId:   flow.BizId,
				Kvs:     kvs,
				Process: process,
			}
			go runNext(flowCtx, &flow, node, map[string]any{
				"agree": "y",
			})
		}
	}
	return nil
}

func (*outerImpl) Disagree(ctx context.Context, reqDTO DisagreeFlowReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	notify, b, err := approvalmd.GetNotifyById(ctx, reqDTO.NotifyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || notify.Account != reqDTO.Operator.Account {
		return util.InvalidArgsError()
	}
	if notify.Done {
		return util.AlreadyExistsError()
	}
	flow, b, err := approvalmd.GetFlowById(ctx, notify.FlowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || flow.FlowStatus != approvalmd.PendingFlowStatus {
		return util.InvalidArgsError()
	}
	if flow.CurrIndex != notify.FlowIndex {
		return util.InvalidArgsError()
	}
	process := flow.ProcessContent
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	node := process.Find(flow.CurrIndex)
	if node == nil {
		return util.InvalidArgsError()
	}
	if node.NodeType != approval.PeopleNode {
		return util.InvalidArgsError()
	}
	findAccount := false
	for _, account := range node.Accounts {
		if account == reqDTO.Operator.Account {
			findAccount = true
			break
		}
	}
	if !findAccount {
		return util.UnauthorizedError()
	}
	b, err = approvalmd.UpdateNotifyDoneWithOldDone(ctx, reqDTO.NotifyId, true, false, approvalmd.DisagreeOp)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		count, err := approvalmd.CountOperate(ctx, flow.Id, flow.CurrIndex, approvalmd.DisagreeOp)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if int(count)+node.AtLeastNum > len(node.Accounts) {
			// 其他人被迫不同意
			_ = approvalmd.UpdateNotifyDoneWithOldDoneByFlowIdAndIndex(ctx, true, false, approvalmd.AutoDisagreeOp, flow.Id, flow.CurrIndex)
			flowCtx := &approval.FlowContext{
				Context: context.Background(),
				FlowId:  flow.Id,
				BizId:   flow.BizId,
				Kvs:     flow.Kvs,
				Process: process,
			}
			go runNext(flowCtx, &flow, node, map[string]any{
				"agree": "n",
			})
		}
	}
	return nil
}

// InsertCustomProcess 创建自定义审批流
func (*outerImpl) InsertCustomProcess(ctx context.Context, reqDTO InsertCustomProcessReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.ApprovalSrvKeysVO.InsertCustomProcess),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 检查权限
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := approvalmd.GetGroupById(ctx, reqDTO.GroupId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	_, b, err = approvalmd.GetProcessByPid(ctx, reqDTO.Pid)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.AlreadyExistsError()
		return
	}
	err = approvalmd.InsertProcess(ctx, approvalmd.InsertProcessReqDTO{
		Pid:        reqDTO.Pid,
		Name:       reqDTO.Name,
		GroupId:    reqDTO.GroupId,
		IconUrl:    reqDTO.IconUrl,
		Process:    reqDTO.Process.Convert(),
		SourceType: approvalmd.CustomSourceType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// UpdateCustomProcess 编辑自定义审批流
func (*outerImpl) UpdateCustomProcess(ctx context.Context, reqDTO UpdateCustomProcessReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.ApprovalSrvKeysVO.UpdateCustomProcess),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 检查权限
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := approvalmd.GetGroupById(ctx, reqDTO.GroupId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	_, err = approvalmd.UpdateProcessById(ctx, approvalmd.UpdateProcessByIdReqDTO{
		Id:      reqDTO.Id,
		Name:    reqDTO.Name,
		GroupId: reqDTO.GroupId,
		IconUrl: reqDTO.IconUrl,
		Process: reqDTO.Process.Convert(),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// InsertCustomFlow 发起自定义审批流
func (*outerImpl) InsertCustomFlow(ctx context.Context, reqDTO InsertCustomFlowReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	process, b, err := approvalmd.GetProcessByPid(ctx, reqDTO.Pid)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	errKeys := process.Content.CheckKvCfgs(reqDTO.Kvs)
	if len(errKeys) > 0 {
		return errKeys, nil
	}
	flow, err := approvalmd.InsertFlow(ctx, approvalmd.InsertFlowReqDTO{
		ProcessId:   process.Id,
		ProcessName: process.Name,
		Process:     *process.Content,
		CurrIndex:   1,
		FlowStatus:  approvalmd.PendingFlowStatus,
		Kvs:         reqDTO.Kvs,
		Creator:     reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	executeFlow(flow)
	return nil, nil
}

// CancelCustomFlow 取消自定义审批流
func (*outerImpl) CancelCustomFlow(ctx context.Context, reqDTO CancelCustomFlowReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	flow, b, err := approvalmd.GetFlowById(ctx, reqDTO.FlowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || flow.FlowStatus != approvalmd.PendingFlowStatus {
		return util.InvalidArgsError()
	}
	if flow.Creator != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err := approvalmd.UpdateFlowStatusWithOldStatus(ctx, reqDTO.FlowId, approvalmd.CanceledFlowStatus, approvalmd.PendingFlowStatus)
		if err != nil {
			return err
		}
		return approvalmd.UpdateNotifyDoneWithOldDoneByFlowId(ctx, true, false, approvalmd.CancelOp, reqDTO.FlowId)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListAllGroupProcess 自定义审批流列表
func (*outerImpl) ListAllGroupProcess(ctx context.Context, reqDTO ListAllGroupProcessReqDTO) ([]GroupProcessDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	groups, err := approvalmd.GetAllGroups(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	processes, err := approvalmd.GetAllCustomProcesses(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	processMap := make(map[int64][]approvalmd.SimpleProcess, 8)
	for _, process := range processes {
		pl, b := processMap[process.GroupId]
		if !b {
			pl = make([]approvalmd.SimpleProcess, 0)
		}
		pl = append(pl, process)
		processMap[process.GroupId] = pl
	}
	groupsRet, _ := listutil.Map(groups, func(t approvalmd.Group) (GroupProcessDTO, error) {
		pl := processMap[t.Id]
		ps, _ := listutil.Map(pl, func(t approvalmd.SimpleProcess) (SimpleProcessDTO, error) {
			return SimpleProcessDTO{
				Id:      t.Id,
				Name:    t.Name,
				IconUrl: t.IconUrl,
			}, nil
		})
		sort.SliceStable(ps, func(i, j int) bool {
			return ps[i].Id < ps[j].Id
		})
		return GroupProcessDTO{
			Id:        t.Id,
			Name:      t.Name,
			Processes: ps,
		}, nil
	})
	sort.SliceStable(groupsRet, func(i, j int) bool {
		return groupsRet[i].Id < groupsRet[j].Id
	})
	return groupsRet, nil
}

// ListCustomFlow 获取申请的审批流
func (*outerImpl) ListCustomFlow(ctx context.Context, reqDTO ListCustomFlowReqDTO) ([]FlowDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	startTime, _ := time.Parse(time.DateTime, reqDTO.DayTime+" 00:00:00")
	endTime, _ := time.Parse(time.DateTime, reqDTO.DayTime+" 23:59:59")
	flows, err := approvalmd.GetFlowByCreatorAndTime(ctx, reqDTO.Operator.Account, startTime, endTime)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(flows, func(t approvalmd.Flow) (FlowDTO, error) {
		return FlowDTO{
			Id:          t.Id,
			ProcessName: t.ProcessName,
			FlowStatus:  t.FlowStatus,
			Creator:     t.Creator,
			Created:     t.Created,
		}, nil
	})
}

// ListOperateFlow 获取所有未审批提醒
func (*outerImpl) ListOperateFlow(ctx context.Context, reqDTO ListOperateFlowReqDTO) ([]FlowDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	startTime, _ := time.Parse(time.DateTime, reqDTO.DayTime+" 00:00:00")
	endTime, _ := time.Parse(time.DateTime, reqDTO.DayTime+" 23:59:59")
	notifies, err := approvalmd.GetNotifyByAccountAndTime(ctx, reqDTO.Operator.Account, startTime, endTime, reqDTO.Done)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	flowIds, _ := listutil.Map(notifies, func(t approvalmd.Notify) (int64, error) {
		return t.FlowId, nil
	})
	flows, err := approvalmd.BatchGetFlows(ctx, flowIds)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(flows, func(t approvalmd.Flow) (FlowDTO, error) {
		return FlowDTO{
			Id:          t.Id,
			ProcessName: t.ProcessName,
			FlowStatus:  t.FlowStatus,
			Creator:     t.Creator,
			Created:     t.Created,
		}, nil
	})
}

// DeleteCustomProcess 删除自定义审批流
func (*outerImpl) DeleteCustomProcess(ctx context.Context, reqDTO DeleteCustomProcessReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.ApprovalSrvKeysVO.DeleteCustomProcess),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 系统管理员
	if !reqDTO.Operator.IsAdmin {
		return util.InvalidArgsError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = approvalmd.DeleteProcessById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// GetFlowDetail 获取审批流详情
func (*outerImpl) GetFlowDetail(ctx context.Context, reqDTO GetFlowDetailReqDTO) (FlowDetailDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return FlowDetailDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	flow, b, err := approvalmd.GetFlowById(ctx, reqDTO.FlowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return FlowDetailDTO{}, util.InternalError(err)
	}
	if !b {
		return FlowDetailDTO{}, util.InvalidArgsError()
	}
	if !reqDTO.Operator.IsAdmin && flow.Creator != reqDTO.Operator.Account {
		b, err = approvalmd.ExistNotifyByAccountAndFlowId(ctx, reqDTO.Operator.Account, reqDTO.FlowId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return FlowDetailDTO{}, util.InternalError(err)
		}
		if !b {
			return FlowDetailDTO{}, util.UnauthorizedError()
		}
	}
	process := flow.ProcessContent
	kvs := flow.Kvs
	notifies, err := approvalmd.GetNotifyByFlowId(ctx, flow.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return FlowDetailDTO{}, util.InternalError(err)
	}
	ret := FlowDetailDTO{
		Id:          flow.Id,
		ProcessName: flow.ProcessName,
		FlowStatus:  flow.FlowStatus,
		Creator:     flow.Creator,
		Created:     flow.Created,
		Kvs:         kvs,
		Process:     *process,
	}
	ret.NotifyList, _ = listutil.Map(notifies, func(t approvalmd.Notify) (NotifyDTO, error) {
		return NotifyDTO{
			Account:   t.Account,
			FlowIndex: t.FlowIndex,
			Done:      t.Done,
			Op:        t.Op,
			Updated:   t.Updated,
		}, nil
	})
	return ret, nil
}

// ListCustomProcess 展示详细审批流列表
func (*outerImpl) ListCustomProcess(ctx context.Context, reqDTO ListCustomProcessReqDTO) ([]ProcessDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	processes, err := approvalmd.GetProcessByGroupId(ctx, reqDTO.GroupId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(processes, func(t approvalmd.Process) (ProcessDTO, error) {
		return ProcessDTO{
			Id:      t.Id,
			Pid:     t.Pid,
			GroupId: t.GroupId,
			Name:    t.Name,
			Content: *t.Content,
			IconUrl: t.IconUrl,
			Created: t.Created,
		}, nil
	})
}

// InsertGroup 插入审批流分组
func (*outerImpl) InsertGroup(ctx context.Context, reqDTO InsertGroupReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.ApprovalSrvKeysVO.InsertGroup),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err = approvalmd.InsertGroup(ctx, approvalmd.InsertGroupReqDTO{
		Name: reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return nil
}

// DeleteGroup 删除分组
func (*outerImpl) DeleteGroup(ctx context.Context, reqDTO DeleteGroupReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.ApprovalSrvKeysVO.DeleteGroup),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = approvalmd.DeleteGroup(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListGroup 展示分组
func (*outerImpl) ListGroup(ctx context.Context, reqDTO ListGroupReqDTO) ([]GroupDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	groups, err := approvalmd.GetAllGroups(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(groups, func(t approvalmd.Group) (GroupDTO, error) {
		return GroupDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
}

// UpdateGroup 修改分组
func (*outerImpl) UpdateGroup(ctx context.Context, reqDTO UpdateGroupReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.ApprovalSrvKeysVO.UpdateGroup),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = approvalmd.UpdateGroup(ctx, approvalmd.UpdateGroupReqDTO{
		Id:   reqDTO.Id,
		Name: reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}
