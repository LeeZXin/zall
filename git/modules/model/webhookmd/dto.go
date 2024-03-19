package webhookmd

type InsertWebhookReqDTO struct {
	RepoId      int64
	HookUrl     string
	HttpHeaders map[string]string
	HookType    HookType
	WildBranch  string
	WildTag     string
}

type UpdateWebhookReqDTO struct {
	Id          int64
	HookUrl     string
	HttpHeaders map[string]string
	WildBranch  string
	WildTag     string
}
