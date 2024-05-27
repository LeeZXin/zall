package oplogapi

import (
	"github.com/LeeZXin/zall/git/modules/service/oplogsrv"
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
		group := e.Group("/api/oplog/", apisession.CheckLogin)
		{
			group.GET("/page", pageLog)
		}
	})
}

func pageLog(c *gin.Context) {
	var req PageLogReqVO
	if util.ShouldBindQuery(&req, c) {
		logs, total, err := oplogsrv.Outer.PageOpLog(c, oplogsrv.PageOpLogReqDTO{
			RepoId:   req.RepoId,
			Account:  req.Account,
			DateStr:  req.DateStr,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(logs, func(t oplogsrv.OpLogDTO) (OpLogVO, error) {
			return OpLogVO{
				Id:      t.Id,
				Account: t.Account,
				Created: t.Created.Format(time.DateTime),
				Content: t.Content,
				ReqBody: t.ReqBody,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[OpLogVO]{
			DataResp: ginutil.DataResp[[]OpLogVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}
