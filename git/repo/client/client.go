package client

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/rpcheader"
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
func InitRepo(ctx context.Context, req reqvo.InitRepoReq) (int64, error) {
	var resp ginutil.DataResp[int64]
	err := postHttp(
		ctx,
		"/api/v1/git/store/initRepo",
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

// DeleteRepo 删除仓库
func DeleteRepo(ctx context.Context, req reqvo.DeleteRepoReq) error {
	var resp ginutil.BaseResp
	err := postHttp(
		ctx,
		"/api/v1/git/store/deleteRepo",
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
func GetAllBranches(ctx context.Context, req reqvo.GetAllBranchesReq) ([]reqvo.RefVO, error) {
	var resp ginutil.DataResp[[]reqvo.RefVO]
	err := postHttp(
		ctx,
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

// DeleteBranch 删除分支
func DeleteBranch(ctx context.Context, req reqvo.DeleteBranchReq) error {
	var resp ginutil.DataResp[[]reqvo.RefVO]
	err := postHttp(
		ctx,
		"/api/v1/git/store/deleteBranch",
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

// PageBranchAndLastCommit 分页获取分支+最后提交信息
func PageBranchAndLastCommit(ctx context.Context, req reqvo.PageRefCommitsReq) ([]reqvo.RefCommitVO, int64, error) {
	var resp ginutil.Page2Resp[reqvo.RefCommitVO]
	err := postHttp(
		ctx,
		"/api/v1/git/store/pageBranchAndLastCommit",
		req,
		&resp,
	)
	if err != nil {
		return nil, 0, err
	}
	if !resp.IsSuccess() {
		return nil, 0, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, resp.TotalCount, nil
}

// PageTagAndCommit 分页获取tag+提交信息
func PageTagAndCommit(ctx context.Context, req reqvo.PageRefCommitsReq) ([]reqvo.RefCommitVO, int64, error) {
	var resp ginutil.Page2Resp[reqvo.RefCommitVO]
	err := postHttp(
		ctx,
		"/api/v1/git/store/pageTagAndCommit",
		req,
		&resp,
	)
	if err != nil {
		return nil, 0, err
	}
	if !resp.IsSuccess() {
		return nil, 0, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, resp.TotalCount, nil
}

// GetAllTags 获取所有的tag
func GetAllTags(ctx context.Context, req reqvo.GetAllTagsReq) ([]reqvo.RefVO, error) {
	var resp ginutil.DataResp[[]reqvo.RefVO]
	err := postHttp(
		ctx,
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
func Gc(ctx context.Context, req reqvo.GcReq) (reqvo.GcResp, error) {
	var resp ginutil.DataResp[reqvo.GcResp]
	err := postHttp(
		ctx,
		"/api/v1/git/store/gc",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.GcResp{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.GcResp{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

// DiffRefs 比较两个refs
func DiffRefs(ctx context.Context, req reqvo.DiffRefsReq) (reqvo.DiffRefsResp, error) {
	var resp ginutil.DataResp[reqvo.DiffRefsResp]
	err := postHttp(
		ctx,
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

// DiffCommits 比较两个commit
func DiffCommits(ctx context.Context, req reqvo.DiffCommitsReq) (reqvo.DiffCommitsResp, error) {
	var resp ginutil.DataResp[reqvo.DiffCommitsResp]
	err := postHttp(
		ctx,
		"/api/v1/git/store/diffCommits",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.DiffCommitsResp{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.DiffCommitsResp{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

// DiffFile 比较两个ref单个文件差异
func DiffFile(ctx context.Context, req reqvo.DiffFileReq) (reqvo.DiffFileResp, error) {
	var resp ginutil.DataResp[reqvo.DiffFileResp]
	err := postHttp(
		ctx,
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

func ShowDiffTextContent(ctx context.Context, req reqvo.ShowDiffTextContentReq) ([]reqvo.DiffLineVO, error) {
	var resp ginutil.DataResp[[]reqvo.DiffLineVO]
	err := postHttp(
		ctx,
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

func HistoryCommits(ctx context.Context, req reqvo.HistoryCommitsReq) (reqvo.HistoryCommitsResp, error) {
	var resp ginutil.DataResp[reqvo.HistoryCommitsResp]
	err := postHttp(
		ctx,
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

func InitRepoHook(ctx context.Context, req reqvo.InitRepoHookReq) error {
	var resp ginutil.BaseResp
	err := postHttp(
		ctx,
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

func DeleteTag(ctx context.Context, req reqvo.DeleteTagReqVO) error {
	var resp ginutil.BaseResp
	err := postHttp(
		ctx,
		"/api/v1/git/store/deleteTag",
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

func EntriesRepo(ctx context.Context, req reqvo.EntriesRepoReq) ([]reqvo.BlobVO, error) {
	var resp ginutil.DataResp[[]reqvo.BlobVO]
	err := postHttp(
		ctx,
		"/api/v1/git/store/entriesRepo",
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

func CatFile(ctx context.Context, req reqvo.CatFileReq) (reqvo.CatFileResp, error) {
	var resp ginutil.DataResp[reqvo.CatFileResp]
	err := postHttp(
		ctx,
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

func IndexRepo(ctx context.Context, req reqvo.IndexRepoReq) (reqvo.IndexRepoResp, error) {
	var resp ginutil.DataResp[reqvo.IndexRepoResp]
	err := postHttp(
		ctx,
		"/api/v1/git/store/indexRepo",
		req,
		&resp,
	)
	if err != nil {
		return reqvo.IndexRepoResp{}, err
	}
	if !resp.IsSuccess() {
		return reqvo.IndexRepoResp{}, bizerr.NewBizErr(resp.Code, resp.Message)
	}
	return resp.Data, nil
}

func Merge(ctx context.Context, req reqvo.MergeReq) (reqvo.DiffRefsResp, error) {
	var resp ginutil.DataResp[reqvo.DiffRefsResp]
	err := postHttp(
		ctx,
		"/api/v1/git/store/merge",
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

func UploadPack(req reqvo.UploadPackReq, repoId int64, pusherAccount, pusherEmail, appUrl string) error {
	return proxyHttp(
		"/api/v1/git/smart/"+req.RepoPath+"/git-upload-pack",
		req.C,
		map[string]string{
			"Repo-PrId":      strconv.FormatInt(repoId, 10),
			"Pusher-Account": pusherAccount,
			"Pusher-Email":   pusherEmail,
			"App-Url":        appUrl,
		},
	)
}

func ReceivePack(req reqvo.ReceivePackReq, repoId int64, pusherAccount, pusherEmail, appUrl string) error {
	return proxyHttp(
		"/api/v1/git/smart/"+req.RepoPath+"/git-receive-pack",
		req.C,
		map[string]string{
			"Repo-RepoId":    strconv.FormatInt(repoId, 10),
			"Pusher-Account": pusherAccount,
			"Pusher-Email":   pusherEmail,
			"AppId-HostUrl":  appUrl,
		},
	)
}

func InfoRefs(req reqvo.InfoRefsReq) error {
	return proxyHttp(
		"/api/v1/git/smart/"+req.RepoPath+"/info/refs?service="+req.Service,
		req.C,
		nil,
	)
}

func LfsUpload(req reqvo.LfsUploadReq) error {
	return proxyHttp(
		"/api/v1/lfs/file/"+req.RepoPath+"/"+req.Oid+"/upload",
		req.C,
		nil,
	)
}

func LfsDownload(req reqvo.LfsDownloadReq) error {
	return proxyHttp(
		"/api/v1/lfs/file/"+req.RepoPath+"/"+req.Oid+"/download",
		req.C,
		nil,
	)
}

func CreateArchive(req reqvo.CreateArchiveReq) error {
	return proxyHttp(
		"/api/v1/git/store/archive/"+req.RepoPath+"/"+req.FileName,
		req.C,
		nil,
	)
}

func LfsExists(ctx context.Context, req reqvo.LfsExistsReq) (bool, error) {
	var resp ginutil.DataResp[bool]
	err := postHttp(
		ctx,
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

func LfsBatchExists(ctx context.Context, req reqvo.LfsBatchExistsReq) (map[string]bool, error) {
	var resp ginutil.DataResp[map[string]bool]
	err := postHttp(
		ctx,
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

func LfsStat(ctx context.Context, req reqvo.LfsStatReq) (reqvo.LfsStatResp, error) {
	var resp ginutil.DataResp[reqvo.LfsStatResp]
	err := postHttp(
		ctx,
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

func Blame(ctx context.Context, req reqvo.BlameReq) ([]reqvo.BlameLineVO, error) {
	var resp ginutil.DataResp[[]reqvo.BlameLineVO]
	err := postHttp(
		ctx,
		"/api/v1/git/store/blame",
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

func CanMerge(ctx context.Context, req reqvo.CanMergeReq) (bool, error) {
	var resp ginutil.DataResp[bool]
	err := postHttp(
		ctx,
		"/api/v1/git/store/canMerge",
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

func postHttp(ctx context.Context, path string, req, resp any) error {
	httpUrl, err := getHttpUrl(ctx)
	if err != nil {
		return err
	}
	err = httputil.Post(ctx,
		httpClient,
		httpUrl+path,
		map[string]string{
			rpcheader.TraceId: rpcheader.GetHeaders(ctx).Get(rpcheader.TraceId),
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

func proxyHttp(path string, ctx *gin.Context, headers map[string]string) error {
	httpUrl, err := getHttpUrl(ctx)
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
	proxyReq.Header.Set(rpcheader.TraceId, rpcheader.GetHeaders(ctx).Get(rpcheader.TraceId))
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

func getHttpUrl(ctx context.Context) (string, error) {
	cfg, err := cfgsrv.GetGitRepoServerCfgFromDB(ctx)
	if err != nil {
		logger.Logger.Error(err)
		return "", bizerr.NewBizErr(apicode.OperationFailedErrCode.Int(), "git repo server url is not set")
	}
	return "http://" + cfg.HttpHost, nil
}
