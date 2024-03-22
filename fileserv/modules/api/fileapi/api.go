package fileapi

import (
	"github.com/LeeZXin/zall/fileserv/modules/service/filesrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

var (
	token string
)

func InitApi() {
	token = static.GetString("files.normal.token")
	filesrv.InitStorage()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/files/icon", apisession.CheckLogin)
		{
			group.POST("/upload/:name", uploadIcon)
			group.GET("/get/:id/:name", getIcon)
		}
		group = e.Group("/api/files/avatar", apisession.CheckLogin)
		{
			group.POST("/upload/:name", uploadAvatar)
			group.GET("/get/:id/:name", getAvatar)
		}
		group = e.Group("/api/files/normal", checkToken)
		{
			group.POST("/upload/:name", uploadNormal)
			group.GET("/get/:id/:name", getNormal)
		}
	})
}

func checkToken(c *gin.Context) {
	if c.Query("t") != token {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.UnauthorizedCode.Int(),
			Message: "invalid token",
		})
		c.Abort()
	}
}

func uploadIcon(c *gin.Context) {
	body, b, err := getBody(c)
	if err != nil {
		logger.Logger.WithContext(c).Error(err)
		util.HandleApiErr(err, c)
		return
	}
	if b {
		defer body.Close()
	}
	path, err := filesrv.Outer.UploadIcon(c, filesrv.UploadIconReqDTO{
		Name:     c.Param("name"),
		Body:     body,
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, UploadIconRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     path,
	})
}

func getIcon(c *gin.Context) {
	name := c.Param("name")
	path, err := filesrv.Outer.GetIcon(c, filesrv.GetIconReqDTO{
		Id:       c.Param("id"),
		Name:     name,
		Operator: apisession.MustGetLoginUser(c),
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

func uploadAvatar(c *gin.Context) {
	body, b, err := getBody(c)
	if err != nil {
		logger.Logger.WithContext(c).Error(err)
		util.HandleApiErr(err, c)
		return
	}
	if b {
		defer body.Close()
	}
	path, err := filesrv.Outer.UploadAvatar(c, filesrv.UploadAvatarReqDTO{
		Name:     c.Param("name"),
		Body:     body,
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, UploadIconRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     path,
	})
}

func getAvatar(c *gin.Context) {
	name := c.Param("name")
	path, err := filesrv.Outer.GetAvatar(c, filesrv.GetAvatarReqDTO{
		Id:       c.Param("id"),
		Name:     name,
		Operator: apisession.MustGetLoginUser(c),
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

func uploadNormal(c *gin.Context) {
	body, b, err := getBody(c)
	if err != nil {
		logger.Logger.WithContext(c).Error(err)
		util.HandleApiErr(err, c)
		return
	}
	if b {
		defer body.Close()
	}
	path, err := filesrv.Outer.UploadNormal(c, filesrv.UploadNormalReqDTO{
		Name: c.Param("name"),
		Body: body,
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, UploadIconRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     path,
	})
}

func getNormal(c *gin.Context) {
	name := c.Param("name")
	path, err := filesrv.Outer.GetNormal(c, filesrv.GetNormalReqDTO{
		Id:   c.Param("id"),
		Name: name,
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

func getBody(c *gin.Context) (io.ReadCloser, bool, error) {
	contentType := strings.ToLower(c.GetHeader("Content-Type"))
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") || strings.HasPrefix(contentType, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			return nil, false, err
		}
		if c.Request.MultipartForm.File == nil {
			return nil, false, http.ErrMissingFile
		}
		for _, files := range c.Request.MultipartForm.File {
			if len(files) > 0 {
				r, err := files[0].Open()
				return r, true, err
			}
		}
		return nil, false, http.ErrMissingFile
	}
	return c.Request.Body, false, nil
}
