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
export {
    createRepoRequest,
    allGitIgnoreTemplateListRequest,
    getRepoListRequest
}