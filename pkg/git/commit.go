package git

import (
	"bufio"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/pkg/git/gitenv"
	"github.com/LeeZXin/zall/pkg/git/signature"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RefType string

const (
	BranchType RefType = "branch"
	CommitType RefType = "commit"
	TagType    RefType = "tag"
)

func (t RefType) IsValid() bool {
	switch t {
	case BranchType, CommitType, TagType:
		return true
	default:
		return false
	}
}

func (t RefType) PackRef(ref string) string {
	switch t {
	case BranchType:
		return BranchPrefix + ref
	case TagType:
		return TagPrefix + ref
	default:
		return ref
	}
}

func (t RefType) String() string {
	return string(t)
}

var (
	ShortCommitIdPattern = regexp.MustCompile(`^[0-9a-f]{7}$`)
)

type Commit struct {
	Id            string
	Tree          Tree
	Parent        []string
	Author        User
	AuthorSigTime time.Time
	Committer     User
	CommitSigTime time.Time
	CommitSig     string
	CommitMsg     string
	Tag           *Tag
	Payload       string
}

func newCommit(id string) Commit {
	return Commit{
		Id:     id,
		Parent: make([]string, 0),
	}
}

type Tag struct {
	Id        string
	Object    string
	Typ       string
	Tag       string
	Tagger    User
	TagTime   time.Time
	CommitMsg string
	Sig       string
	Payload   string
}

type Tree struct {
	Id string `json:"id"`
}

func NewTree(id string) Tree {
	return Tree{
		Id: id,
	}
}

func GetBranchCommitId(ctx context.Context, repoPath string, name string) (string, error) {
	cmd := NewCommand("rev-parse", "--verify").AddDynamicArgs(name)
	result, err := cmd.Run(ctx, WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), nil
}

func CheckRefIsCommit(ctx context.Context, repoPath string, name string) bool {
	return CheckExists(ctx, repoPath, name)
}

func GetCommitByCommitId(ctx context.Context, repoPath string, commitId string) (Commit, error) {
	c := newCommit(commitId)
	return c, CatFileBatch(ctx, repoPath, commitId, func(r io.Reader, closer PipeResultCloser) error {
		defer closer.ClosePipe()
		reader := bufio.NewReader(r)
		var (
			typ  string
			size int64
		)
		for {
			line, isPrefix, err := reader.ReadLine()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("read line err: %v", err)
			}
			if isPrefix {
				continue
			}
			commitId, typ, size, err = readBatchLine(string(line))
			if err != nil {
				return fmt.Errorf("readBatchLine err: %v", err)
			}
			c.Id = commitId
			break
		}
		switch typ {
		case CommitType.String():
			return generateCommit(io.LimitReader(reader, size), &c)
		default:
			return nil
		}
	})
}

func GetCommitByTag(ctx context.Context, repoPath string, tag string) (c Commit, e error) {
	e = CatFileBatch(ctx, repoPath, tag, func(r io.Reader, closer PipeResultCloser) error {
		defer closer.ClosePipe()
		reader := bufio.NewReader(r)
		var (
			typ  string
			size int64
			id   string
		)
		for {
			line, isPrefix, err := reader.ReadLine()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("read line err: %v", err)
			}
			if isPrefix {
				continue
			}
			id, typ, size, err = readBatchLine(string(line))
			if err != nil {
				return fmt.Errorf("%s readBatchLine err: %v", tag, err)
			}
			break
		}
		switch typ {
		case TagType.String():
			t := &Tag{
				Id:  id,
				Tag: tag,
			}
			err := generateTag(io.LimitReader(reader, size), t)
			if err != nil {
				return fmt.Errorf("parse Tag err: %v", err)
			}
			if t.Object == "" {
				return fmt.Errorf("%s object is empty", tag)
			}
			c, err = GetCommitByCommitId(ctx, repoPath, t.Object)
			c.Tag = t
			return err
		default:
			return nil
		}
	})
	return
}

