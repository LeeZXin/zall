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
	TeamUserGroupHasUserWhenDelCode
	CannotUpdateTeamUserAdminGroupCode
	InvalidReviewCountWhenCreatePrCode
	ForcePushForbiddenCode
	MethodUnImplementedCode
	ActionInstanceNotFoundCode
	ThereHasBugErrCode
	EmptyGitNodesErrCode
	OperationFailedErrCode
)

func (c Code) Int() int {
	return int(c)
}
