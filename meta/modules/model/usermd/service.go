package usermd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
)

func IsAccountValid(account string) bool {
	return regexp.MustCompile(`^\w{4,32}$`).MatchString(account)
}

func IsUsernameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsEmailValid(email string) bool {
	return regexp.MustCompile(`^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`).MatchString(email)
}

func IsPasswordValid(password string) bool {
	return regexp.MustCompile("\\S{6,255}").MatchString(password)
}

func InsertUser(ctx context.Context, reqDTO InsertUserReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&User{
			Account:   reqDTO.Account,
			Name:      reqDTO.Name,
			Email:     reqDTO.Email,
			Password:  reqDTO.Password,
			AvatarUrl: reqDTO.AvatarUrl,
			IsAdmin:   reqDTO.IsAdmin,
			IsDba:     reqDTO.IsDba,
		})
	return err
}

func DeleteUser(ctx context.Context, account string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		Delete(new(User))
	return rows == 1, err
}

func GetByAccount(ctx context.Context, account string) (User, bool, error) {
	var ret User
	b, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		Get(&ret)
	return ret, b, err
}

func ExistByAccount(ctx context.Context, account string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		Exist(new(User))
}

func ListUserByAccounts(ctx context.Context, accounts, cols []string) ([]User, error) {
	ret := make([]User, 0)
	session := xormutil.MustGetXormSession(ctx).In("account", accounts)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func CountUserByAccounts(ctx context.Context, accounts []string) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		In("account", accounts).
		Count(new(User))
}

func CountAllUsers(ctx context.Context) (int64, error) {
	return xormutil.MustGetXormSession(ctx).Count(new(User))
}

func PageUser(ctx context.Context, reqDTO PageUserReqDTO) ([]User, int64, error) {
	ret := make([]User, 0)
	session := xormutil.MustGetXormSession(ctx)
	if reqDTO.Account != "" {
		session.And("account like ?", reqDTO.Account+"%")
	}
	total, err := session.
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		Desc("id").
		FindAndCount(&ret)
	return ret, total, err
}

func ListAllUser(ctx context.Context, cols []string) ([]User, error) {
	ret := make([]User, 0)
	err := xormutil.MustGetXormSession(ctx).Cols(cols...).Find(&ret)
	return ret, err
}

func UpdateUser(ctx context.Context, reqDTO UpdateUserReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		Cols("name", "email", "avatar_url").
		Update(&User{
			Name:      reqDTO.Name,
			Email:     reqDTO.Email,
			AvatarUrl: reqDTO.AvatarUrl,
		})
	return rows == 1, err
}

func UpdatePassword(ctx context.Context, reqDTO UpdatePasswordReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		Cols("password").
		Update(&User{
			Password: reqDTO.Password,
		})
	return rows == 1, err
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

func UpdateDba(ctx context.Context, reqDTO UpdateDbaReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		Cols("is_dba").
		Update(&User{
			IsDba: reqDTO.IsDba,
		})
	return rows == 1, err
}

func UpdateProhibited(ctx context.Context, reqDTO SetUserProhibitedReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", reqDTO.Account).
		Cols("is_prohibited").
		Update(&User{
			IsProhibited: reqDTO.IsProhibited,
		})
	return rows == 1, err
}
