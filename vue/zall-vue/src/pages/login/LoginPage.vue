<template>
  <div :class="{'section': true, 'section-wework': loginType === WEWORK_TYPE}">
    <div v-show="loginType === ACCOUNT_PASSWORD_TYPE">
      <div class="title">{{t("login.title")}}</div>
      <div class="text-input">
        <a-input v-model:value="account" :placeholder="t('login.accountPlaceholder')" allow-clear>
          <template #prefix>
            <user-outlined />
          </template>
        </a-input>
      </div>
      <div class="text-input">
        <a-input-password
          v-model:value="password"
          :placeholder="t('login.passwordPlaceholder')"
          @pressEnter="login"
        >
          <template #prefix>
            <key-outlined />
          </template>
        </a-input-password>
      </div>
      <div class="submit-btn">
        <a-button type="primary" style="width:100%" @click="login">{{t("login.loginBtnText")}}</a-button>
      </div>
    </div>
    <div id="wework-login" v-show="loginType === WEWORK_TYPE"></div>
    <ul class="btn-ul">
      <li v-if="hasAccountPasswordType && loginType !== ACCOUNT_PASSWORD_TYPE">
        <a-button style="width:100%" @click="useAccountPassword">{{t("login.loginWithAccountPassword")}}</a-button>
      </li>
      <li v-if="hasWeworkType && loginType !== WEWORK_TYPE">
        <a-button style="width:100%" @click="useWework">{{t("login.loginWithWework")}}</a-button>
      </li>
      <li v-if="hasFeishuType">
        <a-button style="width:100%" @click="useFeishu">{{t("login.loginWithFeishu")}}</a-button>
      </li>
    </ul>
    <div class="sub-section" v-if="allowUserRegister">
      <span class="sub-text" @click="goToRegister">{{t("login.registerText")}}</span>
    </div>
  </div>
