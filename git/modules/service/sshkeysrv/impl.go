package sshkeysrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/sshkeymd"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/git/signature"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/patrickmn/go-cache"
	gossh "golang.org/x/crypto/ssh"
	"strings"
	"time"
)

type innerImpl struct {
	sshKeyCache *cache.Cache
}

func (s *innerImpl) GetAccountByFingerprint(ctx context.Context, fingerprint string) (string, bool, error) {
	sshKey, b := s.sshKeyCache.Get(fingerprint)
	if b {
		ret := sshKey.(string)
		if ret == "" {
			return ret, false, nil
		}
		return ret, true, nil
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	account, b, err := sshkeymd.GetAccountByFingerprint(ctx, fingerprint)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return "", false, util.InternalError(err)
	}
	if !b {
		// 空缓存
		s.sshKeyCache.Set(fingerprint, "", time.Second)
		return "", false, nil
	}
	s.sshKeyCache.Set(fingerprint, account, time.Minute)
	return account, b, nil
}

func (s *innerImpl) GetVerifiedByAccount(ctx context.Context, account string) ([]string, error) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	keys, err := sshkeymd.GetVerifiedByAccount(ctx, account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(keys, func(t sshkeymd.SshKey) (string, error) {
		return t.Content, err
	})
}

type outerImpl struct{}

func (s *outerImpl) DeleteSshKey(ctx context.Context, reqDTO DeleteSshKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	sshKey, b, err := sshkeymd.GetById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 只有拥有人才能删掉公钥
	if sshKey.Account != reqDTO.Operator.Account {
		return util.InvalidArgsError()
	}
	_, err = sshkeymd.DeleteById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (s *outerImpl) InsertSshKey(ctx context.Context, reqDTO InsertSshKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	publicKey, _, _, _, err := gossh.ParseAuthorizedKey([]byte(reqDTO.PubKeyContent))
	if err != nil {
		return util.NewBizErr(apicode.InvalidArgsCode, i18n.SshKeyFormatError)
	}
	//token, b := s.tokenCache.Get(reqDTO.Operator.Account)
	//// token不存在或已失效
	//if !b {
	//	return util.NewBizErr(apicode.SshKeyVerifyTokenExpiredCode, i18n.SshKeyVerifyTokenExpired)
	//}
	//err = signature.VerifySshSignature(reqDTO.Signature, token.(string), reqDTO.PubKeyContent)
	//if err != nil {
	//	// 校验失败
	//	return util.NewBizErr(apicode.SshKeyVerifyFailedCode, i18n.SshKeyVerifyFailed)
	//}
	fingerprint := gossh.FingerprintSHA256(publicKey)
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := sshkeymd.GetAccountByFingerprint(ctx, fingerprint)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.NewBizErr(apicode.InvalidArgsCode, i18n.SshKeyAlreadyExists)
	}
	_, err = sshkeymd.InsertSshKey(ctx, sshkeymd.InsertSshKeyReqDTO{
		Account:     reqDTO.Operator.Account,
		Name:        reqDTO.Name,
		Fingerprint: fingerprint,
		Content:     strings.TrimSpace(reqDTO.PubKeyContent),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) ListSshKey(ctx context.Context, reqDTO ListSshKeyReqDTO) ([]sshkeymd.SshKey, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 展示登录人的sshkey列表
	keyList, err := sshkeymd.GetByAccount(ctx, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return keyList, err
}

func (s *outerImpl) GetToken(ctx context.Context, reqDTO GetTokenReqDTO) (string, []string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return "", nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	token := signature.GetToken(signature.User{
		Account: reqDTO.Operator.Account,
		Email:   reqDTO.Operator.Email,
	})
	// 设置十分钟有效期
	//s.tokenCache.Set(reqDTO.Operator.Account, token, 10*time.Minute)
	// 指引
	guide := fmt.Sprintf("echo -n '%s' | ssh-keygen -Y sign -n git -f /path_to_your_privkey", token)
	return token, []string{
		i18n.GetByKey(i18n.SshKeyVerifyGuide),
		guide,
	}, nil
}
