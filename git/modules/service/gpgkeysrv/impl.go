package gpgkeysrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/gpgkeymd"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/git/signature"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/keybase/go-crypto/openpgp"
	"strings"
)

type innerImpl struct{}

// ListValidByAccount 根据账号获取未过期的
func (*innerImpl) ListValidByAccount(ctx context.Context, account string) []gpgkeymd.GpgKey {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	ret, err := gpgkeymd.ListValidByAccount(ctx, account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return []gpgkeymd.GpgKey{}
	}
	return ret
}

// GetByKeyId 根据keyId获取gpg密钥
func (*innerImpl) GetByKeyId(ctx context.Context, keyId string) (gpgkeymd.GpgKey, bool) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	ret, b, err := gpgkeymd.GetByKeyId(ctx, keyId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return gpgkeymd.GpgKey{}, false
	}
	return ret, b
}

type outerImpl struct{}

/*
CreateGpgKey 创建gpg密钥
`echo '%s' | gpg -a --default-key %s --detach-sig`, token, pubKeyId
*/
func (s *outerImpl) CreateGpgKey(ctx context.Context, reqDTO CreateGpgKeyReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 转换key
	entityList, err := signature.ConvertArmoredGPGKeyString(reqDTO.Content)
	if err != nil || len(entityList) == 0 {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.GpgKeyFormatError)
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	insertReqs := make([]gpgkeymd.InsertGpgKeyReqDTO, 0, len(entityList))
	keyIds := make([]string, 0, len(entityList))
	for _, entity := range entityList {
		keyIds = append(keyIds, entity.PrimaryKey.KeyIdString())
		insertReqs = append(insertReqs, parseGpgKeyInsertReq(entity, reqDTO.Name, reqDTO.Operator.Account))
	}
	// 检查重复
	var count int64
	count, err = gpgkeymd.CountByKeyIds(ctx, keyIds)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if count > 0 {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.GpgKeyAlreadyExists)
		return
	}
	// 插入
	err = gpgkeymd.BatchInsertGpgKeys(ctx, insertReqs)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func parseGpgKeyInsertReq(entity *openpgp.Entity, name, account string) gpgkeymd.InsertGpgKeyReqDTO {
	emailList := make([]string, 0)
	for _, identity := range entity.Identities {
		if identity.Revocation != nil {
			continue
		}
		if identity.UserId != nil {
			emailList = append(emailList, identity.UserId.Email)
		}
	}
	subKeys := make([]gpgkeymd.InsertGpgSubKeyReqDTO, 0)
	for _, subkey := range entity.Subkeys {
		if subkey.PublicKey == nil {
			continue
		}
		subContent, _ := signature.Base64EncGPGPubKey(subkey.PublicKey)
		subKeys = append(subKeys, gpgkeymd.InsertGpgSubKeyReqDTO{
			KeyId:   subkey.PublicKey.KeyIdString(),
			Content: subContent,
		})
	}
	content, _ := signature.Base64EncGPGPubKey(entity.PrimaryKey)
	expired := signature.GetGPGKeyExpiryTime(entity)
	return gpgkeymd.InsertGpgKeyReqDTO{
		Name:    name,
		Account: account,
		KeyId:   entity.PrimaryKey.KeyIdString(),
		Content: content,
		SubKeys: subKeys,
		Email:   strings.Join(emailList, ";"),
		Expired: expired,
	}
}

// DeleteGpgKey 删除密钥
func (s *outerImpl) DeleteGpgKey(ctx context.Context, reqDTO DeleteGpgKeyReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验是否存在
	key, b, err := gpgkeymd.GetById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b || key.Account != reqDTO.Operator.Account {
		err = util.UnauthorizedError()
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

// ListGpgKey 密钥列表
func (s *outerImpl) ListGpgKey(ctx context.Context, reqDTO ListGpgKeyReqDTO) ([]GpgKeyDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	keys, err := gpgkeymd.ListAllByAccount(
		ctx,
		reqDTO.Operator.Account,
		[]string{"id", "name", "key_id", "email", "expired", "created", "sub_keys"},
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(keys, func(t gpgkeymd.GpgKey) (GpgKeyDTO, error) {
		subKeys, _ := listutil.Map(t.SubKeys, func(t gpgkeymd.GpgSubKey) (string, error) {
			return t.KeyId, nil
		})
		return GpgKeyDTO{
			Id:      t.Id,
			Name:    t.Name,
			KeyId:   t.KeyId,
			Email:   t.Email,
			Created: t.Created,
			Expired: t.Expired,
			SubKeys: strings.Join(subKeys, ";"),
		}, nil
	})
}
