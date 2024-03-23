package usersrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	LoginSessionExpiry = 2 * time.Hour
)

type innerImpl struct {
	userCache *cache.Cache
}

func (s *innerImpl) GetByAccount(ctx context.Context, account string) (usermd.UserInfo, bool) {
	user, b := s.getByAccount(ctx, account)
	return user.ToUserInfo(), b
}

func (s *innerImpl) getByAccount(ctx context.Context, account string) (usermd.User, bool) {
	v, b := s.userCache.Get(account)
	if b {
		u := v.(usermd.User)
		return u, u.Account != ""
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByAccount(ctx, account)
	if err != nil || !b {
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		s.userCache.Set(account, user, time.Second)
	} else {
		s.userCache.Set(account, user, time.Minute)
	}
	return user, b
}

func (s *innerImpl) CheckAccountAndPassword(ctx context.Context, reqDTO CheckAccountAndPasswordReqDTO) (usermd.UserInfo, bool) {
	if err := reqDTO.IsValid(); err != nil {
		return usermd.UserInfo{}, false
	}
	user, b := s.getByAccount(ctx, reqDTO.Account)
	if !b {
		return usermd.UserInfo{}, false
	}
	// 检查是否被全局禁用或校验密码
	if user.IsProhibited || user.Password != util.EncryptUserPassword(reqDTO.Password) {
		return usermd.UserInfo{}, false
	}
	return user.ToUserInfo(), true
}

type outerImpl struct{}

func (s *outerImpl) Login(ctx context.Context, reqDTO LoginReqDTO) (sessionId string, expireAt int64, err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Account,
			OpDesc:     i18n.GetByKey(i18n.UserSrvKeysVO.Login),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	user, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
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
	sessionId = apisession.GenSessionId()
	expireAt = time.Now().Add(LoginSessionExpiry).UnixMilli()
	err = sessionStore.PutSession(apisession.Session{
		SessionId: sessionId,
		UserInfo: apisession.UserInfo{
			Account:      user.Account,
			Name:         user.Name,
			Email:        user.Email,
			IsProhibited: user.IsProhibited,
			AvatarUrl:    user.AvatarUrl,
			IsAdmin:      user.IsAdmin,
		},
		ExpireAt: expireAt,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (s *outerImpl) Refresh(ctx context.Context, reqDTO RefreshReqDTO) (string, int64, error) {
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

func (*outerImpl) LoginOut(ctx context.Context, reqDTO LoginOutReqDTO) error {
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

func (s *outerImpl) InsertUser(ctx context.Context, reqDTO InsertUserReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.UserSrvKeysVO.Login),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
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
	// 添加账号
	_, err = usermd.InsertUser(ctx, usermd.InsertUserReqDTO{
		Account:   reqDTO.Account,
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		Password:  util.EncryptUserPassword(reqDTO.Password),
		AvatarUrl: reqDTO.AvatarUrl,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return nil
}

func (*outerImpl) RegisterUser(ctx context.Context, reqDTO RegisterUserReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Account,
			OpDesc:     i18n.GetByKey(i18n.UserSrvKeysVO.RegisterUser),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 获取系统配置检查是否禁用注册功能
	sysCfg, b := cfgsrv.Inner.GetSysCfg(ctx)
	if !b || sysCfg.DisableSelfRegisterUser {
		err = util.UnauthorizedError()
		return
	}
	_, b, err = usermd.GetByAccount(ctx, reqDTO.Account)
	if b {
		err = util.NewBizErr(apicode.UserAlreadyExistsCode, i18n.UserAlreadyExists)
		return
	}
	// 计算企业里面的用户数量，否则第一个注册者就是企业管理员
	countUser, err := usermd.CountUser(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 添加账号
	_, err = usermd.InsertUser(ctx, usermd.InsertUserReqDTO{
		Account:   reqDTO.Account,
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		Password:  util.EncryptUserPassword(reqDTO.Password),
		AvatarUrl: reqDTO.AvatarUrl,
		IsAdmin:   countUser == 0,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

// DeleteUser 注销用户
func (*outerImpl) DeleteUser(ctx context.Context, reqDTO DeleteUserReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.UserSrvKeysVO.DeleteUser),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	// 不是管理员
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 数据库删除用户
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 用户表
		_, err := usermd.DeleteUser(ctx, reqDTO.Operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		// 权限表
		_, err = teammd.DeleteAllTeamUserByAccount(ctx, reqDTO.Account)
		return err
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 删除用户登录状态
	apisession.GetStore().DeleteByAccount(reqDTO.Operator.Account)
	return
}

// ListUser 展示用户列表
func (*outerImpl) ListUser(ctx context.Context, reqDTO ListUserReqDTO) ([]UserDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	// 只有系统管理员才能操作
	if !reqDTO.Operator.IsAdmin {
		return nil, 0, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	userList, err := usermd.ListUser(ctx, usermd.ListUserReqDTO{
		Account: reqDTO.Account,
		Cursor:  reqDTO.Cursor,
		Limit:   reqDTO.Limit,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(userList, func(t usermd.User) (UserDTO, error) {
		return UserDTO{
			Account:      t.Account,
			Name:         t.Name,
			Email:        t.Email,
			IsAdmin:      t.IsAdmin,
			IsProhibited: t.IsProhibited,
			AvatarUrl:    t.AvatarUrl,
			Created:      t.Created,
			Updated:      t.Updated,
		}, nil
	})
	var next int64 = 0
	if len(userList) == reqDTO.Limit {
		next = userList[len(userList)-1].Id
	}
	return data, next, nil
}

func (*outerImpl) UpdateUser(ctx context.Context, reqDTO UpdateUserReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.UserSrvKeysVO.UpdateUser),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 系统管理员或本人才能编辑user
	if reqDTO.Account != reqDTO.Operator.Account {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 账号不存在
	if !b {
		err = util.InvalidArgsError()
		return
	}
	if _, err = usermd.UpdateUser(ctx, usermd.UpdateUserReqDTO{
		Account: reqDTO.Account,
		Name:    reqDTO.Name,
		Email:   reqDTO.Email,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) UpdateAdmin(ctx context.Context, reqDTO UpdateAdminReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.UserSrvKeysVO.UpdateAdmin),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
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
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 账号不存在
	if !b {
		err = util.InvalidArgsError()
		return
	}
	if _, err = usermd.UpdateAdmin(ctx, usermd.UpdateAdminReqDTO{
		Account: reqDTO.Account,
		IsAdmin: reqDTO.IsAdmin,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) UpdatePassword(ctx context.Context, reqDTO UpdatePasswordReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.UserSrvKeysVO.UpdatePassword),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 账号不存在
	if !b {
		err = util.InvalidArgsError()
		return
	}
	if _, err = usermd.UpdatePassword(ctx, usermd.UpdatePasswordReqDTO{
		Account:  reqDTO.Operator.Account,
		Password: reqDTO.Password,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (s *outerImpl) SetProhibited(ctx context.Context, reqDTO SetProhibitedReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.UserSrvKeysVO.SetProhibited),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
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
	_, b, err := usermd.GetByAccount(ctx, reqDTO.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 账号不存在
	if !b {
		return util.InvalidArgsError()
	}
	if _, err = usermd.SetUserProhibited(ctx, usermd.SetUserProhibitedReqDTO{
		Account:      reqDTO.Account,
		IsProhibited: reqDTO.IsProhibited,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}
