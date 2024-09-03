package tpfeishuapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/thirdpart/modules/service/tpfeishusrv"
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
		group := e.Group("/api/feishuAccessToken")
		{
			// token列表
			group.GET("/list", apisession.CheckLogin, listAccessToken)
			// 创建token
			group.POST("/create", apisession.CheckLogin, createAccessToken)
			// 创建token
			group.POST("/update", apisession.CheckLogin, updateAccessToken)
			// 删除
			group.DELETE("/delete/:tokenId", apisession.CheckLogin, deleteAccessToken)
			// 重刷
			group.PUT("/refresh/:tokenId", apisession.CheckLogin, refreshAccessToken)
			// 变更apikey
			group.PUT("/changeApiKey/:tokenId", apisession.CheckLogin, changeAccessTokenApiKey)
			// 通过apikey获取token
			group.GET("/getByApiKey/:apiKey", getAccessTokenByApiKey)
		}
	})
}

func getAccessTokenByApiKey(c *gin.Context) {
	token, tenantToken, err := tpfeishusrv.GetAccessTokenByApiKey(c, c.Param("apiKey"))
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     []string{token, tenantToken},
	})
}

func changeAccessTokenApiKey(c *gin.Context) {
	err := tpfeishusrv.ChangeAccessTokenApiKey(c, tpfeishusrv.ChangeAccessTokenApiKeyReqDTO{
		Id:       cast.ToInt64(c.Param("tokenId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func refreshAccessToken(c *gin.Context) {
	err := tpfeishusrv.RefreshAccessToken(c, tpfeishusrv.RefreshAccessTokenReqDTO{
		Id:       cast.ToInt64(c.Param("tokenId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func deleteAccessToken(c *gin.Context) {
	err := tpfeishusrv.DeleteAccessToken(c, tpfeishusrv.DeleteAccessTokenReqDTO{
		Id:       cast.ToInt64(c.Param("tokenId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func createAccessToken(c *gin.Context) {
	var req CreateAccessTokenReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tpfeishusrv.CreateAccessToken(c, tpfeishusrv.CreateAccessTokenReqDTO{
			TeamId:   req.TeamId,
			Name:     req.Name,
			AppId:    req.AppId,
			Secret:   req.Secret,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateAccessToken(c *gin.Context) {
	var req UpdateAccessTokenReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tpfeishusrv.UpdateAccessToken(c, tpfeishusrv.UpdateAccessTokenReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			AppId:    req.AppId,
			Secret:   req.Secret,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listAccessToken(c *gin.Context) {
	var req ListWeworkAccessTokenReqVO
	if util.ShouldBindQuery(&req, c) {
		tokens, total, err := tpfeishusrv.ListAccessToken(c, tpfeishusrv.ListAccessTokenReqDTO{
			TeamId:   req.TeamId,
			Key:      req.Key,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(tokens, func(t tpfeishusrv.AccessTokenDTO) AccessTokenVO {
			return AccessTokenVO{
				Id:          t.Id,
				TeamId:      t.TeamId,
				Name:        t.Name,
				AppId:       t.AppId,
				Creator:     t.Creator,
				Secret:      t.Secret,
				Token:       t.Token,
				TenantToken: t.TenantToken,
				ApiKey:      t.ApiKey,
				Expired:     t.Expired.Format(time.DateTime),
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[AccessTokenVO]{
			DataResp: ginutil.DataResp[[]AccessTokenVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}
