<template>
  <div style="padding:10px" class="container">
    <div v-if="allowUserCreateTeam" class="flex-right" style="margin-bottom: 10px">
      <a-button
        type="primary"
        @click="toCreateTeam"
        :icon="h(PlusOutlined)"
      >{{t("team.createTeam")}}</a-button>
    </div>
    <div class="team-list">
      <div class="header no-wrap">
        <TeamOutlined />
        <span style="margin-left: 8px">{{t("team.myTeam")}}</span>
      </div>
      <ul class="body" v-if="teamList.length > 0">
        <li
          v-for="item in teamList"
          v-bind:key="item.teamId"
          @click="selectTeam(item)"
        >{{item.name}}</li>
      </ul>
      <div class="no-team" v-if="teamList.length === 0">{{t("team.joinNoTeam")}}</div>
    </div>
  </div>
</template>
<script setup>
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import { getTeamListRequest } from "@/api/team/teamApi";
import { ref, h } from "vue";
import { useTeamStore } from "@/pinia/teamStore";
import { getSysCfgRequest } from "@/api/cfg/cfgApi";
import { PlusOutlined, TeamOutlined } from "@ant-design/icons-vue";
// 团队store
const teamStore = useTeamStore();
// i18n
const { t } = useI18n();
// 团队列表
const teamList = ref([]);
const router = useRouter();
// 跳转创建团队页面
const toCreateTeam = () => {
  router.push("/index/team/create");
};
// 选择团队
const selectTeam = team => {
  teamStore.teamId = team.teamId;
  teamStore.name = team.name;
  router.push(`/team/${team.teamId}/gitRepo/list`);
};
// 是否允许用户创建团队
const allowUserCreateTeam = ref(false);
// 获取团队列表
getTeamListRequest().then(res => {
  teamList.value = res.data;
});
// 获取系统配置
getSysCfgRequest().then(res => {
  allowUserCreateTeam.value = res.data.allowUserCreateTeam;
});
</script>
<style scoped>
.team-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.team-list > .header {
  font-weight: bold;
  font-size: 16px;
  line-height: 44px;
  padding: 0 10px;
}
.team-list > .body {
  border-top: 1px solid #d9d9d9;
}
.team-list > .body > li {
  cursor: pointer;
  font-size: 14px;
  line-height: 40px;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  padding: 0 20px;
}
.team-list > .body > li + li {
  border-top: 1px solid #d9d9d9;
}
.team-list > .body > li:hover {
  background-color: #f0f0f0;
}
.no-team {
  border-top: 1px solid #d9d9d9;
  text-align: center;
  font-size: 16px;
  padding: 24px 0;
}
</style>