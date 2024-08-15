import { defineStore } from 'pinia'

export const usePromScrapeStore = defineStore('promScrape', {
    state: () => ({
        id: 0,
        env: "",
        endpoint: "",
        targetType: 0,
        target: "",
        appId: ""
    })
})