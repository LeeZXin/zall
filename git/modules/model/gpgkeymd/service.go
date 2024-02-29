package gpgkeymd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"strings"
	"time"
)

func IsKeyNameValid(name string) bool {
	return len(name) > 0 && len(name) < 32
}

func GetValidByAccount(ctx context.Context, account string) ([]GpgKey, error) {
	ret := make([]GpgKey, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		And("expiry >= ?", time.Now().UnixMilli()).
		Find(&ret)
	return ret, err
}

func InsertGpgKeys(ctx context.Context, reqDTO InsertGpgKeyReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(GpgKey{
		Account:   reqDTO.Account,
		Name:      reqDTO.Name,
		PubKeyId:  reqDTO.PubKeyId,
		Content:   reqDTO.Content,
		Expiry:    reqDTO.Expiry,
		EmailList: strings.Join(reqDTO.EmailList, ";"),
	})
	return err
}

func ExistsByPubKeyId(ctx context.Context, pubKeyId string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).Where("pub_key_id = ?", pubKeyId).Exist(new(GpgKey))
}

func GetById(ctx context.Context, id int64) (GpgKey, bool, error) {
	var ret GpgKey
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func DeleteById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(GpgKey))
	return rows == 1, err
}

func GetAllSimpleByAccount(ctx context.Context, account string) ([]SimpleGpgKeyDTO, error) {
	ret := make([]SimpleGpgKey, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		Find(&ret)
	if err != nil {
		return nil, err
	}
	return listutil.Map(ret, func(t SimpleGpgKey) (SimpleGpgKeyDTO, error) {
		return SimpleGpgKeyDTO{
			Id:        t.Id,
			Name:      t.Name,
			PubKeyId:  t.PubKeyId,
			Expiry:    t.Expiry,
			EmailList: strings.Split(t.EmailList, ";"),
		}, nil
	})
}
