package apisession

import (
	"context"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

const (
	SessionTableName = "zall_login_session"
)

type SessionModel struct {
	Id        int64     `json:"id" xorm:"pk autoincr"`
	SessionId string    `json:"sessionId"`
	Account   string    `json:"account"`
	UserInfo  *UserInfo `json:"userInfo"`
	ExpireAt  int64     `json:"expireAt"`
	Created   time.Time `json:"created" xorm:"created"`
	Updated   time.Time `json:"updated" xorm:"updated"`
}

func (s *SessionModel) TableName() string {
	return SessionTableName
}

func (s *SessionModel) ToSession() Session {
	userInfo := s.UserInfo
	ret := Session{
		SessionId: s.SessionId,
		ExpireAt:  s.ExpireAt,
		UserInfo:  *userInfo,
	}
	return ret
}

type mysqlStore struct{}

func newMysqlStore() *mysqlStore {
	return new(mysqlStore)
}

func (m *mysqlStore) GetBySessionId(sessionId string) (Session, bool, error) {
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	ret := SessionModel{}
	b, err := xsess.Where("session_id = ?", sessionId).
		And("expire_at > ?", time.Now().UnixMilli()).
		Get(&ret)
	if err != nil || !b {
		return Session{}, b, err
	}
	return ret.ToSession(), b, err
}

func (m *mysqlStore) GetByAccount(account string) (Session, bool, error) {
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	ret := SessionModel{}
	b, err := xsess.Where("account = ?", account).
		And("expire_at > ?", time.Now().UnixMilli()).
		Get(&ret)
	return ret.ToSession(), b, err
}

func (m *mysqlStore) PutSession(session Session) error {
	xsess := xormstore.NewXormSession(context.Background())
	defer xsess.Close()
	var existModel SessionModel
	b, err := xsess.Where("session_id = ?", session.SessionId).Get(&existModel)
	if err != nil {
		return err
	}
	if !b {
		_, err = xsess.Insert(&SessionModel{
			SessionId: session.SessionId,
			Account:   session.UserInfo.Account,
			UserInfo:  &session.UserInfo,
			ExpireAt:  session.ExpireAt,
		})
	} else {
		_, err = xsess.Where("session_id = ?", session.SessionId).
			Cols("user_info", "expire_at").
			Update(&SessionModel{
				UserInfo: &session.UserInfo,
				ExpireAt: session.ExpireAt,
			})
	}
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
