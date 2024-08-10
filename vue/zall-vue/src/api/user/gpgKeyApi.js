import request from '@/utils/request.js'

// 创建gpg密钥
const createGpgKeyRequest = (data) => request.post("/api/gpgKey/create", data);
// gpg密钥列表
const listAllGpgKeyRequest = () => request.get("/api/gpgKey/list");
// 删除密钥
const deleteGpgKeyRequest = (keyId) => request.delete("/api/gpgKey/delete/" + keyId);
export {
    createGpgKeyRequest,
    listAllGpgKeyRequest,
    deleteGpgKeyRequest
}