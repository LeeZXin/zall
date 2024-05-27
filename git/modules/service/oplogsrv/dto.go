package oplogsrv

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"time"
)

type OpLog struct {
	RepoId    int64
	Operator  string
	Log       string
	Req       any
	EventTime time.Time
}

type PageOpLogReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Account  string              `json:"account"`
	DateStr  string              `json:"dayStr"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
	DateTime time.Time           `json:"-"`
}

func (r *PageOpLogReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.Account) > 255 {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	var err error
	r.DateTime, err = time.Parse(time.DateOnly, r.DateStr)
	if err != nil {
		return util.InvalidArgsError()
	}
	return nil
}

type OpLogDTO struct {
	Id      int64
	Account string
	Created time.Time
	Content string
	ReqBody string
}
