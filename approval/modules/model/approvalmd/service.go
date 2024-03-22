package approvalmd

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
	"time"
)

var (
	validPidPattern = regexp.MustCompile(`^\w{1,32}$`)
)

func IsPidValid(pid string) bool {
	return validPidPattern.MatchString(pid)
}

func IsProcessNameValid(name string) bool {
	return len(name) <= 32
}

func IsGroupNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func InsertProcess(ctx context.Context, reqDTO InsertProcessReqDTO) error {
	jsonBytes, _ := json.Marshal(reqDTO.Process)
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Process{
		Pid:        reqDTO.Pid,
		Name:       reqDTO.Name,
		GroupId:    reqDTO.GroupId,
		Content:    string(jsonBytes),
		IconUrl:    reqDTO.IconUrl,
		SourceType: reqDTO.SourceType,
	})
	return err
}

func UpdateProcessById(ctx context.Context, reqDTO UpdateProcessByIdReqDTO) (bool, error) {
	jsonBytes, _ := json.Marshal(reqDTO.Process)
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("approval", "name", "group_id", "icon_url", "content").
		Limit(1).
		Update(&Process{
			Name:    reqDTO.Name,
			GroupId: reqDTO.GroupId,
			Content: string(jsonBytes),
			IconUrl: reqDTO.IconUrl,
		})
	return rows == 1, err
}

func UpdateProcessByPid(ctx context.Context, reqDTO UpdateProcessByPidReqDTO) (bool, error) {
	jsonBytes, _ := json.Marshal(reqDTO.Process)
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("pid = ?", reqDTO.Pid).
		Cols("approval", "name", "group_id").
		Update(&Process{
			Name:    reqDTO.Name,
			GroupId: reqDTO.GroupId,
			Content: string(jsonBytes),
		})
	return rows == 1, err
}

func DeleteProcessById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(Process))
	return rows == 1, err
}

func DeleteProcessByPid(ctx context.Context, pid string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("pid = ?", pid).
		Limit(1).
		Delete(new(Process))
	return rows == 1, err
}

func GetProcessByPid(ctx context.Context, pid string) (Process, bool, error) {
	var ret Process
	b, err := xormutil.MustGetXormSession(ctx).
		Where("pid = ?", pid).
		Get(&ret)
	return ret, b, err
}

func GetProcessById(ctx context.Context, id int64) (Process, bool, error) {
	var ret Process
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func GetProcessByGroupId(ctx context.Context, groupId int64) ([]Process, error) {
	ret := make([]Process, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("group_id = ?", groupId).
		Find(&ret)
	return ret, err
}

func InsertFlow(ctx context.Context, reqDTO InsertFlowReqDTO) (Flow, error) {
	processJson, _ := json.Marshal(reqDTO.Process)
	kvsJson, _ := json.Marshal(reqDTO.Kvs)
	var ret = Flow{
		ProcessName:    reqDTO.ProcessName,
		ProcessId:      reqDTO.ProcessId,
		ProcessContent: string(processJson),
		CurrIndex:      reqDTO.CurrIndex,
		FlowStatus:     reqDTO.FlowStatus,
		Creator:        reqDTO.Creator,
		BizId:          reqDTO.BizId,
		Kvs:            string(kvsJson),
	}
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&ret)
	return ret, err
}

func UpdateFlowStatusWithOldStatus(ctx context.Context, flowId int64, newStatus, oldStatus FlowStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", flowId).
		And("flow_status = ?", oldStatus).
		Cols("flow_status").
		Update(&Flow{
			FlowStatus: newStatus,
		})
	return rows == 1, err
}

func UpdateFlowStatusAndErrMsgWithOldStatus(ctx context.Context, flowId int64, errMsg string, newStatus, oldStatus FlowStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", flowId).
		And("flow_status = ?", oldStatus).
		Cols("flow_status", "err_msg").
		Update(&Flow{
			FlowStatus: newStatus,
			ErrMsg:     errMsg,
		})
	return rows == 1, err
}

func UpdateFlowCurrIndexWithOldCurrIndex(ctx context.Context, flowId int64, newCurrIndex, oldCurrIndex int) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", flowId).
		And("curr_index = ?", oldCurrIndex).
		Cols("curr_index").
		Update(&Flow{
			CurrIndex: newCurrIndex,
		})
	return rows == 1, err
}

