package apisession

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zsf/http/httptask"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"net/url"
	"time"
)

const (
	SessionTableName = "zall_login_session"
)

type SessionModel struct {
	Id        int64     `xorm:"pk autoincr"`
	SessionId string    `json:"sessionId"`
	Account   string    `json:"account"`
	UserInfo  string    `json:"userInfo"`
	ExpireAt  int64     `json:"expireAt"`
	Created   time.Time `json:"created" xorm:"created"`
	Updated   time.Time `json:"updated" xorm:"updated"`
}

func (s *SessionModel) TableName() string {
	return SessionTableName
}

func (s *SessionModel) ToSession() Session {
	ret := Session{
		SessionId: s.SessionId,
		ExpireAt:  s.ExpireAt,
	}
	if s.UserInfo != "" {
		_ = json.Unmarshal([]byte(s.UserInfo), &ret.UserInfo)
	}
	return ret
}

type mysqlStore struct{}

func newMysqlStore() *mysqlStore {
	ret := new(mysqlStore)
	httptask.AppendHttpTask("clearExpiredLoginSession", func(_ []byte, _ url.Values) {
		ret.ClearExpired()
	})
	return ret
}

func (m *mysqlStore) GetBySessionId(sessionId string) (Session, bool, error) {
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	ret := SessionModel{}
	b, err := xsess.Where("session_id = ?", sessionId).And("expire_at > ?", time.Now().UnixMilli()).Get(&ret)
	return ret.ToSession(), b, err
}

func (m *mysqlStore) GetByAccount(account string) (Session, bool, error) {
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	ret := SessionModel{}
	b, err := xsess.Where("account = ?", account).And("expire_at > ?", time.Now().UnixMilli()).Get(&ret)
	return ret.ToSession(), b, err
}

func (m *mysqlStore) PutSession(session Session) error {
	userInfoJson, _ := json.Marshal(session.UserInfo)
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	_, err := xsess.Insert(&SessionModel{
		SessionId: session.SessionId,
		Account:   session.UserInfo.Account,
		UserInfo:  string(userInfoJson),
		ExpireAt:  session.ExpireAt,
	})
	return err
}

func (m *mysqlStore) DeleteByAccount(account string) error {
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	_, err := xsess.Where("account = ?", account).Delete(new(SessionModel))
	return err
}

func (m *mysqlStore) DeleteBySessionId(sessionId string) error {
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	_, err := xsess.Where("session_id = ?", sessionId).Delete(new(SessionModel))
	return err
}

func (m *mysqlStore) RefreshExpiry(sessionId string, expireAt int64) error {
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	_, err := xsess.Where("session_id = ?", sessionId).
		Cols("expire_at").
		Update(&SessionModel{
			ExpireAt: expireAt,
		})
	return err
}

func (m *mysqlStore) ClearExpired() {
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	_, err := xsess.Where("expire_at <= ?", time.Now().UnixMilli()).Delete(new(SessionModel))
	if err != nil {
		logger.Logger.Error(err)
	}
}
