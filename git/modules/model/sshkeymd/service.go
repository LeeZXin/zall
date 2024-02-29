package sshkeymd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func GetAccountByFingerprint(ctx context.Context, fingerprint string) (string, bool, error) {
	var ret SshKey
	b, err := xormutil.MustGetXormSession(ctx).
		Where("fingerprint = ?", fingerprint).
		Cols("account").
		Get(&ret)
	return ret.Account, b, err
}

func GetById(ctx context.Context, id int64) (SshKey, bool, error) {
	var ret SshKey
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func DeleteById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(SshKey))
	return rows == 1, err
}

func InsertSshKey(ctx context.Context, reqDTO InsertSshKeyReqDTO) (SshKey, error) {
	ret := SshKey{
		Account:     reqDTO.Account,
		Name:        reqDTO.Name,
		Fingerprint: reqDTO.Fingerprint,
		Content:     reqDTO.Content,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetByAccount(ctx context.Context, account string) ([]SshKey, error) {
	ret := make([]SshKey, 0)
	err := xormutil.MustGetXormSession(ctx).Where("account = ?", account).Find(&ret)
	return ret, err
}

func GetVerifiedByAccount(ctx context.Context, account string) ([]SshKey, error) {
	ret := make([]SshKey, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		And("verified = 1").
		Find(&ret)
	return ret, err
}
