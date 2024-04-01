package promsrv

import "context"

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	InsertScrape(context.Context, InsertScrapeReqDTO) error
	UpdateScrape(context.Context, UpdateScrapeReqDTO) error
	ListScrape(context.Context, ListScrapeReqDTO) ([]ScrapeDTO, error)
	DeleteScrape(context.Context, DeleteScrapeReqDTO) error
}
