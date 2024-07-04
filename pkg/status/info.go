package status

type Api struct {
	Headers map[string]string `json:"headers,omitempty"`
	Method  string            `json:"method"`
	Url     string            `json:"url"`
}

type Service struct {
	Id      string `json:"id"`
	App     string `json:"app"`
	Status  string `json:"status"`
	Host    string `json:"host"`
	Env     string `json:"env"`
	Created string `json:"created"`
}

type Action struct {
	Label string `json:"label"`
	Api   Api    `json:"api"`
}
