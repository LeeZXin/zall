package util

import (
	"fmt"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/wraperr"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

func InvalidArgsError() error {
	return NewBizErr(apicode.InvalidArgsCode, i18n.SystemInvalidArgs)
}

func InternalError(err error) error {
	return wraperr.WrapErr(err, NewBizErr(apicode.InternalErrorCode, i18n.SystemInternalError))
}

func ApiError(resp ginutil.BaseResp) error {
	return bizerr.NewBizErr(resp.Code, resp.Message)
}

func UnauthorizedError() error {
	return NewBizErr(apicode.UnauthorizedCode, i18n.SystemUnauthorized)
}

func OperationFailedError() error {
	return NewBizErr(apicode.OperationFailedErrCode, i18n.SystemOperationFailed)
}

func AlreadyExistsError() error {
	return NewBizErr(apicode.DataAlreadyExistsCode, i18n.SystemAlreadyExists)
}

func NotExistsError() error {
	return NewBizErr(apicode.DataNotExistsCode, i18n.SystemNotExists)
}

func MethodUnimplementedErr() error {
	return NewBizErr(apicode.MethodUnImplementedCode, i18n.SystemUnimplemented)
}

func NewBizErr(code apicode.Code, key i18n.KeyItem, msg ...string) *bizerr.Err {
	if len(msg) == 0 {
		return bizerr.NewBizErr(code.Int(), i18n.GetByKey(key))
	}
	return bizerr.NewBizErr(code.Int(), fmt.Sprintf(i18n.GetByKey(key), msg))
}

func NewBizErrWithMsg(code apicode.Code, msg string) *bizerr.Err {
	if len(msg) == 0 {
		return bizerr.NewBizErr(code.Int(), msg)
	}
	return bizerr.NewBizErr(code.Int(), msg)
}

func ThereHasBugErr() error {
	return NewBizErr(apicode.ThereHasBugErrCode, i18n.SystemInternalError)
}
