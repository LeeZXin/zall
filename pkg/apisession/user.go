package apisession

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
