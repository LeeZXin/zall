import { defineStore } from 'pinia'

export const userStore = defineStore('userStore', {
    state: () => ({
        account: "",
        name: "",
        email: "",
        isProhibited: false,
        avatarUrl: "",
        isAdmin: true,
        roleType: 0,
        sessionExpireAt: 0,
        sessionId: ""
    })
})