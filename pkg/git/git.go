package git

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/localcache"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type FileMode string

const (
	RegularFileMode    FileMode = "100644"
	SymbolicLinkMode   FileMode = "120000"
	ExecutableFileMode FileMode = "100755"
	DirectoryMode      FileMode = "040000"
	SubModuleMode      FileMode = "160000"
	TreeFileMode       FileMode = "040755"
)

func (m FileMode) String() string {
	return string(m)
}

func (m FileMode) Readable() string {
	switch m {
	case RegularFileMode:
		return "regular"
	case SymbolicLinkMode:
		return "symlink"
	case ExecutableFileMode:
		return "executable"
	case DirectoryMode:
		return "directory"
	case SubModuleMode:
		return "subModule"
	case TreeFileMode:
		return "tree"
	default:
		return "unknown"
	}
}

const (
	RequiredVersion = "2.0.0"
	BranchPrefix    = "refs/heads/"
	TagPrefix       = "refs/tags/"
	TimeLayout      = "Mon Jan _2 15:04:05 2006 -0700"
	PrettyLogFormat = "--pretty=format:%H"

	ZeroCommitId = "0000000000000000000000000000000000000000"
)

const notRegularFileMode = os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice | os.ModeCharDevice | os.ModeIrregular

var (
	versionCache, _ = localcache.NewLazyLoader[*Version](getGitVersion)
)

func Init() {
	initAllSettings()
	if CheckGitVersionAtLeast(RequiredVersion) != nil {
		logger.Logger.Fatal("install git version is not supported, upgrade it before start github.com/LeeZXin/zall")
	}
	if CheckGitVersionAtLeast("2.18") == nil {
		addGlobalCmdArgs("-c", "protocol.version=2")
	}
	if CheckGitVersionAtLeast("2.9") == nil {
		addGlobalCmdArgs("-c", "credential.helper=")
	}
	if CheckGitVersionAtLeast("2.1.2") != nil {
		logger.Logger.Fatal("LFS server support requires Git >= 2.1.2")
	}
	addGlobalCmdArgs("-c", "filter.lfs.required=", "-c", "filter.lfs.smudge=", "-c", "filter.lfs.clean=")
	options := map[string]string{
		"diff.algorithm":  "histogram",
		"gc.reflogExpire": "90",

		"core.logAllRefUpdates": "true",
		"core.quotePath":        "false",
	}
	if CheckGitVersionAtLeast("2.10") == nil {
		options["receive.advertisePushOptions"] = "true"
	}
	if CheckGitVersionAtLeast("2.18") == nil {
		options["core.commitGraph"] = "true"
		options["gc.writeCommitGraph"] = "true"
		options["fetch.writeCommitGraph"] = "true"
	}
	if static.Exists("git.reflog.core.logAllRefUpdates") {
		options["core.logAllRefUpdates"] = strconv.FormatBool(static.GetBool("git.reflog.core.logAllRefUpdates"))
	}
	if static.GetInt("git.reflog.gc.reflogExpire") > 0 {
		options["gc.reflogExpire"] = static.GetString("git.reflog.gc.reflogExpire")
	}
	for k, v := range options {
		mustSetGlobalConfig(k, v)
	}
	mustSetGlobalConfigIfAbsent("user.name", SignUsername())
	mustSetGlobalConfigIfAbsent("user.email", SignEmail())
	mustAddGlobalConfigIfAbsent("safe.directory", "*")
	if util.IsWindows() {
		mustSetGlobalConfig("core.longpaths", "true")
		mustUnsetAllGlobalConfig("core.protectNTFS", "false")
	}
	if CheckGitVersionAtLeast("2.22") == nil {
		mustSetGlobalConfig("uploadpack.allowfilter", "true")
		mustSetGlobalConfig("uploadpack.allowAnySHA1InWant", "true")
	}
}

func mustSetGlobalConfig(k, v string) {
	if err := setGlobalConfig(k, v); err != nil {
		logger.Logger.Fatal(err)
	}
}

func mustSetGlobalConfigIfAbsent(k, v string) {
	if err := setGlobalConfigIfAbsent(k, v); err != nil {
		logger.Logger.Fatal(err)
	}
}

func mustAddGlobalConfigIfAbsent(k, v string) {
	if err := addGlobalConfigIfAbsent(k, v); err != nil {
		logger.Logger.Fatal(err)
	}
}

func mustUnsetAllGlobalConfig(k, v string) {
	if err := unsetAllGlobalConfig(k, v); err != nil {
		logger.Logger.Fatal(err)
	}
}

func setGlobalConfigIfAbsent(k, v string) error {
	return setGlobalConfigCheckOverwrite(k, v, false)
}

func setGlobalConfig(k, v string) error {
	return setGlobalConfigCheckOverwrite(k, v, true)
}

func setGlobalConfigCheckOverwrite(k, v string, overwrite bool) error {
	result, err := NewCommand("config", "--global", "--get", k).Run(nil)
	// fatal error
	if err != nil && !IsExitCode(err, 1) {
		return fmt.Errorf("failed to get git config %s, err: %w", k, err)
	}
	// 如果配置存在但不覆盖
	if err == nil && !overwrite {
		return nil
	}
	var currValue string
	// 配置存在
	if err == nil {
		currValue = strings.TrimSpace(result.ReadAsString())
	}
	if currValue == v {
		return nil
	}
	_, err = NewCommand("config", "--global").AddDynamicArgs(k, v).Run(nil)
	if err != nil {
		return fmt.Errorf("failed to set git global config %s, err: %w", k, err)
	}
	return nil
}

