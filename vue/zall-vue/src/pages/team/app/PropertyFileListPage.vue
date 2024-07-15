<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建配置文件</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <div class="body">
      <ZTable :columns="columns" :dataSource="dataSource">
        <template #bodyCell="{dataIndex, dataItem}">
          <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <div v-else>
            <div class="op-icon" @click="deleteFile(dataItem)">
              <a-tooltip placement="top">
                <template #title>
                  <span>Delete File</span>
                </template>
                <delete-outlined />
              </a-tooltip>
            </div>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <li @click="gotoHistoryListPage(dataItem)">
                    <file-text-outlined />
                    <span style="margin-left:4px">版本列表</span>
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
  </div>
</template>
<script setup>
import {
  DeleteOutlined,
  FileTextOutlined,
  PlusOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref, h, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import { listPropertyFileRequest, deletePropertyFileRequest } from "@/api/app/propertyApi";
import { usePropertyFileStore } from "@/pinia/propertyFileStore";
import EnvSelector from "@/components/app/EnvSelector";
import { Modal, message } from "ant-design-vue";
const propertyFileStore = usePropertyFileStore();
const selectedEnv = ref("");
const route = useRoute();
const router = useRouter();
const dataSource = ref([]);

const columns = ref([
  {
    title: "配置文件",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);

const listPropertyFile = () => {
  listPropertyFileRequest(
    {
      appId: route.params.appId,
      env: selectedEnv.value
    },
    selectedEnv.value
  ).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.name,
        ...item
      };
    });
  });
};

const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/create?env=${selectedEnv.value}`
  );
};

const gotoHistoryListPage = item => {
  propertyFileStore.id = item.id;
  propertyFileStore.name = item.name;
  propertyFileStore.env = item.env;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${item.id}/history/list`
  );
};

const deleteFile = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deletePropertyFileRequest(item.id, item.env).then(() => {
        message.success("删除成功");
        listPropertyFile();
      });
    },
    onCancel() {}
  });
};

const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/list/${e.newVal}`
  );
  listPropertyFile();
};
</script>
<style scoped>
.body {
  width: 100%;
  height: calc(100% - 42px);
}
</style>