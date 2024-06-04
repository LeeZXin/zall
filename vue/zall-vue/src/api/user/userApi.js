import request from '@/utils/request.js'

const listAllUserRequest = () => request.get("/api/user/listAll");

export {
    listAllUserRequest
}