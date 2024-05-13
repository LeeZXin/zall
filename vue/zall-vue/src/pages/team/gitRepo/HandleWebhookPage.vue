<template>
  <div style="padding:14px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">添加Webhook</span>
        <span v-else-if="mode === 'update'">更新Webhook</span>
      </div>
      <div class="section">
        <div class="section-title">
          <span>Hook url</span>
          <span style="color:darkred">*</span>
        </div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.hookUrl" :disabled="mode==='view'" />
          <div class="input-desc">必须以http开头的有效url</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>密钥</span>
        </div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.secret" />
          <div class="input-desc">我们会将通过这个密钥, 经过hmac后, 用于对请求的验证, 具体查看xxx</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>触发事件</span>
        </div>
        <div class="section-body">
          <ul class="event-list">
            <li>
              <a-checkbox v-model:checked="checkboxes.protectedBranch" style="flex-start">保护分支</a-checkbox>
              <div class="checkbox-desc">保护分支的添加、删除、修改事件</div>
            </li>
            <li>
              <a-checkbox v-model:checked="checkboxes.gitPush" style="flex-start">Git push</a-checkbox>
              <div class="checkbox-desc">通过ssh或http对git仓库进行push操作</div>
            </li>
            <li>
              <a-checkbox v-model:checked="checkboxes.pullRequest" style="flex-start">合并请求</a-checkbox>
              <div class="checkbox-desc">合并请求的新增、关闭、合并、评审</div>
            </li>
            <li>
              <a-checkbox v-model:checked="checkboxes.repo" style="flex-start">仓库</a-checkbox>
              <div class="checkbox-desc">仓库的删除、归档</div>
            </li>
          </ul>
        </div>
      </div>
      <div style="width:100%;border-top:1px solid #d9d9d9;margin: 10px 0"></div>
      <div style="margin-bottom:20px">
        <a-button type="primary" @click="createOrUpdateWebhook">立即保存</a-button>
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
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const checkboxes = reactive({
  protectedBranch: false,
  gitPush: false,
  pullRequest: false,
  repo: false
});
const eventMap = {
  protectedBranch: 1,
  gitPush: 2,
  pullRequest: 3,
  repo: 4
};
const webhookStore = useWebhookStore();
const router = useRouter();
const mode = getMode();
const formState = reactive({
  hookUrl: "",
  secret: ""
});
const createEvents = () => {
  let ret = [];
  if (checkboxes.protectedBranch) {
    ret.push(eventMap.protectedBranch);
  }
  if (checkboxes.gitPush) {
    ret.push(eventMap.gitPush);
  }
  if (checkboxes.pullRequest) {
    ret.push(eventMap.pullRequest);
  }
  if (checkboxes.repo) {
    ret.push(eventMap.repo);
  }
  return ret;
};
const createOrUpdateWebhook = () => {
  if (!webhookUrlRegexp.test(formState.hookUrl)) {
    message.warn("url格式错误");
    return;
  }
  if (!webhookSecretRegexp.test(formState.secret)) {
    message.warn("密钥格式错误");
    return;
  }
  const events = createEvents();
  if (events.length === 0) {
    message.warn("至少选择一个事件");
    return;
  }
  if (mode === "create") {
    createWebhookRequest({
      repoId: parseInt(route.params.repoId),
      hookUrl: formState.hookUrl,
      secret: formState.secret,
      events: events
    }).then(() => {
      message.success("添加成功");
      router.push(`/gitRepo/${route.params.repoId}/webhook/list`);
    });
  } else if (mode === "update") {
    updateWebhookRequest({
      webhookId: webhookStore.id,
      hookUrl: formState.hookUrl,
      secret: formState.secret,
      events: events
    }).then(() => {
      message.success("更新成功");
      router.push(`/gitRepo/${route.params.repoId}/webhook/list`);
    });
  }
};
if (mode !== "create") {
  if (
    webhookStore.id === 0 ||
    parseInt(route.params.webhookId) !== webhookStore.id
  ) {
    router.push(`/gitRepo/${route.params.repoId}/webhook/list`);
  } else {
    if (mode !== "create") {
      formState.hookUrl = webhookStore.hookUrl;
      formState.secret = webhookStore.secret;
      if (webhookStore.events && webhookStore.events.length > 0) {
        for (let index in webhookStore.events) {
          let item = webhookStore.events[index];
          switch (item) {
            case 1:
              checkboxes.protectedBranch = true;
              break;
            case 2:
              checkboxes.gitPush = true;
              break;
            case 3:
              checkboxes.pullRequest = true;
              break;
            case 4:
              checkboxes.repo = true;
              break;
          }
        }
      }
    }
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