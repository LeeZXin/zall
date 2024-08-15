<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>超级管理员</span>
      <span class="switch-text" @click="switchIndex">返回首页</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider collapsible>
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="onselect">
          <a-menu-item key="/cfg/list">
            <SettingOutlined />
            <span>系统配置</span>
          </a-menu-item>
          <a-menu-item key="/user/list">
            <UserOutlined />
            <span>用户管理</span>
          </a-menu-item>
          <a-menu-item key="/propertySource/list">
            <BookOutlined />
            <span>配置中心来源</span>
          </a-menu-item>
          <a-menu-item key="/serviceSource/list">
            <BookOutlined />
            <span>服务状态来源</span>
          </a-menu-item>
          <a-menu-item key="/discoverySource/list">
            <BookOutlined />
            <span>注册中心来源</span>
          </a-menu-item>
          <a-menu-item key="/zalletNode/list">
            <ClusterOutlined />
            <span>Zallet代理节点</span>
          </a-menu-item>
          <a-menu-item key="/promScrape/list">
            <AlertOutlined />
            <span>Prometheus</span>
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
import { useRouter, useRoute } from "vue-router";
import { ref, watch } from "vue";
import {
  UserOutlined,
  SettingOutlined,
  BookOutlined,
  ClusterOutlined,
  AlertOutlined
} from "@ant-design/icons-vue";
const route = useRoute();
const router = useRouter();
const selectedKeys = ref([]);
const routeKey = "/sa";
const onselect = event => {
  router.push({
    path: routeKey + event.key,
    force: true
  });
};
const switchIndex = () => {
  router.push("/index");
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
  "/cfg": "/cfg/list",
  "/user": "/user/list",
  "/propertySource": "/propertySource/list",
  "/serviceSource": "/serviceSource/list",
  "/discoverySource": "/discoverySource/list",
  "/zalletNode": "/zalletNode/list",
  "/promScrape": "/promScrape/list"
};
changeSelectedKey(route.path);
watch(
  () => router.currentRoute.value.path,
  newPath => changeSelectedKey(newPath)
);
</script>
<style scoped>
</style>