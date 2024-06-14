import { defineStore } from 'pinia'

export const useProbeStore = defineStore('probe', {
    state: () => ({
        id: 0,
        name: "",
        env: "",
        config: {}
    })
})