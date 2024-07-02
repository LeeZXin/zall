<template>
  <div style="padding:14px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建服务来源</span>
        <span v-else-if="mode === 'update'">编辑服务来源</span>
      </div>
      <div
        style="margin-bottom:10px;font-size:14px;line-height: 20px;color:gray"
      >服务来源将作为服务状态数据源获取, 为了扩展性, 设计成以接口的形式获取服务状态</div>
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
          <a-input v-model:value="formState.name" />
          <div class="input-desc">标识服务来源</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">hosts</div>
        <div class="section-body">
          <a-input v-model:value="formState.hosts" />
          <div class="input-desc">ip:port格式, 多个按;隔开</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">Api Key</div>
        <div class="section-body">
          <a-input v-model:value="formState.apiKey" />
          <div class="input-desc">接口鉴权使用</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateServiceSource">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import {
  serviceSourceNameRegexp,
  serviceSourceHostRegexp,
  serviceSourceApiKeyRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  createServiceSourceRequest,
  updateServiceSourceRequest
} from "@/api/app/serviceSourceApi";
import { useRoute, useRouter } from "vue-router";
import { useServiceSourceStore } from "@/pinia/serviceSourceStore";
const serviceSourceStore = useServiceSourceStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const formState = reactive({
  name: "",
  hosts: "",
  apiKey: "",
  selectedEnv: ""
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
const saveOrUpdateServiceSource = () => {
  if (!serviceSourceNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!formState.hosts) {
    message.warn("hosts格式错误");
    return;
  }
  let hosts = formState.hosts.split(";").filter(item => item);
  for (let index in hosts) {
    if (!serviceSourceHostRegexp.test(hosts[index])) {
      message.warn("hosts格式错误");
      return;
    }
  }
  if (!serviceSourceApiKeyRegexp.test(formState.apiKey)) {
    message.warn("Api Key格式错误");
    return;
  }
  if (mode === "create") {
    createServiceSourceRequest(
      {
        env: formState.selectedEnv,
        appId: route.params.appId,
        apiKey: formState.apiKey,
        name: formState.name,
        hosts: hosts
      },
      formState.selectedEnv
    ).then(() => {
      message.success("创建成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/serviceSource/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "update") {
    updateServiceSourceRequest(
      {
        sourceId: serviceSourceStore.id,
        hosts: hosts,
        name: formState.name,
        apiKey: formState.apiKey
      },
      formState.selectedEnv
    ).then(() => {
      message.success("保存成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/serviceSource/list/${formState.selectedEnv}`
      );
    });
  }
};

if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (serviceSourceStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/serviceSource/list`
    );
  } else {
    formState.name = serviceSourceStore.name;
    formState.selectedEnv = serviceSourceStore.env;
    formState.hosts = serviceSourceStore.hosts;
    formState.apiKey = serviceSourceStore.apiKey;
  }
}
</script>
<style scoped>
</style>