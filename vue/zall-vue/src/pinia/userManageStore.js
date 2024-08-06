import { defineStore } from 'pinia'

export const useUserManageStore = defineStore('userManageStore', {
    state: () => ({
        account: "",
        email: "",
        name: "",
        avatarUrl: ""
    })
})