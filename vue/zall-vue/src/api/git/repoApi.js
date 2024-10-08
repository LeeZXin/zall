import request from '@/utils/request.js'
// 创建仓库
const createRepoRequest = (data) => request.post("/api/gitRepo/create", data);
// gitIgnore模板列表
const allGitIgnoreTemplateListRequest = () => request.get("/api/gitRepo/allGitIgnoreTemplateList")
    // 获取仓库列表
const getRepoListRequest = (teamId) => request.get("/api/gitRepo/list/" + teamId);
// 获取已删除仓库列表
const getDeletedRepoListRequest = (teamId) => request.get("/api/gitRepo/listDeleted/" + teamId);
// 获取仓库信息
const getRepoRequest = (repoId) => request.get("/api/gitRepo/get/" + repoId);
// 获取代码信息
const treeRepoRequest = (data) => request.get("/api/gitRepo/index", { params: data });
// 所有分支
const allBranchesRequest = (repoId) => request.get("/api/gitRepo/allBranches/" + repoId);
// 分页分支+最后提交信息
const listBranchCommitsRequest = (data) => request.get("/api/gitRepo/listBranchCommits", { params: data });
// 分页标签+提交信息
const listTagCommitsRequest = (data) => request.get("/api/gitRepo/listTagCommits", { params: data });
// 基本信息
const getBaseInfoRequest = (repoId) => request.get("/api/gitRepo/base/" + repoId);
// 文件列表
const entriesRepoRequest = (data) => request.get("/api/gitRepo/entries", { params: data });
// 文件详细内容
const catFileRequest = (data) => request.get("/api/gitRepo/catFile", { params: data });
// 获取每一行提交信息
const blameRequest = (data) => request.get("/api/gitRepo/blame", { params: data });
// 比较差异
const diffRefsRequest = (data) => request.get("/api/gitRepo/diffRefs", { params: data });
// 比较差异
const diffCommitsRequest = (data) => request.get("/api/gitRepo/diffCommits", { params: data });
// 文件差异内容
const diffFileRequest = (data) => request.get("/api/gitRepo/diffFile", { params: data });
// 删除分支
const deleteBranchRequest = (data) => request.delete("/api/gitRepo/deleteBranch", { params: data });
// 删除tag
const deleteTagRequest = (data) => request.delete("/api/gitRepo/deleteTag", { params: data });
// 提交历史
const historyCommitsRequest = (data) => request.get("/api/gitRepo/historyCommits", { params: data });
// gc
const gcRequest = (repoId) => request.put("/api/gitRepo/gc/" + repoId);
// 更新仓库配置
const updateRepoRequest = (data) => request.post("/api/gitRepo/update", data);
// 归档仓库
const setArchivedRequest = (repoId) => request.put("/api/gitRepo/setArchived/" + repoId);
// 归档仓库 -> 正常仓库
const setUnArchivedRequest = (repoId) => request.put("/api/gitRepo/setUnArchived/" + repoId);
// 删除仓库
const deleteRepoRequest = (repoId) => request.delete("/api/gitRepo/delete/" + repoId);
// 永久删除仓库
const deleteRepoPermanentlyRequest = (repoId) => request.delete("/api/gitRepo/deletePermanently/" + repoId);
// 从回收站恢复仓库
const reoverFromRecycleRequest = (repoId) => request.put("/api/gitRepo/recoverFromRecycle/" + repoId);
// 管理员展示仓库列表
const getRepoListByAdminRequest = (teamId) => request.get("/api/gitRepo/listByAdmin/" + teamId);
// 详细信息
const getDetailInfoRequest = (repoId) => request.get("/api/gitRepo/detail/" + repoId);
// 迁移团队
const transferRepoRequest = (data) => request.put("/api/gitRepo/transferTeam", data);
export {
    createRepoRequest,
    allGitIgnoreTemplateListRequest,
    getRepoListRequest,
    getRepoRequest,
    treeRepoRequest,
    allBranchesRequest,
    getBaseInfoRequest,
    entriesRepoRequest,
    catFileRequest,
    blameRequest,
    diffRefsRequest,
    diffFileRequest,
    listBranchCommitsRequest,
    deleteBranchRequest,
    historyCommitsRequest,
    diffCommitsRequest,
    listTagCommitsRequest,
    deleteTagRequest,
    gcRequest,
    updateRepoRequest,
    setArchivedRequest,
    setUnArchivedRequest,
    deleteRepoRequest,
    getDeletedRepoListRequest,
    deleteRepoPermanentlyRequest,
    reoverFromRecycleRequest,
    getRepoListByAdminRequest,
    getDetailInfoRequest,
    transferRepoRequest
}