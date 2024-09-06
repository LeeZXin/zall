import { defineStore } from 'pinia'

export const useZalletNodeStore = defineStore('zalletNode', {
    state: () => ({
        id: 0,
        nodeId: "",
        name: "",
        agentHost: "",
        agentToken: "",
    })
})