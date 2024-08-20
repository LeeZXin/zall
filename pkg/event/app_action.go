package event

type AppEventAction string

const (
	AppCreateAction   AppEventAction = "create"
	AppDeleteAction   AppEventAction = "delete"
	AppUpdateAction   AppEventAction = "update"
	AppTransferAction AppEventAction = "transfer"
)

type AppSourceEventAction string

const (
	AppManagePropertySourceAction  AppSourceEventAction = "managePropertySource"
	AppManageServiceSourceAction   AppSourceEventAction = "manageServiceSource"
	AppManageDiscoverySourceAction AppSourceEventAction = "manageDiscoverySource"
)

type AppPropertyFileEventAction string

const (
	AppPropertyFileCreateFileAction AppPropertyFileEventAction = "create"
	AppPropertyFileDeleteFileAction AppPropertyFileEventAction = "delete"
)

type AppPropertyVersionEventAction string

const (
	AppPropertyVersionCreateAction AppPropertyVersionEventAction = "create"
	AppPropertyVersionDeployAction AppPropertyVersionEventAction = "deploy"
)

type AppPipelineEventAction string

const (
	AppPipelineCreatePipelineAction AppPipelineEventAction = "create"
	AppPipelineUpdatePipelineAction AppPipelineEventAction = "update"
	AppPipelineDeletePipelineAction AppPipelineEventAction = "delete"
)

type AppPipelineVarsEventAction string

const (
	AppPipelineCreateVarsAction AppPipelineVarsEventAction = "create"
	AppPipelineUpdateVarsAction AppPipelineVarsEventAction = "update"
	AppPipelineDeleteVarsAction AppPipelineVarsEventAction = "delete"
)

type AppDeployPlanEventAction string

const (
	AppDeployPlanCreateAction   AppDeployPlanEventAction = "create"
	AppDeployPlanCloseAction    AppDeployPlanEventAction = "close"
	AppDeployPlanCompleteAction AppDeployPlanEventAction = "complete"
)

type AppProductEventAction string

const (
	AppProductCreateAction AppProductEventAction = "create"
	AppProductDeleteAction AppProductEventAction = "delete"
)
