import request from '@/utils/request.js';

// 创建工作流
const createWorkflowRequest = (data) => request.post("/api/workflow/create", data);
// 编辑工作流
const updateWorkflowRequest = (data) => request.post("/api/workflow/update", data);
// 工作流列表
const listWorkflowRequest = (repoId) => request.get("/api/workflow/list/" + repoId);
// 删除工作流
const deleteWorkflowRequest = (workflowId) => request.delete("/api/workflow/delete/" + workflowId);
// 工作流任务列表
const listTaskRequest = (workflowId, params) => request.get("/api/workflowTask/list/" + workflowId, { params });
// 获取工作流详情
const getWorkflowDetailRequest = (workflowId) => request.get("/api/workflow/detail/" + workflowId);
// 触发工作流
const triggerWorkflowRequest = (workflowId, branch) => request.put(`/api/workflow/trigger/${workflowId}?branch=${branch}`);
// 停止工作流
const killWorkflowTaskRequest = (taskId) => request.put(`/api/workflowTask/kill/${taskId}`);
// 获取任务详情
const getTaskDetailRequest = (taskId) => request.get("/api/workflowTask/detail/" + taskId);
export {
    createWorkflowRequest,
    listWorkflowRequest,
    deleteWorkflowRequest,
    listTaskRequest,
    getWorkflowDetailRequest,
    updateWorkflowRequest,
    triggerWorkflowRequest,
    killWorkflowTaskRequest,
    getTaskDetailRequest
}