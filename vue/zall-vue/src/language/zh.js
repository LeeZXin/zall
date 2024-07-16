export default {
    app: "ZALL",
    logoutText: "退出登录",
    footer: "Build with zsf",
    createTeamText: "创建团队",
    myTeam: "我的团队",
    switchTeam: "切换团队",
    notFound: "页面丢失了",
    backToIndex: "回到首页",
    previousPage: "上一页",
    nextPage: "下一页",
    ztable: {
        noDataText: "暂无数据",
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
        timerTask: "定时任务",
        applyApproval: "审批申请",
        dbAudit: "数据库审计",
        promScrape: "Prometheus监控",
        roleAndMembers: "角色&成员"
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
        manualTriggerType: "手动触发"
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
    }
}