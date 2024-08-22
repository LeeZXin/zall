import request from '@/utils/request.js';

// 定时任务列表
const listTimerRequest = (data) => request.get("/api/timer/list", { params: data });
// 创建定时任务
const createTimerRequest = (data) => request.post("/api/timer/create", data);
// 启动定时任务
const enableTimerRequest = (timerId) => request.put("/api/timer/enable/" + timerId);
// 停止定时任务
const disableTimerRequest = (timerId) => request.put("/api/timer/disable/" + timerId);
// 删除定时任务
const deleteTimerRequest = (timerId) => request.delete("/api/timer/delete/" + timerId);
// 触发定时任务
const triggerTimerTaskRequest = (timerId) => request.put("/api/timer/trigger/" + timerId);
// 编辑定时任务
const updateTimerRequest = (data) => request.post("/api/timer/update", data);
// 查看定时任务日志
const listTimerTaskLogRequest = (data) => request.get("/api/timerLog/list", { params: data });
export {
    listTimerRequest,
    createTimerRequest,
    enableTimerRequest,
    disableTimerRequest,
    deleteTimerRequest,
    triggerTimerTaskRequest,
    updateTimerRequest,
    listTimerTaskLogRequest,
}