<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建Zallet代理节点</span>
        <span v-else-if="mode === 'update'">编辑Zallet代理节点</span>
      </div>
      <div class="section">
        <div class="section-title">名称</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div class="input-desc">标识代理节点</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">代理host</div>
        <div class="section-body">
          <a-input v-model:value="formState.agentHost" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">代理token</div>
        <div class="section-body">
          <a-input-password v-model:value="formState.agentToken" />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateZalletNode">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import {
  zalletNodeNameRegexp,
  zalletNodeAgentHostRegexp,
  zalletNodeAgentTokenRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import {
  createZalletNodeRequest,
  updateZalletNodeRequest
} from "@/api/zallet/zalletApi";
import { useRoute, useRouter } from "vue-router";
import { useZalletNodeStore } from "@/pinia/zalletNodeStore";
const zalletNodeStore = useZalletNodeStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const formState = reactive({
  name: "",
  agentHost: "",
  agentToken: ""
});
const saveOrUpdateZalletNode = () => {
  if (!zalletNodeNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!zalletNodeAgentHostRegexp.test(formState.agentHost)) {
    message.warn("host格式错误");
    return;
  }
  if (!zalletNodeAgentTokenRegexp.test(formState.agentToken)) {
    message.warn("token格式错误");
    return;
  }
  if (mode === "create") {
    createZalletNodeRequest({
      agentHost: formState.agentHost,
      agentToken: formState.agentToken,
      name: formState.name
    }).then(() => {
      message.success("创建成功");
      router.push(`/sa/zalletNode/list`);
    });
  } else if (mode === "update") {
    updateZalletNodeRequest({
      nodeId: zalletNodeStore.id,
      agentHost: formState.agentHost,
      agentToken: formState.agentToken,
      name: formState.name
    }).then(() => {
      message.success("保存成功");
      router.push(`/sa/zalletNode/list`);
    });
  }
};

if (mode === "update") {
  if (zalletNodeStore.id === 0) {
    router.push(`/sa/zalletNode/list`);
  } else {
    formState.name = zalletNodeStore.name;
    formState.agentHost = zalletNodeStore.agentHost;
    formState.agentToken = zalletNodeStore.agentToken;
  }
}
</script>
<style scoped>
</style>