import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 探针配置列表
const listProbeRequest = (data, env) => request.get("/api/probe/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 创建探针配置
const createProbeRequest = (data, env) => request.post("/api/probe/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 编辑探针配置
const updateProbeRequest = (data, env) => request.post("/api/probe/update", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除探针配置
const deleteProbeRequest = (probeId, env) => request.delete("/api/probe/delete/" + probeId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 启动探针
const enableProbeRequest = (probeId, env) => request.put("/api/probe/enable/" + probeId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 关闭探针
const disableProbeRequest = (probeId, env) => request.put("/api/probe/disable/" + probeId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
export {
    listProbeRequest,
    createProbeRequest,
    updateProbeRequest,
    deleteProbeRequest,
    enableProbeRequest,
    disableProbeRequest
}