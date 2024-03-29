package git

import (
	"context"
	"strings"
)

func GetAllBranchList(ctx context.Context, repoPath string) ([]string, error) {
	cmd := NewCommand("for-each-ref", "--format=%(objectname) %(refname)", BranchPrefix, "--sort=-committerdate")
	pipeResult := cmd.RunWithReadPipe(ctx, WithDir(repoPath))
	ret := make([]string, 0)
	err := pipeResult.RangeStringLines(func(_ int, line string) (bool, error) {
		split := strings.Split(strings.TrimSpace(line), " ")
		if len(split) == 2 {
			ret = append(ret, strings.TrimPrefix(split[1], BranchPrefix))
		}
		return true, nil
	})
	return ret, err
}

func CheckRefIsBranch(ctx context.Context, repoPath string, branch string) bool {
	if !strings.HasPrefix(branch, BranchPrefix) {
		branch = BranchPrefix + branch
	}
	return CheckExists(ctx, repoPath, branch)
}

func CheckCommitIfInBranch(ctx context.Context, repoPath, commitId, branch string) (bool, error) {
	result, err := NewCommand("branch", "--contains", commitId, branch).Run(ctx, WithDir(repoPath))
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(result.ReadAsString())) > 0, nil
}
