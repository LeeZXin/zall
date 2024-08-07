package fileapi

import (
	"fmt"
	"github.com/LeeZXin/zall/fileserv/modules/service/filesrv"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

var (
	productToken string
)

func InitApi() {
	productToken = static.GetString("files.product.token")
	filesrv.InitStorage()
	cfgsrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/files/avatar", apisession.CheckLogin)
		{
			// 上传头像
			group.POST("/upload", uploadAvatar)
			// 获取头像
			group.GET("/get/:name", getAvatar)
		}
		// 简单制品库
		group = e.Group("/api/files/product", checkProductToken)
		{
			// 上传制品
			// curl -F "file=@/Users/lizexin/Desktop/etcd/etcd/README.md" http://127.0.0.1/api/files/product/upload/zall/fuck.md/sit?creator=zxjcli3 -v
			group.POST("/upload/:app/:name/:env", uploadProduct)
			// 下载制品
			group.GET("/get/:app/:name/:env", getProduct)
		}
		// 简单制品库
		group = e.Group("/api/product", apisession.CheckLogin)
		{
			// 制品库列表
			group.GET("/list", listProduct)
			// 删除制品
			group.DELETE("/delete/:productId", deleteProduct)
		}
	})
}

func deleteProduct(c *gin.Context) {
	err := filesrv.Outer.DeleteProduct(c, filesrv.DeleteProductReqDTO{
		ProductId: cast.ToInt64(c.Param("productId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listProduct(c *gin.Context) {
	products, err := filesrv.Outer.ListProduct(c, filesrv.ListProductReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(products, func(t filesrv.ProductDTO) (ProductVO, error) {
		return ProductVO{
			Id:      t.Id,
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

func checkProductToken(c *gin.Context) {
	if c.Query("t") != productToken {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.UnauthorizedCode.Int(),
			Message: "invalid normalToken",
		})
		c.Abort()
	}
}

func uploadAvatar(c *gin.Context) {
	body, b, err := ginutil.GetFile(c)
	if err != nil {
		logger.Logger.WithContext(c).Error(err)
		util.HandleApiErr(err, c)
		return
	}
	if b {
		defer body.Close()
	}
	path, err := filesrv.Outer.UploadAvatar(c, filesrv.UploadAvatarReqDTO{
		Body:     body,
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"filePath": path,
	})
}

func getAvatar(c *gin.Context) {
	path, err := filesrv.Outer.GetAvatar(c, filesrv.GetAvatarReqDTO{
		Name:     c.Param("name"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	if path == "" {
		c.String(http.StatusNotFound, "not found")
		return
	}
	// 设置浏览器缓存策略
	now := time.Now().Unix()
	expires := now + 31536000
	c.Header("Date", fmt.Sprintf("%d", now))
	c.Header("Expires", fmt.Sprintf("%d", expires))
	c.Header("Cache-Control", "public, max-age=31536000")
	c.File(path)
}

func uploadProduct(c *gin.Context) {
	body, b, err := ginutil.GetFile(c)
	if err != nil {
		logger.Logger.WithContext(c).Error(err)
		util.HandleApiErr(err, c)
		return
	}
	if b {
		defer body.Close()
	}
	path, err := filesrv.Outer.UploadProduct(c, filesrv.UploadProductReqDTO{
		AppId:   c.Param("app"),
		Name:    c.Param("name"),
		Creator: c.Query("creator"),
		Env:     c.Param("env"),
		Body:    body,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[string]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     path,
	})
}

func getProduct(c *gin.Context) {
	name := c.Param("name")
	path, err := filesrv.Outer.GetProduct(c, filesrv.GetProductReqDTO{
		AppId: c.Param("app"),
		Env:   c.Param("env"),
		Name:  name,
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	if path == "" {
		c.JSON(http.StatusNotFound, ginutil.BaseResp{
			Code:    apicode.DataNotExistsCode.Int(),
			Message: "file not found",
		})
		return
	}
	if c.Query("a") == "1" {
		c.Header("Content-Disposition", "attachment; filename=\""+name+"\"")
		c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	}
	c.File(path)
}
