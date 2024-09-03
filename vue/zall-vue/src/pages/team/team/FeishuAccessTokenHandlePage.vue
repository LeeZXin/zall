<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建飞书AccessToken任务</span>
        <span v-else-if="mode === 'update'">编辑飞书AccessToken任务</span>
      </div>
      <div class="section">
        <div class="section-title">名称</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div class="input-desc">任务名称, 长度为1-32</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">appId</div>
        <div class="section-body">
          <a-input v-model:value="formState.appId" />
          <div class="input-desc">应用id</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">secret</div>
        <div class="section-body">
          <a-input v-model:value="formState.secret" />
          <div class="input-desc">应用密钥</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateTask">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import {
  feishuAccessTokenNameRegexp,
  feishuAccessTokenAppIdRegexp,
  feishuAccessTokenSecretRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import {
  createAccessTokenRequest,
  updateAccessTokenRequest
} from "@/api/team/feishuApi";
import { useRoute, useRouter } from "vue-router";
import { useFeishuAccessTokenStore } from "@/pinia/feishuAccessTokenStore";
const fsatStore = useFeishuAccessTokenStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
// 表单数据
const formState = reactive({
  name: "",
  appId: "",
  secret: ""
});
// 点击“立即保存”
const saveOrUpdateTask = () => {
  if (!feishuAccessTokenNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!feishuAccessTokenAppIdRegexp.test(formState.appId)) {
    message.warn("appId格式错误");
    return;
  }
  if (!feishuAccessTokenSecretRegexp.test(formState.secret)) {
    message.warn("secret格式错误");
    return;
  }
  if (mode === "create") {
    // 创建
    createAccessTokenRequest({
      name: formState.name,
      teamId: parseInt(route.params.teamId),
      appId: formState.appId,
      secret: formState.secret
    }).then(() => {
      message.success("创建成功");
      router.push(`/team/${route.params.teamId}/feishuAccessToken/list`);
    });
  } else if (mode === "update") {
    // 编辑
    updateAccessTokenRequest({
      id: fsatStore.id,
      name: formState.name,
      appId: formState.appId,
      secret: formState.secret
    }).then(() => {
      message.success("保存成功");
      router.push(`/team/${route.params.teamId}/feishuAccessToken/list`);
    });
  }
};

if (mode === "update") {
  // store没有跳转list
  if (fsatStore.id === 0) {
    router.push(`/team/${route.params.teamId}/feishuAccessToken/list`);
  } else {
    formState.name = fsatStore.name;
    formState.appId = fsatStore.appId;
    formState.secret = fsatStore.secret;
  }
}
</script>
<style scoped>
</style>