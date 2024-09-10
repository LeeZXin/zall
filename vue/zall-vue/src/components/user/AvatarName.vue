<template>
  <div class="avatar-name" :style="props.style" @click="showDrawer">
    <ZAvatar :url="user.avatarUrl" :name="user.name" size="medium" />
    <span style="margin-left: 8px">{{user.name}}</span>
  </div>
  <a-drawer :title="user.name" :closable="false" :open="visible" @close="closeDrawer">
    <div class="drawer-item no-wrap" @click="goto('/personalSetting/profile/info')">
      <SettingOutlined style="font-size:18px" />
      <span>{{t('personalSetting')}}</span>
    </div>
    <div class="drawer-item no-wrap" @click="goto('/sa/cfg/list')" v-if="user.isAdmin">
      <UserOutlined style="font-size:18px" />
      <span>{{t('superAdmin')}}</span>
    </div>
    <a-divider style="margin-bottom: 10px" />
    <div>
      <a-button type="primary" danger style="width:100%" @click="logout">{{t("logoutText")}}</a-button>
    </div>
  </a-drawer>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import { useUserStore } from "@/pinia/userStore";
import { defineProps, ref } from "vue";
import { useI18n } from "vue-i18n";
import { UserOutlined, SettingOutlined } from "@ant-design/icons-vue";
import { useRouter } from "vue-router";
import { logoutRequest } from "@/api/user/loginApi";
const router = useRouter();
const { t } = useI18n();
const user = useUserStore();
const props = defineProps(["style"]);
const visible = ref(false);
const showDrawer = () => {
  visible.value = true;
};
const closeDrawer = () => {
  visible.value = false;
};
const logout = () => {
  logoutRequest().then(() => {
    router.push("/login/login");
  });
};
const goto = url => {
  router.push(url);
};
</script>
<style scoped>
.avatar-name {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: white;
  cursor: pointer;
  height: 64px;
}
.drawer-item {
  display: flex;
  align-items: center;
  width: 100%;
  overflow: hidden;
  font-size: 16px;
  cursor: pointer;
  height: 42px;
  line-height: 42px;
  border-radius: 4px;
  padding: 0 10px;
}
.drawer-item:hover {
  color: #1677ff;
}
.drawer-item > span {
  margin-left: 6px;
}
</style>