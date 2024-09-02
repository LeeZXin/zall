import request from '@/utils/request.js'

// 获取系统配置
const getSysCfgRequest = () => request.get("/api/sysCfg/get");
// 编辑系统配置
const updateSysCfgRequest = (data) => request.post("/api/sysCfg/update", data);
// 获取环境列表
const getEnvCfgRequest = () => request.get("/api/envCfg/get");
// 编辑环境列表
const updateEnvCfgRequest = (data) => request.post("/api/envCfg/update", data);
// 获取git服务端配置
const getGitRepoServerCfgRequest = () => request.get("/api/gitRepoServerCfg/get");
// 编辑git服务端配置
const updateGitRepoServerCfgRequest = (data) => request.post("/api/gitRepoServerCfg/update", data);
// 获取git配置
const getGitCfgRequest = () => request.get("/api/gitCfg/get");
// 编辑git配置
const updateGitCfgRequest = (data) => request.post("/api/gitCfg/update", data);
// 获取登录配置
const getLoginCfgRequest = () => request.get("/api/loginCfg/get");
// 获取登录配置 超级管理员
const getLoginCfgBySaRequest = () => request.get("/api/loginCfg/getBySa");
// 编辑登录配置
const updateLoginCfgRequest = (data) => request.post("/api/loginCfg/update", data);
export {
    getSysCfgRequest,
    updateSysCfgRequest,
    getEnvCfgRequest,
    updateEnvCfgRequest,
    getGitRepoServerCfgRequest,
    updateGitRepoServerCfgRequest,
    getGitCfgRequest,
    updateGitCfgRequest,
    getLoginCfgRequest,
    updateLoginCfgRequest,
    getLoginCfgBySaRequest
}