<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建注册中心来源</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteDiscoverySource(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Source</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <edit-outlined />
                  <span style="margin-left:4px">编辑注册中心来源</span>
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
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, h, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listDiscoverySourceRequest,
  deleteDiscoverySourceRequest
} from "@/api/app/discoveryApi";
import { useDiscoverySourceStore } from "@/pinia/discoverySourceStore";
import { Modal, message } from "ant-design-vue";
const discoverySourceStore = useDiscoverySourceStore();
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
const columns = [
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "endpoints",
    dataIndex: "endpoints",
    key: "endpoints"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];

const deleteDiscoverySource = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteDiscoverySourceRequest(item.id).then(() => {
        message.success("删除成功");
        listDiscoverySource();
      });
    },
    onCancel() {}
  });
};

const listDiscoverySource = () => {
  listDiscoverySourceRequest({
    env: selectedEnv.value
  }).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item,
        endpoints: item.endpoints ? item.endpoints.join(";") : ""
      };
    });
  });
};

const gotoCreatePage = () => {
  router.push(`/sa/discoverySource/create?env=${selectedEnv.value}`);
};

const gotoUpdatePage = item => {
  discoverySourceStore.id = item.id;
  discoverySourceStore.name = item.name;
  discoverySourceStore.env = item.env;
  discoverySourceStore.endpoints = item.endpoints;
  discoverySourceStore.username = item.username;
  discoverySourceStore.password = item.password;
  router.push(`/sa/discoverySource/${item.id}/update`);
};

const onEnvChange = e => {
  router.replace(`/sa/discoverySource/list/${e.newVal}`);
  selectedEnv.value = e.newVal;
  listDiscoverySource();
};
</script>
<style scoped>
</style>