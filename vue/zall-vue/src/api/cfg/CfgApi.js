import request from '@/utils/request.js'

const getSysCfgRequest = () => request({
    url: "/api/sysCfg/get",
    method: "GET"
});

export {
    getSysCfgRequest
}