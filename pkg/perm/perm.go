package perm

import (
	"encoding/json"
	"github.com/LeeZXin/zsf-utils/listutil"
)

var (
	DefaultTeamPerm = TeamPerm{
		CanCreateRepo:          true,
		CanManageDeployConfig:  true,
		CanCreateDeployPlan:    true,
		CanManageTimer:         true,
		CanManagePipeline:      true,
		CanManageServiceSource: true,
	}
	DefaultRepoPerm = RepoPerm{
		CanAccessRepo:              true,
		CanPushRepo:                true,
		CanSubmitPullRequest:       true,
		CanReviewPullRequest:       true,
		CanAddCommentInPullRequest: true,
		CanManageWebhook:           true,
		CanManageWorkflow:          true,
		CanTriggerWorkflow:         true,
	}
	DefaultPermDetail = Detail{
		TeamPerm:        DefaultTeamPerm,
		DefaultRepoPerm: DefaultRepoPerm,
	}
)

type Detail struct {
	// 项目权限
	TeamPerm TeamPerm `json:"teamPerm"`
	// 仓库权限
	DefaultRepoPerm RepoPerm `json:"defaultRepoPerm"`
	// 特殊仓库权限
	RepoPermList []RepoPermWithId `json:"repoPermList,omitempty"`
	// 可开发应用
	DevelopAppList listutil.ComparableList[string] `json:"developAppList,omitempty"`
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

func (d *Detail) IsValid() bool {
	if len(d.RepoPermList) > 1000 || len(d.DevelopAppList) > 1000 {
		return false
	}
	for _, repoPerm := range d.RepoPermList {
		if repoPerm.RepoId <= 0 {
			return false
		}
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

type RepoPerm struct {
	// 可访问
	CanAccessRepo bool `json:"canAccessRepo"`
	// 可推送代码
	CanPushRepo bool `json:"canPushRepo"`
	// 是否可提交合并请求
	CanSubmitPullRequest bool `json:"canSubmitPullRequest"`
	// 是否可评审合并请求
	CanReviewPullRequest bool `json:"canReviewPullRequest"`
	// 合并请求是否可以评论
	CanAddCommentInPullRequest bool `json:"canAddCommentInPullRequest"`
	// 是否可配置webhook
	CanManageWebhook bool `json:"canManageWebhook"`
	// 查看工作流
	CanManageWorkflow bool `json:"canManageWorkflow"`
	// 触发工作流
	CanTriggerWorkflow bool `json:"canTriggerWorkflow"`
}

type TeamPerm struct {
	// 是否可创建仓库
	CanCreateRepo bool `json:"canCreateRepo"`
	// 是否可管理部署配置
	CanManageDeployConfig bool `json:"canManageDeployConfig"`
	// 是否可直接创建发布计划
	CanCreateDeployPlan bool `json:"canCreateDeployPlan"`
	// 是否可管理定时任务
	CanManageTimer bool `json:"canManageTimer"`
	// 是否可管理流水线
	CanManagePipeline bool `json:"canManagePipeline"`
	// 是否可管理服务来源
	CanManageServiceSource bool `json:"canManageServiceSource"`
}
