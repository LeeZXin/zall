package deploy

import (
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"regexp"
	"strings"
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

type Confirm struct {
	// 是否需要交互
	NeedInteract bool       `json:"needInteract" yaml:"needInteract"`
	Message      string     `json:"message" yaml:"message"`
	Action       string     `json:"action" yaml:"action"`
	Form         []FormItem `json:"form,omitempty" yaml:"form,omitempty"`
}

func (c *Confirm) isValid(actions map[string]Action) bool {
	_, b := actions[c.Action]
	if !b {
		return false
	}
	if c.NeedInteract {
		for _, item := range c.Form {
			if !NoSpacePattern.MatchString(item.Key) {
				return false
			}
			// 没有下拉框就是手填数据 需要校验正则
			if len(item.Options) == 0 {
				_, err := regexp.Compile(item.Regexp)
				if err != nil {
					return false
				}
			}
		}
	}
	return true
}

type Rollback struct {
	Action string `json:"action" yaml:"action"`
}

func (r *Rollback) isValid(actions map[string]Action) bool {
	_, b := actions[r.Action]
	return b
}

type Stage struct {
	Name     string    `json:"name" yaml:"name"`
	Agents   []string  `json:"agents,omitempty" yaml:"agents,omitempty"`
	Confirm  Confirm   `json:"confirm" yaml:"confirm"`
	Rollback *Rollback `json:"rollback,omitempty" yaml:"rollback,omitempty"`
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
	if s.Rollback != nil && !s.Rollback.isValid(actions) {
		return false
	}
	return s.Confirm.isValid(actions)
}

type Agent struct {
	Host  string            `json:"host" yaml:"host"`
	Token string            `json:"token" yaml:"token"`
	With  map[string]string `json:"with,omitempty" yaml:"with,omitempty"`
}

func (a *Agent) isValid() bool {
	return IpPortPattern.MatchString(a.Host)
}

func (a *Agent) RunScript(script, service string, env map[string]string) (string, error) {
	if script == "" {
		return "", nil
	}
	args := make(map[string]string, len(a.With)+len(env))
	for k, v := range a.With {
		args[k] = v
	}
	for k, v := range env {
		args[k] = v
	}
	return sshagent.NewServiceCommand(a.Host, a.Token, service).
		Execute(strings.NewReader(script), args)
}
