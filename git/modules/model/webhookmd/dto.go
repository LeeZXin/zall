package webhookmd

type InsertWebhookReqDTO struct {
	RepoId      int64
	HookUrl     string
	HttpHeaders HttpHeaders
	HookType    HookType
	WildBranch  string
	WildTag     string
}

type UpdateWebhookReqDTO struct {
	Id          int64
	HookUrl     string
	HttpHeaders HttpHeaders
	WildBranch  string
	WildTag     string
}
