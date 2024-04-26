<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{team.name}}</span>
      <span class="switch-team-text" @click="switchTeam">{{t("switchTeam")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider v-model:collapsed="collapsed" collapsible>
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @select="onselect">
          <a-menu-item :key="`/team/${route.params.teamId}/gitRepo/list`">
            <branches-outlined />
            <span>{{t("teamMenu.gitRepo")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/app/list">
            <appstore-outlined />
            <span>{{t("teamMenu.app")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/action/list">
            <desktop-outlined />
            <span>{{t("teamMenu.action")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/timerTask">
            <clock-circle-outlined />
            <span>{{t("teamMenu.timerTask")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/applyApproval">
            <form-outlined />
            <span>{{t("teamMenu.applyApproval")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/dbAudit">
            <database-outlined />
            <span>{{t("teamMenu.dbAudit")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/monitorAlert">
            <alert-outlined />
            <span>{{t("teamMenu.monitorAlert")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/monitorAlert" v-if="isAdmin">
            <setting-outlined />
            <span>{{t("teamMenu.teamSettings")}}</span>
          </a-menu-item>
        </a-menu>
      </a-layout-sider>
      <a-layout-content style="height: calc(100vh - 64px); overflow: scroll;background-color:white">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
import AvatarName from "../components/user/AvatarName";
import { useTeamStore } from "../pinia/teamStore";
import { useI18n } from "vue-i18n";
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  DesktopOutlined,
  BranchesOutlined,
  AppstoreOutlined,
  ClockCircleOutlined,
  FormOutlined,
  DatabaseOutlined,
  AlertOutlined,
  SettingOutlined
} from "@ant-design/icons-vue";
import { isTeamAdminRequest } from "@/api/team/teamApi";
const team = useTeamStore();
const router = useRouter();
const { t } = useI18n();
const collapsed = ref(false);

const route = useRoute();
const selectedKeys = ref([]);
const switchTeam = () => {
  router.push("/index");
};
const isAdmin = ref(false);
const onselect = event => {
  router.push(event.key);
};
// 为了子页面能体现在导航栏
const pagesMap = {
  gitRepo: `/team/${route.params.teamId}/gitRepo/list`
};
for (let key in pagesMap) {
  let value = pagesMap[key];
  if (route.path.indexOf(key) >= 0) {
    selectedKeys.value = [value];
    break;
  }
}
isTeamAdminRequest({
  teamId: parseInt(route.params.teamId)
}).then(res => {
  isAdmin.value = res.data;
});
</script>
<style scoped>
.switch-team-text {
  color: white;
  margin-left: 12px;
  font-size: 12px;
  cursor: pointer;
}
.switch-team-text:hover {
  color: #1677ff;
}
</style>