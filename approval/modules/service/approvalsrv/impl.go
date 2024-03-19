package approvalsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/approval/modules/model/approvalmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/approval"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/hashicorp/go-bexpr"
)

func InitDefaultGroup() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, b, err := approvalmd.GetGroupById(ctx, 1)
	if err != nil {
		logger.Logger.Fatalf("get default approval group failed")
	}
	if !b {
		err = approvalmd.InsertGroup(ctx, approvalmd.InsertGroupReqDTO{
			Name: "default",
		})
		if err != nil {
			logger.Logger.Fatalf("init default approval group failed")
		}
	}
}

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
		ProcessId:  process.Id,
		Process:    process.Process,
		CurrIndex:  1,
		FlowStatus: approvalmd.PendingFlowStatus,
		Creator:    reqDTO.Account,
		BizId:      reqDTO.BizId,
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
		Pid:     reqDTO.Pid,
		Name:    reqDTO.Name,
		GroupId: 1,
		Process: reqDTO.Process.Convert(),
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
		GroupId: 1,
		Process: reqDTO.Process.Convert(),
	})
	return err
}

func executeFlow(flow approvalmd.Flow) {
	go func() {
		process, err := flow.GetProcess()
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		flowCtx := &approval.FlowContext{
			Context: context.Background(),
			FlowId:  flow.Id,
			Process: &process,
			BizId:   flow.BizId,
		}
		runFlow(flowCtx, &flow, process.Node)
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
		logger.Logger.Errorf("condition: %s create with err: %v", condition, err)
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
	if !b {
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
	if !b {
		return util.InvalidArgsError()
	}
	if flow.CurrIndex != notify.FlowIndex {
		return util.InvalidArgsError()
	}
	process, err := flow.GetProcess()
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
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
	ctx, committer, err := xormstore.TxContext(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 事务
	{
		err = approvalmd.InsertDetail(ctx, approvalmd.InsertDetailReqDTO{
			FlowId:    flow.Id,
			FlowIndex: flow.CurrIndex,
			FlowOp:    approvalmd.AgreeOp,
			Account:   reqDTO.Operator.Account,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			committer.Rollback()
			return util.InternalError(err)
		}
		_, err = approvalmd.UpdateNotifyDone(ctx, reqDTO.NotifyId, true)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			committer.Rollback()
			return util.InternalError(err)
		}
		committer.Commit()
	}
	count, err := approvalmd.CountDetail(ctx, flow.Id, flow.CurrIndex, approvalmd.AgreeOp)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if int(count) >= node.AtLeastNum {
		flowCtx := &approval.FlowContext{
			Context: context.Background(),
			FlowId:  flow.Id,
			Process: &process,
			BizId:   flow.BizId,
		}
		go runNext(flowCtx, &flow, node, map[string]any{
			"agree": "y",
		})
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
	if !b {
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
	if !b {
		return util.InvalidArgsError()
	}
	if flow.CurrIndex != notify.FlowIndex {
		return util.InvalidArgsError()
	}
	process, err := flow.GetProcess()
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
	ctx, committer, err := xormstore.TxContext(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 事务
	{
		err = approvalmd.InsertDetail(ctx, approvalmd.InsertDetailReqDTO{
			FlowId:    flow.Id,
			FlowIndex: flow.CurrIndex,
			FlowOp:    approvalmd.DisagreeOp,
			Account:   reqDTO.Operator.Account,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			committer.Rollback()
			return util.InternalError(err)
		}
		_, err = approvalmd.UpdateNotifyDone(ctx, reqDTO.NotifyId, true)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			committer.Rollback()
			return util.InternalError(err)
		}
		committer.Commit()
	}
	count, err := approvalmd.CountDetail(ctx, flow.Id, flow.CurrIndex, approvalmd.DisagreeOp)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if int(count)+node.AtLeastNum >= len(node.Accounts) {
		flowCtx := &approval.FlowContext{
			Context: context.Background(),
			FlowId:  flow.Id,
			Process: &process,
			BizId:   flow.BizId,
		}
		go runNext(flowCtx, &flow, node, map[string]any{
			"agree": "n",
		})
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
		Pid:     reqDTO.Pid,
		Name:    reqDTO.Name,
		GroupId: reqDTO.GroupId,
		Process: reqDTO.Process.Convert(),
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
		return nil, err
	}
	// groupId = 1是系统审批流 不是自定义审批流
	if !b || process.GroupId == 1 {
		return nil, util.InvalidArgsError()
	}
	errKeys := process.Process.CheckKvCfgs(reqDTO.Kvs)
	if len(errKeys) > 0 {
		return errKeys, nil
	}
	flow, err := approvalmd.InsertFlow(ctx, approvalmd.InsertFlowReqDTO{
		ProcessId:  process.Id,
		Process:    process.Process,
		CurrIndex:  1,
		FlowStatus: approvalmd.PendingFlowStatus,
		Creator:    reqDTO.Operator.Account,
		BizId:      reqDTO.BizId,
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
	_, err = approvalmd.UpdateFlowStatusWithOldStatus(ctx, reqDTO.FlowId, approvalmd.CanceledFlowStatus, approvalmd.PendingFlowStatus)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}
