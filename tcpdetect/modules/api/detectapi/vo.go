package detectapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type InsertDetectReqVO struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	Name string `json:"name"`
}

type UpdateDetectReqVO struct {
	Id   int64  `json:"id"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	Name string `json:"name"`
}

type ListDetectReqVO struct {
	Name   string `json:"name"`
	Cursor int64  `json:"cursor"`
	Limit  int    `json:"limit"`
}

type DeleteDetectReqVO struct {
	Id int64 `json:"id"`
}

type DetectVO struct {
	Id            int64  `json:"id"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	Name          string `json:"name"`
	HeartbeatTime string `json:"heartbeatTime"`
	Valid         bool   `json:"valid"`
	Enabled       bool   `json:"enabled"`
}

type ListDetectRespVO struct {
	ginutil.BaseResp
	Data   []DetectVO `json:"data"`
	Cursor int64      `json:"cursor"`
}

type ListLogReqVO struct {
	Id     int64 `json:"id"`
	Cursor int64 `json:"cursor"`
	Limit  int   `json:"limit"`
}

type LogVO struct {
	Ip      string `json:"ip"`
	Port    int    `json:"port"`
	Created string `json:"created"`
	Valid   bool   `json:"valid"`
}

type ListLogRespVO struct {
	ginutil.BaseResp
	Data   []LogVO `json:"data"`
	Cursor int64   `json:"cursor"`
}

type EnableDetectReqVO struct {
	Id int64 `json:"id"`
}

type DisableDetectReqVO struct {
	Id int64 `json:"id"`
}
