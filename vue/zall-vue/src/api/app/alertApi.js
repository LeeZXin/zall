import request from '@/utils/request.js';

// 监控告警列表
const listAlertConfigRequest = (data) => request.get("/api/alertConfig/list", { params: data });
// 创建监控告警
const createAlertConfigRequest = (data) => request.post("/api/alertConfig/create", data);
// 启动监控告警
const enableAlertConfigRequest = (configId) => request.put("/api/alertConfig/enable/" + configId);
// 停止监控告警
const disableAlertConfigRequest = (configId) => request.put("/api/alertConfig/disable/" + configId);
// 删除监控告警
const deleteAlertConfigRequest = (configId) => request.delete("/api/alertConfig/delete/" + configId);
// 编辑监控告警
const updateAlertConfigRequest = (data) => request.post("/api/alertConfig/update", data);
export {
    listAlertConfigRequest,
    createAlertConfigRequest,
    enableAlertConfigRequest,
    disableAlertConfigRequest,
    deleteAlertConfigRequest,
    updateAlertConfigRequest,
}