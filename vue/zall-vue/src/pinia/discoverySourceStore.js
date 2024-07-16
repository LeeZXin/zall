import { defineStore } from 'pinia'

export const useDiscoverySourceStore = defineStore('discoverySource', {
    state: () => ({
        id: 0,
        name: "",
        env: "",
        endpoints: "",
        username: "",
        password: ""
    })
})