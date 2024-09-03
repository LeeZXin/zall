import request from '@/utils/request.js';

// token列表
const listAccessTokenRequest = (data) => request.get("/api/weworkAccessToken/list", { params: data });
// 创建token
const createAccessTokenRequest = (data) => request.post("/api/weworkAccessToken/create", data);
// 编辑token
const updateAccessTokenRequest = (data) => request.post("/api/weworkAccessToken/update", data);
// 删除
const deleteAccessTokenRequest = (tokenId) => request.delete("/api/weworkAccessToken/delete/" + tokenId);
// 刷新
const refreshAccessTokenRequest = (tokenId) => request.put("/api/weworkAccessToken/refresh/" + tokenId);
// 变更api key
const changeAccessTokenApiKeyRequest = (tokenId) => request.put("/api/weworkAccessToken/changeApiKey/" + tokenId);

export {
    listAccessTokenRequest,
    createAccessTokenRequest,
    updateAccessTokenRequest,
    deleteAccessTokenRequest,
    refreshAccessTokenRequest,
    changeAccessTokenApiKeyRequest
}