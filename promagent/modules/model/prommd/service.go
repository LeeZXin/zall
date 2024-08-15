package prommd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsTargetValid(target string) bool {
	return len(target) > 0
}

func IsEndpointValid(endpoint string) bool {
	return len(endpoint) > 0 && len(endpoint) <= 32
}

func GetAllScrape(ctx context.Context, reqDTO GetAllScrapeReqDTO) ([]Scrape, error) {
	ret := make([]Scrape, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("endpoint = ?", reqDTO.Endpoint).
		And("env = ?", reqDTO.Env)
	if len(reqDTO.Cols) > 0 {
		session.Cols(reqDTO.Cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func ListScrape(ctx context.Context, reqDTO ListScrapeReqDTO) ([]Scrape, int64, error) {
	ret := make([]Scrape, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("env = ?", reqDTO.Env)
	if reqDTO.AppId != "" {
		session.And("app_id = ?", reqDTO.AppId)
	}
	if reqDTO.Endpoint != "" {
		session.And("endpoint like ?", reqDTO.Endpoint+"%")
	}
	total, err := session.
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		FindAndCount(&ret)
	return ret, total, err
}

func DeleteScrapeById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(Scrape))
	return rows == 1, err
}

func DeleteScrapeByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(Scrape))
	return err
}

func GetScrapeById(ctx context.Context, id int64) (Scrape, bool, error) {
	var ret Scrape
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func UpdateScrapeById(ctx context.Context, reqDTO UpdateScrapeByIdReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("endpoint", "target", "target_type").
		Update(&Scrape{
			Endpoint:   reqDTO.Endpoint,
			Target:     reqDTO.Target,
			TargetType: reqDTO.TargetType,
		})
	return rows == 1, err
}

func InsertScrape(ctx context.Context, reqDTO InsertScrapeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Scrape{
			Endpoint:   reqDTO.Endpoint,
			AppId:      reqDTO.AppId,
			Target:     reqDTO.Target,
			TargetType: reqDTO.TargetType,
			Env:        reqDTO.Env,
		})
	return err
}
