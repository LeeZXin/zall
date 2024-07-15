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
          <a-menu-item key="/propertySource/list">
            <book-outlined />
            <span>配置中心来源</span>
          </a-menu-item>
          <a-menu-item key="/propertyFile/list">
            <pull-request-outlined />
            <span>配置中心</span>
          </a-menu-item>
          <a-menu-item key="/pipeline/list">
            <tag-outlined />
            <span>部署流水线</span>
          </a-menu-item>
          <a-menu-item key="/deployPlan/list">
            <file-outlined />
            <span>发布计划</span>
          </a-menu-item>
          <a-menu-item key="/serviceSource/list">
            <flag-outlined />
            <span>服务状态来源</span>
          </a-menu-item>
          <a-menu-item key="/serviceStatus/list">
            <read-outlined />
            <span>服务状态</span>
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
import { useI18n } from "vue-i18n";
import { ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  BookOutlined,
  FileOutlined,
  PullRequestOutlined,
  TagOutlined,
  FlagOutlined,
  ReadOutlined
} from "@ant-design/icons-vue";
const { t } = useI18n();
const route = useRoute();
const collapsed = ref(false);
const router = useRouter();
const selectedKeys = ref([]);
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
  "/propertySource": "/propertySource/list",
  "/propertyFile": "/propertyFile/list",
  "/deployConfig": "/deployConfig/list",
  "/pipeline": "/pipeline/list",
  "/deployPlan": "/deployPlan/list",
  "/serviceSource": "/serviceSource/list",
  "/serviceStatus": "/serviceStatus/list"
};
watch(
  () => router.currentRoute.value.path,
  newPath => changeSelectedKey(newPath)
);
changeSelectedKey(route.path);
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