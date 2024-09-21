package timer

import (
	"github.com/LeeZXin/zall/pkg/http"
)

const (
	DefaultTrigger = "system"
)

type TaskType string

const (
	HttpTask TaskType = "http"
)

type Task struct {
	TaskType TaskType   `json:"taskType"`
	HttpTask *http.Task `json:"httpTask,omitempty"`
}

func (t *Task) IsValid() bool {
	switch t.TaskType {
	case HttpTask:
		return t.HttpTask != nil && t.HttpTask.IsValid()
	default:
		return false
	}
}

type TriggerType int

func (t TriggerType) IsValid() bool {
	switch t {
	case AutoTriggerType, ManualTriggerType:
		return true
	default:
		return false
	}
}

func (t TriggerType) String() string {
	switch t {
	case AutoTriggerType:
		return "auto"
	case ManualTriggerType:
		return "manuallyTrigger"
	default:
		return ""
	}
}

const (
	AutoTriggerType TriggerType = iota + 1
	ManualTriggerType
)
