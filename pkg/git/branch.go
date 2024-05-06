package git

import (
	"context"
	"strings"
)

type Ref struct {
	LastCommitId string
	Name         string
}

func GetAllBranchList(ctx context.Context, repoPath string) ([]Ref, error) {
	cmd := NewCommand("for-each-ref", "--format=%(objectname) %(refname)", BranchPrefix, "--sort=-committerdate")
	pipeResult := cmd.RunWithReadPipe(ctx, WithDir(repoPath))
	ret := make([]Ref, 0)
	err := pipeResult.RangeStringLines(func(_ int, line string) (bool, error) {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 2 {
			ret = append(ret, Ref{
				LastCommitId: fields[0],
				Name:         strings.TrimPrefix(fields[1], BranchPrefix),
			})
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

func DeleteBranch(ctx context.Context, repoPath, branch string, force bool) error {
	deleteCmd := "-d"
	if force {
		deleteCmd = "-D"
	}
	_, err := NewCommand("branch", deleteCmd, "--", branch).Run(ctx, WithDir(repoPath))
	return err
}
