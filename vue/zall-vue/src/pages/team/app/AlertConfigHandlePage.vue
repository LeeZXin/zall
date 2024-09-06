<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode==='create'">{{t('alertConfig.createConfig')}}</span>
        <span v-else-if="mode==='update'">{{t('alertConfig.updateConfig')}}</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">{{t('alertConfig.selectEnv')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.selectedEnv" :options="envList" />
        </div>
      </div>
      <div class="section" v-if="mode==='update'">
        <div class="section-title">{{t('alertConfig.selectedEnv')}}</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('alertConfig.name')}}</div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('alertConfig.intervalSec')}}</div>
        <div class="section-body">
          <a-input-number
            style="width:100%"
            :min="10"
            :max="3600"
            :step="10"
            v-model:value="formState.intervalSec"
            @change="limitIntervalSecInput"
          />
          <div class="input-desc">{{t('alertConfig.intervalSecDesc')}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>{{t('alertConfig.hookType')}}</span>
        </div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.hookType">
            <a-radio :value="1">{{t('alertConfig.webhook')}}</a-radio>
            <a-radio :value="2">{{t('alertConfig.notification')}}</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.hookType === 1">
        <div class="section-title">{{t('alertConfig.webhook')}}</div>
        <div class="section-body">
          <div>
            <div style="font-size: 12px;margin-bottom: 6px">{{t('alertConfig.hookUrl')}}</div>
            <a-input style="width:100%" v-model:value="formState.hookUrl" />
          </div>
          <div style="margin-top: 10px">
            <div style="font-size: 12px;margin-bottom: 6px">{{t('alertConfig.hookSecret')}}</div>
            <a-input-password style="width:100%" v-model:value="formState.secret" />
          </div>
        </div>
      </div>
      <div class="section" v-else-if="formState.hookType === 2">
        <div class="section-title">{{t('alertConfig.notification')}}</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.tplId"
            :options="tplList"
            show-search
            :filter-option="filterTplListOption"
          />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('alertConfig.sourceType')}}</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.sourceType">
            <a-radio :value="1">{{t('alertConfig.mysql')}}</a-radio>
            <a-radio :value="2">{{t('alertConfig.prom')}}</a-radio>
            <a-radio :value="3">{{t('alertConfig.loki')}}</a-radio>
            <a-radio :value="4">{{t('alertConfig.http')}}</a-radio>
            <a-radio :value="5">{{t('alertConfig.tcp')}}</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.sourceType === 1">
        <div class="section-title">{{t('alertConfig.mysql')}}</div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">{{t('alertConfig.mysqlHost')}}</div>
              <a-input v-model:value="formState.mysqlHost" />
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.mysqlDatabase')}}</div>
              <a-input v-model:value="formState.mysqlDatabase" />
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.mysqlUsername')}}</div>
              <a-input v-model:value="formState.mysqlUsername" />
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.mysqlPassword')}}</div>
              <a-input-password v-model:value="formState.mysqlPassword" />
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.mysqlSelectSql')}}</div>
              <a-input v-model:value="formState.mysqlSelectSql" />
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.mysqlCondition')}}</div>
              <a-input v-model:value="formState.mysqlCondition" />
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="formState.sourceType === 2">
        <div class="section-title">{{t('alertConfig.prom')}}</div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">{{t('alertConfig.promHost')}}</div>
              <a-input v-model:value="formState.promHost" />
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.promQl')}}</div>
              <a-input v-model:value="formState.promQl" />
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.promCondition')}}</div>
              <a-input v-model:value="formState.promCondition" />
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="formState.sourceType === 3">
        <div class="section-title">{{t('alertConfig.loki')}}</div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">{{t('alertConfig.lokiHost')}}</div>
              <a-input v-model:value="formState.lokiHost" />
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.lokiOrgId')}}</div>
              <a-input v-model:value="formState.lokiOrgId" />
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.lokiLogQl')}}</div>
              <a-input v-model:value="formState.lokiLogQl" />
              <div class="input-desc">{{t('alertConfig.lokiLogQlDesc')}}</div>
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.lokiLastDuration')}}</div>
              <a-input v-model:value="formState.lokiLastDuration" />
              <div class="input-desc">{{t('alertConfig.lokiLastDurationDesc')}}</div>
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.lokiStep')}}</div>
              <a-input-number
                style="width:100%"
                :min="1"
                :max="3600"
                v-model:value="formState.lokiStep"
              />
              <div class="input-desc">{{t('alertConfig.lokiStepDesc')}}</div>
            </li>
            <li>
              <div class="input-name">{{t('alertConfig.lokiCondition')}}</div>
              <a-input v-model:value="formState.lokiCondition" />
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="formState.sourceType === 4">
        <div class="section-title">{{t('alertConfig.http')}}</div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">{{t('alertConfig.httpGetUrl')}}</div>
              <a-input v-model:value="formState.httpGetUrl" />
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="formState.sourceType === 5">
        <div class="section-title">{{t('alertConfig.tcp')}}</div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">{{t('alertConfig.tcpHost')}}</div>
              <a-input v-model:value="formState.tcpHost" />
            </li>
          </ul>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateAlertConfig">{{t('alertConfig.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { useRoute, useRouter } from "vue-router";
