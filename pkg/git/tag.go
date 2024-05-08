package git

import (
	"context"
	"strings"
)

func GetAllTagList(ctx context.Context, repoPath string) ([]Ref, error) {
	cmd := NewCommand("for-each-ref", "--format=%(objectname) %(objecttype) %(refname)", TagPrefix, "--sort=-taggerdate")
	pipeResult := cmd.RunWithReadPipe(ctx, WithDir(repoPath))
	ret := make([]Ref, 0)
	err := pipeResult.RangeStringLines(func(_ int, line string) (bool, error) {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 3 {
			if fields[1] == "tag" {
				ret = append(ret, Ref{
					Sha:  fields[0],
					Name: strings.TrimPrefix(fields[2], TagPrefix),
				})
			}
		}
		return true, nil
	})
	return ret, err
}

func DeleteTag(ctx context.Context, repoPath string, tag string) error {
	_, err := NewCommand("tag", "-d").AddDynamicArgs(tag).Run(ctx, WithDir(repoPath))
	return err
}

func CheckRefIsTag(ctx context.Context, repoPath string, tag string) bool {
	if !strings.HasPrefix(tag, TagPrefix) {
		tag = TagPrefix + tag
	}
	return CheckExists(ctx, repoPath, tag)
}
