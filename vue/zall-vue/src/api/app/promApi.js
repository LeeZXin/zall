import request from '@/utils/request.js'

// prometheus抓取任务列表
const listPromScrapeRequest = (data) => request.get("/api/promScrape/list", { params: data });
// 创建prometheus抓取任务
const createPromScrapeRequest = (data) => request.post("/api/promScrape/create", data);
// 编辑prometheus抓取任务
const updatePromScrapeRequest = (data) => request.post("/api/promScrape/update", data);
// 删除prometheus抓取任务
const deletePromScrapeRequest = (scrapeId) => request.delete("/api/promScrape/delete/" + scrapeId);

export {
    listPromScrapeRequest,
    deletePromScrapeRequest,
    createPromScrapeRequest,
    updatePromScrapeRequest
}