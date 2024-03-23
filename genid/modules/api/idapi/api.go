package idapi

import (
	"github.com/LeeZXin/zall/genid/modules/service/idsrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/id")
		{
			group.Any("/snowflake/:batch", genSnowflakeIds)
			group.POST("/generator/insert", insertGenerator)
			group.Any("/generator/incr/:bizName/:step", incrGenerator)
		}
	})
}

func genSnowflakeIds(c *gin.Context) {
	batch := c.Param("batch")
	batchNum, err := strconv.ParseInt(batch, 10, 64)
	if err != nil {
		util.HandleApiErr(util.InvalidArgsError(), c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]int64]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     idsrv.Outer.GenSnowflakeIds(c, int(batchNum)),
	})
}

func insertGenerator(c *gin.Context) {
	var req InsertGeneratorReqVO
	if util.ShouldBindJSON(&req, c) {
		err := idsrv.Outer.InsertGenerator(c, req.BizName, req.CurrentId)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func incrGenerator(c *gin.Context) {
	stepStr := c.Param("step")
	step, err := strconv.ParseInt(stepStr, 10, 64)
	if err != nil {
		util.HandleApiErr(util.InvalidArgsError(), c)
		return
	}
	ids, err := idsrv.Outer.GenerateIdByBizName(c, c.Param("bizName"), int(step))
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]int64]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     ids,
	})
}
