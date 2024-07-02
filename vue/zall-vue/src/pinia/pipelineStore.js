import { defineStore } from 'pinia'

export const usePipelineStore = defineStore('pipeline', {
    state: () => ({
        id: 0,
        name: "",
        env: "",
        config: ""
    })
})