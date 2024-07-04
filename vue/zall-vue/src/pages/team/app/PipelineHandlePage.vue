<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建流水线</span>
        <span v-else-if="mode === 'update'">编辑流水线</span>
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
          <div class="input-desc">标识配置作用</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title flex-between">
          <span>配置内容</span>
          <span @click="formatYaml" class="format-yaml-btn">格式化yaml</span>
        </div>
        <Codemirror
          v-model="formState.content"
          style="height:380px;width:100%"
          :extensions="extensions"
        />
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdatePipeline">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { Codemirror } from "vue-codemirror";
import { yaml } from "@codemirror/lang-yaml";
import { oneDark } from "@codemirror/theme-one-dark";
import { pipelineNameRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  createPipelineRequest,
  updatePipelineRequest
} from "@/api/app/pipelineApi";
import { useRoute, useRouter } from "vue-router";
import { usePipelineStore } from "@/pinia/pipelineStore";
import jsyaml from "js-yaml";
const pipelineStore = usePipelineStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const extensions = ref([yaml(), oneDark]);
const formState = reactive({
  name: "",
  content: "",
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
const formatYaml = () => {
  if (formState.content) {
    try {
      const parsedYaml = jsyaml.load(formState.content);
      formState.content = jsyaml.dump(parsedYaml);
    } catch (e) {
      //
    }
  }
};
const saveOrUpdatePipeline = () => {
  if (!pipelineNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (mode === "create") {
    createPipelineRequest(
      {
        env: formState.selectedEnv,
        appId: route.params.appId,
        config: formState.content,
        name: formState.name
      },
      formState.selectedEnv
    ).then(() => {
      message.success("创建成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "update") {
    updatePipelineRequest(
      {
        pipelineId: pipelineStore.id,
        config: formState.content,
        name: formState.name
      },
      formState.selectedEnv
    ).then(() => {
      message.success("保存成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/list/${formState.selectedEnv}`
      );
    });
  }
};

if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (pipelineStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/list`
    );
  } else {
    formState.name = pipelineStore.name;
    formState.selectedEnv = pipelineStore.env;
    formState.content = pipelineStore.config;
  }
}
</script>
<style scoped>
.format-yaml-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>