package sshkeymd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func IsSshKeyNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 128
}

func GetAccountByFingerprint(ctx context.Context, fingerprint string) (string, bool, error) {
	var ret SshKey
	b, err := xormutil.MustGetXormSession(ctx).
		Where("fingerprint = ?", fingerprint).
		Cols("account").
		Get(&ret)
	return ret.Account, b, err
}

func GetAccountAndIdByFingerprint(ctx context.Context, fingerprint string) (string, int64, bool, error) {
	var ret SshKey
	b, err := xormutil.MustGetXormSession(ctx).
		Where("fingerprint = ?", fingerprint).
		Cols("account", "id").
		Get(&ret)
	return ret.Account, ret.Id, b, err
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
		Account:      reqDTO.Account,
		Name:         reqDTO.Name,
		Fingerprint:  reqDTO.Fingerprint,
		Content:      reqDTO.Content,
		LastOperated: reqDTO.LastOperated,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func ListAllKeyByAccount(ctx context.Context, account string, cols []string) ([]SshKey, error) {
	ret := make([]SshKey, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func UpdateLastOperated(ctx context.Context, id int64, lastOperated time.Time) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("last_operated").
		Update(&SshKey{
			LastOperated: lastOperated,
		})
	return rows == 1, err
}
