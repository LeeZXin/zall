import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
    state: () => ({
        appId: "",
        name: "",
        perm: null
    })
})