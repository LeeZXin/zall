<template>
  <div class="section">
    <div class="title">{{t("register.title")}}</div>
    <div class="text-input">
      <a-input v-model:value="account" :placeholder="t('register.accountPlaceholder')" allow-clear>
        <template #prefix>
          <user-outlined />
        </template>
      </a-input>
    </div>
    <div class="text-input">
      <a-input
        v-model:value="username"
        :placeholder="t('register.usernamePlaceholder')"
        allow-clear
      >
        <template #prefix>
          <highlight-outlined />
        </template>
      </a-input>
    </div>
    <div class="text-input">
      <a-input v-model:value="email" :placeholder="t('register.emailPlaceholder')" allow-clear>
        <template #prefix>
          <mail-outlined />
        </template>
      </a-input>
    </div>
    <div class="text-input">
      <a-input-password v-model:value="password" :placeholder="t('register.passwordPlaceholder')">
        <template #prefix>
          <key-outlined />
        </template>
      </a-input-password>
    </div>
    <div class="text-input">
      <a-input-password
        v-model:value="confirmPassword"
        :placeholder="t('register.confirmPasswordPlaceholder')"
      >
        <template #prefix>
          <key-outlined />
        </template>
      </a-input-password>
    </div>
    <div class="submit-btn">
      <a-button
        type="primary"
        style="width:100%"
        @click="register"
      >{{t("register.registerBtnText")}}</a-button>
    </div>
    <div class="sub-section">
      <span class="sub-text" @click="backToLogin">{{t("register.backToLoginText")}}</span>
    </div>
  </div>
</template>
<script setup>
import {
  UserOutlined,
  KeyOutlined,
  MailOutlined,
  HighlightOutlined
} from "@ant-design/icons-vue";
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import { message } from "ant-design-vue";
const username = ref("");
const email = ref("");
const account = ref("");
const password = ref("");
const confirmPassword = ref("");
const { t } = useI18n();
const router = useRouter();
const backToLogin = () => router.push("/login/login");
const register = () => {
  let inputAccount = account.value;
  if (!inputAccount || inputAccount.length < 4 || inputAccount.length > 32) {
    message.error(t("register.pleaseConfirmAccount"));
    return;
  }
  let inputUsername = username.value;
  if (
    !inputUsername ||
    inputUsername.length > 32 ||
    inputUsername.length === 0
  ) {
    message.error(t("register.pleaseConfirmUsername"));
    return;
  }
  let inputEmail = email.value;
  if (!inputEmail || inputEmail.length === 0) {
    message.error(t("register.pleaseConfirmEmail"));
    return;
  }
  let inputPassword = password.value;
  if (!inputPassword || inputPassword.value < 6) {
    message.error(t("register.pleaseConfirmPassword"));
    return;
  }
  let inputConfirmPassword = confirmPassword.value;
  if (inputConfirmPassword !== inputPassword) {
    message.error(t("register.pleaseConfirmConfirmPassword"));
    return;
  }
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