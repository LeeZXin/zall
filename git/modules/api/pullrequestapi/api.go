package pullrequestapi

import (
	"github.com/LeeZXin/zall/git/modules/service/pullrequestsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/pullRequest", apisession.CheckLogin)
		{
			// 获取合并请求
			group.GET("/get/:prId", getPullRequest)
			// 合并请求详情
			group.GET("/stats/:repoId", statsPullRequest)
			// 合并请求列表
			group.GET("/list", listPullRequest)
			// 创建合并请求
			group.POST("/submit", submitPullRequest)
			// 关闭合并请求
			group.PUT("/close/:prId", closePullRequest)
			// merge合并请求
			group.PUT("/merge/:prId", mergePullRequest)
			// review
			group.POST("/review", reviewPullRequest)
			// 展示时间轴
			group.GET("/listTimeline/:prId", listTimeline)
			// 添加评论
			group.POST("/addComment", addComment)
			// 删除评论
			group.DELETE("/deleteComment/:commentId", deleteComment)
			// 是否可合并
			group.GET("/canMerge/:prId", canMergePullRequest)
			// 获取review列表
			group.POST("/listReview")
			// 修改review
			group.POST("/updateReview")
		}
	})
}

func canMergePullRequest(c *gin.Context) {
	respDTO, err := pullrequestsrv.Outer.CanMergePullRequest(c, pullrequestsrv.CanMergePullRequestReqDTO{
		PrId:     getPrId(c),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[CanMergePullRequestRespVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: CanMergePullRequestRespVO{
			CanMerge:                respDTO.CanMerge,
			IsProtectedBranch:       respDTO.IsProtectedBranch,
			ReviewCountWhenCreatePr: respDTO.ReviewCountWhenCreatePr,
			ReviewerList:            respDTO.ReviewerList,
			DirectPushList:          respDTO.DirectPushList,
			ReviewCount:             respDTO.ReviewCount,
			GitCanMerge:             respDTO.GitCanMerge,
			GitConflictFiles:        respDTO.GitConflictFiles,
			GitCommitCount:          respDTO.GitCommitCount,
		},
	})
}

func addComment(c *gin.Context) {
	var req AddCommentReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.Outer.AddComment(c, pullrequestsrv.AddCommentReqDTO{
			PrId:      req.PrId,
			ReplyFrom: req.ReplyFrom,
			Comment:   req.Comment,
			HasReply:  req.HasReply,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteComment(c *gin.Context) {
	err := pullrequestsrv.Outer.DeleteComment(c, pullrequestsrv.DeleteCommentReqDTO{
		CommentId: cast.ToInt64(c.Param("commentId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listTimeline(c *gin.Context) {
	timelines, err := pullrequestsrv.Outer.ListTimeline(c, pullrequestsrv.ListTimelineReqDTO{
		PrId:     getPrId(c),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(timelines, func(t pullrequestsrv.TimelineDTO) (TimelineVO, error) {
		return TimelineVO{
			Id:      t.Id,
			PrId:    t.PrId,
			Action:  t.Action,
			Account: t.Account,
			Created: t.Created.Format(time.DateTime),
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]TimelineVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func getPullRequest(c *gin.Context) {
	request, err := pullrequestsrv.Outer.GetPullRequest(c, pullrequestsrv.GetPullRequestReqDTO{
		PrId:     getPrId(c),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := pr2Vo(request)
	c.JSON(http.StatusOK, ginutil.DataResp[PullRequestVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func statsPullRequest(c *gin.Context) {
	repoId := cast.ToInt64(c.Param("repoId"))
	stats, err := pullrequestsrv.Outer.GetStats(c, pullrequestsrv.GetStatsReqDTO{
		RepoId:   repoId,
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[StatsPullRequestVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: StatsPullRequestVO{
			TotalCount:  stats.TotalCount,
			OpenCount:   stats.OpenCount,
			ClosedCount: stats.ClosedCount,
			MergedCount: stats.MergedCount,
		},
	})
}

func listPullRequest(c *gin.Context) {
	var req ListPullRequestReqVO
	if util.ShouldBindQuery(&req, c) {
		requests, totalCount, err := pullrequestsrv.Outer.ListPullRequest(c, pullrequestsrv.ListPullRequestReqDTO{
			RepoId:    req.RepoId,
			Status:    req.Status,
			SearchKey: req.SearchKey,
			Page2Req:  req.Page2Req,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(requests, pr2Vo)
		c.JSON(http.StatusOK, ginutil.Page2Resp[PullRequestVO]{
			DataResp: ginutil.DataResp[[]PullRequestVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: totalCount,
		})
	}
}

func pr2Vo(t pullrequestsrv.PullRequestDTO) (PullRequestVO, error) {
	ret := PullRequestVO{
		Id:             t.Id,
		RepoId:         t.RepoId,
		Target:         t.Target,
		TargetType:     t.TargetType,
		TargetCommitId: t.TargetCommitId,
		Head:           t.Head,
		HeadType:       t.HeadType,
		HeadCommitId:   t.HeadCommitId,
		PrStatus:       t.PrStatus,
		CreateBy:       t.CreateBy,
		CloseBy:        t.CloseBy,
		MergeBy:        t.MergeBy,
		PrTitle:        t.PrTitle,
		CommentCount:   t.CommentCount,
		Created:        t.Created.Format(time.DateTime),
	}
	if t.Closed != nil {
		ret.Closed = t.Closed.Format(time.DateTime)
	}
	if t.Merged != nil {
		ret.Merged = t.Merged.Format(time.DateTime)
	}
	return ret, nil
}

func submitPullRequest(c *gin.Context) {
	var req SubmitPullRequestReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.Outer.SubmitPullRequest(c, pullrequestsrv.SubmitPullRequestReqDTO{
			RepoId:     req.RepoId,
			Target:     req.Target,
			TargetType: req.TargetType,
			Head:       req.Head,
			HeadType:   req.HeadType,
			Title:      req.Title,
			Comment:    req.Comment,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func closePullRequest(c *gin.Context) {
	err := pullrequestsrv.Outer.ClosePullRequest(c, pullrequestsrv.ClosePullRequestReqDTO{
		PrId:     getPrId(c),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func mergePullRequest(c *gin.Context) {
	err := pullrequestsrv.Outer.MergePullRequest(c, pullrequestsrv.MergePullRequestReqDTO{
		PrId:     getPrId(c),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func reviewPullRequest(c *gin.Context) {
	var req ReviewPullRequestReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.Outer.ReviewPullRequest(c, pullrequestsrv.ReviewPullRequestReqDTO{
			PrId:      req.Id,
			Status:    req.Status,
			ReviewMsg: req.ReviewMsg,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func getPrId(c *gin.Context) int64 {
	return cast.ToInt64(c.Param("prId"))
}
