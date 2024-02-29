package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/pkg/git/gitenv"
	"github.com/LeeZXin/zall/util"
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
}

// IsMergeAble 是否可合并
func (i *DiffRefsInfo) IsMergeAble() bool {
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
	pr.Commits, err = GetGitDiffCommitList(ctx, repoPath, pr.HeadCommit.Id, pr.TargetCommit.Id)
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
	pr.ConflictFiles, err = findConflictFiles(ctx, repoPath, pr)
	if err != nil {
		return DiffRefsInfo{}, err
	}
	return pr, nil
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
	if err := prepare4Merge(ctx, repoPath, tempDir, pr); err != nil {
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
				),
			),
		); err != nil {
		return fmt.Errorf("git push: %v", err)
	}
	return nil
}

func prepare4Merge(ctx context.Context, repoPath string, tempDir string, pr DiffRefsInfo) error {
	if err := initEmptyRepository(ctx, tempDir, false); err != nil {
		return err
	}
	if _, err := NewCommand("remote", "add", "-t", pr.OriginHead, "-m", pr.OriginHead, "origin", repoPath).
		Run(ctx, WithDir(tempDir)); err != nil {
		return errors.New("add remote failed")
	}
	fetchArgs := make([]string, 0)
	fetchArgs = append(fetchArgs, "--no-tags")
	if CheckGitVersionAtLeast("2.25.0") == nil {
		fetchArgs = append(fetchArgs, "--no-write-commit-graph")
	}
	if _, err := NewCommand("fetch", "origin", pr.OriginHead+":"+MergeBranch, pr.OriginHead+":original_"+pr.OriginHead).AddArgs(fetchArgs...).
		Run(ctx, WithDir(tempDir)); err != nil {
		return err
	}
	if err := SetDefaultBranch(ctx, tempDir, MergeBranch); err != nil {
		return err
	}
	if _, err := NewCommand("fetch", "origin", pr.TargetCommit.Id+":"+TrackingBranch).AddArgs(fetchArgs...).
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
