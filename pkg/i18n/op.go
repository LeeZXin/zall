package i18n

var (
	TeamSrvKeysVO = teamSrvKeys{
		CreateTeam: KeyItem{
			Id:         "teamSrv.CreateTeam",
			DefaultRet: "新增项目组",
		},
		DeleteTeam: KeyItem{
			Id:         "teamSrv.DeleteTeam",
			DefaultRet: "删除项目组",
		},
		UpdateTeam: KeyItem{
			Id:         "teamSrv.UpdateTeam",
			DefaultRet: "编辑项目组",
		},
		DeleteTeamUser: KeyItem{
			Id:         "teamSrv.DeleteUser",
			DefaultRet: "删除项目组用户",
		},
		UpsertTeamUser: KeyItem{
			Id:         "teamSrv.CreateUser",
			DefaultRet: "新增或编辑项目组用户",
		},
		InsertTeamUserGroup: KeyItem{
			Id:         "teamSrv.CreateRole",
			DefaultRet: "添加项目组用户组",
		},
		UpdateTeamUserGroupName: KeyItem{
			Id:         "teamSrv.UpdateRoleName",
			DefaultRet: "编辑项目组用户组名称",
		},
		UpdateTeamUserGroupPerm: KeyItem{
			Id:         "teamSrv.UpdateRole",
			DefaultRet: "编辑项目组用户组权限",
		},
		DeleteTeamUserGroup: KeyItem{
			Id:         "teamSrv.DeleteRole",
			DefaultRet: "删除项目组用户组",
		},
	}

	UserSrvKeysVO = UserSrvKeys{
		Login: KeyItem{
			Id:         "userSrv.Login",
			DefaultRet: "登录",
		},
		LoginOut: KeyItem{
			Id:         "userSrv.Logout",
			DefaultRet: "登出",
		},
		InsertUser: KeyItem{
			Id:         "userSrv.InsertUser",
			DefaultRet: "新增用户",
		},
		RegisterUser: KeyItem{
			Id:         "userSrv.RegisterUser",
			DefaultRet: "注册用户",
		},
		DeleteUser: KeyItem{
			Id:         "userSrv.DeleteUser",
			DefaultRet: "删除用户",
		},
		UpdateUser: KeyItem{
			Id:         "userSrv.UpdateUser",
			DefaultRet: "编辑用户",
		},
		UpdateAdmin: KeyItem{
			Id:         "userSrv.UpdateAdmin",
			DefaultRet: "编辑系统管理员",
		},
		UpdatePassword: KeyItem{
			Id:         "userSrv.UpdatePassword",
			DefaultRet: "编辑密码",
		},
		SetProhibited: KeyItem{
			Id:         "userSrv.SetProhibited",
			DefaultRet: "禁用用户",
		},
	}

	SshKeySrvKeysVO = SshKeySrvKeys{
		DeleteSshKey: KeyItem{
			Id:         "sshKeySrv.DeleteSshKey",
			DefaultRet: "删除ssh key",
		},
		InsertSshKey: KeyItem{
			Id:         "sshKeySrv.InsertSshKey",
			DefaultRet: "添加ssh key",
		},
	}

	RepoSrvKeysVO = RepoSrvKeys{
		CreateRepo: KeyItem{
			Id:         "repoSrv.CreateRepo",
			DefaultRet: "新建仓库",
		},
		DeleteRepo: KeyItem{
			Id:         "repoSrv.DeleteRepo",
			DefaultRet: "删除仓库",
		},
		InsertRepoToken: KeyItem{
			Id:         "repoSrv.InsertRepoToken",
			DefaultRet: "创建访问令牌",
		},
		DeleteRepoToken: KeyItem{
			Id:         "repoSrv.DeleteRepoToken",
			DefaultRet: "删除访问令牌",
		},
		AccessRepo: KeyItem{
			Id:         "repoSrv.AccessRepo",
			DefaultRet: "拉取代码",
		},
		PushRepo: KeyItem{
			Id:         "repoSrv.PushRepo",
			DefaultRet: "提交代码",
		},
		RefreshAllGitHooks: KeyItem{
			Id:         "repoSrv.RefreshAllGitHooks",
			DefaultRet: "刷新所有git hook",
		},
		TransferTeam: KeyItem{
			Id:         "repoSrv.TransferTeam",
			DefaultRet: "迁移team",
		},
	}

	WorkflowSrvKeysVO = ActionSrvKeys{
		CreateWorkflow: KeyItem{
			Id:         "repoSrv.CreateWorkflow",
			DefaultRet: "新增action",
		},
		DeleteWorkflow: KeyItem{
			Id:         "repoSrv.DeleteWorkflow",
			DefaultRet: "删除action",
		},
		UpdateWorkflow: KeyItem{
			Id:         "repoSrv.UpdateWorkflow",
			DefaultRet: "编辑action",
		},
		TriggerAction: KeyItem{
			Id:         "repoSrv.TriggerWorkflow",
			DefaultRet: "手动触发action",
		},
	}

	PullRequestSrvKeysVO = PullRequestSrvKeys{
		SubmitPullRequest: KeyItem{
			Id:         "pullRequestSrv.SubmitPullRequest",
			DefaultRet: "提交合并请求: #%v",
		},
		ClosePullRequest: KeyItem{
			Id:         "pullRequestSrv.ClosePullRequest",
			DefaultRet: "关闭合并请求: #%v",
		},
		MergePullRequest: KeyItem{
			Id:         "pullRequestSrv.MergePullRequest",
			DefaultRet: "触发合并请求: #%v",
		},
		ReviewPullRequest: KeyItem{
			Id:         "pullRequestSrv.AgreeReviewPullRequest",
			DefaultRet: "评审合并请求: #%v",
		},
		AddComment: KeyItem{
			Id:         "pullRequestSrv.AddComment",
			DefaultRet: "添加评论: %v",
		},
		DeleteComment: KeyItem{
			Id:         "pullRequestSrv.DeleteComment",
			DefaultRet: "删除评论: %v",
		},
	}

	CfgSrvKeysVO = CfgSrvKeys{
		UpdateSysCfg: KeyItem{
			Id:         "cfgSrv.UpdateSysCfg",
			DefaultRet: "编辑系统配置",
		},
		UpdateGitCfg: KeyItem{
			Id:         "cfgSrv.UpdateGitCfg",
			DefaultRet: "编辑git配置",
		},
		UpdateGitRepoServerUrl: KeyItem{
			Id:         "cfgSrv.UpdateGitRepoServerCfg",
			DefaultRet: "编辑git仓库服务器地址",
		},
	}

	BranchSrvKeysVO = BranchSrvKeys{
		CreateProtectedBranch: KeyItem{
			Id:         "branchSrv.CreateProtectedBranch",
			DefaultRet: "新增保护分支: %s",
		},
		DeleteProtectedBranch: KeyItem{
			Id:         "branchSrv.DeleteProtectedBranch",
			DefaultRet: "删除保护分支: %s",
		},
		UpdateProtectedBranch: KeyItem{
			Id:         "branchSrv.UpdateProtectedBranch",
			DefaultRet: "编辑保护分支: %s",
		},
	}

	HookSrvKeysVO = HookSrvKeys{
		PreReceive: KeyItem{
			Id:         "hookSrv.PreReceive",
			DefaultRet: "提交代码",
		},
	}

	LfsSrvKeysVO = LfsSrvKeys{
		Download: KeyItem{
			Id:         "lfsSrv.download",
			DefaultRet: "下载lfs文件: %s",
		},
		Upload: KeyItem{
			Id:         "lfsSrv.upload",
			DefaultRet: "上传lfs文件: %s",
		},
	}

	WebhookSrvKeysVO = WebhookSrvKeys{
		InsertWebhook: KeyItem{
			Id:         "webhookSrv.CreateWebhook",
			DefaultRet: "新增webhook",
		},
		UpdateWebhook: KeyItem{
			Id:         "webhookSrv.UpdateWebhook",
			DefaultRet: "编辑webhook",
		},
		DeleteWebhook: KeyItem{
			Id:         "webhookSrv.DeleteWebhook",
			DefaultRet: "删除webhook",
		},
	}

	TimerTaskSrvKeysVO = TimerTaskSrvKeys{
		InsertTask: KeyItem{
			Id:         "timerTaskSrv.CreateTask",
			DefaultRet: "新增定时任务",
		},
		UpdateTask: KeyItem{
			Id:         "timerTaskSrv.UpdateTask",
			DefaultRet: "更新定时任务",
		},
		EnableTask: KeyItem{
			Id:         "timerTaskSrv.EnableTask",
			DefaultRet: "启动定时任务",
		},
		DisableTask: KeyItem{
			Id:         "timerTaskSrv.DisableTask",
			DefaultRet: "关闭定时任务",
		},
		DeleteTask: KeyItem{
			Id:         "timerTaskSrv.DeleteTask",
			DefaultRet: "删除定时任务",
		},
		TriggerTask: KeyItem{
			Id:         "timerTaskSrv.TriggerTask",
			DefaultRet: "手动执行任务",
		},
	}

	GpgSrvKeysVO = GpgSrvKeys{
		InsertGpgKey: KeyItem{
			Id:         "gpgSrv.InsertGpgKey",
			DefaultRet: "新增gpg公钥",
		},
		DeleteGpgKey: KeyItem{
			Id:         "gpgSrv.DeleteGpgKey",
			DefaultRet: "删除gpg公钥",
		},
	}

	PropSrvKeysVO = PropSrvKeys{
		GrantAuth: KeyItem{
			Id:         "propSrv.GrantAuth",
			DefaultRet: "编辑app etcd权限",
		},
		InsertEtcdNode: KeyItem{
			Id:         "propSrv.InsertEtcdNode",
			DefaultRet: "新增prop etcd节点",
		},
		DeleteEtcdNode: KeyItem{
			Id:         "propSrv.DeleteEtcdNode",
			DefaultRet: "删除prop etcd节点",
		},
		UpdateEtcdNode: KeyItem{
			Id:         "propSrv.UpdateEtcdNode",
			DefaultRet: "编辑prop etcd节点",
		},
		InsertPropContent: KeyItem{
			Id:         "propSrv.InsertPropContent",
			DefaultRet: "新增配置",
		},
		UpdatePropContent: KeyItem{
			Id:         "propSrv.UpdatePropContent",
			DefaultRet: "编辑配置",
		},
		DeletePropContent: KeyItem{
			Id:         "propSrv.DeletePropContent",
			DefaultRet: "删除配置",
		},
		DeployPropContent: KeyItem{
			Id:         "propSrv.DeployPropContent",
			DefaultRet: "发布配置",
		},
	}

	AppSrvKeysVO = AppSrvKeys{
		InsertApp: KeyItem{
			Id:         "appSrv.InsertApp",
			DefaultRet: "新增app",
		},
		DeleteApp: KeyItem{
			Id:         "appSrv.DeleteApp",
			DefaultRet: "删除app",
		},
		TransferTeam: KeyItem{
			Id:         "appSrv.TransferTeam",
			DefaultRet: "迁移项目组",
		},
		UpdateApp: KeyItem{
			Id:         "appSrv.UpdateApp",
			DefaultRet: "编辑app",
		},
	}

	TcpDetectSrvKeysVO = TcpDetectSrvKeys{
		InsertDetect: KeyItem{
			Id:         "tcpDetectSrv.InsertDetect",
			DefaultRet: "新增tcp探测",
		},
		DeleteDetect: KeyItem{
			Id:         "tcpDetectSrv.DeleteDetect",
			DefaultRet: "删除tcp探测",
		},
		UpdateDetect: KeyItem{
			Id:         "tcpDetectSrv.UpdateDetect",
			DefaultRet: "编辑tcp探测",
		},
		EnableDetect: KeyItem{
			Id:         "tcpDetectSrv.EnableDetect",
			DefaultRet: "启动tcp探测",
		},
		DisableDetect: KeyItem{
			Id:         "tcpDetectSrv.DisableDetect",
			DefaultRet: "关闭tcp探测",
		},
	}

	ApprovalSrvKeysVO = ApprovalSrvKeys{
		InsertCustomProcess: KeyItem{
			Id:         "approvalSrv.InsertCustomProcess",
			DefaultRet: "新增自定义审批流",
		},
		UpdateCustomProcess: KeyItem{
			Id:         "approvalSrv.UpdateCustomProcess",
			DefaultRet: "编辑自定义审批流",
		},
		DeleteCustomProcess: KeyItem{
			Id:         "approvalSrv.DeleteCustomProcess",
			DefaultRet: "删除自定义审批流",
		},
		InsertGroup: KeyItem{
			Id:         "approvalSrv.InsertGroup",
			DefaultRet: "新增审批流分组",
		},
		DeleteGroup: KeyItem{
			Id:         "approvalSrv.DeleteGroup",
			DefaultRet: "删除审批流分组",
		},
		UpdateGroup: KeyItem{
			Id:         "approvalSrv.UpdateGroup",
			DefaultRet: "编辑审批流分组",
		},
	}

	GitNodeSrvKeysVO = GitNodeSrvKeys{
		InsertNode: KeyItem{
			Id:         "repoSrv.InsertNode",
			DefaultRet: "新增git节点信息",
		},
		DeleteNode: KeyItem{
			Id:         "repoSrv.DeleteNode",
			DefaultRet: "删除git节点信息",
		},
		UpdateNode: KeyItem{
			Id:         "repoSrv.UpdateNode",
			DefaultRet: "编辑git节点信息",
		},
	}

	DeploySrvKeysVO = DeploySrvKeys{
		UpdateConfig: KeyItem{
			Id:         "deploySrv.UpdateConfig",
			DefaultRet: "编辑部署配置",
		},
		InsertConfig: KeyItem{
			Id:         "deploySrv.InsertConfig",
			DefaultRet: "新增部署配置",
		},
		InsertPlan: KeyItem{
			Id:         "deploySrv.InsertPlan",
			DefaultRet: "创建发布计划",
		},
		ClosePlan: KeyItem{
			Id:         "deploySrv.ClosePlan",
			DefaultRet: "关闭发布计划",
		},
		InsertPlanItem: KeyItem{
			Id:         "deploySrv.InsertPlanItem",
			DefaultRet: "添加计划服务",
		},
		ClosePlanItem: KeyItem{
			Id:         "deploySrv.ClosePlanItem",
			DefaultRet: "关闭计划服务",
		},
		ReDeployService: KeyItem{
			Id:         "deploySrv.DeployService",
			DefaultRet: "重建服务",
		},
		StopService: KeyItem{
			Id:         "deploySrv.StopService",
			DefaultRet: "下线服务",
		},
		RestartService: KeyItem{
			Id:         "deploySrv.RestartService",
			DefaultRet: "重启服务",
		},
		DeployServiceWithPlan: KeyItem{
			Id:         "deploySrv.DeployServiceWithPlan",
			DefaultRet: "通过发布计划部署服务",
		},
		RollbackServiceWithPlan: KeyItem{
			Id:         "deploySrv.RollbackServiceWithPlan",
			DefaultRet: "通过发布计划回滚服务",
		},
	}

	PromSrvKeysVO = PromSrvKeys{
		InsertScrape: KeyItem{
			Id:         "promSrv.InsertScrape",
			DefaultRet: "新增prometheus抓取任务",
		},
		UpdateScrape: KeyItem{
			Id:         "promSrv.UpdateScrape",
			DefaultRet: "编辑prometheus抓取任务",
		},
		DeleteScrape: KeyItem{
			Id:         "promSrv.DeleteScrape",
			DefaultRet: "删除prometheus抓取任务",
		},
	}

	DbSrvKeysVO = DbSrvKeys{
		InsertDb: KeyItem{
			Id:         "dbSrv.InsertDb",
			DefaultRet: "新增数据库",
		},
		DeleteDb: KeyItem{
			Id:         "dbSrv.DeleteDb",
			DefaultRet: "删除数据库",
		},
		UpdateDb: KeyItem{
			Id:         "dbSrv.UpdateDb",
			DefaultRet: "编辑数据库",
		},
		DeletePerm: KeyItem{
			Id:         "dbSrv.DeletePerm",
			DefaultRet: "删除数据库权限",
		},
		SearchDb: KeyItem{
			Id:         "dbSrv.SearchDb",
			DefaultRet: "搜索数据库",
		},
	}

	AlertSrvKeysVO = AlertSrvKeys{
		InsertConfig: KeyItem{
			Id:         "alertSrv.InsertConfig",
			DefaultRet: "新增告警配置",
		},
		UpdateConfig: KeyItem{
			Id:         "alertSrv.UpdateConfig",
			DefaultRet: "编辑告警配置",
		},
		DeleteConfig: KeyItem{
			Id:         "alertSrv.InsertConfig",
			DefaultRet: "删除告警配置",
		},
	}
)

