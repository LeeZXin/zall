package tpweworkmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertAccessToken(ctx context.Context, reqDTO InsertAccessTokenReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&AccessToken{
			CorpId:     reqDTO.CorpId,
			Secret:     reqDTO.Secret,
			Token:      reqDTO.Token,
			ExpireTime: reqDTO.ExpireTime,
			ApiKey:     reqDTO.ApiKey,
		})
	return err
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

func GetAccessTokenByApiKey(ctx context.Context, apiKey string) (AccessToken, bool, error) {
	var ret AccessToken
	b, err := xormutil.MustGetXormSession(ctx).
		Where("api_key = ?", apiKey).
		Get(&ret)
	return ret, b, err
}

func DeleteAccessById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(AccessToken))
	return rows == 1, err
}

func BatchGetAccessTokenByIdList(ctx context.Context, idList []int64) ([]AccessToken, error) {
	ret := make([]AccessToken, 0)
	err := xormutil.MustGetXormSession(ctx).In("id", idList).Find(&ret)
	return ret, err
}

func InsertCollaborator(ctx context.Context, reqDTO InsertCollaboratorReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Collaborator{
			TokenId:   reqDTO.TokenId,
			Account:   reqDTO.Account,
			IsCreator: reqDTO.IsCreator,
		})
	return err
}

func ListCollaborator(ctx context.Context, reqDTO ListCollaboratorReqDTO) ([]Collaborator, error) {
	ret := make([]Collaborator, 0)
	session := xormutil.MustGetXormSession(ctx)
	if reqDTO.Account != "" {
		session.And("account = ?", reqDTO.Account)
	}
	err := session.
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		Find(&ret)
	return ret, err
}

func IterateAccessToken(ctx context.Context, expireTime int64, fn func(*AccessToken) error) error {
	return xormutil.MustGetXormSession(ctx).
		Where("expire_time <= ?", expireTime).
		Iterate(new(AccessToken), func(_ int, bean any) error {
			return fn(bean.(*AccessToken))
		})
}
