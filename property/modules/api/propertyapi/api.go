package propertyapi

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/property/modules/service/propertysrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/http/httptask"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"net/url"
	"time"
)

func InitApi() {
	cfgsrv.Init()
	propertysrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/prop/etcdNode", apisession.CheckLogin)
		{
			group.POST("/insert", insertEtcdNode)
			group.GET("/list", listEtcdNode)
			group.GET("/listSimple", listSimpleEtcdNode)
			group.POST("/update", updateEtcdNode)
			group.POST("/delete", deleteEtcdNode)
			group.POST("/getAuth", getAuth)
			group.POST("/grantAuth", grantAuth)
		}
		group = e.Group("/api/propertyFile", apisession.CheckLogin)
		{
			// 创建配置文件
			group.POST("/create", createFile)
			// 配置文件列表
			group.GET("/list", listFile)
			group.POST("/delete", deleteContent)
			group.POST("/deploy", deployContent)
			group.POST("/listDeploy", listDeploy)
		}
		group = e.Group("/api/propertyHistory", apisession.CheckLogin)
		{
			// 获取最新版本的配置
			group.GET("/getByVersion", getHistoryByVersion)
			// 版本列表
			group.GET("/list", pageHistory)
			// 新增版本
			group.POST("/newVersion", newVersion)
		}
	})
	httptask.AppendHttpTask("checkPropDbEtcdConsistent", func(_ []byte, _ url.Values) {
		envs, b := cfgsrv.Inner.GetEnvCfg(context.Background())
		if b {
			for _, env := range envs {
				propertysrv.Inner.CheckConsistent(env)
			}
		}
	})
}

func getHistoryByVersion(c *gin.Context) {
	history, exist, err := propertysrv.Outer.GetHistoryByVersion(c, propertysrv.GetHistoryByVersionReqDTO{
		FileId:   cast.ToInt64(c.Query("fileId")),
		Version:  c.Query("version"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[util.ValueWithExist[HistoryVO]]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: util.ValueWithExist[HistoryVO]{
			Exist: exist,
			Value: history2VO(history),
		},
	})
}

func insertEtcdNode(c *gin.Context) {
	var req InsertEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.InsertEtcdNode(c, propertysrv.InsertEtcdNodeReqDTO{
			NodeId:    req.NodeId,
			Endpoints: req.Endpoints,
			Username:  req.Username,
			Password:  req.Password,
			Env:       req.Env,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateEtcdNode(c *gin.Context) {
	var req UpdateEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.UpdateEtcdNode(c, propertysrv.UpdateEtcdNodeReqDTO{
			NodeId:    req.NodeId,
			Endpoints: req.Endpoints,
			Username:  req.Username,
			Password:  req.Password,
			Env:       req.Env,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteEtcdNode(c *gin.Context) {
	var req DeleteEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.DeleteEtcdNode(c, propertysrv.DeleteEtcdNodeReqDTO{
			NodeId:   req.NodeId,
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

func listEtcdNode(c *gin.Context) {
	var req ListEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		nodes, err := propertysrv.Outer.ListEtcdNode(c, propertysrv.ListEtcdNodeReqDTO{
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(nodes, func(t propertysrv.EtcdNodeDTO) (EtcdNodeVO, error) {
			return EtcdNodeVO{
				NodeId:    t.NodeId,
				Endpoints: t.Endpoints,
				Username:  t.Username,
				Password:  t.Password,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]EtcdNodeVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func listSimpleEtcdNode(c *gin.Context) {
	var req ListSimpleEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		nodes, err := propertysrv.Outer.ListSimpleEtcdNode(c, req.Env)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     nodes,
		})
	}
}

func createFile(c *gin.Context) {
	var req CreateFileReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.CreateFile(c, propertysrv.CreateFileReqDTO{
			AppId:    req.AppId,
			Name:     req.Name,
			Content:  req.Content,
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

func newVersion(c *gin.Context) {
	var req NewVersionReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.NewVersion(c, propertysrv.NewVersionReqDTO{
			FileId:      req.FileId,
			Content:     req.Content,
			LastVersion: req.LastVersion,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteContent(c *gin.Context) {
	var req DeleteContentReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.DeletePropContent(c, propertysrv.DeletePropContentReqDTO{
			Id:       req.Id,
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

func deployContent(c *gin.Context) {
	var req DeployContentReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.DeployPropContent(c, propertysrv.DeployPropContentReqDTO{
			Id:           req.Id,
			Version:      req.Version,
			EtcdNodeList: req.EtcdNodeList,
			Env:          req.Env,
			Operator:     apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listFile(c *gin.Context) {
	contents, err := propertysrv.Outer.ListFile(c, propertysrv.ListFileReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(contents, func(t propertysrv.FileDTO) (FileVO, error) {
		return FileVO{
			Id:    t.Id,
			AppId: t.AppId,
			Name:  t.Name,
			Env:   t.Env,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]FileVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func pageHistory(c *gin.Context) {
	var req PageHistoryReqVO
	if util.ShouldBindQuery(&req, c) {
		histories, cursor, err := propertysrv.Outer.PageHistory(c, propertysrv.PageHistoryReqDTO{
			FileId:   req.FileId,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(histories, func(t propertysrv.HistoryDTO) (HistoryVO, error) {
			return history2VO(t), nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[HistoryVO]{
			DataResp: ginutil.DataResp[[]HistoryVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: cursor,
		})
	}
}

func history2VO(t propertysrv.HistoryDTO) HistoryVO {
	return HistoryVO{
		Id:          t.Id,
		FileId:      t.FileId,
		Content:     t.Content,
		Version:     t.Version,
		Created:     t.Created.Format(time.DateTime),
		Creator:     t.Creator,
		LastVersion: t.LastVersion,
	}
}

func grantAuth(c *gin.Context) {
	var req GrantAuthReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.GrantAuth(c, propertysrv.GrantAuthReqDTO{
			AppId:    req.AppId,
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

func getAuth(c *gin.Context) {
	var req GetAuthReqVO
	if util.ShouldBindJSON(&req, c) {
		username, password, err := propertysrv.Outer.GetAuth(c, propertysrv.GetAuthReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, GetAuthRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Username: username,
			Password: password,
		})
	}
}

func listDeploy(c *gin.Context) {
	var req ListDeployReqVO
	if util.ShouldBindJSON(&req, c) {
		deploys, cursor, err := propertysrv.Outer.ListDeploy(c, propertysrv.ListDeployReqDTO{
			ContentId: req.ContentId,
			Version:   req.Version,
			Cursor:    req.Cursor,
			Limit:     req.Limit,
			NodeId:    req.NodeId,
			Env:       req.Env,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(deploys, func(t propertysrv.DeployDTO) (DeployVO, error) {
			return DeployVO{
				ContentId: t.ContentId,
				Content:   t.Content,
				Version:   t.Version,
				NodeId:    t.NodeId,
				Created:   t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[DeployVO]{
			DataResp: ginutil.DataResp[[]DeployVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: cursor,
		})
	}
}
