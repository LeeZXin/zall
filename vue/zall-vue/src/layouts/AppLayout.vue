<template>
  <a-watermark :content="`${userStore.name}${userStore.account}`" :gap="[200,200]">
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
              <span>{{t("appMenu.propertyFile")}}</span>
            </a-menu-item>
            <a-menu-item key="/pipeline/list" v-if="appStore.perm?.canManagePipeline">
              <TagOutlined />
              <span>{{t("appMenu.deployPipeline")}}</span>
            </a-menu-item>
            <a-menu-item key="/deployPlan/list">
              <FileOutlined />
              <span>{{t("appMenu.deployPlan")}}</span>
            </a-menu-item>
            <a-menu-item key="/serviceStatus/list">
              <ReadOutlined />
              <span>{{t("appMenu.serviceStatus")}}</span>
            </a-menu-item>
            <a-menu-item key="/discoveryService/list">
              <BlockOutlined />
              <span>{{t("appMenu.registryCenter")}}</span>
            </a-menu-item>
            <a-menu-item key="/artifact/list">
              <DatabaseOutlined />
              <span>{{t("appMenu.artifacts")}}</span>
            </a-menu-item>
            <a-menu-item key="/promScrape/list" v-if="teamStore.perm?.canManagePromScrape">
              <AlertOutlined />
              <span>{{t("appMenu.promScrape")}}</span>
            </a-menu-item>
            <a-menu-item key="/alertConfig/list">
              <MonitorOutlined />
              <span>{{t("appMenu.alertConfig")}}</span>
            </a-menu-item>
            <a-menu-item key="/setting" v-if="teamStore.isAdmin">
              <SettingOutlined />
              <span>{{t("appMenu.setting")}}</span>
            </a-menu-item>
          </a-menu>
        </a-layout-sider>
        <a-layout-content
          style="height: calc(100vh - 64px); overflow: scroll;background-color:white"
        >
          <router-view v-if="appLoaded" />
        </a-layout-content>
      </a-layout>
    </a-layout>
  </a-watermark>
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
  AlertOutlined,
  MonitorOutlined
} from "@ant-design/icons-vue";
import { getAppRequest } from "@/api/app/appApi";
import { useAppStore } from "@/pinia/appStore";
import { useTeamStore } from "@/pinia/teamStore";
import { getTeamRequest } from "@/api/team/teamApi";
import { useUserStore } from "@/pinia/userStore";
const userStore = useUserStore();
const teamStore = useTeamStore();
const appStore = useAppStore();
const { t } = useI18n();
const route = useRoute();
// 导航栏是否合上
const collapsed = ref(false);
const router = useRouter();
// 导航栏选择的keys
const selectedKeys = ref([]);
// appStore数据是否加载完
const appLoaded = ref(false);
// 路由前缀
const routeKey = `/team/${route.params.teamId}/app/${route.params.appId}`;
// 切换应用服务
const switchApp = () => {
  router.push(`/team/${route.params.teamId}/app/list`);
};
// 导航栏选择
const onselect = event => {
  router.push({
    path: routeKey + event.key,
    force: true
  });
};
// 导航栏选择后触发
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
  "/promScrape": "/promScrape/list",
  "/setting": "/setting",
  "/artifact": "/artifact/list",
  "/alertConfig": "/alertConfig/list"
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