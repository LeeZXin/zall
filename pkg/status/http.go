package status

type ListServiceReq struct {
	App string `json:"app"`
	Env string `json:"env"`
}

type ListServiceResp struct {
	Services []Service `json:"services"`
	Next     string    `json:"next"`
}
