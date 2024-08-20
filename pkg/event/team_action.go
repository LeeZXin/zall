package event

type TimerEventAction string

const (
	TimerCreateAction TimerEventAction = "create"
	TimerUpdateAction TimerEventAction = "update"
	TimerDeleteAction TimerEventAction = "delete"
)

type TeamEventAction string

const (
	TeamCreateAction TeamEventAction = "create"
	TeamUpdateAction TeamEventAction = "update"
	TeamDeleteAction TeamEventAction = "delete"
)

type TimerTaskEventAction string

const (
	TimerTaskManualTriggerAction TimerTaskEventAction = "manuallyTrigger"
	TimerTaskFailAction          TimerTaskEventAction = "fail"
)

type TeamRoleEventAction string

const (
	TeamRoleCreateAction TeamRoleEventAction = "create"
	TeamRoleUpdateAction TeamRoleEventAction = "update"
	TeamRoleDeleteAction TeamRoleEventAction = "delete"
)

type TeamMemberEventAction string

const (
	TeamMemberCreateAction TeamMemberEventAction = "create"
	TeamMemberUpdateAction TeamMemberEventAction = "update"
	TeamMemberDeleteAction TeamMemberEventAction = "delete"
)