func generateTag(r io.Reader, tag *Tag) error {
	reader := bufio.NewReader(r)
	commitMsg := strings.Builder{}
	payload := strings.Builder{}
	defer func() {
		tag.CommitMsg = commitMsg.String()
		tag.Payload = payload.String()
	}()
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read line err: %v", err)
		}
		if isPrefix {
			continue
		}
		lineStr := strings.TrimSpace(string(line))
		if lineStr == signature.StartGPGSigLineTag ||
			lineStr == signature.StartSSHSigLineTag {
			sigPayLoad := strings.Builder{}
			sigPayLoad.WriteString(lineStr + "\n")
			for {
				line, isPrefix, err = reader.ReadLine()
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return fmt.Errorf("read line err: %v", err)
				}
				if isPrefix {
					continue
				}
				lineStr = strings.TrimSpace(string(line))
				sigPayLoad.WriteString(lineStr + "\n")
				if lineStr == signature.EndGPGSigLineTag || lineStr == signature.EndSSHSigLineTag {
					break
				}
			}
			tag.Sig = sigPayLoad.String()
			continue
		}
		payload.WriteString(lineStr + "\n")
		fields := strings.Fields(lineStr)
		if len(fields) < 1 {
			continue
		}
		switch fields[0] {
		case "object":
			tag.Object = fields[1]
		case "type":
			tag.Typ = fields[1]
		case "tag":
			tag.Tag = fields[1]
		case "tagger":
			tag.Tagger, tag.TagTime = parseUserAndTime(fields[1:])
		default:
			commitMsg.WriteString(lineStr + "\n")
		}
	}
}

func generateCommit(r io.Reader, commit *Commit) error {
	reader := bufio.NewReader(r)
	commitMsg := strings.Builder{}
	payload := strings.Builder{}
	defer func() {
		commit.CommitMsg = commitMsg.String()
		commit.Payload = payload.String()
	}()
	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix {
			continue
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read line err: %v", err)
		}
		rowLine := string(line)
		lineStr := strings.TrimSpace(rowLine)
		fields := strings.Fields(lineStr)
		// 记录非签名数据
		if len(fields) < 1 || fields[0] != "gpgsig" {
			payload.WriteString(lineStr)
			payload.WriteString("\n")
		}
		if len(fields) < 1 {
			continue
		}
		switch fields[0] {
		case "tree":
			commit.Tree = NewTree(fields[1])
		case "parent":
			commit.Parent = append(commit.Parent, fields[1])
		case "author":
			commit.Author, commit.AuthorSigTime = parseUserAndTime(fields[1:])
		case "committer":
			commit.Committer, commit.CommitSigTime = parseUserAndTime(fields[1:])
		case "gpgsig":
			if len(fields) <= 1 {
				continue
			}
			sigPayload := strings.Builder{}
			sigPayload.WriteString(strings.TrimPrefix(lineStr, "gpgsig ") + "\n")
			for {
				line, isPrefix, err = reader.ReadLine()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("read gpgsig err: %v", err)
				}
				if isPrefix {
					continue
				}
				lineStr = strings.TrimSpace(string(line))
				sigPayload.WriteString(lineStr + "\n")
				if lineStr == signature.EndGPGSigLineTag || lineStr == signature.EndSSHSigLineTag {
					break
				}
			}
			commit.CommitSig = sigPayload.String()
		default:
			commitMsg.WriteString(lineStr + "\n")
		}
	}
}

func parseUserAndTime(f []string) (User, time.Time) {
	u := User{}
	l := len(f)
	if l >= 1 {
		u.Account = f[0]
	}
	if l >= 2 {
		u.Email = f[1][1 : len(f[1])-1]
	}
	var eventTime time.Time
	if l >= 3 {
		firstChar := f[2][0]
		if firstChar >= 48 && firstChar <= 57 {
			i, err := strconv.ParseInt(f[2], 10, 64)
			if err == nil {
				eventTime = time.Unix(i, 0)
			}
			if l >= 4 {
				zone := f[3]
				h, herr := strconv.ParseInt(zone[0:3], 10, 64)
				m, merr := strconv.ParseInt(zone[3:], 10, 64)
				if herr == nil && merr == nil {
					if h < 0 {
						m = -m
					}
					eventTime = eventTime.In(time.FixedZone("", int(h*3600+m*60)))
				}
			}
		} else {
			i, err := time.Parse(TimeLayout, f[2])
			if err == nil {
				eventTime = i
			}
		}
	}
	return u, eventTime
}

