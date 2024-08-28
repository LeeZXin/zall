package alertmd

import "github.com/LeeZXin/zall/pkg/alert"

type InsertConfigReqDTO struct {
	Name        string
	Alert       alert.Alert
	AppId       string
	IntervalSec int
	IsEnabled   bool
	Env         string
	Creator     string
}

type UpdateConfigReqDTO struct {
	Id          int64
	Name        string
	Alert       alert.Alert
	IntervalSec int
}

type ListConfigReqDTO struct {
	PageNum  int
	PageSize int
	AppId    string
	Env      string
}

type InsertExecuteReqDTO struct {
	ConfigId  int64
	IsEnabled bool
	NextTime  int64
	Env       string
}
