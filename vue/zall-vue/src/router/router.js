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
            redirect: "/teamList",
            component: () =>
                import ("../layouts/IndexLayout"),
            children: [{
                    path: "/teamList",
                    component: () =>
                        import ("../pages/team/team/TeamListPage")
                },
                {
                    path: "/createTeam",
                    component: () =>
                        import ("../pages/team/team/CreateTeamPage")
                }
            ]
        },
        {
            path: "/team",
            component: () =>
                import ("../layouts/TeamLayout"),
            children: [{
                    path: "/team/:teamId(\\d+)/gitRepo/create",
                    component: () =>
                        import ("../pages/team/gitRepo/RepoCreatePage")
                },
                {
                    path: "/team/:teamId(\\d+)/gitRepo/list",
                    component: () =>
                        import ("../pages/team/gitRepo/RepoListPage")
                },
                {
                    path: "/team/app/list",
                    component: () =>
                        import ("../pages/team/appService/AppListPage")
                },
                {
                    path: "/team/app/create",
                    component: () =>
                        import ("../pages/team/appService/CreateAppPage")
                },
                {
                    path: "/team/action/list",
                    component: () =>
                        import ("../pages/team/action/ActionListPage")
                },
                {
                    path: "/team/action/create",
                    component: () =>
                        import ("../pages/team/action/CreateActionPage")
                },
                {
                    path: "/team/action/update/:actionId(\\d+)",
                    component: () =>
                        import ("../pages/team/action/UpdateActionPage")
                },
                {
                    path: "/team/action/detail/:actionId(\\d+)",
                    component: () =>
                        import ("../pages/team/action/ActionDetailPage")
                },
                {
                    path: "/team/action/task/:actionId(\\d+)",
                    component: () =>
                        import ("../pages/team/action/ActionTaskPage")
                },
                {
                    path: "/team/action/task/:actionId(\\d+)/step/:taskId(\\d+)",
                    component: () =>
                        import ("../pages/team/action/ActionStepPage")
                }
            ]
        },
        {
            path: "/gitRepo",
            component: () =>
                import ("../layouts/GitRepoLayout"),
            children: [{
                path: "/gitRepo/:repoId(\\d+)/index",
                component: () =>
                    import ("../pages/team/gitRepo/RepoIndexPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/tree/:refType/:ref/:files*",
                component: () =>
                    import ("../pages/team/gitRepo/RepoTreePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/pullRequest/list",
                component: () =>
                    import ("../pages/team/gitRepo/PullRequestsPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/pullRequest/create",
                component: () =>
                    import ("../pages/team/gitRepo/PullRequestCreatePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/pullRequest/:prId(\\d+)/detail",
                component: () =>
                    import ("../pages/team/gitRepo/PullRequestDetailPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/branch/list",
                component: () =>
                    import ("../pages/team/gitRepo/BranchesPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/commit/list/:ref?",
                component: () =>
                    import ("../pages/team/gitRepo/HistoryCommitPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/commit/diff/:commitId",
                component: () =>
                    import ("../pages/team/gitRepo/DiffCommitsPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/tag/list",
                component: () =>
                    import ("../pages/team/gitRepo/TagsPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/protectedBranch/list",
                component: () =>
                    import ("../pages/team/gitRepo/ProtectedBranchesPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/protectedBranch/create",
                component: () =>
                    import ("../pages/team/gitRepo/ProtectedBranchHandlePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/protectedBranch/:protectedBranchId(\\d+)/update",
                component: () =>
                    import ("../pages/team/gitRepo/ProtectedBranchHandlePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/webhook/list",
                component: () =>
                    import ("../pages/team/gitRepo/WebhooksPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/webhook/create",
                component: () =>
                    import ("../pages/team/gitRepo/WebhookHandlePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/webhook/:webhookId(\\d+)/update",
                component: () =>
                    import ("../pages/team/gitRepo/WebhookHandlePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/workflow/list",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowsPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/workflow/create",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowHandlePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/workflow/:workflowId(\\d+)/update",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowHandlePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/workflow/:workflowId(\\d+)/tasks",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowTasksPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/workflow/:workflowId(\\d+)/:taskId(\\d+)/steps",
                component: () =>
                    import ("../pages/team/gitRepo/WorkflowStepsPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/workflow/secrets",
                component: () =>
                    import ("../pages/team/gitRepo/SecretsPage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/workflow/secret/create",
                component: () =>
                    import ("../pages/team/gitRepo/SecretHandlePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/workflow/secret/:secretId(\\d+)/update",
                component: () =>
                    import ("../pages/team/gitRepo/SecretHandlePage")
            }, {
                path: "/gitRepo/:repoId(\\d+)/config",
                component: () =>
                    import ("../pages/team/gitRepo/RepoConfigPage")
            }]
        },
        {
            path: "/appService",
            component: () =>
                import ("../layouts/AppServiceLayout"),
            children: [{
                path: "/appService/property/list",
                component: () =>
                    import ("../pages/team/appService/PropertyListPage")

            }, {
                path: "/appService/property/create",
                component: () =>
                    import ("../pages/team/appService/CreatePropertyPage")

            }, {
                path: "/appService/property/update/:id(\\d+)",
                component: () =>
                    import ("../pages/team/appService/UpdatePropertyPage")

            }, {
                path: "/appService/property/deploy/:id(\\d+)",
                component: () =>
                    import ("../pages/team/appService/DeployPropertyPage")

            }, {
                path: "/appService/property/history/:id(\\d+)",
                component: () =>
                    import ("../pages/team/appService/PropertyHistoryPage")

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
            component: import ("../layouts/NotFoundLayout")
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
            user.roleType = res.session.userInfo.roleType;
            user.sessionId = res.session.sessionId;
            user.sessionExpireAt = res.session.expireAt;
        });
    }
    next();
})
export default router;