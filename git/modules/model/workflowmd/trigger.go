package workflowmd

import (
	"encoding/json"
	"github.com/IGLOU-EU/go-wildcard/v2"
)

type SourceType int

const (
	BranchTriggerSource SourceType = iota + 1
	PullRequestTriggerSource
)

func (s SourceType) IsValid() bool {
	switch s {
	case BranchTriggerSource, PullRequestTriggerSource:
		return true
	default:
		return false
	}
}

type PullRequestSource struct {
	Branches []string `json:"branches"`
}

func (s *PullRequestSource) IsValid() bool {
	if len(s.Branches) == 0 {
		return false
	}
	for _, branch := range s.Branches {
		if len(branch) > 32 {
			return false
		}
	}
	return true
}

type BranchSource []string

func (s BranchSource) IsValid() bool {
	if len(s) == 0 {
		return false
	}
	for _, branch := range s {
		if len(branch) > 32 {
			return false
		}
	}
	return true
}

type Source struct {
	SourceType        SourceType         `json:"sourceType"`
	PullRequestSource *PullRequestSource `json:"pullRequestSource,omitempty"`
	BranchSource      BranchSource       `json:"branchSource,omitempty"`
}

func (c *Source) MatchBranchBySource(source SourceType, branch string) bool {
	switch source {
	case PullRequestTriggerSource:
		if c.PullRequestSource != nil {
			for _, wildBranch := range c.PullRequestSource.Branches {
				if wildcard.Match(wildBranch, branch) {
					return true
				}
			}
		}
	case BranchTriggerSource:
		if c.BranchSource != nil {
			for _, wildBranch := range c.BranchSource {
				if wildcard.Match(wildBranch, branch) {
					return true
				}
			}
		}
	}
	return false
}

func (c *Source) IsValid() bool {
	switch c.SourceType {
	case PullRequestTriggerSource:
		if c.PullRequestSource == nil {
			return false
		}
		return c.PullRequestSource.IsValid()
	case BranchTriggerSource:
		if c.BranchSource == nil {
			return false
		}
		return c.BranchSource.IsValid()
	default:
		return false
	}
}

func (c *Source) FromDB(content []byte) error {
	if c == nil {
		*c = Source{}
	}
	return json.Unmarshal(content, c)
}

func (c *Source) ToDB() ([]byte, error) {
	return json.Marshal(c)
}
