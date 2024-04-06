package dbapi

import (
	"github.com/LeeZXin/zall/dbaudit/modules/service/dbsrv"
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
		group := e.Group("/api/db", apisession.CheckLogin)
		{
			group.POST("/insert", insertDb)
			group.POST("/update", updateDb)
			group.POST("/delete", deleteDb)
			group.GET("/list", listDb)
			group.GET("/all", allDb)
		}
		group = e.Group("/api/dbPermOrder", apisession.CheckLogin)
		{
			group.POST("/list", listApprovalOrder)
			group.POST("/apply", applyDbPerm)
			group.POST("/agree", agreeDbPerm)
			group.POST("/disagree", disagreeDbPerm)
			group.POST("/cancel", cancelDbPerm)
		}
		group = e.Group("/api/dbPerm", apisession.CheckLogin)
		{
			group.POST("/list", listDbPerm)
			group.POST("/delete", deleteDbPerm)
			group.POST("/listByAccount", listDbPermByAccount)
		}
		group = e.Group("/api/dbSearch", apisession.CheckLogin)
		{
			group.POST("/allBases", allBases)
			group.POST("/allTables", allTables)
			group.POST("/searchDb", searchDb)
		}
	})
}

func insertDb(c *gin.Context) {
	var req InsertDbReqVO
	if util.ShouldBindJSON(&req, c) {
		err := dbsrv.Outer.InsertDb(c, dbsrv.InsertDbReqDTO{
			Name:     req.Name,
			DbHost:   req.DbHost,
			Username: req.Username,
			Password: req.Password,
			DbType:   req.DbType,
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
		err := dbsrv.Outer.UpdateDb(c, dbsrv.UpdateDbReqDTO{
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
		err := dbsrv.Outer.DeleteDb(c, dbsrv.DeleteDbReqDTO{
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
	dbs, err := dbsrv.Outer.ListDb(c, dbsrv.ListDbReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(dbs, func(t dbsrv.DbDTO) (DbVO, error) {
		return DbVO{
			Id:       t.Id,
			Name:     t.Name,
			DbHost:   t.DbHost,
			Username: t.Username,
			Password: t.Password,
			DbType:   t.DbType.Readable(),
			Created:  t.Created.Format(time.DateTime),
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]DbVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func allDb(c *gin.Context) {
	dbs, err := dbsrv.Outer.ListSimpleDb(c, dbsrv.ListSimpleDbReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(dbs, func(t dbsrv.SimpleDbDTO) (SimpleDbVO, error) {
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
		err := dbsrv.Outer.ApplyDbPerm(c, dbsrv.ApplyDbPermReqDTO{
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
		err := dbsrv.Outer.AgreeDbPerm(c, dbsrv.AgreeDbPermReqDTO{
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
		err := dbsrv.Outer.DisagreeDbPerm(c, dbsrv.DisagreeDbPermReqDTO{
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

func cancelDbPerm(c *gin.Context) {
	var req CancelDbPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := dbsrv.Outer.CancelDbPerm(c, dbsrv.CancelDbPermReqDTO{
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

func listApprovalOrder(c *gin.Context) {
	var req ListApprovalOrderReqVO
	if util.ShouldBindJSON(&req, c) {
		orders, next, err := dbsrv.Outer.ListPermApprovalOrder(c, dbsrv.ListPermApprovalOrderReqDTO{
			Cursor:      req.Cursor,
			Limit:       req.Limit,
			OrderStatus: req.OrderStatus,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(orders, func(t dbsrv.ApprovalOrderDTO) (ApprovalOrderVO, error) {
			return ApprovalOrderVO{
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
		c.JSON(http.StatusOK, ginutil.PageResp[[]ApprovalOrderVO]{
			DataResp: ginutil.DataResp[[]ApprovalOrderVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: next,
		})
	}
}

func listDbPerm(c *gin.Context) {
	var req ListDbPermReqVO
	if util.ShouldBindJSON(&req, c) {
		perms, next, err := dbsrv.Outer.ListDbPerm(c, dbsrv.ListDbPermReqDTO{
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(perms, func(t dbsrv.PermDTO) (PermVO, error) {
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
		err := dbsrv.Outer.CancelDbPerm(c, dbsrv.CancelDbPermReqDTO{
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
		perms, next, err := dbsrv.Outer.ListDbPermByAccount(c, dbsrv.ListDbPermByAccountReqDTO{
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Account:  req.Account,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(perms, func(t dbsrv.PermDTO) (PermVO, error) {
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
		tables, err := dbsrv.Outer.AllTables(c, dbsrv.AllTablesReqDTO{
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
		tables, err := dbsrv.Outer.AllBases(c, dbsrv.AllBasesReqDTO{
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
		columns, result, err := dbsrv.Outer.SearchDb(c, dbsrv.SearchDbReqDTO{
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
