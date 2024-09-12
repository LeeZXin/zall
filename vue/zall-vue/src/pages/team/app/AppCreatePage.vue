<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">{{t('appService.createApp')}}</div>
      <div class="section">
        <div class="section-title">{{t('appService.appId')}}</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.appId" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('appService.name')}}</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.name" />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createApp">{{t('appService.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import { createAppRequest } from "@/api/app/appApi";
import { appNameRegexp, appIdRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const route = useRoute();
// 表单数据
const formState = reactive({
  appId: "",
  name: ""
});
// 创建app
const createApp = () => {
  if (!appIdRegexp.test(formState.appId)) {
    message.warn(t("appService.appIdFormatErr"));
    return;
  }
  if (!appNameRegexp.test(formState.name)) {
    message.warn(t("appService.nameFormatErr"));
    return;
  }
  createAppRequest({
    teamId: parseInt(route.params.teamId),
    appId: formState.appId,
    name: formState.name
  }).then(() => {
    message.success(t("operationSuccess"));
    router.push(`/team/${route.params.teamId}/app/list`);
  });
};
</script>
<style scoped>
</style>