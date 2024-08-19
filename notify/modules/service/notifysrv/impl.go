package notifysrv

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/notify/modules/model/notifymd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/notify/feishu"
	"github.com/LeeZXin/zall/pkg/notify/notify"
	"github.com/LeeZXin/zall/pkg/notify/wework"
	"github.com/LeeZXin/zall/timer/modules/model/taskmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"text/template"
)

var (
	sendExecutor *executor.Executor
)

type innerImpl struct{}

// SendNotificationByTplId 通过模板id发送通知
func (*innerImpl) SendNotificationByTplId(ctx context.Context, tplId int64, params map[string]string) error {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	tpl, b, err := notifymd.GetTplById(ctx, tplId, nil)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return sendNotificationWithNotifyTpl(tpl, params)
}

type outerImpl struct{}

// CreateTpl 创建通知模板
func (*outerImpl) CreateTpl(ctx context.Context, reqDTO CreateTplReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkManageNotifyTplPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return err
	}
	err = notifymd.InsertTpl(ctx, notifymd.InsertTplReqDTO{
		Name:      reqDTO.Name,
		ApiKey:    idutil.RandomUuid(),
		NotifyCfg: reqDTO.Cfg,
		TeamId:    reqDTO.TeamId,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateTpl 编辑通知模板
func (*outerImpl) UpdateTpl(ctx context.Context, reqDTO UpdateTplReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkManageNotifyTplPermByTplId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	_, err = notifymd.UpdateTpl(ctx, notifymd.UpdateTplReqDTO{
		Id:        reqDTO.Id,
		Name:      reqDTO.Name,
		NotifyCfg: reqDTO.Cfg,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteTpl 删除通知模板
func (*outerImpl) DeleteTpl(ctx context.Context, reqDTO DeleteTplReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkManageNotifyTplPermByTplId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除模板
		_, err2 := notifymd.DeleteTpl(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		// 删除失败定时任务通知绑定
		return taskmd.DeleteFailedTaskNotifyTplByTplId(ctx, reqDTO.Id)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListTpl 通知模板列表
func (*outerImpl) ListTpl(ctx context.Context, reqDTO ListTplReqDTO) ([]TplDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkManageNotifyTplPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return nil, 0, err
	}
	tpls, total, err := notifymd.ListTpl(ctx, notifymd.ListTplReqDTO{
		Name:     reqDTO.Name,
		PageNum:  reqDTO.PageNum,
		TeamId:   reqDTO.TeamId,
		PageSize: 10,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret, _ := listutil.Map(tpls, func(t notifymd.Tpl) (TplDTO, error) {
		ret := TplDTO{
			Id:     t.Id,
			Name:   t.Name,
			ApiKey: t.ApiKey,
			TeamId: t.TeamId,
		}
		if t.NotifyCfg != nil {
			ret.NotifyCfg = t.NotifyCfg.Data
		}
		return ret, nil
	})
	return ret, total, nil
}

// ChangeTplApiKey 更换api key
func (*outerImpl) ChangeTplApiKey(ctx context.Context, reqDTO ChangeTplApiKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkManageNotifyTplPermByTplId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	_, err = notifymd.UpdateTplApiKeyById(ctx, reqDTO.Id, idutil.RandomUuid())
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// SendNotificationByApiKey 通过api key发送通知
func (*outerImpl) SendNotificationByApiKey(ctx context.Context, reqDTO SendNotifyByApiKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	tpl, b, err := notifymd.GetTplByApiKey(ctx, reqDTO.ApiKey)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return sendNotificationWithNotifyTpl(tpl, reqDTO.Params)
}

// ListAllTpl 通过团队获取模板列表
func (*outerImpl) ListAllTpl(ctx context.Context, reqDTO ListAllTplReqDTO) ([]SimpleTplDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 只要属于同一个团队的都可以看到数据
	if !reqDTO.Operator.IsAdmin {
		_, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		if !b {
			return nil, util.UnauthorizedError()
		}
	}
	tpls, err := notifymd.ListAllTplByTeamId(ctx, reqDTO.TeamId, []string{"id", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(tpls, func(t notifymd.Tpl) (SimpleTplDTO, error) {
		return SimpleTplDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
}

func checkManageNotifyTplPermByTeamId(ctx context.Context, operator apisession.UserInfo, teamId int64) error {
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanManageNotifyTpl {
		return nil
	}
	return util.UnauthorizedError()
}

func checkManageNotifyTplPermByTplId(ctx context.Context, operator apisession.UserInfo, tplId int64) error {
	tpl, b, err := notifymd.GetTplById(ctx, tplId, []string{"team_id"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkManageNotifyTplPermByTeamId(ctx, operator, tpl.TeamId)
}

func sendNotificationWithNotifyTpl(tpl notifymd.Tpl, params map[string]string) error {
	if tpl.NotifyCfg == nil {
		return util.InvalidArgsError()
	}
	if !tpl.NotifyCfg.Data.IsValid() {
		return util.ThereHasBugErr()
	}
	cfg := tpl.NotifyCfg.Data
	notificationTpl, err := template.New("").Parse(cfg.Template)
	if err != nil {
		// 模板错误
		return util.OperationFailedError()
	}
	msg := new(bytes.Buffer)
	err = notificationTpl.Execute(msg, params)
	if err != nil {
		// 参数错误
		return util.OperationFailedError()
	}
	err = sendNotification(msg.Bytes(), cfg)
	if err != nil {
		return util.OperationFailedError()
	}
	return nil
}

// 发送通知
func sendNotification(msgBytes []byte, cfg notify.Cfg) error {
	switch cfg.NotifyType {
	case notify.Wework:
		var msg wework.Message
		err := json.Unmarshal(msgBytes, &msg)
		if err != nil {
			return err
		}
		if err = msg.IsValid(); err != nil {
			return err
		}
		return sendExecutor.Execute(func() {
			wework.SendMessage(cfg.Url, msg)
		})
	case notify.Feishu:
		var msg feishu.Message
		err := json.Unmarshal(msgBytes, &msg)
		if err != nil {
			return err
		}
		if err = msg.IsValid(); err != nil {
			return err
		}
		return sendExecutor.Execute(func() {
			feishu.SendMessage(cfg.Url, cfg.FeishuSignKey, msg)
		})
	}
	return nil
}
