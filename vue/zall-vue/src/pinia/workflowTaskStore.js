import { defineStore } from 'pinia'

export const useWorkflowTaskStore = defineStore('workflowTask', {
    state: () => ({
        id: 0,
        triggerType: 0,
        yamlContent: "",
        operator: "",
        created: "",
        branch: "",
        prId: 0,
    })
})