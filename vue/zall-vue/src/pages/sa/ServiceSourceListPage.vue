<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建服务来源</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteServiceSource(dataItem)">
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
                  <span style="margin-left:4px">编辑服务来源</span>
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
  listServiceSourceRequest,
  deleteServiceSourceRequest
} from "@/api/app/serviceApi";
import { useServiceSourceStore } from "@/pinia/serviceSourceStore";
import { Modal, message } from "ant-design-vue";
const serviceSourceStore = useServiceSourceStore();
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
    title: "host",
    dataIndex: "host",
    key: "host"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];

const deleteServiceSource = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteServiceSourceRequest(item.id).then(() => {
        message.success("删除成功");
        listServiceSource();
      });
    },
    onCancel() {}
  });
};

const listServiceSource = () => {
  listServiceSourceRequest({
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

const gotoCreatePage = () => {
  router.push(`/sa/serviceSource/create?env=${selectedEnv.value}`);
};

const gotoUpdatePage = item => {
  serviceSourceStore.id = item.id;
  serviceSourceStore.name = item.name;
  serviceSourceStore.env = item.env;
  serviceSourceStore.host = item.host;
  serviceSourceStore.apiKey = item.apiKey;
  router.push(`/sa/serviceSource/${item.id}/update`);
};

const onEnvChange = e => {
  router.replace(`/sa/serviceSource/list/${e.newVal}`);
  selectedEnv.value = e.newVal;
  listServiceSource();
};
</script>
<style scoped>
</style>