import { defineStore } from 'pinia'

export const useServiceSourceStore = defineStore('serviceSource', {
    state: () => ({
        id: 0,
        name: "",
        env: "",
        datasource: ""
    })
})