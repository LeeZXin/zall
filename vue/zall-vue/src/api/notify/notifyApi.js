import request from '@/utils/request.js';

// 通知模板列表
const listNotifyTplRequest = (data) => request.get("/api/notifyTpl/list", { params: data });
// 创建通知模板
const createNotifyTplRequest = (data) => request.post("/api/notifyTpl/create", data);
// 编辑通知模板
const updateNotifyTplRequest = (data) => request.post("/api/notifyTpl/update", data);
// 删除通知模板
const deleteNotifyTplRequest = (tplId) => request.delete("/api/notifyTpl/delete/" + tplId);
// 变更api key
const changeNotifyTplApiKeyRequest = (tplId) => request.put("/api/notifyTpl/changeApiKey/" + tplId);
// 所有通知模板
const listAllTplByTeamIdRequest = (teamId) => request.get("/api/notifyTpl/listAll/" + teamId);
export {
    listNotifyTplRequest,
    createNotifyTplRequest,
    updateNotifyTplRequest,
    deleteNotifyTplRequest,
    changeNotifyTplApiKeyRequest,
    listAllTplByTeamIdRequest
}