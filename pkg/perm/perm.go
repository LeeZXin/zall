package perm

import "encoding/json"

var (
	DefaultTeamPerm = TeamPerm{
		CanCreateRepo:         true,
		CanDeleteRepo:         true,
		CanHandleTimer:        true,
		CanAccessAction:       true,
		CanUpdateAction:       true,
		CanTriggerAction:      true,
		CanHandleDeployConfig: true,
		CanHandleDeployPlan:   true,
	}
	DefaultRepoPerm = RepoPerm{
		CanAccessRepo:        true,
		CanUpdateRepo:        true,
		CanCloseRepo:         true,
		CanHandlePullRequest: true,
		CanHandleWebhook:     true,
		CanAccessToken:       true,
		CanUpdateToken:       true,
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
	RepoPermList []RepoPermWithId `json:"repoPermList"`
	// 可开发应用
	DevelopAppList []string `json:"developAppList"`
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
	return len(d.RepoPermList) < 1000 && len(d.DevelopAppList) < 1000
}

func (d *Detail) FromDB(content []byte) error {
	if d == nil {
		*d = Detail{}
	}
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
	CanUpdateRepo bool `json:"canUpdateRepo"`
	// 是否可归档
	CanCloseRepo bool `json:"canCloseRepo"`
	// 是否可处理pr
	CanHandlePullRequest bool `json:"canHandlePullRequest"`
	// 是否可配置webhook
	CanHandleWebhook bool `json:"canHandleWebhook"`
	// 是否可查看Token
	CanAccessToken bool `json:"canAccessToken"`
	// 是否可编辑Token
	CanUpdateToken bool `json:"canUpdateToken"`
}

type TeamPerm struct {
	// 是否可创建仓库
	CanCreateRepo bool `json:"canCreateRepo"`
	// 是否可删除仓库
	CanDeleteRepo bool `json:"canDeleteRepo"`
	// 是否可处理定时任务
	CanHandleTimer bool `json:"canHandleTimer"`
	// 是否可查看action
	CanAccessAction bool `json:"canAccessAction"`
	// 是否可编辑action
	CanUpdateAction bool `json:"canUpdateAction"`
	// 是否可手动触发action
	CanTriggerAction bool `json:"canTriggerAction"`
	// 是否可编辑部署配置
	CanHandleDeployConfig bool `json:"canHandleDeployConfig"`
	// 是否可直接创建发布计划
	CanHandleDeployPlan bool `json:"canHandleDeployPlan"`
}
