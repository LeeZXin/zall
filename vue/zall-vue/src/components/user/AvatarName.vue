<template>
  <div class="avatar-name" :style="props.style" @click="showDrawer">
    <img v-if="user.avatar" :src="user.avatar" />
    <div v-if="!user.avatar" class="avatar-fake">{{user.name?.substring(0,1)}}</div>
    <span>{{user.name}}</span>
  </div>
  <a-drawer
    :title="user.name"
    :closable="false"
    :open="visible"
    :bodyStyle="bodyStyle"
    @close="closeDrawer"
  >
    <div class="drawer-item">
      <key-outlined style="font-size:18px" />
      <span>{{t("settings.sshAndGpg")}}</span>
    </div>
    <div class="drawer-item">
      <user-outlined style="font-size:18px" />
      <span>cnm</span>
    </div>
    <div class="drawer-item">
      <user-outlined style="font-size:18px" />
      <span>cnm</span>
    </div>
    <div class="drawer-item">
      <user-outlined style="font-size:18px" />
      <span>cnm</span>
    </div>
    <div class="drawer-item">
      <user-outlined style="font-size:18px" />
      <span>cnm</span>
    </div>
    <div class="drawer-item">
      <user-outlined style="font-size:18px" />
      <span>cnm</span>
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
import { UserOutlined, KeyOutlined } from "@ant-design/icons-vue";
import { useRouter } from "vue-router";
import { logoutRequest } from "@/api/user/loginApi";
const bodyStyle = {
  Padding: "10px 20px"
};
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
.avatar-fake {
  width: 42px;
  height: 42px;
  text-align: center;
  background-color: purple;
  color: white;
  font-size: 18px;
  line-height: 42px;
  border-radius: 50%;
  margin-right: 4px;
}
.select-item {
  line-height: 28px;
  width: 80px;
  text-align: center;
  cursor: pointer;
}
.select-item:hover {
  background-color: #f0f0f0;
}
.drawer-item {
  display: flex;
  align-items: center;
  width: 100%;
  overflow: hidden;
  font-size: 16px;
  cursor: pointer;
  height: 32px;
  line-height: 32px;
  border-radius: 4px;
}
.drawer-item:hover {
  background-color: #f0f0f0;
}
.drawer-item > span {
  margin-left: 6px;
}
</style>