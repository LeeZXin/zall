import { defineStore } from 'pinia'

export const useNotifyTplStore = defineStore('notifyTpl', {
    state: () => ({
        id: 0,
        name: "",
        url: "",
        notifyType: "",
        template: "",
        feishuSignKey: ""
    })
})