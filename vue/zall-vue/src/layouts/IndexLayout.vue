<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{t("app")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider collapsible>
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="onselect">
          <a-menu-item key="/index/team/list">
            <TeamOutlined />
            <span>团队协作</span>
          </a-menu-item>
          <a-menu-item key="/db/mysqlReadPermApply/list">
            <DatabaseOutlined />
            <span>Mysql审计</span>
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
import { useRouter, useRoute } from "vue-router";
import { ref, watch } from "vue";
import { TeamOutlined, DatabaseOutlined } from "@ant-design/icons-vue";
const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const selectedKeys = ref([]);
const routeKey = "/index";
const onselect = event => {
  router.push({
    path: event.key,
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
  "/team": "/index/team/list"
};
changeSelectedKey(route.path);
watch(
  () => router.currentRoute.value.path,
  newPath => changeSelectedKey(newPath)
);
</script>
<style scoped>
</style>