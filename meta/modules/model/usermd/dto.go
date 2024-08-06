package usermd

type UserInfo struct {
	Account      string `json:"account"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsProhibited bool   `json:"isProhibited"`
	AvatarUrl    string `json:"avatarUrl"`
	IsAdmin      bool   `json:"isAdmin"`
	IsDba        bool   `json:"isDba"`
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
	IsDba     bool
}

type UpdateUserReqDTO struct {
	Account   string
	Name      string
	Email     string
	AvatarUrl string
}

type UpdateAdminReqDTO struct {
	Account string
	IsAdmin bool
}

type UpdateDbaReqDTO struct {
	Account string
	IsDba   bool
}

type UpdatePasswordReqDTO struct {
	Account  string
	Password string
}

type PageUserReqDTO struct {
	Account  string
	PageNum  int
	PageSize int
}

type SetUserProhibitedReqDTO struct {
	Account      string
	IsProhibited bool
}
