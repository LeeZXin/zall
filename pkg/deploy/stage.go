package deploy

import (
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"regexp"
)

var (
	IpPortPattern  = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}:\d+$`)
	NoSpacePattern = regexp.MustCompile(`^\S+$`)
)

type Option struct {
	Value string `json:"value" yaml:"value"`
	Label string `json:"label" yaml:"label"`
}

type FormItem struct {
	Key     string   `json:"key" yaml:"key"`
	Label   string   `json:"label" yaml:"label"`
	Regexp  string   `json:"regexp" yaml:"regexp"`
	Options []Option `json:"options,omitempty" yaml:"options,omitempty"`
}

func (i *FormItem) IsValid() bool {
	if !NoSpacePattern.MatchString(i.Key) {
		return false
	}
	// 没有下拉框就是手填数据 需要校验正则
	if len(i.Options) == 0 {
		_, err := regexp.Compile(i.Regexp)
		if err != nil {
			return false
		}
	}
	return true
}

type Confirm struct {
	// 是否需要交互
	NeedInteract bool       `json:"needInteract" yaml:"needInteract"`
	Message      string     `json:"message" yaml:"message"`
	Form         []FormItem `json:"form,omitempty" yaml:"form,omitempty"`
}

func (c *Confirm) CheckForm(args map[string]string) (bool, map[string]string) {
	if len(c.Form) == 0 {
		return true, nil
	}
	if args == nil {
		return false, nil
	}
	filteredArgs := make(map[string]string)
	if len(c.Form) == 0 {
		return true, filteredArgs
	}
	for _, item := range c.Form {
		// 下拉框
		if len(item.Options) > 0 {
			find := false
			for _, option := range item.Options {
				if args[item.Key] == option.Value {
					find = true
					filteredArgs[item.Key] = args[item.Key]
				}
			}
			if !find {
				return false, filteredArgs
			}
		} else {
			// 匹配输入框正则
			re, e := regexp.Compile(item.Regexp)
			if e == nil {
				if !re.MatchString(args[item.Key]) {
					return false, filteredArgs
				}
			}
			filteredArgs[item.Key] = args[item.Key]
		}
	}
	return true, filteredArgs
}

func (c *Confirm) isValid() bool {
	if c.NeedInteract {
		for _, item := range c.Form {
			if !item.IsValid() {
				return false
			}
		}
	}
	return true
}

type Stage struct {
	Name     string   `json:"name" yaml:"name"`
	Agents   []string `json:"agents,omitempty" yaml:"agents,omitempty"`
	Confirm  *Confirm `json:"confirm" yaml:"confirm,omitempty"`
	Action   string   `json:"action" yaml:"action"`
	Parallel int      `json:"parallel" yaml:"parallel"`
}

func (s *Stage) isValid(agents map[string]Agent, actions map[string]Action) bool {
	if len(s.Agents) != hashset.NewHashSet(s.Agents...).Size() {
		return false
	}
	for _, agent := range s.Agents {
		_, b := agents[agent]
		if !b {
			return false
		}
	}
	_, b := actions[s.Action]
	if !b {
		return false
	}
	return s.Confirm == nil || s.Confirm.isValid()
}

type Agent struct {
	NodeId string            `json:"nodeId" yaml:"nodeId"`
	With   map[string]string `json:"with,omitempty" yaml:"with,omitempty"`
}

func (a *Agent) isValid() bool {
	return a.NodeId != ""
}
