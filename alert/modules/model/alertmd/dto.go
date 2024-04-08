package alertmd

import "github.com/LeeZXin/zall/pkg/alert"

type InsertConfigReqDTO struct {
	Name        string
	Alert       alert.Alert
	AppId       string
	IntervalSec int
	Enabled     bool
	NextTime    int64
}

type UpdateConfigReqDTO struct {
	Id          int64
	Name        string
	Alert       alert.Alert
	IntervalSec int
	Enabled     bool
	NextTime    int64
}

type ListConfigReqDTO struct {
	Cursor int64
	Limit  int
	TeamId int64
}
