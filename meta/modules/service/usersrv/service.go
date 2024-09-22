package usersrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/feishuapi"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/weworkapi"
	"github.com/LeeZXin/zall/thirdpart/modules/model/tpfeishumd"
	"github.com/LeeZXin/zall/thirdpart/modules/model/tpweworkmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"strings"
	"time"
)

const (
	LoginSessionExpiry = 24 * time.Hour
)

func GetUsersNameAndAvatar(ctx context.Context, accounts ...string) ([]util.User, error) {
	accounts = listutil.Distinct(accounts...)
	if len(accounts) == 0 {
		return []util.User{}, nil
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	users, err := usermd.ListUserByAccounts(ctx, accounts, []string{"account", "avatar_url", "name"})
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]usermd.User, len(users))
	for _, user := range users {
		userMap[user.Account] = user
	}
	return listutil.Map(accounts, func(t string) (util.User, error) {
		return util.User{
			Account:   t,
			Name:      userMap[t].Name,
			AvatarUrl: userMap[t].AvatarUrl,
		}, nil
	})
}

func GetUsersNameAndAvatarMap(ctx context.Context, accounts ...string) (map[string]util.User, error) {
	accounts = listutil.Distinct(accounts...)
	if len(accounts) == 0 {
		return map[string]util.User{}, nil
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	users, err := usermd.ListUserByAccounts(ctx, accounts, []string{"account", "avatar_url", "name"})
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]usermd.User, len(users))
	for _, user := range users {
		userMap[user.Account] = user
	}
	ret := make(map[string]util.User)
	for _, account := range accounts {
		ret[account] = util.User{
			Account:   account,
			Name:      userMap[account].Name,
			AvatarUrl: userMap[account].AvatarUrl,
		}
	}
	// 系统默认账号
	ret["system"] = util.User{
		Account: "system",
		Name:    "system",
	}
	return ret, nil
}

func GetUsersNameAndAvatarMapByEmails(ctx context.Context, emails ...string) (map[string]util.User, error) {
	emails = listutil.Distinct(emails...)
	if len(emails) == 0 {
		return map[string]util.User{}, nil
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	users, err := usermd.ListUserByEmails(ctx, emails, []string{"account", "avatar_url", "name", "email"})
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]usermd.User, len(users))
	for _, user := range users {
		userMap[user.Email] = user
	}
	ret := make(map[string]util.User)
	for _, email := range emails {
		ret[email] = util.User{
			Account:   userMap[email].Account,
			Name:      userMap[email].Name,
			AvatarUrl: userMap[email].AvatarUrl,
		}
	}
	ret["zgit@fake.local"] = util.User{
		Account: "ZGIT",
		Name:    "ZGIT",
	}
	return ret, nil
}

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

