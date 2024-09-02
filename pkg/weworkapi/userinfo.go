package weworkapi

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/httputil"
	"net/http"
)

type GetAuthWeworkUserInfoResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	UserId  string `json:"userid"`
}

func GetAuthWeworkUserInfo(ctx context.Context, url, accessToken, code string) (string, error) {
	if url == "" {
		url = "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo"
	}
	var res GetAuthWeworkUserInfoResp
	err := httputil.Get(ctx, http.DefaultClient, url+fmt.Sprintf("?access_token=%s&code=%s", accessToken, code), nil, &res)
	if err != nil {
		return "", err
	}
	if res.ErrCode != 0 {
		return "", fmt.Errorf("http resp err code: %v msg: %v", res.ErrCode, res.ErrMsg)
	}
	return res.UserId, nil
}
