package event

type AlertEvent struct {
	BaseTeam
	BaseApp
	EventTime string         `json:"eventTime"`
	StartTime string         `json:"startTime"`
	EndTime   string         `json:"endTime"`
	Result    map[string]any `json:"result"`
}

func (*AlertEvent) EventType() string {
	return "alert-event"
}
