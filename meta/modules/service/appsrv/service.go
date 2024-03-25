package appsrv

import "context"

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	InsertApp(context.Context, InsertAppReqDTO) error
	DeleteApp(context.Context, DeleteAppReqDTO) error
	UpdateApp(context.Context, UpdateAppReqDTO) error
	ListApp(context.Context, ListAppReqDTO) ([]AppDTO, error)
	TransferTeam(context.Context, TransferTeamReqDTO) error
}
