package git

import (
	"context"
	"strings"
)

type Ref struct {
	Sha  string
	Name string
}

func GetAllBranchList(ctx context.Context, repoPath string) ([]Ref, error) {
	cmd := NewCommand("for-each-ref", "--format=%(objectname) %(refname)", BranchPrefix, "--sort=-committerdate")
	pipeResult := cmd.RunWithReadPipe(ctx, WithDir(repoPath))
	ret := make([]Ref, 0)
	err := pipeResult.RangeStringLines(func(_ int, line string) (bool, error) {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 2 {
			ret = append(ret, Ref{
				Sha:  fields[0],
				Name: strings.TrimPrefix(fields[1], BranchPrefix),
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

func DeleteBranch(ctx context.Context, repoPath, branch string, force bool) error {
	deleteCmd := "-d"
	if force {
		deleteCmd = "-D"
	}
	_, err := NewCommand("branch", deleteCmd, "--").AddDynamicArgs(branch).Run(ctx, WithDir(repoPath))
	return err
}
