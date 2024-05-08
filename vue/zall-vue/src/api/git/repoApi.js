import request from '@/utils/request.js'
// 创建仓库
const createRepoRequest = (data) => request.post("/api/gitRepo/create", data);
// gitIgnore模板列表
const allGitIgnoreTemplateListRequest = () => request.get("/api/gitRepo/allGitIgnoreTemplateList")
    // 获取仓库列表
const getRepoListRequest = (teamId) => request.get("/api/gitRepo/list/" + teamId);
// 获取仓库信息
const getRepoRequest = (repoId) => request.get("/api/gitRepo/get/" + repoId);
// 获取代码信息
const treeRepoRequest = (data) => request.get("/api/gitRepo/index", { params: data });
// 所有分支
const allBranchesRequest = (repoId) => request.get("/api/gitRepo/allBranches/" + repoId);
// 分页分支+最后提交信息
const pageBranchCommitsRequest = (data) => request.get("/api/gitRepo/pageBranchCommits", { params: data });
// 分页标签+提交信息
const pageTagCommitsRequest = (data) => request.get("/api/gitRepo/pageTagCommits", { params: data });
// 基本信息
const simpleInfoRequest = (repoId) => request.get("/api/gitRepo/simpleInfo/" + repoId);
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
export {
    createRepoRequest,
    allGitIgnoreTemplateListRequest,
    getRepoListRequest,
    getRepoRequest,
    treeRepoRequest,
    allBranchesRequest,
    simpleInfoRequest,
    entriesRepoRequest,
    catFileRequest,
    blameRequest,
    diffRefsRequest,
    diffFileRequest,
    pageBranchCommitsRequest,
    deleteBranchRequest,
    historyCommitsRequest,
    diffCommitsRequest,
    pageTagCommitsRequest,
    deleteTagRequest
}