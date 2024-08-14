import request from '@/utils/request.js'
// 创建节点
const createZalletNodeRequest = (data) => request.post("/api/zalletNode/create", data);
// 编辑节点
const updateZalletNodeRequest = (data) => request.post("/api/zalletNode/update", data);
// 删除节点
const deleteZalletNodeRequest = (nodeId) => request.delete("/api/zalletNode/delete/" + nodeId);
// 节点列表
const listZalletNodeRequest = (data) => request.get("/api/zalletNode/list", { params: data });
// 所有节点
const listAllZalletNodeRequest = () => request.get("/api/zalletNode/listAll");

export {
    createZalletNodeRequest,
    updateZalletNodeRequest,
    deleteZalletNodeRequest,
    listZalletNodeRequest,
    listAllZalletNodeRequest
}