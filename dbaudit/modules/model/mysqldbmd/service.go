package mysqldbmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func IsDbNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsBaseNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 128
}

func IsTableNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 128
}

func IsReasonValid(reason string) bool {
	return len(reason) <= 255
}

func IsUpdateApplyNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertDb(ctx context.Context, reqDTO InsertDbReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Db{
			Name: reqDTO.Name,
			Config: &xormutil.Conversion[Config]{
				Data: reqDTO.Config,
			},
		})
	return err
}

func GetDbById(ctx context.Context, id int64) (Db, bool, error) {
	var ret Db
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func ListDb(ctx context.Context, cols []string) ([]Db, error) {
	ret := make([]Db, 0)
	session := xormutil.MustGetXormSession(ctx)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func PageDb(ctx context.Context, reqDTO PageDbReqDTO) ([]Db, int64, error) {
	session := xormutil.MustGetXormSession(ctx).
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize)
	if reqDTO.Name != "" {
		session.And("name like ?", reqDTO.Name+"%")
	}
	ret := make([]Db, 0)
	total, err := session.FindAndCount(&ret)
	return ret, total, err
}

func DeleteDbById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(Db))
	return rows == 1, err
}

func UpdateDb(ctx context.Context, reqDTO UpdateDbReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "config").
		Update(&Db{
			Name: reqDTO.Name,
			Config: &xormutil.Conversion[Config]{
				Data: reqDTO.Config,
			},
		})
	return rows == 1, err
}

func InsertReadPermApply(ctx context.Context, reqDTO InsertReadPermApplyReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&ReadPermApply{
			Account:      reqDTO.Account,
			DbId:         reqDTO.DbId,
			DbName:       reqDTO.DbName,
			AccessBase:   reqDTO.AccessBase,
			AccessTables: reqDTO.AccessTables,
			ApplyStatus:  reqDTO.OrderStatus,
			ExpireDay:    reqDTO.ExpireDay,
			ApplyReason:  reqDTO.ApplyReason,
		})
	return err
}

func ListReadPermApply(ctx context.Context, reqDTO ListReadPermApplyReqDTO) ([]ReadPermApply, int64, error) {
	session := xormutil.MustGetXormSession(ctx)
	if reqDTO.ApplyStatus > 0 {
		session.And("apply_status = ?", reqDTO.ApplyStatus)
	}
	if reqDTO.Account != "" {
		session.And("account = ?", reqDTO.Account)
	}
	if reqDTO.DbId > 0 {
		session.And("db_id = ?", reqDTO.DbId)
	}
	ret := make([]ReadPermApply, 0)
	total, err := session.
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		Desc("id").
		FindAndCount(&ret)
	return ret, total, err
}

func ListDataUpdateApply(ctx context.Context, reqDTO ListDataUpdateApplyReqDTO) ([]DataUpdateApply, int64, error) {
	session := xormutil.MustGetXormSession(ctx)
	if reqDTO.OrderStatus > 0 {
		session.And("order_status = ?", reqDTO.OrderStatus)
	}
	if reqDTO.Account != "" {
		session.And("account = ?", reqDTO.Account)
	}
	ret := make([]DataUpdateApply, 0)
	total, err := session.Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).FindAndCount(&ret)
	return ret, total, err
}

func GetReadPermApplyById(ctx context.Context, id int64) (ReadPermApply, bool, error) {
	var ret ReadPermApply
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func GetDataUpdateApplyById(ctx context.Context, id int64) (DataUpdateApply, bool, error) {
	var ret DataUpdateApply
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func UpdateReadPermApplyStatus(ctx context.Context, reqDTO UpdateReadPermApplyStatusReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		And("apply_status = ?", reqDTO.OldStatus).
		Cols("apply_status", "auditor", "disagree_reason").
		Update(&ReadPermApply{
			DisagreeReason: reqDTO.DisagreeReason,
			ApplyStatus:    reqDTO.NewStatus,
			Auditor:        reqDTO.Auditor,
		})
	return rows == 1, err
}

func UpdateDataUpdateApplyStatus(ctx context.Context, reqDTO UpdateDataUpdateApplyStatusReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		And("apply_status = ?", reqDTO.OldStatus).
		Cols("apply_status", "auditor", "disagree_reason").
		Update(&DataUpdateApply{
			DisagreeReason: reqDTO.DisagreeReason,
			ApplyStatus:    reqDTO.NewStatus,
			Auditor:        reqDTO.Auditor,
		})
	return rows == 1, err
}

func UpdateDataUpdateApplyExecuteLog(ctx context.Context, id int64, log string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("execute_log").
		Update(&DataUpdateApply{
			ExecuteLog: log,
		})
	return rows == 1, err
}

func BatchGetDbByIdList(ctx context.Context, idList []int64, cols []string) ([]Db, error) {
	ret := make([]Db, 0, len(idList))
	session := xormutil.MustGetXormSession(ctx).In("id", idList)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func BatchInsertReadPerm(ctx context.Context, reqDTOs []InsertReadPermReqDTO) error {
	perms, _ := listutil.Map(reqDTOs, func(reqDTO InsertReadPermReqDTO) (ReadPerm, error) {
		return ReadPerm{
			Account:     reqDTO.Account,
			DbId:        reqDTO.DbId,
			AccessBase:  reqDTO.AccessBase,
			AccessTable: reqDTO.AccessTable,
			ApplyId:     reqDTO.ApplyId,
			Expired:     reqDTO.Expired,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(perms)
	return err
}

func DeleteReadPermById(ctx context.Context, id int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(ReadPerm))
	return err
}

func ListReadPermByAccount(ctx context.Context, reqDTO ListReadPermByAccountReqDTO) ([]ReadPerm, error) {
	ret := make([]ReadPerm, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account)
	if reqDTO.DbId > 0 {
		session.And("db_id = ?", reqDTO.DbId)
	}
	if len(reqDTO.AccessBase) > 0 {
		session.And("access_base = ?", reqDTO.AccessBase)
	}
	if len(reqDTO.Cols) > 0 {
		session.Cols(reqDTO.Cols...)
	}
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}

func PageReadPerm(ctx context.Context, reqDTO PageReadPermReqDTO) ([]ReadPerm, int64, error) {
	ret := make([]ReadPerm, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize)
	if reqDTO.DbId > 0 {
		session.And("db_id = ?", reqDTO.DbId)
	}
	total, err := session.FindAndCount(&ret)
	return ret, total, err
}

func DeleteExpiredReadPermByAccount(ctx context.Context, account string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		And("expired < ?", time.Now().Format(time.DateTime)).
		Delete(new(ReadPerm))
	return err
}

func ExistReadPerm(ctx context.Context, reqDTO ExistReadPermReqDTO) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		And("db_id = ?", reqDTO.DbId).
		And("access_base = ?", reqDTO.AccessBase).
		In("access_table", reqDTO.AccessTables).
		Exist(new(ReadPerm))
}

func InsertUpdateApprovalOrder(ctx context.Context, reqDTO InsertUpdateApprovalOrderReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&DataUpdateApply{
			Name:        reqDTO.Name,
			Account:     reqDTO.Account,
			DbId:        reqDTO.DbId,
			AccessBase:  reqDTO.AccessBase,
			UpdateCmd:   reqDTO.UpdateCmd,
			ApplyStatus: reqDTO.OrderStatus,
		})
	return err
}
