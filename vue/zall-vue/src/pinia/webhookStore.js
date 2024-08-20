import { defineStore } from 'pinia'

export const useWebhookStore = defineStore('webhook', {
    state: () => ({
        events: {},
        hookUrl: "",
        id: 0,
        secret: ""
    })
})