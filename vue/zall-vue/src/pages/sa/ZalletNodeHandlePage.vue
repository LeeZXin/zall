<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('zallet.createNode')}}</span>
        <span v-else-if="mode === 'update'">{{t('zallet.updateNode')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('zallet.nodeId')}}</div>
        <div class="section-body">
          <template v-if="mode==='create'">
            <a-input v-model:value="formState.nodeId" />
          </template>
          <template v-if="mode==='update'">
            <span>{{formState.nodeId}}</span>
          </template>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('zallet.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('zallet.agentHost')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.agentHost" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('zallet.agentToken')}}</div>
        <div class="section-body">
          <a-input-password v-model:value="formState.agentToken" />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateZalletNode">{{t('zallet.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import {
  zalletNodeIdRegexp,
  zalletNameRegexp,
  zalletAgentHostRegexp,
  zalletAgentTokenRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import {
  createZalletNodeRequest,
  updateZalletNodeRequest
} from "@/api/zallet/zalletApi";
import { useRoute, useRouter } from "vue-router";
import { useZalletNodeStore } from "@/pinia/zalletNodeStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const zalletNodeStore = useZalletNodeStore();
const route = useRoute();
const router = useRouter();
// 模式
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
// 表单数据
const formState = reactive({
  name: "",
  nodeId: "",
  agentHost: "",
  agentToken: ""
});
// 新增或保存zallet
const saveOrUpdateZalletNode = () => {
  if (!zalletNameRegexp.test(formState.name)) {
    message.warn(t("zallet.nameFormatErr"));
    return;
  }
  if (!zalletAgentHostRegexp.test(formState.agentHost)) {
    message.warn(t("zallet.agentHostFormatErr"));
    return;
  }
  if (!zalletAgentTokenRegexp.test(formState.agentToken)) {
    message.warn(t("zallet.agentTokenFormatErr"));
    return;
  }
  if (mode === "create") {
    if (!zalletNodeIdRegexp.test(formState.nodeId)) {
      message.warn(t("zallet.nodeIdFormatErr"));
      return;
    }
    createZalletNodeRequest({
      nodeId: formState.nodeId,
      agentHost: formState.agentHost,
      agentToken: formState.agentToken,
      name: formState.name
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/sa/zalletNode/list`);
    });
  } else if (mode === "update") {
    updateZalletNodeRequest({
      id: zalletNodeStore.id,
      agentHost: formState.agentHost,
      agentToken: formState.agentToken,
      name: formState.name
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/sa/zalletNode/list`);
    });
  }
};

if (mode === "update") {
  if (zalletNodeStore.id === 0) {
    router.push(`/sa/zalletNode/list`);
  } else {
    formState.nodeId = zalletNodeStore.nodeId;
    formState.name = zalletNodeStore.name;
    formState.agentHost = zalletNodeStore.agentHost;
    formState.agentToken = zalletNodeStore.agentToken;
  }
}
</script>
<style scoped>
</style>