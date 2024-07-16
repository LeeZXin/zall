package apicode

type Code int

const (
	InternalErrorCode Code = 99999
)

const (
	InvalidArgsCode Code = iota + 40000
	BadRequestCode
	DataNotExistsCode
	DataAlreadyExistsCode
	WrongLoginPasswordCode
	NotLoginCode
	UnauthorizedCode
	UserAlreadyExistsCode
	PullRequestCannotMergeCode
	ForcePushForbiddenCode
	ProtectedBranchNotAllowPushCode
	MethodUnImplementedCode
	ThereHasBugErrCode
	OperationFailedErrCode
	TooManyOperationCode
	OutOfWorkflowCapacityErrCode
	ProxyAbnormalErrCode
)

func (c Code) Int() int {
	return int(c)
}
