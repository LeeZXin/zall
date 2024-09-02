package login

type Wework struct {
	IsEnabled   bool   `json:"isEnabled"`
	State       string `json:"state"`
	AppId       string `json:"appId"`
	AgentId     string `json:"agentId"`
	Secret      string `json:"secret,omitempty"`
	RedirectUrl string `json:"redirectUrl"`
	Lang        string `json:"lang"`
}

func (w *Wework) IsValid() bool {
	if w.State == "" || w.AgentId == "" || w.AppId == "" || w.RedirectUrl == "" || w.Secret == "" {
		return false
	}
	return true
}

func (w *Wework) EraseSensitiveVar() {
	w.Secret = ""
}
