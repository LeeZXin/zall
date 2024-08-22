package event

type TimerEventAction string

func (t TimerEventAction) GetI18nValue() string {
	return "timerEvent." + string(t)
}

const (
	TimerCreateAction        TimerEventAction = "create"
	TimerUpdateAction        TimerEventAction = "update"
	TimerDeleteAction        TimerEventAction = "delete"
	TimerEnableAction        TimerEventAction = "enable"
	TimerDisableAction       TimerEventAction = "disable"
	TimerManualTriggerAction TimerEventAction = "manuallyTrigger"
)

type TeamEventAction string

func (t TeamEventAction) GetI18nValue() string {
	return "teamEvent." + string(t)
}

const (
	TeamCreateAction TeamEventAction = "create"
	TeamUpdateAction TeamEventAction = "update"
	TeamDeleteAction TeamEventAction = "delete"
)

type TimerTaskEventAction string

func (t TimerTaskEventAction) GetI18nValue() string {
	return "timerTaskEvent." + string(t)
}

const (
	TimerTaskFailAction TimerTaskEventAction = "fail"
)

type TeamRoleEventAction string

const (
	TeamRoleCreateAction TeamRoleEventAction = "create"
	TeamRoleUpdateAction TeamRoleEventAction = "update"
	TeamRoleDeleteAction TeamRoleEventAction = "delete"
)

type TeamUserEventAction string

const (
	TeamUserCreateAction     TeamUserEventAction = "create"
	TeamUserChangeRoleAction TeamUserEventAction = "changeRole"
	TeamUserDeleteAction     TeamUserEventAction = "delete"
)
