package usermd

type UserInfo struct {
	Account      string   `json:"account"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	IsProhibited bool     `json:"isProhibited"`
	AvatarUrl    string   `json:"avatarUrl"`
	IsAdmin      bool     `json:"isAdmin"`
	RoleType     RoleType `json:"roleType"`
}

func (i *UserInfo) IsValid() bool {
	return i.Account != ""
}

type InsertUserReqDTO struct {
	Account   string
	Name      string
	Email     string
	Password  string
	AvatarUrl string
	IsAdmin   bool
	RoleType  RoleType
}

type UpdateUserReqDTO struct {
	Account string
	Name    string
	Email   string
}

type UpdateAdminReqDTO struct {
	Account string
	IsAdmin bool
}

type UpdatePasswordReqDTO struct {
	Account  string
	Password string
}

type ListUserReqDTO struct {
	Account string
	Cursor  int64
	Limit   int
}

type SetUserProhibitedReqDTO struct {
	Account      string
	IsProhibited bool
}
