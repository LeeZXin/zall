package git

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type DiffFileType int

const (
	ModifiedFileType DiffFileType = iota + 1
	CreatedFileType
	DeletedFileType
	RenamedFileType
	CopiedFileType
)

const (
	InsertLinePrefix = "+"
	DeleteLinePrefix = "-"
	NormalLinePrefix = " "
	TagLinePrefix    = "@"
)

const EmptyTreeSHA = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

func (t DiffFileType) String() string {
	switch t {
	case CreatedFileType:
		return "created"
	case DeletedFileType:
		return "deleted"
	case RenamedFileType:
		return "renamed"
	case CopiedFileType:
		return "copied"
	case ModifiedFileType:
		return "modified"
	default:
		return "unknown"
	}
}

var (
	EndBytesFlag = []byte("\000")
)

type DiffNumsStatInfo struct {
	FileChangeNums int            `json:"fileChangeNums"`
	InsertNums     int            `json:"insertNums"`
	DeleteNums     int            `json:"deleteNums"`
	Stats          []DiffNumsStat `json:"stats"`
}

type DiffNumsStat struct {
	Path       string
	InsertNums int
	DeleteNums int
}

type DiffFileDetail struct {
	FilePath    string
	OldMode     string
	Mode        string
	IsSubModule bool
	FileType    DiffFileType
	IsBinary    bool
	RenameFrom  string
	RenameTo    string
	CopyFrom    string
	CopyTo      string
	Lines       []DiffLine
}

func newDiffDetail(filePath string) DiffFileDetail {
	return DiffFileDetail{
		FilePath: filePath,
		OldMode:  RegularFileMode.String(),
		Mode:     RegularFileMode.String(),
		FileType: ModifiedFileType,
	}
}

type DiffLine struct {
	LeftNo  int    `json:"leftNo"`
	Prefix  string `json:"prefix"`
	RightNo int    `json:"rightNo"`
	Text    string `json:"text"`
}

func GetFilesDiffCount(ctx context.Context, repoPath, target, head string) (int, error) {
	result, err := NewCommand("diff", "-z", "--name-only").AddDynamicArgs(head+".."+target).Run(ctx, WithDir(repoPath))
	if err != nil {
		if strings.Contains(err.Error(), "no merge base") {
			result, err = NewCommand("diff", "-z", "--name-only").AddDynamicArgs(head, target).Run(ctx, WithDir(repoPath))
		}
	}
	if err != nil {
		return 0, err
	}
	return bytes.Count(result.ReadAsBytes(), EndBytesFlag), nil
}

func GetDiffNumsStat(ctx context.Context, repoPath, target, head string) (DiffNumsStatInfo, error) {
	pipeResult := NewCommand("diff", "--numstat").AddDynamicArgs(head+".."+target).RunWithReadPipe(ctx, WithDir(repoPath))
	stats := make([]DiffNumsStat, 0)
	insertNumsTotal := 0
	deleteNumsTotal := 0
	if err := pipeResult.RangeStringLines(func(_ int, line string) (bool, error) {
		fields := strings.Fields(line)
		if len(fields) == 3 {
			var (
				deleteNums, insertNums int
			)
			if fields[0] != "-" {
				deleteNums, _ = strconv.Atoi(fields[1])
			}
			if fields[1] != "-" {
				insertNums, _ = strconv.Atoi(fields[0])
			}
			insertNumsTotal += insertNums
			deleteNumsTotal += deleteNums
			stats = append(stats, DiffNumsStat{
				Path:       fields[2],
				InsertNums: insertNums,
				DeleteNums: deleteNums,
			})
		}
		return true, nil
	}); err != nil {
		return DiffNumsStatInfo{}, err
	}
	return DiffNumsStatInfo{
		FileChangeNums: len(stats),
		InsertNums:     insertNumsTotal,
		DeleteNums:     deleteNumsTotal,
		Stats:          stats,
	}, nil
}

