package mysqldbsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/dbaudit/modules/model/mysqldbmd"
	"github.com/LeeZXin/zall/dbaudit/modules/service/mysqldbsrv/command"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type outerImpl struct{}

func (*outerImpl) InsertDb(ctx context.Context, reqDTO InsertDbReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DbSrvKeysVO.InsertDb),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if err = checkPerm(reqDTO.Operator); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err = mysqldbmd.InsertDb(ctx, mysqldbmd.InsertDbReqDTO{
		Name:     reqDTO.Name,
		DbHost:   reqDTO.DbHost,
		Username: reqDTO.Username,
		Password: reqDTO.Password,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) UpdateDb(ctx context.Context, reqDTO UpdateDbReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DbSrvKeysVO.UpdateDb),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if err = checkPerm(reqDTO.Operator); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = mysqldbmd.UpdateDb(ctx, mysqldbmd.UpdateDbReqDTO{
		Id:       reqDTO.Id,
		Name:     reqDTO.Name,
		DbHost:   reqDTO.DbHost,
		Username: reqDTO.Username,
		Password: reqDTO.Password,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteDb(ctx context.Context, reqDTO DeleteDbReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DbSrvKeysVO.DeleteDb),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if err = checkPerm(reqDTO.Operator); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = mysqldbmd.DeleteDbById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListDb(ctx context.Context, reqDTO ListDbReqDTO) ([]DbDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	dbs, err := mysqldbmd.ListDb(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(dbs, func(t mysqldbmd.Db) (DbDTO, error) {
		return DbDTO{
			Id:       t.Id,
			Name:     t.Name,
			DbHost:   t.DbHost,
			Username: t.Username,
			Password: t.Password,
			Created:  t.Created,
		}, nil
	})
}

func (*outerImpl) ListSimpleDb(ctx context.Context, reqDTO ListSimpleDbReqDTO) ([]SimpleDbDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	dbs, err := mysqldbmd.ListDb(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(dbs, func(t mysqldbmd.Db) (SimpleDbDTO, error) {
		return SimpleDbDTO{
			Id:     t.Id,
			Name:   t.Name,
			DbHost: t.DbHost,
		}, nil
	})
}

// ApplyDbPerm 申请库表权限
func (*outerImpl) ApplyDbPerm(ctx context.Context, reqDTO ApplyDbPermReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	err = mysqldbmd.InsertPermApprovalOrder(ctx, mysqldbmd.InsertPermApprovalOrderReqDTO{
		Account:      reqDTO.Operator.Account,
		DbId:         reqDTO.DbId,
		AccessBase:   reqDTO.AccessBase,
		AccessTables: reqDTO.AccessTables,
		PermType:     reqDTO.PermType,
		OrderStatus:  mysqldbmd.PendingPermOrderStatus,
		ExpireDay:    reqDTO.ExpireDay,
		Reason:       reqDTO.Reason,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListPermApprovalOrder 展示审批列表
func (*outerImpl) ListPermApprovalOrder(ctx context.Context, reqDTO ListPermApprovalOrderReqDTO) ([]PermApprovalOrderDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	orders, err := mysqldbmd.ListPermApprovalOrder(ctx, mysqldbmd.ListPermApprovalOrderReqDTO{
		Cursor:      reqDTO.Cursor,
		Limit:       reqDTO.Limit,
		OrderStatus: reqDTO.OrderStatus,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(orders) == reqDTO.Limit {
		next = orders[len(orders)-1].Id
	}
	dbIdList, _ := listutil.Map(orders, func(t mysqldbmd.PermApprovalOrder) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, 0, err
	}
	data, _ := listutil.Map(orders, func(t mysqldbmd.PermApprovalOrder) (PermApprovalOrderDTO, error) {
		return PermApprovalOrderDTO{
			Id:           t.Id,
			Account:      t.Account,
			DbId:         t.DbId,
			DbHost:       dbMap[t.DbId].DbHost,
			DbName:       dbMap[t.DbId].Name,
			AccessBase:   t.AccessBase,
			AccessTables: t.AccessTables,
			PermType:     t.PermType,
			OrderStatus:  t.OrderStatus,
			Auditor:      t.Auditor,
			ExpireDay:    t.ExpireDay,
			Reason:       t.Reason,
			Created:      t.Created,
		}, nil
	})
	return data, next, nil
}

// ListAppliedPermApprovalOrder 展示申请的审批列表
func (*outerImpl) ListAppliedPermApprovalOrder(ctx context.Context, reqDTO ListAppliedPermApprovalOrderReqDTO) ([]PermApprovalOrderDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	orders, err := mysqldbmd.ListPermApprovalOrder(ctx, mysqldbmd.ListPermApprovalOrderReqDTO{
		Cursor:      reqDTO.Cursor,
		Limit:       reqDTO.Limit,
		Account:     reqDTO.Operator.Account,
		OrderStatus: reqDTO.OrderStatus,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(orders) == reqDTO.Limit {
		next = orders[len(orders)-1].Id
	}
	dbIdList, _ := listutil.Map(orders, func(t mysqldbmd.PermApprovalOrder) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, 0, err
	}
	data, _ := listutil.Map(orders, func(t mysqldbmd.PermApprovalOrder) (PermApprovalOrderDTO, error) {
		return PermApprovalOrderDTO{
			Id:           t.Id,
			Account:      t.Account,
			DbId:         t.DbId,
			DbHost:       dbMap[t.DbId].DbHost,
			DbName:       dbMap[t.DbId].Name,
			AccessBase:   t.AccessBase,
			AccessTables: t.AccessTables,
			PermType:     t.PermType,
			OrderStatus:  t.OrderStatus,
			Auditor:      t.Auditor,
			ExpireDay:    t.ExpireDay,
			Reason:       t.Reason,
			Created:      t.Created,
		}, nil
	})
	return data, next, nil
}

// AgreeDbPerm 同意审批
func (*outerImpl) AgreeDbPerm(ctx context.Context, reqDTO AgreeDbPermReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetPermApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != mysqldbmd.PendingPermOrderStatus {
		return util.InvalidArgsError()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b, err := mysqldbmd.UpdatePermApprovalOrderStatus(ctx, mysqldbmd.UpdatePermApprovalOrderStatusReqDTO{
			Id:        reqDTO.OrderId,
			NewStatus: mysqldbmd.AgreePermOrderStatus,
			OldStatus: order.OrderStatus,
			Auditor:   reqDTO.Operator.Account,
		})
		if err != nil {
			return err
		}
		if b {
			// 插入权限表
			expired := time.Now().Add(time.Duration(order.ExpireDay) * 24 * time.Hour)
			insertReqs, _ := listutil.Map(order.AccessTables, func(table string) (mysqldbmd.InsertPermReqDTO, error) {
				return mysqldbmd.InsertPermReqDTO{
					Account:     order.Account,
					DbId:        order.DbId,
					AccessBase:  order.AccessBase,
					AccessTable: table,
					PermType:    order.PermType,
					Expired:     expired,
				}, nil
			})
			return mysqldbmd.BatchInsertPerm(ctx, insertReqs)
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DisagreeDbPerm 不同意审批
func (*outerImpl) DisagreeDbPerm(ctx context.Context, reqDTO DisagreeDbPermReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetPermApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != mysqldbmd.PendingPermOrderStatus {
		return util.InvalidArgsError()
	}
	_, err = mysqldbmd.UpdatePermApprovalOrderStatus(ctx, mysqldbmd.UpdatePermApprovalOrderStatusReqDTO{
		Id:             reqDTO.OrderId,
		NewStatus:      mysqldbmd.DisagreePermOrderStatus,
		OldStatus:      order.OrderStatus,
		Auditor:        reqDTO.Operator.Account,
		DisagreeReason: reqDTO.DisagreeReason,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// CancelDbPerm 取消申请
func (*outerImpl) CancelDbPerm(ctx context.Context, reqDTO CancelDbPermReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetPermApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != mysqldbmd.PendingPermOrderStatus {
		return util.InvalidArgsError()
	}
	if order.Account != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	_, err = mysqldbmd.UpdatePermApprovalOrderStatus(ctx, mysqldbmd.UpdatePermApprovalOrderStatusReqDTO{
		Id:        reqDTO.OrderId,
		NewStatus: mysqldbmd.CanceledPermOrderStatus,
		OldStatus: order.OrderStatus,
		Auditor:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListDbPerm 权限列表
func (*outerImpl) ListDbPerm(ctx context.Context, reqDTO ListDbPermReqDTO) ([]PermDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	perms, err := mysqldbmd.ListPerm(ctx, mysqldbmd.ListPermReqDTO{
		Cursor:  reqDTO.Cursor,
		Limit:   reqDTO.Limit,
		Account: reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(perms) == reqDTO.Limit {
		next = perms[len(perms)-1].Id
	}
	dbIdList, _ := listutil.Map(perms, func(t mysqldbmd.Perm) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, 0, err
	}
	data, _ := listutil.Map(perms, func(t mysqldbmd.Perm) (PermDTO, error) {
		return PermDTO{
			Id:          t.Id,
			Account:     t.Account,
			DbId:        t.DbId,
			DbHost:      dbMap[t.DbId].DbHost,
			DbName:      dbMap[t.DbId].Name,
			AccessBase:  t.AccessBase,
			AccessTable: t.AccessTable,
			PermType:    t.PermType,
			Created:     t.Created,
			Expired:     t.Expired,
		}, nil
	})
	return data, next, nil
}

// DeleteDbPerm 删除权限
func (*outerImpl) DeleteDbPerm(ctx context.Context, reqDTO DeleteDbPermReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DbSrvKeysVO.DeleteDb),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if err = checkPerm(reqDTO.Operator); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err = mysqldbmd.DeletePerm(ctx, reqDTO.PermId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListDbPermByAccount 权限列表
func (*outerImpl) ListDbPermByAccount(ctx context.Context, reqDTO ListDbPermByAccountReqDTO) ([]PermDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	perms, err := mysqldbmd.ListPerm(ctx, mysqldbmd.ListPermReqDTO{
		Cursor:  reqDTO.Cursor,
		Limit:   reqDTO.Limit,
		Account: reqDTO.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(perms) == reqDTO.Limit {
		next = perms[len(perms)-1].Id
	}
	dbIdList, _ := listutil.Map(perms, func(t mysqldbmd.Perm) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, 0, err
	}
	data, _ := listutil.Map(perms, func(t mysqldbmd.Perm) (PermDTO, error) {
		return PermDTO{
			Id:          t.Id,
			Account:     t.Account,
			DbId:        t.DbId,
			DbHost:      dbMap[t.DbId].DbHost,
			DbName:      dbMap[t.DbId].Name,
			AccessTable: t.AccessTable,
			PermType:    t.PermType,
			Created:     t.Created,
			Expired:     t.Expired,
		}, nil
	})
	return data, next, nil
}

// AllTables 展示单个数据库所有表
func (*outerImpl) AllTables(ctx context.Context, reqDTO AllTablesReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	ret := make([]string, 0)
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(db.Username),
		url.QueryEscape(db.Password),
		db.DbHost,
		reqDTO.AccessBase,
	)
	_, tables, err := util.MysqlQuery(datasourceName, "show tables")
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	for _, table := range tables {
		if len(table) > 0 {
			ret = append(ret, table[0])
		}
	}
	return ret, nil
}

// AllBases 所有库
func (*outerImpl) AllBases(ctx context.Context, reqDTO AllBasesReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	ret := make([]string, 0)
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/?charset=utf8",
		url.QueryEscape(db.Username),
		url.QueryEscape(db.Password),
		db.DbHost,
	)
	_, tables, err := util.MysqlQuery(datasourceName, "show databases")
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	for _, table := range tables {
		if len(table) > 0 {
			ret = append(ret, table[0])
		}
	}
	return ret, nil
}

// SearchDb 搜索
func (*outerImpl) SearchDb(ctx context.Context, reqDTO SearchDbReqDTO) ([]string, [][]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, nil, util.InternalError(err)
	}
	if !b {
		return nil, nil, util.InvalidArgsError()
	}
	// 检查权限
	needCheckPerm := checkPerm(reqDTO.Operator) != nil
	tableName, sql, err := command.ValidateMysqlSelectSql(reqDTO.Cmd)
	if err != nil {
		return nil, nil, util.NewBizErrWithMsg(apicode.OperationFailedErrCode, err.Error())
	}
	if needCheckPerm {
		perms, err := mysqldbmd.SearchPerm(ctx, mysqldbmd.SearchPermReqDTO{
			Account:      reqDTO.Operator.Account,
			DbId:         reqDTO.DbId,
			AccessBase:   reqDTO.AccessBase,
			AccessTables: []string{tableName, "*"},
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, nil, util.InternalError(err)
		}
		permPass := false
		for _, perm := range perms {
			if perm.PermType.HasReadPermType() {
				permPass = true
				break
			}
		}
		if !permPass {
			return nil, nil, util.UnauthorizedError()
		}
	}
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(db.Username),
		url.QueryEscape(db.Password),
		db.DbHost,
		reqDTO.AccessBase,
	)
	columns, ret, err := util.MysqlQuery(datasourceName, strings.TrimSuffix(sql, ";")+" limit "+strconv.Itoa(reqDTO.Limit))
	if err != nil {
		return nil, nil, util.NewBizErrWithMsg(apicode.OperationFailedErrCode, err.Error())
	}
	return columns, ret, nil
}

// ApplyDbUpdate 提数据库修改单
func (*outerImpl) ApplyDbUpdate(ctx context.Context, reqDTO ApplyDbUpdateReqDTO) ([]command.ValidateUpdateResult, bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, false, util.InternalError(err)
	}
	if !b {
		return nil, false, util.InvalidArgsError()
	}
	// 检查权限
	needCheckPerm := checkPerm(reqDTO.Operator) != nil
	validateResults, allPass, err := command.ValidateMysqlUpdateSql(reqDTO.Cmd)
	if err != nil {
		return nil, false, util.NewBizErrWithMsg(apicode.OperationFailedErrCode, err.Error())
	}
	if needCheckPerm {
		for i := range validateResults {
			if !validateResults[i].Pass {
				continue
			}
			perms, err := mysqldbmd.SearchPerm(ctx, mysqldbmd.SearchPermReqDTO{
				Account:      reqDTO.Operator.Account,
				DbId:         reqDTO.DbId,
				AccessBase:   reqDTO.AccessBase,
				AccessTables: []string{validateResults[i].TableName, "*"},
			})
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return nil, false, util.InternalError(err)
			}
			permPass := false
			for _, perm := range perms {
				if perm.PermType.HasWritePermType() {
					permPass = true
					break
				}
			}
			if !permPass {
				validateResults[i].ErrMsg = i18n.GetByKey(i18n.SystemUnauthorized)
				validateResults[i].Pass = false
				allPass = false
			}
		}
	}
	if allPass {
		// 插入数据库
		err = mysqldbmd.InsertUpdateApprovalOrder(ctx, mysqldbmd.InsertUpdateApprovalOrderReqDTO{
			Name:        reqDTO.Name,
			Account:     reqDTO.Operator.Account,
			DbId:        reqDTO.DbId,
			AccessBase:  reqDTO.AccessBase,
			UpdateCmd:   reqDTO.Cmd,
			OrderStatus: mysqldbmd.PendingUpdateOrderStatus,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, false, util.InternalError(err)
		}
	}
	return validateResults, allPass, nil
}

// ListUpdateApprovalOrder 数据库修改审批单
func (*outerImpl) ListUpdateApprovalOrder(ctx context.Context, reqDTO ListUpdateApprovalOrderReqDTO) ([]UpdateApprovalOrderDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	orders, err := mysqldbmd.ListUpdateApprovalOrder(ctx, mysqldbmd.ListUpdateApprovalOrderReqDTO{
		Cursor:      reqDTO.Cursor,
		Limit:       reqDTO.Limit,
		OrderStatus: reqDTO.OrderStatus,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(orders) == reqDTO.Limit {
		next = orders[len(orders)-1].Id
	}
	data, err := updateApprovalOrdersMd2Dto(ctx, orders)
	if err != nil {
		return nil, 0, err
	}
	return data, next, nil
}

// ListAppliedUpdateApprovalOrder 申请的数据库修改审批单
func (*outerImpl) ListAppliedUpdateApprovalOrder(ctx context.Context, reqDTO ListAppliedUpdateApprovalOrderReqDTO) ([]UpdateApprovalOrderDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	orders, err := mysqldbmd.ListUpdateApprovalOrder(ctx, mysqldbmd.ListUpdateApprovalOrderReqDTO{
		Cursor:      reqDTO.Cursor,
		Limit:       reqDTO.Limit,
		Account:     reqDTO.Operator.Account,
		OrderStatus: reqDTO.OrderStatus,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(orders) == reqDTO.Limit {
		next = orders[len(orders)-1].Id
	}
	data, err := updateApprovalOrdersMd2Dto(ctx, orders)
	if err != nil {
		return nil, 0, err
	}
	return data, next, nil
}

func updateApprovalOrdersMd2Dto(ctx context.Context, orders []mysqldbmd.UpdateApprovalOrder) ([]UpdateApprovalOrderDTO, error) {
	dbIdList, _ := listutil.Map(orders, func(t mysqldbmd.UpdateApprovalOrder) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, err
	}
	data, _ := listutil.Map(orders, func(t mysqldbmd.UpdateApprovalOrder) (UpdateApprovalOrderDTO, error) {
		return UpdateApprovalOrderDTO{
			Id:          t.Id,
			Account:     t.Account,
			Name:        t.Name,
			DbId:        t.DbId,
			DbHost:      dbMap[t.DbId].DbHost,
			DbName:      dbMap[t.DbId].Name,
			AccessBase:  t.AccessBase,
			OrderStatus: t.OrderStatus,
			UpdateCmd:   t.UpdateCmd,
			Auditor:     t.Auditor,
			Created:     t.Created,
			ExecuteLog:  t.ExecuteLog,
		}, nil
	})
	return data, nil
}

// AgreeDbUpdate 同意修改单
func (*outerImpl) AgreeDbUpdate(ctx context.Context, reqDTO AgreeDbUpdateReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetUpdateApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != mysqldbmd.PendingUpdateOrderStatus {
		return util.InvalidArgsError()
	}
	_, err = mysqldbmd.UpdateUpdateApprovalOrderStatus(ctx, mysqldbmd.UpdateUpdateApprovalOrderStatusReqDTO{
		Id:        reqDTO.OrderId,
		NewStatus: mysqldbmd.AgreeUpdateOrderStatus,
		OldStatus: order.OrderStatus,
		Auditor:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DisagreeDbUpdate 不同意修改单
func (*outerImpl) DisagreeDbUpdate(ctx context.Context, reqDTO DisagreeDbUpdateReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetUpdateApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != mysqldbmd.PendingUpdateOrderStatus {
		return util.InvalidArgsError()
	}
	_, err = mysqldbmd.UpdateUpdateApprovalOrderStatus(ctx, mysqldbmd.UpdateUpdateApprovalOrderStatusReqDTO{
		Id:             reqDTO.OrderId,
		NewStatus:      mysqldbmd.DisagreeUpdateOrderStatus,
		OldStatus:      order.OrderStatus,
		Auditor:        reqDTO.Operator.Account,
		DisagreeReason: reqDTO.DisagreeReason,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// CancelDbUpdate 取消修改单
func (*outerImpl) CancelDbUpdate(ctx context.Context, reqDTO CancelDbUpdateReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetUpdateApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != mysqldbmd.PendingUpdateOrderStatus {
		return util.InvalidArgsError()
	}
	if order.Account != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	_, err = mysqldbmd.UpdateUpdateApprovalOrderStatus(ctx, mysqldbmd.UpdateUpdateApprovalOrderStatusReqDTO{
		Id:        reqDTO.OrderId,
		NewStatus: mysqldbmd.CanceledUpdateOrderStatus,
		OldStatus: order.OrderStatus,
		Auditor:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ExecuteDbUpdate 执行修改单
func (*outerImpl) ExecuteDbUpdate(ctx context.Context, reqDTO ExecuteDbUpdateReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetUpdateApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b ||
		(order.OrderStatus != mysqldbmd.AgreeUpdateOrderStatus &&
			order.OrderStatus != mysqldbmd.AskToExecuteUpdateOrderStatus) {
		return util.InvalidArgsError()
	}
	db, b, err := mysqldbmd.GetDbById(ctx, order.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	b, err = mysqldbmd.UpdateUpdateApprovalOrderStatus(ctx, mysqldbmd.UpdateUpdateApprovalOrderStatusReqDTO{
		Id:        reqDTO.OrderId,
		NewStatus: mysqldbmd.ExecutedUpdateOrderStatus,
		OldStatus: order.OrderStatus,
		Auditor:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		go executeUpdateCmd(&order, &db)
	}
	return nil
}

func executeUpdateCmd(order *mysqldbmd.UpdateApprovalOrder, db *mysqldbmd.Db) {
	logMsg := strings.Builder{}
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(db.Username),
		url.QueryEscape(db.Password),
		db.DbHost,
		order.AccessBase,
	)
	results, err := util.MysqlExecute(datasourceName, order.UpdateCmd)
	if err != nil {
		logMsg.WriteString(err.Error())
	} else {
		for _, result := range results {
			logMsg.WriteString("sql: " + result.Sql + "\n")
			logMsg.WriteString("affectedRows: " + strconv.FormatInt(result.AffectedRows, 10) + "\n")
			if result.ErrMsg != "" {
				logMsg.WriteString("errMsg: " + result.ErrMsg + "\n")
			}
			logMsg.WriteString("\n")
			logMsg.WriteString("\n")
		}
	}
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, err = mysqldbmd.UpdateUpdateApprovalOrderExecuteLog(ctx, order.Id, logMsg.String())
	if err != nil {
		logger.Logger.Errorf("db update: %v executed with err: %v", order.Id, err)
	}
}

// AskToExecuteDbUpdate 请求执行修改单
func (*outerImpl) AskToExecuteDbUpdate(ctx context.Context, reqDTO AskToExecuteDbUpdateReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetUpdateApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != mysqldbmd.AgreeUpdateOrderStatus {
		return util.InvalidArgsError()
	}
	if order.Account != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	_, err = mysqldbmd.UpdateUpdateApprovalOrderStatus(ctx, mysqldbmd.UpdateUpdateApprovalOrderStatusReqDTO{
		Id:        reqDTO.OrderId,
		NewStatus: mysqldbmd.AskToExecuteUpdateOrderStatus,
		OldStatus: order.OrderStatus,
		Auditor:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func getDbMap(ctx context.Context, dbIdList []int64) (map[int64]mysqldbmd.Db, error) {
	dbIdList = hashset.NewHashSet(dbIdList...).AllKeys()
	dbList, err := mysqldbmd.BatchGetDbByIdList(ctx, dbIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	dbMap := make(map[int64]mysqldbmd.Db)
	for _, db := range dbList {
		dbMap[db.Id] = db
	}
	return dbMap, nil
}

func checkPerm(operator apisession.UserInfo) error {
	if operator.IsAdmin || operator.RoleType.IsDba() {
		return nil
	}
	return util.UnauthorizedError()
}
