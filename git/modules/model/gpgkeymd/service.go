package gpgkeymd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func IsKeyNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 128
}

func ListValidByAccount(ctx context.Context, account string) ([]GpgKey, error) {
	ret := make([]GpgKey, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		And("expired > ?", time.Now().Format(time.DateTime)).
		Find(&ret)
	return ret, err
}

func BatchInsertGpgKeys(ctx context.Context, reqDTOs []InsertGpgKeyReqDTO) error {
	keys, _ := listutil.Map(reqDTOs, func(reqDTO InsertGpgKeyReqDTO) (GpgKey, error) {
		subKeys, _ := listutil.Map(reqDTO.SubKeys, func(t InsertGpgSubKeyReqDTO) (GpgSubKey, error) {
			return GpgSubKey{
				KeyId:   t.KeyId,
				Content: t.Content,
			}, nil
		})
		return GpgKey{
			Account: reqDTO.Account,
			Name:    reqDTO.Name,
			KeyId:   reqDTO.KeyId,
			Content: reqDTO.Content,
			Expired: reqDTO.Expired,
			Email:   reqDTO.Email,
			SubKeys: subKeys,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(keys)
	return err
}

func CountByKeyIds(ctx context.Context, keyIds []string) (int64, error) {
	return xormutil.MustGetXormSession(ctx).In("key_id", keyIds).Count(new(GpgKey))
}

func GetById(ctx context.Context, id int64) (GpgKey, bool, error) {
	var ret GpgKey
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func GetByKeyId(ctx context.Context, keyId string) (GpgKey, bool, error) {
	var ret GpgKey
	b, err := xormutil.MustGetXormSession(ctx).Where("key_id = ?", keyId).Get(&ret)
	return ret, b, err
}

func DeleteById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(GpgKey))
	return rows == 1, err
}

func ListAllByAccount(ctx context.Context, account string, cols []string) ([]GpgKey, error) {
	ret := make([]GpgKey, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account)
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	err := session.Find(&ret)
	return ret, err
}
