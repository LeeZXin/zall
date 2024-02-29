package cfgapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type GetSysCfgRespVO struct {
	ginutil.BaseResp
	Cfg cfgsrv.SysCfg `json:"cfg"`
}

type UpdateSysCfgReqVO struct {
	cfgsrv.SysCfg
}
