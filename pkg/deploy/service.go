package deploy

import (
	"encoding/json"
)

type Action struct {
	Script string `json:"script" yaml:"script"`
}

func (s *Action) isValid() bool {
	return len(s.Script) > 0
}

type ServiceType string

const (
	ProcessServiceType ServiceType = "process"
	K8sServiceType     ServiceType = "k8s"
)

type Service struct {
	Type    ServiceType       `json:"type" json:"type"`
	Actions map[string]Action `json:"actions,omitempty" yaml:"actions,omitempty"`
	Agents  map[string]Agent  `json:"agents,omitempty" yaml:"agents,omitempty"`
	Deploy  []Stage           `json:"deploy,omitempty" yaml:"deploy,omitempty"`
}

func (s *Service) IsValid() bool {
	if len(s.Actions) == 0 || len(s.Agents) == 0 {
		return false
	}
	for _, action := range s.Actions {
		if !action.isValid() {
			return false
		}
	}
	for _, agent := range s.Agents {
		if !agent.isValid() {
			return false
		}
	}
	switch s.Type {
	case ProcessServiceType:

	case K8sServiceType:

	default:
		return false
	}
	if len(s.Deploy) == 0 {
		return false
	}
	for _, stage := range s.Deploy {
		if !stage.isValid(s.Agents, s.Actions) {
			return false
		}
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
