package approvalapi

import (
	"github.com/LeeZXin/zall/approval/modules/service/approvalsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/customApproval", apisession.CheckLogin)
		{
			group.POST("/agree", agreeApproval)
			group.POST("/disagree", disagreeApproval)
			group.Any("/list", listApproval)
		}
		group = e.Group("/api/customApprovalProcess", apisession.CheckLogin)
		{
			group.POST("/insert", insertProcess)
			group.POST("/update", updateProcess)
			group.POST("/delete", deleteProcess)
			group.POST("/list", listProcess)
		}
		group = e.Group("/api/customApprovalFlow", apisession.CheckLogin)
		{
			group.POST("/insert", insertFlow)
			group.POST("/cancel", cancelFlow)
			group.POST("/get", getFlow)
			group.POST("/list", listFlow)
			group.POST("/listOperate", listOperate)
		}
		group = e.Group("/api/approvalGroup", apisession.CheckLogin)
		{
			group.Any("/list", listGroup)
			group.POST("/insert", insertGroup)
			group.POST("/update", updateGroup)
			group.POST("/delete", deleteGroup)
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
			IconUrl:  req.IconUrl,
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
			IconUrl:  req.IconUrl,
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

func deleteProcess(c *gin.Context) {
	var req DeleteCustomProcessReqVO
	if util.ShouldBindJSON(&req, c) {
		err := approvalsrv.Outer.DeleteCustomProcess(c, approvalsrv.DeleteCustomProcessReqDTO{
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

func listProcess(c *gin.Context) {
	var req ListCustomProcessReqVO
	if util.ShouldBindJSON(&req, c) {
		processes, err := approvalsrv.Outer.ListCustomProcess(c, approvalsrv.ListCustomProcessReqDTO{
			GroupId:  req.GroupId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := ListCustomProcessRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
		}
		ret.Data, _ = listutil.Map(processes, func(t approvalsrv.ProcessDTO) (ProcessVO, error) {
			return ProcessVO{
				Id:      t.Id,
				Pid:     t.Pid,
				GroupId: t.GroupId,
				Name:    t.Name,
				Content: t.Content,
				IconUrl: t.IconUrl,
				Created: t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ret)
	}
}

func listApproval(c *gin.Context) {
	groups, err := approvalsrv.Outer.ListAllGroupProcess(c, approvalsrv.ListAllGroupProcessReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	ret := ListCustomGroupsRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
	}
	ret.Data, _ = listutil.Map(groups, func(t approvalsrv.GroupProcessDTO) (GroupProcessVO, error) {
		processes, _ := listutil.Map(t.Processes, func(t approvalsrv.SimpleProcessDTO) (SimpleProcessVO, error) {
			return SimpleProcessVO{
				Id:      t.Id,
				Name:    t.Name,
				IconUrl: t.IconUrl,
			}, nil
		})
		return GroupProcessVO{
			Id:        t.Id,
			Name:      t.Name,
			Processes: processes,
		}, nil
	})
	c.JSON(http.StatusOK, ret)
}

func getFlow(c *gin.Context) {
	var req GetFlowReqVO
	if util.ShouldBindJSON(&req, c) {
		flow, err := approvalsrv.Outer.GetFlowDetail(c, approvalsrv.GetFlowDetailReqDTO{
			FlowId:   req.FlowId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := FlowDetailVO{
			Id:          flow.Id,
			ProcessName: flow.ProcessName,
			FlowStatus:  flow.FlowStatus.Readable(),
			Creator:     flow.Creator,
			Created:     flow.Created.Format(time.DateTime),
			Kvs:         flow.Kvs,
			Process:     flow.Process,
		}
		data.NotifyList, _ = listutil.Map(flow.NotifyList, func(t approvalsrv.NotifyDTO) (NotifyVO, error) {
			return NotifyVO{
				Account:   t.Account,
				FlowIndex: t.FlowIndex,
				Done:      t.Done,
				Op:        t.Op.Readable(),
				Updated:   t.Updated.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, GetFlowRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func insertGroup(c *gin.Context) {
	var req InsertGroupReqVO
	if util.ShouldBindJSON(&req, c) {
		err := approvalsrv.Outer.InsertGroup(c, approvalsrv.InsertGroupReqDTO{
			Name:     req.Name,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteGroup(c *gin.Context) {
	var req DeleteGroupReqVO
	if util.ShouldBindJSON(&req, c) {
		err := approvalsrv.Outer.DeleteGroup(c, approvalsrv.DeleteGroupReqDTO{
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

func updateGroup(c *gin.Context) {
	var req UpdateGroupReqVO
	if util.ShouldBindJSON(&req, c) {
		err := approvalsrv.Outer.UpdateGroup(c, approvalsrv.UpdateGroupReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listGroup(c *gin.Context) {
	groups, err := approvalsrv.Outer.ListGroup(c, approvalsrv.ListGroupReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	ret := ListGroupRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
	}
	ret.Data, _ = listutil.Map(groups, func(t approvalsrv.GroupDTO) (GroupVO, error) {
		return GroupVO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, ret)
}

func listFlow(c *gin.Context) {
	var req ListCustomFlowReqVO
	if util.ShouldBindJSON(&req, c) {
		flows, err := approvalsrv.Outer.ListCustomFlow(c, approvalsrv.ListCustomFlowReqDTO{
			DayTime:  req.DayTime,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := ListCustomFlowRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
		}
		ret.Data, _ = listutil.Map(flows, func(t approvalsrv.FlowDTO) (FlowVO, error) {
			return FlowVO{
				Id:          t.Id,
				ProcessName: t.ProcessName,
				FlowStatus:  t.FlowStatus.Readable(),
				Creator:     t.Creator,
				Created:     t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ret)
	}
}

func listOperate(c *gin.Context) {
	var req ListOperateFlowReqVO
	if util.ShouldBindJSON(&req, c) {
		flows, err := approvalsrv.Outer.ListOperateFlow(c, approvalsrv.ListOperateFlowReqDTO{
			DayTime:  req.DayTime,
			Done:     req.Done,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := ListOperateFlowRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
		}
		ret.Data, _ = listutil.Map(flows, func(t approvalsrv.FlowDTO) (FlowVO, error) {
			return FlowVO{
				Id:          t.Id,
				ProcessName: t.ProcessName,
				FlowStatus:  t.FlowStatus.Readable(),
				Creator:     t.Creator,
				Created:     t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ret)
	}
}
