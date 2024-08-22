package cmd

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/git/gitenv"
	"github.com/LeeZXin/zall/pkg/githook"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// subHookPreReceive 可用于仓库大小检查提交权限和分支
var subHookPreReceive = &cli.Command{
	Name:        "pre-receive",
	Usage:       "pre-receive hook",
	Description: "This command should only be called by zgit",
	Action:      runPreReceive,
}

// subHookPostReceive 用于发送通知等
var subHookPostReceive = &cli.Command{
	Name:        "post-receive",
	Usage:       "post-receive hook",
	Description: "This command should only be called by zgit",
	Action:      runPostReceive,
}

var GitHook = &cli.Command{
	Name:        "gitHook",
	Usage:       "This command for zgit hook",
	Description: "This command should only be called by zgit",
	Subcommands: []*cli.Command{
		subHookPreReceive,
		subHookPostReceive,
	},
}

func runPreReceive(c *cli.Context) error {
	if isInternal, _ := strconv.ParseBool(os.Getenv(gitenv.EnvIsInternal)); isInternal {
		return nil
	}
	ctx, cancel := initWaitContext(c.Context)
	defer cancel()
	err := scanStdinAndDoHttp(ctx, githook.ApiPreReceiveUrl)
	if err != nil {
		if bizerr.IsBizErr(err) {
			return err
		}
		return errors.New("internal error")
	}
	return nil
}

// scanStdinAndDoHttp 处理输入并发送http
func scanStdinAndDoHttp(ctx context.Context, httpUrl string) error {
	infoList := make([]githook.RevInfo, 0)
	// the environment is set by serv command
	pusherId := os.Getenv(gitenv.EnvPusherAccount)
	pusherName := os.Getenv(gitenv.EnvPusherName)
	pusherEmail := os.Getenv(gitenv.EnvPusherEmail)
	prId, _ := strconv.ParseInt(os.Getenv(gitenv.EnvPrId), 10, 64)
	repoId, _ := strconv.ParseInt(os.Getenv(gitenv.EnvRepoId), 10, 64)
	aod := os.Getenv(gitenv.EnvAlternativeObjectDirectories)
	qp := os.Getenv(gitenv.EnvQuarantinePath)
	od := os.Getenv(gitenv.EnvObjectDirectory)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Bytes()))
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		refName := git.RefName(fields[2])
		if refName.IsBranch() || refName.IsTag() {
			infoList = append(infoList, githook.RevInfo{
				OldCommitId: fields[0],
				NewCommitId: fields[1],
				Ref:         fields[2],
			})
		}
	}
	client := newHttpClient()
	defer client.CloseIdleConnections()
	partitionList := listutil.Partition(infoList, 30)
	for _, partition := range partitionList {
		reqVO := githook.Opts{
			RevInfoList:                  partition,
			RepoId:                       repoId,
			PrId:                         prId,
			PusherAccount:                pusherId,
			PusherName:                   pusherName,
			PusherEmail:                  pusherEmail,
			ObjectDirectory:              od,
			AlternativeObjectDirectories: aod,
			QuarantinePath:               qp,
		}
		if err := doHttp(ctx, client, reqVO, httpUrl); err != nil {
			return err
		}
	}
	return nil
}

func runPostReceive(c *cli.Context) error {
	if isInternal, _ := strconv.ParseBool(os.Getenv(gitenv.EnvIsInternal)); isInternal {
		return nil
	}
	ctx, cancel := initWaitContext(c.Context)
	defer cancel()
	err := scanStdinAndDoHttp(ctx, githook.ApiPostReceiveUrl)
	if err != nil {
		if bizerr.IsBizErr(err) {
			return err
		}
		return errors.New("internal error")
	}
	return nil
}

func doHttp(ctx context.Context, client *http.Client, reqVO githook.Opts, url string) error {
	resp := ginutil.BaseResp{}
	err := httputil.Post(
		ctx,
		client,
		fmt.Sprintf("%s/%s", os.Getenv(gitenv.EnvHookUrl), url),
		map[string]string{
			"Authorization": os.Getenv(gitenv.EnvHookToken),
		},
		reqVO,
		&resp,
	)
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return nil
}

func newHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        0, // 禁用连接池
			MaxIdleConnsPerHost: 0, // 禁用连接池
			IdleConnTimeout:     0, // 禁用连接池
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 5 * time.Second,
	}
}
