package mysqldbapi

import (
	"github.com/LeeZXin/zall/dbaudit/modules/service/mysqldbsrv"
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
		group := e.Group("/api/mysqldb", apisession.CheckLogin)
		{
			group.POST("/insert", insertDb)
			group.POST("/update", updateDb)
			group.POST("/delete", deleteDb)
			group.GET("/list", listDb)
			group.GET("/all", allDb)
		}
		group = e.Group("/api/mysqldbPermOrder", apisession.CheckLogin)
		{
			group.POST("/list", listPermOrder)
			group.POST("/listApplied", listAppliedPermOrder)
			group.POST("/apply", applyDbPerm)
			group.POST("/agree", agreeDbPerm)
			group.POST("/disagree", disagreeDbPerm)
			group.POST("/cancel", cancelDbPerm)
		}
		group = e.Group("/api/mysqldbPerm", apisession.CheckLogin)
		{
			group.POST("/list", listDbPerm)
			group.POST("/delete", deleteDbPerm)
			group.POST("/listByAccount", listDbPermByAccount)
		}
		group = e.Group("/api/mysqldbSearch", apisession.CheckLogin)
		{
			group.POST("/allBases", allBases)
			group.POST("/allTables", allTables)
			group.POST("/searchDb", searchDb)
		}
		group = e.Group("/api/mysqldbUpdateOrder", apisession.CheckLogin)
		{
			group.POST("/apply", applyDbUpdate)
			group.POST("/list", listUpdateOrder)
			group.POST("/listApplied", listAppliedUpdateOrder)
			group.POST("/agree", agreeDbUpdate)
			group.POST("/disagree", disagreeDbUpdate)
			group.POST("/cancel", cancelDbUpdate)
			group.POST("/execute", executeDbUpdate)
			group.POST("/askToExecute", askToExecuteDbUpdate)
		}
	})
}

