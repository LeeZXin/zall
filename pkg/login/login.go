package login

type Cfg struct {
	AccountPassword AccountPassword `json:"accountPassword"`
	Wework          Wework          `json:"wework"`
	Feishu          Feishu          `json:"feishu"`
}

func (c *Cfg) IsValid() bool {
	if !c.Wework.IsEnabled && !c.AccountPassword.IsEnabled && !c.Feishu.IsEnabled {
		return false
	}
	if c.Wework.IsEnabled {
		if !c.Wework.IsValid() {
			return false
		}
	}
	if c.AccountPassword.IsEnabled {
		// nothing
	}
	if c.Feishu.IsValid() {
		if !c.Feishu.IsValid() {
			return false
		}
	}
	return true
}
