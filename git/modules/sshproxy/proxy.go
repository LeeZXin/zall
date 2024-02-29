package sshproxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/git/gitnode"
	"github.com/LeeZXin/zall/git/modules/model/sshkeymd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/i18n"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/gliderlabs/ssh"
	"github.com/kballard/go-shellquote"
	gossh "golang.org/x/crypto/ssh"
	"path/filepath"
	"strings"
)

var (
	fingerprintKey = zssh.ContextKey("ZGIT_FINGERPRINT")

	proxyDialer *zssh.Dialer
	hiWords     = "Hi there! You've successfully authenticated with the deploy key named %v, but zgit does not provide shell access."
)

func InitProxy() {
	hostKey := filepath.Join(git.DataDir(), "ssh", "proxy.rsa")
	server, err := zssh.NewServer(&zssh.ServerOpts{
		Port:    static.GetInt("git.proxy.server.port"),
		HostKey: filepath.Join(git.DataDir(), "ssh", "proxy.rsa"),
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			if ctx.User() != "git" {
				return false
			}
			ctx.SetValue(fingerprintKey, gossh.FingerprintSHA256(key))
			return true
		},
		SessionHandler: func(session ssh.Session) {
			ctx, cancel := context.WithCancel(session.Context())
			defer cancel()
			if err := handleGitCommand(ctx, session); err != nil {
				util.ExitWithErrMsg(session, "internal error")
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
	zsf.RegisterApplicationLifeCycle(server)
}

func handleGitCommand(ctx context.Context, session ssh.Session) error {
	fingerprint := gossh.FingerprintSHA256(session.PublicKey())
	cmd := session.RawCommand()
	// 命令为空
	if cmd == "" {
		user := getUserByFingerprint(ctx, fingerprint)
		if user == nil {
			return errors.New("user not found")
		}
		fmt.Fprintln(session, fmt.Sprintf(hiWords, user.Name))
		return nil
	}
	words, err := shellquote.Split(cmd)
	if err != nil {
		return errors.New("error parsing arguments")
	}
	if len(words) < 2 {
		if git.CheckGitVersionAtLeast("2.29") == nil {
			if cmd == "ssh_info" {
				fmt.Fprintln(session, `{"type":"zgit","version":1}`)
				return nil
			}
		}
		return errors.New(i18n.GetByKey(i18n.SystemInvalidArgs))
	}
	repoPath := strings.TrimPrefix(words[1], "/")
	// 寻找仓库存储节点
	sshHost := pickRepoSshHost(ctx, repoPath)
	if sshHost == "" {
		return fmt.Errorf("ssh host not found: %s", repoPath)
	}
	err = proxyDialer.ProxySession(sshHost, session, &zssh.ProxyOpts{
		SrcFingerprint: ctx.Value(fingerprintKey).(string),
	})
	if err != nil {
		return err
	}
	return nil
}

func getUserByFingerprint(ctx context.Context, fingerprint string) *usermd.UserInfo {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	account, b, err := sshkeymd.GetAccountByFingerprint(ctx, fingerprint)
	if err != nil {
		logger.Logger.Error(err)
		return nil
	}
	if !b {
		return nil
	}
	user, b := usersrv.Inner.GetByAccount(ctx, account)
	// 账号不存在或被禁用
	if !b || user.IsProhibited {
		return nil
	}
	return &user
}

func pickRepoSshHost(ctx context.Context, repoPath string) string {
	repo, b := reposrv.Inner.GetByRepoPath(ctx, repoPath)
	if !b {
		return ""
	}
	host, err := gitnode.PickSshHost(ctx, repo.NodeId)
	if err != nil {
		logger.Logger.Error(err)
	}
	return host
}
