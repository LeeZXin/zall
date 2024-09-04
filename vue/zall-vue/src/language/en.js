export default {
    app: "ZALL",
    logoutText: "logout",
    footer: "build with zsf",
    createTeamText: "Create team",
    myTeam: "My Teams",
    switchTeam: "switch team",
    notFound: "The page is lost",
    backToIndex: "Back to index page",
    noData: "No data found",
    superAdmin: "Super Admin",
    ztable: {
        noDataText: "No data found",
    },
    operationSuccess: "Successfully",
    operationFail: "Failed",
    appMenu: {
        propertyFile: "Property file",
        deployPipeline: "Deploy pipeline",
        deployPlan: "Deploy plan",
        serviceStatus: "Service status",
        registryCenter: "Registry center",
        artifacts: "Artifacts",
        promScrape: "Prom Scrape",
        alertConfig: "Alert config",
        setting: "Setting"
    },
    gitRepoMenu: {
        index: "Repo code",
        pullRequest: "PullRequest",
        branch: "Branch",
        tag: "Tag",
        commitHistory: "Commit history",
        workflow: "Workflow",
        protectedBranch: "Protected branch",
        webhook: "Webhook",
        setting: "Setting"
    },
    indexMenu: {
        team: "Team Collaborate",
        mysqlAudit: "Mysql Audit"
    },
    mysqlAuditMenu: {
        databaseSource: "Database source",
        readPermApply: "Apply read perm",
        readPermAudit: "Audit read perm",
        readPermList: "Read perm List",
        readPermManage: "Manage read perm",
        dataUpdateApply: "Apply data update",
        dataUpdateAudit: "Audit data update",
        dataSearch: "Search data"
    },
    personalSettingMenu: {
        profile: "Profile",
        changePassword: "Change password",
        sshAndGpg: "SSH&GPG"
    },
    saMenu: {
        sysCfg: "System config",
        userManage: "Manage users",
        propertyCenterSource: "Property source",
        serviceStatusSource: "Service source",
        registryCenterSource: "Registry source",
        zallet: "Zallet",
        promScrape: "Prom scrape"
    },
    system: {
        requestFailed: "something wrong happened, please retry in a while",
        notLogin: "you do not login, will back to login page",
        request404: "can not find request, please retry in a while",
        request403: "you are unauthorized to access",
        request400: "bad request, check your request",
        internalError: "internal error"
    },
    login: {
        title: "Login",
        accountPlaceholder: "please input account",
        passwordPlaceholder: "please input password",
        loginBtnText: "Login Now",
        registerText: "Register account",
        pleaseConfirmAccount: "account length should should between 4 and 32",
        pleaseConfirmPassword: "password length should greater than 5"
    },
    register: {
        title: "Register",
        accountPlaceholder: "please input account",
        passwordPlaceholder: "please input password",
        usernamePlaceholder: "please input username",
        emailPlaceholder: "please input email",
        confirmPasswordPlaceholder: "please confirm password",
        registerBtnText: "Register Now",
        backToLoginText: "Back to login",
        pleaseConfirmAccount: "account length should should between 4 and 32",
        pleaseConfirmPassword: "password length should greater than 5",
        pleaseConfirmUsername: "username length should should between 0 and 32",
        pleaseConfirmEmail: "email length should greater than 0",
        pleaseConfirmConfirmPassword: "confirm password should be equal with password"
    },
    teamMenu: {
        gitRepo: "Git repo",
        app: "App service",
        action: "Action",
        propertyCenter: "Property center",
        timer: "Timer task",
        dbAudit: "Database audit",
        promScrape: "Prometheus scrape",
        roleAndMembers: "Role&Member",
        setting: "Team setting",
        notifyTpl: "Notification",
        teamHook: "Team hook",
        weworkAccessToken: "Wework token",
        feishuAccessToken: "Feishu token"
    },
    roleListPage: {
        roleList: "Role List",
        roleName: "Role Name"
    },
    settings: {
        sshAndGpg: "SSH and GPG keys"
    },
    gitRepo: {
        searchText: "Search Repository",
        createRepoText: "Create Repository",
        switchRepo: "Switch Repository"
    },
    appService: {
        switchApp: "Switch App"
    },
    createGitRepo: {
        backToRepoList: "Back to Git repo list",
        createText: "Create a new Repository",
        starText: "Item marked with * is required",
        owner: "Owner",
        team: "Team",
        storeNode: "Storage node",
        repoName: "Repository name* (excluding special characters, the length must not exceed 32)",
        repoDesc: "Description(add brief description for your repository, the length must not exceed 255)",
        gitignore: "Gitignore template(let git ignores some files)",
        defaultBranch: "Default branch(the length must not exceed 32)",
        addReadme: "Add Readme",
        createBtn: "Create"
    },
    createTeam: {
        teamName: "Team name* (excluding special characters, the length must not exceed 32)"
    },
    workflow: {
        manualTriggerType: "Manual",
        hookTriggerType: "Automatic",
        status: {
            success: "Successful",
            fail: "Failed",
            cancel: "Cancelled",
            timeout: "Timeout",
            running: "Running",
            queue: "Queued",
            unknown: "Unknown",
            unExecuted: "UnExecuted"
        }
    },
    secondBefore: " second ago",
    minuteBefore: " minute ago",
    hourBefore: " hour ago",
    dayBefore: " day ago",
    monthBefore: " month ago",
    yearBefore: " year ago",
    yes: "Yes",
    no: "No",
    timerTask: {
        autoTriggerType: "Automatically",
        manualTriggerType: "Manually",
        successful: "Success",
        failed: "Failed",
        logColumns: {
            triggerBy: "Trigger by",
            triggerType: "Type",
            isSuccess: "Status",
            created: "Execute Time",
            operation: "Operation"
        },
        viewErrLog: "View Error Log",
        viewTaskCfg: "View Configuration",
        errLog: "Error Log",
        taskCfg: "Task Configuration",
        executeLog: "Execute Log",
        searchMonthly: "Search monthly"
    },
    deployPlan: {
        pendingStatus: "Pending",
        runningStatus: "Running",
        successfulStatus: "Successful",
        closedStatus: "Closed"
    },
    mysqlReadPermApply: {
        pendingStatus: "Pending",
        agreeStatus: "Agree",
        disagreeStatus: "Disagree",
        canceledStatus: "Canceled",
        unknownStatus: "Unknown"
    },
    mysqlDataUpdateApply: {
        title: "Apply for data update",
        selectHost: "Select host",
        fillAccessBaseName: "Fill access base name",
        fillApplyReason: "Fill reason for application",
        fillSql: "Fill sql",
        executeImmediatelyAfterApproval: "Execute immediately after approval",
        apply: "Apply",
        pendingStatus: "Pending",
        agreeStatus: "Agree",
        disagreeStatus: "Disagree",
        canceledStatus: "Canceled",
        askToExecuteStatus: "Ask To Execute",
        executedStatus: "Executed",
        unknownStatus: "Unknown",
        yes: "Yes",
        no: "No",
        cancel: "Cancel",
        viewExplain: "View explain",
        viewSql: "View sql",
        viewLog: "View log",
        askToExecute: "Ask to execute",
        dbName: "Database",
        accessBase: "Access base",
        applyStatus: "Status",
        applyReason: "Reason",
        executeImmediatelyAfterApprovalCol: "Execute immediately after approval",
        created: "Application time",
        account: "Applicant",
        operation: "Operation",
        auditor: "Auditor",
        auditTime: "Audit time",
        disagreeReason: "Disagree reason",
        cancelTime: "Cancel time",
        applyTime: "Application time",
        executor: "Executor",
        executeTime: "Execute time",
        accessBaseFormatErr: "Access base format is wrong",
        sqlFormatErr: "Sql format is wrong",
        applyReasonFormatErr: "Reason format is wrong",
        disagreeReasonFormatErr: "Disagree reason format is wrong",
        confirmCancel: "Do you confirm to cancel",
        confirmAskToExecute: "Do you confirm to ask to execute",
        confirmAgree: "Do you confirm to agree",
        confirmExecute: "Do you confirm to execute",
        agree: "Agree",
        disagree: "Disagree",
        allDatabases: "All databases",
        executeApply: "Execute",
        fillDisagreeReason: "Fill disagree reason"
    }
}