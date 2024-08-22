package event

type Event interface {
	EventType() string
}

type BaseEvent struct {
	Operator     string `json:"operator"`
	OperatorName string `json:"operatorName"`
	EventTime    string `json:"eventTime"`
	ActionName   string `json:"actionName"`
	ActionNameEn string `json:"actionNameEn"`
}
