<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('deployPipeline.createVars')}}</span>
        <span v-else-if="mode === 'update'">{{t('deployPipeline.updateVars')}}</span>
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
        <div class="section-title">{{t('deployPipeline.varsKey')}}</div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.name" :disabled="mode==='update'" />
          <div class="input-desc">{{t('deployPipeline.varsKeyDesc')}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('deployPipeline.varsContent')}}</div>
        <div class="section-body">
          <a-textarea
            style="width:100%"
            v-model:value="formState.content"
            :auto-size="{ minRows: 8, maxRows: 15 }"
            @keydown.tab="handleTab"
          />
        </div>
      </div>
      <div style="width:100%;border-top:1px solid #d9d9d9;margin: 10px 0"></div>
      <div style="margin-bottom:20px">
        <a-button type="primary" @click="saveOrUpdateVars">{{t('deployPipeline.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import {
  createPipelineVarsRequest,
  updatePipelineVarsRequest,
  getPipelineVarsRequest
} from "@/api/app/pipelineApi";
import { useRoute, useRouter } from "vue-router";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  pipelineVarsNameRegexp,
  pipelineVarsContentRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import { usePipelineVarsStore } from "@/pinia/pipelineVarsStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const varsStore = usePipelineVarsStore();
const route = useRoute();
// 模式
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const router = useRouter();
const mode = getMode();
// 表单数据
const formState = reactive({
  name: "",
  content: "",
  selectedEnv: ""
});
// 环境列表
const envList = ref([]);
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
// tab键处理
const handleTab = event => {
  event.preventDefault();
  let inputElement = event.target;
  let value = inputElement.value;
  let selectionStart = inputElement.selectionStart;
  let leftValue = value.substring(0, selectionStart);
  let rightValue = value.substring(selectionStart);
  inputElement.value = leftValue + "    " + rightValue;
  inputElement.selectionStart = selectionStart + 4;
  inputElement.selectionEnd = inputElement.selectionStart;
};
// 新增或编辑变量
const saveOrUpdateVars = () => {
  if (!pipelineVarsNameRegexp.test(formState.name)) {
    message.warn(t("deployPipeline.varsKeyFormatErr"));
    return;
  }
  if (!pipelineVarsContentRegexp.test(formState.content)) {
    message.warn(t("deployPipeline.varsContentFormatErr"));
    return;
  }
  if (mode === "create") {
    createPipelineVarsRequest({
      appId: route.params.appId,
      name: formState.name,
      content: formState.content,
      env: formState.selectedEnv
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/vars/${formState.selectedEnv}`
      );
    });
  } else if (mode === "update") {
    updatePipelineVarsRequest({
      id: varsStore.id,
      content: formState.content
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/vars/${formState.selectedEnv}`
      );
    });
  }
};
if (mode === "update") {
  if (varsStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/vars`
    );
  } else {
    formState.selectedEnv = varsStore.env;
    getPipelineVarsRequest(varsStore.id).then(res => {
      formState.name = res.data.name;
      formState.content = res.data.content;
    });
  }
} else if (mode === "create") {
  getEnvCfg();
}
</script>
<style scoped>
.header {
  font-size: 18px;
  margin-bottom: 10px;
  font-weight: bold;
}
.event-list {
  font-size: 14px;
  display: flex;
  flex-wrap: wrap;
}
.event-list > li {
  padding-right: 10px;
  width: 50%;
  margin-bottom: 16px;
}
</style>