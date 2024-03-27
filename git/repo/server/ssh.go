package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/sshkeymd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/git/gitenv"
	"github.com/LeeZXin/zall/pkg/git/lfs"
	"github.com/LeeZXin/zall/pkg/git/process"
	"github.com/LeeZXin/zall/pkg/i18n"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/gliderlabs/ssh"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kballard/go-shellquote"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	noneRepo = iota
	accessRepo
	updateRepo
)

const (
	lfsAuthenticateVerb = "git-lfs-authenticate"
)

var (
	allowedCommands = map[string]int{
		"git-upload-pack":    accessRepo,
		"git-upload-archive": accessRepo,
		"git-receive-pack":   updateRepo,
		lfsAuthenticateVerb:  noneRepo,
	}
)

func InitSshServer() zsf.LifeCycle {
	ret, err := zssh.NewServer(&zssh.ServerOpts{
		Port:    static.GetInt("git.repo.server.port"),
		HostKey: filepath.Join(git.DataDir(), "ssh", "gitRepoServer.rsa"),
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			if ctx.User() != "git" {
				return false
			}
			return true
		},
		SessionHandler: func(session ssh.Session) {
			ctx, cancel := context.WithCancel(session.Context())
			defer cancel()
			envs := util.CutEnv(session.Environ())
			// 检查proxy的fingerprint
			fingerprint := envs["ZGIT_SRC_FINGERPRINT"]
			user := getUserByFingerprint(ctx, fingerprint)
			if user == nil {
				// 找不到用户信息
				util.ExitWithErrMsg(session, "not found")
				return
			}
			if err := handleGitCommand(ctx, user, session); err != nil {
				util.ExitWithErrMsg(session, err.Error())
			} else {
				session.Exit(0)
			}
		},
	})
	if err != nil {
		logger.Logger.Fatalf("new ssh server: %v", err)
	}
	return ret
}

func handleGitCommand(ctx context.Context, user *usermd.UserInfo, session ssh.Session) error {
	gitCfg, b := cfgsrv.Inner.GetGitCfg(ctx)
	if !b {
		return errors.New(i18n.GetByKey(i18n.SystemInternalError))
	}
	words, err := shellquote.Split(session.RawCommand())
	if err != nil {
		return err
	}
	verb := words[0]
	repoPath := strings.TrimPrefix(words[1], "/")
	// 校验 repoPath
	sp := strings.Split(repoPath, "/")
	if len(sp) != 2 && strings.HasSuffix(sp[1], ".git") {
		return errors.New(i18n.GetByKey(i18n.SystemInvalidArgs))
	}
	var lfsVerb string
	if verb == lfsAuthenticateVerb {
		if len(words) > 2 {
			lfsVerb = words[2]
		}
	}
	accessMode, b := allowedCommands[verb]
	if !b {
		return errors.New(i18n.GetByKey(i18n.SshCmdNotSupported))
	}
	if verb == lfsAuthenticateVerb {
		if lfsVerb == "upload" {
			accessMode = updateRepo
		} else if lfsVerb == "download" {
			accessMode = accessRepo
		} else {
			return errors.New(i18n.GetByKey(i18n.SshCmdNotSupported))
		}
	}
	repo, err := checkAccessMode(ctx, user.Account, repoPath, accessMode)
	if err != nil {
		return err
	}
	// LFS token authentication
	if verb == lfsAuthenticateVerb {
		url := fmt.Sprintf("%s/%s/info/lfs", gitCfg.AppUrl, repoPath)
		now := time.Now()
		claims := lfs.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(gitCfg.GetLfsJwtExpiry())),
				NotBefore: jwt.NewNumericDate(now),
			},
			RepoPath: repoPath,
			Op:       lfsVerb,
			Account:  user.Account,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Sign and get the complete encoded token as a string using the secret
		tokenStr, err := token.SignedString(gitCfg.GetLfsJwtSecretBytes())
		if err != nil {
			return fmt.Errorf("failed to sign JWT token: %v", err)
		}
		tokenAuthentication := &lfs.TokenRespVO{
			Header: map[string]string{
				"Authorization": tokenStr,
			},
			Href: url,
		}
		err = json.NewEncoder(session).Encode(tokenAuthentication)
		if err != nil {
			return fmt.Errorf("failed to encode LFS json response: %v", err)
		}
		return nil
	}
	if accessMode == accessRepo {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    user.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.AccessRepo),
			ReqContent: repo,
			Err:        err,
		})
	} else {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    user.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.PushRepo),
			ReqContent: repo,
			Err:        err,
		})
	}
	var gitCmd *exec.Cmd
	gitBinPath := filepath.Dir(git.ExecutablePath()) // e.g. /usr/bin
	gitBinVerb := filepath.Join(gitBinPath, verb)    // e.g. /usr/bin/git-upload-pack
	if _, err = os.Stat(gitBinVerb); err != nil {
		verbFields := strings.SplitN(verb, "-", 2)
		if len(verbFields) == 2 {
			gitCmd = exec.CommandContext(ctx, git.ExecutablePath(), verbFields[1], repoPath)
		}
	}
	if gitCmd == nil {
		gitCmd = exec.CommandContext(ctx, gitBinVerb, repoPath)
	}
	process.SetSysProcAttribute(gitCmd)
	gitCmd.Dir = git.RepoDir()
	gitCmd.Stdout = session
	gitCmd.Stdin = session
	gitCmd.Stderr = session.Stderr()
	gitCmd.Env = append(gitCmd.Env, os.Environ()...)
	gitCmd.Env = append(gitCmd.Env,
		util.JoinFields(
			gitenv.EnvHookUrl, fmt.Sprintf("http://127.0.0.1:%d", common.HttpServerPort()),
			gitenv.EnvRepoId, strconv.FormatInt(repo.Id, 10),
			gitenv.EnvPusherAccount, user.Account,
			gitenv.EnvPusherEmail, user.Email,
			gitenv.EnvAppUrl, gitCfg.AppUrl,
			gitenv.EnvHookToken, git.HookToken(),
		)...,
	)
	gitCmd.Env = append(gitCmd.Env, git.CommonEnvs()...)
	return gitCmd.Run()
}

func checkAccessMode(ctx context.Context, account, repoPath string, permCode int) (repomd.RepoInfo, error) {
	repo, b := reposrv.Inner.GetByRepoPath(ctx, repoPath)
	if !b {
		return repomd.RepoInfo{}, util.InvalidArgsError()
	}
	// 获取权限
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, repo.TeamId, account)
	if !b {
		return repomd.RepoInfo{}, util.UnauthorizedError()
	}
	pass := false
	switch permCode {
	case accessRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanUpdateRepo
	case updateRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanAccessRepo
	case noneRepo:
		pass = true
	}
	if !pass {
		return repomd.RepoInfo{}, util.UnauthorizedError()
	}
	return repo, nil
}

func getUserByFingerprint(ctx context.Context, fingerprint string) *usermd.UserInfo {
	ctx, closer := xormstore.Context(ctx)
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
