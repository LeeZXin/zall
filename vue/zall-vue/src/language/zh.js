export default {
    app: "ZALL",
    logoutText: "退出登录",
    footer: "Build with zsf",
    createTeamText: "创建团队",
    myTeam: "我的团队",
    switchTeam: "切换团队",
    notFound: "页面丢失了",
    backToIndex: "回到首页",
    noData: "暂无数据",
    superAdmin: "超级管理员",
    ztable: {
        noDataText: "暂无数据",
    },
    appMenu: {
        propertyFile: "配置文件",
        deployPipeline: "部署流水线",
        deployPlan: "发布计划",
        serviceStatus: "服务状态",
        registryCenter: "注册中心",
        artifacts: "制品库",
        promScrape: "Prom抓取任务",
        alertConfig: "告警配置",
        setting: "设置"
    },
    gitRepoMenu: {
        index: "代码文件",
        pullRequest: "合并请求",
        branch: "分支",
        tag: "标签",
        commitHistory: "提交历史",
        workflow: "工作流",
        protectedBranch: "保护分支",
        webhook: "Webhook",
        setting: "设置"
    },
    indexMenu: {
        team: "团队协作",
        mysqlAudit: "Mysql审计"
    },
    mysqlAuditMenu: {
        databaseSource: "数据源",
        readPermApply: "读权限申请",
        readPermAudit: "读权限审批",
        readPermList: "读权限列表",
        readPermManage: "读权限管理",
        dataUpdateApply: "数据修改",
        dataUpdateAudit: "数据修改审批",
        dataSearch: "数据查询"
    },
    personalSettingMenu: {
        profile: "个人信息",
        changePassword: "修改密码",
        sshAndGpg: "SSH&GPG"
    },
    saMenu: {
        sysCfg: "系统配置",
        userManage: "用户管理",
        propertyCenterSource: "配置中心来源",
        serviceStatusSource: "服务状态来源",
        registryCenterSource: "注册中心来源",
        zallet: "Zallet",
        promScrape: "Prom抓取任务"
    },
    system: {
        requestFailed: "请求发送似乎有点问题, 请稍后重试",
        notLogin: "未登录, 即将跳转登录页面",
        request404: "服务器似乎有点问题, 请稍后重试",
        request403: "暂时无权限访问",
        request400: "请求格式错误",
        internalError: "请求发送似乎有点问题, 请稍后重试"
    },
    login: {
        title: "登录",
        accountPlaceholder: "请输入账号",
        passwordPlaceholder: "请输入密码",
        loginBtnText: "立即登录",
        registerText: "注册用户",
        pleaseConfirmAccount: "账号长度在4-32之间",
        pleaseConfirmPassword: "密码长度必须大于5"
    },
    register: {
        title: "注册",
        accountPlaceholder: "请输入账号",
        usernamePlaceholder: "请输入用户名",
        emailPlaceholder: "请输入邮箱",
        passwordPlaceholder: "请输入密码",
        confirmPasswordPlaceholder: "请再次输入密码",
        registerBtnText: "立即注册",
        backToLoginText: "返回登录",
        pleaseConfirmAccount: "账号长度在4-32之间",
        pleaseConfirmPassword: "密码长度必须大于5",
        pleaseConfirmUsername: "用户名长度在0-32之间",
        pleaseConfirmEmail: "请输入正确的邮箱",
        pleaseConfirmConfirmPassword: "确认密码和输入密码不一致"
    },
    teamMenu: {
        gitRepo: "代码仓库",
        app: "应用服务",
        action: "工作流",
        propertyCenter: "配置中心",
        timer: "定时任务",
        applyApproval: "审批申请",
        dbAudit: "数据库审计",
        promScrape: "Prometheus监控",
        roleAndMembers: "角色&成员",
        setting: "设置",
        notifyTpl: "外部通知",
        teamHook: "Team Hook"
    },
    roleListPage: {
        roleList: "角色列表",
        roleName: "角色名称"
    },
    settings: {
        sshAndGpg: "SSH和GPG密钥"
    },
    gitRepo: {
        searchText: "搜索仓库",
        createRepoText: "创建仓库",
        switchRepo: "切换仓库"
    },
    appService: {
        switchApp: "切换应用"
    },
    createGitRepo: {
        backToRepoList: "返回仓库列表",
        createText: "新建仓库",
        starText: "标记*的为必填项",
        owner: "拥有人",
        team: "所属团队",
        storeNode: "存储服务器",
        repoName: "仓库名称*(不包含特殊字符,长度不得超过32)",
        repoDesc: "仓库描述(为仓库添加一段简短的描述,长度不得超过255)",
        gitignore: "Gitignore模版(让git忽略某些文件)",
        defaultBranch: "默认分支(长度不得超过32)",
        addReadme: "添加Readme",
        createBtn: "新建"
    },
    createTeam: {
        teamName: "团队名称*(不包含特殊字符,长度不得超过32)"
    },
    workflow: {
        manualTriggerType: "手动触发",
        hookTriggerType: "自动触发",
        status: {
            success: "成功",
            fail: "失败",
            cancel: "中止",
            timeout: "超时",
            running: "运行中",
            queue: "排队中",
            unknown: "未知",
            unExecuted: "未执行"
        }
    },
    secondBefore: "秒前",
    minuteBefore: "分钟前",
    hourBefore: "小时前",
    dayBefore: "天前",
    monthBefore: "月前",
    yearBefore: "年前",
    yes: "是",
    no: "否",
    timerTask: {
        autoTriggerType: "自动触发",
        manualTriggerType: "手动触发",
        successful: "成功",
        failed: "失败",
        logColumns: {
            triggerBy: "操作人",
            triggerType: "操作类型",
            isSuccess: "状态",
            created: "执行时间",
            operation: "操作"
        },
        viewErrLog: "查看错误日志",
        viewTaskCfg: "查看任务配置",
        errLog: "错误日志",
        taskCfg: "任务配置",
        executeLog: "执行日志",
        searchMonthly: "搜索月份"
    },
    deployPlan: {
        pendingStatus: "待发布",
        runningStatus: "发布中",
        successfulStatus: "发布成功",
        closedStatus: "已关闭",
        unknownStatus: "未知"
    },
    mysqlReadPermApply: {
        pendingStatus: "等待审核",
        agreeStatus: "同意",
        disagreeStatus: "不同意",
        canceledStatus: "已取消",
        unknownStatus: "未知"
    },
    mysqlDataUpdateApply: {
        pendingStatus: "等待审批",
        agreeStatus: "同意",
        disagreeStatus: "不同意",
        canceledStatus: "已取消",
        askToExecuteStatus: "请求执行",
        executedStatus: "已执行",
        unknownStatus: "未知"
    }
}