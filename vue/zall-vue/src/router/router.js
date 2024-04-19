import { createRouter, createWebHashHistory } from 'vue-router'
const router = createRouter({
    history: createWebHashHistory(),
    routes: [{
            path: "",
            redirect: "/login/login"
        },
        {
            path: "/index",
            component: () =>
                import ("../layouts/IndexLayout")
        },
        {
            path: "/createTeam",
            component: () =>
                import ("../layouts/CreateTeamLayout"),
        },
        {
            path: "/team",
            redirect: "/team/gitRepo/list",
            component: () =>
                import ("../layouts/TeamLayout"),
            children: [{
                    path: "/team/gitRepo/create",
                    component: () =>
                        import ("../pages/team/gitRepo/CreateRepoPage")
                },
                {
                    path: "/team/gitRepo/list",
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
                path: "/gitRepo/index",
                component: () =>
                    import ("../pages/team/gitRepo/RepoIndexPage")
            }, {
                path: "/gitRepo/tree/:file+",
                component: () =>
                    import ("../pages/team/gitRepo/RepoTreePage")
            }, {
                path: "/gitRepo/pullRequests/list",
                component: () =>
                    import ("../pages/team/gitRepo/PullRequestsListPage")
            }, {
                path: "/gitRepo/pullRequests/create",
                component: () =>
                    import ("../pages/team/gitRepo/CreatePullRequestPage")
            }, {
                path: "/gitRepo/Branches/list",
                component: () =>
                    import ("../pages/team/gitRepo/BranchesPage")
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

            }]
        },
        {
            path: "/login",
            redirect: "/login/login",
            component: () =>
                import ("../layouts/LoginLayout"),
            children: [{
                path: "/login/login",
                component: () =>
                    import ("../pages/login/LoginPage")
            }, {
                path: "/login/register",
                component: () =>
                    import ("../pages/login/RegisterPage")
            }]

        }, {
            path: '/:pathMatch(.*)',
            component: import ("../layouts/NotFoundLayout")
        }
    ]
})
export default router;