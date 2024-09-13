package git

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/pkg/git/gitenv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/common"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	WikiDefaultBranch = "main"
	GitDefaultBranch  = "main"
)

var (
	gitIgnoreResourcesPath = filepath.Join(common.ResourcesDir, "gitignore")

	RepoExistsErr = errors.New("repo exists")
)

type Repository struct {
	Id    string `json:"LockId"`
	Owner User   `json:"owner"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}

type InitRepoOpts struct {
	Owner         User
	RepoName      string
	RepoPath      string
	AddReadme     bool
	GitIgnoreName string
	DefaultBranch string
}

type InitWikiOpts struct {
	Owner    User
	WikiPath string
}

type CommitAndPushOpts struct {
	RepoPath          string
	Owner, Committer  User
	Branch, CommitMsg string
}

func initEmptyRepository(ctx context.Context, repoPath string, bare bool) error {
	if repoPath == "" {
		return errors.New("repoPath is empty")
	}
	isExist, err := util.IsExist(repoPath)
	if err != nil {
		return err
	}
	if isExist {
		return RepoExistsErr
	}
	err = os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		return err
	}
	cmd := NewCommand("init")
	if bare {
		cmd.AddArgs(BareFlag)
	}
	_, err = cmd.Run(ctx, WithDir(repoPath))
	return err
}

func InitRepository(ctx context.Context, opts InitRepoOpts) error {
	if err := initEmptyRepository(ctx, opts.RepoPath, true); err != nil {
		return err
	}
	if opts.DefaultBranch == "" {
		opts.DefaultBranch = GitDefaultBranch
	}
	if opts.AddReadme || opts.GitIgnoreName != "" {
		tmpDir, err := os.MkdirTemp(TempDir(), "init-"+util.RandomIdWithTime())
		if err != nil {
			return fmt.Errorf("failed to create temp dir for repository %s: %w", opts.RepoPath, err)
		}
		defer util.RemoveAll(tmpDir)
		if err = initTemporaryRepository(ctx, tmpDir, opts); err != nil {
			return err
		}
		SetDefaultBranch(ctx, opts.RepoPath, opts.DefaultBranch)
	}
	return InitRepoHook(opts.RepoPath)
}

func initTemporaryRepository(ctx context.Context, tmpDir string, opts InitRepoOpts) error {
	if _, err := NewCommand("clone").AddDynamicArgs(opts.RepoPath, tmpDir).Run(ctx); err != nil {
		return fmt.Errorf("failed to clone original repository %s: %w", opts.RepoPath, err)
	}
	if opts.AddReadme {
		util.WriteFile(filepath.Join(tmpDir, "README.md"), []byte(fmt.Sprintf("# %s  \n", opts.RepoName)))
	}
	if opts.GitIgnoreName != "" {
		content, err := os.ReadFile(filepath.Join(gitIgnoreResourcesPath, opts.GitIgnoreName))
		if err == nil {
			util.WriteFile(filepath.Join(tmpDir, ".gitignore"), content)
		}
	}
	return commitAndPushRepository(ctx, CommitAndPushOpts{
		RepoPath:  tmpDir,
		Owner:     opts.Owner,
		Committer: opts.Owner,
		Branch:    opts.DefaultBranch,
		CommitMsg: "first commit",
	})
}

func EnsureValidRepository(ctx context.Context, repoPath string) error {
	cmd := NewCommand("rev-parse")
	_, err := cmd.Run(ctx, WithDir(repoPath))
	return err
}

func SetDefaultBranch(ctx context.Context, repoPath, branch string) error {
	if !strings.HasPrefix(branch, BranchPrefix) {
		branch = BranchPrefix + branch
	}
	cmd := NewCommand("symbolic-ref", "HEAD").AddDynamicArgs(branch)
	_, err := cmd.Run(ctx, WithDir(repoPath))
	return err
}

func commitAndPushRepository(ctx context.Context, opts CommitAndPushOpts) error {
	commitTimeStr := time.Now().Format(time.RFC3339)
	env := append(
		os.Environ(),
		util.JoinFields(
			"GIT_AUTHOR_NAME", opts.Owner.Account,
			"GIT_AUTHOR_EMAIL", opts.Owner.Email,
			"GIT_AUTHOR_DATE", commitTimeStr,
			"GIT_COMMITTER_DATE", commitTimeStr,
		)...,
	)
	_, err := NewCommand("add", "--all").Run(ctx, WithDir(opts.RepoPath))
	if err != nil {
		return fmt.Errorf("git add -all failed repo:%s err: %v", opts.RepoPath, err)
	}
	commitCmd := NewCommand(
		"commit",
		"--no-gpg-sign",
		fmt.Sprintf("--message=%s", opts.CommitMsg),
		fmt.Sprintf("--author='%s <%s>'", opts.Committer.Account, opts.Committer.Email),
	)
	env = append(env,
		util.JoinFields(
			"GIT_COMMITTER_NAME", opts.Committer.Account,
			"GIT_COMMITTER_EMAIL", opts.Committer.Email,
		)...,
	)
	_, err = commitCmd.Run(ctx, WithDir(opts.RepoPath), WithEnv(env))
	if err != nil {
		return fmt.Errorf("git commit failed repo:%s err: %v", opts.RepoPath, err)
	}
	_, err = NewCommand("push", "origin").AddDynamicArgs("HEAD:"+opts.Branch).
		Run(
			ctx,
			WithDir(opts.RepoPath),
			WithEnv(util.JoinFields(gitenv.EnvIsInternal, "true")),
		)
	return err
}

func GetRepoUsername(repoPath string) (string, error) {
	return getRepoProperty("user.name", repoPath)
}

func GetRepoSignKey(repoPath string) (string, error) {
	return getRepoProperty("user.signingkey", repoPath)
}

func GetRepoUserEmail(repoPath string) (string, error) {
	return getRepoProperty("user.email", repoPath)
}

func getRepoProperty(name, repoPath string) (string, error) {
	run, err := NewCommand("config", "--get").AddDynamicArgs(name).Run(nil, WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return run.ReadAsString(), nil
}
