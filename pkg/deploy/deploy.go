package deploy

import "github.com/LeeZXin/zsf-utils/collections/hashset"

type Deploy struct {
	Agents []Agent `json:"agents,omitempty" yaml:"agents,omitempty"`
	Deploy []Stage `json:"deploy,omitempty" yaml:"deploy,omitempty"`
}

func (p *Deploy) IsValid() bool {
	if len(p.Agents) == 0 || len(p.Deploy) == 0 {
		return false
	}
	idSet := hashset.NewHashSet[string]()
	for _, agent := range p.Agents {
		if !agent.isValid() {
			return false
		}
		// id重复
		if idSet.Contains(agent.Id) {
			return false
		}
		idSet.Add(agent.Id)
	}
	for _, stage := range p.Deploy {
		if !stage.isValid(idSet) {
			return false
		}
	}
	return true
}

func (p *Deploy) GetAgentMap() map[string]Agent {
	ret := make(map[string]Agent, len(p.Agents))
	for _, agent := range p.Agents {
		ret[agent.Id] = agent
	}
	return ret
}
