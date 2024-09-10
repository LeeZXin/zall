<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button
        type="primary"
        :icon="h(PlusOutlined)"
        @click="gotoCreatePage"
      >{{t('discoverySource.createSource')}}</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteDiscoverySource(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">{{t('discoverySource.updateSource')}}</span>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const discoverySourceStore = useDiscoverySourceStore();
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
// 数据项
const columns = [
  {
    i18nTitle: "discoverySource.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "discoverySource.endpoints",
    dataIndex: "endpoints",
    key: "endpoints"
  },
  {
    i18nTitle: "discoverySource.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 删除来源
const deleteDiscoverySource = item => {
  Modal.confirm({
    title: `${t("discoverySource.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteDiscoverySourceRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listDiscoverySource();
      });
    },
    onCancel() {}
  });
};
// 获取列表
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
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(`/sa/discoverySource/create?env=${selectedEnv.value}`);
};
// 跳转编辑页面
const gotoUpdatePage = item => {
  discoverySourceStore.id = item.id;
  discoverySourceStore.name = item.name;
  discoverySourceStore.env = item.env;
  discoverySourceStore.endpoints = item.endpoints;
  discoverySourceStore.username = item.username;
  discoverySourceStore.password = item.password;
  router.push(`/sa/discoverySource/${item.id}/update`);
};
// 环境变化
const onEnvChange = e => {
  router.replace(`/sa/discoverySource/list/${e.newVal}`);
  selectedEnv.value = e.newVal;
  listDiscoverySource();
};
</script>
<style scoped>
</style>