<template>
  <div style="padding:10px">
    <div class="header">
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('gitWebhook.createWebhook')}}</a-button>
    </div>
    <ul class="webhook-list" v-if="webhooks.length > 0">
      <li v-for="item in webhooks" v-bind:key="item.id">
        <div class="webhook-pattern no-wrap">{{item.hookUrl}}</div>
        <ul class="op-btns">
          <li class="ping-btn" @click="pingWebhook(item)">{{t('gitWebhook.ping')}}</li>
          <li class="update-btn" @click="updateWebhook(item)">{{t('gitWebhook.update')}}</li>
          <li class="del-btn" @click="deleteWebhook(item)">{{t('gitWebhook.delete')}}</li>
        </ul>
      </li>
    </ul>
    <ZNoData v-else />
  </div>
</template>
<script setup>
import { ref, createVNode, h } from "vue";
import { useRouter, useRoute } from "vue-router";
import ZNoData from "@/components/common/ZNoData";
import {
  listWebhookRequest,
  pingWebhookRequest,
  deleteWebhookRequest
} from "@/api/git/webhookApi";
import { ExclamationCircleOutlined, PlusOutlined } from "@ant-design/icons-vue";
import { message, Modal } from "ant-design-vue";
import { useWebhookStore } from "@/pinia/webhookStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const webhooks = ref([]);
const webhookStore = useWebhookStore();
// 跳转新增webhook页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/webhook/create`
  );
};
// 删除webhook
const deleteWebhook = item => {
  Modal.confirm({
    title: `${t("gitWebhook.confirmDelete")} ${item.hookUrl}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteWebhookRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listWebhook();
      });
    },
    onCancel() {}
  });
};
// 获取webhook列表
const listWebhook = () => {
  listWebhookRequest(route.params.repoId).then(res => {
    webhooks.value = res.data;
  });
};
// 编辑webhook
const updateWebhook = item => {
  webhookStore.id = item.id;
  webhookStore.hookUrl = item.hookUrl;
  webhookStore.events = item.events;
  webhookStore.secret = item.secret;
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/webhook/${item.id}/update`
  );
};
// ping
const pingWebhook = item => {
  pingWebhookRequest(item.id).then(() => {
    message.success(t("operationSuccess"));
  });
};
listWebhook();
</script>
<style scoped>
.webhook-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.webhook-list > li {
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.webhook-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.webhook-pattern {
  font-size: 14px;
  line-height: 32px;
  width: 60%;
  padding-left: 10px;
}
.op-btns {
  display: flex;
  align-items: center;
}
.op-btns > li {
  line-height: 32px;
  font-size: 14px;
  padding: 0 10px;
  cursor: pointer;
}
.op-btns > li:first-child {
  border-top: 1px solid #d9d9d9;
  border-left: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
}
.op-btns > li:not(:first-child, :last-child) {
  border-top: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
}
.op-btns > li:last-child {
  border-top: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
  border-right: 1px solid #d9d9d9;
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
}
.op-btns > li + li {
  border-left: 1px solid #d9d9d9;
}
.header {
  margin-bottom: 10px;
}
.header > span {
  font-size: 18px;
  font-weight: bold;
  line-height: 32px;
  padding-left: 8px;
}
.del-btn {
  color: darkred;
}
.del-btn:hover {
  color: white;
  background-color: darkred;
}
.update-btn:hover,
.ping-btn:hover {
  background-color: #f0f0f0;
}
.no-data-text {
  font-size: 14px;
  line-height: 18px;
  padding: 10px;
  text-align: center;
  word-break: break-all;
}
</style>