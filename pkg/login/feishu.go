package login

type Feishu struct {
	IsEnabled   bool   `json:"isEnabled"`
	ClientId    string `json:"clientId"`
	RedirectUrl string `json:"redirectUrl"`
	State       string `json:"state"`
	Secret      string `json:"secret,omitempty"`
}

func (f *Feishu) IsValid() bool {
	return f.ClientId != "" && f.RedirectUrl != "" && f.State != ""
}

func (f *Feishu) EraseSensitiveVar() {
	f.Secret = ""
}
