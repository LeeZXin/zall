import request from '@/utils/request.js';

// webhook列表
const listWebhookRequest = (repoId) => request.get("/api/webhook/list/" + repoId);
// 创建webhook
const createWebhookRequest = (data) => request.post("/api/webhook/create", data);
// 编辑webhook
const updateWebhookRequest = (data) => request.post("/api/webhook/update", data);
// 删除webhook
const deleteWebhookRequest = (id) => request.delete("/api/webhook/delete/" + id);
// ping webhook
const pingWebhookRequest = (id) => request.put("/api/webhook/ping/" + id);
export {
    listWebhookRequest,
    createWebhookRequest,
    updateWebhookRequest,
    deleteWebhookRequest,
    pingWebhookRequest
}