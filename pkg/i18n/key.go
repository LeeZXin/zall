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
	RepoPermsContainsTargetRepoId = KeyItem{
		Id:         "repo.permsContainsTargetRepoId",
		DefaultRet: "该项目组仍包含该仓库的特殊权限配置",
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

var (
	TimerTaskPendingStatus = KeyItem{
		Id:         "timerTask.pending",
		DefaultRet: "等待执行",
	}
	TimerTaskRunningStatus = KeyItem{
		Id:         "timerTask.running",
		DefaultRet: "正在执行",
	}
	TimerTaskSuccessfulStatus = KeyItem{
		Id:         "timerTask.successful",
		DefaultRet: "执行成功",
	}
	TimerTaskFailedStatus = KeyItem{
		Id:         "timerTask.failed",
		DefaultRet: "执行失败",
	}
	TimerTaskClosedStatus = KeyItem{
		Id:         "timerTask.pending",
		DefaultRet: "已关闭",
	}
	TimerTaskUnknownStatus = KeyItem{
		Id:         "timerTask.pending",
		DefaultRet: "未知",
	}
	TimerTaskAutoTriggerType = KeyItem{
		Id:         "timerTask.autoTriggerType",
		DefaultRet: "自动触发",
	}
	TimerTaskManualTriggerType = KeyItem{
		Id:         "timerTask.manualTriggerType",
		DefaultRet: "手动触发",
	}
	TimerTaskUnknownTriggerType = KeyItem{
		Id:         "timerTask.unknownTriggerType",
		DefaultRet: "未知",
	}
)

var (
	AppPermsContainerTargetAppId = KeyItem{
		Id:         "app.permsContainerTargetAppId",
		DefaultRet: "该项目组仍包含该app的特殊权限配置",
	}
	AppRemainCountGreaterThanZero = KeyItem{
		Id:         "app.remainCountGreaterThanZero",
		DefaultRet: "app数量仍大于0",
	}
)

var (
	FlowPendingStatus = KeyItem{
		Id:         "flowStatus.pending",
		DefaultRet: "执行中",
	}
	FlowAgreeStatus = KeyItem{
		Id:         "flowStatus.agree",
		DefaultRet: "同意",
	}
	FlowDisagreeStatus = KeyItem{
		Id:         "flowStatus.disagree",
		DefaultRet: "不同意",
	}
	FlowErrStatus = KeyItem{
		Id:         "flowStatus.err",
		DefaultRet: "出现错误",
	}
	FlowCanceledStatus = KeyItem{
		Id:         "flowStatus.canceled",
		DefaultRet: "已取消",
	}
	FlowUnknownStatus = KeyItem{
		Id:         "flowStatus.unknown",
		DefaultRet: "未知",
	}
	FlowPendingOp = KeyItem{
		Id:         "flowOp.pending",
		DefaultRet: "等待执行",
	}
	FlowAgreeOp = KeyItem{
		Id:         "flowOp.agree",
		DefaultRet: "同意",
	}
	FlowDisagreeOp = KeyItem{
		Id:         "flowOp.disagree",
		DefaultRet: "不同意",
	}
	FlowCancelOp = KeyItem{
		Id:         "flowOp.cancel",
		DefaultRet: "取消",
	}
	FlowAutoAgreeOp = KeyItem{
		Id:         "flowOp.autoAgree",
		DefaultRet: "自动同意",
	}
	FlowAutoDisagreeOp = KeyItem{
		Id:         "flowOp.autoDisagree",
		DefaultRet: "自动不同意",
	}
	FlowUnknownOp = KeyItem{
		Id:         "flowOp.unknown",
		DefaultRet: "未知",
	}
)
