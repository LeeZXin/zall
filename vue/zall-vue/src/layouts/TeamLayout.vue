<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{team.teamName}}</span>
      <span class="switch-team-text" @click="switchTeam">{{t("switchTeam")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider v-model:collapsed="collapsed" collapsible>
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @select="onselect">
          <a-menu-item key="/team/gitRepo/list">
            <branches-outlined />
            <span>{{t("menu.gitRepo")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/app/list">
            <appstore-outlined />
            <span>{{t("menu.app")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/action/list">
            <desktop-outlined />
            <span>{{t("menu.action")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/propertyCenter">
            <file-zip-outlined />
            <span>{{t("menu.propertyCenter")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/timerTask">
            <clock-circle-outlined />
            <span>{{t("menu.timerTask")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/applyApproval">
            <form-outlined />
            <span>{{t("menu.applyApproval")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/dbAudit">
            <database-outlined />
            <span>{{t("menu.dbAudit")}}</span>
          </a-menu-item>
          <a-menu-item key="/team/monitorAlert">
            <alert-outlined />
            <span>{{t("menu.monitorAlert")}}</span>
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
import { useTeamStore } from "../pinia/TeamStore";
import { useI18n } from "vue-i18n";
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  DesktopOutlined,
  BranchesOutlined,
  AppstoreOutlined,
  FileZipOutlined,
  ClockCircleOutlined,
  FormOutlined,
  DatabaseOutlined,
  AlertOutlined
} from "@ant-design/icons-vue";
const { t } = useI18n();
const collapsed = ref(false);
const router = useRouter();
const route = useRoute();
const selectedKeys = ref([]);
const team = useTeamStore();
const switchTeam = () => {
  router.push("/index");
};
const onselect = event => {
  router.push(event.key);
};
// 为了子页面能体现在导航栏
const pagesMap = {
  "/team/gitRepo": "/team/gitRepo/list",
  "/team/app": "/team/app/list",
  "/team/action": "/team/action/list"
};
for (let key in pagesMap) {
  let value = pagesMap[key];
  if (route.path.startsWith(key)) {
    selectedKeys.value = [value];
    break;
  }
}
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