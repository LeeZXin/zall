package apisession

import (
	"context"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"sync"
	"time"
)

// memStore 内存
type memStore struct {
	sync.RWMutex
	session           map[string]Session
	userSession       map[string]Session
	clearTaskStopFunc taskutil.StopFunc
}

func newMemStore() Store {
	m := &memStore{
		RWMutex:     sync.RWMutex{},
		session:     make(map[string]Session, 8),
		userSession: make(map[string]Session, 8),
	}
	m.clearTaskStopFunc, _ = taskutil.RunPeriodicalTask(
		10*time.Minute,
		10*time.Minute,
		func(context.Context) {
			m.ClearExpired()
		})
	quit.AddShutdownHook(quit.ShutdownHook(m.clearTaskStopFunc))
	return m
}

func (s *memStore) ClearExpired() {
	s.Lock()
	defer s.Unlock()
	now := time.Now().UnixMilli()
	for _, session := range s.session {
		if session.ExpireAt < now {
			delete(s.session, session.SessionId)
			delete(s.session, session.UserInfo.Account)
		}
	}
}

func (s *memStore) GetBySessionId(sessionId string) (Session, bool, error) {
	s.RLock()
	defer s.RUnlock()
	ret, b := s.session[sessionId]
	if ret.ExpireAt < time.Now().UnixMilli() {
		return Session{}, false, nil
	}
	return ret, b, nil
}

func (s *memStore) GetByAccount(account string) (Session, bool, error) {
	s.RLock()
	defer s.RUnlock()
	ret, b := s.userSession[account]
	if ret.ExpireAt < time.Now().UnixMilli() {
		return Session{}, false, nil
	}
	return ret, b, nil
}

func (s *memStore) PutSession(session Session) error {
	s.Lock()
	defer s.Unlock()
	s.session[session.SessionId] = session
	s.userSession[session.UserInfo.Account] = session
	return nil
}

func (s *memStore) DeleteByAccount(account string) error {
	s.Lock()
	defer s.Unlock()
	session, b := s.userSession[account]
	if b {
		delete(s.userSession, account)
		delete(s.session, session.SessionId)
	}
	return nil
}

func (s *memStore) DeleteBySessionId(sessionId string) error {
	s.Lock()
	defer s.Unlock()
	session, b := s.session[sessionId]
	if b {
		delete(s.userSession, session.UserInfo.Account)
		delete(s.session, sessionId)
	}
	return nil
}

func (s *memStore) RefreshExpiry(sessionId string, expireAt int64) error {
	s.Lock()
	defer s.Unlock()
	session, b := s.session[sessionId]
	if b {
		session.ExpireAt = expireAt
		s.session[sessionId] = session
		s.userSession[session.UserInfo.Account] = session
	}
	return nil
}
