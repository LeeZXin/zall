import request from '@/utils/request.js';

// 提交合并请求
const submitPullRequestRequest = (data) => request.post("/api/pullRequest/submit", data);
// 合并请求列表
const listPullRequestRequest = (data) => request.get("/api/pullRequest/list", { params: data });
// 获取合并请求
const getPullRequestRequest = (prId) => request.get("/api/pullRequest/get/" + prId);
// 关闭合并请求
const closePullRequestRequest = (prId) => request.put("/api/pullRequest/close/" + prId);
// 统计数据
const statsPullRequestRequest = (repoId) => request.get("/api/pullRequest/stats/" + repoId);
// 展示时间轴
const listTimelineRequest = (prId) => request.get("/api/pullRequest/listTimeline/" + prId);
// 提交评论
const addCommentRequest = (data) => request.post("/api/pullRequest/addComment", data);
// 删除评论
const deleteCommentRequest = (commentId) => request.delete("/api/pullRequest/deleteComment/" + commentId);
// 是否可提交合并请求
const canMergeRequest = (prId) => request.get("/api/pullRequest/canMerge/" + prId);
// 提交合并请求
const mergePullRequestRequest = (prId) => request.put("/api/pullRequest/merge/" + prId);
export {
    submitPullRequestRequest,
    listPullRequestRequest,
    getPullRequestRequest,
    closePullRequestRequest,
    statsPullRequestRequest,
    listTimelineRequest,
    addCommentRequest,
    deleteCommentRequest,
    canMergeRequest,
    mergePullRequestRequest
}