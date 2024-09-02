package feishuapi

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/httputil"
	"net/http"
)

type GetAppAccessTokenResp struct {
	AppAccessToken    string `json:"app_access_token"`
	Code              int    `json:"code"`
	Expire            int    `json:"expire"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}

func GetAppAccessToken(ctx context.Context, url, appId, appSecret string) (string, string, int, error) {
	if url == "" {
		url = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"
	}
	var res GetAppAccessTokenResp
	err := httputil.Post(ctx, http.DefaultClient, url, nil, map[string]string{
		"app_id":     appId,
		"app_secret": appSecret,
	}, &res)
	if err != nil {
		return "", "", 0, err
	}
	if res.Code != 0 {
		return "", "", 0, fmt.Errorf("http resp err code: %v msg: %v", res.Code, res.Msg)
	}
	return res.AppAccessToken, res.TenantAccessToken, res.Expire, nil
}

type GetUserAccessTokenResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data UserAccessToken `json:"data"`
}

type UserAccessToken struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
}

func GetUserAccessToken(ctx context.Context, url, appAccessToken, code string) (UserAccessToken, error) {
	if url == "" {
		url = "https://open.feishu.cn/open-apis/authen/v1/oidc/access_token"
	}
	var res GetUserAccessTokenResp
	err := httputil.Post(ctx, http.DefaultClient, url, map[string]string{
		"Authorization": "Bearer " + appAccessToken,
	}, map[string]string{
		"grant_type": "authorization_code",
		"code":       code,
	}, &res)
	if err != nil {
		return UserAccessToken{}, err
	}
	if res.Code != 0 {
		return UserAccessToken{}, fmt.Errorf("http resp err code: %v msg: %v", res.Code, res.Msg)
	}
	return res.Data, nil
}
