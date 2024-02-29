package webhookmd

type InsertWebhookReqDTO struct {
	RepoId      int64
	HookUrl     string
	HttpHeaders map[string]string
	HookType    HookType
}

type UpdateWebhookReqDTO struct {
	Id          int64
	HookUrl     string
	HttpHeaders map[string]string
}

type WebhookDTO struct {
	Id          int64
	RepoId      int64
	HookUrl     string
	HttpHeaders map[string]string
	HookType    HookType
}
