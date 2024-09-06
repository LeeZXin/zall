<template>
  <div style="padding:10px">
    <div class="header flex-between">
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('deployPipeline.createVars')}}</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ul class="vars-list" v-if="varsList.length > 0">
      <li v-for="item in varsList" v-bind:key="item.id">
        <div class="vars-pattern no-wrap">{{item.name}}</div>
        <ul class="op-btns">
          <li class="update-btn" @click="handleVars(item)">{{t('deployPipeline.updateVars')}}</li>
          <li class="del-btn" @click="deleteVars(item)">{{t('deployPipeline.deleteVars')}}</li>
        </ul>
      </li>
    </ul>
    <ZNoData v-else />
  </div>
</template>
<script setup>
import ZNoData from "@/components/common/ZNoData";
import EnvSelector from "@/components/app/EnvSelector";
import { ref, createVNode, h } from "vue";
import { useRouter, useRoute } from "vue-router";
import { ExclamationCircleOutlined, PlusOutlined } from "@ant-design/icons-vue";
import { message, Modal } from "ant-design-vue";
import {
  listPipelineVarsRequest,
  deletePipelineVarsRequest
} from "@/api/app/pipelineApi";
import { usePipelineVarsStore } from "@/pinia/pipelineVarsStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const varsStore = usePipelineVarsStore();
const router = useRouter();
const route = useRoute();
// 数据
const varsList = ref([]);
// 选择的环境
const selectedEnv = ref("");
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/vars/create?env=${selectedEnv.value}`
  );
};
// 删除变量
const deleteVars = item => {
  Modal.confirm({
    title: `${t('deployPipeline.confirmDelete')} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deletePipelineVarsRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listVars();
      });
    },
    onCancel() {}
  });
};
// 变量列表
const listVars = () => {
  listPipelineVarsRequest({
    appId: route.params.appId,
    env: selectedEnv.value
  }).then(res => {
    varsList.value = res.data;
  });
};
// 编辑变量
const handleVars = item => {
  varsStore.id = item.id;
  varsStore.name = item.name;
  varsStore.env = item.env;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/vars/${item.id}/update`
  );
};
// 环境变化
const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/vars/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  listVars();
};
</script>
<style scoped>
.vars-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.vars-list > li {
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.vars-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.vars-pattern {
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
.update-btn:hover {
  background-color: #f0f0f0;
}
</style>