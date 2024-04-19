<template>
  <div class="section">
    <div class="title">{{t("login.title")}}</div>
    <div class="text-input">
      <a-input v-model:value="account" :placeholder="t('login.accountPlaceholder')" allow-clear>
        <template #prefix>
          <user-outlined />
        </template>
      </a-input>
    </div>
    <div class="text-input">
      <a-input-password v-model:value="password" :placeholder="t('login.passwordPlaceholder')">
        <template #prefix>
          <key-outlined />
        </template>
      </a-input-password>
    </div>
    <div class="submit-btn">
      <a-button type="primary" style="width:100%" @click="login">{{t("login.loginBtnText")}}</a-button>
    </div>
    <div class="sub-section">
      <span class="sub-text" @click="goToRegister">{{t("login.registerText")}}</span>
    </div>
  </div>
</template>
<script setup>
import { UserOutlined, KeyOutlined } from "@ant-design/icons-vue";
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import api from "@/api/user/LoginApi";
import { message } from "ant-design-vue";
import { useUserStore } from "@/pinia/UserStore";

const user = useUserStore();
const account = ref("");
const password = ref("");
const { t } = useI18n();
const router = useRouter();
// 跳转注册页面
const goToRegister = () => router.push("/login/register");
// 登录请求
const login = () => {
  let inputAccount = account.value;
  if (!inputAccount || inputAccount.length < 4 || inputAccount.length > 32) {
    message.error(t("login.pleaseConfirmAccount"));
    return;
  }
  let inputPassword = password.value;
  if (!inputPassword || inputPassword.length < 6) {
    message.error(t("login.pleaseConfirmPassword"));
    return;
  }
  api
    .Login({
      account: account.value,
      password: password.value
    })
    .then(res => {
      user.$patch(state => {
        state.account = res.user.account;
        state.avatarUrl = res.user.avatarUrl;
        state.email = res.user.email;
        state.isAdmin = res.user.isAdmin;
        state.isProhibited = res.user.isProhibited;
        state.name = res.user.name;
        state.roleType = res.user.roleType;
        state.sessionId = res.sessionId;
        state.sessionExpireAt = res.expireAt;
      });
    });
};
</script>
<style scoped>
.section {
  padding: 18px;
  width: 24%;
  overflow: hidden;
  margin-top: calc(50vh - 64px);
  margin-left: 38%;
  background-color: white;
  border-radius: 4px;
  box-shadow: 0 0 15px #2f2f2f;
  transform: translateY(-50%);
}
.title {
  font-size: 20px;
  color: black;
  margin-bottom: 24px;
}
.text-input {
  margin-bottom: 14px;
}
.submit-btn {
  margin-bottom: 14px;
}
.sub-section {
  display: flex;
  justify-content: flex-end;
}
.sub-text {
  font-size: 14px;
  color: #1677ff;
  line-height: 24px;
  cursor: pointer;
}
.sub-text:hover {
  color: #1677ff;
}
</style>