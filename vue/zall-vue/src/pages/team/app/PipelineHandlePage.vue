<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('deployPipeline.createPipeline')}}</span>
        <span v-else-if="mode === 'update'">{{t('deployPipeline.updatePipeline')}}</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">{{t('deployPipeline.selectEnv')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.selectedEnv" :options="envList" />
        </div>
      </div>
      <div class="section" v-if="mode==='update'">
        <div class="section-title">{{t('deployPipeline.selectedEnv')}}</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('deployPipeline.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title flex-between">
          <span>{{t('deployPipeline.yaml')}}</span>
          <span @click="formatYaml" class="format-yaml-btn">{{t('deployPipeline.formatYaml')}}</span>
        </div>
        <Codemirror
          v-model="formState.content"
          style="height:380px;width:100%"
          :extensions="extensions"
        />
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdatePipeline">{{t('deployPipeline.save')}}</a-button>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const pipelineStore = usePipelineStore();
const route = useRoute();
const router = useRouter();
// 模式
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const extensions = ref([yaml(), oneDark]);
// 表单数据
const formState = reactive({
  name: "",
  content: "",
  selectedEnv: ""
});
// 环境列表
const envList = ref([]);
// 获取环境配置
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
// 格式化yaml
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
// 新增或编辑流水线
const saveOrUpdatePipeline = () => {
  if (!pipelineNameRegexp.test(formState.name)) {
    message.warn(t("deployPipeline.nameFormatErr"));
    return;
  }
  if (mode === "create") {
    createPipelineRequest({
      env: formState.selectedEnv,
      appId: route.params.appId,
      config: formState.content,
      name: formState.name
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "update") {
    updatePipelineRequest({
      pipelineId: pipelineStore.id,
      config: formState.content,
      name: formState.name
    }).then(() => {
      message.success(t("operationSuccess"));
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