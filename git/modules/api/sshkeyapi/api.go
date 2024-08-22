package sshkeyapi

import (
	"github.com/LeeZXin/zall/git/modules/service/sshkeysrv"
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
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/sshKey", apisession.CheckLogin)
		{
			// 删除
			group.DELETE("/delete/:keyId", deleteSshKey)
			// 插入
			group.POST("/create", createSshKey)
			// 列表展示
			group.GET("/list", listSshKey)
		}
	})
}

func createSshKey(c *gin.Context) {
	var req CreateSshKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := sshkeysrv.CreateSshKey(c, sshkeysrv.CreateSshKeyReqDTO{
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

func deleteSshKey(c *gin.Context) {
	err := sshkeysrv.DeleteSshKey(c, sshkeysrv.DeleteSshKeyReqDTO{
		Id:       cast.ToInt64(c.Param("keyId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listSshKey(c *gin.Context) {
	keys, err := sshkeysrv.ListSshKey(c, sshkeysrv.ListSshKeyReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(keys, func(t sshkeysrv.SshKeyDTO) (SshKeyVO, error) {
		return SshKeyVO{
			Id:           t.Id,
			Name:         t.Name,
			Fingerprint:  t.Fingerprint,
			Created:      t.Created.Format(time.DateOnly),
			LastOperated: t.LastOperated.Format(time.DateOnly),
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SshKeyVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
