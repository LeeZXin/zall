<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('discoverySource.createSource')}}</span>
        <span v-else-if="mode === 'update'">{{t('discoverySource.updateSource')}}</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">{{t('discoverySource.selectEnv')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.selectedEnv" :options="envList" />
        </div>
      </div>
      <div class="section" v-if="mode==='update'">
        <div class="section-title">{{t('discoverySource.selectedEnv')}}</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('discoverySource.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('discoverySource.endpoints')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.endpoints" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('discoverySource.etcdUsername')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.username" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('discoverySource.etcdPwd')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.password" type="password" />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateDiscoverySource">{{t('discoverySource.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { discoverySourceNameRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  createDiscoverySourceRequest,
  updateDiscoverySourceRequest
} from "@/api/app/discoveryApi";
import { useRoute, useRouter } from "vue-router";
import { useDiscoverySourceStore } from "@/pinia/discoverySourceStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const discoverySourceStore = useDiscoverySourceStore();
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
  name: "",
  endpoints: "",
  username: "",
  password: "",
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
// 新增或编辑来源
const saveOrUpdateDiscoverySource = () => {
  if (!discoverySourceNameRegexp.test(formState.name)) {
    message.warn(t("discoverySource.nameFormatErr"));
    return;
  }
  if (!formState.endpoints) {
    message.warn(t("discoverySource.pleaseFillEndpoints"));
    return;
  }
  let endpoints = formState.endpoints.split(";").filter(item => item);
  if (mode === "create") {
    createDiscoverySourceRequest({
      env: formState.selectedEnv,
      endpoints: endpoints,
      username: formState.username,
      password: formState.password,
      name: formState.name
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/sa/discoverySource/list/${formState.selectedEnv}`);
    });
  } else if (mode === "update") {
    updateDiscoverySourceRequest({
      sourceId: discoverySourceStore.id,
      endpoints: endpoints,
      username: formState.username,
      password: formState.password,
      name: formState.name
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/sa/discoverySource/list/${formState.selectedEnv}`);
    });
  }
};

if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (discoverySourceStore.id === 0) {
    router.push(`/sa/discoverySource/list`);
  } else {
    formState.name = discoverySourceStore.name;
    formState.selectedEnv = discoverySourceStore.env;
    formState.username = discoverySourceStore.username;
    formState.password = discoverySourceStore.password;
    formState.endpoints = discoverySourceStore.endpoints;
  }
}
</script>
<style scoped>
</style>