import { defineStore } from 'pinia'

export const usePipelineVarsStore = defineStore('pipelineVars', {
    state: () => ({
        id: 0,
        name: "",
        env: "",
    })
})