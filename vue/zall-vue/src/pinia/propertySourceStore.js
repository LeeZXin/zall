import { defineStore } from 'pinia'

export const usePropertySourceStore = defineStore('propertySource', {
    state: () => ({
        id: 0,
        name: "",
        env: "",
        endpoints: "",
        username: "",
        password: ""
    })
})