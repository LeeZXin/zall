<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode==='create'">创建监控告警</span>
        <span v-else-if="mode==='update'">编辑监控告警</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">选择环境</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            placeholder="选择环境"
            v-model:value="formState.selectedEnv"
            :options="envList"
          />
          <div class="input-desc">多环境选择, 选择其中一个环境</div>
        </div>
      </div>
      <div class="section" v-if="mode==='update'">
        <div class="section-title">已选环境</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">名称</div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.name" />
          <div class="input-desc">描述监控告警作用, 长度为32以内</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">自动触发间隔</div>
        <div class="section-body">
          <a-input-number
            style="width:100%"
            :min="10"
            :max="3600"
            :step="10"
            v-model:value="formState.intervalSec"
            @change="limitIntervalSecInput"
          />
          <div class="input-desc">每隔一段时间自动触发检测, 10的倍数, 单位为秒, 最小为10, 最大为3600</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>触发类型</span>
          <span style="color:darkred">*</span>
        </div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.hookType">
            <a-radio :value="1">Webhook</a-radio>
            <a-radio :value="2">外部通知</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.hookType === 1">
        <div class="section-title">
          <span>Webhook</span>
        </div>
        <div class="section-body">
          <div>
            <div style="font-size: 12px;margin-bottom: 6px">hook url</div>
            <a-input style="width:100%" v-model:value="formState.hookUrl" placeholder="请填写" />
          </div>
          <div style="margin-top: 10px">
            <div style="font-size: 12px;margin-bottom: 6px">签名密钥</div>
            <a-input-password
              style="width:100%"
              v-model:value="formState.secret"
              placeholder="请填写"
            />
          </div>
        </div>
      </div>
      <div class="section" v-else-if="formState.hookType === 2">
        <div class="section-title">
          <span>外部通知模板</span>
        </div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.tplId"
            :options="tplList"
            show-search
            :filter-option="filterTplListOption"
            placeholder="请选择"
          />
        </div>
      </div>
      <div class="section">
        <div class="section-title">来源类型</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.sourceType">
            <a-radio :value="1">Mysql</a-radio>
            <a-radio :value="2">Prometheus</a-radio>
            <a-radio :value="3">Loki</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.sourceType === 1">
        <div class="section-title">Mysql配置</div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">Host</div>
              <a-input v-model:value="formState.mysqlHost" />
              <div class="input-desc">mysql主机ip ip:port格式</div>
            </li>
            <li>
              <div class="input-name">数据库</div>
              <a-input v-model:value="formState.mysqlDatabase" />
              <div class="input-desc">mysql数据库</div>
            </li>
            <li>
              <div class="input-name">账号</div>
              <a-input v-model:value="formState.mysqlUsername" />
              <div class="input-desc">访问mysql的账号</div>
            </li>
            <li>
              <div class="input-name">密码</div>
              <a-input-password v-model:value="formState.mysqlPassword" />
              <div class="input-desc">访问mysql的密码</div>
            </li>
            <li>
              <div class="input-name">Select sql</div>
              <a-input v-model:value="formState.mysqlSelectSql" />
              <div class="input-desc">执行查询的select的sql</div>
            </li>
            <li>
              <div class="input-name">判断条件</div>
              <a-input v-model:value="formState.mysqlCondition" />
              <div class="input-desc">查询结果满足判断条件则会触发告警</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="formState.sourceType === 2">
        <div class="section-title">Prometheus配置</div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">Host</div>
              <a-input v-model:value="formState.promHost" />
              <div class="input-desc">prometheus主机url http开头</div>
            </li>
            <li>
              <div class="input-name">promQl</div>
              <a-input v-model:value="formState.promQl" />
              <div class="input-desc">执行查询的promQl</div>
            </li>
            <li>
              <div class="input-name">判断条件</div>
              <a-input v-model:value="formState.promCondition" />
              <div class="input-desc">查询结果满足判断条件则会触发告警</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="formState.sourceType === 3">
        <div class="section-title">Loki配置</div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">loki host</div>
              <a-input v-model:value="formState.lokiHost" />
              <div class="input-desc">loki请求地址 http开头</div>
            </li>
            <li>
              <div class="input-name">loki OrgId</div>
              <a-input v-model:value="formState.lokiOrgId" />
              <div class="input-desc">loki的X-Scope-OrgID 非必填</div>
            </li>
            <li>
              <div class="input-name">logQl</div>
              <a-input v-model:value="formState.lokiLogQl" />
              <div class="input-desc">执行查询的logQl, 只支持matrix结果的logQl</div>
            </li>
            <li>
              <div class="input-name">过去多少时间</div>
              <a-input v-model:value="formState.lokiLastDuration" />
              <div class="input-desc">开始时间查询时间为过去多少时间, 例如1m, 代表从过去一分钟到现在的时间段, 不允许超过一个小时</div>
            </li>
            <li>
              <div class="input-name">步长step</div>
              <a-input-number
                style="width:100%"
                :min="1"
                :max="3600"
                v-model:value="formState.lokiStep"
              />
              <div class="input-desc">查询步长, 数字, 单位为秒, 最小为1, 最大为3600</div>
            </li>
            <li>
              <div class="input-name">判断条件</div>
              <a-input v-model:value="formState.lokiCondition" />
              <div class="input-desc">查询结果满足判断条件则会触发告警</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateTimerTask">立即保存</a-button>
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
  alertHttpHostRegexp
} from "@/utils/regexp";
import { useAlertConfigStore } from "@/pinia/alertConfigStore";
import {
  createAlertConfigRequest,
  updateAlertConfigRequest
} from "@/api/app/alertApi";
const alertConfigStore = useAlertConfigStore();
const router = useRouter();
const tplList = ref([]);
const formState = reactive({
  selectedEnv: null,
  name: "",
  intervalSec: 10,
  sourceType: 1,
  mysqlHost: "",
  mysqlDatabase: "",
  mysqlUsername: "",
  mysqlPassword: "",
  mysqlSelectSql: "",
  mysqlCondition: "",
  promHost: "",
  promQl: "",
  promCondition: "",
  lokiHost: "",
  lokiOrgId: "",
  lokiLogQl: "",
  lokiLastDuration: "",
  lokiStep: 60,
  lokiCondition: "",
  hookUrl: "",
  secret: "",
  tplId: null,
  hookType: 1
});
const envList = ref([]);
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
const saveOrUpdateTimerTask = () => {
  if (!formState.selectedEnv) {
    message.warn("请选择环境");
    return;
  }
  if (!alertConfigNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!formState.intervalSec) {
    message.warn("触发间隔错误");
    return;
  }
  let hookCfg = {};
  if (formState.hookType === 1) {
    if (!alertConfigHookUrlRegexp.test(formState.hookUrl)) {
      message.warn("hook url错误");
      return;
    }
    if (!alertConfigSecretRegexp.test(formState.secret)) {
      message.warn("hook签名密钥错误");
      return;
    }
    hookCfg.hookUrl = formState.hookUrl;
    hookCfg.secret = formState.secret;
  } else if (formState.hookType === 2) {
    if (!formState.tplId) {
      message.warn("请选择外部通知模板");
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
      message.warn("mysql host格式錯誤");
      return;
    }
    if (!formState.mysqlDatabase) {
      message.warn("请填写mysql数据库");
      return;
    }
    if (!formState.mysqlUsername) {
      message.warn("请填写mysql账号");
      return;
    }
    if (!formState.mysqlPassword) {
      message.warn("请填写mysql密码");
      return;
    }
    if (!formState.mysqlSelectSql) {
      message.warn("请填写mysql sql");
      return;
    }
    if (!formState.mysqlCondition) {
      message.warn("请填写mysql触发条件");
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
      message.warn("prometheus host格式錯誤");
      return;
    }
    if (!formState.promQl) {
      message.warn("请填写promQl");
      return;
    }
    if (!formState.promCondition) {
      message.warn("请填写prometheus触发条件");
      return;
    }
    alert.prom = {
      host: formState.promHost,
      promQl: formState.promQl,
      condition: formState.promCondition
    };
  } else if (formState.sourceType === 3) {
    if (!alertHttpHostRegexp.test(formState.lokiHost)) {
      message.warn("loki host格式錯誤");
      return;
    }
    if (!formState.lokiLogQl) {
      message.warn("请填写logQl");
      return;
    }
    if (!formState.lokiLastDuration) {
      message.warn("请填写loki过去多少时间");
      return;
    }
    if (!formState.lokiStep) {
      message.warn("请填写loki过去多少时间");
      return;
    }
    if (!formState.lokiCondition) {
      message.warn("请填写loki触发条件");
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
      message.success("创建成功");
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
      message.success("编辑成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/alertConfig/list/${formState.selectedEnv}`
      );
    });
  }
};
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
    formState.intervalSec = parseInt(formState.intervalSec / 10) * 10;
  }, 1000);
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
.headers-ul > li {
  height: 32px;
  line-height: 32px;
  width: 100%;
  display: flex;
  align-items: center;
  font-size: 14px;
}
.headers-ul > li + li {
  margin-top: 6px;
}
</style>