package gpgkeyapi

import (
	"github.com/LeeZXin/zall/git/modules/service/gpgkeysrv"
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
	gpgkeysrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/gpgKey", apisession.CheckLogin)
		{
			// gpg密钥列表
			group.GET("/list", listGpgKey)
			// 新增gpg密钥
			group.POST("/create", createGpgKey)
			// 删除gpg密钥
			group.DELETE("/delete/:keyId", deleteGpgKey)
		}
	})
}

func createGpgKey(c *gin.Context) {
	var req CreateGpgKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := gpgkeysrv.Outer.CreateGpgKey(c, gpgkeysrv.CreateGpgKeyReqDTO{
			Name:     req.Name,
			Content:  req.Content,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteGpgKey(c *gin.Context) {
	err := gpgkeysrv.Outer.DeleteGpgKey(c, gpgkeysrv.DeleteGpgKeyReqDTO{
		Id:       cast.ToInt64(c.Param("keyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listGpgKey(c *gin.Context) {
	keys, err := gpgkeysrv.Outer.ListGpgKey(c, gpgkeysrv.ListGpgKeyReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(keys, func(t gpgkeysrv.GpgKeyDTO) (GpgKeyVO, error) {
		return GpgKeyVO{
			Id:      t.Id,
			Name:    t.Name,
			KeyId:   t.KeyId,
			Expired: t.Expired.Format(time.DateOnly),
			Created: t.Created.Format(time.DateOnly),
			Email:   t.Email,
			SubKeys: t.SubKeys,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]GpgKeyVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
