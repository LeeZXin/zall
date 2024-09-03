<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建企微AccessToken任务</span>
        <span v-else-if="mode === 'update'">编辑企微AccessToken任务</span>
      </div>
      <div class="section">
        <div class="section-title">名称</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div class="input-desc">任务名称, 长度为1-32</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">corpId</div>
        <div class="section-body">
          <a-input v-model:value="formState.corpId" />
          <div class="input-desc">企业id</div>
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
  weworkAccessTokenNameRegexp,
  weworkAccessTokenCorpIdRegexp,
  weworkAccessTokenSecretRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import {
  createAccessTokenRequest,
  updateAccessTokenRequest
} from "@/api/team/weworkApi";
import { useRoute, useRouter } from "vue-router";
import { useWeworkAccessTokenStore } from "@/pinia/weworkAccessTokenStore";
const wwatStore = useWeworkAccessTokenStore();
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
  corpId: "",
  secret: ""
});
// 点击“立即保存”
const saveOrUpdateTask = () => {
  if (!weworkAccessTokenNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!weworkAccessTokenCorpIdRegexp.test(formState.corpId)) {
    message.warn("corpId格式错误");
    return;
  }
  if (!weworkAccessTokenSecretRegexp.test(formState.secret)) {
    message.warn("secret格式错误");
    return;
  }
  if (mode === "create") {
    // 创建
    createAccessTokenRequest({
      name: formState.name,
      teamId: parseInt(route.params.teamId),
      corpId: formState.corpId,
      secret: formState.secret
    }).then(() => {
      message.success("创建成功");
      router.push(`/team/${route.params.teamId}/weworkAccessToken/list`);
    });
  } else if (mode === "update") {
    // 编辑
    updateAccessTokenRequest({
      id: wwatStore.id,
      name: formState.name,
      corpId: formState.corpId,
      secret: formState.secret
    }).then(() => {
      message.success("保存成功");
      router.push(`/team/${route.params.teamId}/weworkAccessToken/list`);
    });
  }
};

if (mode === "update") {
  // store没有跳转list
  if (wwatStore.id === 0) {
    router.push(`/team/${route.params.teamId}/weworkAccessToken/list`);
  } else {
    formState.name = wwatStore.name;
    formState.corpId = wwatStore.corpId;
    formState.secret = wwatStore.secret;
  }
}
</script>
<style scoped>
</style>