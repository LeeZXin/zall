<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <div>
        <a-button
          type="primary"
          :icon="h(PlusOutlined)"
          @click="gotoCreatePage"
        >{{t('deployPipeline.createPipeline')}}</a-button>
        <a-button
          type="primary"
          @click="gotoVarsPage"
          :icon="h(KeyOutlined)"
          style="margin-left:8px"
        >{{t('deployPipeline.manageVars')}}</a-button>
      </div>

      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deletePipeline(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">{{t('deployPipeline.updatePipeline')}}</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </div>
      </template>
    </ZTable>
  </div>
</template>
<script setup>
import {
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  KeyOutlined
} from "@ant-design/icons-vue";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, h, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listPipelineRequest,
  deletePipelineRequest
} from "@/api/app/pipelineApi";
import { usePipelineStore } from "@/pinia/pipelineStore";
import { Modal, message } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const pipelineStore = usePipelineStore();
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
// 数据项
const columns = [
  {
    i18nTitle: "deployPipeline.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "deployPipeline.operation",
    dataIndex: "operation",
    key: "operation"
  }
];
// 删除流水线
const deletePipeline = item => {
  Modal.confirm({
    title: `${t("deployPipeline.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deletePipelineRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listPipeline();
      });
    },
    onCancel() {}
  });
};
// 流水列表
const listPipeline = () => {
  listPipelineRequest({
    appId: route.params.appId,
    env: selectedEnv.value
  }).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/create?env=${selectedEnv.value}`
  );
};
// 跳转管理变量页面
const gotoVarsPage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/vars/${selectedEnv.value}`
  );
};
// 跳转编辑页面
const gotoUpdatePage = item => {
  pipelineStore.id = item.id;
  pipelineStore.name = item.name;
  pipelineStore.env = item.env;
  pipelineStore.config = item.config;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/${item.id}/update`
  );
};
// 环境变化
const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/list/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  listPipeline();
};
</script>
<style scoped>
</style>