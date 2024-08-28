import request from '@/utils/request.js'

// 制品列表
const listProductRequest = (data) => request.get("/api/product/list", { params: data });
// 删除制品
const deleteProductRequest = (productId) => request.delete("/api/product/delete/" + productId);
export {
    listProductRequest,
    deleteProductRequest
}