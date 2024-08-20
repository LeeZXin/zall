import { defineStore } from 'pinia'

export const useTeamHookStore = defineStore('teamHook', {
    state: () => ({
        events: {},
        hookCfg: {},
        hookType: 0,
        id: 0,
        name: "",
        teamId: 0,
    })
})