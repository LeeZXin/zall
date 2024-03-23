package sshkeyapi

import (
	"github.com/LeeZXin/zall/git/modules/model/sshkeymd"
	"github.com/LeeZXin/zall/git/modules/service/sshkeysrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/sshKey", apisession.CheckLogin)
		{
			// 删除
			group.POST("/delete", deleteSshKey)
			// 插入
			group.POST("/insert", insertSshKey)
			// 列表展示
			group.POST("/list", listSshKey)
			// 获取校验token
			group.POST("/getToken", getToken)
		}
	})
}

func insertSshKey(c *gin.Context) {
	var req InsertSshKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := sshkeysrv.Outer.InsertSshKey(c, sshkeysrv.InsertSshKeyReqDTO{
			Name:          req.Name,
			PubKeyContent: req.PubKeyContent,
			Operator:      apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteSshKey(c *gin.Context) {
	var req DeleteSshKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := sshkeysrv.Outer.DeleteSshKey(c, sshkeysrv.DeleteSshKeyReqDTO{
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

func listSshKey(c *gin.Context) {
	respDTO, err := sshkeysrv.Outer.ListSshKey(c, sshkeysrv.ListSshKeyReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(respDTO, func(t sshkeymd.SshKey) (SshKeyVO, error) {
		return SshKeyVO{
			Id:          t.Id,
			Name:        t.Name,
			Fingerprint: t.Fingerprint,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SshKeyVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func getToken(c *gin.Context) {
	var req GetTokenReqVO
	if util.ShouldBindJSON(&req, c) {
		token, guides, err := sshkeysrv.Outer.GetToken(c, sshkeysrv.GetTokenReqDTO{
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		tokenVO := TokenVO{
			Token:  token,
			Guides: guides,
		}
		c.JSON(http.StatusOK, ginutil.DataResp[TokenVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     tokenVO,
		})
	}
}
