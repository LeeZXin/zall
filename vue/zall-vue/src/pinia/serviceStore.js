import { defineStore } from 'pinia'

export const useServiceStore = defineStore('service', {
    state: () => ({
        id: 0,
        name: "",
        env: "",
        config: ""
    })
})