package event

type AppEventCfg struct {
	Create   bool `json:"create"`
	Delete   bool `json:"delete"`
	Update   bool `json:"update"`
	Transfer bool `json:"transfer"`
}

type AppSourceEventCfg struct {
	ManagePropertySource  bool `json:"managePropertySource"`
	ManageDiscoverySource bool `json:"manageDiscoverySource"`
	ManageServiceSource   bool `json:"manageServiceSource"`
}

type AppPropertyFileEventCfg struct {
	Create bool `json:"create"`
	Delete bool `json:"delete"`
}

type AppPropertyVersionEventCfg struct {
	Create bool `json:"create"`
	Deploy bool `json:"deploy"`
}

type AppPipelineEventCfg struct {
	Create bool `json:"create"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type AppPipelineVarsEventCfg struct {
	Create bool `json:"create"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type AppDeployPlanEventCfg struct {
	Create   bool `json:"create"`
	Close    bool `json:"close"`
	Complete bool `json:"complete"`
}

type AppProductEventCfg struct {
	Create bool `json:"create"`
	Delete bool `json:"delete"`
}
