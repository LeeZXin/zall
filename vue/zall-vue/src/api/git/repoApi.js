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
// 文件差异内容
const diffFileRequest = (data) => request.get("/api/gitRepo/diffFile", { params: data });
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
    diffFileRequest
}