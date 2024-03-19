package gitactionmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertTask(ctx context.Context, reqDTO InsertTaskReqDTO) (Task, error) {
	ret := Task{
		TaskName:      reqDTO.TaskName,
		ActionId:      reqDTO.ActionId,
		TriggerType:   reqDTO.TriggerType,
		TaskStatus:    reqDTO.TaskStatus,
		ActionContent: reqDTO.ActionContent,
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

func UpdateTaskStatusWithOldStatus(ctx context.Context, instanceId string, newStatus, oldStatus TaskStatus) (int64, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		And("task_status = ?", oldStatus).
		Cols("task_status").
		Update(&Task{
			TaskStatus: newStatus,
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

func IsNodeNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertNode(ctx context.Context, reqDTO InsertNodeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Node{
			Name:     reqDTO.Name,
			HttpHost: reqDTO.HttpHost,
		})
	return err
}

func GetNodeById(ctx context.Context, id int64) (Node, bool, error) {
	var ret Node
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func BatchGetNode(ctx context.Context, idList []int64) ([]Node, error) {
	ret := make([]Node, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", idList).
		Find(&ret)
	return ret, err
}

func DeleteNode(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(Node))
	return rows == 1, err
}

func UpdateNode(ctx context.Context, reqDTO UpdateNodeReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "http_host").
		Limit(1).
		Update(Node{
			HttpHost: reqDTO.HttpHost,
			Name:     reqDTO.Name,
		})
	return rows == 1, err
}

func GetAllNodes(ctx context.Context) ([]Node, error) {
	ret := make([]Node, 0)
	session := xormutil.MustGetXormSession(ctx)
	err := session.OrderBy("id asc").Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, err
}

func InsertAction(ctx context.Context, reqDTO InsertActionReqDTO) error {
	ret := Action{
		RepoId:     reqDTO.RepoId,
		Content:    reqDTO.Content,
		NodeId:     reqDTO.NodeId,
		WildBranch: reqDTO.PushBranch,
		Name:       reqDTO.Name,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func UpdateAction(ctx context.Context, reqDTO UpdateActionReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("content", "node_id", "push_branch", "name").
		Update(&Action{
			Content:    reqDTO.Content,
			NodeId:     reqDTO.NodeId,
			WildBranch: reqDTO.PushBranch,
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
