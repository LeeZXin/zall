package propsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/prop/modules/model/propmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/strutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type outerImpl struct {
}

func (*outerImpl) ListSimpleEtcdNode(ctx context.Context) ([]string, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	nodes, err := propmd.ListEtcdNode(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(nodes, func(t propmd.EtcdNode) (string, error) {
		return t.NodeId, nil
	})
}

func (*outerImpl) ListEtcdNode(ctx context.Context, reqDTO ListEtcdNodeReqDTO) ([]EtcdNodeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	// 系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	nodes, err := propmd.ListEtcdNode(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(nodes, func(t propmd.EtcdNode) (EtcdNodeDTO, error) {
		return EtcdNodeDTO{
			NodeId:    t.NodeId,
			Endpoints: strings.Split(t.Endpoints, ";"),
			Username:  t.Username,
			Password:  t.Password,
		}, nil
	})
}

func (*outerImpl) InsertEtcdNode(ctx context.Context, reqDTO InsertEtcdNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.InsertEtcdNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, b, err := propmd.GetEtcdNodeByNodeId(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.AlreadyExistsError()
		return
	}
	err = propmd.InsertEtcdNode(ctx, propmd.InsertEtcdNodeReqDTO{
		NodeId:    reqDTO.NodeId,
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
	})
	if err != nil {
		if xormutil.IsDuplicatedEntryError(err) {
			err = util.AlreadyExistsError()
			return
		}
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteEtcdNode(ctx context.Context, reqDTO DeleteEtcdNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.DeleteEtcdNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	// 系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	err = propmd.DeleteEtcdNode(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) UpdateEtcdNode(ctx context.Context, reqDTO UpdateEtcdNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.UpdateEtcdNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	// 系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	_, err = propmd.UpdateEtcdNode(ctx, propmd.UpdateEtcdNodeReqDTO{
		NodeId:    reqDTO.NodeId,
		Endpoints: strings.Join(reqDTO.Endpoints, ";"),
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) GrantAuth(ctx context.Context, reqDTO GrantAuthReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.GrantAuth),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err = checkPropContentPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return
	}
	var (
		auth propmd.EtcdAuth
		b    bool
	)
	auth, b, err = propmd.GetAuthByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		insertReq := propmd.InsertAuthReqDTO{
			AppId:    reqDTO.AppId,
			Username: reqDTO.AppId,
			Password: strutil.RandomStr(16),
		}
		err = propmd.InsertAuth(ctx, insertReq)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		auth = propmd.EtcdAuth{
			AppId:    insertReq.AppId,
			Username: insertReq.Username,
			Password: insertReq.Password,
		}
	}
	go grantAuthToEtcd(auth)
	return
}

