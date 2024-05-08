package git

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"path/filepath"
	"strings"
)

func InitWiki(ctx context.Context, opts InitWikiOpts) error {
	if err := initEmptyRepository(ctx, opts.WikiPath, true); err != nil {
		return err
	}
	return SetDefaultBranch(ctx, opts.WikiPath, WikiDefaultBranch)
}

func UpdateWikiPage(ctx context.Context, wikiPath, pageName, content, message string) error {
	tempDir := filepath.Join(TempDir(), "wiki-"+idutil.RandomUuid())
	defer util.RemoveAll(tempDir)
	hasMasterBranch, err := prepareUpdateWikiPage(ctx, wikiPath, tempDir)
	if err != nil {
		return err
	}
	object, err := HashObjectByStdin(ctx, tempDir, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("hash content failed with err:%v", err)
	}
	if err = AddObjectToIndex(ctx, tempDir, RegularFileMode.String(), object, pageName); err != nil {
		return fmt.Errorf("addObjectToIndex failed with err:%v", err)
	}
	return afterUpdateWikiPage(ctx, tempDir, message, hasMasterBranch)
}

func DeleteWikiPage(ctx context.Context, wikiPath, pageName, message string) error {
	tempDir := filepath.Join(TempDir(), "wiki-"+idutil.RandomUuid())
	defer util.RemoveAll(tempDir)
	hasMasterBranch, err := prepareUpdateWikiPage(ctx, wikiPath, tempDir)
	if err != nil {
		return err
	}
	if err = RemoveFilesFromIndex(ctx, tempDir, pageName); err != nil {
		return fmt.Errorf("RemoveFilesFromIndex failed with err:%v", err)
	}
	return afterUpdateWikiPage(ctx, tempDir, message, hasMasterBranch)
}

func afterUpdateWikiPage(ctx context.Context, tempDir string, message string, hasMasterBranch bool) error {
	tree, err := WriteTree(ctx, tempDir)
	if err != nil {
		return fmt.Errorf("write tree failed with err:%v", err)
	}
	opts := CommitTreeOpts{
		Message: message,
	}
	if hasMasterBranch {
		opts.Parents = []string{"HEAD"}
	}
	commitHash, err := CommitTree(ctx, tempDir, tree, opts)
	if err != nil {
		return fmt.Errorf("commit tree failed with err:%v", err)
	}
	if _, err = NewCommand("push", DefaultRemote).AddDynamicArgs(fmt.Sprintf("%s:%s", commitHash, BranchPrefix+WikiDefaultBranch)).
		Run(ctx, WithDir(tempDir)); err != nil {
		return fmt.Errorf("push failed with err:%v", err)
	}
	return nil
}

func prepareUpdateWikiPage(ctx context.Context, wikiPath string, tempDir string) (bool, error) {
	hasMasterBranch := IsBranchExist(ctx, wikiPath, WikiDefaultBranch)
	cloneCmd := NewCommand("clone", "-s", "--bare").AddDynamicArgs(wikiPath, tempDir)
	if hasMasterBranch {
		cloneCmd.AddArgs("-b", WikiDefaultBranch)
	}
	if _, err := cloneCmd.Run(ctx); err != nil {
		return false, fmt.Errorf("clone tempDir:%s failed with err:%v", tempDir, err)
	}
	if hasMasterBranch {
		commitId, err := GetBranchCommitId(ctx, tempDir, "HEAD")
		if err != nil {
			return false, fmt.Errorf("get head commitId failed with err:%v", err)
		}
		if _, err = NewCommand("read-tree").AddDynamicArgs(commitId).
			Run(ctx, WithDir(tempDir)); err != nil {
			return false, fmt.Errorf("read tree failed with err:%v", err)
		}
	}
	return hasMasterBranch, nil
}
