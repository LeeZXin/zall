<template>
  <div class="avatar-name" :style="props.style" @click="showDrawer">
    <img v-if="showAvatar" :src="avatar" />
    <div v-if="!showAvatar" class="avatar-fake">草</div>
    <span>{{username}}</span>
  </div>
  <a-drawer
    :title="username"
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
    <div class="loginout-text">{{t("loginoutText")}}</div>
  </a-drawer>
</template>
<script setup>
import { useUserStore } from "@/pinia/UserStore";
import { defineProps, ref } from "vue";
import { useI18n } from "vue-i18n";
import { UserOutlined, KeyOutlined } from "@ant-design/icons-vue";
const bodyStyle = {
    Padding: "10px 20px"
}
const { t } = useI18n();
const user = useUserStore();
const username = user.name;
const avatar = user.avatar;
const showAvatar = avatar ? true : false;
const props = defineProps(["style"]);
const visible = ref(false);
const showDrawer = () => {
  visible.value = true;
};
const closeDrawer = () => {
  visible.value = false;
};
</script>
<style scoped>
.avatar-name {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: white;
  cursor: pointer;
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
.loginout-text {
  height: 32px;
  line-height: 32px;
  font-size: 14px;
  color: darkred;
  cursor: pointer;
}
.loginout-text:hover {
  color: red;
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