<template>
  <div class="section">
    <div class="title">{{t("register.title")}}</div>
    <div class="text-input">
      <a-input
        v-model:value="formState.account"
        :placeholder="t('register.accountPlaceholder')"
        allow-clear
      >
        <template #prefix>
          <user-outlined />
        </template>
      </a-input>
    </div>
    <div class="text-input">
      <a-input
        v-model:value="formState.name"
        :placeholder="t('register.usernamePlaceholder')"
        allow-clear
      >
        <template #prefix>
          <highlight-outlined />
        </template>
      </a-input>
    </div>
    <div class="text-input">
      <a-input
        v-model:value="formState.email"
        :placeholder="t('register.emailPlaceholder')"
        allow-clear
      >
        <template #prefix>
          <mail-outlined />
        </template>
      </a-input>
    </div>
    <div class="text-input">
      <a-input-password
        v-model:value="formState.password"
        :placeholder="t('register.passwordPlaceholder')"
      >
        <template #prefix>
          <key-outlined />
        </template>
      </a-input-password>
    </div>
    <div class="text-input">
      <a-input-password
        v-model:value="formState.confirmPassword"
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
import { reactive } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import { message } from "ant-design-vue";
import {
  accountRegexp,
  passwordRegexp,
  usernameRegexp,
  emailRegexp
} from "@/utils/regexp";
import { registerRequest, loginRequest } from "@/api/user/loginApi";
import { useUserStore } from "@/pinia/userStore";
const user = useUserStore();
const formState = reactive({
  account: "",
  name: "",
  email: "",
  password: "",
  confirmPassword: ""
});
const { t } = useI18n();
const router = useRouter();
const backToLogin = () => router.push("/login/login");
const register = () => {
  if (!accountRegexp.test(formState.account)) {
    message.error(t("register.pleaseConfirmAccount"));
    return;
  }
  if (!usernameRegexp.test(formState.name)) {
    message.error(t("register.pleaseConfirmUsername"));
    return;
  }
  if (!emailRegexp.test(formState.email)) {
    message.error(t("register.pleaseConfirmEmail"));
    return;
  }
  if (!passwordRegexp.test(formState.password)) {
    message.error(t("register.pleaseConfirmPassword"));
    return;
  }
  if (formState.confirmPassword !== formState.password) {
    message.error(t("register.pleaseConfirmConfirmPassword"));
    return;
  }
  registerRequest({
    account: formState.account,
    name: formState.name,
    password: formState.password,
    email: formState.email
  }).then(() => {
    loginRequest({
      account: formState.account,
      password: formState.password
    }).then(res => {
      user.account = res.session.userInfo.account;
      user.avatarUrl = res.session.userInfo.avatarUrl;
      user.email = res.session.userInfo.email;
      user.isAdmin = res.session.userInfo.isAdmin;
      user.name = res.session.userInfo.name;
      user.roleType = res.session.userInfo.roleType;
      user.sessionId = res.session.sessionId;
      user.sessionExpireAt = res.session.expireAt;
      router.push("/index");
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