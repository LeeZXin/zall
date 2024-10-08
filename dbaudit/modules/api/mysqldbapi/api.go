package mysqldbapi

import (
	"github.com/LeeZXin/zall/dbaudit/modules/service/mysqldbsrv"
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
		group := e.Group("/api/mysqldb", apisession.CheckLogin)
		{
			// 创建数据库
			group.POST("/create", createDb)
			// 编辑数据库
			group.POST("/update", updateDb)
			// 删除数据库
			group.DELETE("/delete/:dbId", deleteDb)
			// 数据库列表
			group.GET("/list", listDb)
			// 数据库列表
			group.GET("/all", allDb)
		}
		group = e.Group("/api/mysqlReadPermApply", apisession.CheckLogin)
		{
			// dba查看申请单
			group.GET("/list", listReadPermApplyByDba)
			// 同意申请
			group.PUT("/agree/:applyId", agreeReadPerm)
			// 不同意申请
			group.PUT("/disagree", disagreeReadPerm)
		}
		group = e.Group("/api/mysqlReadPerm", apisession.CheckLogin)
		{
			// 申请列表
			group.GET("/listApply", listReadPermApplyByOperator)
			// 用户读权限列表
			group.GET("/list", listReadPermByOperator)
			// 申请读权限
			group.POST("/apply", applyReadPerm)
			// 撤销申请
			group.PUT("/cancel/:applyId", cancelReadPerm)
			// 查看审批单
			group.GET("/getApply/:applyId", getReadPermApply)
			// 删除读权限
			group.DELETE("/delete/:permId", deleteReadPerm)
			// dba查看申请权限
			group.GET("/listManage", listReadPermByDba)
		}
		group = e.Group("/api/mysqlSearch", apisession.CheckLogin)
		{
			// 展示授权的数据库
			group.GET("/listAuthorizedDb", listAuthorizedDb)
			// 展示授权库
			group.GET("/listAuthorizedBase/:dbId", listAuthorizedBase)
			// 展示授权表
			group.GET("/listAuthorizedTable", listAuthorizedTable)
			// 展示建表语句
			group.GET("/getCreateTableSql", getCreateTableSql)
			// 展示建表语句
			group.GET("/showTableIndex", showTableIndex)
			// 执行查询sql
			group.POST("/executeSelectSql", executeSelectSql)
		}
		group = e.Group("/api/mysqlDataUpdate", apisession.CheckLogin)
		{
			// 审批列表
			group.GET("/listApply", listDataUpdateApplyByOperator)
			// 申请数据库修改
			group.POST("/apply", applyDataUpdate)
			// 查看执行计划
			group.GET("/explainApply/:applyId", explainDataUpdate)
			// 取消数据修改单申请
			group.PUT("/cancelApply/:applyId", cancelDataUpdate)
		}
		group = e.Group("/api/mysqlDataUpdateApply", apisession.CheckLogin)
		{
			// dba查看申请单
			group.GET("/list", listDataUpdateApplyByDba)
			// 同意申请
			group.PUT("/agree/:applyId", agreeDataUpdate)
			// 请求执行
			group.PUT("/askToExecute/:applyId", askToExecuteDataUpdate)
			// 不同意申请
			group.PUT("/disagree", disagreeDataUpdate)
			// 执行数据库修改单
			group.PUT("/execute/:applyId", executeDataUpdate)
		}
	})
}

