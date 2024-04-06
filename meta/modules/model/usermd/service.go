package usermd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
)

var (
	validUserAccountPattern = regexp.MustCompile(`^\w{4,32}$`)
)

func IsAccountValid(account string) bool {
	return validUserAccountPattern.MatchString(account)
}

func InsertUser(ctx context.Context, reqDTO InsertUserReqDTO) (User, error) {
	u := User{
		Account:   reqDTO.Account,
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		Password:  reqDTO.Password,
		AvatarUrl: reqDTO.AvatarUrl,
		IsAdmin:   reqDTO.IsAdmin,
		RoleType:  reqDTO.RoleType,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&u)
	return u, err
}

func DeleteUser(ctx context.Context, account string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("account = ?", account).Delete(new(User))
	return rows == 1, err
}

func GetByAccount(ctx context.Context, account string) (User, bool, error) {
	var ret User
	b, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		Get(&ret)
	return ret, b, err
}

func ListUserByAccounts(ctx context.Context, accounts []string) ([]User, error) {
	ret := make([]User, 0)
	err := xormutil.MustGetXormSession(ctx).In("account", accounts).Find(&ret)
	return ret, err
}

func ListUser(ctx context.Context, reqDTO ListUserReqDTO) ([]User, error) {
	ret := make([]User, 0)
	session := xormutil.MustGetXormSession(ctx).OrderBy("id asc").Limit(reqDTO.Limit)
	if reqDTO.Account != "" {
		session.And("account like ?", "%"+reqDTO.Account+"%")
	}
	if reqDTO.Cursor > 0 {
		session.And("id > ?", reqDTO.Cursor)
	}
	err := session.Find(&ret)
	return ret, err
}

func UpdateUser(ctx context.Context, reqDTO UpdateUserReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		Limit(1).
		Cols("name", "email").
		Update(&User{
			Name:  reqDTO.Name,
			Email: reqDTO.Email,
		})
	return rows == 1, err
}

func UpdatePassword(ctx context.Context, reqDTO UpdatePasswordReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		Limit(1).
		Cols("password").
		Update(&User{
			Password: reqDTO.Password,
		})
	return rows == 1, err
}

func CountUser(ctx context.Context) (int64, error) {
	return xormutil.MustGetXormSession(ctx).Count(new(User))
}

func UpdateAdmin(ctx context.Context, reqDTO UpdateAdminReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		Cols("is_admin").
		Update(&User{
			IsAdmin: reqDTO.IsAdmin,
		})
	return rows == 1, err
}

func SetUserProhibited(ctx context.Context, reqDTO SetUserProhibitedReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		Cols("is_prohibited").
		Update(&User{
			IsProhibited: reqDTO.IsProhibited,
		})
	return rows == 1, err
}
