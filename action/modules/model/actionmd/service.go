package actionmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsActionNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertTask(ctx context.Context, reqDTO InsertTaskReqDTO) (Task, error) {
	ret := Task{
		ActionId:      reqDTO.ActionId,
		TaskStatus:    reqDTO.TaskStatus,
		TriggerType:   reqDTO.TriggerType,
		ActionContent: reqDTO.ActionContent,
		AgentUrl:      reqDTO.AgentUrl,
		AgentToken:    reqDTO.AgentToken,
		Operator:      reqDTO.Operator,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateTaskStatusWithOldStatus(ctx context.Context, taskId int64, oldStatus, newStatus TaskStatus) (bool, error) {
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

func InsertAction(ctx context.Context, reqDTO InsertActionReqDTO) error {
	ret := Action{
		Aid:        reqDTO.Aid,
		TeamId:     reqDTO.TeamId,
		Content:    reqDTO.Content,
		AgentHost:  reqDTO.AgentHost,
		AgentToken: reqDTO.AgentToken,
		Name:       reqDTO.Name,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func UpdateAction(ctx context.Context, reqDTO UpdateActionReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("content", "agent_host", "agent_token", "name").
		Update(&Action{
			Content:    reqDTO.Content,
			AgentHost:  reqDTO.AgentHost,
			AgentToken: reqDTO.AgentToken,
			Name:       reqDTO.Name,
		})
	return rows == 1, err
}

func DeleteAction(ctx context.Context, actionId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", actionId).
		Delete(new(Action))
	return rows == 1, err
}

func ListAction(ctx context.Context, repoId int64) ([]Action, error) {
	ret := make([]Action, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		Find(&ret)
	return ret, err
}

func GetActionById(ctx context.Context, actionId int64) (Action, bool, error) {
	var ret Action
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", actionId).
		Get(&ret)
	return ret, b, err
}

func GetActionByAid(ctx context.Context, aid string) (Action, bool, error) {
	var ret Action
	b, err := xormutil.MustGetXormSession(ctx).
		Where("aid = ?", aid).
		Get(&ret)
	return ret, b, err
}

func GetTask(ctx context.Context, reqDTO GetTaskReqDTO) ([]Task, error) {
	ret := make([]Task, 0)
	session := xormutil.MustGetXormSession(ctx).Where("action_id = ?", reqDTO.ActionId)
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	err := session.Find(&ret)
	return ret, err
}

func GetTaskById(ctx context.Context, id int64) (Task, bool, error) {
	var ret Task
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func GetStepByTaskId(ctx context.Context, taskId int64) ([]Step, error) {
	ret := make([]Step, 0)
	err := xormutil.MustGetXormSession(ctx).Where("task_id = ?", taskId).Find(&ret)
	return ret, err
}
