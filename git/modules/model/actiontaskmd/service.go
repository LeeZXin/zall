package actiontaskmd

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
	"time"
)

var (
	instanceIdRegexp = regexp.MustCompile(`^\w{1,32}$`)
)

func IsInstanceIdValid(instanceId string) bool {
	return instanceIdRegexp.MatchString(instanceId)
}

func InsertTask(ctx context.Context, reqDTO InsertTaskReqDTO) (Task, error) {
	hookContent, _ := json.Marshal(reqDTO.Hook)
	ret := Task{
		TaskName:    reqDTO.TaskName,
		RepoId:      reqDTO.RepoId,
		InstanceId:  reqDTO.InstanceId,
		TaskType:    reqDTO.TaskType,
		TaskStatus:  reqDTO.TaskStatus,
		HookContent: string(hookContent),
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateTaskStatus(ctx context.Context, taskId int64, oldStatus, newStatus TaskStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", taskId).
		And("task_status = ?", oldStatus).
		Cols("task_status").
		Limit(1).
		Update(&Task{
			TaskStatus: newStatus,
		})
	return rows == 1, err
}

func BatchInsertSteps(ctx context.Context, reqDTO []InsertStepReqDTO) ([]Step, error) {
	ret, _ := listutil.Map(reqDTO, func(t InsertStepReqDTO) (*Step, error) {
		return &Step{
			TaskId:     t.TaskId,
			JobName:    t.JobName,
			StepName:   t.StepName,
			StepIndex:  t.StepIndex,
			StepStatus: t.StepStatus,
		}, nil
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(ret)
	if err != nil {
		return nil, err
	}
	return listutil.Map(ret, func(t *Step) (Step, error) {
		return *t, nil
	})
}

func UpdateInstanceHeartbeat(ctx context.Context, instanceId string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		Cols("heartbeat").
		Limit(1).
		Update(&Instance{
			Heartbeat: time.Now().UnixMilli(),
		})
	return rows == 1, err
}

func InsertInstance(ctx context.Context, instanceId, name, host string) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Instance{
		InstanceId:   instanceId,
		InstanceHost: host,
		Name:         name,
		Heartbeat:    time.Now().UnixMilli(),
	})
	return err
}

func DeleteInstance(ctx context.Context, instanceId string) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("instance_id = ?", instanceId).Delete(new(Instance))
	return err
}

func UpdateInstanceTaskFailStatus(ctx context.Context, instanceId string) (int64, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		And("task_status = ?", TaskRunningStatus).
		Cols("task_status").
		Update(&Task{
			TaskStatus: TaskFailStatus,
		})
	return rows, err
}

func IterateTask(ctx context.Context, instanceId string, taskStatus TaskStatus, fn func(*Task) error) error {
	return xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		And("task_status = ?", taskStatus).
		Iterate(new(Task), func(_ int, bean interface{}) error {
			return fn(bean.(*Task))
		})
}

func UpdateStepStatus(ctx context.Context, taskId int64, jobName string, stepIndex int, oldStatus, newStatus StepStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("task_id = ?", taskId).
		And("job_name = ?", jobName).
		And("step_status = ?", oldStatus).
		And("step_index = ?", stepIndex).
		Cols("step_status").
		Limit(1).
		Update(&Step{
			StepStatus: newStatus,
		})
	return rows == 1, err
}

func UpdateStepLogContent(ctx context.Context, taskId int64, jobName string, stepIndex int, content string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("task_id = ?", taskId).
		And("job_name = ?", jobName).
		And("step_index = ?", stepIndex).
		Cols("log_content").
		Limit(1).
		Update(&Step{
			LogContent: content,
		})
	return rows == 1, err
}

func SelectLeastJobCountInstancesForUpdate(ctx context.Context, leastHeartbeatTime int64, instanceId string) (Instance, bool, error) {
	ret := Instance{}
	session := xormutil.MustGetXormSession(ctx).Where("heartbeat >= ?", leastHeartbeatTime)
	if instanceId != "" {
		session.And("instance_id = ?", instanceId)
	}
	b, err := session.OrderBy("job_count asc, id asc").ForUpdate().Get(&ret)
	return ret, b, err
}

func IncrJobCount(ctx context.Context, instanceId string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		Cols("job_count").
		Incr("job_count").
		Limit(1).
		Update(new(Instance))
	return rows == 1, err
}
