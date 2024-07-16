import request from '@/utils/request.js';

// 定时任务列表
const listTimerTaskRequest = (data) => request.get("/api/timerTask/list", { params: data });
// 创建定时任务
const createTimerTaskRequest = (data) => request.post("/api/timerTask/create", data);
// 启动定时任务
const enableTimerTaskRequest = (taskId) => request.put("/api/timerTask/enable/" + taskId);
// 停止定时任务
const disableTimerTaskRequest = (taskId) => request.put("/api/timerTask/disable/" + taskId);
// 删除定时任务
const deleteTimerTaskRequest = (taskId) => request.delete("/api/timerTask/delete/" + taskId);
// 触发定时任务
const triggerTimerTaskRequest = (taskId) => request.put("/api/timerTask/trigger/" + taskId);
// 编辑定时任务
const updateTimerTaskRequest = (data) => request.post("/api/timerTask/update", data);
// 查看定时任务日志
const pageTimerTaskLogRequest = (data) => request.get("/api/timerLog/list", { params: data });
export {
    listTimerTaskRequest,
    createTimerTaskRequest,
    enableTimerTaskRequest,
    disableTimerTaskRequest,
    deleteTimerTaskRequest,
    triggerTimerTaskRequest,
    updateTimerTaskRequest,
    pageTimerTaskLogRequest
}