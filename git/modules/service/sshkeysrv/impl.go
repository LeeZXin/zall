package sshkeysrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/sshkeymd"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	gossh "golang.org/x/crypto/ssh"
	"strings"
	"time"
)

type innerImpl struct {
}

func (s *innerImpl) ListAllPubKeyByAccount(ctx context.Context, account string) []string {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	keys, err := sshkeymd.ListAllKeyByAccount(ctx, account, []string{"content"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return []string{}
	}
	ret, _ := listutil.Map(keys, func(t sshkeymd.SshKey) (string, error) {
		return t.Content, nil
	})
	return ret
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
		return util.UnauthorizedError()
	}
	_, err = sshkeymd.DeleteById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// CreateSshKey 创建ssh密钥
func (s *outerImpl) CreateSshKey(ctx context.Context, reqDTO CreateSshKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	publicKey, _, _, _, err := gossh.ParseAuthorizedKey([]byte(reqDTO.Content))
	if err != nil {
		return util.NewBizErr(apicode.InvalidArgsCode, i18n.SshKeyFormatError)
	}
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
		Account:      reqDTO.Operator.Account,
		Name:         reqDTO.Name,
		Fingerprint:  fingerprint,
		Content:      strings.TrimSpace(reqDTO.Content),
		LastOperated: time.Now(),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) ListSshKey(ctx context.Context, reqDTO ListSshKeyReqDTO) ([]SshKeyDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	keys, err := sshkeymd.ListAllKeyByAccount(
		ctx,
		reqDTO.Operator.Account,
		[]string{"id", "name", "fingerprint", "created", "last_operated"},
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(keys, func(t sshkeymd.SshKey) (SshKeyDTO, error) {
		return SshKeyDTO{
			Id:           t.Id,
			Name:         t.Name,
			Fingerprint:  t.Fingerprint,
			Created:      t.Created,
			LastOperated: t.LastOperated,
		}, nil
	})
}
