package approvalapi

import (
	"github.com/LeeZXin/zall/approval/modules/service/approvalsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/approval", apisession.CheckLogin)
		{
			group.POST("/agree", agreeApproval)
			group.POST("/disagree", disagreeApproval)
			group.POST("/insert", insertProcess)
			group.POST("/update", updateProcess)
			group.POST("/insertFlow", insertFlow)
			group.POST("/cancelFlow", cancelFlow)
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

func insertProcess(c *gin.Context) {
	var req InsertCustomProcessReqVO
	if util.ShouldBindJSON(&req, c) {
		err := approvalsrv.Outer.InsertCustomProcess(c, approvalsrv.InsertCustomProcessReqDTO{
			Pid:      req.Pid,
			Name:     req.Name,
			GroupId:  req.GroupId,
			Process:  req.Process,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateProcess(c *gin.Context) {
	var req UpdateCustomProcessReqVO
	if util.ShouldBindJSON(&req, c) {
		err := approvalsrv.Outer.UpdateCustomProcess(c, approvalsrv.UpdateCustomProcessReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			GroupId:  req.GroupId,
			Process:  req.Process,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func insertFlow(c *gin.Context) {
	var req InsertCustomFlowReqVO
	if util.ShouldBindJSON(&req, c) {
		errKeys, err := approvalsrv.Outer.InsertCustomFlow(c, approvalsrv.InsertCustomFlowReqDTO{
			Pid:      req.Pid,
			BizId:    req.BizId,
			Kvs:      req.Kvs,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		if len(errKeys) > 0 {
			c.JSON(http.StatusOK, InsertCustomFlowRespVO{
				BaseResp: ginutil.DefaultSuccessResp,
				ErrKeys:  errKeys,
			})
			return
		}
		util.DefaultOkResponse(c)
	}
}

func cancelFlow(c *gin.Context) {
	var req CancelCustomFlowReqVO
	if util.ShouldBindJSON(&req, c) {
		err := approvalsrv.Outer.CancelCustomFlow(c, approvalsrv.CancelCustomFlowReqDTO{
			FlowId:   req.FlowId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
