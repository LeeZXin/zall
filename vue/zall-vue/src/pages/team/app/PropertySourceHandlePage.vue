<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建配置来源</span>
        <span v-else-if="mode === 'update'">编辑配置来源</span>
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
          <a-input v-model:value="formState.name" />
          <div class="input-desc">标识配置来源</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">endpoints</div>
        <div class="section-body">
          <a-input v-model:value="formState.endpoints" />
          <div class="input-desc">etcd endpoints ip:port格式, 以;隔开</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">账号</div>
        <div class="section-body">
          <a-input v-model:value="formState.username" />
          <div class="input-desc">etcd账号</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">密码</div>
        <div class="section-body">
          <a-input-password v-model:value="formState.password" />
          <div class="input-desc">etcd密码</div>
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
import { propertySourceNameRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  createPropertySourceRequest,
  updatePropertySourceRequest
} from "@/api/app/propertyApi";
import { useRoute, useRouter } from "vue-router";
import { usePropertySourceStore } from "@/pinia/propertySourceStore";
const propertySourceStore = usePropertySourceStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const formState = reactive({
  name: "",
  endpoints: "",
  username: "",
  password: "",
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
  if (!propertySourceNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!formState.endpoints) {
    message.warn("endpoints格式错误");
    return;
  }
  let endpoints = formState.endpoints.split(";").filter(item => item);
  if (mode === "create") {
    createPropertySourceRequest({
      env: formState.selectedEnv,
      appId: route.params.appId,
      endpoints: endpoints,
      username: formState.username,
      password: formState.password,
      name: formState.name
    }).then(() => {
      message.success("创建成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/propertySource/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "update") {
    updatePropertySourceRequest({
      sourceId: propertySourceStore.id,
      endpoints: endpoints,
      username: formState.username,
      password: formState.password,
      name: formState.name
    }).then(() => {
      message.success("保存成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/propertySource/list/${formState.selectedEnv}`
      );
    });
  }
};

if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (propertySourceStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/propertySource/list`
    );
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