func addGlobalConfigIfAbsent(k, v string) error {
	_, err := NewCommand("config", "--global", "--get", k, regexp.QuoteMeta(v)).Run(nil)
	if err == nil {
		return nil
	}
	if IsExitCode(err, 1) {
		_, err = NewCommand("config", "--global", "--add").AddDynamicArgs(k, v).Run(nil)
		if err != nil {
			return fmt.Errorf("failed to add git global config %s, err: %w", k, err)
		}
		return nil
	}
	return fmt.Errorf("failed to get git config %s, err: %w", k, err)
}

func setLocalConfig(ctx context.Context, repoPath, k, v string) error {
	_, err := NewCommand("config", "--local", k, v).Run(ctx, WithDir(repoPath))
	if err != nil {
		return fmt.Errorf("failed to set local config %s, err: %w", k, err)
	}
	return nil
}

func unsetAllGlobalConfig(k, v string) error {
	_, err := NewCommand("config", "--global", "--get").AddDynamicArgs(k).Run(nil)
	if err == nil {
		_, err = NewCommand("config", "--global", "--unset-all").AddDynamicArgs(k, regexp.QuoteMeta(v)).Run(nil)
		if err != nil {
			return fmt.Errorf("failed to unset git global config %s, err: %w", k, err)
		}
		return nil
	}
	if IsExitCode(err, 1) {
		return nil
	}
	return fmt.Errorf("failed to get git config %s, err: %w", k, err)
}

// IsReferenceExist returns true if given reference exists in the repository.
func IsReferenceExist(ctx context.Context, repoPath, name string) bool {
	_, err := NewCommand("show-ref", "--verify", "--").AddDynamicArgs(name).Run(ctx, WithDir(repoPath))
	return err == nil
}

// IsBranchExist returns true if given branch exists in the repository.
func IsBranchExist(ctx context.Context, repoPath, name string) bool {
	if !strings.HasPrefix(name, BranchPrefix) {
		name = BranchPrefix + name
	}
	return IsReferenceExist(ctx, repoPath, name)
}

func HashObjectByStdin(ctx context.Context, repoPath string, reader io.Reader) (string, error) {
	result, err := NewCommand("hash-object", "-w", "--stdin").
		Run(ctx, WithDir(repoPath), WithStdin(reader))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), nil
}

func HashObjectByPath(ctx context.Context, repoPath, relativePath, absolutePath string) (string, error) {
	result, err := NewCommand("hash-object", "-w", "--path").AddDynamicArgs(relativePath, absolutePath).Run(ctx, WithDir(repoPath))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), nil
}

func AddObjectToIndex(ctx context.Context, repoPath, mode, object, filename string) error {
	_, err := NewCommand("update-index", "--add", "--replace", "--cacheinfo").AddDynamicArgs(mode, object, filename).
		Run(ctx, WithDir(repoPath))
	return err
}

// RemoveFilesFromIndex removes given filenames from the index - it does not check whether they are present.
func RemoveFilesFromIndex(ctx context.Context, repoPath string, filenames ...string) error {
	buffer := new(bytes.Buffer)
	for _, file := range filenames {
		if file != "" {
			buffer.WriteString("0 0000000000000000000000000000000000000000\t")
			buffer.WriteString(file)
			buffer.WriteByte('\000')
		}
	}
	_, err := NewCommand("update-index", "--remove", "-z", "--index-info").
		Run(ctx, WithDir(repoPath), WithStdin(bytes.NewReader(buffer.Bytes())))
	return err
}

// WriteTree writes the current index as a tree to the object db and returns its hash
func WriteTree(ctx context.Context, repoPath string) (Tree, error) {
	result, err := NewCommand("write-tree").Run(ctx, WithDir(repoPath))
	if err != nil {
		return Tree{}, err
	}
	return NewTree(strings.TrimSpace(result.ReadAsString())), nil
}

type CommitTreeOpts struct {
	Parents []string
	Message string
}

// CommitTree creates a commit from a given tree id for the user with provided message
func CommitTree(ctx context.Context, repoPath string, tree Tree, opts CommitTreeOpts) (string, error) {
	cmd := NewCommand("commit-tree", "--no-gpg-sign").AddDynamicArgs(tree.Id)
	for _, parent := range opts.Parents {
		cmd.AddArgs("-p", parent)
	}
	message := new(bytes.Buffer)
	message.WriteString(opts.Message)
	message.WriteString("\n")
	result, err := cmd.Run(ctx, WithDir(repoPath), WithStdin(message))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.ReadAsString()), nil
}

func GetRepoSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, info os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) { // ignore the error because the file maybe deleted during traversing.
				return nil
			}
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := info.Info()
		if err != nil {
			return err
		}
		if (f.Mode() & notRegularFileMode) == 0 {
			size += f.Size()
		}
		return err
	})
	return size, err
}

func GetLfsSize(path string) (int64, error) {
	return GetRepoSize(path)
}

type RefName string

func (n RefName) IsBranch() bool {
	return strings.HasPrefix(string(n), BranchPrefix)
}

func (n RefName) IsTag() bool {
	return strings.HasPrefix(string(n), TagPrefix)
}

func Gc(ctx context.Context, repoPath string) error {
	_, err := NewCommand("gc").Run(ctx, WithDir(repoPath))
	return err
}

func JoinRelativeRepoPath(str ...string) string {
	return filepath.Join(str...) + ".git"
}

func JoinAbsRepoPath(str ...string) string {
	return filepath.Join(RepoDir(), JoinRelativeRepoPath(str...))
}
