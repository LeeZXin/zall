package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/gitnodemd"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
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
func InitRepo(ctx context.Context, req reqvo.InitRepoReq, nodeId int64) error {
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
func DeleteRepo(ctx context.Context, req reqvo.DeleteRepoReq, nodeId int64) error {
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
func GetAllBranches(ctx context.Context, req reqvo.GetAllBranchesReq, nodeId int64) ([]string, error) {
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
func GetAllTags(ctx context.Context, req reqvo.GetAllTagsReq, nodeId int64) ([]string, error) {
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
func Gc(ctx context.Context, req reqvo.GcReq, nodeId int64) error {
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
func DiffRefs(ctx context.Context, req reqvo.DiffRefsReq, nodeId int64) (reqvo.DiffRefsResp, error) {
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
func DiffFile(ctx context.Context, req reqvo.DiffFileReq, nodeId int64) (reqvo.DiffFileResp, error) {
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
func GetRepoSize(ctx context.Context, req reqvo.GetRepoSizeReq, nodeId int64) (int64, error) {
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

func ShowDiffTextContent(ctx context.Context, req reqvo.ShowDiffTextContentReq, nodeId int64) ([]reqvo.DiffLineVO, error) {
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

func HistoryCommits(ctx context.Context, req reqvo.HistoryCommitsReq, nodeId int64) (reqvo.HistoryCommitsResp, error) {
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

func InitRepoHook(ctx context.Context, req reqvo.InitRepoHookReq, nodeId int64) error {
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

func EntriesRepo(ctx context.Context, req reqvo.EntriesRepoReq, nodeId int64) (reqvo.TreeVO, error) {
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

func CatFile(ctx context.Context, req reqvo.CatFileReq, nodeId int64) (reqvo.CatFileResp, error) {
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

func TreeRepo(ctx context.Context, req reqvo.TreeRepoReq, nodeId int64) (reqvo.TreeRepoResp, error) {
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

func Merge(ctx context.Context, req reqvo.MergeReq, nodeId int64) error {
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

func UploadPack(req reqvo.UploadPackReq, nodeId int64, repoId int64, pusherAccount, pusherEmail, appUrl string) error {
	return proxyHttp(
		nodeId,
		"/api/v1/git/smart/"+req.RepoPath+"/git-upload-pack",
		req.C,
		map[string]string{
			"Repo-TeamId":    strconv.FormatInt(repoId, 10),
			"Pusher-Account": pusherAccount,
			"Pusher-Email":   pusherEmail,
			"AppId-HostUrl":  appUrl,
		},
	)
}

func ReceivePack(req reqvo.ReceivePackReq, nodeId int64, repoId int64, pusherAccount, pusherEmail, appUrl string) error {
	return proxyHttp(
		nodeId,
		"/api/v1/git/smart/"+req.RepoPath+"/git-receive-pack",
		req.C,
		map[string]string{
			"Repo-TeamId":    strconv.FormatInt(repoId, 10),
			"Pusher-Account": pusherAccount,
			"Pusher-Email":   pusherEmail,
			"AppId-HostUrl":  appUrl,
		},
	)
}

func InfoRefs(req reqvo.InfoRefsReq, nodeId int64) error {
	return proxyHttp(
		nodeId,
		"/api/v1/git/smart/"+req.RepoPath+"/info/refs?service="+req.C.Query("service"),
		req.C,
		nil,
	)
}

func LfsUpload(req reqvo.LfsUploadReq, nodeId int64) error {
	return proxyHttp(
		nodeId,
		"/api/v1/lfs/file/"+req.RepoPath+"/"+req.Oid+"/upload",
		req.C,
		nil,
	)
}

func LfsDownload(req reqvo.LfsDownloadReq, nodeId int64) error {
	return proxyHttp(
		nodeId,
		"/api/v1/lfs/file/"+req.RepoPath+"/"+req.Oid+"/download",
		req.C,
		nil,
	)
}

func LfsExists(ctx context.Context, req reqvo.LfsExistsReq, nodeId int64) (bool, error) {
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

func LfsBatchExists(ctx context.Context, req reqvo.LfsBatchExistsReq, nodeId int64) (map[string]bool, error) {
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

func LfsStat(ctx context.Context, req reqvo.LfsStatReq, nodeId int64) (reqvo.LfsStatResp, error) {
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

func postHttp(ctx context.Context, nodeId int64, path string, req, resp any) error {
	httpUrl, err := getHttpUrl(nodeId)
	if err != nil {
		return err
	}
	err = httputil.Post(ctx,
		httpClient,
		httpUrl+path,
		nil,
		req,
		resp,
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return HttpFailedErr
	}
	return nil
}

func proxyHttp(nodeId int64, path string, ctx *gin.Context, headers map[string]string) error {
	httpUrl, err := getHttpUrl(nodeId)
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

func getHttpUrl(nodeId int64) (string, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	node, b, err := gitnodemd.GetById(ctx, nodeId)
	if err != nil {
		return "", err
	}
	if !b {
		return "", fmt.Errorf("unknown git node id: %d", nodeId)
	}
	return "http://" + node.HttpHost, nil
}
