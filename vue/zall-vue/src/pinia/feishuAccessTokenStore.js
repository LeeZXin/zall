import { defineStore } from 'pinia'

export const useFeishuAccessTokenStore = defineStore('feishuAccessToken', {
    state: () => ({
        id: 0,
        name: "",
        appId: "",
        secret: ""
    })
})