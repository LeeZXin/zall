package sshproxy

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/git"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"path/filepath"
)

var (
	proxyDialer *zssh.Dialer
)

func InitProxy() zsf.LifeCycle {
	hostKey := filepath.Join(git.DataDir(), "ssh", "proxy.rsa")
	server, err := zssh.NewServer(&zssh.ServerOpts{
		Port:    static.GetInt("git.proxy.server.port"),
		HostKey: hostKey,
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			if ctx.User() != "git" {
				return false
			}
			return true
		},
		SessionHandler: func(session ssh.Session) {
			if err := handleGitCommand(session); err != nil {
				session.Exit(1)
			} else {
				session.Exit(0)
			}
		},
	})
	if err != nil {
		logger.Logger.Fatalf("new ssh server: %v", err)
	}
	proxyDialer, err = zssh.NewDialer(&zssh.DialerOpts{
		UserName: "git",
		HostKey:  hostKey,
	})
	if err != nil {
		logger.Logger.Fatalf("new ssh proxy dialer: %v", err)
	}
	return server
}

func handleGitCommand(session ssh.Session) error {
	fingerprint := gossh.FingerprintSHA256(session.PublicKey())
	// 寻找仓库存储节点
	sshHost, err := pickRepoSshHost(session.Context())
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	err = proxyDialer.ProxySession(sshHost, session, &zssh.ProxyOpts{
		SrcFingerprint: fingerprint,
	})
	if err != nil {
		return err
	}
	return nil
}

func pickRepoSshHost(ctx context.Context) (string, error) {
	cfg, b := cfgsrv.Inner.GetGitRepoServerCfg(ctx)
	if !b {
		return "", errors.New("git repo server ssh host is not set")
	}
	return cfg.SshHost, nil
}