func createDb(c *gin.Context) {
	var req CreateDbReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.CreateDb(c, mysqldbsrv.CreateDbReqDTO{
			Name:     req.Name,
			Config:   req.Config,
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
		err := mysqldbsrv.UpdateDb(c, mysqldbsrv.UpdateDbReqDTO{
			DbId:     req.DbId,
			Name:     req.Name,
			Config:   req.Config,
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
	err := mysqldbsrv.DeleteDb(c, mysqldbsrv.DeleteDbReqDTO{
		DbId:     cast.ToInt64(c.Param("dbId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listDb(c *gin.Context) {
	var req ListDbReqVO
	if util.ShouldBindQuery(&req, c) {
		dbs, total, err := mysqldbsrv.ListDb(c, mysqldbsrv.ListDbReqDTO{
			PageNum:  req.PageNum,
			Name:     req.Name,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(dbs, func(t mysqldbsrv.DbDTO) DbVO {
			return DbVO{
				Id:      t.Id,
				Name:    t.Name,
				Config:  t.Config,
				Created: t.Created.Format(time.DateTime),
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[DbVO]{
			DataResp: ginutil.DataResp[[]DbVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func allDb(c *gin.Context) {
	dbs, err := mysqldbsrv.ListSimpleDb(c, mysqldbsrv.ListSimpleDbReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(dbs, func(t mysqldbsrv.SimpleDbDTO) SimpleDbVO {
		return SimpleDbVO{
			Id:   t.Id,
			Name: t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleDbVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func applyReadPerm(c *gin.Context) {
	var req ApplyDbPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.ApplyReadPerm(c, mysqldbsrv.ApplyReadPermReqDTO{
			DbId:         req.DbId,
			AccessBase:   req.AccessBase,
			AccessTables: req.AccessTables,
			ApplyReason:  req.ApplyReason,
			ExpireDay:    req.ExpireDay,
			Operator:     apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func agreeReadPerm(c *gin.Context) {
	err := mysqldbsrv.AgreeReadPermApply(c, mysqldbsrv.AgreeReadPermApplyReqDTO{
		ApplyId:  cast.ToInt64(c.Param("applyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func disagreeReadPerm(c *gin.Context) {
	var req DisagreeReadPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.DisagreeReadPermApply(c, mysqldbsrv.DisagreeReadPermApplyReqDTO{
			ApplyId:        req.ApplyId,
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

func cancelReadPerm(c *gin.Context) {
	err := mysqldbsrv.CancelReadPermApply(c, mysqldbsrv.CancelReadPermApplyReqDTO{
		ApplyId:  cast.ToInt64(c.Param("applyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func getReadPermApply(c *gin.Context) {
	apply, err := mysqldbsrv.GetReadPermApply(c, mysqldbsrv.GetReadPermApplyReqDTO{
		ApplyId:  cast.ToInt64(c.Param("applyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[ReadPermApplyVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     readPermApply2Vo(apply),
	})
}

func listReadPermApplyByDba(c *gin.Context) {
	var req ListReadPermApplyByDbaReqVO
	if util.ShouldBindQuery(&req, c) {
		applies, total, err := mysqldbsrv.ListReadPermApplyByDba(c, mysqldbsrv.ListReadPermApplyByDbaReqDTO{
			DbId:        req.DbId,
			PageNum:     req.PageNum,
			ApplyStatus: req.ApplyStatus,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.Page2Resp[ReadPermApplyVO]{
			DataResp: ginutil.DataResp[[]ReadPermApplyVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     listutil.MapNe(applies, readPermApply2Vo),
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func readPermApplies2Vo(applies []mysqldbsrv.ReadPermApplyDTO) []ReadPermApplyVO {
	data := listutil.MapNe(applies, readPermApply2Vo)
	return data
}

func readPermApply2Vo(t mysqldbsrv.ReadPermApplyDTO) ReadPermApplyVO {
	return ReadPermApplyVO{
		Id:             t.Id,
		Account:        t.Account,
		DbId:           t.DbId,
		DbName:         t.DbName,
		AccessBase:     t.AccessBase,
		AccessTables:   t.AccessTables,
		ApplyStatus:    t.ApplyStatus,
		Auditor:        t.Auditor,
		ExpireDay:      t.ExpireDay,
		ApplyReason:    t.ApplyReason,
		DisagreeReason: t.DisagreeReason,
		Created:        t.Created.Format(time.DateTime),
		Updated:        t.Updated.Format(time.DateTime),
	}
}

func listReadPermApplyByOperator(c *gin.Context) {
	var req listReadPermApplyByOperatorReqVO
	if util.ShouldBindQuery(&req, c) {
		applies, next, err := mysqldbsrv.ListReadPermApplyByOperator(c, mysqldbsrv.ListReadPermApplyByOperatorReqDTO{
			PageNum:     req.PageNum,
			ApplyStatus: req.ApplyStatus,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.PageResp[ReadPermApplyVO]{
			DataResp: ginutil.DataResp[[]ReadPermApplyVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     readPermApplies2Vo(applies),
			},
			Next: next,
		})
	}
}

func listDataUpdateApplyByDba(c *gin.Context) {
	var req ListDataUpdateApplyByDbaReqVO
	if util.ShouldBindQuery(&req, c) {
		applies, total, err := mysqldbsrv.ListDataUpdateApplyByDba(c, mysqldbsrv.ListDataUpdateApplyByDbaReqDTO{
			PageNum:     req.PageNum,
			DbId:        req.DbId,
			ApplyStatus: req.ApplyStatus,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.Page2Resp[DataUpdateApplyVO]{
			DataResp: ginutil.DataResp[[]DataUpdateApplyVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     dataUpdateApplyDto2Vo(applies),
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func dataUpdateApplyDto2Vo(orders []mysqldbsrv.DataUpdateApplyDTO) []DataUpdateApplyVO {
	data := listutil.MapNe(orders, func(t mysqldbsrv.DataUpdateApplyDTO) DataUpdateApplyVO {
		return DataUpdateApplyVO{
			Id:                              t.Id,
			Account:                         t.Account,
			DbId:                            t.DbId,
			DbName:                          t.DbName,
			AccessBase:                      t.AccessBase,
			UpdateCmd:                       t.UpdateCmd,
			ApplyStatus:                     t.ApplyStatus,
			Executor:                        t.Executor,
			Auditor:                         t.Auditor,
			ApplyReason:                     t.ApplyReason,
			DisagreeReason:                  t.DisagreeReason,
			ExecuteLog:                      t.ExecuteLog,
			ExecuteImmediatelyAfterApproval: t.ExecuteImmediatelyAfterApproval,
			Created:                         t.Created.Format(time.DateTime),
			Updated:                         t.Updated.Format(time.DateTime),
			IsUnExecuted:                    t.ApplyStatus.IsUnExecuted(),
		}
	})
	return data
}

func listDataUpdateApplyByOperator(c *gin.Context) {
	var req ListDataUpdateApplyByOperatorReqVO
	if util.ShouldBindQuery(&req, c) {
		applies, total, err := mysqldbsrv.ListDataUpdateApplyByOperator(c, mysqldbsrv.ListDataUpdateApplyByOperatorReqDTO{
			PageNum:     req.PageNum,
			ApplyStatus: req.ApplyStatus,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.Page2Resp[DataUpdateApplyVO]{
			DataResp: ginutil.DataResp[[]DataUpdateApplyVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     dataUpdateApplyDto2Vo(applies),
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func listReadPermByOperator(c *gin.Context) {
	var req ListReadPermByOperatorReqVO
	if util.ShouldBindQuery(&req, c) {
		perms, total, err := mysqldbsrv.ListReadPermByOperator(c, mysqldbsrv.ListReadPermByOperatorReqDTO{
			PageNum:  req.PageNum,
			DbId:     req.DbId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(perms, func(t mysqldbsrv.ReadPermDTO) ReadPermVO {
			return ReadPermVO{
				Id:          t.Id,
				Account:     t.Account,
				DbId:        t.DbId,
				DbName:      t.DbName,
				AccessBase:  t.AccessBase,
				AccessTable: t.AccessTable,
				Created:     t.Created.Format(time.DateTime),
				Expired:     t.Expired.Format(time.DateTime),
				ApplyId:     t.ApplyId,
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[ReadPermVO]{
			DataResp: ginutil.DataResp[[]ReadPermVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func deleteReadPerm(c *gin.Context) {
	err := mysqldbsrv.DeleteReadPermByDba(c, mysqldbsrv.DeleteReadPermByDbaReqDTO{
		PermId:   cast.ToInt64(c.Param("permId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listReadPermByDba(c *gin.Context) {
	var req ListDbPermByDbaReqVO
	if util.ShouldBindQuery(&req, c) {
		perms, total, err := mysqldbsrv.ListReadPermByDba(c, mysqldbsrv.ListReadPermByDbaReqDTO{
			PageNum:  req.PageNum,
			DbId:     req.DbId,
			Account:  req.Account,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(perms, func(t mysqldbsrv.ReadPermDTO) ReadPermVO {
			return ReadPermVO{
				Id:          t.Id,
				Account:     t.Account,
				DbId:        t.DbId,
				DbName:      t.DbName,
				AccessBase:  t.AccessBase,
				AccessTable: t.AccessTable,
				Created:     t.Created.Format(time.DateTime),
				Expired:     t.Expired.Format(time.DateTime),
				ApplyId:     t.ApplyId,
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[ReadPermVO]{
			DataResp: ginutil.DataResp[[]ReadPermVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			TotalCount: total,
		})
	}
}

func listAuthorizedDb(c *gin.Context) {
	dbs, err := mysqldbsrv.ListAuthorizedDb(c, mysqldbsrv.ListAuthorizedDbReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(dbs, func(t mysqldbsrv.SimpleDbDTO) SimpleDbVO {
		return SimpleDbVO{
			Id:   t.Id,
			Name: t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleDbVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listAuthorizedBase(c *gin.Context) {
	bases, err := mysqldbsrv.ListAuthorizedBase(c, mysqldbsrv.ListAuthorizedBaseReqDTO{
		DbId:     cast.ToInt64(c.Param("dbId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     bases,
	})
}

func listAuthorizedTable(c *gin.Context) {
	tables, err := mysqldbsrv.ListAuthorizedTable(c, mysqldbsrv.ListAuthorizedTableReqDTO{
		DbId:       cast.ToInt64(c.Query("dbId")),
		AccessBase: c.Query("accessBase"),
		Operator:   apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	ret := listutil.MapNe(tables, func(t mysqldbsrv.AuthorizedTableDTO) AuthorizedTableVO {
		return AuthorizedTableVO{
			Table: t.Table,
			Size:  t.Size,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]AuthorizedTableVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     ret,
	})
}

func getCreateTableSql(c *gin.Context) {
	var req GetCreateTableSqlReqVO
	if util.ShouldBindQuery(&req, c) {
		sql, err := mysqldbsrv.GetCreateTableSql(c, mysqldbsrv.GetCreateSqlReqDTO{
			DbId:        req.DbId,
			AccessBase:  req.AccessBase,
			AccessTable: req.AccessTable,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[string]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     sql,
		})
	}
}

func showTableIndex(c *gin.Context) {
	var req ShowTableIndexReqVO
	if util.ShouldBindQuery(&req, c) {
		columns, data, err := mysqldbsrv.ShowTableIndex(c, mysqldbsrv.ShowTableIndexReqDTO{
			DbId:        req.DbId,
			AccessBase:  req.AccessBase,
			AccessTable: req.AccessTable,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := make([]map[string]string, 0)
		for _, datum := range data {
			if len(datum) != len(columns) {
				continue
			}
			item := make(map[string]string)
			for i := range columns {
				item[columns[i]] = datum[i]
			}
			ret = append(ret, item)
		}
		c.JSON(http.StatusOK, ginutil.DataResp[ShowTableIndexRespVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data: ShowTableIndexRespVO{
				Columns: columns,
				Data:    ret,
			},
		})
	}
}

func executeSelectSql(c *gin.Context) {
	var req ExecuteSelectSqlReqVO
	if util.ShouldBindJSON(&req, c) {
		result, err := mysqldbsrv.ExecuteSelectSql(c, mysqldbsrv.ExecuteSelectSqlReqDTO{
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
		c.JSON(http.StatusOK, ginutil.DataResp[ExecuteSelectSqlResultVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data: ExecuteSelectSqlResultVO{
				Columns:  result.Columns,
				Data:     result.ToMap(),
				Duration: result.Duration.String(),
			},
		})
	}
}

func applyDataUpdate(c *gin.Context) {
	var req ApplyDataUpdateReqVO
	if util.ShouldBindJSON(&req, c) {
		results, allPass, err := mysqldbsrv.ApplyDataUpdate(c, mysqldbsrv.ApplyDataUpdateReqDTO{
			DbId:                            req.DbId,
			AccessBase:                      req.AccessBase,
			Cmd:                             req.Cmd,
			ApplyReason:                     req.ApplyReason,
			ExecuteImmediatelyAfterApproval: req.ExecuteImmediatelyAfterApproval,
			Operator:                        apisession.MustGetLoginUser(c),
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

func explainDataUpdate(c *gin.Context) {
	ret, err := mysqldbsrv.ExplainDataUpdate(c, mysqldbsrv.ExplainDataUpdateReqDTO{
		ApplyId:  cast.ToInt64(c.Param("applyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[string]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     ret,
	})
}

func agreeDataUpdate(c *gin.Context) {
	err := mysqldbsrv.AgreeDataUpdateApply(c, mysqldbsrv.AgreeDbUpdateReqDTO{
		ApplyId:  cast.ToInt64(c.Param("applyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func disagreeDataUpdate(c *gin.Context) {
	var req DisagreeDataUpdateReqVO
	if util.ShouldBindJSON(&req, c) {
		err := mysqldbsrv.DisagreeDataUpdateApply(c, mysqldbsrv.DisagreeDataUpdateApplyReqDTO{
			ApplyId:        req.ApplyId,
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

func cancelDataUpdate(c *gin.Context) {
	err := mysqldbsrv.CancelDataUpdateApply(c, mysqldbsrv.CancelDataUpdateApplyReqDTO{
		ApplyId:  cast.ToInt64(c.Param("applyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func askToExecuteDataUpdate(c *gin.Context) {
	err := mysqldbsrv.AskToExecuteDataUpdateApply(c, mysqldbsrv.AskToExecuteDataUpdateApplyReqDTO{
		ApplyId:  cast.ToInt64(c.Param("applyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func executeDataUpdate(c *gin.Context) {
	err := mysqldbsrv.ExecuteDataUpdateApply(c, mysqldbsrv.ExecuteDataUpdateApplyReqDTO{
		ApplyId:  cast.ToInt64(c.Param("applyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}
