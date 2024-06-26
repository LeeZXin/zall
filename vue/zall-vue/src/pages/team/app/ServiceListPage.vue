<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建服务</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env"/>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex === 'serviceType'">{{t(`service.${dataItem[dataIndex]}`)}}</span>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteService(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Service</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <edit-outlined />
                  <span style="margin-left:4px">编辑服务</span>
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
import { listServiceRequest, deleteServiceRequest } from "@/api/app/serviceApi";
import { useServiceStore } from "@/pinia/serviceStore";
import { Modal, message } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const serviceStore = useServiceStore();
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
    title: "服务类型",
    dataIndex: "serviceType",
    key: "serviceType"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);

const deleteService = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteServiceRequest(item.id, item.env).then(() => {
        message.success("删除成功");
        listService();
      });
    },
    onCancel() {}
  });
};

const listService = () => {
  listServiceRequest(
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
    `/team/${route.params.teamId}/app/${route.params.appId}/service/create?env=${selectedEnv.value}`
  );
};

const gotoUpdatePage = item => {
  serviceStore.id = item.id;
  serviceStore.name = item.name;
  serviceStore.env = item.env;
  serviceStore.config = item.config;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/service/${item.id}/update`
  );
};

const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/service/list/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  listService();
};
</script>
<style scoped>
</style>