func GenDiffShortStat(ctx context.Context, repoPath, target, head string) (int, int, int, error) {
	result, err := NewCommand("diff", "--shortstat").AddDynamicArgs(target+".."+head).Run(ctx, WithDir(repoPath))
	if err != nil {
		return 0, 0, 0, err
	}
	var (
		fileChangeNums, insertNums, deleteNums int
	)
	line := strings.TrimSpace(result.ReadAsString())
	lineSplit := strings.Split(line, ",")
	for _, item := range lineSplit {
		fields := strings.Fields(item)
		if strings.Contains(item, "files changed") {
			fileChangeNums, err = strconv.Atoi(fields[0])
			if err != nil {
				return 0, 0, 0, fmt.Errorf("parseInt err:%v", err)
			}
		} else if strings.Contains(item, "insertions") {
			insertNums, err = strconv.Atoi(fields[0])
			if err != nil {
				return 0, 0, 0, fmt.Errorf("parseInt err:%v", err)
			}
		} else if strings.Contains(item, "deletions") {
			deleteNums, err = strconv.Atoi(fields[0])
			if err != nil {
				return 0, 0, 0, fmt.Errorf("parseInt err:%v", err)
			}
		}
	}
	return fileChangeNums, insertNums, deleteNums, nil
}

func GetDiffFileDetail(ctx context.Context, repoPath, target, head, filePath string) (DiffFileDetail, error) {
	pipeResult := NewCommand("diff", "--src-prefix=a/", "--dst-prefix=b/").AddDynamicArgs(head+".."+target).AddArgs("--").AddDynamicArgs(filePath).RunWithReadPipe(ctx, WithDir(repoPath))
	defer pipeResult.ClosePipe()
	reader := bufio.NewReader(pipeResult.Reader())
	c := newDiffDetail(filePath)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return DiffFileDetail{}, err
		}
		if isPrefix {
			continue
		}
		lineStr := strings.TrimSpace(string(line))
		if strings.HasPrefix(lineStr, "diff --git") {
			// nothing
		} else if strings.HasPrefix(lineStr, "index") {
			if strings.HasSuffix(lineStr, SubModuleMode.String()) {
				c.IsSubModule = true
			}
		} else if strings.HasPrefix(lineStr, "---") {
			// nothing
		} else if strings.HasPrefix(lineStr, "+++") {
			// parse hunks
			c.Lines, err = parseHunks(reader)
			if err != nil {
				return DiffFileDetail{}, fmt.Errorf("parse hunks err: %v", err)
			}
		} else if strings.HasPrefix(lineStr, "new mode") {
			c.Mode = strings.TrimSpace(strings.TrimPrefix(lineStr, "new mode"))
			if strings.HasSuffix(lineStr, SubModuleMode.String()) {
				c.IsSubModule = true
			}
		} else if strings.HasPrefix(lineStr, "old mode") {
			c.OldMode = strings.TrimSpace(strings.TrimPrefix(lineStr, "old mode"))
			if strings.HasSuffix(lineStr, SubModuleMode.String()) {
				c.IsSubModule = true
			}
		} else if strings.HasPrefix(lineStr, "new file mode") {
			c.FileType = CreatedFileType
			c.Mode = strings.TrimSpace(strings.TrimPrefix(lineStr, "new file mode"))
			if strings.HasSuffix(lineStr, SubModuleMode.String()) {
				c.IsSubModule = true
			}
		} else if strings.HasPrefix(lineStr, "rename from") {
			c.RenameFrom = strings.TrimSpace(strings.TrimPrefix(lineStr, "rename from"))
			c.FileType = RenamedFileType
		} else if strings.HasPrefix(lineStr, "rename to") {
			c.RenameTo = strings.TrimSpace(strings.TrimPrefix(lineStr, "rename to"))
			c.FileType = RenamedFileType
		} else if strings.HasPrefix(lineStr, "copy from") {
			c.CopyFrom = strings.TrimSpace(strings.TrimPrefix(lineStr, "copy from"))
			c.FileType = CopiedFileType
		} else if strings.HasPrefix(lineStr, "copy to") {
			c.CopyTo = strings.TrimSpace(strings.TrimPrefix(lineStr, "copy to"))
			c.FileType = CopiedFileType
		} else if strings.HasPrefix(lineStr, "deleted") {
			c.FileType = DeletedFileType
		} else if strings.HasPrefix(lineStr, "Binary") {
			c.IsBinary = true
		}
	}
	return c, nil
}

