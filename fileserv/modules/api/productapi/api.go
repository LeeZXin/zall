package productapi

import (
	"github.com/LeeZXin/zall/fileserv/modules/service/productsrv"
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
		// 简单制品库
		group := e.Group("/api/product", apisession.CheckLogin)
		{
			group.GET("/list/:app", listProduct)
		}
	})
}

func listProduct(c *gin.Context) {
	products, err := productsrv.Outer.ListProduct(c, productsrv.ListProductReqDTO{
		AppId:    c.Param("app"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(products, func(t productsrv.ProductDTO) (ProductVO, error) {
		return ProductVO{
			Name:    t.Name,
			Creator: t.Creator,
			Created: t.Created.Format(time.DateTime),
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ProductVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
