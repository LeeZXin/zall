package usersrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

const (
	LoginSessionExpiry = 2 * time.Hour
)

func GetByAccount(ctx context.Context, account string) (usermd.UserInfo, bool) {
	user, b := getByAccount(ctx, account)
	return user.ToUserInfo(), b
}

func getByAccount(ctx context.Context, account string) (usermd.User, bool) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByAccount(ctx, account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	return user, b
}

func CheckAccountAndPassword(ctx context.Context, reqDTO CheckAccountAndPasswordReqDTO) (usermd.UserInfo, bool) {
	if err := reqDTO.IsValid(); err != nil {
		return usermd.UserInfo{}, false
	}
	user, b := getByAccount(ctx, reqDTO.Account)
	if !b {
		return usermd.UserInfo{}, false
	}
	// 检查是否被全局禁用或校验密码
	if user.IsProhibited || user.Password != util.EncryptUserPassword(reqDTO.Password) {
		return usermd.UserInfo{}, false
	}
	return user.ToUserInfo(), true
}

func Login(ctx context.Context, reqDTO LoginReqDTO) (session apisession.Session, err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		b    bool
		user usermd.User
	)
	user, b, err = usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.NewBizErr(apicode.DataNotExistsCode, i18n.UserNotFound)
		return
	}
	// 校验密码
	if user.Password != util.EncryptUserPassword(reqDTO.Password) {
		err = util.NewBizErr(apicode.WrongLoginPasswordCode, i18n.UserWrongPassword)
		return
	}
	// 检查是否被全局禁用
	if user.IsProhibited {
		err = util.UnauthorizedError()
		return
	}
	sessionStore := apisession.GetStore()
	// 删除原有的session
	err = sessionStore.DeleteByAccount(user.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	// 生成sessionId
	sessionId := apisession.GenSessionId()
	expireAt := time.Now().Add(LoginSessionExpiry).UnixMilli()
	session = apisession.Session{
		SessionId: sessionId,
		UserInfo: apisession.UserInfo{
			Account:      user.Account,
			Name:         user.Name,
			Email:        user.Email,
			IsProhibited: user.IsProhibited,
			AvatarUrl:    user.AvatarUrl,
			IsAdmin:      user.IsAdmin,
			IsDba:        user.IsDba,
		},
		ExpireAt: expireAt,
	}
	err = sessionStore.PutSession(session)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func Refresh(ctx context.Context, reqDTO RefreshReqDTO) (string, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", 0, err
	}
	sessionStore := apisession.GetStore()
	// 生成sessionId
	sessionId := apisession.GenSessionId()
	expireAt := time.Now().Add(LoginSessionExpiry).UnixMilli()
	err := sessionStore.PutSession(apisession.Session{
		SessionId: sessionId,
		UserInfo:  reqDTO.Operator,
		ExpireAt:  expireAt,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", 0, util.InternalError(err)
	}
	// 删除原有的session
	err = sessionStore.DeleteBySessionId(reqDTO.SessionId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	return sessionId, expireAt, nil
}

func Logout(ctx context.Context, reqDTO LogoutReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	sessionStore := apisession.GetStore()
	err := sessionStore.DeleteByAccount(reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// CreateUser 管理员创建用户
func CreateUser(ctx context.Context, reqDTO CreateUserReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 不是企业管理员
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	var b bool
	b, err = usermd.ExistByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.AlreadyExistsError()
		return
	}
	// 添加账号
	err = usermd.InsertUser(ctx, usermd.InsertUserReqDTO{
		Account:   reqDTO.Account,
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		Password:  util.EncryptUserPassword(reqDTO.Password),
		AvatarUrl: reqDTO.AvatarUrl,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// RegisterUser 注册用户
func RegisterUser(ctx context.Context, reqDTO RegisterUserReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 获取系统配置检查是否禁用注册功能
	sysCfg, err := cfgsrv.GetSysCfg(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if sysCfg.DisableSelfRegisterUser {
		err = util.UnauthorizedError()
		return
	}
	var (
		b         bool
		userCount int64
	)
	_, b, err = usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	if b {
		err = util.NewBizErr(apicode.UserAlreadyExistsCode, i18n.UserAlreadyExists)
		return
	}
	userCount, err = usermd.CountAllUsers(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 添加账号
	err = usermd.InsertUser(ctx, usermd.InsertUserReqDTO{
		Account:  reqDTO.Account,
		Name:     reqDTO.Name,
		Email:    reqDTO.Email,
		Password: util.EncryptUserPassword(reqDTO.Password),
		IsAdmin:  userCount == 0,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

// DeleteUser 注销用户
func DeleteUser(ctx context.Context, reqDTO DeleteUserReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	// 不是管理员
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	// 不能自己删自己
	if reqDTO.Operator.Account == reqDTO.Account {
		err = util.InvalidArgsError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 数据库删除用户
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 用户表
		_, err2 := usermd.DeleteUser(ctx, reqDTO.Account)
		if err2 != nil {
			return err2
		}
		// 权限表
		_, err2 = teammd.DeleteAllTeamUserByAccount(ctx, reqDTO.Account)
		return err2
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 删除用户登录状态
	err2 := apisession.GetStore().DeleteByAccount(reqDTO.Account)
	if err2 != nil {
		logger.Logger.WithContext(ctx).Error(err2)
	}
	return
}

// ListUser 展示用户列表
func ListUser(ctx context.Context, reqDTO ListUserReqDTO) ([]UserDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	// 只有系统管理员才能操作
	if !reqDTO.Operator.IsAdmin {
		return nil, 0, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	users, total, err := usermd.PageUser(ctx, usermd.PageUserReqDTO{
		Account:  reqDTO.Account,
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data := listutil.MapNe(users, func(t usermd.User) UserDTO {
		return UserDTO{
			Account:      t.Account,
			Name:         t.Name,
			Email:        t.Email,
			IsAdmin:      t.IsAdmin,
			IsProhibited: t.IsProhibited,
			AvatarUrl:    t.AvatarUrl,
			Created:      t.Created,
			IsDba:        t.IsDba,
		}
	})
	return data, total, nil
}

// UpdateUser 编辑用户
func UpdateUser(ctx context.Context, reqDTO UpdateUserReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 系统管理员或本人才能编辑user
	if !reqDTO.Operator.IsAdmin && reqDTO.Account != reqDTO.Operator.Account {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = usermd.UpdateUser(ctx, usermd.UpdateUserReqDTO{
		Account:   reqDTO.Account,
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		AvatarUrl: reqDTO.AvatarUrl,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	} else {
		// 修改信息后 需重新登录 登录session里有旧信息
		err2 := apisession.GetStore().DeleteByAccount(reqDTO.Account)
		if err2 != nil {
			logger.Logger.WithContext(ctx).Error(err2)
		}
	}
	return
}

// SetAdmin 设置系统管理员角色
func SetAdmin(ctx context.Context, reqDTO SetAdminReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才能设置系统管理员
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	// 系统管理员不能处理自己
	if reqDTO.Operator.Account == reqDTO.Account {
		err = util.InvalidArgsError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var b bool
	b, err = usermd.UpdateAdmin(ctx, usermd.UpdateAdminReqDTO{
		Account: reqDTO.Account,
		IsAdmin: reqDTO.IsAdmin,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	if !b {
		err = util.OperationFailedError()
		return
	}
	return
}

// SetDba 设置dba角色
func SetDba(ctx context.Context, reqDTO SetDbaReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 只有系统管理员才能设置系统管理员
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	b, err := usermd.UpdateDba(ctx, usermd.UpdateDbaReqDTO{
		Account: reqDTO.Account,
		IsDba:   reqDTO.IsDba,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.OperationFailedError()
	}
	return nil

}

// UpdatePassword 修改密码
func UpdatePassword(ctx context.Context, reqDTO UpdatePasswordReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 账号不存在
	if !b {
		err = util.ThereHasBugErr()
		return
	}
	// 原密码不正确
	if user.Password != util.EncryptUserPassword(reqDTO.Origin) {
		err = util.NewBizErr(apicode.OperationFailedErrCode, i18n.UserWrongOriginPassword)
		return
	}
	_, err = usermd.UpdatePassword(ctx, usermd.UpdatePasswordReqDTO{
		Account:  reqDTO.Operator.Account,
		Password: util.EncryptUserPassword(reqDTO.Password),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func SetProhibited(ctx context.Context, reqDTO SetProhibitedReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才能设置系统管理员
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	// 系统管理员不能处理自己
	if reqDTO.Operator.Account == reqDTO.Account {
		err = util.InvalidArgsError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var b bool
	b, err = usermd.UpdateProhibited(ctx, usermd.SetUserProhibitedReqDTO{
		Account:      reqDTO.Account,
		IsProhibited: reqDTO.IsProhibited,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	if !b {
		err = util.OperationFailedError()
		return
	}
	// 禁用用户下线登录token
	if reqDTO.IsProhibited {
		err2 := apisession.GetStore().DeleteByAccount(reqDTO.Account)
		if err2 != nil {
			logger.Logger.WithContext(ctx).Error(err2)
		}
	}
	return
}

// ListAllUser 所有用户列表
func ListAllUser(ctx context.Context, reqDTO ListAllUserReqDTO) ([]SimpleUserDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	users, err := usermd.ListAllUser(ctx, []string{"account", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(users, func(t usermd.User) (SimpleUserDTO, error) {
		return SimpleUserDTO{
			Account: t.Account,
			Name:    t.Name,
		}, nil
	})
}

// ResetPassword 重置密码
func ResetPassword(ctx context.Context, reqDTO ResetPasswordReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 只有系统管理员才能重置密码
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := usermd.UpdatePassword(ctx, usermd.UpdatePasswordReqDTO{
		Account:  reqDTO.Account,
		Password: util.EncryptUserPassword("123456"),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}
