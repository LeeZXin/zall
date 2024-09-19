package event

type AppEventAction string

func (a AppEventAction) GetI18nValue() string {
	return "appEvent." + string(a)
}

const (
	AppCreateAction   AppEventAction = "create"
	AppDeleteAction   AppEventAction = "delete"
	AppUpdateAction   AppEventAction = "update"
	AppTransferAction AppEventAction = "transfer"
)

type AppSourceEventAction string

func (a AppSourceEventAction) GetI18nValue() string {
	return "appSourceEvent." + string(a)
}

const (
	AppManagePropertySourceAction  AppSourceEventAction = "managePropertySource"
	AppManageServiceSourceAction   AppSourceEventAction = "manageServiceSource"
	AppManageDiscoverySourceAction AppSourceEventAction = "manageDiscoverySource"
)

type AppPropertyFileEventAction string

func (a AppPropertyFileEventAction) GetI18nValue() string {
	return "appPropertyFileEvent." + string(a)
}

const (
	AppPropertyFileCreateFileAction AppPropertyFileEventAction = "create"
	AppPropertyFileDeleteFileAction AppPropertyFileEventAction = "delete"
)

type AppPropertyVersionEventAction string

func (a AppPropertyVersionEventAction) GetI18nValue() string {
	return "appPropertyVersionEvent." + string(a)
}

const (
	AppPropertyVersionNewAction    AppPropertyVersionEventAction = "new"
	AppPropertyVersionDeployAction AppPropertyVersionEventAction = "deploy"
)

type AppDeployPipelineEventAction string

func (a AppDeployPipelineEventAction) GetI18nValue() string {
	return "appDeployPipelineEvent." + string(a)
}

const (
	AppDeployPipelineCreatePipelineAction AppDeployPipelineEventAction = "create"
	AppDeployPipelineUpdatePipelineAction AppDeployPipelineEventAction = "update"
	AppDeployPipelineDeletePipelineAction AppDeployPipelineEventAction = "delete"
)

type AppDeployPipelineVarsEventAction string

func (a AppDeployPipelineVarsEventAction) GetI18nValue() string {
	return "appDeployPipelineVarsEvent." + string(a)
}

const (
	AppDeployPipelineVarsCreateAction AppDeployPipelineVarsEventAction = "create"
	AppDeployPipelineVarsUpdateAction AppDeployPipelineVarsEventAction = "update"
	AppDeployPipelineVarsDeleteAction AppDeployPipelineVarsEventAction = "delete"
)

type AppDeployPlanEventAction string

func (a AppDeployPlanEventAction) GetI18nValue() string {
	return "appDeployPlanEvent." + string(a)
}

const (
	AppDeployPlanCreateAction AppDeployPlanEventAction = "create"
	AppDeployPlanStartAction  AppDeployPlanEventAction = "start"
	AppDeployPlanCloseAction  AppDeployPlanEventAction = "close"
)

type AppArtifactEventAction string

func (a AppArtifactEventAction) GetI18nValue() string {
	return "appArtifactEvent." + string(a)
}

const (
	AppArtifactUploadAction AppArtifactEventAction = "upload"
	AppArtifactDeleteAction AppArtifactEventAction = "delete"
)

type AppDiscoveryEventAction string

func (a AppDiscoveryEventAction) GetI18nValue() string {
	return "appDiscoveryEvent." + string(a)
}

const (
	AppDiscoveryMarkAsDownAction AppDiscoveryEventAction = "markAsDown"
	AppDiscoveryMarkAsUpAction   AppDiscoveryEventAction = "markAsUp"
)

type AppDeployServiceEventAction string

func (a AppDeployServiceEventAction) GetI18nValue() string {
	return "appDeployServiceEvent." + string(a)
}

const (
	AppDeployServiceKillAction    AppDeployServiceEventAction = "kill"
	AppDeployServiceRestartAction AppDeployServiceEventAction = "restart"
)

type AppPromScrapeEventAction string

func (a AppPromScrapeEventAction) GetI18nValue() string {
	return "appPromScrapeEvent." + string(a)
}

const (
	AppPromScrapeCreateAction AppPromScrapeEventAction = "create"
	AppPromScrapeUpdateAction AppPromScrapeEventAction = "update"
	AppPromScrapeDeleteAction AppPromScrapeEventAction = "delete"
)

type AppAlertConfigEventAction string

func (a AppAlertConfigEventAction) GetI18nValue() string {
	return "appAlertConfigEvent." + string(a)
}

const (
	AppAlertConfigCreateAction  AppAlertConfigEventAction = "create"
	AppAlertConfigUpdateAction  AppAlertConfigEventAction = "update"
	AppAlertConfigDeleteAction  AppAlertConfigEventAction = "delete"
	AppAlertConfigEnableAction  AppAlertConfigEventAction = "enable"
	AppAlertConfigDisableAction AppAlertConfigEventAction = "disable"
)
