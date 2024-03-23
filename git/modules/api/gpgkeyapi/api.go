package gpgkeyapi

import (
	"github.com/LeeZXin/zall/git/modules/service/gpgkeysrv"
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
		group := e.Group("/api/gpgKey", apisession.CheckLogin)
		{
			group.GET("/list", listGpgKey)
			group.POST("/getToken", getToken)
			group.POST("/insert", insertGpgKey)
			group.POST("/delete", deleteGpgKey)
		}
	})
}

func getToken(c *gin.Context) {
	var req GetTokenReqVO
	if util.ShouldBindJSON(&req, c) {
		token, guides, err := gpgkeysrv.Outer.GetToken(c, gpgkeysrv.GetTokenReqDTO{
			Content:  req.Content,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, GetTokenRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Token:    token,
			Guides:   guides,
		})
	}
}

func insertGpgKey(c *gin.Context) {
	var req InsertGpgKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := gpgkeysrv.Outer.InsertGpgKey(c, gpgkeysrv.InsertGpgKeyReqDTO{
			Name:      req.Name,
			Content:   req.Content,
			Signature: req.Signature,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteGpgKey(c *gin.Context) {
	var req DeleteGpgKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := gpgkeysrv.Outer.DeleteGpgKey(c, gpgkeysrv.DeleteGpgKeyReqDTO{
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

func listGpgKey(c *gin.Context) {
	dtoList, err := gpgkeysrv.Outer.ListGpgKey(c, gpgkeysrv.ListGpgKeyReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(dtoList, func(t gpgkeysrv.GpgKeyDTO) (GpgKeyVO, error) {
		return GpgKeyVO{
			Id:         t.Id,
			Name:       t.Name,
			PubKeyId:   t.PubKeyId,
			ExpireTime: t.ExpireTime.Format(time.DateTime),
			EmailList:  t.EmailList,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]GpgKeyVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
