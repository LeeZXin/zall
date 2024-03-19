package approvalmd

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
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

func InsertProcess(ctx context.Context, reqDTO InsertProcessReqDTO) error {
	jsonBytes, _ := json.Marshal(reqDTO.Process)
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Process{
		Pid:     reqDTO.Pid,
		Name:    reqDTO.Name,
		GroupId: reqDTO.GroupId,
		Content: string(jsonBytes),
	})
	return err
}

func UpdateProcessById(ctx context.Context, reqDTO UpdateProcessByIdReqDTO) (bool, error) {
	jsonBytes, _ := json.Marshal(reqDTO.Process)
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("approval", "name", "group_id").
		Update(&Process{
			Name:    reqDTO.Name,
			GroupId: reqDTO.GroupId,
			Content: string(jsonBytes),
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

func DeleteProcess(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("approval").
		Delete(new(Process))
	return rows == 1, err
}

func GetProcessByPid(ctx context.Context, pid string) (ProcessDTO, bool, error) {
	var ret Process
	b, err := xormutil.MustGetXormSession(ctx).
		Where("pid = ?", pid).
		Get(&ret)
	if err != nil || !b {
		return ProcessDTO{}, b, err
	}
	return processMd2Dto(ret), true, nil
}

func GetProcessById(ctx context.Context, id int64) (ProcessDTO, bool, error) {
	var ret Process
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	if err != nil || !b {
		return ProcessDTO{}, b, err
	}
	return processMd2Dto(ret), true, nil
}

func processMd2Dto(p Process) ProcessDTO {
	process, _ := p.GetUnmarshalProcess()
	return ProcessDTO{
		Id:      p.Id,
		Pid:     p.Pid,
		Process: process,
		GroupId: p.GroupId,
		Created: p.Created,
	}
}

func InsertFlow(ctx context.Context, reqDTO InsertFlowReqDTO) (Flow, error) {
	processJson, _ := json.Marshal(reqDTO.Process)
	kvsJson, _ := json.Marshal(reqDTO.Kvs)
	var ret = Flow{
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

func CountDetail(ctx context.Context, flowId int64, flowIndex int, flowOp FlowOp) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("flow_id = ?", flowId).
		And("flow_index = ?", flowIndex).
		And("op = ?", flowOp).Count(new(Detail))
}

func ExistsDetailByAccount(ctx context.Context, flowId int64, flowIndex int, account string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("flow_id = ?", flowId).
		And("flow_index = ?", flowIndex).
		And("account = ?", account).
		Exist(new(Detail))
}

func InsertDetail(ctx context.Context, reqDTO InsertDetailReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Detail{
		FlowId:    reqDTO.FlowId,
		Account:   reqDTO.Account,
		FlowIndex: reqDTO.FlowIndex,
		Op:        reqDTO.FlowOp,
	})
	return err
}

func UpdateNotifyDone(ctx context.Context, notifyId int64, done bool) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", notifyId).
		Cols("done").
		Update(&Notify{
			Done: done,
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
