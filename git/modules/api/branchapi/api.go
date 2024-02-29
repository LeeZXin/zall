package branchapi

import (
	"github.com/LeeZXin/zall/git/modules/service/branchsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		// 保护分支
		group := e.Group("/api/protectedBranch", apisession.CheckLogin)
		{
			// 新增保护分支
			group.POST("/insert", insertProtectedBranch)
			group.POST("/delete", deleteProtectedBranch)
			group.POST("/list", listProtectedBranch)
		}
	})
}

func insertProtectedBranch(c *gin.Context) {
	var req InsertProtectedBranchReqVO
	if util.ShouldBindJSON(&req, c) {
		err := branchsrv.Outer.InsertProtectedBranch(c, branchsrv.InsertProtectedBranchReqDTO{
			RepoId:   req.RepoId,
			Branch:   req.Branch,
			Cfg:      req.Cfg,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteProtectedBranch(c *gin.Context) {
	var req DeleteProtectedBranchReqVO
	if util.ShouldBindJSON(&req, c) {
		err := branchsrv.Outer.DeleteProtectedBranch(c, branchsrv.DeleteProtectedBranchReqDTO{
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

func listProtectedBranch(c *gin.Context) {
	var req ListProtectedBranchReqVO
	if util.ShouldBindJSON(&req, c) {
		branchList, err := branchsrv.Outer.ListProtectedBranch(c, branchsrv.ListProtectedBranchReqDTO{
			RepoId:   req.RepoId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		respVO := ListProtectedBranchRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
		}
		respVO.Branches, _ = listutil.Map(branchList, func(t branchsrv.ProtectedBranchDTO) (ProtectedBranchVO, error) {
			return ProtectedBranchVO{
				Id:     t.Id,
				Branch: t.Branch,
				Cfg:    t.Cfg,
			}, nil
		})
		c.JSON(http.StatusOK, respVO)
	}
}
