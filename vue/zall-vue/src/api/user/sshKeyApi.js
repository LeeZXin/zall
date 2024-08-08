import request from '@/utils/request.js'

// 创建ssh密钥
const createSshKeyRequest = (data) => request.post("/api/sshKey/create", data);
// ssh密钥列表
const listAllSshKeyRequest = () => request.get("/api/sshKey/list");
// 删除密钥
const deleteSshKeyRequest = (keyId) => request.delete("/api/sshKey/delete/" + keyId);
export {
    createSshKeyRequest,
    listAllSshKeyRequest,
    deleteSshKeyRequest
}