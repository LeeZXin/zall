package approvalapi

type AgreeApprovalReqVO struct {
	NotifyId int64 `json:"notifyId"`
}

type DisagreeApprovalReqVO struct {
	NotifyId int64 `json:"notifyId"`
}
