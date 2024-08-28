import request from '@/utils/request.js'

// 制品列表
const listArtifactRequest = (data) => request.get("/api/artifact/list", { params: data });
// 删除制品
const deleteArtifactRequest = (artifactId) => request.delete("/api/artifact/delete/" + artifactId);
export {
    listArtifactRequest,
    deleteArtifactRequest
}