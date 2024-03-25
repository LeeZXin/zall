package actionapi

type InsertActionReqVO struct {
	Name          string `json:"name"`
	TeamId        int64  `json:"teamId"`
	AgentUrl      string `json:"agentUrl"`
	AgentToken    string `json:"agentToken"`
	ActionContent string `json:"actionContent"`
}

type DeleteActionReqVO struct {
	Id int64 `json:"id"`
}

type ListActionReqVO struct {
	TeamId int64 `json:"teamId"`
}

type UpdateActionReqVO struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	AgentUrl      string `json:"agentUrl"`
	AgentToken    string `json:"agentToken"`
	ActionContent string `json:"actionContent"`
}

type ActionVO struct {
	Id            int64  `json:"id"`
	Aid           string `json:"aid"`
	AgentUrl      string `json:"agentUrl"`
	AgentToken    string `json:"agentToken"`
	ActionContent string `json:"actionContent"`
	Created       string `json:"created"`
}

type TriggerActionReqVO struct {
	Id int64 `json:"id"`
}

type ListTaskReqVO struct {
	ActionId int64 `json:"actionId"`
	Cursor   int64 `json:"cursor"`
	Limit    int   `json:"limit"`
}

type TaskVO struct {
	TaskStatus    string `json:"taskStatus"`
	TriggerType   string `json:"triggerType"`
	ActionContent string `json:"actionContent"`
	Operator      string `json:"operator"`
	Created       string `json:"created"`
}

type StepVO struct {
	JobName    string `json:"jobName"`
	StepName   string `json:"stepName"`
	StepIndex  int    `json:"stepIndex"`
	LogContent string `json:"logContent"`
	StepStatus string `json:"stepStatus"`
	Created    string `json:"created"`
	Updated    string `json:"updated"`
}

type GetTaskStepsReqVO struct {
	TaskId int64 `json:"taskId"`
}
