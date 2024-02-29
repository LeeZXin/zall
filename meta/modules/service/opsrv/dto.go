package opsrv

import (
	"time"
)

type InsertOpLogReqDTO struct {
	Account    string
	OpDesc     string
	EventTime  time.Time
	ReqContent any
	Err        error
}
