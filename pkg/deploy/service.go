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

type Pipeline struct {
	Actions map[string]Action `json:"actions,omitempty" yaml:"actions,omitempty"`
	Agents  map[string]Agent  `json:"agents,omitempty" yaml:"agents,omitempty"`
	Deploy  []Stage           `json:"deploy,omitempty" yaml:"deploy,omitempty"`
}

func (s *Pipeline) IsValid() bool {
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

func (s *Pipeline) FromDB(content []byte) error {
	if s == nil {
		*s = Pipeline{}
	}
	return json.Unmarshal(content, s)
}

func (s *Pipeline) ToDB() ([]byte, error) {
	return json.Marshal(s)
}
