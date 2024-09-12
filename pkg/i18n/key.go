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
	SshKeyVerifyGuide = KeyItem{
		Id:         "sshKey.verifyGuide",
		DefaultRet: "请在十分钟内执行以下命令，输出签名并提交",
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
	PullRequestAgreeReviewStatus = KeyItem{
		Id:         "pullRequest.agreeMerge",
		DefaultRet: "同意合并",
	}
	PullRequestCanceledReviewStatus = KeyItem{
		Id:         "pullRequest.canceledMerge",
		DefaultRet: "被撤销",
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
	ReviewAlreadyExists = KeyItem{
		Id:         "review.alreadyExists",
		DefaultRet: "已经评审过",
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

var (
	WorkflowHookTriggerType = KeyItem{
		Id:         "action.hookTriggerType",
		DefaultRet: "webhook触发",
	}
	WorkflowManualTriggerType = KeyItem{
		Id:         "action.manualTriggerType",
		DefaultRet: "手动触发",
	}
	WorkflowUnknownTriggerType = KeyItem{
		Id:         "action.unknownTriggerType",
		DefaultRet: "未知类型触发",
	}
)

var (
	ServiceAbnormalStatus = KeyItem{
		Id:         "service.abnormalStatus",
		DefaultRet: "异常",
	}
	ServiceStartingStatus = KeyItem{
		Id:         "service.startingStatus",
		DefaultRet: "启动中",
	}
	ServiceStartedStatus = KeyItem{
		Id:         "service.startedStatus",
		DefaultRet: "正常",
	}
	ServiceStoppingStatus = KeyItem{
		Id:         "service.stoppingStatus",
		DefaultRet: "下线中",
	}
	ServiceStoppedStatus = KeyItem{
		Id:         "service.stoppedStatus",
		DefaultRet: "已下线",
	}
	ServiceUnknownStatus = KeyItem{
		Id:         "service.unknownStatus",
		DefaultRet: "未知状态",
	}
	ServiceRestartOp = KeyItem{
		Id:         "service.restartOp",
		DefaultRet: "重启服务",
	}
	ServiceStopOp = KeyItem{
		Id:         "service.stopOp",
		DefaultRet: "下线服务",
	}
	ServiceUnknownOp = KeyItem{
		Id:         "service.unknownOp",
		DefaultRet: "未知类型",
	}
	ServiceAddServiceAfterPlanCreatingType = KeyItem{
		Id:         "service.addServiceAfterPlanCreatingType",
		DefaultRet: "常规发布",
	}
	ServiceAddServiceBeforePlanCreatingType = KeyItem{
		Id:         "service.addServiceBeforePlanCreatingType",
		DefaultRet: "临时发布",
	}
	ServiceUnknownPlanType = KeyItem{
		Id:         "service.unknownPlanType",
		DefaultRet: "未知类型",
	}
	ServiceRunningPlanStatus = KeyItem{
		Id:         "service.runningPlanStatus",
		DefaultRet: "发布中",
	}
	ServiceClosedPlanStatus = KeyItem{
		Id:         "service.closedPlanStatus",
		DefaultRet: "已关闭",
	}
	ServiceUnknownPlanStatus = KeyItem{
		Id:         "service.unknownPlanStatus",
		DefaultRet: "已关闭",
	}
	ServiceWaitPlanItemStatus = KeyItem{
		Id:         "service.waitPlanItemStatus",
		DefaultRet: "等待部署",
	}
	ServiceDeployedPlanItemStatus = KeyItem{
		Id:         "service.deployedPlanItemStatus",
		DefaultRet: "已部署",
	}
	ServiceRollbackPlanItemStatus = KeyItem{
		Id:         "service.rollbackPlanItemStatus",
		DefaultRet: "已回滚",
	}
	ServiceClosedPlanItemStatus = KeyItem{
		Id:         "service.closedPlanItemStatus",
		DefaultRet: "未知状态",
	}
	ServiceUnknownPlanItemStatus = KeyItem{
		Id:         "service.unknownPlanItemStatus",
		DefaultRet: "未知状态",
	}
)

var (
	PromScrapeDiscoveryTargetType = KeyItem{
		Id:         "promScrape.discoveryTargetType",
		DefaultRet: "服务发现类型",
	}
	PromScrapeHostTargetType = KeyItem{
		Id:         "promScrape.hostTargetType",
		DefaultRet: "主机服务类型",
	}
	PromScrapeUnknownTargetType = KeyItem{
		Id:         "promScrape.unknownTargetType",
		DefaultRet: "未知类型",
	}
)

var (
	DbReadPermType = KeyItem{
		Id:         "db.readPermType",
		DefaultRet: "读权限",
	}
	DbWritePermType = KeyItem{
		Id:         "db.writePermType",
		DefaultRet: "写权限",
	}
	DbReadWritePermType = KeyItem{
		Id:         "db.readPermType",
		DefaultRet: "读写权限",
	}
	DbUnknownPermType = KeyItem{
		Id:         "db.unknownPermType",
		DefaultRet: "未知类型",
	}
	DbPendingPermOrderStatus = KeyItem{
		Id:         "db.pendingPermOrderStatus",
		DefaultRet: "审批中",
	}
	DbAgreePermOrderStatus = KeyItem{
		Id:         "db.agreePermOrderStatus",
		DefaultRet: "同意",
	}
	DbCanceledPermOrderStatus = KeyItem{
		Id:         "db.canceledPermOrderStatus",
		DefaultRet: "已取消",
	}
	DbDisagreePermOrderStatus = KeyItem{
		Id:         "db.disagreePermOrderStatus",
		DefaultRet: "不同意",
	}
	DbUnknownPermOrderStatus = KeyItem{
		Id:         "db.unknownPermOrderStatus",
		DefaultRet: "未知状态",
	}
	DbPendingUpdateOrderStatus = KeyItem{
		Id:         "db.pendingUpdateOrderStatus",
		DefaultRet: "审批中",
	}
	DbAgreeUpdateOrderStatus = KeyItem{
		Id:         "db.agreeUpdateOrderStatus",
		DefaultRet: "同意",
	}
	DbCanceledUpdateOrderStatus = KeyItem{
		Id:         "db.canceledUpdateOrderStatus",
		DefaultRet: "已取消",
	}
	DbDisagreeUpdateOrderStatus = KeyItem{
		Id:         "db.disagreeUpdateOrderStatus",
		DefaultRet: "不同意",
	}
	DbExecutedUpdateOrderStatus = KeyItem{
		Id:         "db.executedUpdateOrderStatus",
		DefaultRet: "执行完成",
	}
	DbUnknownUpdateOrderStatus = KeyItem{
		Id:         "db.unknownUpdateOrderStatus",
		DefaultRet: "未知状态",
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

var (
	TimerAutoTriggerType = KeyItem{
		Id:         "timerTask.autoTriggerType",
		DefaultRet: "自动触发",
	}
	TimerManualTriggerType = KeyItem{
		Id:         "timerTask.manualTriggerType",
		DefaultRet: "手动触发",
	}
)
