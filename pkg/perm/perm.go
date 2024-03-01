package perm

var (
	DefaultTeamPerm = TeamPerm{
		CanInitRepo:    true,
		CanDeleteRepo:  true,
		CanHandleTimer: true,
	}
	DefaultRepoPerm = RepoPerm{
		CanAccessRepo:            true,
		CanUpdateRepo:            true,
		CanCloseRepo:             true,
		CanHandleProtectedBranch: true,
		CanHandlePullRequest:     true,
		CanHandleWebhook:         true,
		CanAccessWiki:            true,
		CanUpdateWiki:            true,
		CanAccessAction:          true,
		CanUpdateAction:          true,
		CanAccessToken:           true,
		CanUpdateToken:           true,
	}
	DefaultPermDetail = Detail{
		TeamPerm:             DefaultTeamPerm,
		ApplyDefaultRepoPerm: true,
		DefaultRepoPerm:      DefaultRepoPerm,
	}
)

type Detail struct {
	// 项目权限
	TeamPerm TeamPerm `json:"teamPerm"`
	// 使用仓库全局默认权限
	ApplyDefaultRepoPerm bool `json:"applyDefaultRepoPerm"`
	// 默认仓库权限
	DefaultRepoPerm RepoPerm `json:"defaultRepoPerm"`
	// 可访问仓库权限列表
	RepoPermList []RepoPermWithId `json:"repoPermList,omitempty"`
}

func (d *Detail) GetRepoPerm(repoId int64) RepoPerm {
	if d.ApplyDefaultRepoPerm {
		return d.DefaultRepoPerm
	}
	for _, perm := range d.RepoPermList {
		if perm.RepoId == repoId {
			return perm.RepoPerm
		}
	}
	return RepoPerm{}
}

type RepoPermWithId struct {
	RepoId int64 `json:"repoId"`
	RepoPerm
}

type RepoPerm struct {
	// 可访问
	CanAccessRepo bool `json:"canAccessRepo"`
	// 可推送代码
	CanUpdateRepo bool `json:"canUpdateRepo"`
	// 是否可归档
	CanCloseRepo bool `json:"canCloseRepo"`
	// 是否可处理保护分支
	CanHandleProtectedBranch bool `json:"canHandleProtectedBranch"`
	// 是否可处理pr
	CanHandlePullRequest bool `json:"canHandlePullRequest"`
	// 是否可配置webhook
	CanHandleWebhook bool `json:"canHandleWebhook"`
	// 是否可访问wiki
	CanAccessWiki bool `json:"canAccessWiki"`
	// 是否可编辑wiki
	CanUpdateWiki bool `json:"canUpdateWiki"`
	// 是否可查看action
	CanAccessAction bool `json:"canAccessAction"`
	// 是否可编辑action
	CanUpdateAction bool `json:"canUpdateAction"`
	// 是否可手动触发action
	CanTriggerAction bool `json:"canTriggerAction"`
	// 是否可查看accessToken
	CanAccessToken bool `json:"canAccessToken"`
	// 是否可编辑accessToken
	CanUpdateToken bool `json:"canUpdateToken"`
}

type TeamPerm struct {
	// 是否可创建仓库
	CanInitRepo bool `json:"canInitRepo"`
	// 是否可删除仓库
	CanDeleteRepo bool `json:"canDeleteRepo"`
	// 是否可处理定时任务
	CanHandleTimer bool `json:"canHandleTimer"`
}
