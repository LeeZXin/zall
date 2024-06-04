package timer

import (
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/http"
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

func (t *Task) FromDB(content []byte) error {
	if t == nil {
		*t = Task{}
	}
	return json.Unmarshal(content, t)
}

func (t *Task) ToDB() ([]byte, error) {
	return json.Marshal(t)
}
