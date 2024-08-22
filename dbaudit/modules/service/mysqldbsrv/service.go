package mysqldbsrv

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/dbaudit/modules/model/mysqldbmd"
	"github.com/LeeZXin/zall/dbaudit/modules/service/mysqldbsrv/command"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	defaultBases = hashset.NewHashSet("information_schema", "mysql", "sys", "performance_schema")
)

// CreateDb 创建数据库
func CreateDb(ctx context.Context, reqDTO CreateDbReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if err = checkDbaPerm(reqDTO.Operator); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err = mysqldbmd.InsertDb(ctx, mysqldbmd.InsertDbReqDTO{
		Name:   reqDTO.Name,
		Config: reqDTO.Config,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// UpdateDb 编辑数据库
func UpdateDb(ctx context.Context, reqDTO UpdateDbReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if err = checkDbaPerm(reqDTO.Operator); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = mysqldbmd.UpdateDb(ctx, mysqldbmd.UpdateDbReqDTO{
		Id:     reqDTO.DbId,
		Name:   reqDTO.Name,
		Config: reqDTO.Config,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// DeleteDb 删除数据库
func DeleteDb(ctx context.Context, reqDTO DeleteDbReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if err = checkDbaPerm(reqDTO.Operator); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除数据库
		_, err2 := mysqldbmd.DeleteDbById(ctx, reqDTO.DbId)
		if err2 != nil {
			return err2
		}
		// 删除读权限
		_, err2 = mysqldbmd.DeleteReadPermByDbId(ctx, reqDTO.DbId)
		if err2 != nil {
			return err2
		}
		// 删除读权限申请
		_, err2 = mysqldbmd.DeleteReadPermApplyByDbId(ctx, reqDTO.DbId)
		if err2 != nil {
			return err2
		}
		// 删除数据库修改单
		_, err2 = mysqldbmd.DeleteDataUpdateApplyByDbId(ctx, reqDTO.DbId)
		return err2
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListDb 数据库列表
func ListDb(ctx context.Context, reqDTO ListDbReqDTO) ([]DbDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	if err := checkDbaPerm(reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	dbs, total, err := mysqldbmd.PageDb(ctx, mysqldbmd.PageDbReqDTO{
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
		Name:     reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(dbs, func(t mysqldbmd.Db) (DbDTO, error) {
		ret := DbDTO{
			Id:      t.Id,
			Name:    t.Name,
			Created: t.Created,
		}
		if t.Config != nil {
			ret.Config = t.Config.Data
		}
		return ret, nil
	})
	return data, total, nil
}

// ListSimpleDb 数据库列表
func ListSimpleDb(ctx context.Context, reqDTO ListSimpleDbReqDTO) ([]SimpleDbDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	dbs, err := mysqldbmd.ListDb(ctx, []string{"id", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(dbs, func(t mysqldbmd.Db) (SimpleDbDTO, error) {
		return SimpleDbDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
}

// ApplyReadPerm 申请库表读权限
func ApplyReadPerm(ctx context.Context, reqDTO ApplyReadPermReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	err = mysqldbmd.InsertReadPermApply(ctx, mysqldbmd.InsertReadPermApplyReqDTO{
		Account:      reqDTO.Operator.Account,
		DbId:         reqDTO.DbId,
		DbName:       db.Name,
		AccessBase:   reqDTO.AccessBase,
		AccessTables: reqDTO.AccessTables,
		OrderStatus:  mysqldbmd.PendingReadPermApplyStatus,
		ExpireDay:    reqDTO.ExpireDay,
		ApplyReason:  reqDTO.ApplyReason,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// GetReadPermApply 查看读权限审批单
func GetReadPermApply(ctx context.Context, reqDTO GetReadPermApplyReqDTO) (ReadPermApplyDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return ReadPermApplyDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	apply, b, err := mysqldbmd.GetReadPermApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ReadPermApplyDTO{}, util.InternalError(err)
	}
	if !b {
		return ReadPermApplyDTO{}, util.InvalidArgsError()
	}
	if checkDbaPerm(reqDTO.Operator) != nil && apply.Account != reqDTO.Operator.Account {
		return ReadPermApplyDTO{}, util.UnauthorizedError()
	}
	return readPermApply2Dto(apply)
}

// ListReadPermApplyByDba dba查看展示审批列表
func ListReadPermApplyByDba(ctx context.Context, reqDTO ListReadPermApplyByDbaReqDTO) ([]ReadPermApplyDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	if err := checkDbaPerm(reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	applies, total, err := mysqldbmd.ListReadPermApply(ctx, mysqldbmd.ListReadPermApplyReqDTO{
		PageNum:     reqDTO.PageNum,
		PageSize:    10,
		ApplyStatus: reqDTO.ApplyStatus,
		DbId:        reqDTO.DbId,
	},
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(applies, readPermApply2Dto)
	return data, total, nil
}

func readPermApply2Dto(t mysqldbmd.ReadPermApply) (ReadPermApplyDTO, error) {
	return ReadPermApplyDTO{
		Id:             t.Id,
		Account:        t.Account,
		DbId:           t.DbId,
		DbName:         t.DbName,
		AccessBase:     t.AccessBase,
		AccessTables:   t.AccessTables,
		ApplyStatus:    t.ApplyStatus,
		Auditor:        t.Auditor,
		ExpireDay:      t.ExpireDay,
		ApplyReason:    t.ApplyReason,
		DisagreeReason: t.DisagreeReason,
		Created:        t.Created,
		Updated:        t.Updated,
	}, nil
}

// ListReadPermApplyByOperator 展示申请的审批列表
func ListReadPermApplyByOperator(ctx context.Context, reqDTO ListReadPermApplyByOperatorReqDTO) ([]ReadPermApplyDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	applies, total, err := mysqldbmd.ListReadPermApply(ctx, mysqldbmd.ListReadPermApplyReqDTO{
		PageNum:     reqDTO.PageNum,
		PageSize:    10,
		Account:     reqDTO.Operator.Account,
		ApplyStatus: reqDTO.ApplyStatus,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(applies, readPermApply2Dto)
	return data, total, nil
}

// AgreeReadPermApply 同意审批
func AgreeReadPermApply(ctx context.Context, reqDTO AgreeReadPermApplyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkDbaPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	apply, b, err := mysqldbmd.GetReadPermApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || apply.ApplyStatus != mysqldbmd.PendingReadPermApplyStatus {
		return util.InvalidArgsError()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b2, err2 := mysqldbmd.UpdateReadPermApplyStatus(ctx, mysqldbmd.UpdateReadPermApplyStatusReqDTO{
			Id:        reqDTO.ApplyId,
			NewStatus: mysqldbmd.AgreeReadPermApplyStatus,
			OldStatus: apply.ApplyStatus,
			Auditor:   reqDTO.Operator.Account,
		})
		if err2 != nil {
			return err2
		}
		if b2 {
			// 插入权限表
			expired := time.Now().Add(time.Duration(apply.ExpireDay) * 24 * time.Hour)
			tables, _ := listutil.Filter(listutil.Distinct(strings.Split(apply.AccessTables, ";")...), func(t string) (bool, error) {
				return len(t) > 0, nil
			})
			insertReqs, _ := listutil.Map(tables, func(t string) (mysqldbmd.InsertReadPermReqDTO, error) {
				return mysqldbmd.InsertReadPermReqDTO{
					Account:     apply.Account,
					DbId:        apply.DbId,
					AccessBase:  apply.AccessBase,
					AccessTable: t,
					ApplyId:     apply.Id,
					Expired:     expired,
				}, nil
			})
			return mysqldbmd.BatchInsertReadPerm(ctx, insertReqs)
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListAuthorizedDb 展示授权的库表
func ListAuthorizedDb(ctx context.Context, reqDTO ListAuthorizedDbReqDTO) ([]SimpleDbDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 先删除过期数据
	err := mysqldbmd.DeleteExpiredReadPermByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	isDba := checkDbaPerm(reqDTO.Operator) == nil
	if isDba {
		dbs, err := mysqldbmd.ListDb(ctx, []string{"id", "name"})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		return listutil.Map(dbs, func(t mysqldbmd.Db) (SimpleDbDTO, error) {
			return SimpleDbDTO{
				Id:   t.Id,
				Name: t.Name,
			}, nil
		})
	}
	perms, err := mysqldbmd.ListReadPermByAccount(ctx, mysqldbmd.ListReadPermByAccountReqDTO{
		Account: reqDTO.Operator.Account,
		Cols:    []string{"db_id"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	dbIdList, _ := listutil.Map(perms, func(t mysqldbmd.ReadPerm) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, err
	}
	ret := make([]SimpleDbDTO, 0)
	for dbId, db := range dbMap {
		ret = append(ret, SimpleDbDTO{
			Id:   dbId,
			Name: db.Name,
		})
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i].Id < ret[j].Id
	})
	return ret, nil
}

// ListAuthorizedBase 展示授权的库
func ListAuthorizedBase(ctx context.Context, reqDTO ListAuthorizedBaseReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 先删除过期数据
	err := mysqldbmd.DeleteExpiredReadPermByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	isDba := checkDbaPerm(reqDTO.Operator) == nil
	if isDba {
		db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		if !b {
			return nil, util.InvalidArgsError()
		}
		datasourceName := fmt.Sprintf(
			"%s:%s@tcp(%s)/?charset=utf8",
			url.QueryEscape(db.Config.Data.ReadNode.Username),
			url.QueryEscape(db.Config.Data.ReadNode.Password),
			db.Config.Data.ReadNode.Host,
		)
		result, err := util.MysqlQuery(datasourceName, "show databases")
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		bases := result.Data
		ret := make([]string, 0, len(bases))
		for _, base := range bases {
			if len(base) > 0 {
				ret = append(ret, base[0])
			}
		}
		ret, _ = listutil.Filter(ret, func(t string) (bool, error) {
			return !defaultBases.Contains(t), nil
		})
		return ret, nil
	}
	perms, err := mysqldbmd.ListReadPermByAccount(ctx, mysqldbmd.ListReadPermByAccountReqDTO{
		DbId:    reqDTO.DbId,
		Account: reqDTO.Operator.Account,
		Cols:    []string{"access_base"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret, _ := listutil.Map(perms, func(t mysqldbmd.ReadPerm) (string, error) {
		return t.AccessBase, nil
	})
	ret = listutil.Distinct(ret...)
	ret, _ = listutil.Filter(ret, func(t string) (bool, error) {
		return !defaultBases.Contains(t), nil
	})
	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return ret, nil
}

// ListAuthorizedTable 展示授权的表
func ListAuthorizedTable(ctx context.Context, reqDTO ListAuthorizedTableReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 先删除过期数据
	err := mysqldbmd.DeleteExpiredReadPermByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	isDba := checkDbaPerm(reqDTO.Operator) == nil
	if isDba {
		db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		if !b {
			return nil, util.InvalidArgsError()
		}
		datasourceName := fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8",
			url.QueryEscape(db.Config.Data.ReadNode.Username),
			url.QueryEscape(db.Config.Data.ReadNode.Password),
			db.Config.Data.ReadNode.Host,
			reqDTO.AccessBase,
		)
		result, err := util.MysqlQuery(datasourceName, "show tables")
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		tables := result.Data
		ret := make([]string, 0, len(tables))
		for _, table := range tables {
			if len(table) > 0 {
				ret = append(ret, table[0])
			}
		}
		return ret, nil
	}
	perms, err := mysqldbmd.ListReadPermByAccount(ctx, mysqldbmd.ListReadPermByAccountReqDTO{
		DbId:       reqDTO.DbId,
		AccessBase: reqDTO.AccessBase,
		Account:    reqDTO.Operator.Account,
		Cols:       []string{"access_table"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	retTables, _ := listutil.Map(perms, func(t mysqldbmd.ReadPerm) (string, error) {
		return t.AccessTable, nil
	})
	retTables = listutil.Distinct(retTables...)
	hasStar := false
	for _, table := range retTables {
		if table == "*" {
			hasStar = true
			break
		}
	}
	if hasStar {
		db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		if !b {
			return nil, util.ThereHasBugErr()
		}
		datasourceName := fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8",
			url.QueryEscape(db.Config.Data.ReadNode.Username),
			url.QueryEscape(db.Config.Data.ReadNode.Password),
			db.Config.Data.ReadNode.Host,
			reqDTO.AccessBase,
		)
		result, err := util.MysqlQuery(datasourceName, "show tables")
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		tables := result.Data
		retTables = make([]string, 0, len(tables))
		for _, table := range tables {
			if len(table) > 0 {
				retTables = append(retTables, table[0])
			}
		}
	}
	sort.SliceStable(retTables, func(i, j int) bool {
		return retTables[i] < retTables[j]
	})
	return retTables, nil
}

// GetCreateTableSql 展示建表语句
func GetCreateTableSql(ctx context.Context, reqDTO GetCreateSqlReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 先删除过期数据
	err := mysqldbmd.DeleteExpiredReadPermByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", util.ThereHasBugErr()
	}
	needCheckPerm := checkDbaPerm(reqDTO.Operator) != nil
	if needCheckPerm {
		hasPerm, err := mysqldbmd.ExistReadPerm(ctx, mysqldbmd.ExistReadPermReqDTO{
			Account:      reqDTO.Operator.Account,
			DbId:         reqDTO.DbId,
			AccessBase:   reqDTO.AccessBase,
			AccessTables: []string{reqDTO.AccessTable, "*"},
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return "", util.InternalError(err)
		}
		if !hasPerm {
			return "", util.UnauthorizedError()
		}
	}
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(db.Config.Data.ReadNode.Username),
		url.QueryEscape(db.Config.Data.ReadNode.Password),
		db.Config.Data.ReadNode.Host,
		reqDTO.AccessBase,
	)
	result, err := util.MysqlQuery(datasourceName, "show create table "+reqDTO.AccessTable)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	data := result.Data
	if len(data) == 1 && len(data[0]) == 2 {
		return data[0][1], nil
	}
	return "", nil
}

// ShowTableIndex 展示索引语句
func ShowTableIndex(ctx context.Context, reqDTO ShowTableIndexReqDTO) ([]string, [][]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 先删除过期数据
	err := mysqldbmd.DeleteExpiredReadPermByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, nil, util.InternalError(err)
	}
	db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, nil, util.InternalError(err)
	}
	if !b {
		return nil, nil, util.ThereHasBugErr()
	}
	needCheckPerm := checkDbaPerm(reqDTO.Operator) != nil
	if needCheckPerm {
		hasPerm, err := mysqldbmd.ExistReadPerm(ctx, mysqldbmd.ExistReadPermReqDTO{
			Account:      reqDTO.Operator.Account,
			DbId:         reqDTO.DbId,
			AccessBase:   reqDTO.AccessBase,
			AccessTables: []string{reqDTO.AccessTable, "*"},
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, nil, util.InternalError(err)
		}
		if !hasPerm {
			return nil, nil, util.UnauthorizedError()
		}
	}
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(db.Config.Data.ReadNode.Username),
		url.QueryEscape(db.Config.Data.ReadNode.Password),
		db.Config.Data.ReadNode.Host,
		reqDTO.AccessBase,
	)
	result, err := util.MysqlQuery(datasourceName, "show index from "+reqDTO.AccessTable)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, nil, util.InternalError(err)
	}
	return result.Columns, result.Data, nil
}

// DisagreeReadPermApply 不同意审批
func DisagreeReadPermApply(ctx context.Context, reqDTO DisagreeReadPermApplyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkDbaPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetReadPermApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.ApplyStatus != mysqldbmd.PendingReadPermApplyStatus {
		return util.InvalidArgsError()
	}
	_, err = mysqldbmd.UpdateReadPermApplyStatus(ctx, mysqldbmd.UpdateReadPermApplyStatusReqDTO{
		Id:             reqDTO.ApplyId,
		NewStatus:      mysqldbmd.DisagreeReadPermApplyStatus,
		OldStatus:      order.ApplyStatus,
		Auditor:        reqDTO.Operator.Account,
		DisagreeReason: reqDTO.DisagreeReason,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// CancelReadPermApply 取消申请
func CancelReadPermApply(ctx context.Context, reqDTO CancelReadPermApplyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	apply, b, err := mysqldbmd.GetReadPermApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || apply.ApplyStatus != mysqldbmd.PendingReadPermApplyStatus {
		return util.InvalidArgsError()
	}
	if apply.Account != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	_, err = mysqldbmd.UpdateReadPermApplyStatus(ctx, mysqldbmd.UpdateReadPermApplyStatusReqDTO{
		Id:        reqDTO.ApplyId,
		NewStatus: mysqldbmd.CanceledReadPermApplyStatus,
		OldStatus: apply.ApplyStatus,
		Auditor:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListReadPermByOperator 权限列表
func ListReadPermByOperator(ctx context.Context, reqDTO ListReadPermByOperatorReqDTO) ([]ReadPermDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 先删除过期的
	err := mysqldbmd.DeleteExpiredReadPermByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	perms, total, err := mysqldbmd.PageReadPerm(ctx, mysqldbmd.PageReadPermReqDTO{
		Account:  reqDTO.Operator.Account,
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
		DbId:     reqDTO.DbId,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	dbIdList, _ := listutil.Map(perms, func(t mysqldbmd.ReadPerm) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, 0, err
	}
	data, _ := listutil.Map(perms, func(t mysqldbmd.ReadPerm) (ReadPermDTO, error) {
		return readPerm2Dto(t, dbMap), nil
	})
	return data, total, nil
}

func readPerm2Dto(t mysqldbmd.ReadPerm, dbMap map[int64]mysqldbmd.Db) ReadPermDTO {
	return ReadPermDTO{
		Id:          t.Id,
		Account:     t.Account,
		DbId:        t.DbId,
		DbName:      dbMap[t.DbId].Name,
		AccessBase:  t.AccessBase,
		AccessTable: t.AccessTable,
		Created:     t.Created,
		Expired:     t.Expired,
		ApplyId:     t.ApplyId,
	}
}

// DeleteReadPermByDba 删除权限
func DeleteReadPermByDba(ctx context.Context, reqDTO DeleteReadPermByDbaReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if err = checkDbaPerm(reqDTO.Operator); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err = mysqldbmd.DeleteReadPermById(ctx, reqDTO.PermId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListReadPermByDba 权限列表
func ListReadPermByDba(ctx context.Context, reqDTO ListReadPermByDbaReqDTO) ([]ReadPermDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	if err := checkDbaPerm(reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	perms, total, err := mysqldbmd.ListReadPerm(ctx, mysqldbmd.ListReadPermReqDTO{
		DbId:     reqDTO.DbId,
		Account:  reqDTO.Account,
		PageSize: 10,
		PageNum:  reqDTO.PageNum,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	dbIdList, _ := listutil.Map(perms, func(t mysqldbmd.ReadPerm) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, 0, err
	}
	data, _ := listutil.Map(perms, func(t mysqldbmd.ReadPerm) (ReadPermDTO, error) {
		return readPerm2Dto(t, dbMap), nil
	})
	return data, total, nil
}

// ExecuteSelectSql 搜索
func ExecuteSelectSql(ctx context.Context, reqDTO ExecuteSelectSqlReqDTO) (util.MysqlQueryResult, error) {
	if err := reqDTO.IsValid(); err != nil {
		return util.MysqlQueryResult{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	db, b, err := mysqldbmd.GetDbById(ctx, reqDTO.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.MysqlQueryResult{}, util.InternalError(err)
	}
	if !b || db.Config == nil {
		return util.MysqlQueryResult{}, util.InvalidArgsError()
	}
	tableName, sql, isExplain, err := command.ValidateMysqlSelectSql(reqDTO.Cmd)
	if err != nil {
		return util.MysqlQueryResult{}, util.NewBizErrWithMsg(apicode.OperationFailedErrCode, err.Error())
	}
	// 检查权限
	needCheckPerm := checkDbaPerm(reqDTO.Operator) != nil
	if needCheckPerm {
		hasPerm, err := mysqldbmd.ExistReadPerm(ctx, mysqldbmd.ExistReadPermReqDTO{
			Account:      reqDTO.Operator.Account,
			DbId:         reqDTO.DbId,
			AccessBase:   reqDTO.AccessBase,
			AccessTables: []string{tableName, "*"},
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.MysqlQueryResult{}, util.InternalError(err)
		}
		if !hasPerm {
			return util.MysqlQueryResult{}, util.UnauthorizedError()
		}
	}
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(db.Config.Data.ReadNode.Username),
		url.QueryEscape(db.Config.Data.ReadNode.Password),
		db.Config.Data.ReadNode.Host,
		reqDTO.AccessBase,
	)
	sql = strings.TrimSuffix(sql, ";")
	if !isExplain {
		sql = sql + " limit " + strconv.Itoa(reqDTO.Limit)
	}
	result, err := util.MysqlQuery(datasourceName, sql)
	if err != nil {
		return util.MysqlQueryResult{}, util.NewBizErrWithMsg(apicode.OperationFailedErrCode, err.Error())
	}
	return result, nil
}

// ApplyDataUpdate 提数据库修改单
func ApplyDataUpdate(ctx context.Context, reqDTO ApplyDataUpdateReqDTO) ([]command.ValidateUpdateResult, bool, error) {
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
	validateResults, allPass, err := command.ValidateMysqlUpdateSql(reqDTO.Cmd)
	if err != nil {
		return nil, false, util.NewBizErrWithMsg(apicode.OperationFailedErrCode, err.Error())
	}
	if allPass {
		// 插入数据库
		err = mysqldbmd.InsertDataUpdateApply(ctx, mysqldbmd.InsertDataUpdateApplyReqDTO{
			Account:          reqDTO.Operator.Account,
			DbId:             reqDTO.DbId,
			AccessBase:       reqDTO.AccessBase,
			UpdateCmd:        reqDTO.Cmd,
			ApplyReason:      reqDTO.ApplyReason,
			ApplyStatus:      mysqldbmd.PendingDataUpdateApplyStatus,
			ExecuteWhenApply: reqDTO.ExecuteWhenApply,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, false, util.InternalError(err)
		}
	}
	return validateResults, allPass, nil
}

// ExplainDataUpdate 数据库修改单的执行计划
func ExplainDataUpdate(ctx context.Context, reqDTO ExplainDataUpdateReqDTO) (string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	apply, b, err := mysqldbmd.GetDataUpdateApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b || !apply.ApplyStatus.IsUnExecuted() {
		return "", util.InvalidArgsError()
	}
	// 校验权限
	if apply.Account != reqDTO.Operator.Account && checkDbaPerm(reqDTO.Operator) != nil {
		return "", util.UnauthorizedError()
	}
	db, b, err := mysqldbmd.GetDbById(ctx, apply.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.InternalError(err)
	}
	if !b {
		return "", util.ThereHasBugErr()
	}
	validateResults, allPass, err := command.ValidateMysqlUpdateSql(apply.UpdateCmd)
	if err != nil || !allPass {
		return "", util.ThereHasBugErr()
	}
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(db.Config.Data.ReadNode.Username),
		url.QueryEscape(db.Config.Data.ReadNode.Password),
		db.Config.Data.ReadNode.Host,
		apply.AccessBase,
	)
	explainableSqls := make([]string, 0)
	unexplainableSqls := make([]string, 0)
	for _, result := range validateResults {
		if result.IsExplainable {
			explainableSqls = append(explainableSqls, result.Sql)
		} else {
			unexplainableSqls = append(unexplainableSqls, result.Sql)
		}
	}
	explainableSqls, _ = listutil.Map(explainableSqls, func(t string) (string, error) {
		return "explain " + t, nil
	})
	explainableResults, err := util.MysqlQueries(datasourceName, explainableSqls...)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", util.OperationFailedError()
	}
	if len(explainableSqls) != len(explainableResults) {
		return "", util.ThereHasBugErr()
	}
	ret := new(bytes.Buffer)
	separator := "\n---------------\n"
	for _, sql := range unexplainableSqls {
		ret.WriteString(strings.TrimSpace(sql) + "\n\n")
		ret.WriteString(separator)
	}
	for i := range explainableSqls {
		ret.WriteString(strings.TrimSpace(strings.TrimPrefix(explainableSqls[i], "explain ")) + "\n\n")
		result := explainableResults[i]
		if result.Err != nil {
			ret.WriteString("err: " + result.Err.Error() + "\n")
		} else {
			m := result.ToMap()
			if len(m) > 0 {
				ret.WriteString("rows: " + m[0]["rows"] + "\n")
				ret.WriteString("type: " + m[0]["type"] + "\n")
				ret.WriteString("possible_keys: " + m[0]["possible_keys"] + "\n")
			}
		}
		ret.WriteString(separator)
	}
	return ret.String(), nil
}

// ListDataUpdateApplyByDba 数据库修改审批单
func ListDataUpdateApplyByDba(ctx context.Context, reqDTO ListDataUpdateApplyByDbaReqDTO) ([]DataUpdateApplyDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	if err := checkDbaPerm(reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	applies, total, err := mysqldbmd.ListDataUpdateApply(ctx, mysqldbmd.ListDataUpdateApplyReqDTO{
		PageNum:     reqDTO.PageNum,
		PageSize:    10,
		DbId:        reqDTO.DbId,
		ApplyStatus: reqDTO.ApplyStatus,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, err := dataUpdateApplyMd2Dto(ctx, applies)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

// ListDataUpdateApplyByOperator 申请的数据库修改审批单
func ListDataUpdateApplyByOperator(ctx context.Context, reqDTO ListDataUpdateApplyByOperatorReqDTO) ([]DataUpdateApplyDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	applies, total, err := mysqldbmd.ListDataUpdateApply(ctx, mysqldbmd.ListDataUpdateApplyReqDTO{
		PageNum:     reqDTO.PageNum,
		PageSize:    10,
		Account:     reqDTO.Operator.Account,
		ApplyStatus: reqDTO.ApplyStatus,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, err := dataUpdateApplyMd2Dto(ctx, applies)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func dataUpdateApplyMd2Dto(ctx context.Context, applies []mysqldbmd.DataUpdateApply) ([]DataUpdateApplyDTO, error) {
	dbIdList, _ := listutil.Map(applies, func(t mysqldbmd.DataUpdateApply) (int64, error) {
		return t.DbId, nil
	})
	dbMap, err := getDbMap(ctx, dbIdList)
	if err != nil {
		return nil, err
	}
	data, _ := listutil.Map(applies, func(t mysqldbmd.DataUpdateApply) (DataUpdateApplyDTO, error) {
		ret := DataUpdateApplyDTO{
			Id:               t.Id,
			Account:          t.Account,
			DbId:             t.DbId,
			DbName:           dbMap[t.DbId].Name,
			AccessBase:       t.AccessBase,
			UpdateCmd:        t.UpdateCmd,
			ApplyStatus:      t.ApplyStatus,
			Auditor:          t.Auditor,
			Executor:         t.Executor,
			ApplyReason:      t.ApplyReason,
			DisagreeReason:   t.DisagreeReason,
			ExecuteLog:       t.ExecuteLog,
			ExecuteWhenApply: t.ExecuteWhenApply,
			Created:          t.Created,
			Updated:          t.Updated,
		}
		return ret, nil
	})
	return data, nil
}

// AgreeDataUpdateApply 同意修改单
func AgreeDataUpdateApply(ctx context.Context, reqDTO AgreeDbUpdateReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkDbaPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetDataUpdateApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.ApplyStatus != mysqldbmd.PendingDataUpdateApplyStatus {
		return util.InvalidArgsError()
	}
	_, err = mysqldbmd.UpdateDataUpdateApplyStatusWithAuditor(ctx, mysqldbmd.UpdateDataUpdateApplyStatusWithAuditorReqDTO{
		Id:        reqDTO.ApplyId,
		NewStatus: mysqldbmd.AgreeDataUpdateApplyStatus,
		OldStatus: order.ApplyStatus,
		Auditor:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DisagreeDataUpdateApply 不同意修改单
func DisagreeDataUpdateApply(ctx context.Context, reqDTO DisagreeDataUpdateApplyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkDbaPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetDataUpdateApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.ApplyStatus != mysqldbmd.PendingDataUpdateApplyStatus {
		return util.InvalidArgsError()
	}
	_, err = mysqldbmd.UpdateDataUpdateApplyStatusWithAuditor(ctx, mysqldbmd.UpdateDataUpdateApplyStatusWithAuditorReqDTO{
		Id:             reqDTO.ApplyId,
		NewStatus:      mysqldbmd.DisagreeDataUpdateApplyStatus,
		OldStatus:      order.ApplyStatus,
		Auditor:        reqDTO.Operator.Account,
		DisagreeReason: reqDTO.DisagreeReason,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// CancelDataUpdateApply 取消修改单
func CancelDataUpdateApply(ctx context.Context, reqDTO CancelDataUpdateApplyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	apply, b, err := mysqldbmd.GetDataUpdateApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || apply.ApplyStatus != mysqldbmd.PendingDataUpdateApplyStatus {
		return util.InvalidArgsError()
	}
	if apply.Account != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	_, err = mysqldbmd.UpdateDataUpdateApplyStatusWithAuditor(ctx, mysqldbmd.UpdateDataUpdateApplyStatusWithAuditorReqDTO{
		Id:        reqDTO.ApplyId,
		NewStatus: mysqldbmd.CanceledDataUpdateApplyStatus,
		OldStatus: apply.ApplyStatus,
		Auditor:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ExecuteDataUpdateApply 执行修改单
func ExecuteDataUpdateApply(ctx context.Context, reqDTO ExecuteDataUpdateApplyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if err := checkDbaPerm(reqDTO.Operator); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	apply, b, err := mysqldbmd.GetDataUpdateApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || !apply.ApplyStatus.IsExecutable() {
		return util.InvalidArgsError()
	}
	db, b, err := mysqldbmd.GetDbById(ctx, apply.DbId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || db.Config == nil {
		return util.ThereHasBugErr()
	}
	b, err = mysqldbmd.UpdateDataUpdateApplyStatusWithExecutor(ctx, mysqldbmd.UpdateDataUpdateApplyStatusWithExecutorReqDTO{
		Id:        reqDTO.ApplyId,
		NewStatus: mysqldbmd.ExecutedDataUpdateApplyStatus,
		OldStatus: apply.ApplyStatus,
		Executor:  reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return err
	}
	if b {
		go executeUpdateCmd(&apply, &db)
	}
	return nil
}

func executeUpdateCmd(apply *mysqldbmd.DataUpdateApply, db *mysqldbmd.Db) {
	logMsg := strings.Builder{}
	datasourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8",
		url.QueryEscape(db.Config.Data.WriteNode.Username),
		url.QueryEscape(db.Config.Data.WriteNode.Password),
		db.Config.Data.WriteNode.Host,
		apply.AccessBase,
	)
	results, err := util.MysqlExecute(datasourceName, apply.UpdateCmd)
	if err != nil {
		logMsg.WriteString(err.Error())
	} else {
		for _, result := range results {
			logMsg.WriteString("sql: " + strings.TrimSpace(result.Sql) + "\n")
			logMsg.WriteString("affectedRows: " + strconv.FormatInt(result.AffectedRows, 10) + "\n")
			if result.ErrMsg != "" {
				logMsg.WriteString("errMsg: " + result.ErrMsg + "\n")
			}
			logMsg.WriteString("\n\n")
		}
	}
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, err = mysqldbmd.UpdateDataUpdateApplyExecuteLog(ctx, apply.Id, logMsg.String())
	if err != nil {
		logger.Logger.Errorf("db update: %v executed with err: %v", apply.Id, err)
	}
}

// AskToExecuteDataUpdateApply 请求执行修改单
func AskToExecuteDataUpdateApply(ctx context.Context, reqDTO AskToExecuteDataUpdateApplyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	order, b, err := mysqldbmd.GetDataUpdateApplyById(ctx, reqDTO.ApplyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || order.ApplyStatus != mysqldbmd.AgreeDataUpdateApplyStatus {
		return util.InvalidArgsError()
	}
	if order.Account != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	_, err = mysqldbmd.UpdateDataUpdateApplyStatusWithAuditor(ctx, mysqldbmd.UpdateDataUpdateApplyStatusWithAuditorReqDTO{
		Id:        reqDTO.ApplyId,
		NewStatus: mysqldbmd.AskToExecuteDataUpdateApplyStatus,
		OldStatus: order.ApplyStatus,
		Auditor:   reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func getDbMap(ctx context.Context, idList []int64) (map[int64]mysqldbmd.Db, error) {
	idList = listutil.Distinct(idList...)
	ret := make(map[int64]mysqldbmd.Db, len(idList))
	dbs, err := mysqldbmd.BatchGetDbByIdList(ctx, idList, []string{"id", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	for _, db := range dbs {
		ret[db.Id] = db
	}
	return ret, nil
}

func checkDbaPerm(operator apisession.UserInfo) error {
	if operator.IsAdmin || operator.IsDba {
		return nil
	}
	return util.UnauthorizedError()
}