type KeyItem struct {
	Id         string
	DefaultRet string
}

type teamSrvKeys struct {
	CreateTeam              KeyItem
	DeleteTeam              KeyItem
	UpdateTeam              KeyItem
	DeleteTeamUser          KeyItem
	UpsertTeamUser          KeyItem
	InsertTeamUserGroup     KeyItem
	UpdateTeamUserGroupName KeyItem
	UpdateTeamUserGroupPerm KeyItem
	DeleteTeamUserGroup     KeyItem
}

type UserSrvKeys struct {
	Login          KeyItem
	LoginOut       KeyItem
	InsertUser     KeyItem
	RegisterUser   KeyItem
	DeleteUser     KeyItem
	UpdateUser     KeyItem
	ListUser       KeyItem
	UpdateAdmin    KeyItem
	UpdatePassword KeyItem
	SetProhibited  KeyItem
}

type SshKeySrvKeys struct {
	DeleteSshKey KeyItem
	InsertSshKey KeyItem
}

type RepoSrvKeys struct {
	CreateRepo         KeyItem
	DeleteRepo         KeyItem
	InsertRepoToken    KeyItem
	DeleteRepoToken    KeyItem
	AccessRepo         KeyItem
	PushRepo           KeyItem
	RefreshAllGitHooks KeyItem
	TransferTeam       KeyItem
}

