import { defineStore } from 'pinia'

export const useProtectedBranchStore = defineStore('protectedBranch', {
    state: () => ({
        id: 0,
        pattern: "",
        cfg: null
    })
})