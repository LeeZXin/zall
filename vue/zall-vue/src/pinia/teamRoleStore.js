import { defineStore } from 'pinia'

export const useTeamRoleStore = defineStore('teamRoleStore', {
    state: () => ({
        roleId: 0,
        name: "",
        teamId: 0,
        teamPerm: {},
        defaultRepoPerm: {},
        repoPermList: [],
        developAppList: []
    })
})