func (*outerImpl) InsertPropContent(ctx context.Context, reqDTO InsertPropContentReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.InsertPropContent),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err = checkPropContentPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return
	}
	var b bool
	_, b, err = propmd.GetPropContentByAppIdAndName(ctx, reqDTO.AppId, reqDTO.Name)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.AlreadyExistsError()
		return
	}
	err = mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		content, err := propmd.InsertPropContent(ctx, propmd.InsertPropContentReqDTO{
			AppId: reqDTO.AppId,
			Name:  reqDTO.Name,
		})
		if err != nil {
			return err
		}
		return propmd.InsertHistory(ctx, propmd.InsertHistoryReqDTO{
			ContentId: content.Id,
			Content:   reqDTO.Content,
			Version:   genVersion(),
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) UpdatePropContent(ctx context.Context, reqDTO UpdatePropContentReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.UpdatePropContent),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if _, err = checkPropContentPerm(ctx, reqDTO.Operator, reqDTO.Id); err != nil {
		return
	}
	err = propmd.InsertHistory(ctx, propmd.InsertHistoryReqDTO{
		ContentId: reqDTO.Id,
		Content:   reqDTO.Content,
		Version:   genVersion(),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeletePropContent(ctx context.Context, reqDTO DeletePropContentReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.DeletePropContent),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	var content propmd.PropContent
	if content, err = checkPropContentPerm(ctx, reqDTO.Operator, reqDTO.Id); err != nil {
		return
	}
	err = mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		_, err := propmd.DeletePropContent(ctx, reqDTO.Id)
		if err != nil {
			return err
		}
		return propmd.DeleteHistory(ctx, reqDTO.Id)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	go deleteFromEtcd(content)
	return
}

func (*outerImpl) ListPropContent(ctx context.Context, reqDTO ListPropContentReqDTO) ([]PropContentDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err := checkPropContentPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return nil, err
	}
	contents, err := propmd.ListPropContent(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(contents, func(t propmd.PropContent) (PropContentDTO, error) {
		return PropContentDTO{
			Id:    t.Id,
			AppId: t.AppId,
			Name:  t.Name,
		}, nil
	})
}

func (*outerImpl) DeployPropContent(ctx context.Context, reqDTO DeployPropContentReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PropSrvKeysVO.DeployPropContent),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if _, err = checkPropContentPerm(ctx, reqDTO.Operator, reqDTO.Id); err != nil {
		return
	}
	nodes, err := propmd.BatchGetEtcdNodes(ctx, reqDTO.EtcdNodeList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if len(nodes) == 0 {
		return util.InvalidArgsError()
	}
	// 获取历史版本记录
	var (
		history propmd.PropHistory
		b       bool
	)
	history, b, err = propmd.GetHistoryByVersion(ctx, reqDTO.Id, reqDTO.Version)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	go func() {
		for _, node := range nodes {
			deployToEtcd(reqDTO.Id, history.Content, history.Version, node)
		}
	}()
	return nil
}

func (*outerImpl) ListHistory(ctx context.Context, reqDTO ListHistoryReqDTO) ([]HistoryDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if _, err := checkPropContentPerm(ctx, reqDTO.Operator, reqDTO.ContentId); err != nil {
		return nil, 0, err
	}
	histories, err := propmd.ListHistory(ctx, propmd.ListHistoryReqDTO{
		ContentId: reqDTO.ContentId,
		Version:   reqDTO.Version,
		Cursor:    reqDTO.Cursor,
		Limit:     reqDTO.Limit,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var cursor int64 = 0
	if len(histories) == reqDTO.Limit {
		cursor = histories[len(histories)-1].Id
	}
	ret, _ := listutil.Map(histories, func(t propmd.PropHistory) (HistoryDTO, error) {
		return HistoryDTO{
			ContentId: t.ContentId,
			Content:   t.Content,
			Version:   t.Version,
			Created:   t.Created,
		}, nil
	})
	return ret, cursor, nil
}

func (*outerImpl) ListDeploy(ctx context.Context, reqDTO ListDeployReqDTO) ([]DeployDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if _, err := checkPropContentPerm(ctx, reqDTO.Operator, reqDTO.ContentId); err != nil {
		return nil, 0, err
	}
	deploys, err := propmd.ListDeploy(ctx, propmd.ListDeployReqDTO{
		ContentId: reqDTO.ContentId,
		Version:   reqDTO.Version,
		Cursor:    reqDTO.Cursor,
		Limit:     reqDTO.Limit,
		NodeId:    reqDTO.NodeId,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var cursor int64 = 0
	if len(deploys) == reqDTO.Limit {
		cursor = deploys[len(deploys)-1].Id
	}
	ret, _ := listutil.Map(deploys, func(t propmd.PropDeploy) (DeployDTO, error) {
		return DeployDTO{
			ContentId: t.ContentId,
			Content:   t.Content,
			Version:   t.Version,
			NodeId:    t.NodeId,
			Created:   t.Created,
		}, nil
	})
	return ret, cursor, nil
}

func deleteFromEtcd(content propmd.PropContent) {
	ctx, closer := mysqlstore.Context(context.Background())
	defer closer.Close()
	_, err := propmd.ListEtcdNode(ctx)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	// todo 操作etcd
}

func grantAuthToEtcd(auth propmd.EtcdAuth) {
	//todo 操作etcd
}

func deployToEtcd(id int64, content, version string, node propmd.EtcdNode) {
	ctx, closer := mysqlstore.Context(context.Background())
	defer closer.Close()
	err := mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		err := propmd.InsertDeploy(ctx, propmd.InsertDeployReqDTO{
			ContentId: id,
			Content:   content,
			Version:   version,
			NodeId:    node.NodeId,
			Endpoints: node.Endpoints,
			Username:  node.Username,
			Password:  node.Password,
		})
		if err != nil {
			return err
		}
		// todo 操作etcd
		return nil
	})
	if err != nil {
		logger.Logger.Error(err)
	}
}

func genVersion() string {
	now := time.Now()
	rint := strconv.Itoa(rand.Intn(1000000))
	if len(rint) < 6 {
		rint = "000000" + rint
		rint = rint[len(rint)-6:]
	} else if len(rint) > 6 {
		rint = rint[len(rint)-6:]
	}
	return now.Format("20060102150405") + rint
}

func checkPropContentPerm(ctx context.Context, operator apisession.UserInfo, id int64) (propmd.PropContent, error) {
	content, b, err := propmd.GetPropContentById(ctx, id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return propmd.PropContent{}, util.InternalError(err)
	}
	if !b {
		return propmd.PropContent{}, util.InvalidArgsError()
	}
	return content, checkPropContentPermByAppId(ctx, operator, content.AppId)
}

func checkPropContentPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, app.TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if !p.IsAdmin && !p.PermDetail.GetAppPerm(appId).CanHandleProp {
		return util.UnauthorizedError()
	}
	return nil
}
