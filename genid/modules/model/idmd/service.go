package idmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func GetByBizName(ctx context.Context, bizName string) (Generator, bool, error) {
	ret := Generator{}
	b, err := xormutil.MustGetXormSession(ctx).Where("biz_name = ?", bizName).Get(&ret)
	return ret, b, err
}

func UpdateCurrentId(ctx context.Context, reqDTO UpdateCurrentIdReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("biz_name = ?", reqDTO.BizName).
		And("version = ?", reqDTO.Version).
		Cols("current_id", "version").
		Limit(1).
		Update(&Generator{
			Version:   reqDTO.Version + 1,
			CurrentId: reqDTO.CurrentId,
		})
	return rows == 1, err
}

func InsertGenerator(ctx context.Context, reqDTO InsertGeneratorReqDTO) error {
	ret := Generator{
		BizName:   reqDTO.BizName,
		CurrentId: reqDTO.CurrentId,
		Version:   0,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}
