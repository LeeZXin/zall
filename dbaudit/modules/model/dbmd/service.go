package dbmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsDbNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsBaseNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsTableNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsUsernameValid(name string) bool {
	return len(name) <= 32
}

func IsPasswordValid(password string) bool {
	return len(password) <= 32
}

func IsReasonValid(reason string) bool {
	return len(reason) <= 255
}

func InsertDb(ctx context.Context, reqDTO InsertDbReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Db{
		Name:     reqDTO.Name,
		DbHost:   reqDTO.DbHost,
		Username: reqDTO.Username,
		Password: reqDTO.Password,
		DbType:   reqDTO.DbType,
	})
	return err
}

func GetDbById(ctx context.Context, id int64) (Db, bool, error) {
	var ret Db
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func BatchGetDbByIdList(ctx context.Context, idList []int64) ([]Db, error) {
	ret := make([]Db, 0)
	err := xormutil.MustGetXormSession(ctx).In("id", idList).Find(&ret)
	return ret, err
}

func ListDb(ctx context.Context) ([]Db, error) {
	ret := make([]Db, 0)
	err := xormutil.MustGetXormSession(ctx).Find(&ret)
	return ret, err
}

func DeleteDbById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(Db))
	return rows == 1, err
}

func UpdateDb(ctx context.Context, reqDTO UpdateDbReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "db_host", "username", "password").
		Update(&Db{
			Name:     reqDTO.Name,
			DbHost:   reqDTO.DbHost,
			Username: reqDTO.Username,
			Password: reqDTO.Password,
		})
	return rows == 1, err
}

func InsertApprovalOrder(ctx context.Context, reqDTO InsertApprovalOrderReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&ApprovalOrder{
			Account:      reqDTO.Account,
			DbId:         reqDTO.DbId,
			AccessBase:   reqDTO.AccessBase,
			AccessTables: reqDTO.AccessTables,
			PermType:     reqDTO.PermType,
			OrderStatus:  reqDTO.OrderStatus,
			ExpireDay:    reqDTO.ExpireDay,
			Reason:       reqDTO.Reason,
		})
	return err
}

func ListApprovalOrder(ctx context.Context, reqDTO ListApprovalOrderReqDTO) ([]ApprovalOrder, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("order_status = ?", reqDTO.OrderStatus)
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	ret := make([]ApprovalOrder, 0)
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}

func GetApprovalOrderById(ctx context.Context, id int64) (ApprovalOrder, bool, error) {
	var ret ApprovalOrder
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func UpdateApprovalOrderStatus(ctx context.Context, reqDTO UpdateApprovalOrderStatusReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		And("order_status = ?", reqDTO.OldStatus).
		Cols("order_status", "auditor").
		Update(&ApprovalOrder{
			OrderStatus: reqDTO.NewStatus,
			Auditor:     reqDTO.Auditor,
		})
	return rows == 1, err
}

func InsertPerm(ctx context.Context, reqDTO InsertPermReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Perm{
			Account:     reqDTO.Account,
			DbId:        reqDTO.DbId,
			AccessBase:  reqDTO.AccessBase,
			AccessTable: reqDTO.AccessTable,
			PermType:    reqDTO.PermType,
			Expired:     reqDTO.Expired,
		})
	return err
}

func BatchInsertPerm(ctx context.Context, reqDTOs []InsertPermReqDTO) error {
	perms, _ := listutil.Map(reqDTOs, func(reqDTO InsertPermReqDTO) (Perm, error) {
		return Perm{
			Account:     reqDTO.Account,
			DbId:        reqDTO.DbId,
			AccessBase:  reqDTO.AccessBase,
			AccessTable: reqDTO.AccessTable,
			PermType:    reqDTO.PermType,
			Expired:     reqDTO.Expired,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(perms)
	return err
}

func DeletePerm(ctx context.Context, id int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(Perm))
	return err
}

func ListPerm(ctx context.Context, reqDTO ListPermReqDTO) ([]Perm, error) {
	ret := make([]Perm, 0)
	session := xormutil.MustGetXormSession(ctx).Where("account = ?", reqDTO.Account)
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}

func SearchPerm(ctx context.Context, reqDTO SearchPermReqDTO) ([]Perm, error) {
	ret := make([]Perm, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		And("db_id = ?", reqDTO.DbId).
		And("access_base = ?", reqDTO.AccessBase).
		In("access_table", reqDTO.AccessTables).
		Find(&ret)
	return ret, err
}
