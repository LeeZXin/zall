import request from '@/utils/request.js'

// mysql数据源列表
const listMysqlDbRequest = (data) => request.get("/api/mysqldb/list", { params: data });
// 所有数据库列表
const getAllMysqlDbRequest = () => request.get("/api/mysqldb/all");
// 删除mysql数据源
const deleteMysqlDbRequest = (dbId) => request.delete("/api/mysqldb/delete/" + dbId);
// 创建mysql数据源
const createMysqlDbRequest = (data) => request.post("/api/mysqldb/create", data);
// 编辑mysql数据源
const updateMysqlDbRequest = (data) => request.post("/api/mysqldb/update", data);
// 用户读权限列表
const listReadPermByOperatorRequest = (data) => request.get("/api/mysqlReadPerm/list", { params: data });
// 申请列表
const listReadPermApplyByOperatorRequest = (data) => request.get("/api/mysqlReadPerm/listApply", { params: data });
// dba查看申请列表
const listReadPermApplyByDbaRequest = (data) => request.get("/api/mysqlReadPermApply/list", { params: data });
// 申请读权限
const applyReadPermRequest = (data) => request.post("/api/mysqlReadPerm/apply", data);
// 撤销申请读权限
const cancelReadPermRequest = (applyId) => request.put("/api/mysqlReadPerm/cancel/" + applyId);
// 同意申请读权限
const agreeReadPermRequest = (applyId) => request.put("/api/mysqlReadPermApply/agree/" + applyId);
// 不同意申请读权限
const disagreeReadPermRequest = (data) => request.put("/api/mysqlReadPermApply/disagree", data);
// 查看审批单
const getReadPermApplyRequest = (applyId) => request.get("/api/mysqlReadPerm/getApply/" + applyId);
// 查看授权过的数据库
const listAuthorizedDbRequest = () => request.get("/api/mysqlSearch/listAuthorizedDb");
// 查看授权过的库
const listAuthorizedBaseRequest = (dbId) => request.get("/api/mysqlSearch/listAuthorizedBase/" + dbId);
// 查看授权过的表
const listAuthorizedTableRequest = (data) => request.get("/api/mysqlSearch/listAuthorizedTable", { params: data });
// 查看建表语句
const getCreateTableSqlRequest = (data) => request.get("/api/mysqlSearch/getCreateTableSql", { params: data });
// 查看索引语句
const showTableIndexRequest = (data) => request.get("/api/mysqlSearch/showTableIndex", { params: data });
// 执行select语句
const executeSelectSqlRequest = (data) => request.post("/api/mysqlSearch/executeSelectSql", data);
export {
    listMysqlDbRequest,
    getAllMysqlDbRequest,
    deleteMysqlDbRequest,
    createMysqlDbRequest,
    updateMysqlDbRequest,
    listReadPermByOperatorRequest,
    listReadPermApplyByDbaRequest,
    listReadPermApplyByOperatorRequest,
    applyReadPermRequest,
    cancelReadPermRequest,
    agreeReadPermRequest,
    disagreeReadPermRequest,
    getReadPermApplyRequest,
    listAuthorizedDbRequest,
    listAuthorizedBaseRequest,
    listAuthorizedTableRequest,
    getCreateTableSqlRequest,
    showTableIndexRequest,
    executeSelectSqlRequest
}