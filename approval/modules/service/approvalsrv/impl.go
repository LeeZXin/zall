package approvalsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/approval/approval"
	"github.com/LeeZXin/zall/approval/modules/model/approvalmd"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/hashicorp/go-bexpr"
)

type innerImpl struct {
}

func (*innerImpl) InsertFlow(ctx context.Context, pid string) error {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	process, b, err := approvalmd.GetProcessByPid(ctx, pid)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("process not found: %s", pid)
	}
	flow, err := approvalmd.InsertFlow(ctx, approvalmd.InsertFlowReqDTO{
		ProcessId:  process.Id,
		Approval:   process.Approval,
		CurrIndex:  1,
		FlowStatus: approvalmd.PendingFlowStatus,
		Creator:    "system",
	})
	if err != nil {
		return err
	}
	executeFlow(flow)
	return nil
}

func (*innerImpl) InsertProcess(ctx context.Context, reqDTO InsertProcessReqDTO) error {
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
		Pid:      reqDTO.Pid,
		Name:     reqDTO.Name,
		Approval: *reqDTO.Approval.ToApproval(),
	})
}

func (*innerImpl) UpdateProcess(ctx context.Context, reqDTO UpdateProcessReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := approvalmd.UpdateProcessByPid(ctx, approvalmd.UpdateProcessByPidReqDTO{
		Pid:      reqDTO.Pid,
		Name:     reqDTO.Name,
		Approval: *reqDTO.Approval.ToApproval(),
	})
	return err
}

func executeFlow(flow approvalmd.Flow) {
	go func() {
		p, err := flow.GetApproval()
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		runFlow(&flow, &p)
	}()
}

func runFlow(flow *approvalmd.Flow, p *approval.Approval) {
	p.FindAndDo(flow.CurrIndex, map[approval.NodeType]func(*approval.Approval){
		approval.PeopleNode: func(p *approval.Approval) {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			err := xormstore.WithTx(ctx, func(ctx context.Context) error {
				return approvalmd.InsertNotify(ctx, approvalmd.InsertNotifyReqDTO{
					FlowId:    flow.Id,
					Accounts:  p.Accounts,
					Done:      false,
					FlowIndex: flow.CurrIndex,
				})
			})
			if err != nil {
				logger.Logger.Error(err)
			}
		},
		approval.ApiNode: func(p *approval.Approval) {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			if p.Api == nil {
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
			response, err := p.Api.DoRequest()
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
			runNext(flow, p, response)
		},
		approval.MethodNode: func(p *approval.Approval) {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			if p.Method == nil {
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
			response, err := p.Method.DoMethod()
			if err != nil {
				logger.Logger.Errorf("flowId: %d currIndex: %d its method: %s err: %v", flow.Id, flow.CurrIndex, p.Method.Name, err)
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
			runNext(flow, p, response)
		},
		approval.DisagreeNode: func(p *approval.Approval) {
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
		approval.AgreeNode: func(p *approval.Approval) {
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

func runNext(flow *approvalmd.Flow, p *approval.Approval, response map[string]any) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	if len(p.Next) == 0 {
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
	}
	for _, c := range p.Next {
		if doCondition(c.Condition, response) {
			// 更新currIndex
			b, err := approvalmd.UpdateFlowCurrIndexWithOldCurrIndex(ctx, flow.Id, c.Node.NodeId, flow.CurrIndex)
			if err != nil {
				logger.Logger.Error(err)
				return
			}
			if b {
				flow.CurrIndex = c.Node.NodeId
				runFlow(flow, c.Node)
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
	p, err := flow.GetApproval()
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	n := p.Find(notify.FlowIndex)
	if n == nil {
		return util.InvalidArgsError()
	}
	if n.NodeType != approval.PeopleNode {
		return util.InvalidArgsError()
	}
	findAccount := false
	for _, account := range n.Accounts {
		if account == reqDTO.Operator.Account {
			findAccount = true
			break
		}
	}
	if !findAccount {
		return util.UnauthorizedError()
	}
	b, err = approvalmd.ExistsDetailByAccount(ctx, flow.Id, flow.CurrIndex, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	err = approvalmd.InsertDetail(ctx, approvalmd.InsertDetailReqDTO{
		FlowId:    flow.Id,
		FlowIndex: flow.CurrIndex,
		FlowOp:    approvalmd.AgreeOp,
		Account:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	count, err := approvalmd.CountDetail(ctx, flow.Id, flow.CurrIndex, approvalmd.AgreeOp)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if int(count) >= p.AtLeastNum {
		go runNext(&flow, n, map[string]any{
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
	p, err := flow.GetApproval()
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	n := p.Find(flow.CurrIndex)
	if n == nil {
		return util.InvalidArgsError()
	}
	if n.NodeType != approval.PeopleNode {
		return util.InvalidArgsError()
	}
	findAccount := false
	for _, account := range n.Accounts {
		if account == reqDTO.Operator.Account {
			findAccount = true
			break
		}
	}
	if !findAccount {
		return util.UnauthorizedError()
	}
	b, err = approvalmd.ExistsDetailByAccount(ctx, flow.Id, flow.CurrIndex, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	err = approvalmd.InsertDetail(ctx, approvalmd.InsertDetailReqDTO{
		FlowId:    flow.Id,
		FlowIndex: flow.CurrIndex,
		FlowOp:    approvalmd.DisagreeOp,
		Account:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	count, err := approvalmd.CountDetail(ctx, flow.Id, flow.CurrIndex, approvalmd.DisagreeOp)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if int(count)+p.AtLeastNum >= len(p.Accounts) {
		go runNext(&flow, n, map[string]any{
			"agree": "n",
		})
	}
	return nil
}
