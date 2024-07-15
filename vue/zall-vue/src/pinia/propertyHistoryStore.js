import { defineStore } from 'pinia'

export const usePropertyHistoryStore = defineStore('propertyHistory', {
    state: () => ({
        id: 0,
        fileName: "",
        fileId: 0,
        content: "",
        version: "",
        created: "",
        creator: "",
        lastVersion: "",
        env: ""
    })
})