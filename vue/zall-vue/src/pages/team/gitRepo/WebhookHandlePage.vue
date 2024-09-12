<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('gitWebhook.createWebhook')}}</span>
        <span v-else-if="mode === 'update'">{{t('gitWebhook.updateWebhook')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('gitWebhook.webhookUrl')}}</div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.hookUrl" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('gitWebhook.webhookSecret')}}</div>
        <div class="section-body">
          <a-input-password style="width:100%" v-model:value="formState.secret" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>{{t('gitWebhook.webhookEvents')}}</span>
        </div>
        <div class="section-body">
          <ul class="event-list">
            <li>
              <a-checkbox
                v-model:checked="checkboxes.protectedBranch"
              >{{t('gitWebhook.protectedBranch')}}</a-checkbox>
              <div class="checkbox-desc">{{t('gitWebhook.protectedBranchDesc')}}</div>
            </li>
            <li>
              <a-checkbox v-model:checked="checkboxes.gitPush">{{t('gitWebhook.gitPush')}}</a-checkbox>
              <div class="checkbox-desc">{{t('gitWebhook.gitPushDesc')}}</div>
            </li>
            <li>
              <a-checkbox v-model:checked="checkboxes.pullRequest">{{t('gitWebhook.pullRequest')}}</a-checkbox>
              <div class="checkbox-desc">{{t('gitWebhook.pullRequestDesc')}}</div>
            </li>
            <li>
              <a-checkbox v-model:checked="checkboxes.gitRepo">{{t('gitWebhook.gitRepo')}}</a-checkbox>
              <div class="checkbox-desc">{{t('gitWebhook.gitRepoDesc')}}</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createOrUpdateWebhook">{{t('gitWebhook.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import {
  createWebhookRequest,
  updateWebhookRequest
} from "@/api/git/webhookApi";
import { useRoute, useRouter } from "vue-router";
import { webhookUrlRegexp, webhookSecretRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useWebhookStore } from "@/pinia/webhookStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const route = useRoute();
// 模式
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
// 四个checkbox数据
const checkboxes = reactive({
  protectedBranch: false,
  gitPush: false,
  pullRequest: false,
  gitRepo: false
});
const webhookStore = useWebhookStore();
const router = useRouter();
const mode = getMode();
// 表单数据
const formState = reactive({
  hookUrl: "",
  secret: ""
});
// 新增或编辑webhook
const createOrUpdateWebhook = () => {
  if (!webhookUrlRegexp.test(formState.hookUrl)) {
    message.warn(t("gitWebhook.webhookUrlFormatErr"));
    return;
  }
  if (!webhookSecretRegexp.test(formState.secret)) {
    message.warn(t("gitWebhook.webhookSecretFormatErr"));
    return;
  }
  if (
    !checkboxes.pullRequest &&
    !checkboxes.gitRepo &&
    !checkboxes.gitPush &&
    !checkboxes.protectedBranch
  ) {
    message.warn(t("gitWebhook.atLeastChooseOneEvent"));
    return;
  }
  if (mode === "create") {
    createWebhookRequest({
      repoId: parseInt(route.params.repoId),
      hookUrl: formState.hookUrl,
      secret: formState.secret,
      events: checkboxes
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/webhook/list`
      );
    });
  } else if (mode === "update") {
    updateWebhookRequest({
      webhookId: webhookStore.id,
      hookUrl: formState.hookUrl,
      secret: formState.secret,
      events: checkboxes
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/webhook/list`
      );
    });
  }
};
if (mode === "update") {
  if (webhookStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/webhook/list`
    );
  } else {
    formState.hookUrl = webhookStore.hookUrl;
    formState.secret = webhookStore.secret;
    checkboxes.protectedBranch = webhookStore.events?.protectedBranch;
    checkboxes.gitPush = webhookStore.events?.gitPush;
    checkboxes.pullRequest = webhookStore.events?.pullRequest;
    checkboxes.gitRepo = webhookStore.events?.gitRepo;
  }
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