func Login(ctx context.Context, reqDTO LoginReqDTO) (apisession.Session, error) {
	if err := reqDTO.IsValid(); err != nil {
		return apisession.Session{}, err
	}
	hasA := reqDTO.A == "zsf"
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	loginCfg, err := cfgsrv.GetLoginCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	var (
		user usermd.User
		b    bool
	)
	// 如果开了后门且不允许账号密码登录 则只允许超级管理员登录
	if !hasA && !loginCfg.AccountPassword.IsEnabled {
		return apisession.Session{}, util.UnauthorizedError()
	}
	user, b, err = usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	if !b {
		return apisession.Session{}, util.NewBizErr(apicode.DataNotExistsCode, i18n.UserNotFound)
	}
	// 校验密码
	if user.Password != util.EncryptUserPassword(reqDTO.Password) {
		return apisession.Session{}, util.NewBizErr(apicode.WrongLoginPasswordCode, i18n.UserWrongPassword)
	}
	// 检查是否被全局禁用
	if user.IsProhibited {
		return apisession.Session{}, util.UnauthorizedError()
	}
	// 检查后门管理员
	if hasA && !user.IsAdmin {
		return apisession.Session{}, util.UnauthorizedError()
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
	session := apisession.Session{
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
		return apisession.Session{}, util.InternalError(err)
	}
	return session, nil
}

func WeworkLogin(ctx context.Context, reqDTO WeworkLoginReqDTO) (apisession.Session, error) {
	if err := reqDTO.IsValid(); err != nil {
		return apisession.Session{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	loginCfg, err := cfgsrv.GetLoginCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	if !loginCfg.Wework.IsEnabled {
		return apisession.Session{}, util.UnauthorizedError()
	}
	cfg := loginCfg.Wework
	if cfg.State != reqDTO.State {
		return apisession.Session{}, util.InvalidArgsError()
	}
	at, b, err := tpweworkmd.GetAccessTokenByCorpIdAndSecret(ctx, cfg.AppId, cfg.Secret)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	if !b {
		return apisession.Session{}, util.OperationFailedError()
	}
	account, err := weworkapi.GetAuthWeworkUserInfo(ctx, static.GetString("wework.auth.getUserInfoUrl"), at.Token, reqDTO.Code)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	user, b, err := usermd.GetByAccount(ctx, account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	if !b {
		return apisession.Session{}, util.NewBizErr(apicode.DataNotExistsCode, i18n.UserNotFound)
	}
	// 检查是否被全局禁用
	if user.IsProhibited {
		return apisession.Session{}, util.UnauthorizedError()
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
	session := apisession.Session{
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
		return apisession.Session{}, util.InternalError(err)
	}
	return session, nil
}

func FeishuLogin(ctx context.Context, reqDTO FeishuLoginReqDTO) (apisession.Session, error) {
	if err := reqDTO.IsValid(); err != nil {
		return apisession.Session{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	loginCfg, err := cfgsrv.GetLoginCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	if !loginCfg.Feishu.IsEnabled {
		return apisession.Session{}, util.UnauthorizedError()
	}
	cfg := loginCfg.Feishu
	if cfg.State != reqDTO.State {
		return apisession.Session{}, util.InvalidArgsError()
	}
	at, b, err := tpfeishumd.GetAccessTokenByAppIdAndSecret(ctx, cfg.ClientId, cfg.Secret)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	if !b {
		return apisession.Session{}, util.OperationFailedError()
	}
	ut, err := feishuapi.GetUserAccessToken(ctx, static.GetString("feishu.auth.getUserAccessTokenUrl"), at.Token, reqDTO.Code)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	userInfo, err := feishuapi.GetUserInfo(ctx, static.GetString("feishu.auth.getUserInfoUrl"), ut.AccessToken)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	user, b, err := usermd.GetByAccount(ctx, userInfo.UserId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return apisession.Session{}, util.InternalError(err)
	}
	if !b {
		return apisession.Session{}, util.NewBizErr(apicode.DataNotExistsCode, i18n.UserNotFound)
	}
	// 检查是否被全局禁用
	if user.IsProhibited {
		return apisession.Session{}, util.UnauthorizedError()
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
	session := apisession.Session{
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
		return apisession.Session{}, util.InternalError(err)
	}
	return session, nil
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
func CreateUser(ctx context.Context, reqDTO CreateUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 不是企业管理员
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	if strings.ToLower(reqDTO.Account) == "system" {
		return util.InvalidArgsError()
	}
	b, err := usermd.ExistByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	// 检查邮箱是否存在
	b, err = usermd.ExistByEmail(ctx, reqDTO.Email)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.NewBizErr(apicode.UserAlreadyExistsCode, i18n.UserEmailAlreadyExists)
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
		return util.InternalError(err)
	}
	return nil
}

// RegisterUser 注册用户
func RegisterUser(ctx context.Context, reqDTO RegisterUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 获取系统配置检查是否禁用注册功能
	sysCfg, err := cfgsrv.GetSysCfg(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if sysCfg.DisableSelfRegisterUser {
		return util.UnauthorizedError()
	}
	if strings.ToLower(reqDTO.Account) == "system" {
		return util.InvalidArgsError()
	}
	// 检查账号是否存在
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.NewBizErr(apicode.UserAlreadyExistsCode, i18n.UserAlreadyExists)
	}
	// 检查邮箱是否存在
	b, err = usermd.ExistByEmail(ctx, reqDTO.Email)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.NewBizErr(apicode.UserAlreadyExistsCode, i18n.UserEmailAlreadyExists)
	}
	userCount, err := usermd.CountAllUsers(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
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
		return util.InternalError(err)
	}
	return nil
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
func UpdateUser(ctx context.Context, reqDTO UpdateUserReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 系统管理员或本人才能编辑user
	if !reqDTO.Operator.IsAdmin &&
		reqDTO.Account != reqDTO.Operator.Account {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByEmail(ctx, reqDTO.Email)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b && user.Account != reqDTO.Account {
		return util.NewBizErr(apicode.UserAlreadyExistsCode, i18n.UserEmailAlreadyExists)
	}
	_, err = usermd.UpdateUser(ctx, usermd.UpdateUserReqDTO{
		Account:   reqDTO.Account,
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		AvatarUrl: reqDTO.AvatarUrl,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 修改信息后 需重新登录 登录session里有旧信息
	err2 := apisession.GetStore().DeleteByAccount(reqDTO.Account)
	if err2 != nil {
		logger.Logger.WithContext(ctx).Error(err2)
	}
	return nil
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
