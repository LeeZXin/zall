<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{t("indexMenu.mysqlAudit")}}</span>
      <span class="switch-text" @click="switchIndex">{{t("backToIndex")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider collapsible>
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="onselect">
          <a-menu-item key="/mysqlDb/list" v-if="userStore.isDba">
            <DatabaseOutlined />
            <span>{{t("mysqlAuditMenu.databaseSource")}}</span>
          </a-menu-item>
          <a-menu-item key="/mysqlReadPermApply/list">
            <CarryOutOutlined />
            <span>{{t("mysqlAuditMenu.readPermApply")}}</span>
          </a-menu-item>
          <a-menu-item key="/mysqlReadPermAudit/list" v-if="userStore.isDba">
            <AuditOutlined />
            <span>{{t("mysqlAuditMenu.readPermAudit")}}</span>
          </a-menu-item>
          <a-menu-item key="/mysqlReadPermDetail/list">
            <UnorderedListOutlined />
            <span>{{t("mysqlAuditMenu.readPermList")}}</span>
          </a-menu-item>
          <a-menu-item key="/mysqlReadPermManage/list" v-if="userStore.isDba">
            <BlockOutlined />
            <span>{{t("mysqlAuditMenu.readPermManage")}}</span>
          </a-menu-item>
          <a-menu-item key="/mysqlDataUpdateApply/list">
            <BookOutlined />
            <span>{{t("mysqlAuditMenu.dataUpdateApply")}}</span>
          </a-menu-item>
          <a-menu-item key="/mysqlDataUpdateAudit/list" v-if="userStore.isDba">
            <AuditOutlined />
            <span>{{t("mysqlAuditMenu.dataUpdateAudit")}}</span>
          </a-menu-item>
          <a-menu-item key="/mysqlSearch">
            <FileSearchOutlined />
            <span>{{t("mysqlAuditMenu.dataSearch")}}</span>
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
import { useI18n } from "vue-i18n";
import { useUserStore } from "@/pinia/userStore";
const { t } = useI18n();
const userStore = useUserStore();
const route = useRoute();
const router = useRouter();
const selectedKeys = ref([]);
// 路由前缀
const routeKey = "/db";
// 导航栏选择key
const onselect = event => {
  router.push({
    path: routeKey + event.key,
    force: true
  });
};
// 返回首页
const switchIndex = () => {
  router.push("/index");
};
// 选择key变化
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