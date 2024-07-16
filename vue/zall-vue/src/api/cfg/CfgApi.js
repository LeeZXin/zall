import request from '@/utils/request.js'

// 获取系统配置
const getSysCfgRequest = () => request.get("/api/sysCfg/get");
// 获取环境列表
const getEnvCfgRequest = () => request.get("/api/envCfg/get");

export {
    getSysCfgRequest,
    getEnvCfgRequest
}