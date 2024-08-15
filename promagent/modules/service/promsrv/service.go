package promsrv

import "context"

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		Outer = new(outerImpl)
	}
}

type OuterService interface {
	// CreateScrape 创建抓取配置
	CreateScrape(context.Context, CreateScrapeReqDTO) error
	// UpdateScrape 编辑抓取配置
	UpdateScrape(context.Context, UpdateScrapeReqDTO) error
	// ListScrape 展示抓取列表
	ListScrape(context.Context, ListScrapeReqDTO) ([]ScrapeDTO, int64, error)
	// DeleteScrape 删除抓取配置
	DeleteScrape(context.Context, DeleteScrapeReqDTO) error
}
