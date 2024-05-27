import request from '@/utils/request.js';

// 分页查询日志
const pageLogRequest = (data) => request.get("/api/oplog/page", { params: data });

export {
    pageLogRequest
}