import request from '@/utils/request.js'
// 展示所有用户列表
const listAllUserRequest = () => request.get("/api/user/listAll");
// 管理员查看用户列表
const listUserByAdminRequest = (data) => request.get("/api/user/list", { params: data });
// 设置dba角色
const setDbaRequest = (data) => request.put("/api/user/setDba", data);
// 设置系统管理员角色
const setAdminRequest = (data) => request.put("/api/user/setAdmin", data);
// 设置禁用
const setProhibitedRequest = (data) => request.put("/api/user/setProhibited", data);
// 创建用户
const createUserRequest = (data) => request.post("/api/user/create", data);
// 编辑用户
const updateUserRequest = (data) => request.post("/api/user/update", data);
// 重置密码
const resetPasswordRequest = (account) => request.put("/api/user/resetPassword/" + account);
// 删除用户
const deleteUserRequest = (account) => request.delete("/api/user/delete/" + account);
export {
    listAllUserRequest,
    listUserByAdminRequest,
    setDbaRequest,
    setAdminRequest,
    setProhibitedRequest,
    createUserRequest,
    updateUserRequest,
    resetPasswordRequest,
    deleteUserRequest
}