import request from '@/utils/request.js'

// prometheus抓取任务列表
const listPromScrapeBySaRequest = (data) => request.get("/api/promScrapeBySa/list", { params: data });
// 创建prometheus抓取任务
const createPromScrapeBySaRequest = (data) => request.post("/api/promScrapeBySa/create", data);
// 编辑prometheus抓取任务
const updatePromScrapeBySaRequest = (data) => request.post("/api/promScrapeBySa/update", data);
// 删除prometheus抓取任务
const deletePromScrapeBySaRequest = (scrapeId) => request.delete("/api/promScrapeBySa/delete/" + scrapeId);
// prometheus抓取任务列表
const listPromScrapeByTeamRequest = (data) => request.get("/api/promScrapeByTeam/list", { params: data });
// 创建prometheus抓取任务
const createPromScrapeByTeamRequest = (data) => request.post("/api/promScrapeByTeam/create", data);
// 编辑prometheus抓取任务
const updatePromScrapeByTeamRequest = (data) => request.post("/api/promScrapeByTeam/update", data);
// 删除prometheus抓取任务
const deletePromScrapeByTeamRequest = (scrapeId) => request.delete("/api/promScrapeByTeam/delete/" + scrapeId);

export {
    listPromScrapeBySaRequest,
    deletePromScrapeBySaRequest,
    createPromScrapeBySaRequest,
    updatePromScrapeBySaRequest,
    listPromScrapeByTeamRequest,
    createPromScrapeByTeamRequest,
    updatePromScrapeByTeamRequest,
    deletePromScrapeByTeamRequest
}