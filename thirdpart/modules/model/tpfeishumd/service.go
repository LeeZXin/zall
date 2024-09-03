package tpfeishumd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"xorm.io/builder"
)

func IsAccessTokenNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsAccessTokenAppIdValid(appId string) bool {
	return len(appId) > 0 && len(appId) <= 32
}

func IsAccessTokenSecretValid(secret string) bool {
	return len(secret) > 0 && len(secret) <= 64
}

func InsertAccessToken(ctx context.Context, reqDTO InsertAccessTokenReqDTO) (AccessToken, error) {
	ret := AccessToken{
		TeamId:      reqDTO.TeamId,
		Name:        reqDTO.Name,
		AppId:       reqDTO.AppId,
		Secret:      reqDTO.Secret,
		Token:       reqDTO.Token,
		ExpireTime:  reqDTO.ExpireTime,
		ApiKey:      reqDTO.ApiKey,
		TenantToken: reqDTO.TenantToken,
		Creator:     reqDTO.Creator,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateAccessTokenToken(ctx context.Context, reqDTO UpdateAccessTokenTokenReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("token", "expire_time", "tenant_token").
		Update(&AccessToken{
			Token:       reqDTO.Token,
			ExpireTime:  reqDTO.ExpireTime,
			TenantToken: reqDTO.TenantToken,
		})
	return rows == 1, err
}

func UpdateAccessTokenApiKeyById(ctx context.Context, id int64, apiKey string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Cols("api_key").
		Update(&AccessToken{
			ApiKey: apiKey,
		})
	return rows == 1, err
}

func GetAccessTokenByAppIdAndSecret(ctx context.Context, appId, secret string) (AccessToken, bool, error) {
	var ret AccessToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("secret = ?", secret).
		Get(&ret)
	return ret, b, err
}

func GetAccessTokenByApiKey(ctx context.Context, apiKey string) (AccessToken, bool, error) {
	var ret AccessToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("api_key = ?", apiKey).
		Get(&ret)
	return ret, b, err
}

func DeleteAccessTokenById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(AccessToken))
	return rows == 1, err
}

func IterateAccessToken(ctx context.Context, expireTime int64, fn func(*AccessToken) error) error {
	return xormutil.MustGetXormSession(ctx).
		Where("expire_time <= ?", expireTime).
		Iterate(new(AccessToken), func(_ int, bean any) error {
			return fn(bean.(*AccessToken))
		})
}

func GetAccessTokenById(ctx context.Context, id int64) (AccessToken, bool, error) {
	var ret AccessToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func UpdateAccessToken(ctx context.Context, reqDTO UpdateAccessTokenReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("app_id", "secret", "name").
		Update(&AccessToken{
			AppId:  reqDTO.AppId,
			Secret: reqDTO.Secret,
			Name:   reqDTO.Name,
		})
	return rows == 1, err
}

func ExistAccessTokenByAppIdAndSecret(ctx context.Context, appId, secret string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("secret = ?", secret).
		Exist(new(AccessToken))
}

func ListAccessToken(ctx context.Context, reqDTO ListAccessTokenReqDTO) ([]AccessToken, int64, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId)
	if reqDTO.Key != "" {
		session.And(
			builder.Or(
				builder.Expr("name like ?", reqDTO.Key+"%"),
				builder.Expr("app_id like ?", reqDTO.Key+"%"),
			),
		)
	}
	ret := make([]AccessToken, 0)
	total, err := session.
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		FindAndCount(&ret)
	return ret, total, err
}
