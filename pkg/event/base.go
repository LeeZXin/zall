package event

type Event interface {
	EventType() string
}

type BaseEvent struct {
	Operator  string `json:"operator"`
	EventTime int64  `json:"eventTime"`
}
