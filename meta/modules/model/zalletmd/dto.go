package zalletmd

type InsertZalletNodeReqDTO struct {
	Name       string
	AgentHost  string
	AgentToken string
}

type ListZalletNodeReqDTO struct {
	PageNum  int
	PageSize int
	Name     string
	Cols     []string
}

type UpdateZalletNodeReqDTO struct {
	Id         int64
	Name       string
	AgentHost  string
	AgentToken string
}
