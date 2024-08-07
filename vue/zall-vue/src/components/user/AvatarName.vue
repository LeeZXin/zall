<template>
  <div class="avatar-name" :style="props.style" @click="showDrawer">
    <img v-if="user.avatarUrl" :src="user.avatarUrl" class="avatar" />
    <div v-else class="avatar-fake">{{user.name?.substring(0,1)}}</div>
    <span>{{user.name}}</span>
  </div>
  <a-drawer :title="user.name" :closable="false" :open="visible" @close="closeDrawer">
    <div class="drawer-item no-wrap" @click="goto('/personalSetting/profile/info')">
      <SettingOutlined style="font-size:18px" />
      <span>个人设置</span>
    </div>
    <div class="drawer-item no-wrap" @click="goto('/sa/cfg/list')" v-if="user.isAdmin">
      <UserOutlined style="font-size:18px" />
      <span>超级管理员</span>
    </div>
    <a-divider style="margin-bottom: 10px" />
    <div>
      <a-button type="primary" danger style="width:100%" @click="logout">{{t("logoutText")}}</a-button>
    </div>
  </a-drawer>
</template>
<script setup>
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
.avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  margin-right: 8px;
  object-fit: contain;
}
.avatar-name {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: white;
  cursor: pointer;
  height: 64px;
}
.avatar-fake {
  width: 42px;
  height: 42px;
  text-align: center;
  background-color: purple;
  color: white;
  font-size: 18px;
  line-height: 42px;
  border-radius: 50%;
  margin-right: 8px;
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
  background-color: #f0f0f0;
}
.drawer-item > span {
  margin-left: 6px;
}
</style>