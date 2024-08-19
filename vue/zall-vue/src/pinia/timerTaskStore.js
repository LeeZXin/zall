import { defineStore } from 'pinia'

// 定时任务
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