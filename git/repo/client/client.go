package client

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/git/gitnode"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

var (
	HttpFailedErr = errors.New("http failed")
	httpClient    = httputil.NewRetryableHttpClient()
)

// InitRepo 初始化仓库
func InitRepo(ctx context.Context, req reqvo.InitRepoReq, nodeId string) error {
	var resp ginutil.BaseResp
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/initRepo",
		req,
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

// DeleteRepo 删除仓库
func DeleteRepo(ctx context.Context, req reqvo.DeleteRepoReq, nodeId string) error {
	var resp ginutil.BaseResp
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/delRepo",
		req,
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

// GetAllBranches 获取所有的分支
func GetAllBranches(ctx context.Context, req reqvo.GetAllBranchesReq, nodeId string) ([]string, error) {
	var resp ginutil.DataResp[[]string]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/getAllBranches",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

// GetAllTags 获取所有的tag
func GetAllTags(ctx context.Context, req reqvo.GetAllTagsReq, nodeId string) ([]string, error) {
	var resp ginutil.DataResp[[]string]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/getAllTags",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

// Gc 仓库gc
func Gc(ctx context.Context, req reqvo.GcReq, nodeId string) error {
	var resp ginutil.BaseResp
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/gc",
		req,
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

// DiffRefs 比较两个refs
func DiffRefs(ctx context.Context, req reqvo.DiffRefsReq, nodeId string) (reqvo.DiffRefsResp, error) {
	var resp ginutil.DataResp[reqvo.DiffRefsResp]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/diffRefs",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.DiffRefsResp{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.DiffRefsResp{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

// DiffFile 比较两个ref单个文件差异
func DiffFile(ctx context.Context, req reqvo.DiffFileReq, nodeId string) (reqvo.DiffFileResp, error) {
	var resp ginutil.DataResp[reqvo.DiffFileResp]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/diffFile",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.DiffFileResp{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.DiffFileResp{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

// GetRepoSize 获取仓库大小
func GetRepoSize(ctx context.Context, req reqvo.GetRepoSizeReq, nodeId string) (int64, error) {
	var resp ginutil.DataResp[int64]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/diffFile",
		req,
		&resp,
	)
	if err != nil {
		return 0, err
	}
	if !resp.IsSuccess() {
		return 0, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func ShowDiffTextContent(ctx context.Context, req reqvo.ShowDiffTextContentReq, nodeId string) ([]reqvo.DiffLineVO, error) {
	var resp ginutil.DataResp[[]reqvo.DiffLineVO]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/showDiffTextContent",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func HistoryCommits(ctx context.Context, req reqvo.HistoryCommitsReq, nodeId string) (reqvo.HistoryCommitsResp, error) {
	var resp ginutil.DataResp[reqvo.HistoryCommitsResp]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/historyCommits",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.HistoryCommitsResp{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.HistoryCommitsResp{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func InitRepoHook(ctx context.Context, req reqvo.InitRepoHookReq, nodeId string) error {
	var resp ginutil.BaseResp
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/initRepoHook",
		req,
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

func EntriesRepo(ctx context.Context, req reqvo.EntriesRepoReq, nodeId string) (reqvo.TreeVO, error) {
	var resp ginutil.DataResp[reqvo.TreeVO]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/entriesRepo",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.TreeVO{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.TreeVO{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func CatFile(ctx context.Context, req reqvo.CatFileReq, nodeId string) (reqvo.CatFileResp, error) {
	var resp ginutil.DataResp[reqvo.CatFileResp]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/catFile",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.CatFileResp{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.CatFileResp{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func TreeRepo(ctx context.Context, req reqvo.TreeRepoReq, nodeId string) (reqvo.TreeRepoResp, error) {
	var resp ginutil.DataResp[reqvo.TreeRepoResp]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/treeRepo",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.TreeRepoResp{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.TreeRepoResp{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func Merge(ctx context.Context, req reqvo.MergeReq, nodeId string) error {
	var resp ginutil.BaseResp
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/git/store/merge",
		req,
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

func UploadPack(req reqvo.UploadPackReq, nodeId string, repoId int64, pusherAccount, pusherEmail, appUrl string) error {
	return proxyHttp(
		nodeId,
		"/api/v1/git/smart/"+req.RepoPath+"/git-upload-pack",
		req.C,
		map[string]string{
			"Repo-RepoId":    strconv.FormatInt(repoId, 10),
			"Pusher-Account": pusherAccount,
			"Pusher-Email":   pusherEmail,
			"AppId-Url":      appUrl,
		},
	)
}

func ReceivePack(req reqvo.ReceivePackReq, nodeId string, repoId int64, pusherAccount, pusherEmail, appUrl string) error {
	return proxyHttp(
		nodeId,
		"/api/v1/git/smart/"+req.RepoPath+"/git-receive-pack",
		req.C,
		map[string]string{
			"Repo-RepoId":    strconv.FormatInt(repoId, 10),
			"Pusher-Account": pusherAccount,
			"Pusher-Email":   pusherEmail,
			"AppId-Url":      appUrl,
		},
	)
}

func InfoRefs(req reqvo.InfoRefsReq, nodeId string) error {
	return proxyHttp(
		nodeId,
		"/api/v1/git/smart/"+req.RepoPath+"/info/refs?service="+req.C.Query("service"),
		req.C,
		nil,
	)
}

func LfsUpload(req reqvo.LfsUploadReq, nodeId string) error {
	return proxyHttp(
		nodeId,
		"/api/v1/lfs/file/"+req.RepoPath+"/"+req.Oid+"/upload",
		req.C,
		nil,
	)
}

func LfsDownload(req reqvo.LfsDownloadReq, nodeId string) error {
	return proxyHttp(
		nodeId,
		"/api/v1/lfs/file/"+req.RepoPath+"/"+req.Oid+"/download",
		req.C,
		nil,
	)
}

func LfsExists(ctx context.Context, req reqvo.LfsExistsReq, nodeId string) (bool, error) {
	var resp ginutil.DataResp[bool]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/lfs/prop/exists",
		req,
		&resp,
	)
	if err != nil {
		return false, err
	}
	if !resp.IsSuccess() {
		return false, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func LfsBatchExists(ctx context.Context, req reqvo.LfsBatchExistsReq, nodeId string) (map[string]bool, error) {
	var resp ginutil.DataResp[map[string]bool]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/lfs/prop/batchExists",
		req,
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func LfsStat(ctx context.Context, req reqvo.LfsStatReq, nodeId string) (reqvo.LfsStatResp, error) {
	var resp ginutil.DataResp[reqvo.LfsStatResp]
	err := postHttp(
		ctx,
		nodeId,
		"/api/v1/lfs/prop/stat",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.LfsStatResp{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.LfsStatResp{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func postHttp(ctx context.Context, nodeId string, path string, req, resp any) error {
	httpUrl, err := getHttpUrl(ctx, nodeId)
	if err != nil {
		return err
	}
	err = httputil.Post(ctx,
		httpClient,
		httpUrl+path,
		map[string]string{
			"Authorization": getRepoToken(ctx),
		},
		req,
		resp,
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return HttpFailedErr
	}
	return nil
}

func proxyHttp(nodeId, path string, ctx *gin.Context, headers map[string]string) error {
	httpUrl, err := getHttpUrl(ctx, nodeId)
	if err != nil {
		return err
	}
	proxyReq, err := http.NewRequestWithContext(ctx, ctx.Request.Method, httpUrl+path, ctx.Request.Body)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return HttpFailedErr
	}
	proxyReq.Header = ctx.Request.Header.Clone()
	for k, v := range headers {
		proxyReq.Header.Set(k, v)
	}
	proxyReq.Header.Set("Authorization", getRepoToken(ctx))
	proxyResp, err := httpClient.Do(proxyReq)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return HttpFailedErr
	}
	defer proxyResp.Body.Close()
	for k := range proxyResp.Header {
		ctx.Header(k, proxyResp.Header.Get(k))
	}
	ctx.Writer.WriteHeader(proxyResp.StatusCode)
	io.Copy(ctx.Writer, proxyResp.Body)
	return nil
}

func getHttpUrl(ctx context.Context, nodeId string) (string, error) {
	ret, err := gitnode.PickHttpHost(ctx, nodeId)
	if err != nil {
		return "", err
	}
	return "http://" + ret, nil
}

func getRepoToken(ctx context.Context) string {
	cfg, b := cfgsrv.Inner.GetGitCfg(ctx)
	if !b {
		return ""
	}
	return cfg.RepoToken
}
