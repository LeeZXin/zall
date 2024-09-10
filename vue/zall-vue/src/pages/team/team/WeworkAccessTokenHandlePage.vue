<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('weworkAccessToken.createTask')}}</span>
        <span v-else-if="mode === 'update'">{{t('weworkAccessToken.updateTask')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('weworkAccessToken.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('weworkAccessToken.corpId')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.corpId" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('weworkAccessToken.secret')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.secret" />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateTask">{{t('weworkAccessToken.save')}}</a-button>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
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
    message.warn(t("weworkAccessToken.nameFormatErr"));
    return;
  }
  if (!weworkAccessTokenCorpIdRegexp.test(formState.corpId)) {
    message.warn(t("weworkAccessToken.corpIdFormatErr"));
    return;
  }
  if (!weworkAccessTokenSecretRegexp.test(formState.secret)) {
    message.warn(t("weworkAccessToken.secretFormatErr"));
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
      message.success(t("operationSuccess"));
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
      message.success(t("operationSuccess"));
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