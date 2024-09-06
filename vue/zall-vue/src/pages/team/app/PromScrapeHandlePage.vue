<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('promScrape.createScrape')}}</span>
        <span v-else-if="mode === 'update'">{{t('promScrape.updateScrape')}}</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">{{t('promScrape.selectEnv')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.selectedEnv" :options="envList" />
        </div>
      </div>
      <div class="section" v-if="mode==='update'">
        <div class="section-title">{{t('promScrape.selectedEnv')}}</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('promScrape.endpoint')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.endpoint" />
          <div class="input-desc">{{t('promScrape.endpointDesc')}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('promScrape.targetType')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.targetType">
            <a-select-option
              v-for="item in targetTypeList"
              v-bind:key="item.value"
              :value="item.value"
            >{{t(item.label)}}</a-select-option>
          </a-select>
          <div class="input-desc">{{t('promScrape.targetTypeDesc')}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('promScrape.target')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.target" />
          <div class="input-desc">{{t('promScrape.targetDesc')}}</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateScrape">{{t('promScrape.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import {
  promScrapeEndpointRegexp,
  promScrapeTargetRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  createPromScrapeByTeamRequest,
  updatePromScrapeByTeamRequest
} from "@/api/app/promApi";
import { useRoute, useRouter } from "vue-router";
import { usePromScrapeStore } from "@/pinia/promScrapeStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const promScrapeStore = usePromScrapeStore();
const route = useRoute();
const router = useRouter();
// 模式
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
// 表单数据
const formState = reactive({
  endpoint: "",
  target: "",
  targetType: 2
});
// 环境列表
const envList = ref([]);
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
// 目标类型
const targetTypeList = [
  {
    value: 1,
    label: "promScrape.discoveryType"
  },
  {
    value: 2,
    label: "promScrape.hostType"
  }
];
// 新增或编辑任务
const saveOrUpdateScrape = () => {
  if (!promScrapeEndpointRegexp.test(formState.endpoint)) {
    message.warn(t("promScrape.endpointFormatErr"));
    return;
  }
  if (!promScrapeTargetRegexp.test(formState.target)) {
    message.warn(t("promScrape.targetFormatErr"));
    return;
  }
  if (mode === "create") {
    createPromScrapeByTeamRequest({
      env: formState.selectedEnv,
      appId: route.params.appId,
      endpoint: formState.endpoint,
      target: formState.target,
      targetType: formState.targetType
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/promScrape/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "update") {
    updatePromScrapeByTeamRequest({
      scrapeId: promScrapeStore.id,
      target: formState.target,
      targetType: formState.targetType,
      endpoint: formState.endpoint
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/promScrape/list/${formState.selectedEnv}`
      );
    });
  }
};
if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (promScrapeStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/promScrape/list`
    );
  } else {
    formState.endpoint = promScrapeStore.endpoint;
    formState.selectedEnv = promScrapeStore.env;
    formState.targetType = promScrapeStore.targetType;
    formState.target = promScrapeStore.target;
    formState.appId = promScrapeStore.appId;
  }
}
</script>
<style scoped>
</style>