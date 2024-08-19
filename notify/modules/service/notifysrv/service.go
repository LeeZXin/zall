package notifysrv

import (
	"context"
	"github.com/LeeZXin/zsf-utils/executor"
	"time"
)

var (
	Outer OuterService
	Inner InnerService
)

func Init() {
	InitOuter()
	InitInner()
	if sendExecutor == nil {
		sendExecutor, _ = executor.NewExecutor(10, 1024, time.Minute, executor.AbortStrategy)
	}
}

func InitOuter() {
	if Outer == nil {
		Outer = new(outerImpl)
	}
	initSendExecutor()
}

func InitInner() {
	if Inner == nil {
		Inner = new(innerImpl)
	}
	initSendExecutor()
}

func initSendExecutor() {
	if sendExecutor == nil {
		sendExecutor, _ = executor.NewExecutor(10, 1024, time.Minute, executor.AbortStrategy)
	}
}

type InnerService interface {
	// SendNotificationByTplId 通过模板id发送通知
	SendNotificationByTplId(context.Context, int64, map[string]string) error
}

type OuterService interface {
	// CreateTpl 创建通知模板
	CreateTpl(context.Context, CreateTplReqDTO) error
	// UpdateTpl 编辑通知模板
	UpdateTpl(context.Context, UpdateTplReqDTO) error
	// DeleteTpl 删除通知模板
	DeleteTpl(context.Context, DeleteTplReqDTO) error
	// ListTpl 通知模板列表
	ListTpl(context.Context, ListTplReqDTO) ([]TplDTO, int64, error)
	// ChangeTplApiKey 更换api key
	ChangeTplApiKey(context.Context, ChangeTplApiKeyReqDTO) error
	// SendNotificationByApiKey 通过api key发送通知
	SendNotificationByApiKey(context.Context, SendNotifyByApiKeyReqDTO) error
	// ListAllTpl 通过团队获取模板列表
	ListAllTpl(context.Context, ListAllTplReqDTO) ([]SimpleTplDTO, error)
}
