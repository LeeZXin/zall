package dbsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/dbaudit/modules/model/dbmd"
	"github.com/LeeZXin/zall/dbaudit/modules/service/dbsrv/command"
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
	err = dbmd.InsertDb(ctx, dbmd.InsertDbReqDTO{
		Name:     reqDTO.Name,
		DbHost:   reqDTO.DbHost,
		Username: reqDTO.Username,
		Password: reqDTO.Password,
		DbType:   reqDTO.DbType,
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
	_, err = dbmd.UpdateDb(ctx, dbmd.UpdateDbReqDTO{
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
	_, err = dbmd.DeleteDbById(ctx, reqDTO.Id)
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
	dbs, err := dbmd.ListDb(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(dbs, func(t dbmd.Db) (DbDTO, error) {
		return DbDTO{
			Id:       t.Id,
			Name:     t.Name,
			DbHost:   t.DbHost,
			Username: t.Username,
			Password: t.Password,
			DbType:   t.DbType,
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
	dbs, err := dbmd.ListDb(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(dbs, func(t dbmd.Db) (SimpleDbDTO, error) {
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
	_, b, err := dbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	err = dbmd.InsertApprovalOrder(ctx, dbmd.InsertApprovalOrderReqDTO{
		Account:      reqDTO.Operator.Account,
		DbId:         reqDTO.DbId,
		AccessBase:   reqDTO.AccessBase,
		AccessTables: reqDTO.AccessTables,
		PermType:     reqDTO.PermType,
		OrderStatus:  dbmd.PendingOrderStatus,
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
func (*outerImpl) ListPermApprovalOrder(ctx context.Context, reqDTO ListPermApprovalOrderReqDTO) ([]ApprovalOrderDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	if err := checkPerm(reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	orders, err := dbmd.ListApprovalOrder(ctx, dbmd.ListApprovalOrderReqDTO{
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
	dbIdList, _ := listutil.Map(orders, func(t dbmd.ApprovalOrder) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, 0, err
	}
	data, _ := listutil.Map(orders, func(t dbmd.ApprovalOrder) (ApprovalOrderDTO, error) {
		return ApprovalOrderDTO{
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
	order, b, err := dbmd.GetApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != dbmd.PendingOrderStatus {
		return util.InvalidArgsError()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b, err := dbmd.UpdateApprovalOrderStatus(ctx, dbmd.UpdateApprovalOrderStatusReqDTO{
			Id:        reqDTO.OrderId,
			NewStatus: dbmd.AgreeOrderStatus,
			OldStatus: order.OrderStatus,
			Auditor:   reqDTO.Operator.Account,
		})
		if err != nil {
			return err
		}
		if b {
			// 插入权限表
			expired := time.Now().Add(time.Duration(order.ExpireDay) * 24 * time.Hour)
			insertReqs, _ := listutil.Map(order.AccessTables, func(table string) (dbmd.InsertPermReqDTO, error) {
				return dbmd.InsertPermReqDTO{
					Account:     order.Account,
					DbId:        order.DbId,
					AccessBase:  order.AccessBase,
					AccessTable: table,
					PermType:    order.PermType,
					Expired:     expired,
				}, nil
			})
			return dbmd.BatchInsertPerm(ctx, insertReqs)
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
	order, b, err := dbmd.GetApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != dbmd.PendingOrderStatus {
		return util.InvalidArgsError()
	}
	_, err = dbmd.UpdateApprovalOrderStatus(ctx, dbmd.UpdateApprovalOrderStatusReqDTO{
		Id:        reqDTO.OrderId,
		NewStatus: dbmd.DisagreeOrderStatus,
		OldStatus: order.OrderStatus,
		Auditor:   reqDTO.Operator.Account,
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
	order, b, err := dbmd.GetApprovalOrderById(ctx, reqDTO.OrderId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.OrderStatus != dbmd.PendingOrderStatus {
		return util.InvalidArgsError()
	}
	if order.Account != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	_, err = dbmd.UpdateApprovalOrderStatus(ctx, dbmd.UpdateApprovalOrderStatusReqDTO{
		Id:        reqDTO.OrderId,
		NewStatus: dbmd.CanceledOrderStatus,
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
	perms, err := dbmd.ListPerm(ctx, dbmd.ListPermReqDTO{
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
	dbIdList, _ := listutil.Map(perms, func(t dbmd.Perm) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, 0, err
	}
	data, _ := listutil.Map(perms, func(t dbmd.Perm) (PermDTO, error) {
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
	err = dbmd.DeletePerm(ctx, reqDTO.PermId)
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
	perms, err := dbmd.ListPerm(ctx, dbmd.ListPermReqDTO{
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
	dbIdList, _ := listutil.Map(perms, func(t dbmd.Perm) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, 0, err
	}
	data, _ := listutil.Map(perms, func(t dbmd.Perm) (PermDTO, error) {
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
	db, b, err := dbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	ret := make([]string, 0)
	switch db.DbType {
	case dbmd.MysqlDbType:
		datasourceName := fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8",
			url.QueryEscape(db.Username),
			url.QueryEscape(db.Password),
			db.DbHost,
			reqDTO.AccessBase,
		)
		_, tables, err := command.MysqlQuery(datasourceName, "show tables")
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		for _, table := range tables {
			if len(table) > 0 {
				ret = append(ret, table[0])
			}
		}
	case dbmd.RedisDbType:
	case dbmd.MongoDbType:
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
	db, b, err := dbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	ret := make([]string, 0)
	switch db.DbType {
	case dbmd.MysqlDbType:
		datasourceName := fmt.Sprintf(
			"%s:%s@tcp(%s)/?charset=utf8",
			url.QueryEscape(db.Username),
			url.QueryEscape(db.Password),
			db.DbHost,
		)
		_, tables, err := command.MysqlQuery(datasourceName, "show databases")
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		for _, table := range tables {
			if len(table) > 0 {
				ret = append(ret, table[0])
			}
		}
	case dbmd.RedisDbType:
	case dbmd.MongoDbType:
	}
	return ret, nil
}

// SearchDb 搜索
func (*outerImpl) SearchDb(ctx context.Context, reqDTO SearchDbReqDTO) (columns []string, ret [][]string, err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DbSrvKeysVO.SearchDb),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		db dbmd.Db
		b  bool
	)
	db, b, err = dbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 检查权限
	needCheckPerm := checkPerm(reqDTO.Operator) != nil
	switch db.DbType {
	case dbmd.MysqlDbType:
		var (
			tableName, sql string
		)
		tableName, sql, err = command.ValidateMysqlSelectSql(reqDTO.Cmd)
		if err != nil {
			err = util.NewBizErrWithMsg(apicode.OperationFailedErrCode, err.Error())
			return
		}
		if needCheckPerm {
			var perms []dbmd.Perm
			perms, err = dbmd.SearchPerm(ctx, dbmd.SearchPermReqDTO{
				Account:      reqDTO.Operator.Account,
				DbId:         reqDTO.DbId,
				AccessBase:   reqDTO.AccessBase,
				AccessTables: []string{tableName, "*"},
			})
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				err = util.InternalError(err)
				return
			}
			permPass := false
			for _, perm := range perms {
				if perm.PermType.HasReadPermType() {
					permPass = true
					break
				}
			}
			if !permPass {
				err = util.UnauthorizedError()
				return
			}
		}
		datasourceName := fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8",
			url.QueryEscape(db.Username),
			url.QueryEscape(db.Password),
			db.DbHost,
			reqDTO.AccessBase,
		)
		columns, ret, err = command.MysqlQuery(datasourceName, strings.TrimSuffix(sql, ";")+" limit "+strconv.Itoa(reqDTO.Limit))
		if err != nil {
			err = util.NewBizErrWithMsg(apicode.OperationFailedErrCode, err.Error())
			return
		}
	case dbmd.RedisDbType:
	case dbmd.MongoDbType:
	}
	return
}

func getDbMap(ctx context.Context, dbIdList []int64) (map[int64]dbmd.Db, error) {
	dbIdList = hashset.NewHashSet(dbIdList...).AllKeys()
	dbList, err := dbmd.BatchGetDbByIdList(ctx, dbIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	dbMap := make(map[int64]dbmd.Db)
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
