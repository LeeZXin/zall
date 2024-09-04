import request from '@/utils/request.js';

// webhook列表
const listWebhookRequest = (repoId) => request.get("/api/gitWebhook/list/" + repoId);
// 创建webhook
const createWebhookRequest = (data) => request.post("/api/gitWebhook/create", data);
// 编辑webhook
const updateWebhookRequest = (data) => request.post("/api/gitWebhook/update", data);
// 删除webhook
const deleteWebhookRequest = (id) => request.delete("/api/gitWebhook/delete/" + id);
// ping webhook
const pingWebhookRequest = (id) => request.put("/api/gitWebhook/ping/" + id);
export {
    listWebhookRequest,
    createWebhookRequest,
    updateWebhookRequest,
    deleteWebhookRequest,
    pingWebhookRequest
}