func GetFullShaCommitId(ctx context.Context, repoPath, short string) (string, error) {
	if ShortCommitIdPattern.MatchString(short) {
		line, _, _, err := CatFileBatchCheck(ctx, repoPath, short)
		return line, err
	}
	return short, nil
}

func GetGitDiffCommitList(ctx context.Context, repoPath, target, head string) ([]Commit, error) {
	result, err := NewCommand("log", PrettyLogFormat).AddDynamicArgs(head+".."+target).AddArgs("--max-count=100", "--").
		Run(ctx, WithDir(repoPath))
	if err != nil {
		return nil, err
	}
	idList := strings.Fields(strings.TrimSpace(result.ReadAsString()))
	return listutil.Map(idList, func(t string) (Commit, error) {
		return GetCommitByCommitId(ctx, repoPath, t)
	})
}

func GetGitLogCommitList(ctx context.Context, repoPath, ref string, skip, limit int) ([]Commit, error) {
	result, err := NewCommand("log", PrettyLogFormat).
		AddDynamicArgs(ref).
		AddArgs(
			fmt.Sprintf("--max-count=%d", limit),
			fmt.Sprintf("--skip=%d", skip), "--",
		).Run(ctx, WithDir(repoPath))
	if err != nil {
		return nil, err
	}
	idList := strings.Fields(strings.TrimSpace(result.ReadAsString()))
	return listutil.Map(idList, func(t string) (Commit, error) {
		return GetCommitByCommitId(ctx, repoPath, t)
	})
}

func GetFileLatestCommit(ctx context.Context, repoPath, ref, filePath string) (Commit, error) {
	result, err := NewCommand("log", PrettyLogFormat, "-1").AddDynamicArgs(ref).AddArgs("--").AddDynamicArgs(filePath).
		Run(ctx, WithDir(repoPath))
	if err != nil {
		return Commit{}, err
	}
	commitId := strings.TrimSpace(result.ReadAsString())
	return GetCommitByCommitId(ctx, repoPath, commitId)
}

func GetCommit(ctx context.Context, repoPath string, ref string) (Commit, string, error) {
	if CheckRefIsTag(ctx, repoPath, ref) {
		if !strings.HasPrefix(ref, TagPrefix) {
			ref = TagPrefix + ref
		}
		commit, err := GetCommitByTag(ctx, repoPath, ref)
		return commit, ref, err
	}
	if CheckRefIsBranch(ctx, repoPath, ref) {
		if !strings.HasPrefix(ref, BranchPrefix) {
			ref = BranchPrefix + ref
		}
		commitId, err := GetBranchCommitId(ctx, repoPath, ref)
		if err != nil {
			return Commit{}, "", err
		}
		commit, err := GetCommitByCommitId(ctx, repoPath, commitId)
		return commit, ref, err
	}
	if CheckRefIsCommit(ctx, repoPath, ref) {
		commit, err := GetCommitByCommitId(ctx, repoPath, ref)
		return commit, ref, err
	}
	return Commit{}, "", fmt.Errorf("%s unsupported type", ref)
}

type DetectForcePushEnv struct {
	ObjectDirectory              string
	AlternativeObjectDirectories string
	QuarantinePath               string
}

func DetectForcePush(ctx context.Context, repoPath, oldCommitId, newCommitId string, env DetectForcePushEnv) (bool, error) {
	result, err := NewCommand("rev-list", "--max-count=1").
		AddDynamicArgs(oldCommitId, "^"+newCommitId).
		Run(ctx, WithDir(repoPath),
			WithEnv(util.JoinFields(
				gitenv.EnvObjectDirectory, env.ObjectDirectory,
				gitenv.EnvAlternativeObjectDirectories, env.AlternativeObjectDirectories,
				gitenv.EnvQuarantinePath, env.QuarantinePath,
			)))
	if err != nil {
		return false, err
	}
	return len(result.ReadAsBytes()) > 0, nil
}
