import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from "@/pinia/userStore";
import { getUserInfoRequest } from "@/api/user/loginApi";

const router = createRouter({
    history: createWebHistory(),
    routes: [{
            path: "",
            redirect: "/index"
        },
        {
            path: "/index",
            redirect: "/index/team/list",
            component: () =>
                import ("../layouts/IndexLayout"),
            children: [{
                    path: "/index/team/list",
                    component: () =>
                        import ("../pages/team/team/TeamListPage")
                },
                {
                    path: "/index/team/create",
                    component: () =>
                        import ("../pages/team/team/TeamCreatePage")
                }
            ]
        },
        {
            path: "/sa",
            redirect: "/sa/cfg/list",
            component: () =>
                import ("../layouts/SaLayout"),
            children: [{
                path: "/sa/cfg/list",
                component: () =>
                    import ("../pages/sa/CfgListPage")
            }, {
                path: "/sa/propertySource/list/:env?",
                component: () =>
                    import ("../pages/sa/PropertySourceListPage")

            }, {
                path: "/sa/propertySource/create",
                component: () =>
                    import ("../pages/sa/PropertySourceHandlePage")

            }, {
                path: "/sa/propertySource/:sourceId(\\d+)/update",
                component: () =>
                    import ("../pages/sa/PropertySourceHandlePage")

            }, {
                path: "/sa/serviceSource/list/:env?",
                component: () =>
                    import ("../pages/sa/ServiceSourceListPage")

            }, {
                path: "/sa/serviceSource/create",
                component: () =>
                    import ("../pages/sa/ServiceSourceHandlePage")

            }, {
                path: "/sa/serviceSource/:sourceId(\\d+)/update",
                component: () =>
                    import ("../pages/sa/ServiceSourceHandlePage")

            }, {
                path: "/sa/discoverySource/list/:env?",
                component: () =>
                    import ("../pages/sa/DiscoverySourceListPage")

            }, {
                path: "/sa/discoverySource/create",
                component: () =>
                    import ("../pages/sa/DiscoverySourceHandlePage")

            }, {
                path: "/sa/discoverySource/:sourceId(\\d+)/update",
                component: () =>
                    import ("../pages/sa/DiscoverySourceHandlePage")

            }, {
                path: "/sa/user/list",
                component: () =>
                    import ("../pages/sa/UserManageListPage")
            }, {
                path: "/sa/user/create",
                component: () =>
                    import ("../pages/sa/UserManageHandlePage")
            }, {
                path: "/sa/user/:account(.+)/update",
                component: () =>
                    import ("../pages/sa/UserManageHandlePage")
            }, {
                path: "/sa/zalletNode/list/:env?",
                component: () =>
                    import ("../pages/sa/ZalletNodeListPage")
            }, {
                path: "/sa/zalletNode/create",
                component: () =>
                    import ("../pages/sa/ZalletNodeHandlePage")

            }, {
                path: "/sa/zalletNode/:nodeId(\\d+)/update",
                component: () =>
                    import ("../pages/sa/ZalletNodeHandlePage")

            }, {
                path: "/sa/promScrape/list/:env?",
                component: () =>
                    import ("../pages/sa/PromScrapeListPage")

            }, {
                path: "/sa/promScrape/create",
                component: () =>
                    import ("../pages/sa/PromScrapeHandlePage")

            }, {
                path: "/sa/promScrape/:scrapeId(\\d+)/update",
                component: () =>
                    import ("../pages/sa/PromScrapeHandlePage")

            }]
        },
        {
            path: "/personalSetting",
            redirect: "/personalSetting/profile/info",
            component: () =>
                import ("../layouts/PersonalSettingLayout"),
            children: [{
                path: "/personalSetting/profile/info",
                component: () =>
                    import ("../pages/personalSetting/ProfilePage")
            }, {
                path: "/personalSetting/changePassword",
                component: () =>
                    import ("../pages/personalSetting/ChangePasswordPage")
            }, {
                path: "/personalSetting/sshAndGpg/list",
                component: () =>
                    import ("../pages/personalSetting/SshAndGpgPage")
            }, {
                path: "/personalSetting/sshAndGpg/createSshKey",
                component: () =>
                    import ("../pages/personalSetting/SshKeyCreatePage")
            }, {
                path: "/personalSetting/sshAndGpg/createGpgKey",
                component: () =>
                    import ("../pages/personalSetting/GpgKeyCreatePage")
            }]
        },
        {
            path: "/db",
            redirect: "/db/mysqlDb/list",
            component: () =>
                import ("../layouts/MysqlAuditLayout"),
            children: [{
                path: "/db/mysqlDb/list",
                component: () =>
                    import ("../pages/db/mysql/MysqlDbListPage")
            }, {
                path: "/db/mysqlDb/create",
                component: () =>
                    import ("../pages/db/mysql/MysqlDbHandlePage")
            }, {
                path: "/db/mysqlDb/:dbId(\\d+)/update",
                component: () =>
                    import ("../pages/db/mysql/MysqlDbHandlePage")
            }, {
                path: "/db/mysqlReadPermApply/list",
                component: () =>
                    import ("../pages/db/mysql/MysqlReadPermApplyListPage")
            }, {
                path: "/db/mysqlReadPermApply/apply",
                component: () =>
                    import ("../pages/db/mysql/MysqlReadPermApplyCreatePage")
            }, {
                path: "/db/mysqlReadPermAudit/list",
                component: () =>
                    import ("../pages/db/mysql/MysqlReadPermAuditListPage")
            }, {
                path: "/db/mysqlReadPermDetail/list",
                component: () =>
                    import ("../pages/db/mysql/MysqlReadPermDetailListPage")
            }, {
                path: "/db/mysqlSearch",
                component: () =>
                    import ("../pages/db/mysql/MysqlSearchPage")
            }, {
                path: "/db/mysqlDataUpdateApply/list",
                component: () =>
                    import ("../pages/db/mysql/MysqlDataUpdateApplyListPage")
            }, {
                path: "/db/mysqlDataUpdateApply/apply",
                component: () =>
                    import ("../pages/db/mysql/MysqlDataUpdateApplyCreatePage")
            }, {
                path: "/db/mysqlDataUpdateAudit/list",
                component: () =>
                    import ("../pages/db/mysql/MysqlDataUpdateAuditListPage")
            }, {
                path: "/db/mysqlReadPermManage/list",
                component: () =>
                    import ("../pages/db/mysql/MysqlReadPermManageListPage")
            }]
        },
        {
            path: "/team",
            component: () =>
                import ("../layouts/TeamLayout"),
            children: [{
                path: "/team/:teamId(\\d+)/gitRepo/create",
                component: () =>
                    import ("../pages/team/gitRepo/RepoCreatePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/list",
                component: () =>
                    import ("../pages/team/gitRepo/RepoListPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/recycle",
                component: () =>
                    import ("../pages/team/gitRepo/RepoRecyclePage")
            }, {
                path: "/team/:teamId(\\d+)/role/list",
                component: () =>
                    import ("../pages/team/team/RoleListPage")
            }, {
                path: "/team/:teamId(\\d+)/role/create",
                component: () =>
                    import ("../pages/team/team/RoleHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/role/:roleId(\\d+)/update",
                component: () =>
                    import ("../pages/team/team/RoleHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/role/user/list",
                component: () =>
                    import ("../pages/team/team/UserListPage")
            }, {
                path: "/team/:teamId(\\d+)/timer/list/:env?",
                component: () =>
                    import ("../pages/team/team/TimerListPage")
            }, {
                path: "/team/:teamId(\\d+)/timer/create",
                component: () =>
                    import ("../pages/team/team/TimerHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/timer/:timerId(\\d+)/update",
                component: () =>
                    import ("../pages/team/team/TimerHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/timer/:timerId(\\d+)/log",
                component: () =>
                    import ("../pages/team/team/TimerLogListPage")
            }, {
                path: "/team/:teamId(\\d+)/app/list",
                component: () =>
                    import ("../pages/team/app/AppListPage.vue")
            }, {
                path: "/team/:teamId(\\d+)/app/create",
                component: () =>
                    import ("../pages/team/app/AppCreatePage")
            }, {
                path: "/team/:teamId(\\d+)/setting",
                component: () =>
                    import ("../pages/team/team/TeamSettingPage")
            }, {
                path: "/team/:teamId(\\d+)/notifyTpl/list",
                component: () =>
                    import ("../pages/team/team/NotifyTplListPage")
            }, {
                path: "/team/:teamId(\\d+)/notifyTpl/create",
                component: () =>
                    import ("../pages/team/team/NotifyTplHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/notifyTpl/:tplId(\\d+)/update",
                component: () =>
                    import ("../pages/team/team/NotifyTplHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/teamHook/create",
                component: () =>
                    import ("../pages/team/team/TeamHookHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/teamHook/list",
                component: () =>
                    import ("../pages/team/team/TeamHookListPage")
            }, {
                path: "/team/:teamId(\\d+)/teamHook/:hookId(\\d+)/update",
                component: () =>
                    import ("../pages/team/team/TeamHookHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/weworkAccessToken/list",
                component: () =>
                    import ("../pages/team/team/WeworkAccessTokenListPage")
            }, {
                path: "/team/:teamId(\\d+)/weworkAccessToken/create",
                component: () =>
                    import ("../pages/team/team/WeworkAccessTokenHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/weworkAccessToken/:tokenId(\\d+)/update",
                component: () =>
                    import ("../pages/team/team/WeworkAccessTokenHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/feishuAccessToken/list",
                component: () =>
                    import ("../pages/team/team/FeishuAccessTokenListPage")
            }, {
                path: "/team/:teamId(\\d+)/feishuAccessToken/create",
                component: () =>
                    import ("../pages/team/team/FeishuAccessTokenHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/feishuAccessToken/:tokenId(\\d+)/update",
                component: () =>
                    import ("../pages/team/team/FeishuAccessTokenHandlePage")
            }]
        },
        {
            path: "/team/:teamId(\\d+)/gitRepo",
            component: () =>
                import ("../layouts/GitRepoLayout"),
            children: [{
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/index",
                component: () =>
                    import ("../pages/team/gitRepo/RepoIndexPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/tree/:refType/:ref/:files*",
                component: () =>
                    import ("../pages/team/gitRepo/RepoTreePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/pullRequest/list",
                component: () =>
                    import ("../pages/team/gitRepo/PullRequestListPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/pullRequest/create",
                component: () =>
                    import ("../pages/team/gitRepo/PullRequestCreatePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/pullRequest/:index(\\d+)/detail",
                component: () =>
                    import ("../pages/team/gitRepo/PullRequestDetailPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/branch/list",
                component: () =>
                    import ("../pages/team/gitRepo/BranchListPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/commit/list/:ref?",
                component: () =>
                    import ("../pages/team/gitRepo/HistoryCommitPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/commit/diff/:commitId",
                component: () =>
                    import ("../pages/team/gitRepo/DiffCommitsPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/tag/list",
                component: () =>
                    import ("../pages/team/gitRepo/TagListPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/protectedBranch/list",
                component: () =>
                    import ("../pages/team/gitRepo/ProtectedBranchListPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/protectedBranch/create",
                component: () =>
                    import ("../pages/team/gitRepo/ProtectedBranchHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/protectedBranch/:protectedBranchId(\\d+)/update",
                component: () =>
                    import ("../pages/team/gitRepo/ProtectedBranchHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/webhook/list",
                component: () =>
                    import ("../pages/team/gitRepo/WebhookListPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/webhook/create",
                component: () =>
                    import ("../pages/team/gitRepo/WebhookHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/webhook/:webhookId(\\d+)/update",
                component: () =>
                    import ("../pages/team/gitRepo/WebhookHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/workflow/list",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowListPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/workflow/create",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/workflow/:workflowId(\\d+)/update",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/workflow/:workflowId(\\d+)/tasks",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowTasksPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/workflow/:workflowId(\\d+)/:taskId(\\d+)/steps",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowStepsPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/workflow/vars",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowVarsPage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/workflow/vars/create",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowVarsHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/workflow/vars/:varsId(\\d+)/update",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowVarsHandlePage")
            }, {
                path: "/team/:teamId(\\d+)/gitRepo/:repoId(\\d+)/setting",
                component: () =>
                    import ("../pages/team/gitRepo/RepoSettingPage")
            }]
        },
        {
            path: "/team/:teamId(\\d+)/app",
            component: () =>
                import ("../layouts/AppLayout"),
            children: [{
                path: "/team/:teamId(\\d+)/app/:appId/propertyFile/list/:env?",
                component: () =>
                    import ("../pages/team/app/PropertyFileListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/propertyFile/create",
                component: () =>
                    import ("../pages/team/app/PropertyFileHandlePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/propertyFile/:fileId(\\d+)/new",
                component: () =>
                    import ("../pages/team/app/PropertyFileHandlePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/propertyFile/:fileId(\\d+)/deploy/:version",
                component: () =>
                    import ("../pages/team/app/PropertyHistoryDeployPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/propertyFile/:fileId(\\d+)/history/list",
                component: () =>
                    import ("../pages/team/app/PropertyHistoryListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/pipeline/list/:env?",
                component: () =>
                    import ("../pages/team/app/PipelineListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/pipeline/create",
                component: () =>
                    import ("../pages/team/app/PipelineHandlePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/pipeline/:pipelineId(\\d+)/update",
                component: () =>
                    import ("../pages/team/app/PipelineHandlePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/deployPlan/list/:env?",
                component: () =>
                    import ("../pages/team/app/DeployPlanListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/deployPlan/create",
                component: () =>
                    import ("../pages/team/app/DeployPlanCreatePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/deployPlan/:planId(\\d+)/view",
                component: () =>
                    import ("../pages/team/app/DeployPlanViewPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/pipeline/vars/:env?",
                component: () =>
                    import ("../pages/team/app/PipelineVarsListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/pipeline/vars/create",
                component: () =>
                    import ("../pages/team/app/PipelineVarsHandlePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/pipeline/vars/:varsId(\\d+)/update",
                component: () =>
                    import ("../pages/team/app/PipelineVarsHandlePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/serviceStatus/list/:env?",
                component: () =>
                    import ("../pages/team/app/ServiceStatusListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/discoveryService/list/:env?",
                component: () =>
                    import ("../pages/team/app/DiscoveryServiceListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/setting",
                component: () =>
                    import ("../pages/team/app/AppSettingPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/artifact/list/:env?",
                component: () =>
                    import ("../pages/team/app/ArtifactListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/promScrape/list/:env?",
                component: () =>
                    import ("../pages/team/app/PromScrapeListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/promScrape/create",
                component: () =>
                    import ("../pages/team/app/PromScrapeHandlePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/promScrape/:scrapeId(\\d+)/update",
                component: () =>
                    import ("../pages/team/app/PromScrapeHandlePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/alertConfig/list/:env?",
                component: () =>
                    import ("../pages/team/app/AlertConfigListPage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/alertConfig/create",
                component: () =>
                    import ("../pages/team/app/AlertConfigHandlePage")

            }, {
                path: "/team/:teamId(\\d+)/app/:appId/alertConfig/:configId(\\d+)/update",
                component: () =>
                    import ("../pages/team/app/AlertConfigHandlePage")

            }]
        },
        {
            path: "/login",
            redirect: "/login/login",
            component: () =>
                import ("../layouts/LoginLayout"),
            children: [{
                path: "/login/login",
                name: "login",
                component: () =>
                    import ("../pages/login/LoginPage")
            }, {
                path: "/login/register",
                name: "register",
                component: () =>
                    import ("../pages/login/RegisterPage")
            }]
        }, {
            path: '/:pathMatch(.*)',
            component: () =>
                import ("../layouts/NotFoundLayout")
        }
    ]
});
router.beforeEach((to, from, next) => {
    // 登录或注册页面无需关注登录用户信息
    if (to.name === "login" || to.name === "register") {
        next();
        return;
    }
    // 获取登录人信息
    const user = useUserStore();
    if (!user.account) {
        getUserInfoRequest().then(res => {
            user.account = res.session.userInfo.account;
            user.avatarUrl = res.session.userInfo.avatarUrl;
            user.email = res.session.userInfo.email;
            user.isAdmin = res.session.userInfo.isAdmin;
            user.name = res.session.userInfo.name;
            user.isDba = res.session.userInfo.isDba;
            user.sessionId = res.session.sessionId;
            user.sessionExpireAt = res.session.expireAt;
        });
    }
    next();
})
export default router;