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
	FileChangeNums int `json:"fileChangeNums"`
	InsertNums     int `json:"insertNums"`
	DeleteNums     int `json:"deleteNums"`
	Stats          []DiffNumsStat
}

type DiffNumsStat struct {
	Path       string
	TotalNums  int
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
	Index   int    `json:"index"`
	LeftNo  int    `json:"leftNo"`
	Prefix  string `json:"prefix"`
	RightNo int    `json:"rightNo"`
	Text    string `json:"text"`
}

func GetFilesDiffCount(ctx context.Context, repoPath, target, head string) (int, error) {
	result, err := NewCommand("diff", "-z", "--name-only", target+".."+head, "--").Run(ctx, WithDir(repoPath))
	if err != nil {
		if strings.Contains(err.Error(), "no merge base") {
			result, err = NewCommand("diff", "-z", "--name-only", target, head, "--").Run(ctx, WithDir(repoPath))
		}
	}
	if err != nil {
		return 0, err
	}
	return bytes.Count(result.ReadAsBytes(), EndBytesFlag), nil
}

func GetDiffNumsStat(ctx context.Context, repoPath, target, head string) (DiffNumsStatInfo, error) {
	pipeResult := NewCommand("diff", "--numstat", target+".."+head, "--").RunWithReadPipe(ctx, WithDir(repoPath))
	ret := make([]DiffNumsStat, 0)
	insertNumsTotal := 0
	deleteNumsTotal := 0
	if err := pipeResult.RangeStringLines(func(_ int, line string) (bool, error) {
		fields := strings.Fields(line)
		if len(fields) == 3 {
			deleteNums, err := strconv.Atoi(fields[0])
			if err != nil {
				return false, fmt.Errorf("parseInt err: %v", deleteNums)
			}
			insertNums, err := strconv.Atoi(fields[1])
			if err != nil {
				return false, fmt.Errorf("parseInt err: %v", insertNums)
			}
			insertNumsTotal += insertNums
			deleteNumsTotal += deleteNums
			ret = append(ret, DiffNumsStat{
				Path:       fields[2],
				InsertNums: insertNums,
				DeleteNums: deleteNums,
				TotalNums:  insertNums + deleteNums,
			})
		}
		return true, nil
	}); err != nil {
		return DiffNumsStatInfo{}, err
	}
	return DiffNumsStatInfo{
		FileChangeNums: len(ret),
		InsertNums:     insertNumsTotal,
		DeleteNums:     deleteNumsTotal,
		Stats:          ret,
	}, nil
}

func GenDiffShortStat(ctx context.Context, repoPath, target, head string) (int, int, int, error) {
	result, err := NewCommand("diff", "--shortstat", target+".."+head, "--").Run(ctx, WithDir(repoPath))
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
	pipeResult := NewCommand("diff", "--src-prefix=a/", "--dst-prefix=b/", head+".."+target, "--", filePath).RunWithReadPipe(ctx, WithDir(repoPath))
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
		index, leftNo, rightNo int
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
				Index:   index,
				LeftNo:  leftNo - 1,
				Prefix:  TagLinePrefix,
				RightNo: rightNo - 1,
				Text:    lineStr,
			})
		} else if strings.HasPrefix(lineStr, "+") {
			ret = append(ret, DiffLine{
				Index:   index,
				LeftNo:  leftNo,
				Prefix:  InsertLinePrefix,
				RightNo: rightNo,
				Text:    lineStr[1:],
			})
			rightNo++
			insertionNums++
		} else if strings.HasPrefix(lineStr, "-") {
			ret = append(ret, DiffLine{
				Index:   index,
				LeftNo:  leftNo,
				Prefix:  DeleteLinePrefix,
				RightNo: rightNo,
				Text:    lineStr[1:],
			})
			leftNo++
			deletionNums++
		} else {
			ret = append(ret, DiffLine{
				Index:   index,
				LeftNo:  leftNo,
				Prefix:  NormalLinePrefix,
				RightNo: rightNo,
				Text:    lineStr[1:],
			})
			leftNo++
			rightNo++
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
	result, err := NewCommand("diff", "--src-prefix=a/", "--dst-prefix=b/", target+".."+head, "--", filePath).Run(ctx, WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return result.ReadAsString(), nil
}
