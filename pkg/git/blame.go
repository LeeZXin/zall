package git

import (
	"context"
	"strings"
)

type BlameLine struct {
	Number int     `json:"number"`
	Commit *Commit `json:"commit"`
}

func Blame(ctx context.Context, repoPath, ref, filePath string) ([]BlameLine, error) {
	// 判断是否是二进制文件 这里造成多余的计算 没想到其他的方法
	_, content, _, _, err := GetFileTextContentByRef(ctx, repoPath, ref, filePath)
	if err != nil {
		return nil, err
	}
	if content == "" {
		return []BlameLine{}, nil
	}
	//^d2ba af760 50358 98bd7 11835 459db ddc76 6e719
	cmd := NewCommand("blame", "--date=short").AddDynamicArgs(ref).AddArgs("--").AddDynamicArgs(filePath)
	pipeResult := cmd.RunWithReadPipe(ctx, WithDir(repoPath))
	ret := make([]BlameLine, 0)
	commitMap := make(map[string]*Commit)
	err = pipeResult.RangeStringLines(func(n int, line string) (bool, error) {
		if n > LsTreeLimit {
			return false, nil
		}
		fields := strings.Fields(line)
		// 90549747 (committer 2024-04-26 lineNumber) xxxx
		if len(fields) > 0 {
			sha1 := strings.TrimPrefix(fields[0], "^")
			var (
				c *Commit
				b bool
			)
			c, b = commitMap[sha1]
			if !b {
				commit, _, err := GetCommit(ctx, repoPath, sha1)
				if err != nil {
					return false, err
				}
				c = &commit
				commitMap[sha1] = c
			}
			ret = append(ret, BlameLine{
				Number: n + 1,
				Commit: c,
			})
		}
		return true, nil
	})
	return ret, err
}
