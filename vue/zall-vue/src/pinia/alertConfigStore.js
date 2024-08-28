import { defineStore } from 'pinia'

export const useAlertConfigStore = defineStore('alertConfig', {
    state: () => ({
        id: 0,
        name: "",
        content: {},
        intervalSec: 10,
        isEnabled: false,
        env: ""
    })
})