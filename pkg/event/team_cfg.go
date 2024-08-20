package event

type TimerEventCfg struct {
	Create bool `json:"create"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type TimerTaskEventCfg struct {
	ManuallyTrigger bool `json:"manuallyTrigger"`
	Fail            bool `json:"fail"`
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

type TeamMemberEventCfg struct {
	Create bool `json:"create"`
	Delete bool `json:"delete"`
	Update bool `json:"update"`
}
