<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px;text-align:right">
      <div>
        <span style="margin-right:6px">环境:</span>
        <a-select
          style="width: 200px"
          placeholder="选择环境"
          v-model:value="selectedEnv"
          :options="envList"
        />
      </div>
    </div>
    <div class="body">
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
import ZTable from "@/components/common/ZTable";
import { ref, h, watch, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  listServiceRequest,
  deleteServiceRequest,
  enableServiceRequest,
  disableServiceRequest
} from "@/api/app/serviceApi";
import { useServiceStore } from "@/pinia/serviceStore";
import { Modal, message } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const serviceStore = useServiceStore();
const selectedEnv = ref("");
const envList = ref([]);
const route = useRoute();
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
    title: "是否启用",
    dataIndex: "isEnabled",
    key: "isEnabled"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);

const getEnvList = () => {
  getEnvCfgRequest().then(res => {
    envList.value = res.data.map(item => {
      return {
        value: item,
        label: item
      };
    });
    if (route.params.env && res.data?.includes(route.params.env)) {
      selectedEnv.value = route.params.env;
    } else if (res.data.length > 0) {
      selectedEnv.value = res.data[0];
    }
  });
};

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

const gotoUpdatePage = item => {
  serviceStore.id = item.id;
  serviceStore.name = item.name;
  serviceStore.env = item.env;
  serviceStore.config = item.config;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/service/${item.id}/update`
  );
};

getEnvList();

watch(
  () => selectedEnv.value,
  newVal => {
    router.replace(
      `/team/${route.params.teamId}/app/${route.params.appId}/service/list/${newVal}`
    );
    listService();
  }
);
</script>
<style scoped>
.body {
  width: 100%;
  height: calc(100% - 42px);
}
</style>