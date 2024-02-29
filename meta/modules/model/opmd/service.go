package opmd

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertLog(ctx context.Context, reqDTO InsertLogReqDTO) error {
	ret := OpLog{
		Operator: reqDTO.Operator,
		OpDesc:   reqDTO.OpDesc,
		ErrMsg:   reqDTO.ErrMsg,
	}
	if reqDTO.ReqContent != nil {
		req, err := json.Marshal(reqDTO.ReqContent)
		if err == nil {
			ret.ReqContent = string(req)
		} else {
			ret.ReqContent = err.Error()
		}
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(ret)
	return err
}
