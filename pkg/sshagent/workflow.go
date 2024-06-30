package sshagent

type BaseStatus struct {
	Status    Status `json:"status"`
	Duration  int64  `json:"duration"`
	ErrLog    string `json:"errLog"`
	BeginTime int64  `json:"beginTime"`
}

type TaskStatus struct {
	BaseStatus
	JobStatus []JobStatus `json:"jobStatus"`
}

type JobStatus struct {
	JobName string `json:"jobName"`
	BaseStatus
	Steps []StepStatus `json:"steps"`
}

type StepStatus struct {
	StepName string `json:"stepName"`
	BaseStatus
}

type TaskStatusCallbackReq struct {
	Status   Status      `json:"status"`
	Duration int64       `json:"duration"`
	Task     *TaskStatus `json:"task"`
}

type Status string

const (
	RunningStatus Status = "running"
	SuccessStatus Status = "success"
	FailStatus    Status = "fail"
	TimeoutStatus Status = "timeout"
	CancelStatus  Status = "cancel"
	QueueStatus   Status = "queue"
)

func (s Status) IsFinalType() bool {
	switch s {
	case FailStatus, TimeoutStatus, CancelStatus:
		return true
	default:
		return false
	}
}
