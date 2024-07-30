<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>Mysql审计</span>
      <span class="switch-text" @click="switchIndex">返回首页</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider collapsible>
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="onselect">
          <a-menu-item key="/mysqlDb/list">
            <DatabaseOutlined />
            <span>数据源</span>
          </a-menu-item>
          <a-menu-item key="/mysqlReadPermApply/list">
            <CarryOutOutlined />
            <span>读权限申请</span>
          </a-menu-item>
          <a-menu-item key="/mysqlReadPermAudit/list" v-if="userStore.isDba">
            <AuditOutlined />
            <span>读权限审批</span>
          </a-menu-item>
          <a-menu-item key="/mysqlReadPermDetail/list">
            <UnorderedListOutlined />
            <span>读权限列表</span>
          </a-menu-item>
          <a-menu-item key="/mysqlReadPermManage/list" v-if="userStore.isDba">
            <BlockOutlined />
            <span>读权限管理</span>
          </a-menu-item>
          <a-menu-item key="/mysqlDataUpdateApply/list">
            <BookOutlined />
            <span>数据修改单</span>
          </a-menu-item>
          <a-menu-item key="/mysqlDataUpdateAudit/list" v-if="userStore.isDba">
            <AuditOutlined />
            <span>数据修改单审批</span>
          </a-menu-item>
          <a-menu-item key="/mysqlSearch">
            <FileSearchOutlined />
            <span>数据查询</span>
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
  DatabaseOutlined,
  UnorderedListOutlined,
  BookOutlined,
  FileSearchOutlined,
  AuditOutlined,
  BlockOutlined,
  CarryOutOutlined
} from "@ant-design/icons-vue";
import { useUserStore } from "@/pinia/userStore";
const userStore = useUserStore();
const route = useRoute();
const router = useRouter();
const selectedKeys = ref([]);
const routeKey = "/db";
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
  "/mysqlDb": "/mysqlDb/list",
  "/mysqlReadPermApply": "/mysqlReadPermApply/list",
  "/mysqlReadPermAudit": "/mysqlReadPermAudit/list",
  "/mysqlReadPermDetail": "/mysqlReadPermDetail/list",
  "/mysqlSearch": "/mysqlSearch",
  "/mysqlDataUpdateApply": "/mysqlDataUpdateApply/list",
  "/mysqlDataUpdateAudit": "/mysqlDataUpdateAudit/list",
  "/mysqlReadPermManage": "/mysqlReadPermManage/list"
};
changeSelectedKey(route.path);
watch(
  () => router.currentRoute.value.path,
  newPath => changeSelectedKey(newPath)
);
</script>
<style scoped>
</style>