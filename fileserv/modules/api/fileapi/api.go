package fileapi

import (
	"fmt"
	"github.com/LeeZXin/zall/fileserv/modules/service/filesrv"
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
	artifactToken string
)

func InitApi() {
	artifactToken = static.GetString("files.artifact.token")
	filesrv.InitStorage()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/files/avatar", apisession.CheckLogin)
		{
			// 上传头像
			group.POST("/upload", uploadAvatar)
			// 获取头像
			group.GET("/get/:name", getAvatar)
		}
		// 简单制品库
		group = e.Group("/api/files/artifact", checkArtifactToken)
		{
			// 上传制品
			// curl -F "file=@/Users/lizexin/Desktop/etcd/etcd/README.md" http://127.0.0.1/api/files/artifact/upload/zall/fuck.md/sit?creator=zxjcli3 -v
			group.POST("/upload/:app/:name/:env", uploadArtifact)
			// 下载制品
			group.GET("/get/:app/:name/:env", getArtifact)
		}
		// 简单制品库
		group = e.Group("/api/artifact", apisession.CheckLogin)
		{
			// 制品库列表
			group.GET("/list", listArtifact)
			// 删除制品
			group.DELETE("/delete/:artifactId", deleteArtifact)
		}
	})
}

func deleteArtifact(c *gin.Context) {
	err := filesrv.DeleteArtifact(c, filesrv.DeleteArtifactReqDTO{
		Id:       cast.ToInt64(c.Param("artifactId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listArtifact(c *gin.Context) {
	var req ListArtifactReqVO
	if util.ShouldBindQuery(&req, c) {
		artifacts, total, err := filesrv.ListArtifact(c, filesrv.ListArtifactReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(artifacts, func(t filesrv.ArtifactDTO) ArtifactVO {
			return ArtifactVO{
				Id:      t.Id,
				Name:    t.Name,
				Creator: t.Creator,
				Created: t.Created.Format(time.DateTime),
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[ArtifactVO]{
			DataResp: ginutil.DataResp[[]ArtifactVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func checkArtifactToken(c *gin.Context) {
	if c.Query("t") != artifactToken {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.UnauthorizedCode.Int(),
			Message: "invalid token",
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
	path, err := filesrv.UploadAvatar(c, filesrv.UploadAvatarReqDTO{
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
	path, err := filesrv.GetAvatar(c, filesrv.GetAvatarReqDTO{
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

func uploadArtifact(c *gin.Context) {
	body, b, err := ginutil.GetFile(c)
	if err != nil {
		logger.Logger.WithContext(c).Error(err)
		util.HandleApiErr(err, c)
		return
	}
	if b {
		defer body.Close()
	}
	path, err := filesrv.UploadArtifact(c, filesrv.UploadArtifactReqDTO{
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

func getArtifact(c *gin.Context) {
	name := c.Param("name")
	path, err := filesrv.GetArtifact(c, filesrv.GetArtifactReqDTO{
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
