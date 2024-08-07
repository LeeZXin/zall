<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{route.params.appId}}</span>
      <span class="switch-app-text" @click="switchApp">{{t("appService.switchApp")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider v-model:collapsed="collapsed" collapsible>
        <a-menu theme="dark" mode="inline" @click="onselect" v-model:selectedKeys="selectedKeys">
          <a-menu-item key="/propertyFile/list">
            <ContainerOutlined />
            <span>配置中心</span>
          </a-menu-item>
          <a-menu-item key="/pipeline/list" v-if="appStore.perm?.canManagePipeline">
            <TagOutlined />
            <span>部署流水线</span>
          </a-menu-item>
          <a-menu-item key="/deployPlan/list">
            <FileOutlined />
            <span>发布计划</span>
          </a-menu-item>
          <a-menu-item key="/serviceStatus/list">
            <ReadOutlined />
            <span>服务状态</span>
          </a-menu-item>
          <a-menu-item key="/discoveryService/list">
            <BlockOutlined />
            <span>注册中心</span>
          </a-menu-item>
          <a-menu-item key="/product/list">
            <DatabaseOutlined />
            <span>制品库</span>
          </a-menu-item>
          <a-menu-item key="/promScrape/list" v-if="appStore.perm?.canManagePromAgent">
            <alert-outlined />
            <span>Prometheus</span>
          </a-menu-item>
          <a-menu-item key="/settings" v-if="teamStore.isAdmin">
            <SettingOutlined />
            <span>设置</span>
          </a-menu-item>
        </a-menu>
      </a-layout-sider>
      <a-layout-content style="height: calc(100vh - 64px); overflow: scroll;background-color:white">
        <router-view v-if="appLoaded" />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
import AvatarName from "../components/user/AvatarName";
import { useI18n } from "vue-i18n";
import { ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  FileOutlined,
  TagOutlined,
  ReadOutlined,
  ContainerOutlined,
  BlockOutlined,
  SettingOutlined,
  DatabaseOutlined,
  AlertOutlined
} from "@ant-design/icons-vue";
import { getAppRequest } from "@/api/app/appApi";
import { useAppStore } from "@/pinia/appStore";
import { useTeamStore } from "@/pinia/teamStore";
import { getTeamRequest } from "@/api/team/teamApi";
const teamStore = useTeamStore();
const appStore = useAppStore();
const { t } = useI18n();
const route = useRoute();
const collapsed = ref(false);
const router = useRouter();
const selectedKeys = ref([]);
const appLoaded = ref(false);
const routeKey = `/team/${route.params.teamId}/app/${route.params.appId}`;
const switchApp = () => {
  router.push(`/team/${route.params.teamId}/app/list`);
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
  "/propertyFile": "/propertyFile/list",
  "/deployConfig": "/deployConfig/list",
  "/pipeline": "/pipeline/list",
  "/deployPlan": "/deployPlan/list",
  "/serviceStatus": "/serviceStatus/list",
  "/discoveryService": "/discoveryService/list",
  "/settings": "/settings",
  "/product": "/product/list",
  "/promScrape": "/promScrape/list"
};
watch(
  () => router.currentRoute.value.path,
  newPath => changeSelectedKey(newPath)
);
changeSelectedKey(route.path);
// 获取app信息和权限
const getApp = () => {
  getAppRequest(route.params.appId).then(res => {
    appStore.appId = res.data.appId;
    appStore.name = res.data.name;
    appStore.perm = res.data.perm;
    appLoaded.value = true;
  });
};
// 获取团队信息和权限
const getTeam = () => {
  getTeamRequest(route.params.teamId).then(res => {
    teamStore.teamId = res.data.teamId;
    teamStore.name = res.data.name;
    teamStore.isAdmin = res.data.isAdmin;
    teamStore.perm = res.data.perm;
  });
};
getApp();
if (teamStore.teamId === 0) {
  getTeam();
}
</script>
<style scoped>
.switch-app-text {
  color: white;
  margin-left: 12px;
  font-size: 12px;
  cursor: pointer;
}
.switch-app-text:hover {
  color: #1677ff;
}
</style>