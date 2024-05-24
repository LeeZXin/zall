<template>
  <div style="padding:14px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">添加密钥</span>
        <span v-else-if="mode === 'update'">更新密钥</span>
      </div>
      <div class="section">
        <div class="section-title">
          <span>Key</span>
        </div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.name" :disabled="mode==='update'" />
          <div class="input-desc">key用来唯一标识密钥, 不为空, 长度不得超过32</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>内容</span>
        </div>
        <div class="section-body">
          <a-textarea
            style="width:100%"
            v-model:value="formState.content"
            :auto-size="{ minRows: 5, maxRows: 10 }"
          />
          <div class="input-desc">密钥的具体内容</div>
        </div>
      </div>
      <div style="width:100%;border-top:1px solid #d9d9d9;margin: 10px 0"></div>
      <div style="margin-bottom:20px">
        <a-button type="primary" @click="createOrUpdateSecret">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import {
  createSecretRequest,
  updateSecretRequest,
  getSecretContentRequest
} from "@/api/git/workflowApi";
import { useRoute, useRouter } from "vue-router";
import { workflowSecretNameRegexp, workflowSecretContentRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const router = useRouter();
const mode = getMode();
const formState = reactive({
  name: "",
  content: ""
});
const createOrUpdateSecret = () => {
  if (!workflowSecretNameRegexp.test(formState.name)) {
    message.warn("key格式错误");
    return;
  }
  if (!workflowSecretContentRegexp.test(formState.content)) {
    message.warn("内容格式错误");
    return;
  }
  if (mode === "create") {
    createSecretRequest({
      repoId: parseInt(route.params.repoId),
      name: formState.name,
      content: formState.content
    }).then(() => {
      message.success("添加成功");
      router.push(`/gitRepo/${route.params.repoId}/workflow/secrets`);
    });
  } else if (mode === "update") {
    updateSecretRequest({
      secretId: parseInt(route.params.secretId),
      content: formState.content
    }).then(() => {
      message.success("更新成功");
      router.push(`/gitRepo/${route.params.repoId}/workflow/secrets`);
    });
  }
};
if (mode === "update") {
  getSecretContentRequest(route.params.secretId).then(res => {
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