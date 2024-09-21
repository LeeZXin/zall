package i18n

type KeyItem struct {
	Id         string
	DefaultRet string
}

var (
	SystemInternalError = KeyItem{
		Id:         "system.internalErr",
		DefaultRet: "系统异常",
	}
	SystemInvalidArgs = KeyItem{
		Id:         "system.invalidArgs",
		DefaultRet: "参数错误",
	}
	SystemNotLogin = KeyItem{
		Id:         "system.notLogin",
		DefaultRet: "未登录",
	}
	SystemUnauthorized = KeyItem{
		Id:         "system.unauthorized",
		DefaultRet: "权限不足",
	}
	SystemForbidden = KeyItem{
		Id:         "system.forbidden",
		DefaultRet: "禁止操作",
	}
	SystemProxyAbnormal = KeyItem{
		Id:         "system.proxyAbnormal",
		DefaultRet: "代理执行异常",
	}
	SystemAlreadyExists = KeyItem{
		Id:         "system.dataAlreadyExists",
		DefaultRet: "数据已存在",
	}
	SystemNotExists = KeyItem{
		Id:         "system.dataNotExists",
		DefaultRet: "数据不存在",
	}
	SystemUnimplemented = KeyItem{
		Id:         "system.unimplemented",
		DefaultRet: "未实现",
	}
	SystemOperationFailed = KeyItem{
		Id:         "system.operationFailed",
		DefaultRet: "操作失败",
	}
	SystemTooManyOperation = KeyItem{
		Id:         "system.tooManyOperation",
		DefaultRet: "太多人操作了，请稍后",
	}
	SystemHasBug = KeyItem{
		Id:         "system.hasBug",
		DefaultRet: "系统数据异常",
	}
)

var (
	UserNotFound = KeyItem{
		Id:         "user.notFound",
		DefaultRet: "用户不存在",
	}
	UserWrongPassword = KeyItem{
		Id:         "user.wrongPassword",
		DefaultRet: "密码不正确",
	}
	UserWrongOriginPassword = KeyItem{
		Id:         "user.wrongOriginPassword",
		DefaultRet: "原密码不正确",
	}
	UserAlreadyExists = KeyItem{
		Id:         "user.alreadyExists",
		DefaultRet: "用户已存在",
	}
	UserEmailAlreadyExists = KeyItem{
		Id:         "user.emailAlreadyExists",
		DefaultRet: "用户邮箱已存在",
	}
)

var (
	AvatarNotImageError = KeyItem{
		Id:         "avatar.notImage",
		DefaultRet: "不是图片",
	}
)

var (
	SshKeyFormatError = KeyItem{
		Id:         "sshKey.formatError",
		DefaultRet: "ssh公钥格式错误",
	}
	SshKeyAlreadyExists = KeyItem{
		Id:         "sshKey.alreadyExists",
		DefaultRet: "ssh公钥已存在",
	}
)

var (
	CronExpError = KeyItem{
		Id:         "cronExp.valueError",
		DefaultRet: "cron表达式错误",
	}
)

var (
	GpgKeyFormatError = KeyItem{
		Id:         "gpgKey.formatError",
		DefaultRet: "gpg公钥格式错误",
	}
	GpgKeyAlreadyExists = KeyItem{
		Id:         "gpg.alreadyExists",
		DefaultRet: "gpg公钥已存在",
	}
)

var (
	TeamAdminUserGroupName = KeyItem{
		Id:         "team.adminUserGroupName",
		DefaultRet: "项目组管理员",
	}
)

var (
	RepoAlreadyExists = KeyItem{
		Id:         "repo.alreadyExists",
		DefaultRet: "仓库已存在",
	}
	RepoSizeExceedLimit = KeyItem{
		Id:         "repo.sizeExceedLimit",
		DefaultRet: "大小超过上限%s",
	}
	RepoRemainCountGreaterThanZero = KeyItem{
		Id:         "repo.remainCountGreaterThanZero",
		DefaultRet: "仓库数量仍大于0",
	}
)

var (
	PullRequestCannotMerge = KeyItem{
		Id:         "pullRequest.cannotMerge",
		DefaultRet: "无法合并",
	}
	PullRequestAlreadyExists = KeyItem{
		Id:         "pullRequest.alreadyExists",
		DefaultRet: "合并请求已存在",
	}
	PullRequestReviewerCountLowerThanCfg = KeyItem{
		Id:         "pullRequest.reviewerCountLowerThanCfg",
		DefaultRet: "代码评审数量小于配置数量",
	}
)

var (
	ProtectedBranchNotAllowForcePush = KeyItem{
		Id:         "protectedBranch.notAllowForcePush",
		DefaultRet: "保护分支禁止强制push",
	}
	ProtectedBranchNotAllowDelete = KeyItem{
		Id:         "protectedBranch.notAllowDelete",
		DefaultRet: "保护分支不可删除",
	}
	ProtectedBranchNotAllowPush = KeyItem{
		Id:         "protectedBranch.notAllowPush",
		DefaultRet: "保护分支不可推送",
	}
)

var (
	SshCmdNotSupported = KeyItem{
		Id:         "ssh.notSupportedCmd",
		DefaultRet: "不支持该命令",
	}
)

var (
	InvalidWorkflowContent = KeyItem{
		Id:         "action.invalidContent",
		DefaultRet: "action yaml不合法",
	}
)

var (
	TimerTaskRemainCountGreaterThanZero = KeyItem{
		Id:         "timerTask.remainCountGreaterThanZero",
		DefaultRet: "定时任务数量仍大于0",
	}
)

var (
	AppRemainCountGreaterThanZero = KeyItem{
		Id:         "app.remainCountGreaterThanZero",
		DefaultRet: "app数量仍大于0",
	}
)

var (
	SqlWrongSyntaxMsg = KeyItem{
		Id:         "sql.wrongSyntaxMsg",
		DefaultRet: "sql语法错误",
	}
	SqlUnsupportedMsg = KeyItem{
		Id:         "sql.unsupportedMsg",
		DefaultRet: "不支持该sql",
	}
	SqlNotAllowHasLimitMsg = KeyItem{
		Id:         "sql.notAllowHasLimit",
		DefaultRet: "sql不允许带limit",
	}
	SqlNotAllowNoWhereMsg = KeyItem{
		Id:         "sql.notAllowNoWhere",
		DefaultRet: "sql没有带where关键词",
	}
)
