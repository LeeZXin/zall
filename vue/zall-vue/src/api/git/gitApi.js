import request from '@/utils/request.js'
// 创建仓库
const createRepoRequest = (data) => request({
    url: "/api/gitRepo/create",
    method: "POST",
    data: data
});
// gitIgnore模板列表
const allGitIgnoreTemplateListRequest = () => request({
    url: "/api/gitRepo/allGitIgnoreTemplateList",
    method: "GET"
});
// 获取仓库列表
const getRepoListRequest = (data) => request({
    url: "/api/gitRepo/list",
    method: "POST",
    data: data
});
// 获取仓库信息
const getRepoRequest = (data) => request({
    url: "/api/gitRepo/get",
    method: "POST",
    data: data
});
// 获取代码信息
const treeRepoRequest = (data) => request({
    url: "/api/gitRepo/index",
    method: "POST",
    data: data
});
// 所有分支
const allBranchesRequest = (data) => request({
    url: "/api/gitRepo/allBranches",
    method: "POST",
    data: data
});
// 基本信息
const simpleInfoRequest = (data) => request({
    url: "/api/gitRepo/simpleInfo",
    method: "POST",
    data: data
});
// 文件列表
const entriesRepoRequest = (data) => request({
    url: "/api/gitRepo/entries",
    method: "POST",
    data: data
});
// 文件详细内容
const catFileRequest = (data) => request({
    url: "/api/gitRepo/catFile",
    method: "POST",
    data: data
});
// 获取每一行提交信息
const blameRequest = (data) => request({
    url: "/api/gitRepo/blame",
    method: "POST",
    data: data
});
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
    blameRequest
}