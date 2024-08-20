import request from '@/utils/request.js';

// team hook列表
const listTeamHookRequest = (teamId) => request.get("/api/teamHook/list/" + teamId);
// 创建team hook
const createTeamHookRequest = (data) => request.post("/api/teamHook/create", data);
// 编辑team hook
const updateTeamHookRequest = (data) => request.post("/api/teamHook/update", data);
// 删除team hook
const deleteTeamHookRequest = (hookId) => request.delete("/api/teamHook/delete/" + hookId);

export {
    listTeamHookRequest,
    createTeamHookRequest,
    updateTeamHookRequest,
    deleteTeamHookRequest
}