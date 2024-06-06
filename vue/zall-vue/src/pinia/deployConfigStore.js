import { defineStore } from 'pinia'

export const useDeloyConfigStore = defineStore('deployConfig', {
    state: () => ({
        id: 0,
        name: "",
        env: "",
        content: ""
    })
})