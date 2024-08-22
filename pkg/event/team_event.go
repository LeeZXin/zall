package event

type BaseTeam struct {
	TeamId   int64  `json:"teamId"`
	TeamName string `json:"teamName"`
}

type BaseTimer struct {
	TimerId   int64  `json:"timerId"`
	TimerName string `json:"timerName"`
	Env       string `json:"env"`
}

type TeamEvent struct {
	BaseTeam
	BaseEvent
	Action TeamEventAction `json:"action"`
}

func (*TeamEvent) EventType() string {
	return "team-event"
}

type TimerEvent struct {
	BaseTeam
	BaseEvent
	BaseTimer
	Action TimerEventAction `json:"action"`
}

func (*TimerEvent) EventType() string {
	return "timer-event"
}

type TimerTaskEvent struct {
	BaseTeam
	BaseEvent
	BaseTimer
	Action      TimerTaskEventAction `json:"action"`
	TriggerType string               `json:"triggerType"`
	TaskStatus  string               `json:"taskStatus"`
}

func (*TimerTaskEvent) EventType() string {
	return "timer-task-event"
}

type BaseRole struct {
	RoleId   int64  `json:"roleId"`
	RoleName string `json:"roleName"`
}

type TeamRoleEvent struct {
	BaseTeam
	BaseEvent
	BaseRole
	Action TeamRoleEventAction `json:"action"`
}

func (*TeamRoleEvent) EventType() string {
	return "team-role-event"
}

type TeamUserEvent struct {
	BaseTeam
	BaseEvent
	BaseRole
	User   string              `json:"user"`
	Action TeamUserEventAction `json:"action"`
}

func (*TeamUserEvent) EventType() string {
	return "team-member-event"
}
