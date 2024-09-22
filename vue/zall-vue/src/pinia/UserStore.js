import { defineStore } from 'pinia'

export const useUserStore = defineStore('userStore', {
    state: () => ({
        account: "",
        name: "",
        email: "",
        isProhibited: false,
        avatarUrl: "",
        isAdmin: false,
        isDba: false,
    })
})