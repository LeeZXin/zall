package event

type TimerEventCfg struct {
	Create          bool `json:"create"`
	Update          bool `json:"update"`
	Delete          bool `json:"delete"`
	Enable          bool `json:"enable"`
	Disable         bool `json:"disable"`
	ManuallyTrigger bool `json:"manuallyTrigger"`
}

type TimerTaskEventCfg struct {
	Fail bool `json:"fail"`
}

type TeamEventCfg struct {
	Create bool `json:"create"`
	Delete bool `json:"delete"`
	Update bool `json:"update"`
}

type TeamRoleEventCfg struct {
	Create bool `json:"create"`
	Delete bool `json:"delete"`
	Update bool `json:"update"`
}

type TeamUserEventCfg struct {
	Create     bool `json:"create"`
	Delete     bool `json:"delete"`
	ChangeRole bool `json:"changeRole"`
}
