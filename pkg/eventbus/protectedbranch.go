package eventbus

import (
	"github.com/LeeZXin/zall/pkg/branch"
	"time"
)

const ProtectedBranchEventTopic = "protected-branch"

type ProtectedBranchObj struct {
	Pattern string `json:"pattern"`
	branch.ProtectedBranchCfg
}

type ProtectedBranchEvent struct {
	RepoId    int64               `json:"repoId"`
	Name      string              `json:"name"`
	Path      string              `json:"path"`
	Operator  string              `json:"operator"`
	Action    string              `json:"action"`
	Before    *ProtectedBranchObj `json:"before,omitempty"`
	After     *ProtectedBranchObj `json:"after,omitempty"`
	EventTime time.Time           `json:"eventTime"`
}
