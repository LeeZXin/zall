import request from '@/utils/request.js'
// 获取团队列表
const getTeamListRequest = () => request.get("/api/team/list");
// 创建团队
const createTeamRequest = (data) => request.post("/api/team/create", data);
// 编辑团队
const updateTeamRequest = (data) => request.post("/api/team/update", data);
// 获取单个团队信息
const getTeamRequest = (teamId) => request.get("/api/team/get/" + teamId);
// 获取团队成员账号
const listUserByTeamIdRequest = (teamId) => request.get("/api/teamUser/listByTeamId/" + teamId);
// 团队角色列表
const listRolesRequest = (teamId) => request.get("/api/teamRole/list/" + teamId);
// 创建角色
const createRoleRequest = (data) => request.post("/api/teamRole/create", data);
// 编辑角色
const updateRoleRequest = (data) => request.post("/api/teamRole/update", data);
// 删除角色
const deleteRoleRequest = (roleId) => request.delete("/api/teamRole/delete/" + roleId);
// 展示角色成员
const listRoleUserRequest = (teamId) => request.get("/api/teamUser/listRoleUser/" + teamId);
// 创建角色成员关系
const createTeamUserRequest = (data) => request.post("/api/teamUser/create", data);
// 删除团队成员绑定关系
const deleteTeamUserRequest = (relationId) => request.delete("/api/teamUser/delete/" + relationId);
// 更换角色
const changeRoleRequest = (data) => request.post("/api/teamUser/change", data);
// 删除团队
const deleteTeamRequest = (teamId) => request.delete("/api/team/delete/" + teamId);
// 管理员获取所有团队列表
const listAllByAdminRequest = () => request.get("/api/team/listAllByAdmin");
export {
    getTeamListRequest,
    createTeamRequest,
    getTeamRequest,
    listUserByTeamIdRequest,
    listRolesRequest,
    createRoleRequest,
    updateRoleRequest,
    deleteRoleRequest,
    listRoleUserRequest,
    createTeamUserRequest,
    deleteTeamUserRequest,
    changeRoleRequest,
    updateTeamRequest,
    deleteTeamRequest,
    listAllByAdminRequest
}