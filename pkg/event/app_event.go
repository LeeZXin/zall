package event

type BaseApp struct {
	AppId   string `json:"appId"`
	AppName string `json:"appName"`
}

type AppEvent struct {
	BaseApp
	BaseEvent
	Action  AppEventAction `json:"action"`
	Sources []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"sources"`
}

func (*AppEvent) EventType() string {
	return "app-event"
}

type BasePropertyFile struct {
	FileId   int64  `json:"fileId"`
	FileName string `json:"fileName"`
	Env      string `json:"env"`
}

type AppPropertyFileEvent struct {
	BaseApp
	BaseEvent
	BasePropertyFile
	Action AppPropertyFileEventAction `json:"action"`
}

func (*AppPropertyFileEvent) EventType() string {
	return "app-property-file-event"
}

type AppPropertyVersionEvent struct {
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
}

type AppPipelineEvent struct {
	BaseApp
	BaseEvent
	BasePipeline
	Action AppPipelineEventAction `json:"action"`
}

func (*AppPipelineEvent) EventType() string {
	return "app-pipeline-event"
}

type AppPipelineVarsEvent struct {
	BaseApp
	BaseEvent
	BasePipeline
	VarsId   int64                      `json:"varsId"`
	VarsName string                     `json:"varsName"`
	Action   AppPipelineVarsEventAction `json:"action"`
}

func (*AppPipelineVarsEvent) EventType() string {
	return "app-pipeline-vars-event"
}

type AppDeployPlanEvent struct {
	BaseApp
	BaseEvent
	Action AppDeployPlanEventAction `json:"action"`
	Status string                   `json:"status"`
}

func (*AppDeployPlanEvent) EventType() string {
	return "app-deploy-plan-event"
}

type AppProductEvent struct {
	BaseApp
	BaseEvent
	ProductPath string                `json:"productPath"`
	ProductName string                `json:"productName"`
	Action      AppProductEventAction `json:"action"`
}

func (*AppProductEvent) EventType() string {
	return "app-product-event"
}
