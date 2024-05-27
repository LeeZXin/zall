package oplogsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/pkg/i18n"
)

var (
	Inner              = newInnerService()
	Outer OuterService = new(outerImpl)
)

type InnerService interface {
	InsertOpLog(OpLog)
}

type OuterService interface {
	PageOpLog(context.Context, PageOpLogReqDTO) ([]OpLogDTO, int64, error)
}

func FormatI18n(key i18n.KeyItem, str ...any) string {
	return fmt.Sprintf(i18n.GetByKey(key), str...)
}
