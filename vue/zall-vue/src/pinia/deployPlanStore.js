import { defineStore } from 'pinia'

export const useDeloyPlanStore = defineStore('deployPlan', {
    state: () => ({
        id: 0,
        env: ""
    })
})