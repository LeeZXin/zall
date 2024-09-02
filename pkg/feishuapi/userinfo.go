package feishuapi

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/httputil"
	"net/http"
)

type GetUserInfoResp struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data UserInfo `json:"data"`
}

type UserInfo struct {
	Name            string `json:"name"`
	EnName          string `json:"en_name"`
	AvatarUrl       string `json:"avatar_url"`
	AvatarThumb     string `json:"avatar_thumb"`
	AvatarMiddle    string `json:"avatar_middle"`
	AvatarBig       string `json:"avatar_big"`
	OpenId          string `json:"open_id"`
	UnionId         string `json:"union_id"`
	Email           string `json:"email"`
	EnterpriseEmail string `json:"enterprise_email"`
	UserId          string `json:"user_id"`
	Mobile          string `json:"mobile"`
	TenantKey       string `json:"tenant_key"`
	EmployeeNo      string `json:"employee_no"`
}

func GetUserInfo(ctx context.Context, url, tenantToken string) (UserInfo, error) {
	if url == "" {
		url = "https://open.feishu.cn/open-apis/authen/v1/user_info"
	}
	var res GetUserInfoResp
	err := httputil.Get(ctx, http.DefaultClient, url, map[string]string{
		"Authorization": "Bearer " + tenantToken,
	}, &res)
	if err != nil {
		return UserInfo{}, err
	}
	if res.Code != 0 {
		return UserInfo{}, fmt.Errorf("http resp err code: %v msg: %v", res.Code, res.Msg)
	}
	return res.Data, nil
}
