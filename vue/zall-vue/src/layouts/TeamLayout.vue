<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{teamStore.name}}</span>
      <span class="switch-text" @click="switchTeam">{{t("switchTeam")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider v-model:collapsed="collapsed" collapsible>
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="onselect">
          <a-menu-item key="/gitRepo/list">
            <branches-outlined />
            <span>{{t("teamMenu.gitRepo")}}</span>
          </a-menu-item>
          <a-menu-item key="/app/list">
            <appstore-outlined />
            <span>{{t("teamMenu.app")}}</span>
          </a-menu-item>
          <a-menu-item key="/timerTask/list" v-if="teamStore.perm?.canManageTimer">
            <clock-circle-outlined />
            <span>{{t("teamMenu.timerTask")}}</span>
          </a-menu-item>
          <a-menu-item key="/role/list" v-if="teamStore.isAdmin">
            <user-outlined />
            <span>{{t("teamMenu.roleAndMembers")}}</span>
          </a-menu-item>
          <a-menu-item key="/settings" v-if="teamStore.isAdmin">
            <setting-outlined />
            <span>设置</span>
          </a-menu-item>
        </a-menu>
      </a-layout-sider>
      <a-layout-content style="height: calc(100vh - 64px); overflow: scroll;background-color:white">
        <router-view v-if="teamInfoLoaded" />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
import AvatarName from "../components/user/AvatarName";
import { useTeamStore } from "../pinia/teamStore";
import { useI18n } from "vue-i18n";
import { ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  BranchesOutlined,
  AppstoreOutlined,
  ClockCircleOutlined,
  SettingOutlined,
  UserOutlined
} from "@ant-design/icons-vue";
import { getTeamRequest } from "@/api/team/teamApi";
const teamStore = useTeamStore();
const teamInfoLoaded = ref(false);
const router = useRouter();
const { t } = useI18n();
const collapsed = ref(false);
const route = useRoute();
const selectedKeys = ref([]);
const routeKey = `/team/${route.params.teamId}`;
const switchTeam = () => {
  router.push("/index");
};
const onselect = event => {
  router.push({
    path: routeKey + event.key,
    force: true
  });
};
const changeSelectedKey = path => {
  const routeSuffix = path.replace(new RegExp(`^${routeKey}`), "");
  for (let key in pagesMap) {
    let value = pagesMap[key];
    if (routeSuffix.startsWith(key)) {
      selectedKeys.value = [value];
      break;
    }
  }
};
// 为了子页面能体现在导航栏
const pagesMap = {
  "/gitRepo": "/gitRepo/list",
  "/role": "/role/list",
  "/timerTask": "/timerTask/list",
  "/app": "/app/list",
  "/settings": "/settings"
};
changeSelectedKey(route.path);
getTeamRequest(route.params.teamId).then(res => {
  teamStore.teamId = res.data.teamId;
  teamStore.name = res.data.name;
  teamStore.isAdmin = res.data.isAdmin;
  teamStore.perm = res.data.perm;
  teamInfoLoaded.value = true;
});
watch(
  () => router.currentRoute.value.path,
  newPath => changeSelectedKey(newPath)
);
</script>
<style scoped>
</style>