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
	New    bool `json:"new"`
	Deploy bool `json:"deploy"`
}

type AppDeployPipelineEventCfg struct {
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
	Create bool `json:"create"`
	Close  bool `json:"close"`
	Start  bool `json:"start"`
}

type AppProductEventCfg struct {
	Delete bool `json:"delete"`
}

type AppDiscoveryEventCfg struct {
	Deregister        bool `json:"deregister"`
	ReRegister        bool `json:"reRegister"`
	DeleteDownService bool `json:"deleteDownService"`
}

type AppDeployServiceEventCfg struct {
	TriggerAction bool `json:"triggerAction"`
}
