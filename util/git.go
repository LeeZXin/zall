package util

import (
	"github.com/IGLOU-EU/go-wildcard/v2"
	"strings"
)

func LongCommitId2ShortId(commitId string) string {
	if len(commitId) < 7 {
		return commitId
	}
	return commitId[:7]
}

func BaseRefName(ref string) string {
	split := strings.Split(ref, "/")
	if len(split) > 0 {
		return split[len(split)-1]
	}
	return ref
}

func WildcardMatchBranches(branches []string, ref string) bool {
	for _, branch := range branches {
		if wildcard.Match(branch, ref) {
			return true
		}
	}
	return false
}
