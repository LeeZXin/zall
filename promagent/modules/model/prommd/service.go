package prommd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsTargetValid(target string) bool {
	return len(target) > 0 && len(target) <= 32
}

func GetAllScrapeByServerUrl(ctx context.Context, serverUrl, env string) ([]Scrape, error) {
	ret := make([]Scrape, 0)
	err := xormutil.MustGetXormSession(ctx).
		Table("zprom_scrape_"+env).
		Where("server_url = ?", serverUrl).
		Find(&ret)
	return ret, err
}

func GetAllScrapeByAppId(ctx context.Context, appId, env string) ([]Scrape, error) {
	ret := make([]Scrape, 0)
	err := xormutil.MustGetXormSession(ctx).
		Table("zprom_scrape_"+env).
		Where("app_id = ?", appId).
		Find(&ret)
	return ret, err
}

func DeleteById(ctx context.Context, id int64, env string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zprom_scrape_"+env).
		Where("id = ?", id).
		Delete(new(Scrape))
	return rows == 1, err
}

func GetById(ctx context.Context, id int64, env string) (Scrape, bool, error) {
	var ret Scrape
	b, err := xormutil.MustGetXormSession(ctx).
		Table("zprom_scrape_"+env).
		Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func UpdateScrapeById(ctx context.Context, reqDTO UpdateScrapeByIdReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Table("zprom_scrape_"+reqDTO.Env).
		Where("id = ?", reqDTO.Id).
		Cols("server_url", "target", "target_type").
		Update(&Scrape{
			ServerUrl:  reqDTO.ServerUrl,
			Target:     reqDTO.Target,
			TargetType: reqDTO.TargetType,
		})
	return rows == 1, err
}

func InsertScrape(ctx context.Context, reqDTO InsertScrapeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zprom_scrape_" + reqDTO.Env).
		Insert(&Scrape{
			ServerUrl:  reqDTO.ServerUrl,
			AppId:      reqDTO.AppId,
			Target:     reqDTO.Target,
			TargetType: reqDTO.TargetType,
		})
	return err
}
