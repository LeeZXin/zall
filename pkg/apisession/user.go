package apisession

import "encoding/json"

type UserInfo struct {
	Account      string `json:"account"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsProhibited bool   `json:"isProhibited"`
	AvatarUrl    string `json:"avatarUrl"`
	IsAdmin      bool   `json:"isAdmin"`
}

func (i *UserInfo) IsValid() bool {
	return i.Account != ""
}

func (i *UserInfo) FromDB(content []byte) error {
	if i == nil {
		*i = UserInfo{}
	}
	return json.Unmarshal(content, i)
}

func (i *UserInfo) ToDB() ([]byte, error) {
	return json.Marshal(i)
}
