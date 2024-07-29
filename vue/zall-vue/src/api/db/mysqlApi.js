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
const cancelReadPermApplyRequest = (applyId) => request.put("/api/mysqlReadPerm/cancel/" + applyId);
// 同意申请读权限
const agreeReadPermApplyRequest = (applyId) => request.put("/api/mysqlReadPermApply/agree/" + applyId);
// 不同意申请读权限
const disagreeReadPermApplyRequest = (data) => request.put("/api/mysqlReadPermApply/disagree", data);
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
// 查看数据修改审批单列表
const listDataUpdateApplyByOperatorRequest = (data) => request.get("/api/mysqlDataUpdate/listApply", { params: data });
// 申请数据库修改单
const applyDataUpdateRequest = (data) => request.post("/api/mysqlDataUpdate/apply", data);
// 申请数据库修改单的执行计划
const explainDataUpdateApplyRequest = (applyId) => request.get("/api/mysqlDataUpdate/explainApply/" + applyId);
// 取消数据修改单申请
const cancelDataUpdateApplyRequest = (applyId) => request.put("/api/mysqlDataUpdate/cancelApply/" + applyId);
// dba查看数据修改审批单列表
const listDataUpdateApplyByDbaRequest = (data) => request.get("/api/mysqlDataUpdateApply/list", { params: data });
// 同意数据修改单
const agreeDataUpdateApplyRequest = (applyId) => request.put("/api/mysqlDataUpdateApply/agree/" + applyId);
// 请求执行
const askToExecuteDataUpdateApplyRequest = (applyId) => request.put("/api/mysqlDataUpdateApply/askToExecute/" + applyId);
// 执行数据库修改单
const executeDataUpdateApplyRequest = (applyId) => request.put("/api/mysqlDataUpdateApply/execute/" + applyId);
// 不同意数据库修改单
const disagreedataUpdateApplyRequest = (data) => request.put("/api/mysqlDataUpdateApply/disagree", data);
// dba查看读权限列表
const listReadPermByDbaRequest = (data) => request.get("/api/mysqlReadPerm/listManage", { params: data });
// 删除读权限
const deleteReadPermRequest = (permId) => request.delete("/api/mysqlReadPerm/delete/" + permId);
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
    cancelReadPermApplyRequest,
    agreeReadPermApplyRequest,
    disagreeReadPermApplyRequest,
    getReadPermApplyRequest,
    listAuthorizedDbRequest,
    listAuthorizedBaseRequest,
    listAuthorizedTableRequest,
    getCreateTableSqlRequest,
    showTableIndexRequest,
    executeSelectSqlRequest,
    listDataUpdateApplyByOperatorRequest,
    applyDataUpdateRequest,
    explainDataUpdateApplyRequest,
    cancelDataUpdateApplyRequest,
    listDataUpdateApplyByDbaRequest,
    agreeDataUpdateApplyRequest,
    askToExecuteDataUpdateApplyRequest,
    executeDataUpdateApplyRequest,
    disagreedataUpdateApplyRequest,
    listReadPermByDbaRequest,
    deleteReadPermRequest
}