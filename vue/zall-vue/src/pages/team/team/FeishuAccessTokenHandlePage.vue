<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('feishuAccessToken.createTask')}}</span>
        <span v-else-if="mode === 'update'">{{t('feishuAccessToken.updateTask')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('feishuAccessToken.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('feishuAccessToken.appId')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.appId" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('feishuAccessToken.secret')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.secret" />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateTask">{{t('feishuAccessToken.save')}}</a-button>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
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
    message.warn(t("feishuAccessToken.nameFormatErr"));
    return;
  }
  if (!feishuAccessTokenAppIdRegexp.test(formState.appId)) {
    message.warn(t("feishuAccessToken.appIdFormatErr"));
    return;
  }
  if (!feishuAccessTokenSecretRegexp.test(formState.secret)) {
    message.warn(t("feishuAccessToken.secretFormatErr"));
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
      message.success(t("operationSuccess"));
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
      message.success(t("operationSuccess"));
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