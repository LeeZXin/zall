package actionagentapi

type RunScriptReqVO struct {
	Workdir string   `json:"workdir"`
	Envs    []string `json:"envs"`
	Script  string   `json:"script"`
}

type RunScriptRespVO struct {
	Stderr string `json:"stderr"`
	Stdout string `json:"stdout"`
}
