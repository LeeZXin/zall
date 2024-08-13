package perm

import (
	"encoding/json"
	"github.com/LeeZXin/zsf-utils/listutil"
)

var (
	DefaultAppPerm = AppPerm{
		CanDevelop:               true,
		CanManagePipeline:        true,
		CanManageServiceSource:   true,
		CanManagePropertySource:  true,
		CanDeployProperty:        true,
		CanManageDiscoverySource: true,
		CanManagePromAgent:       true,
		CanCreateDeployPlan:      true,
	}
	DefaultTeamPerm = TeamPerm{
		CanManageTimer: true,
	}
	DefaultRepoPerm = RepoPerm{
		CanAccessRepo:              true,
		CanPushRepo:                true,
		CanSubmitPullRequest:       true,
		CanAddCommentInPullRequest: true,
		CanManageWebhook:           true,
		CanManageWorkflow:          true,
		CanTriggerWorkflow:         true,
		CanManageProtectedBranch:   true,
	}
	DefaultPermDetail = Detail{
		TeamPerm:        DefaultTeamPerm,
		DefaultRepoPerm: DefaultRepoPerm,
		DefaultAppPerm:  DefaultAppPerm,
	}
)

type Detail struct {
	// 项目权限
	TeamPerm TeamPerm `json:"teamPerm"`
	// 仓库权限
	DefaultRepoPerm RepoPerm `json:"defaultRepoPerm"`
	// 特殊仓库权限
	RepoPermList []RepoPermWithId `json:"repoPermList,omitempty"`
	// 应用服务权限
	DefaultAppPerm AppPerm `json:"defaultAppPerm"`
	// 特殊应用服务权限
	AppPermList []AppPermWithId `json:"appPermList,omitempty"`
}

func (d *Detail) GetRepoPerm(repoId int64) RepoPerm {
	if len(d.RepoPermList) == 0 {
		return d.DefaultRepoPerm
	}
	for _, p := range d.RepoPermList {
		if p.RepoId == repoId {
			return p.RepoPerm
		}
	}
	return RepoPerm{}
}

func (d *Detail) GetAppPerm(appId string) AppPerm {
	if len(d.AppPermList) == 0 {
		return d.DefaultAppPerm
	}
	for _, p := range d.AppPermList {
		if p.AppId == appId {
			return p.AppPerm
		}
	}
	return AppPerm{}
}

func (d *Detail) IsValid() bool {
	if len(d.RepoPermList) > 1000 || len(d.AppPermList) > 1000 {
		return false
	}
	for _, repoPerm := range d.RepoPermList {
		if repoPerm.RepoId <= 0 {
			return false
		}
	}
	repoIdList, _ := listutil.Map(d.RepoPermList, func(t RepoPermWithId) (int64, error) {
		return t.RepoId, nil
	})
	if len(listutil.Distinct(repoIdList...)) != len(repoIdList) {
		return false
	}
	for _, appPerm := range d.AppPermList {
		if appPerm.AppId == "" {
			return false
		}
	}
	appIdList, _ := listutil.Map(d.AppPermList, func(t AppPermWithId) (string, error) {
		return t.AppId, nil
	})
	if len(listutil.Distinct(appIdList...)) != len(appIdList) {
		return false
	}
	return true
}

func (d *Detail) FromDB(content []byte) error {
	return json.Unmarshal(content, d)
}

func (d *Detail) ToDB() ([]byte, error) {
	return json.Marshal(d)
}

type RepoPermWithId struct {
	RepoPerm
	RepoId int64 `json:"repoId"`
}

type AppPermWithId struct {
	AppPerm
	AppId string `json:"appId"`
}

type RepoPerm struct {
	// 可访问
	CanAccessRepo bool `json:"canAccessRepo"`
	// 可推送代码
	CanPushRepo bool `json:"canPushRepo"`
	// 是否可提交合并请求
	CanSubmitPullRequest bool `json:"canSubmitPullRequest"`
	// 合并请求是否可以评论
	CanAddCommentInPullRequest bool `json:"canAddCommentInPullRequest"`
	// 是否可配置webhook
	CanManageWebhook bool `json:"canManageWebhook"`
	// 查看工作流
	CanManageWorkflow bool `json:"canManageWorkflow"`
	// 触发工作流
	CanTriggerWorkflow bool `json:"canTriggerWorkflow"`
	// 管理保护分支
	CanManageProtectedBranch bool `json:"canManageProtectedBranch"`
}

type TeamPerm struct {
	// 是否可管理定时任务
	CanManageTimer bool `json:"canManageTimer"`
}

type AppPerm struct {
	// 是否可开发
	CanDevelop bool `json:"canDevelop"`
	// 是否可管理流水线
	CanManagePipeline bool `json:"canManagePipeline"`
	// 是否可管理服务来源
	CanManageServiceSource bool `json:"canManageServiceSource"`
	// 是否可管理配置来源
	CanManagePropertySource bool `json:"canManagePropertySource"`
	// 是否可发布配置
	CanDeployProperty bool `json:"canDeployProperty"`
	// 是否可管理注册中心来源
	CanManageDiscoverySource bool `json:"canManageDiscoverySource"`
	// 是否可管理监控告警
	CanManagePromAgent bool `json:"canManagePromAgent"`
	// 是否可直接创建发布计划
	CanCreateDeployPlan bool `json:"canCreateDeployPlan"`
}
