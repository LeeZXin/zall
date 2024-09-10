<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('propertySource.createSource')}}</span>
        <span v-else-if="mode === 'update'">{{t('propertySource.updateSource')}}</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">{{t('propertySource.selectEnv')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.selectedEnv" :options="envList" />
        </div>
      </div>
      <div class="section" v-if="mode==='update'">
        <div class="section-title">{{t('propertySource.selectedEnv')}}</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('propertySource.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('propertySource.endpoints')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.endpoints" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('propertySource.etcdUsername')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.username" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('propertySource.etcdPwd')}}</div>
        <div class="section-body">
          <a-input-password v-model:value="formState.password" />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdatePropertySource">{{t('propertySource.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { propertySourceNameRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  createPropertySourceRequest,
  updatePropertySourceRequest
} from "@/api/app/propertyApi";
import { useRoute, useRouter } from "vue-router";
import { usePropertySourceStore } from "@/pinia/propertySourceStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const propertySourceStore = usePropertySourceStore();
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
const saveOrUpdatePropertySource = () => {
  if (!propertySourceNameRegexp.test(formState.name)) {
    message.warn(t("propertySource.nameFormatErr"));
    return;
  }
  if (!formState.endpoints) {
    message.warn("propertySource.pleaseFillEndpoints");
    return;
  }
  let endpoints = formState.endpoints.split(";").filter(item => item);
  if (mode === "create") {
    createPropertySourceRequest({
      env: formState.selectedEnv,
      endpoints: endpoints,
      username: formState.username,
      password: formState.password,
      name: formState.name
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/sa/propertySource/list/${formState.selectedEnv}`);
    });
  } else if (mode === "update") {
    updatePropertySourceRequest({
      sourceId: propertySourceStore.id,
      endpoints: endpoints,
      username: formState.username,
      password: formState.password,
      name: formState.name
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/sa/propertySource/list/${formState.selectedEnv}`);
    });
  }
};

if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (propertySourceStore.id === 0) {
    router.push(`/sa/propertySource/list`);
  } else {
    formState.name = propertySourceStore.name;
    formState.selectedEnv = propertySourceStore.env;
    formState.username = propertySourceStore.username;
    formState.password = propertySourceStore.password;
    formState.endpoints = propertySourceStore.endpoints;
  }
}
</script>
<style scoped>
</style>