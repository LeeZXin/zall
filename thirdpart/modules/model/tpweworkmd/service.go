package tpweworkmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"xorm.io/builder"
)

func IsAccessTokenNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func IsAccessTokenCorpIdValid(corpId string) bool {
	return len(corpId) > 0 && len(corpId) <= 32
}

func IsAccessTokenSecretValid(secret string) bool {
	return len(secret) > 0 && len(secret) <= 64
}

func InsertAccessToken(ctx context.Context, reqDTO InsertAccessTokenReqDTO) (AccessToken, error) {
	ret := AccessToken{
		Name:       reqDTO.Name,
		TeamId:     reqDTO.TeamId,
		CorpId:     reqDTO.CorpId,
		Secret:     reqDTO.Secret,
		Token:      reqDTO.Token,
		ExpireTime: reqDTO.ExpireTime,
		ApiKey:     reqDTO.ApiKey,
		Creator:    reqDTO.Creator,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateAccessTokenToken(ctx context.Context, reqDTO UpdateAccessTokenTokenReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("token", "expire_time").
		Update(&AccessToken{
			Token:      reqDTO.Token,
			ExpireTime: reqDTO.ExpireTime,
		})
	return rows == 1, err
}

func UpdateAccessToken(ctx context.Context, reqDTO UpdateAccessTokenReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("corp_id", "secret", "name").
		Update(&AccessToken{
			CorpId: reqDTO.CorpId,
			Secret: reqDTO.Secret,
			Name:   reqDTO.Name,
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

func GetAccessTokenByCorpIdAndSecret(ctx context.Context, corpId, secret string) (AccessToken, bool, error) {
	var ret AccessToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("corp_id = ?", corpId).
		And("secret = ?", secret).
		Get(&ret)
	return ret, b, err
}

func ExistAccessTokenByCorpIdAndSecret(ctx context.Context, corpId, secret string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("corp_id = ?", corpId).
		And("secret = ?", secret).
		Exist(new(AccessToken))
}

func GetAccessTokenByApiKey(ctx context.Context, apiKey string) (AccessToken, bool, error) {
	var ret AccessToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("api_key = ?", apiKey).
		Get(&ret)
	return ret, b, err
}

func GetAccessTokenById(ctx context.Context, id int64) (AccessToken, bool, error) {
	var ret AccessToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
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

func ListAccessToken(ctx context.Context, reqDTO ListAccessTokenReqDTO) ([]AccessToken, int64, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId)
	if reqDTO.Key != "" {
		session.And(
			builder.Or(
				builder.Expr("name like ?", reqDTO.Key+"%"),
				builder.Expr("corp_id like ?", reqDTO.Key+"%"),
			),
		)
	}
	ret := make([]AccessToken, 0)
	total, err := session.
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		FindAndCount(&ret)
	return ret, total, err
}
