package productmd

type InsertProductReqDTO struct {
	AppId   string
	Name    string
	Creator string
	Env     string
}

type GetProductReqDTO struct {
	AppId string
	Name  string
	Env   string
}

type ListProductReqDTO struct {
	AppId    string
	Env      string
	PageNum  int
	PageSize int
}
