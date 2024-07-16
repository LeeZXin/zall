import { defineStore } from 'pinia'

export const useMysqldbStore = defineStore('mysqldb', {
    state: () => ({
        id: 0,
        name: "",
        writeHost: "",
        wirteUsername: "",
        writePassword: "",
        readHost: "",
        readUsername: "",
        readPassword: ""
    })
})