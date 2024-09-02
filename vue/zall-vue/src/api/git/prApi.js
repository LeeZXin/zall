import request from '@/utils/request.js';

// 提交合并请求
const submitPullRequestRequest = (data) => request.post("/api/pullRequest/submit", data);
// 合并请求列表
const listPullRequestRequest = (data) => request.get("/api/pullRequest/list", { params: data });
// 获取合并请求
const getPullRequestRequest = (prId) => request.get("/api/pullRequest/get/" + prId);
// 关闭合并请求
const closePullRequestRequest = (prId) => request.put("/api/pullRequest/close/" + prId);
// 展示时间轴
const listTimelineRequest = (prId) => request.get("/api/pullRequest/listTimeline/" + prId);
// 提交评论
const addCommentRequest = (data) => request.post("/api/pullRequest/addComment", data);
// 删除评论
const deleteCommentRequest = (commentId) => request.delete("/api/pullRequest/deleteComment/" + commentId);
// 是否可提交合并请求
const canMergeRequest = (prId) => request.get("/api/pullRequest/canMerge/" + prId);
// 是否可评审合并请求
const canReviewRequest = (prId) => request.get("/api/pullRequest/canReview/" + prId);
// 同意合并请求
const agreeReviewRequest = (prId) => request.put("/api/pullRequest/agreeReview/" + prId);
// 提交合并请求
const mergePullRequestRequest = (prId) => request.put("/api/pullRequest/merge/" + prId);
// 评审记录 
const listReviewRequest = (prId) => request.get("/api/pullRequest/listReview/" + prId);
export {
    submitPullRequestRequest,
    listPullRequestRequest,
    getPullRequestRequest,
    closePullRequestRequest,
    listTimelineRequest,
    addCommentRequest,
    deleteCommentRequest,
    canMergeRequest,
    mergePullRequestRequest,
    canReviewRequest,
    agreeReviewRequest,
    listReviewRequest
}