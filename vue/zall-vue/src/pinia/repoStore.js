import { defineStore } from 'pinia'

export const useRepoStore = defineStore('repoStore', {
    state: () => ({
        repoId: 0,
        name: "",
        teamId: 0
    })
})