import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 制品列表
const listProductRequest = (data, env) => request.get("/api/product/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除制品
const deleteProductRequest = (productId, env) => request.delete("/api/product/delete/" + productId, {
    headers: {
        [ENV_HEADER]: env
    }
});
export {
    listProductRequest,
    deleteProductRequest
}