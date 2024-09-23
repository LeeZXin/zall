<template>
  <div style="padding:10px">
    <div class="container">
      <div class="section">
        <div class="section-title">{{t('cfgList.systemCfg')}}</div>
        <div class="section-body">
          <div>
            <a-checkbox
              v-model:checked="sysCfgState.disableSelfRegisterUser"
            >{{t('cfgList.disallowUserRegisterAccount')}}</a-checkbox>
            <div class="checkbox-desc">{{t('cfgList.disallowUserRegisterAccountDesc')}}</div>
          </div>
          <div style="margin-top: 10px">
            <a-checkbox
              v-model:checked="sysCfgState.allowUserCreateTeam"
            >{{t('cfgList.allowUserCreateTeam')}}</a-checkbox>
            <div class="checkbox-desc">{{t('cfgList.allowUserCreateTeamDesc')}}</div>
          </div>
          <div style="margin-top: 10px">
            <a-button type="primary" @click="updateSysCfg">{{t('cfgList.save')}}</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('cfgList.envCfg')}}</div>
        <div class="section-body">
          <div>
            <a-input v-model:value="envCfgState.envs" />
          </div>
          <div class="input-desc">{{t('cfgList.envCfgDesc')}}</div>
          <div style="margin-top: 10px">
            <a-button type="primary" @click="updateEnvCfg">{{t('cfgList.save')}}</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('cfgList.loginCfg')}}</div>
        <div class="section-body">
          <div>
            <a-checkbox
              v-model:checked="loginCfgState.allowAccountPassword"
            >{{t('cfgList.loginWithAccountPassword')}}</a-checkbox>
          </div>
          <div style="margin-top:10px;">
            <a-checkbox v-model:checked="loginCfgState.allowWework">{{t('cfgList.loginWithWework')}}</a-checkbox>
            <ul class="login-inputs">
              <li>
                <div class="login-input-title">appId</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.weworkAppId"
                    :disabled="!loginCfgState.allowWework"
                  />
                </div>
              </li>
              <li>
                <div class="login-input-title">agentId</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.weworkAgentId"
                    :disabled="!loginCfgState.allowWework"
                  />
                </div>
              </li>
              <li>
                <div class="login-input-title">secret</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.weworkSecret"
                    :disabled="!loginCfgState.allowWework"
                  />
                </div>
              </li>
              <li>
                <div class="login-input-title">state</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.weworkState"
                    :disabled="!loginCfgState.allowWework"
                  />
                </div>
              </li>
              <li>
                <div class="login-input-title">redirect_url</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.weworkRedirectUrl"
                    :disabled="!loginCfgState.allowWework"
                  />
                </div>
              </li>
            </ul>
          </div>
          <div style="margin-top:10px;">
            <a-checkbox v-model:checked="loginCfgState.allowFeishu">{{t('cfgList.loginWithFeishu')}}</a-checkbox>
            <ul class="login-inputs">
              <li>
                <div class="login-input-title">clientId</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.feishuClientId"
                    :disabled="!loginCfgState.allowFeishu"
                  />
                </div>
              </li>
              <li>
                <div class="login-input-title">secret</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.feishuSecret"
                    :disabled="!loginCfgState.allowFeishu"
                  />
                </div>
              </li>
              <li>
                <div class="login-input-title">state</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.feishuState"
                    :disabled="!loginCfgState.allowFeishu"
                  />
                </div>
              </li>
              <li>
                <div class="login-input-title">redirect_url</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.feishuRedirectUrl"
                    :disabled="!loginCfgState.allowFeishu"
                  />
                </div>
              </li>
            </ul>
          </div>
          <div style="margin-top: 10px">
            <a-button type="primary" @click="updateLoginCfg">{{t('cfgList.save')}}</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('cfgList.gitCfg')}}</div>
        <div class="section-body">
          <div>
            <div class="input-title">HTTP URL</div>
            <a-input v-model:value="gitCfgState.httpUrl" />
          </div>
          <div>
            <div class="input-title" style="margin-top: 10px">SSH URL</div>
            <a-input v-model:value="gitCfgState.sshUrl" />
          </div>
          <div>
            <div class="input-title" style="margin-top: 10px">{{t('cfgList.lfsJwtExpiredDuration')}}</div>
            <a-input-number v-model:value="gitCfgState.lfsJwtExpiry" style="width:100%" />
          </div>
          <div>
            <div class="input-title" style="margin-top: 10px">{{t('cfgList.lfsJwtKey')}}</div>
            <a-input v-model:value="gitCfgState.lfsJwtSecret" />
          </div>
          <div style="margin-top: 10px">
            <a-button type="primary" @click="updateGitCfg">{{t('cfgList.save')}}</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('cfgList.gitServerCfg')}}</div>
        <div class="section-body">
          <div>
            <div class="input-title">HTTP HOST</div>
            <a-input v-model:value="gitRepoServerCfgState.httpHost" />
          </div>
          <div>
            <div class="input-title" style="margin-top: 10px">SSH HOST</div>
            <a-input v-model:value="gitRepoServerCfgState.sshHost" />
          </div>
          <div style="margin-top: 10px">
            <a-button type="primary" @click="updateGitRepoServerCfg">{{t('cfgList.save')}}</a-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {
  getSysCfgRequest,
  updateSysCfgRequest,
  getEnvCfgRequest,
  updateEnvCfgRequest,
  getGitRepoServerCfgRequest,
  updateGitRepoServerCfgRequest,
  getGitCfgRequest,
  updateGitCfgRequest,
  getLoginCfgBySaRequest,
  updateLoginCfgRequest
} from "@/api/cfg/cfgApi";
import { reactive } from "vue";
import { message } from "ant-design-vue";
import { gitRepoServerHostRegexp, envRegexp } from "@/utils/regexp";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
// 系统配置
const sysCfgState = reactive({
  disableSelfRegisterUser: false,
  allowUserCreateTeam: false
});
// 环境配置
const envCfgState = reactive({
  envs: ""
});
// git url配置
const gitCfgState = reactive({
  httpUrl: "",
  sshUrl: "",
  lfsJwtExpiry: 600,
  lfsJwtSecret: ""
});
// 登录配置
const loginCfgState = reactive({
  allowAccountPassword: false,
  allowWework: false,
  allowFeishu: false,
  weworkAppId: "",
  weworkAgentId: "",
  weworkSecret: "",
  weworkState: "",
  weworkRedirectUrl: "",
  feishuClientId: "",
  feishuSecret: "",
  feishuState: "",
  feishuRedirectUrl: ""
});
// 仓库位置配置
const gitRepoServerCfgState = reactive({
  httpHost: "",
  sshHost: ""
});
// 获取系统配置
const getSysCfg = () => {
  getSysCfgRequest().then(res => {
    sysCfgState.disableSelfRegisterUser = res.data.disableSelfRegisterUser;
    sysCfgState.allowUserCreateTeam = res.data.allowUserCreateTeam;
  });
};
// 编辑系统配置
const updateSysCfg = () => {
  updateSysCfgRequest({
    disableSelfRegisterUser: sysCfgState.disableSelfRegisterUser,
    allowUserCreateTeam: sysCfgState.allowUserCreateTeam
  }).then(() => {
    message.success(t("operationSuccess"));
  });
};
// 环境配置
const getEnvCfg = () => {
  getEnvCfgRequest().then(res => {
    envCfgState.envs = res.data.join(";");
  });
};
// 编辑环境配置
const updateEnvCfg = () => {
  let envs = envCfgState.envs.split(";").filter(item => item);
  if (envs.length === 0) {
    message.warn(t("cfgList.pleaseFillEnv"));
    return;
  }
  for (let i in envs) {
    if (!envRegexp.test(envs[i])) {
      message.warn(t("cfgList.envFormatErr"));
      return;
    }
  }
  updateEnvCfgRequest({
    envs: envs
  }).then(() => {
    message.success(t("operationSuccess"));
  });
};
// git服务端配置
const getGitRepoServerCfg = () => {
  getGitRepoServerCfgRequest().then(res => {
    gitRepoServerCfgState.httpHost = res.data.httpHost;
    gitRepoServerCfgState.sshHost = res.data.sshHost;
  });
};
// 编辑git服务端配置
const updateGitRepoServerCfg = () => {
  if (!gitRepoServerHostRegexp.test(gitRepoServerCfgState.httpHost)) {
    message.warn(t("cfgList.httpHostFormatErr"));
    return;
  }
  if (!gitRepoServerHostRegexp.test(gitRepoServerCfgState.sshHost)) {
    message.warn(t("cfgList.sshHostFormatErr"));
    return;
  }
  updateGitRepoServerCfgRequest({
    httpHost: gitRepoServerCfgState.httpHost,
    sshHost: gitRepoServerCfgState.sshHost
  }).then(() => {
    message.success(t("operationSuccess"));
  });
};
// git配置
const getGitCfg = () => {
  getGitCfgRequest().then(res => {
    gitCfgState.httpUrl = res.data.httpUrl;
    gitCfgState.sshUrl = res.data.sshUrl;
    gitCfgState.lfsJwtExpiry = res.data.lfsJwtExpiry;
    gitCfgState.lfsJwtSecret = res.data.lfsJwtSecret;
  });
};
// 编辑git配置
const updateGitCfg = () => {
  updateGitCfgRequest({
    httpUrl: gitCfgState.httpUrl,
    sshUrl: gitCfgState.sshUrl,
    lfsJwtExpiry: gitCfgState.lfsJwtExpiry,
    lfsJwtSecret: gitCfgState.lfsJwtSecret
  }).then(() => {
    message.success(t("operationSuccess"));
  });
};
// 获取登录配置
const getLoginCfg = () => {
  getLoginCfgBySaRequest().then(res => {
    loginCfgState.allowAccountPassword = res.data.accountPassword.isEnabled;
    loginCfgState.allowWework = res.data.wework.isEnabled;
    loginCfgState.allowFeishu = res.data.feishu.isEnabled;
    // wework
    loginCfgState.weworkAppId = res.data.wework.appId;
    loginCfgState.weworkAgentId = res.data.wework.agentId;
    loginCfgState.weworkState = res.data.wework.state;
    loginCfgState.weworkSecret = res.data.wework.secret;
    loginCfgState.weworkRedirectUrl = res.data.wework.redirectUrl;
    // feishu
    loginCfgState.feishuClientId = res.data.feishu.clientId;
    loginCfgState.feishuSecret = res.data.feishu.secret;
    loginCfgState.feishuState = res.data.feishu.state;
    loginCfgState.feishuRedirectUrl = res.data.feishu.redirectUrl;
  });
};
// 编辑登录配置
const updateLoginCfg = () => {
  if (
    !loginCfgState.allowAccountPassword &&
    !loginCfgState.allowWework &&
    !loginCfgState.allowFeishu
  ) {
    message.warn(t("cfgList.pleaseChooseAtLeastOneLoginWay"));
    return;
  }
  let req = {};
  req.accountPassword = {
    isEnabled: loginCfgState.allowAccountPassword
  };
  req.wework = {
    appId: loginCfgState.weworkAppId,
    agentId: loginCfgState.weworkAgentId,
    state: loginCfgState.weworkState,
    secret: loginCfgState.weworkSecret,
    redirectUrl: loginCfgState.weworkRedirectUrl
  };
  req.feishu = {
    clientId: loginCfgState.feishuClientId,
    secret: loginCfgState.feishuSecret,
    state: loginCfgState.feishuState,
    redirectUrl: loginCfgState.feishuRedirectUrl
  };
  if (loginCfgState.allowWework) {
    if (!loginCfgState.weworkAppId) {
      message.warn(`${t("cfgList.pleaseFill")} appId`);
      return;
    }
    if (!loginCfgState.weworkAgentId) {
      message.warn(`${t("cfgList.pleaseFill")} agentId`);
      return;
    }
    if (!loginCfgState.weworkSecret) {
      message.warn(`${t("cfgList.pleaseFill")} secret`);
      return;
    }
    if (!loginCfgState.weworkState) {
      message.warn(`${t("cfgList.pleaseFill")} state`);
      return;
    }
    if (!loginCfgState.weworkRedirectUrl) {
      message.warn(`${t("cfgList.pleaseFill")} redirect_url`);
      return;
    }
    req.wework.isEnabled = true;
  }
  if (loginCfgState.allowFeishu) {
    if (!loginCfgState.feishuClientId) {
      message.warn(`${t("cfgList.pleaseFill")} clientId`);
      return;
    }
    if (!loginCfgState.feishuSecret) {
      message.warn(`${t("cfgList.pleaseFill")} secret`);
      return;
    }
    if (!loginCfgState.feishuState) {
      message.warn(`${t("cfgList.pleaseFill")} state`);
      return;
    }
    if (!loginCfgState.feishuRedirectUrl) {
      message.warn(`${t("cfgList.pleaseFill")} redirect_url`);
      return;
    }
    req.feishu.isEnabled = true;
  }
  updateLoginCfgRequest(req).then(() => {
    message.success(t("operationSuccess"));
  });
};
getSysCfg();
getEnvCfg();
getGitCfg();
getLoginCfg();
getGitRepoServerCfg();
</script>

<style scoped>
.login-inputs {
  padding-left: 24px;
  margin-top: 4px;
}
.login-inputs > li + li {
  margin-top: 4px;
}
.login-input-title {
  font-size: 12px;
  margin-bottom: 4px;
}
.login-input-title-gray {
  color: gray;
}
</style>