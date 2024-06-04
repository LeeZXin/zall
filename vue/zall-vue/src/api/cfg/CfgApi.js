import request from '@/utils/request.js'

// 获取系统配置
const getSysCfgRequest = () => request.get("/api/sysCfg/get");
// 获取环境列表
const getEnvCfgRequest = () => request.get("/api/envCfg/get");
// 单元配置
const getZonesCfgRequest = () => request.get("/api/zonesCfg/get");

export {
    getSysCfgRequest,
    getEnvCfgRequest,
    getZonesCfgRequest
}