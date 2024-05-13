import request from '@/utils/request.js';

// 创建工作流
const createWorkflowRequest = (data) => request.post("/api/workflow/create", data);
// 工作流列表
const listWorkflowRequest = (repoId) => request.get("/api/workflow/list/" + repoId);
// 删除工作流
const DeleteWorkflowRequest = (workflowId) => request.delete("/api/workflow/delete/" + workflowId);
// 工作流任务列表
const listTaskRequest = (workflowId, params) => request.get("/api/workflow/list/" + workflowId, { params });
export {
    createWorkflowRequest,
    listWorkflowRequest,
    DeleteWorkflowRequest,
    listTaskRequest
}