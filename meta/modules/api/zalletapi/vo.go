package zalletapi

type CreateZalletNodeReqVO struct {
	NodeId     string `json:"nodeId"`
	Name       string `json:"name"`
	AgentHost  string `json:"agentHost"`
	AgentToken string `json:"agentToken"`
}

type UpdateZalletNodeReqVO struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	AgentHost  string `json:"agentHost"`
	AgentToken string `json:"agentToken"`
}

type ListZalletNodeReqVO struct {
	PageNum int    `json:"pageNum"`
	Name    string `json:"name"`
}

type ZalletNodeVO struct {
	Id         int64  `json:"id"`
	NodeId     string `json:"nodeId"`
	Name       string `json:"name"`
	AgentHost  string `json:"agentHost"`
	AgentToken string `json:"agentToken"`
}

type SimpleZalletNodeVO struct {
	Id     int64  `json:"id"`
	NodeId string `json:"nodeId"`
	Name   string `json:"name"`
}
