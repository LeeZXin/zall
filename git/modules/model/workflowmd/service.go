package workflowmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
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
		YamlContent: reqDTO.YamlContent,
		AgentHost:   reqDTO.AgentHost,
		AgentToken:  reqDTO.AgentToken,
		Branch:      reqDTO.Branch,
		Operator:    reqDTO.Operator,
		PrId:        reqDTO.PrId,
		BizId:       reqDTO.BizId,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateTaskStatusAndDuration(ctx context.Context, taskId int64, oldStatus, newStatus TaskStatus, duration time.Duration) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", taskId).
		And("task_status = ?", oldStatus).
		Cols("task_status", "duration").
		Limit(1).
		Update(&Task{
			TaskStatus: newStatus,
			Duration:   duration.Milliseconds(),
		})
	return rows == 1, err
}

func InsertWorkflow(ctx context.Context, reqDTO InsertWorkflowReqDTO) error {
	ret := Workflow{
		Name:        reqDTO.Name,
		Description: reqDTO.Desc,
		RepoId:      reqDTO.RepoId,
		YamlContent: reqDTO.YamlContent,
		Source:      &reqDTO.Source,
		AgentHost:   reqDTO.AgentHost,
		AgentToken:  reqDTO.AgentToken,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func UpdateWorkflow(ctx context.Context, reqDTO UpdateWorkflowReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("yaml_content", "agent_host", "agent_token", "name", "description", "source").
		Update(&Workflow{
			YamlContent: reqDTO.Content,
			AgentHost:   reqDTO.AgentHost,
			AgentToken:  reqDTO.AgentToken,
			Name:        reqDTO.Name,
			Description: reqDTO.Desc,
			Source:      &reqDTO.Source,
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

func BatchGetWorkflowNameById(ctx context.Context, idList []int64) (map[int64]string, error) {
	wfList := make([]Workflow, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", idList).
		Cols("name", "id").
		Find(&wfList)
	if err != nil {
		return nil, err
	}
	ret := make(map[int64]string, len(wfList))
	for _, wf := range wfList {
		ret[wf.Id] = wf.Name
	}
	return ret, nil
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

func ListTaskByPrId(ctx context.Context, prId int64) ([]Task, error) {
	ret := make([]Task, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("pr_id = ?", prId).
		Find(&ret)
	return ret, err
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

func GetTaskByBizId(ctx context.Context, bizId string) (Task, bool, error) {
	var ret Task
	b, err := xormutil.MustGetXormSession(ctx).Where("biz_id = ?", bizId).Get(&ret)
	return ret, b, err
}

func DeleteTaskById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(Task))
	return rows == 1, err
}
