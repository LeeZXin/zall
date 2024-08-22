package branchapi

import (
	"github.com/LeeZXin/zall/git/modules/service/branchsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		// 保护分支
		group := e.Group("/api/protectedBranch", apisession.CheckLogin)
		{
			// 新增保护分支
			group.POST("/create", createProtectedBranch)
			// 删除保护分支
			group.DELETE("/delete/:id", deleteProtectedBranch)
			// 查看保护分支
			group.GET("/list/:repoId", listProtectedBranch)
			// 编辑保护分支
			group.POST("/update", updateProtectedBranch)
		}
	})
}

func createProtectedBranch(c *gin.Context) {
	var req InsertProtectedBranchReqVO
	if util.ShouldBindJSON(&req, c) {
		err := branchsrv.CreateProtectedBranch(c, branchsrv.CreateProtectedBranchReqDTO{
			RepoId:   req.RepoId,
			Pattern:  req.Pattern,
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

func updateProtectedBranch(c *gin.Context) {
	var req UpdateProtectedBranchVO
	if util.ShouldBindJSON(&req, c) {
		err := branchsrv.UpdateProtectedBranch(c, branchsrv.UpdateProtectedBranchReqDTO{
			ProtectedBranchId: req.ProtectedBranchId,
			Pattern:           req.Pattern,
			Cfg:               req.Cfg,
			Operator:          apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteProtectedBranch(c *gin.Context) {
	err := branchsrv.DeleteProtectedBranch(c, branchsrv.DeleteProtectedBranchReqDTO{
		ProtectedBranchId: cast.ToInt64(c.Param("id")),
		Operator:          apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listProtectedBranch(c *gin.Context) {
	branchList, err := branchsrv.ListProtectedBranch(c, branchsrv.ListProtectedBranchReqDTO{
		RepoId:   cast.ToInt64(c.Param("repoId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(branchList, func(t branchsrv.ProtectedBranchDTO) (ProtectedBranchVO, error) {
		return ProtectedBranchVO{
			Id:      t.Id,
			Pattern: t.Pattern,
			Cfg:     t.Cfg,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ProtectedBranchVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
