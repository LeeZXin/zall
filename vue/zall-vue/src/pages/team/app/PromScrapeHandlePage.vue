<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建抓取任务</span>
        <span v-else-if="mode === 'update'">编辑抓取任务</span>
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
        <div class="section-title">endpoint</div>
        <div class="section-body">
          <a-input v-model:value="formState.endpoint" />
          <div class="input-desc">prometheus标识</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">目标类型</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.targetType"
            :options="targetTypeList"
          />
          <div class="input-desc">选择服务发现类型, 则利用注册中心发现服务后抓取, 适合配置了注册中心的服务. 直连类型适合ip不经常变动的服务.</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">目标</div>
        <div class="section-body">
          <a-input v-model:value="formState.target" />
          <div class="input-desc">抓取目标, 若选择服务发现类型, 则填写服务发现的key, 例如xxx-http. 若选择直连类型, 多个ip:port用;隔开</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateScrape">立即保存</a-button>
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
const promScrapeStore = usePromScrapeStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const formState = reactive({
  endpoint: "",
  target: "",
  targetType: 2
});
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
const targetTypeList = [
  {
    value: 1,
    label: "服务发现类型"
  },
  {
    value: 2,
    label: "直连类型"
  }
];
const saveOrUpdateScrape = () => {
  if (!promScrapeEndpointRegexp.test(formState.endpoint)) {
    message.warn("endpoint格式错误");
    return;
  }
  if (!promScrapeTargetRegexp.test(formState.target)) {
    message.warn("目标格式错误");
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
      message.success("创建成功");
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
      message.success("保存成功");
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