</template>
<script setup>
import { UserOutlined, KeyOutlined } from "@ant-design/icons-vue";
import { ref, nextTick } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter, useRoute } from "vue-router";
import { loginRequest, weworkLoginRequest } from "@/api/user/loginApi";
import { getSysCfgRequest, getLoginCfgRequest } from "@/api/cfg/cfgApi";
import { message } from "ant-design-vue";
import { useUserStore } from "@/pinia/userStore";
import { accountRegexp, passwordRegexp } from "@/utils/regexp";
import * as ww from "@wecom/jssdk";
import { setLoginUser } from "@/utils/login";
const { locale } = useI18n();
const ACCOUNT_PASSWORD_TYPE = "accountPassword";
const WEWORK_TYPE = "wework";
const FEISHU_TYPE = "feishu";
const user = useUserStore();
const account = ref("");
const password = ref("");
const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const allowUserRegister = ref(false);
const loginCfg = ref({});
const loginType = ref("");
const hasAccountPasswordType = ref(false);
const hasWeworkType = ref(false);
const hasFeishuType = ref(false);
const wwPanel = ref(null);
// 跳转注册页面
const goToRegister = () => router.push("/login/register");
// 登录请求
const login = () => {
  // 验证账号
  if (!accountRegexp.test(account.value)) {
    message.error(t("login.pleaseConfirmAccount"));
    return;
  }
  // 验证密码
  if (!passwordRegexp.test(password.value)) {
    message.error(t("login.pleaseConfirmPassword"));
    return;
  }
  loginRequest({
    account: account.value,
    password: password.value,
    a: route.query.a ? route.query.a : ""
  }).then(res => {
    user.account = res.session.userInfo.account;
    user.avatarUrl = res.session.userInfo.avatarUrl;
    user.email = res.session.userInfo.email;
    user.isAdmin = res.session.userInfo.isAdmin;
    user.name = res.session.userInfo.name;
    user.isDba = res.session.userInfo.isDba;
    setLoginUser({
      ...res.session.userInfo,
      sessionExpireAt: res.session.expireAt,
      sessionId: res.session.sessionId
    });
    if (route.query && route.query.redirect_uri) {
      window.location.href = decodeURI(route.query.redirect_uri);
    } else {
      router.push("/index");
    }
  });
};
// 获取系统配置 检查是否允许用户注册
const getSysCfg = () => {
  getSysCfgRequest().then(res => {
    allowUserRegister.value = !res.data.disableSelfRegisterUser;
  });
};
// 获取登录配置
const getLoginCfg = () => {
  getLoginCfgRequest().then(res => {
    let hasA = route.query.a === "zsf";
    loginCfg.value = res.data;
    hasAccountPasswordType.value = res.data.accountPassword.isEnabled || hasA;
    hasWeworkType.value = res.data.wework.isEnabled;
    hasFeishuType.value = res.data.feishu.isEnabled;
    // 留给后门用于强制账号密码登录
    if (hasA || hasAccountPasswordType.value) {
      loginType.value = ACCOUNT_PASSWORD_TYPE;
    } else if (hasWeworkType.value) {
      loginType.value = WEWORK_TYPE;
      useWework();
    } else if (hasFeishuType.value) {
      loginType.value = FEISHU_TYPE;
    }
  });
};
// 使用账号密码登录
const useAccountPassword = () => {
  loginType.value = ACCOUNT_PASSWORD_TYPE;
};
// 使用企业微信登录
const useWework = () => {
  loginType.value = WEWORK_TYPE;
  let cfg = loginCfg.value?.wework;
  nextTick(() => {
    createWeworkLoginPanel(
      cfg?.appId,
      cfg?.agentId,
      cfg?.redirectUrl,
      cfg?.state,
      locale.value
    );
  });
};
// 使用飞书登录
const useFeishu = () => {
  let cfg = loginCfg.value.feishu;
  window.location.href = `https://passport.feishu.cn/suite/passport/oauth/authorize?client_id=${
    cfg?.clientId
  }&redirect_uri=${encodeURI(cfg?.redirectUrl)}&response_type=code&state=${
    cfg?.state
  }&scope=contact:user.employee_id:readonly`;
};
// 企微登录框
const createWeworkLoginPanel = (appId, agentId, redirectUrl, state, lang) => {
  if (wwPanel.value) {
    wwPanel.value.unmount();
  }
  wwPanel.value = ww.createWWLoginPanel({
    el: "#wework-login",
    params: {
      login_type: "CorpApp",
      appid: appId,
      agentid: agentId,
      redirect_uri: redirectUrl,
      state: state,
      redirect_type: "callback",
      panel_size: "small",
      lang: lang
    },
    onLoginSuccess({ code }) {
      weworkLoginRequest({ code, state }).then(res => {
        user.account = res.session.userInfo.account;
        user.avatarUrl = res.session.userInfo.avatarUrl;
        user.email = res.session.userInfo.email;
        user.isAdmin = res.session.userInfo.isAdmin;
        user.name = res.session.userInfo.name;
        user.isDba = res.session.userInfo.isDba;
        user.sessionId = res.session.sessionId;
        user.sessionExpireAt = res.session.expireAt;
        setLoginUser({
          ...res.session.userInfo,
          sessionExpireAt: res.session.expireAt,
          sessionId: res.session.sessionId
        });
        if (route.query && route.query.redirect_uri) {
          window.location.href = decodeURI(route.query.redirect_uri);
        } else {
          router.push("/index");
        }
      });
      createWeworkLoginPanel(appId, agentId, redirectUrl, state);
    }
  });
};
getSysCfg();
getLoginCfg();
</script>
<style scoped>
.section {
  padding: 18px;
  width: 356px;
  overflow: hidden;
  margin-top: calc(50vh - 64px);
  margin-left: auto;
  margin-right: auto;
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
.sub-section {
  display: flex;
  justify-content: flex-end;
  margin-top: 4px;
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
.btn-ul {
  margin-top: 14px;
}
.btn-ul > li + li {
  margin-top: 14px;
}
</style>