package i18n

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
	SystemAlreadyExists = KeyItem{
		Id:         "system.dataAlreadyExists",
		DefaultRet: "数据已存在",
	}
	SystemUnimplemented = KeyItem{
		Id:         "system.unimplemented",
		DefaultRet: "未实现",
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
	UserAlreadyExists = KeyItem{
		Id:         "user.alreadyExists",
		DefaultRet: "用户已存在",
	}
	UserAccountNotFoundWarnFormat = KeyItem{
		Id:         "user.notFoundWarnFormat",
		DefaultRet: "用户%s不存在",
	}
	UserAccountUnauthorizedReviewCodeWarnFormat = KeyItem{
		Id:         "user.notFoundWarnFormat",
		DefaultRet: "该用户%s无评审代码的权限",
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
	SshKeyVerifyGuide = KeyItem{
		Id:         "sshKey.verifyGuide",
		DefaultRet: "请在十分钟内执行以下命令，输出签名并提交",
	}
)

var (
	EmptyGitNodesError = KeyItem{
		Id:         "gitNodes.emptyError",
		DefaultRet: "可用存储节点为空",
	}
)

var (
	GpgKeyFormatError = KeyItem{
		Id:         "gpgKey.formatError",
		DefaultRet: "gpg公钥格式错误",
	}
	GpgKeyVerifiedFailedError = KeyItem{
		Id:         "gpgKey.verifiedFailedError",
		DefaultRet: "gpg校验失败",
	}
	GpgTokenExpiredError = KeyItem{
		Id:         "gpgKey.tokenExpiredError",
		DefaultRet: "token过期",
	}
	GpgVerifyGuide = KeyItem{
		Id:         "gpgKey.verifyGuide",
		DefaultRet: "请在十分钟内执行以下命令，输出签名并提交",
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
	TeamUserGroupHasUserWhenDel = KeyItem{
		Id:         "team.userGroupHasUserWhenDel",
		DefaultRet: "该项目组存在关联用户, 无法删除",
	}
	TeamUserGroupUpdateAdminNotAllow = KeyItem{
		Id:         "team.updateTeamUserAdminGroupNotAllow",
		DefaultRet: "不允许编辑项目组管理员权限",
	}
)

var (
	RepoAlreadyExists = KeyItem{
		Id:         "repo.alreadyExists",
		DefaultRet: "仓库已存在",
	}
	RepoSizeExceedLimit = KeyItem{
		Id:         "repo.sizeExceedLimit",
		DefaultRet: "仓库大小超过上限%s",
	}
	RepoRemainCountGreaterThanZero = KeyItem{
		Id:         "repo.remainCountGreaterThanZero",
		DefaultRet: "仓库数量仍大于0",
	}
)

var (
	TimeBeforeSecondUnit = KeyItem{
		Id:         "time.beforeSecondUnit",
		DefaultRet: "秒前",
	}
	TimeBeforeMinuteUnit = KeyItem{
		Id:         "time.beforeMinuteUnit",
		DefaultRet: "分钟前",
	}
	TimeBeforeHourUnit = KeyItem{
		Id:         "time.beforeHourUnit",
		DefaultRet: "小时前",
	}
	TimeBeforeDayUnit = KeyItem{
		Id:         "time.beforeDdayUnit",
		DefaultRet: "天前",
	}
	TimeBeforeMonthUnit = KeyItem{
		Id:         "time.beforeMonthUnit",
		DefaultRet: "月前",
	}
	TimeBeforeYearUnit = KeyItem{
		Id:         "time.beforeYearUnit",
		DefaultRet: "年前",
	}
)

var (
	RepoOpenStatus = KeyItem{
		Id:         "repo.openStatus",
		DefaultRet: "打开状态",
	}
	RepoClosedStatus = KeyItem{
		Id:         "repo.closedStatus",
		DefaultRet: "归档状态",
	}
	RepoDeletedStatus = KeyItem{
		Id:         "repo.deletedStatus",
		DefaultRet: "删除状态",
	}
	RepoUnknownStatus = KeyItem{
		Id:         "repo.unknownStatus",
		DefaultRet: "未知状态",
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
	PullRequestOpenStatus = KeyItem{
		Id:         "pullRequest.openStatus",
		DefaultRet: "已打开",
	}
	PullRequestClosedStatus = KeyItem{
		Id:         "pullRequest.closedStatus",
		DefaultRet: "已关闭",
	}
	PullRequestMergedStatus = KeyItem{
		Id:         "pullRequest.mergedStatus",
		DefaultRet: "已合并",
	}
	PullRequestUnknownStatus = KeyItem{
		Id:         "pullRequest.unknownStatus",
		DefaultRet: "未知",
	}
	PullRequestMergeMessage = KeyItem{
		Id:         "pullRequest.mergeMessage",
		DefaultRet: "合并请求: %v, 申请人: %s, 合并人: %s",
	}
	PullRequestAgreeMergeStatus = KeyItem{
		Id:         "pullRequest.agreeMerge",
		DefaultRet: "同意合并",
	}
	PullRequestDisagreeMergeStatus = KeyItem{
		Id:         "pullRequest.disagreeMerge",
		DefaultRet: "不同意合并",
	}
	PullRequestUnknownReviewStatus = KeyItem{
		Id:         "pullRequest.unknownReviewStatus",
		DefaultRet: "未知状态",
	}
	PullRequestReviewerCountLowerThanCfg = KeyItem{
		Id:         "pullRequest.reviewerCountLowerThanCfg",
		DefaultRet: "代码评审数量小于配置数量",
	}
)

var (
	SshKeyVerifyTokenExpired = KeyItem{
		Id:         "sshKey.verifyTokenExpired",
		DefaultRet: "token已失效",
	}
	SshKeyVerifyFailed = KeyItem{
		Id:         "sshKey.verifyFailed",
		DefaultRet: "校验失败",
	}
)

var (
	ProtectedBranchInvalidReviewCountWhenCreatePr = KeyItem{
		Id:         "protectedBranch.invalidReviewCountWhenCreatePr",
		DefaultRet: "保护分支代码评审者数量不合法",
	}
	ProtectedBranchNotAllowForcePush = KeyItem{
		Id:         "protectedBranch.notAllowForcePush",
		DefaultRet: "保护分支禁止强制push",
	}
	ProtectedBranchNotAllowDelete = KeyItem{
		Id:         "protectedBranch.notAllowDelete",
		DefaultRet: "保护分支禁止删除",
	}
	ProtectedBranchNotAllowDirectPush = KeyItem{
		Id:         "protectedBranch.notAllowDirectPush",
		DefaultRet: "保护分支不可直接push",
	}
)

var (
	LfsNotSupported = KeyItem{
		Id:         "lfs.notSupported",
		DefaultRet: "不支持lfs",
	}
)

var (
	SshCmdNotSupported = KeyItem{
		Id:         "ssh.notSupportedCmd",
		DefaultRet: "不支持该命令",
	}
)

var (
	ReviewAlreadyExists = KeyItem{
		Id:         "review.alreadyExists",
		DefaultRet: "已经评审过",
	}
)

var (
	InvalidActionContent = KeyItem{
		Id:         "action.invalidContent",
		DefaultRet: "action yaml不合法",
	}
	ActionInstanceNotFound = KeyItem{
		Id:         "action.instanceNotFound",
		DefaultRet: "无可用执行节点",
	}
)
