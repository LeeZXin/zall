<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('gitWorkflow.createVars')}}</span>
        <span v-else-if="mode === 'update'">{{t('gitWorkflow.updateVars')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('gitWorkflow.varsKey')}}</div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.name" :disabled="mode==='update'" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('gitWorkflow.varsContent')}}</div>
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
        <a-button type="primary" @click="createOrUpdateVars">{{t('gitWorkflow.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import {
  createVarsRequest,
  updateVarsRequest,
  getVarsContentRequest
} from "@/api/git/workflowApi";
import { useRoute, useRouter } from "vue-router";
import {
  workflowVarsNameRegexp,
  workflowVarsContentRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const router = useRouter();
// 模式
const mode = getMode();
// 表单数据
const formState = reactive({
  name: "",
  content: ""
});
// tab默认行为
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
const createOrUpdateVars = () => {
  if (!workflowVarsNameRegexp.test(formState.name)) {
    message.warn(t("gitWorkflow.varsKeyFormatErr"));
    return;
  }
  if (!workflowVarsContentRegexp.test(formState.content)) {
    message.warn(t("gitWorkflow.varsContentFormatErr"));
    return;
  }
  if (mode === "create") {
    createVarsRequest({
      repoId: parseInt(route.params.repoId),
      name: formState.name,
      content: formState.content
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/vars`
      );
    });
  } else if (mode === "update") {
    updateVarsRequest({
      varsId: parseInt(route.params.varsId),
      content: formState.content
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/vars`
      );
    });
  }
};
if (mode === "update") {
  getVarsContentRequest(route.params.varsId).then(res => {
    formState.name = res.data.name;
    formState.content = res.data.content;
  });
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