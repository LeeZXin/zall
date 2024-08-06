<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建配置来源</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deletePropertySource(dataItem)">
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
                  <span style="margin-left:4px">编辑配置来源</span>
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
  listPropertySourceRequest,
  deletePropertySourceRequest
} from "@/api/app/propertyApi";
import { usePropertySourceStore } from "@/pinia/propertySourceStore";
import { Modal, message } from "ant-design-vue";
const propertySourceStore = usePropertySourceStore();
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

const deletePropertySource = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deletePropertySourceRequest(item.id).then(() => {
        message.success("删除成功");
        listPropertySource();
      });
    }
  });
};

const listPropertySource = () => {
  listPropertySourceRequest({
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
  router.push(`/sa/propertySource/create?env=${selectedEnv.value}`);
};

const gotoUpdatePage = item => {
  propertySourceStore.id = item.id;
  propertySourceStore.name = item.name;
  propertySourceStore.env = item.env;
  propertySourceStore.endpoints = item.endpoints;
  propertySourceStore.username = item.username;
  propertySourceStore.password = item.password;
  router.push(`/sa/propertySource/${item.id}/update`);
};

const onEnvChange = e => {
  router.replace(`/sa/propertySource/list/${e.newVal}`);
  selectedEnv.value = e.newVal;
  listPropertySource();
};
</script>
<style scoped>
</style>