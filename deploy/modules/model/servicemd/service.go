package servicemd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsServiceNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertService(ctx context.Context, reqDTO InsertServiceReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Service{
			AppId:       reqDTO.AppId,
			Name:        reqDTO.Name,
			ServiceType: reqDTO.ServiceType,
			Config:      reqDTO.Config,
			Env:         reqDTO.Env,
			Probed:      reqDTO.Probed,
			IsEnabled:   reqDTO.IsEnabled,
		})
	return err
}

func ListService(ctx context.Context, reqDTO ListServiceReqDTO) ([]Service, int64, error) {
	ret := make([]Service, 0)
	total, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("env = ?", reqDTO.Env).
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		Desc("probed", "id").
		FindAndCount(&ret)
	return ret, total, err
}

func DeleteService(ctx context.Context, serviceId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", serviceId).
		Delete(new(Service))
	return rows == 1, err
}

func UpdateService(ctx context.Context, reqDTO UpdateServiceReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.ServiceId).
		Cols("config", "name", "service_type").
		Update(&Service{
			ServiceType: reqDTO.ServiceType,
			Name:        reqDTO.Name,
			Config:      reqDTO.Config,
		})
	return rows == 1, err
}

func IterateService(ctx context.Context, env string, fn func(*Service) error) error {
	return xormutil.MustGetXormSession(ctx).
		And("env = ?", env).
		And("is_enabled = 1").
		Iterate(new(Service), func(idx int, obj interface{}) error {
			if err := fn(obj.(*Service)); err != nil {
				return err
			}
			return nil
		})
}

func UpdateProbed(ctx context.Context, serviceId int64, probed int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", serviceId).
		Cols("probed").
		Update(&Service{
			Probed: probed,
		})
	return rows == 1, err
}

func GetServiceById(ctx context.Context, serviceId int64) (Service, bool, error) {
	var ret Service
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", serviceId).Get(&ret)
	return ret, b, err
}

func UpdateServiceIsEnabled(ctx context.Context, serviceId int64, isEnabled bool) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", serviceId).
		Cols("is_enabled").
		Update(&Service{
			IsEnabled: isEnabled,
		})
	return rows == 1, err
}
