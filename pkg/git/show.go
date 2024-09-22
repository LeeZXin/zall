package git

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf-utils/sortutil"
	"github.com/LeeZXin/zsf-utils/typesniffer"
	"strconv"
	"strings"
)

const (
	// FileBlobSizeLimit 文件大小限制
	FileBlobSizeLimit = 1024 * 1024 * 2
	LsTreeLimit       = 500
)

func ShowFileTextContentByCommitId(ctx context.Context, repoPath, commitId, filePath string, startLine, limit int) ([]string, error) {
	pipeResult := NewCommand("show", "--text").AddDynamicArgs(fmt.Sprintf("%s:%s", commitId, filePath)).
		RunWithReadPipe(ctx, WithDir(repoPath))
	ret := make([]string, 0)
	endLine := startLine + limit
	if err := pipeResult.RangeStringLines(func(i int, line string) (bool, error) {
		if i < startLine {
			return true, nil
		}
		if limit < 0 || (i >= startLine && i < endLine) {
			ret = append(ret, line)
			return true, nil
		}
		return false, nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func GetRefFilesCountAndSize(ctx context.Context, repoPath, refName string) (int, int64, error) {
	result := NewCommand("ls-tree", "--full-tree", "-r", "-l").AddDynamicArgs(refName).
		RunWithReadPipe(ctx, WithDir(repoPath))
	var (
		fileCount = 0

		fileSize int64 = 0
	)
	if err := result.RangeStringLines(func(_ int, line string) (bool, error) {
		fields := strings.Fields(line)
		if len(fields) == 5 {
			fileCount++
			size, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				fileSize += size
			}
		}
		return true, nil
	}); err != nil {
		return 0, 0, err
	}
	return fileCount, fileSize, nil
}

type LsTreeRet struct {
	Mode FileMode
	Path string
	Size int64
	Blob string
}

func (c LsTreeRet) CompareTo(c2 LsTreeRet) bool {
	if c.Mode == DirectoryMode && c2.Mode == DirectoryMode {
		return c.Path <= c2.Path
	} else if c.Mode == DirectoryMode {
		return true
	} else if c2.Mode == DirectoryMode {
		return false
	} else {
		return c.Path <= c2.Path
	}
}

// LsTreeWithoutRecurse git ls-tree --full-tree -l master
func LsTreeWithoutRecurse(ctx context.Context, repoPath, ref, dir string) ([]LsTreeRet, bool, error) {
	cmd := NewCommand("ls-tree", "--full-tree", "-l").AddDynamicArgs(ref)
	if dir != "" {
		cmd.AddArgs("--", dir)
	}
	result := cmd.RunWithReadPipe(ctx, WithDir(repoPath))
	ret := make([]LsTreeRet, 0)
	if err := result.RangeStringLines(func(n int, line string) (bool, error) {
		fields := strings.Fields(line)
		if len(fields) == 5 {
			size, _ := strconv.ParseInt(fields[3], 10, 64)
			ret = append(ret, LsTreeRet{
				Mode: FileMode(fields[0]),
				Path: fields[4],
				Size: size,
				Blob: fields[2],
			})
		}
		// 限制展示500的数量
		if n >= LsTreeLimit {
			return false, nil
		}
		return true, nil
	}); err != nil {
		return nil, false, err
	}
	return ret, false, nil
}

func CountObjects(ctx context.Context, repoPath string) (int64, float64, error) {
	result, err := NewCommand("count-objects").Run(ctx, WithDir(repoPath))
	if err != nil {
		return 0, 0, err
	}
	line := result.ReadAsString()
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) != 4 {
		return 0, 0, errors.New("unknown count-objects output")
	}
	objectsCount, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, 0, errors.New("unknown count-objects fields[0]")
	}
	objectsSize, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return 0, 0, errors.New("unknown count-objects fields[2]")
	}
	return objectsCount, objectsSize, nil
}

type FileCommit struct {
	LsTreeRet
	Commit
}

func LsTreeCommit(ctx context.Context, repoPath, ref string, dir string) ([]FileCommit, error) {
	lsRet, _, err := LsTreeWithoutRecurse(ctx, repoPath, ref, dir)
	if err != nil {
		return nil, err
	}
	if len(lsRet) > 1 {
		sortutil.SliceStable(lsRet)
	}
	commits := make([]FileCommit, 0, len(lsRet))
	for _, ret := range lsRet {
		commit, err := GetFileLatestCommit(ctx, repoPath, ref, ret.Path)
		if err != nil {
			return nil, err
		}
		commits = append(commits, FileCommit{
			LsTreeRet: ret,
			Commit:    commit,
		})
	}
	return commits, nil
}

func LsTreeBlob(ctx context.Context, repoPath, ref string, dir string) ([]LsTreeRet, error) {
	lsRet, _, err := LsTreeWithoutRecurse(ctx, repoPath, ref, dir)
	if err != nil {
		return nil, err
	}
	sortutil.SliceStable(lsRet)
	return lsRet, nil
}

func GetFileContentByBlob(ctx context.Context, repoPath, blob string) ([]byte, error) {
	result, err := NewCommand("show").AddDynamicArgs(blob).Run(ctx, WithDir(repoPath))
	if err != nil {
		return nil, err
	}
	return result.ReadAsBytes(), nil
}

func GetFileTextContentByRef(ctx context.Context, repoPath, ref, filePath string) (FileMode, string, int64, bool, bool, error) {
	result, err := NewCommand("ls-tree", "--full-tree", "-l").
		AddDynamicArgs(ref).
		AddArgs("--").
		AddDynamicArgs(filePath).
		Run(ctx, WithDir(repoPath))
	if err != nil {
		return "", "", 0, false, false, err
	}
	ret := result.ReadAsString()
	if ret == "" {
		return "", "", 0, false, false, nil
	}
	fields := strings.Fields(strings.TrimSpace(ret))
	if len(fields) < 4 {
		return "", "", 0, false, false, errors.New("unknown format")
	}
	blob := fields[2]
	mode := fields[0]
	size, err := strconv.ParseInt(fields[3], 10, 64)
	var (
		content string
		isText  bool
	)
	if err == nil && size < FileBlobSizeLimit {
		bc, err := GetFileContentByBlob(ctx, repoPath, blob)
		if err != nil {
			return "", "", 0, false, false, err
		}
		isText = typesniffer.DetectContentType(bc).IsText()
		if isText {
			content = string(bc)
		}
	}
	return FileMode(mode), content, size, isText, true, nil
}
