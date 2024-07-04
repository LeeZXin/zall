package statusmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func ListService(ctx context.Context, reqDTO ListServiceReqDTO) ([]Service, error) {
	ret := make([]Service, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app = ?", reqDTO.AppId).
		And("env = ?", reqDTO.Env).
		Asc("id").
		Find(&ret)
	return ret, err
}

func GetServiceByServiceId(ctx context.Context, serviceId string) (Service, bool, error) {
	var ret Service
	b, err := xormutil.MustGetXormSession(ctx).Where("service_id = ?", serviceId).Get(&ret)
	return ret, b, err
}
