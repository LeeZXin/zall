package sshagent

import (
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	gossh "golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"sync"
)

var (
	pwdDir string
)

var (
	clientCfg  *gossh.ClientConfig
	clientOnce sync.Once
)

func initClientCfg() {
	clientOnce.Do(func() {
		pwd, err := os.Getwd()
		if err != nil {
			logger.Logger.Fatal(err)
		}
		hostKey, err := util.ReadOrGenRsaKey(filepath.Join(pwd, "data", "ssh", "sshAgent.rsa"))
		if err != nil {
			logger.Logger.Fatal(err)
		}
		privateKey, err := os.ReadFile(hostKey)
		if err != nil {
			logger.Logger.Fatal(err)
		}
		keySigner, err := gossh.ParsePrivateKey(privateKey)
		if err != nil {
			logger.Logger.Fatal(err)
		}
		clientCfg = zssh.NewCommonClientConfig("zall", keySigner)
	})
}
