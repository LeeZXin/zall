import { defineStore } from 'pinia'

export const useWeworkAccessTokenStore = defineStore('weworkAccessToken', {
    state: () => ({
        id: 0,
        name: "",
        corpId: "",
        secret: ""
    })
})