func insertDb(c *gin.Context) {
	var req InsertDbReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.InsertDb(c, mysqldbsrv.InsertDbReqDTO{
			Name:     req.Name,
			DbHost:   req.DbHost,
			Username: req.Username,
			Password: req.Password,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateDb(c *gin.Context) {
	var req UpdateDbReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.UpdateDb(c, mysqldbsrv.UpdateDbReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			DbHost:   req.DbHost,
			Username: req.Username,
			Password: req.Password,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteDb(c *gin.Context) {
	var req DeleteDbReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.DeleteDb(c, mysqldbsrv.DeleteDbReqDTO{
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

func listDb(c *gin.Context) {
	dbs, err := mysqldbsrv.Outer.ListDb(c, mysqldbsrv.ListDbReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(dbs, func(t mysqldbsrv.DbDTO) (DbVO, error) {
		return DbVO{
			Id:       t.Id,
			Name:     t.Name,
			DbHost:   t.DbHost,
			Username: t.Username,
			Password: t.Password,
			Created:  t.Created.Format(time.DateTime),
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]DbVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func allDb(c *gin.Context) {
	dbs, err := mysqldbsrv.Outer.ListSimpleDb(c, mysqldbsrv.ListSimpleDbReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(dbs, func(t mysqldbsrv.SimpleDbDTO) (SimpleDbVO, error) {
		return SimpleDbVO{
			Id:     t.Id,
			Name:   t.Name,
			DbHost: t.DbHost,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleDbVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func applyDbPerm(c *gin.Context) {
	var req ApplyDbPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.ApplyDbPerm(c, mysqldbsrv.ApplyDbPermReqDTO{
			DbId:         req.DbId,
			AccessBase:   req.AccessBase,
			AccessTables: req.AccessTables,
			Reason:       req.Reason,
			ExpireDay:    req.ExpireDay,
			PermType:     req.PermType,
			Operator:     apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func agreeDbPerm(c *gin.Context) {
	var req AgreeDbPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.AgreeDbPerm(c, mysqldbsrv.AgreeDbPermReqDTO{
			OrderId:  req.OrderId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func disagreeDbPerm(c *gin.Context) {
	var req DisagreeDbPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.DisagreeDbPerm(c, mysqldbsrv.DisagreeDbPermReqDTO{
			OrderId:        req.OrderId,
			DisagreeReason: req.DisagreeReason,
			Operator:       apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func cancelDbPerm(c *gin.Context) {
	var req CancelDbPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.CancelDbPerm(c, mysqldbsrv.CancelDbPermReqDTO{
			OrderId:  req.OrderId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listPermOrder(c *gin.Context) {
	var req ListPermApprovalOrderReqVO
	if util.ShouldBindJSON(&req, c) {
		orders, next, err := mysqldbsrv.Outer.ListPermApprovalOrder(c, mysqldbsrv.ListPermApprovalOrderReqDTO{
			Cursor:      req.Cursor,
			Limit:       req.Limit,
			OrderStatus: req.OrderStatus,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.PageResp[[]PermApprovalOrderVO]{
			DataResp: ginutil.DataResp[[]PermApprovalOrderVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     permOrdersDto2Vo(orders),
			},
			Next: next,
		})
	}
}

func permOrdersDto2Vo(orders []mysqldbsrv.PermApprovalOrderDTO) []PermApprovalOrderVO {
	data, _ := listutil.Map(orders, func(t mysqldbsrv.PermApprovalOrderDTO) (PermApprovalOrderVO, error) {
		return PermApprovalOrderVO{
			Id:           t.Id,
			Account:      t.Account,
			DbId:         t.DbId,
			DbHost:       t.DbHost,
			DbName:       t.DbName,
			AccessBase:   t.AccessBase,
			AccessTables: t.AccessTables,
			PermType:     t.PermType.Readable(),
			OrderStatus:  t.OrderStatus.Readable(),
			Auditor:      t.Auditor,
			ExpireDay:    t.ExpireDay,
			Reason:       t.Reason,
			Created:      t.Created.Format(time.DateTime),
		}, nil
	})
	return data
}

func listAppliedPermOrder(c *gin.Context) {
	var req ListAppliedPermApprovalOrderReqVO
	if util.ShouldBindJSON(&req, c) {
		orders, next, err := mysqldbsrv.Outer.ListAppliedPermApprovalOrder(c, mysqldbsrv.ListAppliedPermApprovalOrderReqDTO{
			Cursor:      req.Cursor,
			Limit:       req.Limit,
			OrderStatus: req.OrderStatus,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.PageResp[[]PermApprovalOrderVO]{
			DataResp: ginutil.DataResp[[]PermApprovalOrderVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     permOrdersDto2Vo(orders),
			},
			Next: next,
		})
	}
}

func listUpdateOrder(c *gin.Context) {
	var req ListUpdateApprovalOrderReqVO
	if util.ShouldBindJSON(&req, c) {
		orders, next, err := mysqldbsrv.Outer.ListUpdateApprovalOrder(c, mysqldbsrv.ListUpdateApprovalOrderReqDTO{
			Cursor:      req.Cursor,
			Limit:       req.Limit,
			OrderStatus: req.OrderStatus,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.PageResp[[]UpdateApprovalOrderVO]{
			DataResp: ginutil.DataResp[[]UpdateApprovalOrderVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     updateOrdersDto2Vo(orders),
			},
			Next: next,
		})
	}
}

func updateOrdersDto2Vo(orders []mysqldbsrv.UpdateApprovalOrderDTO) []UpdateApprovalOrderVO {
	data, _ := listutil.Map(orders, func(t mysqldbsrv.UpdateApprovalOrderDTO) (UpdateApprovalOrderVO, error) {
		return UpdateApprovalOrderVO{
			Id:          t.Id,
			Name:        t.Name,
			Account:     t.Account,
			DbId:        t.DbId,
			DbHost:      t.DbHost,
			DbName:      t.DbName,
			AccessBase:  t.AccessBase,
			UpdateCmd:   t.UpdateCmd,
			OrderStatus: t.OrderStatus.Readable(),
			Auditor:     t.Auditor,
			ExecuteLog:  t.ExecuteLog,
			Created:     t.Created.Format(time.DateTime),
		}, nil
	})
	return data
}

func listAppliedUpdateOrder(c *gin.Context) {
	var req ListAppliedUpdateApprovalOrderReqVO
	if util.ShouldBindJSON(&req, c) {
		orders, next, err := mysqldbsrv.Outer.ListAppliedUpdateApprovalOrder(c, mysqldbsrv.ListAppliedUpdateApprovalOrderReqDTO{
			Cursor:      req.Cursor,
			Limit:       req.Limit,
			OrderStatus: req.OrderStatus,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.PageResp[[]UpdateApprovalOrderVO]{
			DataResp: ginutil.DataResp[[]UpdateApprovalOrderVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     updateOrdersDto2Vo(orders),
			},
			Next: next,
		})
	}
}

func listDbPerm(c *gin.Context) {
	var req ListDbPermReqVO
	if util.ShouldBindJSON(&req, c) {
		perms, next, err := mysqldbsrv.Outer.ListDbPerm(c, mysqldbsrv.ListDbPermReqDTO{
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(perms, func(t mysqldbsrv.PermDTO) (PermVO, error) {
			return PermVO{
				Id:          t.Id,
				Account:     t.Account,
				DbId:        t.DbId,
				DbHost:      t.DbHost,
				DbName:      t.DbName,
				AccessBase:  t.AccessBase,
				AccessTable: t.AccessTable,
				PermType:    t.PermType.Readable(),
				Created:     t.Created.Format(time.DateTime),
				Expired:     t.Expired.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[[]PermVO]{
			DataResp: ginutil.DataResp[[]PermVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: next,
		})
	}
}

func deleteDbPerm(c *gin.Context) {
	var req CancelDbPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.CancelDbPerm(c, mysqldbsrv.CancelDbPermReqDTO{
			OrderId:  req.OrderId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listDbPermByAccount(c *gin.Context) {
	var req ListDbPermByAccountReqVO
	if util.ShouldBindJSON(&req, c) {
		perms, next, err := mysqldbsrv.Outer.ListDbPermByAccount(c, mysqldbsrv.ListDbPermByAccountReqDTO{
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Account:  req.Account,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(perms, func(t mysqldbsrv.PermDTO) (PermVO, error) {
			return PermVO{
				Id:          t.Id,
				Account:     t.Account,
				DbId:        t.DbId,
				DbHost:      t.DbHost,
				DbName:      t.DbName,
				AccessTable: t.AccessTable,
				PermType:    t.PermType.Readable(),
				Created:     t.Created.Format(time.DateTime),
				Expired:     t.Expired.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[[]PermVO]{
			DataResp: ginutil.DataResp[[]PermVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: next,
		})
	}
}

func allTables(c *gin.Context) {
	var req AllTablesReqVO
	if util.ShouldBindJSON(&req, c) {
		tables, err := mysqldbsrv.Outer.AllTables(c, mysqldbsrv.AllTablesReqDTO{
			DbId:       req.DbId,
			AccessBase: req.AccessBase,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     tables,
		})
	}
}

func allBases(c *gin.Context) {
	var req AllBasesReqVO
	if util.ShouldBindJSON(&req, c) {
		tables, err := mysqldbsrv.Outer.AllBases(c, mysqldbsrv.AllBasesReqDTO{
			DbId:     req.DbId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     tables,
		})
	}
}

func searchDb(c *gin.Context) {
	var req SearchDbReqVO
	if util.ShouldBindJSON(&req, c) {
		columns, result, err := mysqldbsrv.Outer.SearchDb(c, mysqldbsrv.SearchDbReqDTO{
			DbId:       req.DbId,
			AccessBase: req.AccessBase,
			Cmd:        req.Cmd,
			Limit:      req.Limit,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[SearchDbResultVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data: SearchDbResultVO{
				Columns: columns,
				Result:  result,
			},
		})
	}
}

func applyDbUpdate(c *gin.Context) {
	var req ApplyDbUpdateReqVO
	if util.ShouldBindJSON(&req, c) {
		results, allPass, err := mysqldbsrv.Outer.ApplyDbUpdate(c, mysqldbsrv.ApplyDbUpdateReqDTO{
			Name:       req.Name,
			DbId:       req.DbId,
			AccessBase: req.AccessBase,
			Cmd:        req.Cmd,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[ApplyDbUpdateResultVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data: ApplyDbUpdateResultVO{
				Results: results,
				AllPass: allPass,
			},
		})
	}
}

func agreeDbUpdate(c *gin.Context) {
	var req AgreeDbUpdateReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.AgreeDbUpdate(c, mysqldbsrv.AgreeDbUpdateReqDTO{
			OrderId:  req.OrderId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func disagreeDbUpdate(c *gin.Context) {
	var req DisagreeDbUpdateReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.DisagreeDbUpdate(c, mysqldbsrv.DisagreeDbUpdateReqDTO{
			OrderId:        req.OrderId,
			DisagreeReason: req.DisagreeReason,
			Operator:       apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func cancelDbUpdate(c *gin.Context) {
	var req CancelDbUpdateReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.CancelDbUpdate(c, mysqldbsrv.CancelDbUpdateReqDTO{
			OrderId:  req.OrderId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func askToExecuteDbUpdate(c *gin.Context) {
	var req AskToExecuteDbUpdateReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.AskToExecuteDbUpdate(c, mysqldbsrv.AskToExecuteDbUpdateReqDTO{
			OrderId:  req.OrderId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func executeDbUpdate(c *gin.Context) {
	var req ExecuteDbUpdateReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.Outer.ExecuteDbUpdate(c, mysqldbsrv.ExecuteDbUpdateReqDTO{
			OrderId:  req.OrderId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
