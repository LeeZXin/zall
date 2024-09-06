<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">{{t('deployPlan.createPlan')}}</div>
      <div class="section">
        <div class="section-title">{{t('deployPlan.selectEnv')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.selectedEnv" :options="envList" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('deployPlan.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('deployPlan.selectPipeline')}}</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.pipelineId"
            :options="pipelineList"
          />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('deployPlan.selectArtifact')}}</div>
        <div class="section-body">
          <a-select
            v-model:value="formState.artifactVersion"
            :options="artifactList"
            style="width:100%"
          />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createDeployPlan">{{t('deployPlan.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref, watch } from "vue";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  listPipelineWhenCreateDeployPlanRequest,
  createDeployPlanRequest
} from "@/api/app/deployPlanApi";
import { listLatestArtifactRequest } from "@/api/app/artifactApi";
import { useRoute, useRouter } from "vue-router";
import { deployPlanNameRegexp } from "@/utils/regexp";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const route = useRoute();
const router = useRouter();
// 表单数据
const formState = reactive({
  name: "",
  selectedEnv: undefined,
  pipelineId: undefined,
  artifactVersion: ""
});
// 流水线列表
const pipelineList = ref([]);
// 环境列表
const envList = ref([]);
// 制品列表
const artifactList = ref([]);
// 获取环境
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
// 创建发布计划
const createDeployPlan = () => {
  if (!deployPlanNameRegexp.test(formState.name)) {
    message.warn(t("deployPlan.nameFormatErr"));
    return;
  }
  if (!formState.pipelineId) {
    message.warn(t("deployPlan.pleaseSelectPipeline"));
    return;
  }
  if (!formState.artifactVersion) {
    message.warn(t("deployPlan.pleaseSelectArtifact"));
    return;
  }
  createDeployPlanRequest({
    pipelineId: formState.pipelineId,
    artifactVersion: formState.artifactVersion,
    name: formState.name
  }).then(() => {
    message.success(t("operationSuccess"));
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/deployPlan/list/${formState.selectedEnv}`
    );
  });
};
// 获取流水线列表
const listPipeline = () => {
  listPipelineWhenCreateDeployPlanRequest({
    appId: route.params.appId,
    env: formState.selectedEnv
  }).then(res => {
    pipelineList.value = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
    if (res.data?.length > 0) {
      formState.pipelineId = res.data[0].id;
    } else {
      formState.pipelineId = undefined;
    }
  });
};
// 获取制品
const listArtifact = () => {
  listLatestArtifactRequest({
    appId: route.params.appId,
    env: formState.selectedEnv
  }).then(res => {
    artifactList.value = res.data.map(item => {
      return {
        value: item.name,
        label: `${item.name}(${item.created})`
      };
    });
    if (artifactList.value.length > 0) {
      formState.artifactVersion = artifactList.value[0].value;
    } else {
      formState.artifactVersion = undefined;
    }
  });
};

watch(
  () => formState.selectedEnv,
  () => {
    listPipeline();
    listArtifact();
  }
);

getEnvCfg();
</script>
<style scoped>
</style>