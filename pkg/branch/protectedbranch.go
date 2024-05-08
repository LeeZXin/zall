package branch

import "encoding/json"

type PushOption int

const (
	AllowPush PushOption = iota
	NotAllowPush
	WhiteListPush
)

func (o PushOption) IsValid() bool {
	switch o {
	case AllowPush, NotAllowPush, WhiteListPush:
		return true
	default:
		return false
	}
}

type ProtectedBranchCfg struct {
	// 受保护文件模式
	ProtectedFilePatterns []string `json:"protectedFilePatterns"`
	// 推送模式
	PushOption PushOption `json:"pushOption"`
	// 当推送人模式为白名单
	PushWhiteList []string `json:"pushWhiteList"`
	// 合并请求时代码评审数量大于该数量才能合并
	ReviewCountWhenCreatePr int `json:"reviewCountWhenCreatePr"`
	// 代码评审员名单
	ReviewerList []string `json:"reviewerList"`
	// 撤销过时的审批
	CancelOldReviewApprovalWhenNewCommit bool `json:"cancelOldReviewApprovalWhenNewCommit"`
}

func (c *ProtectedBranchCfg) IsValid() bool {
	// 如果限制了白名单 白名单数量肯定大于等于审核数量
	if len(c.ReviewerList) > 0 && c.ReviewCountWhenCreatePr > len(c.ReviewerList) {
		return false
	}
	return c.PushOption.IsValid() &&
		len(c.PushWhiteList) <= 50 && len(c.ReviewerList) <= 50 &&
		len(c.ProtectedFilePatterns) <= 1000
}

func (c *ProtectedBranchCfg) ToString() string {
	m, _ := json.Marshal(c)
	return string(m)
}

func (c *ProtectedBranchCfg) FromDB(content []byte) error {
	if c == nil {
		*c = ProtectedBranchCfg{}
	}
	return json.Unmarshal(content, c)
}

func (c *ProtectedBranchCfg) ToDB() ([]byte, error) {
	return json.Marshal(c)
}