func InsertNotify(ctx context.Context, reqDTO InsertNotifyReqDTO) error {
	notifyList := make([]Notify, 0, len(reqDTO.Accounts))
	for _, account := range reqDTO.Accounts {
		notifyList = append(notifyList, Notify{
			FlowId:    reqDTO.FlowId,
			Account:   account,
			FlowIndex: reqDTO.FlowIndex,
			Done:      reqDTO.Done,
		})
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(notifyList)
	return err
}

func GetFlowById(ctx context.Context, flowId int64) (Flow, bool, error) {
	var ret Flow
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", flowId).Get(&ret)
	return ret, b, err
}

func GetNotifyById(ctx context.Context, notifyId int64) (Notify, bool, error) {
	var ret Notify
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", notifyId).Get(&ret)
	return ret, b, err
}

func CountOperate(ctx context.Context, flowId int64, flowIndex int, flowOp FlowOp) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("flow_id = ?", flowId).
		And("flow_index = ?", flowIndex).
		And("op = ?", flowOp).Count(new(Notify))
}

func UpdateNotifyDoneWithOldDone(ctx context.Context, notifyId int64, newDone, oldDone bool, op FlowOp) (bool, error) {
	oldDoneInt := 0
	if oldDone {
		oldDoneInt = 1
	}
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", notifyId).
		And("done = ?", oldDoneInt).
		Cols("done", "op").
		Limit(1).
		Update(&Notify{
			Done: newDone,
			Op:   op,
		})
	return rows == 1, err
}

func InsertGroup(ctx context.Context, reqDTO InsertGroupReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Group{
			Name: reqDTO.Name,
		})
	return err
}

func DeleteGroup(ctx context.Context, groupId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", groupId).Limit(1).
		Delete(new(Group))
	return rows == 1, err
}

func UpdateGroup(ctx context.Context, reqDTO UpdateGroupReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name").
		Update(&Group{
			Name: reqDTO.Name,
		})
	return rows == 1, err
}

func GetGroupById(ctx context.Context, id int64) (Group, bool, error) {
	var ret Group
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func GetFlowByCreatorAndTime(ctx context.Context, creator string, startTime, endTime time.Time) ([]Flow, error) {
	ret := make([]Flow, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("creator = ?", creator).
		And("created between ? and ?",
			startTime.Format(time.DateTime),
			endTime.Format(time.DateTime)).
		OrderBy("updated desc").
		Find(&ret)
	return ret, err
}

func BatchGetFlows(ctx context.Context, flowIds []int64) ([]Flow, error) {
	ret := make([]Flow, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", flowIds).
		OrderBy("updated desc").
		Find(&ret)
	return ret, err
}

func GetNotifyByAccountAndTime(ctx context.Context, account string, startTime, endTime time.Time, done bool) ([]Notify, error) {
	ret := make([]Notify, 0)
	doneInt := 0
	if done {
		doneInt = 1
	}
	err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		And("done = ?", doneInt).
		And("created between ? and ?",
			startTime.Format(time.DateTime),
			endTime.Format(time.DateTime)).
		Find(&ret)
	return ret, err
}

func UpdateNotifyDoneWithOldDoneByFlowId(ctx context.Context, newDone, oldDone bool, op FlowOp, flowId int64) error {
	oldDoneInt := 0
	if oldDone {
		oldDoneInt = 1
	}
	_, err := xormutil.MustGetXormSession(ctx).Where("flow_id = ?", flowId).
		Where("done = ?", oldDoneInt).
		Cols("done", "op").
		Update(&Notify{
			Done: newDone,
			Op:   op,
		})
	return err
}

func UpdateNotifyDoneWithOldDoneByFlowIdAndIndex(ctx context.Context, newDone, oldDone bool, op FlowOp, flowId int64, flowIndex int) error {
	oldDoneInt := 0
	if oldDone {
		oldDoneInt = 1
	}
	_, err := xormutil.MustGetXormSession(ctx).
		Where("flow_id = ?", flowId).
		Where("done = ?", oldDoneInt).
		And("flow_index = ?", flowIndex).
		Cols("done", "op").
		Update(&Notify{
			Done: newDone,
			Op:   op,
		})
	return err
}

func GetNotifyByFlowId(ctx context.Context, flowId int64) ([]Notify, error) {
	ret := make([]Notify, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("flow_id = ?", flowId).
		Find(&ret)
	return ret, err
}

func ExistNotifyByAccountAndFlowId(ctx context.Context, account string, flowId int64) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		And("flow_id = ?", flowId).
		Exist(new(Notify))
}

func GetAllGroups(ctx context.Context) ([]Group, error) {
	ret := make([]Group, 0)
	err := xormutil.MustGetXormSession(ctx).Find(&ret)
	return ret, err
}

func GetAllCustomProcesses(ctx context.Context) ([]SimpleProcess, error) {
	ret := make([]SimpleProcess, 0)
	err := xormutil.MustGetXormSession(ctx).Where("source_type = ?", CustomSourceType).Find(&ret)
	return ret, err
}
