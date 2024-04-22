import { defineStore } from 'pinia'

export const useTeamStore = defineStore('teamStore', {
    state: () => ({
        teamId: 0,
        name: ""
    })
})