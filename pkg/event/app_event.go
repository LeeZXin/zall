package event

type BaseApp struct {
	AppId   string `json:"appId"`
	AppName string `json:"appName"`
}

type AppEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	Action       AppEventAction `json:"action"`
	TransferTeam *BaseTeam      `json:"transferTeam,omitempty"`
}

func (*AppEvent) EventType() string {
	return "app-event"
}

type AppSource struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Env  string `json:"env"`
}

type AppSourceEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	Env     string               `json:"env"`
	Action  AppSourceEventAction `json:"action"`
	Sources []AppSource          `json:"sources"`
}

func (*AppSourceEvent) EventType() string {
	return "app-source-event"
}

type BasePropertyFile struct {
	FileId   int64  `json:"fileId"`
	FileName string `json:"fileName"`
	Env      string `json:"env"`
}

type AppPropertyFileEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	BasePropertyFile
	Action AppPropertyFileEventAction `json:"action"`
}

func (*AppPropertyFileEvent) EventType() string {
	return "app-property-file-event"
}

type AppPropertyVersionEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	BasePropertyFile
	OldVersion string                        `json:"oldVersion"`
	Version    string                        `json:"version"`
	Action     AppPropertyVersionEventAction `json:"action"`
}

func (*AppPropertyVersionEvent) EventType() string {
	return "app-property-version-event"
}

type BasePipeline struct {
	PipelineId   int64  `json:"pipelineId"`
	PipelineName string `json:"pipelineName"`
	Env          string `json:"env"`
}

type AppDeployPipelineEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	BasePipeline
	Action AppDeployPipelineEventAction `json:"action"`
}

func (*AppDeployPipelineEvent) EventType() string {
	return "app-deploy-pipeline-event"
}

type AppDeployPipelineVarsEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	VarsId   int64                            `json:"varsId"`
	VarsName string                           `json:"varsName"`
	Env      string                           `json:"env"`
	Action   AppDeployPipelineVarsEventAction `json:"action"`
}

func (*AppDeployPipelineVarsEvent) EventType() string {
	return "app-deploy-pipeline-vars-event"
}

type AppDeployPlanEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	BasePipeline
	Action   AppDeployPlanEventAction `json:"action"`
	PlanId   int64                    `json:"planId"`
	PlanName string                   `json:"planName"`
	Env      string                   `json:"env"`
}

func (*AppDeployPlanEvent) EventType() string {
	return "app-deploy-plan-event"
}

type AppArtifactEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	ArtifactId   int64                  `json:"productId"`
	ArtifactName string                 `json:"productName"`
	Env          string                 `json:"env"`
	Action       AppArtifactEventAction `json:"action"`
}

func (*AppArtifactEvent) EventType() string {
	return "app-product-event"
}

type AppDiscoveryEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	Source AppSource               `json:"source"`
	Action AppDiscoveryEventAction `json:"action"`
}

func (*AppDiscoveryEvent) EventType() string {
	return "app-discovery-event"
}

type AppDeployServiceEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	Source        AppSource                   `json:"source"`
	Action        AppDeployServiceEventAction `json:"action"`
	TriggerAction string                      `json:"triggerAction"`
}

func (*AppDeployServiceEvent) EventType() string {
	return "app-deploy-service-event"
}

type AppPromScrapeEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	Action   AppPromScrapeEventAction `json:"action"`
	Env      string                   `json:"env"`
	Endpoint string                   `json:"endpoint"`
}

func (*AppPromScrapeEvent) EventType() string {
	return "app-prom-scrape-event"
}

type AppAlertConfigEvent struct {
	BaseTeam
	BaseApp
	BaseEvent
	Action AppAlertConfigEventAction `json:"action"`
	Name   string                    `json:"name"`
	Env    string                    `json:"env"`
}

func (*AppAlertConfigEvent) EventType() string {
	return "app-alert-config-event"
}
