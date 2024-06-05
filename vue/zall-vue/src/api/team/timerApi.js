import request from '@/utils/request.js';

const ENV_HEADER = "ZALL-ENV";

// 定时任务列表
const listTimerTaskRequest = (data, env) => request.get("/api/timerTask/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 创建定时任务
const createTimerTaskRequest = (data, env) => request.post("/api/timerTask/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 启动定时任务
const enableTimerTaskRequest = (taskId, env) => request.put("/api/timerTask/enable/" + taskId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 停止定时任务
const disableTimerTaskRequest = (taskId, env) => request.put("/api/timerTask/disable/" + taskId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除定时任务
const deleteTimerTaskRequest = (taskId, env) => request.delete("/api/timerTask/delete/" + taskId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 触发定时任务
const triggerTimerTaskRequest = (taskId, env) => request.put("/api/timerTask/trigger/" + taskId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 编辑定时任务
const updateTimerTaskRequest = (data, env) => request.post("/api/timerTask/update", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 查看定时任务日志
const pageTimerTaskLogRequest = (data, env) => request.get("/api/timerLog/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
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