type ActionSrvKeys struct {
	CreateWorkflow KeyItem
	DeleteWorkflow KeyItem
	UpdateWorkflow KeyItem
	TriggerAction  KeyItem
}

type PullRequestSrvKeys struct {
	SubmitPullRequest KeyItem
	ClosePullRequest  KeyItem
	MergePullRequest  KeyItem
	ReviewPullRequest KeyItem
	AddComment        KeyItem
	DeleteComment     KeyItem
}

type CfgSrvKeys struct {
	UpdateSysCfg           KeyItem
	UpdateGitCfg           KeyItem
	UpdateGitRepoServerUrl KeyItem
}

type BranchSrvKeys struct {
	CreateProtectedBranch KeyItem
	DeleteProtectedBranch KeyItem
	UpdateProtectedBranch KeyItem
}

type HookSrvKeys struct {
	PreReceive KeyItem
}

type LfsSrvKeys struct {
	Download KeyItem
	Upload   KeyItem
}

type WebhookSrvKeys struct {
	InsertWebhook KeyItem
	UpdateWebhook KeyItem
	DeleteWebhook KeyItem
}

type GpgSrvKeys struct {
	InsertGpgKey KeyItem
	DeleteGpgKey KeyItem
}

type TimerTaskSrvKeys struct {
	InsertTask  KeyItem
	UpdateTask  KeyItem
	EnableTask  KeyItem
	DisableTask KeyItem
	DeleteTask  KeyItem
	TriggerTask KeyItem
}

