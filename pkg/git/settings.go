package git

import (
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/logger"
	"os"
	"path/filepath"
)

var (
	dataDir, homeDir, lfsDir, tempDir, repoDir, actionDir string

	hookToken string
)

func initAllSettings() {
	initCommand()
	pwd, err := os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("os.Getwd err: %v", err)
	}
	dataDir = filepath.Join(pwd, "data")
	homeDir = filepath.Join(dataDir, "home")
	lfsDir = filepath.Join(dataDir, "lfs")
	tempDir = filepath.Join(dataDir, "temp")
	repoDir = filepath.Join(dataDir, "repo")
	actionDir = filepath.Join(dataDir, "action")
	util.MkdirAll(
		homeDir,
		lfsDir,
		tempDir,
		repoDir,
		actionDir,
	)
	hookToken = idutil.RandomUuid()
}

func HomeDir() string {
	return homeDir
}

func LfsDir() string {
	return lfsDir
}

func TempDir() string {
	return tempDir
}

func RepoDir() string {
	return repoDir
}

func SignUsername() string {
	return "zgit"
}

func SignEmail() string {
	return "zgit@fake.local"
}

func HookToken() string {
	return hookToken
}

func ActionDir() string {
	return actionDir
}

func DataDir() string {
	return dataDir
}