func parseHunks(reader *bufio.Reader) ([]DiffLine, error) {
	insertionNums := 0
	deletionNums := 0
	ret := make([]DiffLine, 0)
	var (
		leftNo, rightNo, index int

		lastLeftNo = -1
	)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if isPrefix {
			continue
		}
		lineStr := string(line)
		if strings.HasPrefix(lineStr, "@@") {
			leftNo, _, rightNo, _, err = parseHunkString(lineStr)
			if err != nil {
				return nil, err
			}
			ret = append(ret, DiffLine{
				LeftNo:  -1,
				Prefix:  TagLinePrefix,
				RightNo: -1,
				Text:    lineStr,
			})
		} else if strings.HasPrefix(lineStr, InsertLinePrefix) {
			left := lastLeftNo
			if index > 0 && ret[index-1].Prefix == InsertLinePrefix {
				left = ret[index-1].LeftNo + 1
			}
			ret = append(ret, DiffLine{
				LeftNo:  left,
				Prefix:  InsertLinePrefix,
				RightNo: rightNo,
				Text:    lineStr[1:],
			})
			rightNo++
			insertionNums++
		} else if strings.HasPrefix(lineStr, DeleteLinePrefix) {
			right := rightNo
			if index > 0 {
				if ret[index-1].Prefix == DeleteLinePrefix {
					right = ret[index-1].RightNo + 1
				} else {
					lastLeftNo = leftNo
				}
			} else {
				lastLeftNo = leftNo
			}
			ret = append(ret, DiffLine{
				LeftNo:  leftNo,
				Prefix:  DeleteLinePrefix,
				RightNo: right,
				Text:    lineStr[1:],
			})
			leftNo++
			deletionNums++
		} else {
			ret = append(ret, DiffLine{
				LeftNo:  leftNo,
				Prefix:  NormalLinePrefix,
				RightNo: rightNo,
				Text:    lineStr[1:],
			})
			leftNo++
			rightNo++
			lastLeftNo = leftNo
		}
		index++
	}
	return ret, nil
}

func parseHunkString(line string) (int, int, int, int, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 || fields[0] != "@@" || fields[3] != "@@" {
		return 0, 0, 0, 0, errors.New("invalid @@ format")
	}
	var (
		o1, o2, n1, n2 int
		err            error
	)
	o := strings.Split(fields[1][1:], ",")
	if len(o) >= 2 {
		o1, err = strconv.Atoi(o[0])
		if err != nil {
			return 0, 0, 0, 0, err
		}
		o2, err = strconv.Atoi(o[1])
		if err != nil {
			return 0, 0, 0, 0, err
		}
	} else {
		o1, err = strconv.Atoi(o[0])
		if err != nil {
			return 0, 0, 0, 0, err
		}
	}
	n := strings.Split(fields[2][1:], ",")
	if len(n) >= 2 {
		n1, err = strconv.Atoi(n[0])
		if err != nil {
			return 0, 0, 0, 0, err
		}
		n2, err = strconv.Atoi(n[1])
		if err != nil {
			return 0, 0, 0, 0, err
		}
	} else {
		n1, err = strconv.Atoi(n[0])
		if err != nil {
			return 0, 0, 0, 0, err
		}
	}
	return o1, o2, n1, n2, nil
}

func GenDiffDetailRowData(ctx context.Context, repoPath, target, head, filePath string) (string, error) {
	result, err := NewCommand("diff", "--src-prefix=a/", "--dst-prefix=b/").
		AddDynamicArgs(target+".."+head).
		AddArgs("--").
		AddDynamicArgs(filePath).
		Run(ctx, WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return result.ReadAsString(), nil
}
