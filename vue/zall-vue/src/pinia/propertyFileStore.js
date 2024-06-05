import { defineStore } from 'pinia'

export const usePropertyFileStore = defineStore('propertyFile', {
    state: () => ({
        id: 0,
        name: "",
        env: ""
    })
})