import { defineStore } from 'pinia'

export const useUserStore = defineStore('userStore', {
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