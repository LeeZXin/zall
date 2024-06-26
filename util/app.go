package util

import (
	"errors"
	"github.com/LeeZXin/zsf/logger"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	isWindows bool
	appPath   string
)

func init() {
	isWindows = runtime.GOOS == "windows"
	appPath = getAppPath()
}

func getAppPath() string {
	var path string
	var err error
	defer func() {
		if err != nil {
			logger.Logger.Fatalf("getAppPath failed: %v", err)
		}
	}()
	if runtime.GOOS == "windows" && filepath.IsAbs(os.Args[0]) {
		path = filepath.Clean(os.Args[0])
	} else {
		path, err = exec.LookPath(os.Args[0])
	}
	if err != nil {
		if !errors.Is(err, exec.ErrDot) {
			return ""
		}
		path, err = filepath.Abs(os.Args[0])
	}
	if err != nil {
		return ""
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(path, "\\", "/")
}

func GetAppPath() string {
	return appPath
}
