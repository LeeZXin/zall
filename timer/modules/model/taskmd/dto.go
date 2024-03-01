package taskmd

type InsertTaskReqDTO struct {
	Name       string
	CronExp    string
	Content    string
	NextTime   int64
	TaskStatus TaskStatus
	TeamId     int64
}

type UpdateTaskReqDTO struct {
	Id         int64
	Name       string
	CronExp    string
	Content    string
	NextTime   int64
	TaskStatus TaskStatus
	Version    int64
}

type ListTaskReqDTO struct {
	TeamId int64
	Name   string
	Cursor int64
	Limit  int
}

type InsertTaskLogReqDTO struct {
	TaskId      int64
	TaskContent string
	LogContent  string
	TriggerType TriggerType
	TriggerBy   string
	TaskStatus  TaskStatus
}

type ListTaskLogReqDTO struct {
	Id     int64
	Cursor int64
	Limit  int
}

type UpdateLogContentAndStatusReqDTO struct {
	Id      int64
	Content string
	Status  TaskStatus
}
