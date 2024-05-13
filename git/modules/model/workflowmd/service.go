package workflowmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func IsWorkflowNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsWorkflowDescValid(desc string) bool {
	return len(desc) > 0 && len(desc) <= 255
}

func InsertTask(ctx context.Context, reqDTO InsertTaskReqDTO) (Task, error) {
	ret := Task{
		WorkflowId:  reqDTO.WorkflowId,
		TaskStatus:  reqDTO.TaskStatus,
		TriggerType: reqDTO.TriggerType,
		Workflow:    &reqDTO.Workflow,
		Operator:    reqDTO.Operator,
		Branch:      reqDTO.Branch,
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
			WorkflowId: t.WorkflowId,
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

func InsertWorkflow(ctx context.Context, reqDTO InsertWorkflowReqDTO) error {
	ret := Workflow{
		RepoId:      reqDTO.RepoId,
		YamlContent: reqDTO.YamlContent,
		Agent:       &reqDTO.Agent,
		Name:        reqDTO.Name,
		Source:      &reqDTO.Source,
		Description: reqDTO.Desc,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func UpdateWorkflow(ctx context.Context, reqDTO UpdateWorkflowReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("content", "agent_host", "agent_secret", "name").
		Update(&Workflow{
			YamlContent: reqDTO.Content,
			Agent:       &reqDTO.Agent,
			Name:        reqDTO.Name,
		})
	return rows == 1, err
}

func DeleteWorkflow(ctx context.Context, workflowId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", workflowId).
		Delete(new(Workflow))
	return rows == 1, err
}

func ListWorkflow(ctx context.Context, repoId int64) ([]Workflow, error) {
	ret := make([]Workflow, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		Find(&ret)
	return ret, err
}

func GetWorkflowById(ctx context.Context, id int64) (Workflow, bool, error) {
	var ret Workflow
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ListTaskByWorkflowId(ctx context.Context, reqDTO ListTaskByWorkflowIdReqDTO) ([]Task, int64, error) {
	ret := make([]Task, 0)
	total, err := xormutil.MustGetXormSession(ctx).
		Where("workflow_id = ?", reqDTO.WorkflowId).
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		OrderBy("id desc").
		FindAndCount(&ret)
	return ret, total, err
}

func GetTaskById(ctx context.Context, id int64) (Task, bool, error) {
	var ret Task
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func GetTasksByIdList(ctx context.Context, idList []int64) ([]Task, error) {
	ret := make([]Task, 0)
	err := xormutil.MustGetXormSession(ctx).In("id", idList).Find(&ret)
	return ret, err
}

func GetStepByTaskId(ctx context.Context, taskId int64) ([]Step, error) {
	ret := make([]Step, 0)
	err := xormutil.MustGetXormSession(ctx).Where("task_id = ?", taskId).Find(&ret)
	return ret, err
}

func UpdateLastTaskIdByWorkflowId(ctx context.Context, workflowId int64, lastTaskId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", workflowId).
		Cols("last_task_id").
		Update(&Workflow{
			LastTaskId: lastTaskId,
		})
	return rows == 1, err
}

func DeleteTasksByWorkflowId(ctx context.Context, workflowId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("workflow_id = ?", workflowId).Delete(new(Task))
	return err
}

func DeleteStepsByWorkflowId(ctx context.Context, workflowId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("workflow_id = ?", workflowId).Delete(new(Step))
	return err
}
