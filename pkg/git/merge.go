package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/pkg/git/gitenv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/common"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	MergeBranch    = "base"
	TrackingBranch = "tracking"
)

var (
	escapedSymbols = regexp.MustCompile(`([*[?! \\])`)
)

type DiffRefsInfo struct {
	OriginHead    string           `json:"originHead"`
	OriginTarget  string           `json:"originTarget"`
	Target        string           `json:"target"`
	Head          string           `json:"head"`
	TargetCommit  Commit           `json:"targetCommit"`
	HeadCommit    Commit           `json:"headCommit"`
	Commits       []Commit         `json:"commits"`
	NumFiles      int              `json:"numFiles"`
	MergeBase     string           `json:"mergeBase"`
	DiffNumsStats DiffNumsStatInfo `json:"diffNumsStats"`
	ConflictFiles []string         `json:"conflictFiles"`
	CanMerge      bool             `json:"canMerge"`
}

type DiffCommitsInfo struct {
	Commit        Commit           `json:"commit"`
	NumFiles      int              `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfo `json:"diffNumsStats"`
}

// IsMergeAble 是否可合并
func (i *DiffRefsInfo) IsMergeAble() bool {
	if !strings.HasPrefix(i.OriginHead, BranchPrefix) {
		return false
	}
	return len(i.Commits) > 0 && len(i.ConflictFiles) == 0
}

type MergeRepoOpts struct {
	RepoId        int64
	PrId          int64
	PusherAccount string
	PusherEmail   string
	Message       string
	AppUrl        string
}

func GetDiffRefsInfo(ctx context.Context, repoPath, target, head string) (DiffRefsInfo, error) {
	pr := DiffRefsInfo{}
	pr.OriginTarget, pr.OriginHead = target, head
	var err error
	pr.HeadCommit, pr.Head, err = GetCommit(ctx, repoPath, head)
	if err != nil {
		return DiffRefsInfo{}, err
	}
	pr.TargetCommit, pr.Target, err = GetCommit(ctx, repoPath, target)
	if err != nil {
		return DiffRefsInfo{}, err
	}
	// 这里要反过来 git log 查看target的提交记录 不是head的提交记录
	pr.Commits, err = GetGitDiffCommitList(ctx, repoPath, pr.TargetCommit.Id, pr.HeadCommit.Id)
	if err != nil {
		return DiffRefsInfo{}, err
	}
	pr.NumFiles, err = GetFilesDiffCount(ctx, repoPath, pr.TargetCommit.Id, pr.HeadCommit.Id)
	if err != nil {
		return DiffRefsInfo{}, err
	}
	pr.DiffNumsStats, err = GetDiffNumsStat(ctx, repoPath, pr.TargetCommit.Id, pr.HeadCommit.Id)
	if err != nil {
		return DiffRefsInfo{}, err
	}
	pr.MergeBase, err = MergeBase(ctx, repoPath, pr.TargetCommit.Id, pr.HeadCommit.Id)
	if err != nil {
		return DiffRefsInfo{}, err
	}
	if strings.HasPrefix(pr.OriginHead, BranchPrefix) {
		pr.ConflictFiles, err = findConflictFiles(ctx, repoPath, pr.OriginHead, pr.TargetCommit.Id, pr.MergeBase)
		if err != nil {
			return DiffRefsInfo{}, err
		}
	} else {
		pr.ConflictFiles = []string{}
	}
	return pr, nil
}

func GetDiffCommitsInfo(ctx context.Context, repoPath, commitId string) (DiffCommitsInfo, error) {
	pr := DiffCommitsInfo{}
	var err error
	pr.Commit, err = GetCommitByCommitId(ctx, repoPath, commitId)
	if err != nil {
		return DiffCommitsInfo{}, err
	}
	var (
		targetCommitId string
	)
	parentLen := len(pr.Commit.Parent)
	if parentLen == 0 {
		targetCommitId = EmptyTreeSHA
	} else {
		targetCommitId = pr.Commit.Parent[0]
	}
	pr.NumFiles, err = GetFilesDiffCount(ctx, repoPath, commitId, targetCommitId)
	if err != nil {
		return DiffCommitsInfo{}, err
	}
	pr.DiffNumsStats, err = GetDiffNumsStat(ctx, repoPath, commitId, targetCommitId)
	if err != nil {
		return DiffCommitsInfo{}, err
	}
	return pr, nil
}

func CanMerge(ctx context.Context, repoPath, target, head string) (bool, error) {
	headCommit, _, err := GetCommit(ctx, repoPath, head)
	if err != nil {
		return false, err
	}
	targetCommit, _, err := GetCommit(ctx, repoPath, target)
	if err != nil {
		return false, err
	}
	// 这里要反过来 git log 查看target的提交记录 不是head的提交记录
	commits, err := GetGitDiffCommitList(ctx, repoPath, headCommit.Id, targetCommit.Id)
	if err != nil {
		return false, err
	}
	mergeBase, err := MergeBase(ctx, repoPath, targetCommit.Id, headCommit.Id)
	if err != nil {
		return false, err
	}
	if strings.HasPrefix(head, BranchPrefix) {
		conflictFiles, err := findConflictFiles(ctx, repoPath, head, targetCommit.Id, mergeBase)
		if err != nil {
			return false, err
		}
		return len(commits) > 0 && len(conflictFiles) == 0, nil
	}
	return false, nil
}

func Merge(ctx context.Context, repoPath, target, head string, opts MergeRepoOpts) error {
	info, err := GetDiffRefsInfo(ctx, repoPath, target, head)
	if err != nil {
		return err
	}
	return doMerge(ctx, repoPath, info, opts)
}

func doMerge(ctx context.Context, repoPath string, pr DiffRefsInfo, opts MergeRepoOpts) error {
	if len(pr.Commits) == 0 {
		return errors.New("nothing to commit")
	}
	tempDir := filepath.Join(TempDir(), "merge-"+util.RandomIdWithTime())
	defer util.RemoveAll(tempDir)
	if err := prepare4Merge(ctx, repoPath, tempDir, pr.OriginHead, pr.TargetCommit.Id); err != nil {
		return err
	}
	infoPath := filepath.Join(tempDir, ".git", "info")
	if err := os.MkdirAll(infoPath, 0o700); err != nil {
		return fmt.Errorf("unable to create .git/info in tmpBasePath: %w", err)
	}
	sparseCheckout, err := os.OpenFile(filepath.Join(infoPath, "sparse-checkout"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("unable to write .git/info/sparse-checkout file in tmpBasePath: %w", err)
	}
	defer sparseCheckout.Close()
	trees, err := getDiffTreeForMerge(ctx, tempDir, TrackingBranch, MergeBranch)
	if err != nil {
		return fmt.Errorf("unable to get diff tree in tmpBasePath: %w", err)
	}
	for _, tree := range trees {
		if _, err = sparseCheckout.WriteString(tree); err != nil {
			return fmt.Errorf("unable to write to sparseCheckout in tmpBasePath: %w", err)
		}
	}
	if err = setLocalConfig(ctx, tempDir, "filter.lfs.process", ""); err != nil {
		return err
	}
	if err = setLocalConfig(ctx, tempDir, "filter.lfs.required", "false"); err != nil {
		return err
	}
	if err = setLocalConfig(ctx, tempDir, "filter.lfs.clean", ""); err != nil {
		return err
	}
	if err = setLocalConfig(ctx, tempDir, "filter.lfs.smudge", ""); err != nil {
		return err
	}
	if err = setLocalConfig(ctx, tempDir, "core.sparseCheckout", "true"); err != nil {
		return err
	}
	if _, err = NewCommand("read-tree", "HEAD").Run(ctx, WithDir(tempDir)); err != nil {
		return err
	}
	if _, err = NewCommand("merge", "--no-ff", "--no-commit", TrackingBranch).
		Run(ctx, WithDir(tempDir)); err != nil {
		return fmt.Errorf("git merge err: %v", err)
	}
	mergeCmd := NewCommand("commit", "--no-gpg-sign", "-m", opts.Message)
	if _, err = mergeCmd.Run(ctx, WithDir(tempDir)); err != nil {
		return err
	}
	if _, err = NewCommand("push", "origin", MergeBranch+":"+pr.Head).
		Run(ctx,
			WithDir(tempDir),
			WithEnv(
				util.JoinFields(
					gitenv.EnvRepoId, strconv.FormatInt(opts.RepoId, 10),
					gitenv.EnvAppUrl, opts.AppUrl,
					gitenv.EnvHookToken, HookToken(),
					gitenv.EnvPrId, strconv.FormatInt(opts.PrId, 10),
					gitenv.EnvPusherAccount, opts.PusherAccount,
					gitenv.EnvPusherEmail, opts.PusherEmail,
					gitenv.EnvHookUrl, fmt.Sprintf("http://127.0.0.1:%d", common.HttpServerPort()),
				),
			),
		); err != nil {
		return fmt.Errorf("git push: %v", err)
	}
	return nil
}

func prepare4Merge(ctx context.Context, repoPath string, tempDir string, originHead, targetCommitId string) error {
	if err := initEmptyRepository(ctx, tempDir, false); err != nil {
		return err
	}
	if _, err := NewCommand("remote", "add", "-t", originHead, "-m", originHead, "origin", repoPath).
		Run(ctx, WithDir(tempDir)); err != nil {
		return errors.New("add remote failed")
	}
	fetchArgs := make([]string, 0)
	fetchArgs = append(fetchArgs, "--no-tags")
	if CheckGitVersionAtLeast("2.25.0") == nil {
		fetchArgs = append(fetchArgs, "--no-write-commit-graph")
	}
	if _, err := NewCommand("fetch", "origin", originHead+":"+MergeBranch, originHead+":original_"+originHead).AddArgs(fetchArgs...).
		Run(ctx, WithDir(tempDir)); err != nil {
		return err
	}
	if err := SetDefaultBranch(ctx, tempDir, MergeBranch); err != nil {
		return err
	}
	if _, err := NewCommand("fetch", "origin", targetCommitId+":"+TrackingBranch).AddArgs(fetchArgs...).
		Run(ctx, WithDir(tempDir)); err != nil {
		return err
	}
	return nil
}

func getDiffTreeForMerge(ctx context.Context, repoPath, target, head string) ([]string, error) {
	diffTreeResult, err := NewCommand("diff-tree", "--no-commit-id", "--name-only", "-r", "-r", "-z", "--root", target, head).
		Run(ctx, WithDir(repoPath))
	if err != nil {
		return nil, fmt.Errorf("unable to diff tree in tmpBasePath: %w", err)
	}
	treeResult := bytes.Split(diffTreeResult.ReadAsBytes(), []byte{'\x00'})
	ret := make([]string, 0)
	for _, r := range treeResult {
		line := strings.TrimSpace(string(r))
		if len(line) > 0 {
			ret = append(ret, fmt.Sprintf("/%s\n", escapedSymbols.ReplaceAllString(line, `\$1`)))
		}
	}
	return ret, nil
}

func MergeBase(ctx context.Context, repoPath, target, head string) (string, error) {
	result, err := NewCommand("merge-base", "--", target, head).Run(ctx, WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), nil
}

func MergeFile(ctx context.Context, repoPath string, files ...string) error {
	_, err := NewCommand("merge-file").AddArgs(files...).Run(ctx, WithDir(repoPath))
	return err
}
