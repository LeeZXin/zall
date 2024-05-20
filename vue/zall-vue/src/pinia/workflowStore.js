import { defineStore } from 'pinia'

export const useWorkflowStore = defineStore('workflow', {
    state: () => ({
        id: 0,
        name: "",
        desc: ""
    })
})