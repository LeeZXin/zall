package serviceapi

import (
	"github.com/LeeZXin/zall/deploy/modules/service/servicesrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func InitApi() {
	servicesrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/service", apisession.CheckLogin)
		{
			// 创建探针
			group.POST("/create", createService)
			// 编辑探针
			group.POST("/update", updateService)
			// 启动探针
			group.PUT("/enable/:serviceId", enableService)
			// 关闭探针
			group.PUT("/disable/:serviceId", disableService)
			// 删除探针
			group.DELETE("/delete/:serviceId", deleteService)
			// 探针列表
			group.GET("/list", listService)
		}
	})
}

func createService(c *gin.Context) {
	var req CreateServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := servicesrv.Outer.CreateService(c, servicesrv.CreateServiceReqDTO{
			AppId:    req.AppId,
			Name:     req.Name,
			Config:   req.Config,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateService(c *gin.Context) {
	var req UpdateServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := servicesrv.Outer.UpdateService(c, servicesrv.UpdateServiceReqDTO{
			ServiceId: req.ServiceId,
			Name:      req.Name,
			Config:    req.Config,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func enableService(c *gin.Context) {
	err := servicesrv.Outer.EnableService(c, servicesrv.EnableServiceReqDTO{
		ServiceId: cast.ToInt64(c.Param("serviceId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func disableService(c *gin.Context) {
	err := servicesrv.Outer.DisableService(c, servicesrv.DisableServiceReqDTO{
		ServiceId: cast.ToInt64(c.Param("serviceId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func deleteService(c *gin.Context) {
	err := servicesrv.Outer.DeleteService(c, servicesrv.DeleteServiceReqDTO{
		ServiceId: cast.ToInt64(c.Param("serviceId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listService(c *gin.Context) {
	var req ListServiceReqVO
	if util.ShouldBindQuery(&req, c) {
		probes, total, err := servicesrv.Outer.ListService(c, servicesrv.ListServiceReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(probes, func(t servicesrv.ServiceDTO) (ServiceVO, error) {
			return ServiceVO{
				Id:        t.Id,
				AppId:     t.AppId,
				Config:    t.Config,
				Env:       t.Env,
				IsEnabled: t.IsEnabled,
				Name:      t.Name,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[ServiceVO]{
			DataResp: ginutil.DataResp[[]ServiceVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}

}
