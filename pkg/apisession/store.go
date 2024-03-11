package apisession

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/property/static"
	"strconv"
	"time"
)

const (
	MemStoreType  = "mem"
	xormstoreType = "mysql"
)

var (
	storeImpl Store
)

func init() {
	switch static.GetString("apisession.store.type") {
	case MemStoreType:
		storeImpl = newMemStore()
	case xormstoreType:
		storeImpl = newMysqlStore()
	default:
		storeImpl = newMysqlStore()
	}
}

type Session struct {
	SessionId string   `json:"sessionId"`
	UserInfo  UserInfo `json:"userInfo"`
	ExpireAt  int64    `json:"expireAt"`
}

type Store interface {
	GetBySessionId(string) (Session, bool, error)
	GetByAccount(string) (Session, bool, error)
	PutSession(Session) error
	DeleteByAccount(string) error
	DeleteBySessionId(string) error
	RefreshExpiry(string, int64) error
	ClearExpired()
}

func GetStore() Store {
	return storeImpl
}

func GenSessionId() string {
	h := sha256.New()
	h.Write([]byte(idutil.RandomUuid() + strconv.FormatInt(time.Now().UnixNano(), 10)))
	return hex.EncodeToString(h.Sum(nil))
}
