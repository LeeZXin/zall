import { defineStore } from 'pinia'

export const useTimerTaskStore = defineStore('timerTaskStore', {
    state: () => ({
        cronExp: "",
        env: "0",
        id: 0,
        isEnabled: false,
        name: "",
        task: {},
        teamId: 0,

    })
})