package approvalapi

import (
	"github.com/LeeZXin/zall/approval/modules/service/approvalsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/approval", apisession.CheckLogin)
		{
			group.POST("/agree", agreeApproval)
			group.POST("/disagree", disagreeApproval)
		}
	})
}

func agreeApproval(c *gin.Context) {
	var req AgreeApprovalReqVO
	if util.ShouldBindJSON(&req, c) {
		err := approvalsrv.Outer.Agree(c, approvalsrv.AgreeFlowReqDTO{
			NotifyId: req.NotifyId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func disagreeApproval(c *gin.Context) {
	var req DisagreeApprovalReqVO
	if util.ShouldBindJSON(&req, c) {
		err := approvalsrv.Outer.Disagree(c, approvalsrv.DisagreeFlowReqDTO{
			NotifyId: req.NotifyId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
