package deploy

import (
	"encoding/json"
	"github.com/LeeZXin/zall/util"
)

type Action struct {
	Name   string `json:"name" yaml:"name"`
	Script string `json:"script" yaml:"script"`
}

func (s *Action) IsValid() bool {
	return len(s.Name) > 0 && len(s.Script) > 0
}

type Script struct {
	AgentHost  string   `json:"agentHost" yaml:"agentHost"`
	AgentToken string   `json:"agentToken" yaml:"agentToken"`
	Actions    []Action `json:"actions" yaml:"actions"`
}

func (s *Script) IsValid() bool {
	if !util.IpPortPattern.MatchString(s.AgentHost) {
		return false
	}
	for _, action := range s.Actions {
		if !action.IsValid() {
			return false
		}
	}
	return true
}

type ServiceType string

const (
	ProcessServiceType ServiceType = "process"
	K8sServiceType     ServiceType = "k8s"
)

type Service struct {
	Type   ServiceType       `json:"type" json:"type"`
	K8s    *K8sServiceConfig `json:"k8s,omitempty" yaml:"k8s,omitempty"`
	Probe  *Probe            `json:"probe,omitempty" yaml:"probe,omitempty"`
	Script *Script           `json:"script,omitempty" yaml:"script,omitempty"`
}

func (s *Service) IsValid() bool {
	switch s.Type {
	case ProcessServiceType:
	case K8sServiceType:
		if s.K8s == nil || !s.K8s.IsValid() {
			return false
		}
	default:
		return false
	}
	if s.Probe != nil && !s.Probe.IsValid() {
		return false
	}
	if s.Script != nil && !s.Script.IsValid() {
		return false
	}
	return true
}

func (s *Service) FromDB(content []byte) error {
	if s == nil {
		*s = Service{}
	}
	return json.Unmarshal(content, s)
}

func (s *Service) ToDB() ([]byte, error) {
	return json.Marshal(s)
}

type K8sServiceConfig struct {
}

func (*K8sServiceConfig) IsValid() bool {
	return true
}
