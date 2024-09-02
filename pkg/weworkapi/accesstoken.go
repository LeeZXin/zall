package weworkapi

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/httputil"
	"net/http"
)

type GetAccessTokenResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetAccessToken(ctx context.Context, url, corpId, secret string) (string, int, error) {
	if url == "" {
		url = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	}
	var res GetAccessTokenResp
	err := httputil.Get(ctx, http.DefaultClient, url+fmt.Sprintf("?corpid=%s&corpsecret=%s", corpId, secret), nil, &res)
	if err != nil {
		return "", 0, err
	}
	if res.ErrCode != 0 {
		return "", 0, fmt.Errorf("http resp err code: %v msg: %v", res.ErrCode, res.ErrMsg)
	}
	return res.AccessToken, res.ExpiresIn, nil
}