type PropSrvKeys struct {
	GrantAuth         KeyItem
	InsertEtcdNode    KeyItem
	DeleteEtcdNode    KeyItem
	UpdateEtcdNode    KeyItem
	InsertPropContent KeyItem
	UpdatePropContent KeyItem
	DeletePropContent KeyItem
	DeployPropContent KeyItem
}

type AppSrvKeys struct {
	InsertApp    KeyItem
	DeleteApp    KeyItem
	UpdateApp    KeyItem
	TransferTeam KeyItem
}

type TcpDetectSrvKeys struct {
	InsertDetect  KeyItem
	DeleteDetect  KeyItem
	UpdateDetect  KeyItem
	EnableDetect  KeyItem
	DisableDetect KeyItem
}

type ApprovalSrvKeys struct {
	InsertCustomProcess KeyItem
	UpdateCustomProcess KeyItem
	DeleteCustomProcess KeyItem
	InsertGroup         KeyItem
	DeleteGroup         KeyItem
	UpdateGroup         KeyItem
}

type GitNodeSrvKeys struct {
	InsertNode KeyItem
	DeleteNode KeyItem
	UpdateNode KeyItem
}

type DeploySrvKeys struct {
	UpdateConfig            KeyItem
	InsertConfig            KeyItem
	InsertPlan              KeyItem
	ClosePlan               KeyItem
	InsertPlanItem          KeyItem
	ClosePlanItem           KeyItem
	ReDeployService         KeyItem
	StopService             KeyItem
	RestartService          KeyItem
	DeployServiceWithPlan   KeyItem
	RollbackServiceWithPlan KeyItem
}

type PromSrvKeys struct {
	InsertScrape KeyItem
	UpdateScrape KeyItem
	DeleteScrape KeyItem
}

type DbSrvKeys struct {
	InsertDb   KeyItem
	DeleteDb   KeyItem
	UpdateDb   KeyItem
	DeletePerm KeyItem
	SearchDb   KeyItem
}

type AlertSrvKeys struct {
	InsertConfig KeyItem
	UpdateConfig KeyItem
	DeleteConfig KeyItem
}
