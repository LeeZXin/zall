<template>
  <a-watermark :content="`${userStore.name}${userStore.account}`" :gap="[200,200]">
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
              <BranchesOutlined />
              <span>{{t("teamMenu.gitRepo")}}</span>
            </a-menu-item>
            <a-menu-item key="/app/list">
              <AppstoreOutlined />
              <span>{{t("teamMenu.app")}}</span>
            </a-menu-item>
            <a-menu-item key="/timer/list" v-if="teamStore.perm?.canManageTimer">
              <ClockCircleOutlined />
              <span>{{t("teamMenu.timer")}}</span>
            </a-menu-item>
            <a-menu-item key="/role/list" v-if="teamStore.isAdmin">
              <UserOutlined />
              <span>{{t("teamMenu.roleAndMembers")}}</span>
            </a-menu-item>
            <a-menu-item key="/setting" v-if="teamStore.isAdmin">
              <SettingOutlined />
              <span>{{t("teamMenu.setting")}}</span>
            </a-menu-item>
            <a-menu-item key="/notifyTpl/list" v-if="teamStore.perm?.canManageNotifyTpl">
              <NotificationOutlined />
              <span>{{t("teamMenu.notifyTpl")}}</span>
            </a-menu-item>
            <a-menu-item key="/teamHook/list" v-if="teamStore.perm?.canManageTeamHook">
              <ApiOutlined />
              <span>{{t("teamMenu.teamHook")}}</span>
            </a-menu-item>
            <a-menu-item
              key="/weworkAccessToken/list"
              v-if="teamStore.perm?.canManageWeworkAccessToken"
            >
              <KeyOutlined />
              <span>{{t("teamMenu.weworkAccessToken")}}</span>
            </a-menu-item>
            <a-menu-item
              key="/feishuAccessToken/list"
              v-if="teamStore.perm?.canManageFeishuAccessToken"
            >
              <KeyOutlined />
              <span>{{t("teamMenu.feishuAccessToken")}}</span>
            </a-menu-item>
          </a-menu>
        </a-layout-sider>
        <a-layout-content
          style="height: calc(100vh - 64px); overflow: scroll;background-color:white"
        >
          <router-view v-if="teamInfoLoaded" />
        </a-layout-content>
      </a-layout>
    </a-layout>
  </a-watermark>
</template>
<script setup>
import I18nSelect from "@/components/i18n/I18nSelect";
import AvatarName from "@/components/user/AvatarName";
import { useTeamStore } from "@/pinia/teamStore";
import { useUserStore } from "@/pinia/userStore";
import { useI18n } from "vue-i18n";
import { ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  BranchesOutlined,
  AppstoreOutlined,
  ClockCircleOutlined,
  SettingOutlined,
  UserOutlined,
  NotificationOutlined,
  ApiOutlined,
  KeyOutlined
} from "@ant-design/icons-vue";
import { getTeamRequest } from "@/api/team/teamApi";
const userStore = useUserStore();
const teamStore = useTeamStore();
// 获取团队信息完才开始加载页面
const teamInfoLoaded = ref(false);
const router = useRouter();
const { t } = useI18n();
const collapsed = ref(false);
const route = useRoute();
const selectedKeys = ref([]);
// 路由前缀
const routeKey = `/team/${route.params.teamId}`;
// 切换团队
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
  "/timer": "/timer/list",
  "/app": "/app/list",
  "/setting": "/setting",
  "/notifyTpl": "/notifyTpl/list",
  "/teamHook": "/teamHook/list",
  "/weworkAccessToken": "/weworkAccessToken/list",
  "/feishuAccessToken": "/feishuAccessToken/list"
};
// 获取团队信息
const getTeamInfo = () => {
  getTeamRequest(route.params.teamId).then(res => {
    teamStore.teamId = res.data.teamId;
    teamStore.name = res.data.name;
    teamStore.isAdmin = res.data.isAdmin;
    teamStore.perm = res.data.perm;
    teamInfoLoaded.value = true;
  });
};
watch(
  () => router.currentRoute.value.path,
  newPath => changeSelectedKey(newPath)
);
changeSelectedKey(route.path);
getTeamInfo();
</script>
<style scoped>
</style>