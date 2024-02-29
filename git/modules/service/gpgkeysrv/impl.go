package gpgkeysrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/gpgkeymd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/git/signature"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"github.com/keybase/go-crypto/openpgp"
	"github.com/patrickmn/go-cache"
	"time"
)

type innerImpl struct{}

func (*innerImpl) GetVerifiedByAccount(ctx context.Context, account string) ([]openpgp.EntityList, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	keys, err := gpgkeymd.GetValidByAccount(ctx, account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret := make([]openpgp.EntityList, 0, len(keys))
	for _, key := range keys {
		entity, err := signature.ConvertArmoredGPGKeyString(key.Content)
		if err == nil {
			ret = append(ret, entity)
		}
	}
	return ret, nil
}

type outerImpl struct {
	tokenCache *cache.Cache
}

func (s *outerImpl) InsertGpgKey(ctx context.Context, reqDTO InsertGpgKeyReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.GpgSrvKeysVO.InsertGpgKey),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	token, b := s.tokenCache.Get(reqDTO.Operator.Account)
	if !b {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.GpgTokenExpiredError)
		return
	}
	// 转换key
	entityList, err := signature.ConvertArmoredGPGKeyString(reqDTO.Content)
	if err != nil || len(entityList) == 0 {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.GpgKeyFormatError)
		return
	}
	// 校验token signature
	_, err = signature.CheckArmoredDetachedSignature(entityList, token.(string), reqDTO.Signature)
	if err != nil {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.GpgKeyVerifiedFailedError)
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	keyReqDTO := parseGpgKey(entityList[0], reqDTO.Content, reqDTO.Name, reqDTO.Operator.Account)
	// 检查重复
	b, err = gpgkeymd.ExistsByPubKeyId(ctx, keyReqDTO.PubKeyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.GpgKeyAlreadyExists)
		return
	}
	// 插入
	err = gpgkeymd.InsertGpgKeys(ctx, keyReqDTO)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func parseGpgKey(entity *openpgp.Entity, content, name, account string) gpgkeymd.InsertGpgKeyReqDTO {
	pubKeyId := entity.PrimaryKey.KeyIdString()
	expiry := signature.GetGPGKeyExpiryTime(entity).UnixMilli()
	emailList := make([]string, 0)
	for _, identity := range entity.Identities {
		if identity.Revocation != nil {
			continue
		}
		if identity.UserId != nil {
			emailList = append(emailList, identity.UserId.Email)
		}
	}
	return gpgkeymd.InsertGpgKeyReqDTO{
		Name:      name,
		Account:   account,
		PubKeyId:  pubKeyId,
		Content:   content,
		EmailList: emailList,
		Expiry:    expiry,
	}
}

func (s *outerImpl) GetToken(_ context.Context, reqDTO GetTokenReqDTO) (string, []string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", nil, err
	}
	// 转换key
	entityList, err := signature.ConvertArmoredGPGKeyString(reqDTO.Content)
	if err != nil || len(entityList) == 0 {
		return "", nil, util.NewBizErr(apicode.InvalidArgsCode, i18n.GpgKeyFormatError)
	}
	token := signature.GetToken(signature.User{
		Account: reqDTO.Operator.Account,
		Email:   reqDTO.Operator.Email,
	})
	// 设置十分钟有效期
	s.tokenCache.Set(reqDTO.Operator.Account, token, 10*time.Minute)
	pubKeyId := entityList[0].PrimaryKey.KeyIdString()
	return "", []string{
		i18n.GetByKey(i18n.GpgVerifyGuide),
		fmt.Sprintf(`echo '%s' | gpg -a --default-key %s --detach-sig`, token, pubKeyId),
	}, nil
}

func (s *outerImpl) DeleteGpgKey(ctx context.Context, reqDTO DeleteGpgKeyReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.GpgSrvKeysVO.DeleteGpgKey),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验是否存在
	key, b, err := gpgkeymd.GetById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b || key.Account != reqDTO.Operator.Account {
		err = util.InvalidArgsError()
		return
	}
	_, err = gpgkeymd.DeleteById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (s *outerImpl) ListGpgKey(ctx context.Context, reqDTO ListGpgKeyReqDTO) ([]GpgKeyDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	keys, err := gpgkeymd.GetAllSimpleByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(keys, func(t gpgkeymd.SimpleGpgKeyDTO) (GpgKeyDTO, error) {
		return GpgKeyDTO{
			Id:         t.Id,
			Name:       t.Name,
			PubKeyId:   t.PubKeyId,
			EmailList:  t.EmailList,
			ExpireTime: time.UnixMilli(t.Expiry),
		}, nil
	})
}
