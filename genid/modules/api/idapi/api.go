package idapi

import (
	"github.com/LeeZXin/zall/genid/modules/service/idsrv"
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
		}
	})
}

func genSnowflakeIds(c *gin.Context) {
	batch := c.Param("batch")
	batchNum, err := strconv.ParseInt(batch, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "wrong batchNum")
		return
	}
	ids := idsrv.Outer.GenSnowflakeIds(c, int(batchNum))
	c.JSON(http.StatusOK, SnowFlakeRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     ids,
	})
}
