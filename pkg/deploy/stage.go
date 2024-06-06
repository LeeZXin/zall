package deploy

import (
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"gopkg.in/yaml.v3"
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

type Confirm struct {
	Message string     `json:"message" yaml:"message"`
	Agents  []string   `json:"agents,omitempty" yaml:"agents,omitempty"`
	Script  string     `json:"script" yaml:"script"`
	Form    []FormItem `json:"form,omitempty" yaml:"form,omitempty"`
}

func (c *Confirm) IsValid(agentIdSet hashset.Set[string]) bool {
	for _, agent := range c.Agents {
		if !agentIdSet.Contains(agent) {
			return false
		}
	}
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
	return true
}

type Stage struct {
	Name    string  `json:"name" yaml:"name"`
	Confirm Confirm `json:"confirm" yaml:"confirm"`
	// 回滚脚本
	Rollback string `json:"rollback" yaml:"rollback"`
}

func (s *Stage) IsValid(agentIdSet hashset.Set[string]) bool {
	return s.Confirm.IsValid(agentIdSet)
}

type TcpProbe struct {
	Addr string `json:"addr" yaml:"addr"`
}

type HttpProbe struct {
	Url string `json:"url" yaml:"url"`
}

type Probe struct {
	Type string     `json:"type" yaml:"type"`
	Tcp  *TcpProbe  `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	Http *HttpProbe `json:"http,omitempty" yaml:"http,omitempty"`
}

type Agent struct {
	Id    string            `json:"id" yaml:"id"`
	Host  string            `json:"host" yaml:"host"`
	Token string            `json:"token" yaml:"token"`
	Probe *Probe            `json:"probe,omitempty" yaml:"probe,omitempty"`
	With  map[string]string `json:"with,omitempty" yaml:"with,omitempty"`
}

func (a *Agent) IsValid() bool {
	return a.Id != "" && IpPortPattern.MatchString(a.Host)
}

type Process struct {
	Name    string  `json:"name" yaml:"name"`
	Agents  []Agent `json:"agents,omitempty" yaml:"agents,omitempty"`
	Deploy  []Stage `json:"deploy,omitempty" yaml:"deploy,omitempty"`
	Actions []Stage `json:"actions,omitempty" yaml:"actions,omitempty"`
}

func (p *Process) IsValid() bool {
	if len(p.Agents) == 0 || len(p.Deploy) == 0 {
		return false
	}
	idSet := hashset.NewHashSet[string]()
	for _, agent := range p.Agents {
		if !agent.IsValid() {
			return false
		}
		// id重复
		if idSet.Contains(agent.Id) {
			return false
		}
		idSet.Add(agent.Id)
	}
	for _, stage := range p.Deploy {
		if !stage.IsValid(idSet) {
			return false
		}
	}
	for _, action := range p.Actions {
		if !action.IsValid(idSet) {
			return false
		}
	}
	return true
}

type Deploy struct {
	Processes []Process `json:"processes,omitempty" yaml:"processes,omitempty"`
}

func (c *Deploy) IsValid() bool {
	if len(c.Processes) == 0 {
		return false
	}
	for _, process := range c.Processes {
		if !process.IsValid() {
			return false
		}
	}
	return true
}

func (c *Deploy) FromDB(content []byte) error {
	if c == nil {
		*c = Deploy{}
	}
	return yaml.Unmarshal(content, c)
}

func (c *Deploy) ToDB() ([]byte, error) {
	return yaml.Marshal(c)
}
