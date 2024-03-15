package taskmd

type InsertTaskReqDTO struct {
	Name       string
	CronExp    string
	Content    string
	NextTime   int64
	TaskStatus TaskStatus
	TeamId     int64
	Env        string
}

type UpdateTaskReqDTO struct {
	Id         int64
	Name       string
	CronExp    string
	Content    string
	NextTime   int64
	TaskStatus TaskStatus
	Version    int64
	Env        string
}

type ListTaskReqDTO struct {
	TeamId int64
	Name   string
	Cursor int64
	Limit  int
	Env    string
}

type InsertTaskLogReqDTO struct {
	TaskId      int64
	TaskContent string
	LogContent  string
	TriggerType TriggerType
	TriggerBy   string
	TaskStatus  TaskStatus
	Env         string
}

type ListTaskLogReqDTO struct {
	Id     int64
	Cursor int64
	Limit  int
	Env    string
}

type UpdateLogContentAndStatusReqDTO struct {
	Id      int64
	Content string
	Status  TaskStatus
}

type UpdateTaskStatusReqDTO struct {
	TaskId    int64
	NewStatus TaskStatus
	Version   int64
	Env       string
}

type UpdateTaskNextTimeAndStatusReqDTO struct {
	TaskId   int64
	Status   TaskStatus
	NextTime int64
	Version  int64
	Env      string
}
