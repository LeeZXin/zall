<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <div>
        <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建流水线</a-button>
        <a-button
          type="primary"
          @click="gotoVarsPage"
          :icon="h(KeyOutlined)"
          style="margin-left:8px"
        >管理变量</a-button>
      </div>

      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deletePipeline(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Pipeline</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <edit-outlined />
                  <span style="margin-left:4px">编辑流水线</span>
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
const pipelineStore = usePipelineStore();
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
const columns = ref([
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);

const deletePipeline = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deletePipelineRequest(item.id, item.env).then(() => {
        message.success("删除成功");
        listPipeline();
      });
    },
    onCancel() {}
  });
};

const listPipeline = () => {
  listPipelineRequest(
    {
      appId: route.params.appId,
      env: selectedEnv.value
    },
    selectedEnv.value
  ).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};

const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/create?env=${selectedEnv.value}`
  );
};

const gotoVarsPage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/vars/${selectedEnv.value}`
  );
};

const gotoUpdatePage = item => {
  pipelineStore.id = item.id;
  pipelineStore.name = item.name;
  pipelineStore.env = item.env;
  pipelineStore.config = item.config;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/pipeline/${item.id}/update`
  );
};

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