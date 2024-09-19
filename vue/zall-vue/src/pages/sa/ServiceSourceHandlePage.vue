<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('serviceSource.createSource')}}</span>
        <span v-else-if="mode === 'update'">{{t('serviceSource.updateSource')}}</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">{{t('serviceSource.selectEnv')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.selectedEnv" :options="envList" />
        </div>
      </div>
      <div class="section" v-if="mode==='update'">
        <div class="section-title">{{t('serviceSource.selectedEnv')}}</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('serviceSource.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('serviceSource.datasource')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.datasource" />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateServiceSource">{{t('serviceSource.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import {
  serviceSourceNameRegexp,
  serviceSourceDatasourceRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  createServiceSourceRequest,
  updateServiceSourceRequest
} from "@/api/app/serviceApi";
import { useRoute, useRouter } from "vue-router";
import { useServiceSourceStore } from "@/pinia/serviceSourceStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const serviceSourceStore = useServiceSourceStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
// 模式
const mode = getMode();
const formState = reactive({
  name: "",
  datasource: "",
  selectedEnv: ""
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
// 新增或编辑服务来源
const saveOrUpdateServiceSource = () => {
  if (!serviceSourceNameRegexp.test(formState.name)) {
    message.warn(t("serviceSource.nameFormatErr"));
    return;
  }
  if (!serviceSourceDatasourceRegexp.test(formState.datasource)) {
    message.warn(t("serviceSource.datasourceFormatErr"));
    return;
  }
  if (mode === "create") {
    createServiceSourceRequest({
      env: formState.selectedEnv,
      name: formState.name,
      datasource: formState.datasource
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/sa/serviceSource/list/${formState.selectedEnv}`);
    });
  } else if (mode === "update") {
    updateServiceSourceRequest({
      sourceId: serviceSourceStore.id,
      name: formState.name,
      datasource: formState.datasource
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/sa/serviceSource/list/${formState.selectedEnv}`);
    });
  }
};

if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (serviceSourceStore.id === 0) {
    router.push(`/sa/serviceSource/list`);
  } else {
    formState.name = serviceSourceStore.name;
    formState.selectedEnv = serviceSourceStore.env;
    formState.datasource = serviceSourceStore.datasource;
  }
}
</script>
<style scoped>
</style>