import request from '@/utils/request.js';

// 创建保护分支
const createProtectedBranchRequest = (data) => request.post("/api/protectedBranch/create", data);
// 保护分支列表
const listProtectedBranchRequest = (repoId) => request.get("/api/protectedBranch/list/" + repoId);
// 删除保护分支
const deleteProtectedBranchRequest = (id) => request.delete("/api/protectedBranch/delete/" + id);
// 编辑保护分支
const updateProtecteddBranchRequest = (data) => request.post("/api/protectedBranch/update", data);
export {
    createProtectedBranchRequest,
    listProtectedBranchRequest,
    deleteProtectedBranchRequest,
    updateProtecteddBranchRequest
}