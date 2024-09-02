<template>
  <div style="padding:10px">
    <div class="container">
      <div class="section">
        <div class="section-title">系统配置</div>
        <div class="section-body">
          <div>
            <a-checkbox v-model:checked="sysCfgState.disableSelfRegisterUser">禁止用户注册帐号</a-checkbox>
            <div class="checkbox-desc">用户不可自行注册帐号, 只允许超级管理员创建用户帐号</div>
          </div>
          <div style="margin-top: 10px">
            <a-checkbox v-model:checked="sysCfgState.allowUserCreateTeam">允许用户创建团队</a-checkbox>
            <div class="checkbox-desc">允许用户自行创建团队并成为团队的管理员</div>
          </div>
          <div style="margin-top: 10px">
            <a-button type="primary" @click="updateSysCfg">保存系统配置</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">环境配置</div>
        <div class="section-body">
          <div>
            <a-input v-model:value="envCfgState.envs" placeholder="请填写" />
          </div>
          <div class="input-desc">多个环境使用;隔开</div>
          <div style="margin-top: 10px">
            <a-button type="primary" @click="updateEnvCfg">保存环境配置</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">登录配置</div>
        <div class="section-body">
          <div>
            <a-checkbox v-model:checked="loginCfgState.allowAccountPassword">允许账号密码登录</a-checkbox>
            <div class="checkbox-desc">使用账号密码登录</div>
          </div>
          <div style="margin-top:10px;">
            <a-checkbox v-model:checked="loginCfgState.allowWework">允许企业微信登录</a-checkbox>
            <div class="checkbox-desc">使用企业微信扫码登录</div>
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
              <li>
                <div class="login-input-title">lang</div>
                <div>
                  <a-input
                    style="width:100%"
                    v-model:value="loginCfgState.weworkLang"
                    :disabled="!loginCfgState.allowWework"
                  />
                </div>
              </li>
            </ul>
          </div>
          <div style="margin-top:10px;">
            <a-checkbox v-model:checked="loginCfgState.allowFeishu">允许飞书登录</a-checkbox>
            <div class="checkbox-desc">使用飞书扫码登录</div>
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
            <a-button type="primary" @click="updateLoginCfg">保存登录配置</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">GIT配置</div>
        <div class="section-body">
          <div>
            <div class="input-title">HTTP URL</div>
            <a-input v-model:value="gitCfgState.httpUrl" placeholder="请填写" />
          </div>
          <div>
            <div class="input-title" style="margin-top: 10px">SSH URL</div>
            <a-input v-model:value="gitCfgState.sshUrl" placeholder="请填写" />
          </div>
          <div>
            <div class="input-title" style="margin-top: 10px">LFS JWT过期时间</div>
            <a-input-number
              v-model:value="gitCfgState.lfsJwtExpiry"
              placeholder="请填写"
              style="width:100%"
            />
            <div class="input-desc">单位为秒</div>
          </div>
          <div>
            <div class="input-title" style="margin-top: 10px">LFS JWT密钥</div>
            <a-input v-model:value="gitCfgState.lfsJwtSecret" placeholder="请填写" />
          </div>
          <div style="margin-top: 10px">
            <a-button type="primary" @click="updateGitCfg">保存GIT配置</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">GIT仓库服务端配置</div>
        <div class="section-body">
          <div>
            <div class="input-title">HTTP HOST</div>
            <a-input v-model:value="gitRepoServerCfgState.httpHost" placeholder="请填写" />
            <div class="input-desc">ip:port格式</div>
          </div>
          <div>
            <div class="input-title" style="margin-top: 10px">SSH HOST</div>
            <a-input v-model:value="gitRepoServerCfgState.sshHost" placeholder="请填写" />
            <div class="input-desc">ip:port格式</div>
          </div>
          <div style="margin-top: 10px">
            <a-button type="primary" @click="updateGitRepoServerCfg">保存GIT仓库服务端配置</a-button>
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
  weworkLang: "zh",
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
    message.success("编辑成功");
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
    message.warn("请填写环境");
    return;
  }
  for (let i in envs) {
    if (!envRegexp.test(envs[i])) {
      message.warn("环境格式错误");
      return;
    }
  }
  updateEnvCfgRequest({
    envs: envs
  }).then(() => {
    message.success("编辑成功");
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
    message.warn("http host格式错误");
    return;
  }
  if (!gitRepoServerHostRegexp.test(gitRepoServerCfgState.sshHost)) {
    message.warn("ssh host格式错误");
    return;
  }
  updateGitRepoServerCfgRequest({
    httpHost: gitRepoServerCfgState.httpHost,
    sshHost: gitRepoServerCfgState.sshHost
  }).then(() => {
    message.success("编辑成功");
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
    message.success("编辑成功");
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
    loginCfgState.weworkLang = res.data.wework.lang;
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
    message.warn("至少选择一种登录方式");
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
    redirectUrl: loginCfgState.weworkRedirectUrl,
    lang: loginCfgState.weworkLang
  };
  req.feishu = {
    clientId: loginCfgState.feishuClientId,
    secret: loginCfgState.feishuSecret,
    state: loginCfgState.feishuState,
    redirectUrl: loginCfgState.feishuRedirectUrl
  };
  if (loginCfgState.allowWework) {
    if (!loginCfgState.weworkAppId) {
      message.warn("请填写appId");
      return;
    }
    if (!loginCfgState.weworkAgentId) {
      message.warn("请填写agentId");
      return;
    }
    if (!loginCfgState.weworkSecret) {
      message.warn("请填写secret");
      return;
    }
    if (!loginCfgState.weworkState) {
      message.warn("请填写state");
      return;
    }
    if (!loginCfgState.weworkRedirectUrl) {
      message.warn("请填写redirect_url");
      return;
    }
    req.wework.isEnabled = true;
  }
  if (loginCfgState.allowFeishu) {
    if (!loginCfgState.feishuClientId) {
      message.warn("请填写clientId");
      return;
    }
    if (!loginCfgState.feishuSecret) {
      message.warn("请填写secret");
      return;
    }
    if (!loginCfgState.feishuState) {
      message.warn("请填写state");
      return;
    }
    if (!loginCfgState.feishuRedirectUrl) {
      message.warn("请填写redirect_url");
      return;
    }
    req.feishu.isEnabled = true;
  }
  updateLoginCfgRequest(req).then(() => {
    message.success("编辑成功");
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