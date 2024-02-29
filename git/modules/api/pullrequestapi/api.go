package pullrequestapi

import (
	"github.com/LeeZXin/zall/git/modules/service/pullrequestsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/pullRequest", apisession.CheckLogin)
		{
			// 创建合并请求
			group.POST("/submit", submitPullRequest)
			// 关闭合并请求
			group.POST("/close", closePullRequest)
			// merge合并请求
			group.POST("/merge", mergePullRequest)
			// review
			group.POST("/review", reviewPullRequest)
			// 获取review列表
			group.POST("/listReview")
			// 修改review
			group.POST("/updateReview")
		}
	})
}

func submitPullRequest(c *gin.Context) {
	var req SubmitPullRequestReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.Outer.SubmitPullRequest(c, pullrequestsrv.SubmitPullRequestReqDTO{
			RepoId:   req.RepoId,
			Target:   req.Target,
			Head:     req.Head,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func closePullRequest(c *gin.Context) {
	var req ClosePullRequestReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.Outer.ClosePullRequest(c, pullrequestsrv.ClosePullRequestReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func mergePullRequest(c *gin.Context) {
	var req MergePullRequestReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.Outer.MergePullRequest(c, pullrequestsrv.MergePullRequestReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func reviewPullRequest(c *gin.Context) {
	var req ReviewPullRequestReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.Outer.ReviewPullRequest(c, pullrequestsrv.ReviewPullRequestReqDTO{
			Id:        req.Id,
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
