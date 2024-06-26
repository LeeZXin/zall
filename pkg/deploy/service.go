package deploy

import (
	"encoding/json"
)

type Action struct {
	Alias  string `json:"alias" yaml:"alias"`
	Show   bool   `json:"show" yaml:"show"`
	Script string `json:"script" yaml:"script"`
}

func (s *Action) isValid() bool {
	return len(s.Script) > 0
}

type ProbeFail struct {
	Times  int    `json:"times" yaml:"times"`
	Action string `json:"action" yaml:"action"`
}

func (f *ProbeFail) isValid(actions map[string]Action) bool {
	if f.Times <= 0 {
		return false
	}
	_, b := actions[f.Action]
	return b
}

type ServiceType string

const (
	ProcessServiceType ServiceType = "process"
	K8sServiceType     ServiceType = "k8s"
)

type Service struct {
	Type    ServiceType       `json:"type" json:"type"`
	K8s     *K8s              `json:"k8s,omitempty" yaml:"k8s,omitempty"`
	Process []Process         `json:"process,omitempty" yaml:"process,omitempty"`
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
		if len(s.Process) == 0 {
			return false
		}
		for _, p := range s.Process {
			if !p.isValid(s.Agents) {
				return false
			}
		}
	case K8sServiceType:
		if s.K8s == nil || !s.K8s.isValid(s.Agents) {
			return false
		}
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

type K8sService struct {
	GetStatusScript string `json:"getStatusScript" yaml:"getStatusScript"`
}

func (s *K8sService) isValid() bool {
	return len(s.GetStatusScript) > 0
}

type K8s struct {
	Agent           string `json:"agent" yaml:"agent"`
	GetStatusScript string `json:"getStatusScript" yaml:"getStatusScript"`
}

func (k *K8s) isValid(agents map[string]Agent) bool {
	_, b := agents[k.Agent]
	if !b {
		return false
	}
	return k.GetStatusScript != ""
}

type Process struct {
	Name  string `json:"name" yaml:"name"`
	Agent string `json:"agent" yaml:"agent"`
}

func (c *Process) isValid(agents map[string]Agent) bool {
	// 不存在agent
	_, b := agents[c.Agent]
	if !b {
		return false
	}
	return true
}
