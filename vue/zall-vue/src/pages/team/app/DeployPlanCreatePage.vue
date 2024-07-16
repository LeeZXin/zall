<template>
  <div style="padding:14px">
    <div class="container">
      <div class="title">创建发布计划</div>
      <div class="section">
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
      <div class="section">
        <div class="section-title">计划名称</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div class="input-desc">标识发布计划作用</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">选择流水线</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.pipelineId"
            :options="pipelineList"
          />
          <div class="input-desc">选择服务, 使用其发布流水线发布服务</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">制品号</div>
        <div class="section-body">
          <a-input v-model:value="formState.productVersion" />
          <div class="input-desc">选择或填写发布的制品</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createDeployPlan">立即创建</a-button>
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
import { useRoute, useRouter } from "vue-router";
import {
  deployProductVersionRegexp,
  deployPlanNameRegexp
} from "@/utils/regexp";
const route = useRoute();
const router = useRouter();
const formState = reactive({
  name: "",
  selectedEnv: undefined,
  pipelineId: undefined,
  productVersion: ""
});
const pipelineList = ref([]);
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

const createDeployPlan = () => {
  if (!deployPlanNameRegexp.test(formState.name)) {
    message.warn("名称格式不正确");
    return;
  }
  if (!deployProductVersionRegexp.test(formState.productVersion)) {
    message.warn("制品号格式不正确");
    return;
  }
  if (!formState.pipelineId) {
    message.warn("请选择流水线");
    return;
  }
  createDeployPlanRequest({
    pipelineId: formState.pipelineId,
    productVersion: formState.productVersion,
    name: formState.name
  }).then(() => {
    message.success("创建成功");
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/deployPlan/list/${formState.selectedEnv}`
    );
  });
};

const listService = () => {
  listPipelineWhenCreateDeployPlanRequest({
    appId: route.params.appId,
    env: formState.selectedEnv
  }).then(res => {
    if (res.data?.length > 0) {
      pipelineList.value = res.data.map(item => {
        return {
          value: item.id,
          label: item.name
        };
      });
      formState.pipelineId = res.data[0].id;
    } else {
      pipelineList.value = [];
      formState.pipelineId = undefined;
    }
  });
};

watch(
  () => formState.selectedEnv,
  () => {
    listService();
  }
);

getEnvCfg();
</script>
<style scoped>
.diff-btn:hover {
  cursor: pointer;
  color: #1677ff;
}
</style>