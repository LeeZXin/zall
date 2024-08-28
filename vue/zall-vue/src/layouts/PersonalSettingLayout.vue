<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{userStore.account}}</span>
      <span class="switch-text" @click="switchIndex">{{t("backToIndex")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider collapsible>
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="onselect">
          <a-menu-item key="/profile/info">
            <UserOutlined />
            <span>{{t("personalSettingMenu.profile")}}</span>
          </a-menu-item>
          <a-menu-item key="/changePassword">
            <LockOutlined />
            <span>{{t("personalSettingMenu.changePassword")}}</span>
          </a-menu-item>
          <a-menu-item key="/sshAndGpg/list">
            <KeyOutlined />
            <span>{{t("personalSettingMenu.sshAndGpg")}}</span>
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
import { LockOutlined, KeyOutlined, UserOutlined } from "@ant-design/icons-vue";
import { useUserStore } from "@/pinia/userStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const userStore = useUserStore();
const route = useRoute();
const router = useRouter();
const selectedKeys = ref([]);
const routeKey = "/personalSetting";
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
  "/profile": "/profile/info",
  "/changePassword": "/changePassword",
  "/sshAndGpg": "/sshAndGpg/list"
};
changeSelectedKey(route.path);
watch(
  () => router.currentRoute.value.path,
  newPath => changeSelectedKey(newPath)
);
</script>
<style scoped>
</style>