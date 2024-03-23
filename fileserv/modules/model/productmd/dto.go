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
