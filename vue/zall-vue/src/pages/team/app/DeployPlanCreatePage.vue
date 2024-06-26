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
        <div class="section-title">选择服务</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            placeholder="选择环境"
            v-model:value="formState.serviceId"
            :options="serviceList"
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
  listServiceWhenCreateDeployPlanRequest,
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
  serviceId: undefined,
  productVersion: ""
});
const serviceList = ref([]);
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
  if (!formState.serviceId) {
    message.warn("请选择服务");
    return;
  }
  createDeployPlanRequest(
    {
      serviceId: formState.serviceId,
      productVersion: formState.productVersion,
      name: formState.name
    },
    formState.selectedEnv
  ).then(() => {
    message.success("创建成功");
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/deployPlan/list/${formState.selectedEnv}`
    );
  });
};

const listService = () => {
  listServiceWhenCreateDeployPlanRequest(
    {
      appId: route.params.appId,
      env: formState.selectedEnv
    },
    formState.selectedEnv
  ).then(res => {
    if (res.data?.length > 0) {
      serviceList.value = res.data.map(item => {
        return {
          value: item.id,
          label: item.name
        };
      });
      formState.serviceId = res.data[0].id;
    } else {
      serviceList.value = [];
      formState.serviceId = undefined;
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