import { ref, reactive } from "vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import { listAllTplByTeamIdRequest } from "@/api/team/notifyApi";
import { message } from "ant-design-vue";
import {
  alertConfigNameRegexp,
  alertConfigHookUrlRegexp,
  alertConfigSecretRegexp,
  alertMysqlHostRegexp,
  alertHttpHostRegexp,
  alertIpPortHostRegexp
} from "@/utils/regexp";
import { useAlertConfigStore } from "@/pinia/alertConfigStore";
import {
  createAlertConfigRequest,
  updateAlertConfigRequest
} from "@/api/app/alertApi";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const alertConfigStore = useAlertConfigStore();
const router = useRouter();
// 通知列表
const tplList = ref([]);
// 表单数据
const formState = reactive({
  selectedEnv: null,
  name: "",
  intervalSec: 10,
  sourceType: 1,
  // mysql
  mysqlHost: "",
  mysqlDatabase: "",
  mysqlUsername: "",
  mysqlPassword: "",
  mysqlSelectSql: "",
  mysqlCondition: "",
  // prom
  promHost: "",
  promQl: "",
  promCondition: "",
  // loki
  lokiHost: "",
  lokiOrgId: "",
  lokiLogQl: "",
  lokiLastDuration: "",
  lokiStep: 60,
  lokiCondition: "",
  // hook
  hookUrl: "",
  secret: "",
  tplId: null,
  hookType: 1,
  // http
  httpGetUrl: "",
  // tcp
  tcpHost: ""
});
// 环境列表
const envList = ref([]);
// 模式
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const getEnvCfg = () => {
  getEnvCfgRequest().then(res => {
    envList.value = res.data.map(item => {
      return {
        value: item,
        label: item
      };
    });
    if (route.query.env && res.data?.includes(route.query.env)) {
      formState.selectedEnv = route.query.env;
    } else if (res.data.length > 0) {
      formState.selectedEnv = res.data[0];
    }
  });
};
// 新增或编辑告警配置
const saveOrUpdateAlertConfig = () => {
  if (!alertConfigNameRegexp.test(formState.name)) {
    message.warn(t("alertConfig.nameFormatErr"));
    return;
  }
  if (!formState.intervalSec) {
    message.warn(t("alertConfig.pleaseFillIntervalSec"));
    return;
  }
  let hookCfg = {};
  if (formState.hookType === 1) {
    if (!alertConfigHookUrlRegexp.test(formState.hookUrl)) {
      message.warn(t("alertConfig.hookUrlFormatErr"));
      return;
    }
    if (!alertConfigSecretRegexp.test(formState.secret)) {
      message.warn(t("alertConfig.hookSecretFormatErr"));
      return;
    }
    hookCfg.hookUrl = formState.hookUrl;
    hookCfg.secret = formState.secret;
  } else if (formState.hookType === 2) {
    if (!formState.tplId) {
      message.warn(t("alertConfig.pleaseSelectTplId"));
      return;
    }
    hookCfg.notifyTplId = formState.tplId;
  } else {
    // 不可能走这
    return;
  }
  let alert = {
    sourceType: formState.sourceType,
    hookType: formState.hookType,
    hookCfg
  };
  if (formState.sourceType === 1) {
    if (!alertMysqlHostRegexp.test(formState.mysqlHost)) {
      message.warn(t("alertConfig.mysqlHostFormatErr"));
      return;
    }
    if (!formState.mysqlDatabase) {
      message.warn(t("alertConfig.pleaseFillMysqlDatabase"));
      return;
    }
    if (!formState.mysqlUsername) {
      message.warn(t("alertConfig.pleaseFillMysqlUsername"));
      return;
    }
    if (!formState.mysqlPassword) {
      message.warn(t("alertConfig.pleaseFillMysqlPassword"));
      return;
    }
    if (!formState.mysqlSelectSql) {
      message.warn(t("alertConfig.pleaseFillMysqlSelectSql"));
      return;
    }
    if (!formState.mysqlCondition) {
      message.warn(t("alertConfig.pleaseFillMysqlCondition"));
      return;
    }
    alert.mysql = {
      host: formState.mysqlHost,
      database: formState.mysqlDatabase,
      username: formState.mysqlUsername,
      password: formState.mysqlPassword,
      selectSql: formState.mysqlSelectSql,
      condition: formState.mysqlCondition
    };
  } else if (formState.sourceType === 2) {
    if (!alertHttpHostRegexp.test(formState.promHost)) {
      message.warn(t("alertConfig.promHostFormatErr"));
      return;
    }
    if (!formState.promQl) {
      message.warn(t("alertConfig.pleaseFillPromQl"));
      return;
    }
    if (!formState.promCondition) {
      message.warn(t("alertConfig.pleaseFillPromCondition"));
      return;
    }
    alert.prom = {
      host: formState.promHost,
      promQl: formState.promQl,
      condition: formState.promCondition
    };
  } else if (formState.sourceType === 3) {
    if (!alertHttpHostRegexp.test(formState.lokiHost)) {
      message.warn(t("alertConfig.lokiHostFormatErr"));
      return;
    }
    if (!formState.lokiLogQl) {
      message.warn(t("alertConfig.pleaseFillLokiLogQl"));
      return;
    }
    if (!formState.lokiLastDuration) {
      message.warn(t("alertConfig.pleaseFillLokiLastDuration"));
      return;
    }
    if (!formState.lokiStep) {
      message.warn(t("alertConfig.pleaseFillLokiStep"));
      return;
    }
    if (!formState.lokiCondition) {
      message.warn(t("alertConfig.pleaseFillLokiCodition"));
      return;
    }
    alert.loki = {
      host: formState.lokiHost,
      orgId: formState.lokiOrgId,
      logQl: formState.lokiLogQl,
      lastDuration: formState.lokiLastDuration,
      step: formState.lokiStep,
      condition: formState.lokiCondition
    };
  } else if (formState.sourceType === 4) {
    if (!alertHttpHostRegexp.test(formState.httpGetUrl)) {
      message.warn(t("alertConfig.httpGetUrlFormatErr"));
      return;
    }
    alert.http = {
      getUrl: formState.httpGetUrl
    };
  } else if (formState.sourceType === 5) {
    if (!alertIpPortHostRegexp.test(formState.tcpHost)) {
      message.warn(t("alertConfig.tcpHostFormatErr"));
      return;
    }
    alert.tcp = {
      host: formState.tcpHost
    };
  } else {
    return;
  }
  if (mode === "create") {
    createAlertConfigRequest({
      env: formState.selectedEnv,
      appId: route.params.appId,
      name: formState.name,
      intervalSec: formState.intervalSec,
      alert
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/alertConfig/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "update") {
    updateAlertConfigRequest({
      name: formState.name,
      intervalSec: formState.intervalSec,
      alert,
      id: alertConfigStore.id
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/alertConfig/list/${formState.selectedEnv}`
      );
    });
  }
};
// 获取外部通知列表
const getTplList = () => {
  listAllTplByTeamIdRequest(route.params.teamId).then(res => {
    tplList.value = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
  });
};
// 下拉框过滤
const filterTplListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 限制intervalsec输入
const limitIntervalSecInput = () => {
  setTimeout(() => {
    if (formState.intervalSec > 9 && formState.intervalSec < 3601) {
      formState.intervalSec = parseInt(formState.intervalSec / 10) * 10;
    }
  }, 2000);
};
if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (alertConfigStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/alertConfig/list`
    );
  } else {
    formState.selectedEnv = alertConfigStore.env;
    formState.intervalSec = alertConfigStore.intervalSec;
    formState.name = alertConfigStore.name;
    formState.sourceType = alertConfigStore.content?.sourceType;
    if (formState.sourceType === 1) {
      formState.mysqlHost = alertConfigStore.content?.mysql?.host;
      formState.mysqlDatabase = alertConfigStore.content?.mysql?.database;
      formState.mysqlUsername = alertConfigStore.content?.mysql?.username;
      formState.mysqlPassword = alertConfigStore.content?.mysql?.password;
      formState.mysqlSelectSql = alertConfigStore.content?.mysql?.selectSql;
      formState.mysqlCondition = alertConfigStore.content?.mysql?.condition;
    } else if (formState.sourceType === 2) {
      formState.promHost = alertConfigStore.content?.prom?.host;
      formState.promQl = alertConfigStore.content?.prom?.promQl;
      formState.promCondition = alertConfigStore.content?.prom?.condition;
    } else if (formState.sourceType === 3) {
      formState.lokiHost = alertConfigStore.content?.loki?.host;
      formState.lokiOrgId = alertConfigStore.content?.loki?.orgId;
      formState.lokiLogQl = alertConfigStore.content?.loki?.logQl;
      formState.lokiLastDuration = alertConfigStore.content?.loki?.lastDuration;
      formState.lokiStep = alertConfigStore.content?.loki?.step;
      formState.lokiCondition = alertConfigStore.content?.loki?.condition;
    } else if (formState.sourceType === 4) {
      formState.httpGetUrl = alertConfigStore.content?.http?.getUrl;
    } else if (formState.sourceType === 5) {
      formState.tcpHost = alertConfigStore.content?.tcp?.host;
    }
    formState.hookType = alertConfigStore.content?.hookType;
    if (formState.hookType === 1) {
      formState.hookUrl = alertConfigStore.content?.hookCfg?.hookUrl;
      formState.secret = alertConfigStore.content?.hookCfg?.secret;
    } else if (formState.hookType === 2) {
      formState.tplId = alertConfigStore.content?.hookCfg?.notifyTplId;
    }
    formState.hookType = alertConfigStore.content?.hookType;
  }
}
getTplList();
</script>
<style scoped>
</style>