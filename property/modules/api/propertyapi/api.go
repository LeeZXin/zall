package propertyapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/property/modules/service/propertysrv"
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
	cfgsrv.Init()
	propertysrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/propertySource", apisession.CheckLogin)
		{
			// 新增配置来源
			group.POST("/create", createPropertySource)
			// 配置来源列表
			group.GET("/list", listPropertySource)
			// 编辑配置列表
			group.POST("/update", updatePropertySource)
			// 删除配置来源
			group.DELETE("/delete/:sourceId", deletePropertySource)
			// 所有配置来源列表
			group.GET("/listAll/:env", listAllPropertySource)
			// 获取应用服务绑定的配置来源
			group.GET("/listBind", listBindPropertySource)
			// 绑定应用服务和配置来源
			group.POST("/bindApp", bindAppAndPropertySource)
		}
		group = e.Group("/api/propertyFile", apisession.CheckLogin)
		{
			// 创建配置文件
			group.POST("/create", createFile)
			// 配置文件列表
			group.GET("/list", listFile)
			// 展示来源列表
			group.GET("/listSource/:fileId", listPropertySourceByFileId)
			// 删除文件
			group.DELETE("/delete/:fileId", deleteFile)
		}
		group = e.Group("/api/propertyHistory", apisession.CheckLogin)
		{
			// 获取最新版本的配置
			group.GET("/getByVersion", getHistoryByVersion)
			// 版本列表
			group.GET("/list", pageHistory)
			// 新增版本
			group.POST("/newVersion", newVersion)
			// 发布配置
			group.POST("/deploy", deployHistory)
			// 查看发布记录
			group.GET("/listDeploy/:historyId", listDeploy)
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

func createPropertySource(c *gin.Context) {
	var req CreatePropertySourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.CreatePropertySource(c, propertysrv.CreatePropertySourceReqDTO{
			Endpoints: req.Endpoints,
			Username:  req.Username,
			Password:  req.Password,
			Env:       req.Env,
			Name:      req.Name,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updatePropertySource(c *gin.Context) {
	var req UpdatePropertySourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.UpdatePropertySource(c, propertysrv.UpdatePropertySourceReqDTO{
			SourceId:  req.SourceId,
			Name:      req.Name,
			Endpoints: req.Endpoints,
			Username:  req.Username,
			Password:  req.Password,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deletePropertySource(c *gin.Context) {
	err := propertysrv.Outer.DeletePropertySource(c, propertysrv.DeletePropertySourceReqDTO{
		SourceId: cast.ToInt64(c.Param("sourceId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listPropertySource(c *gin.Context) {
	nodes, err := propertysrv.Outer.ListPropertySource(c, propertysrv.ListPropertySourceReqDTO{
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(nodes, func(t propertysrv.PropertySourceDTO) (PropertySourceVO, error) {
		return PropertySourceVO{
			Id:        t.Id,
			Name:      t.Name,
			Endpoints: t.Endpoints,
			Username:  t.Username,
			Password:  t.Password,
			Env:       t.Env,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]PropertySourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listAllPropertySource(c *gin.Context) {
	nodes, err := propertysrv.Outer.ListAllPropertySource(c, propertysrv.ListAllPropertySourceReqDTO{
		Env:      c.Param("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimplePropertySourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     sourceDto2SimpleVo(nodes),
	})
}

func sourceDto2SimpleVo(sources []propertysrv.SimplePropertySourceDTO) []SimplePropertySourceVO {
	data, _ := listutil.Map(sources, func(t propertysrv.SimplePropertySourceDTO) (SimplePropertySourceVO, error) {
		return SimplePropertySourceVO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
	return data
}

func listBindPropertySource(c *gin.Context) {
	nodes, err := propertysrv.Outer.ListBindPropertySource(c, propertysrv.ListBindPropertySourceReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimplePropertySourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     sourceDto2SimpleVo(nodes),
	})
}

func bindAppAndPropertySource(c *gin.Context) {
	var req BindAppAndPropertySourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.BindAppAndPropertySource(c, propertysrv.BindAppAndPropertySourceReqDTO{
			AppId:        req.AppId,
			SourceIdList: req.SourceIdList,
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

func listPropertySourceByFileId(c *gin.Context) {
	nodes, err := propertysrv.Outer.ListPropertySourceByFileId(c, propertysrv.ListPropertySourceByFileIdReqDTO{
		FileId:   cast.ToInt64(c.Param("fileId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimplePropertySourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     sourceDto2SimpleVo(nodes),
	})
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

func deleteFile(c *gin.Context) {
	err := propertysrv.Outer.DeleteFile(c, propertysrv.DeleteFileReqDTO{
		FileId:   cast.ToInt64(c.Param("fileId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func deployHistory(c *gin.Context) {
	var req DeployHistoryReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propertysrv.Outer.DeployHistory(c, propertysrv.DeployHistoryReqDTO{
			HistoryId:    req.HistoryId,
			SourceIdList: req.SourceIdList,
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
		FileName:    t.FileName,
		FileId:      t.FileId,
		Content:     t.Content,
		Version:     t.Version,
		Created:     t.Created.Format(time.DateTime),
		Creator:     t.Creator,
		LastVersion: t.LastVersion,
		Env:         t.Env,
	}
}

func listDeploy(c *gin.Context) {
	deploys, err := propertysrv.Outer.ListDeploy(c, propertysrv.ListDeployReqDTO{
		HistoryId: cast.ToInt64(c.Param("historyId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(deploys, func(t propertysrv.DeployDTO) (DeployVO, error) {
		return DeployVO{
			NodeName:  t.NodeName,
			Endpoints: t.Endpoints,
			Created:   t.Created.Format(time.DateTime),
			Creator:   t.Creator,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]DeployVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
