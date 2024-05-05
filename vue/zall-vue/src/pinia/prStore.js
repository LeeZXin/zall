import { defineStore } from 'pinia'

export const usePrStore = defineStore('prStore', {
    state: () => ({
        commentCount: 0,
        createBy: "",
        created: "",
        head: "",
        headCommitId: "",
        headType: "",
        id: 0,
        prComment: "",
        prStatus: 0,
        prTitle: "",
        repoId: 0,
        target: "",
        targetCommitId: "",
        targetType: ""
    })
})