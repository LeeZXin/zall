<template>
  <div style="padding:14px">
    <div class="container">
      <div class="title">创建应用服务</div>
      <div class="section">
        <div class="section-title">AppId</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.appId" />
          <div class="input-desc">应用服务的唯一标识, 不能包含空格等特殊字符</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">应用名称</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.name" />
          <div class="input-desc">简单的话来描述应用服务</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createApp">立即创建</a-button>
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
const router = useRouter();
const route = useRoute();
const formState = reactive({
  appId: "",
  name: ""
});
const createApp = () => {
  if (!appIdRegexp.test(formState.appId)) {
    message.warn("appId格式不正确");
    return;
  }
  if (!appNameRegexp.test(formState.name)) {
    message.warn("名称格式不正确");
    return;
  }
  createAppRequest({
    teamId: parseInt(route.params.teamId),
    appId: formState.appId,
    name: formState.name
  }).then(() => {
    message.success("创建成功");
    router.push(`/team/${route.params.teamId}/app/list`);
  });
};
</script>
<style scoped>
</style>