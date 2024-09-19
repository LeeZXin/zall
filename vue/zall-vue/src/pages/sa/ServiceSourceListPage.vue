<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button
        type="primary"
        :icon="h(PlusOutlined)"
        @click="gotoCreatePage"
      >{{t('serviceSource.createSource')}}</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteServiceSource(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">{{t('serviceSource.updateSource')}}</span>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const serviceSourceStore = useServiceSourceStore();
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
// 数据项
const columns = [
  {
    i18nTitle: "serviceSource.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "serviceSource.operation",
    dataIndex: "operation",
    key: "operation",
    fixed: "right",
    width: 130
  }
];
// 删除来源
const deleteServiceSource = item => {
  Modal.confirm({
    title: `${t("serviceSource.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteServiceSourceRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listServiceSource();
      });
    },
    onCancel() {}
  });
};
// 获取列表
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
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(`/sa/serviceSource/create?env=${selectedEnv.value}`);
};
// 跳转编辑页面
const gotoUpdatePage = item => {
  serviceSourceStore.id = item.id;
  serviceSourceStore.name = item.name;
  serviceSourceStore.env = item.env;
  serviceSourceStore.datasource = item.datasource;
  router.push(`/sa/serviceSource/${item.id}/update`);
};
// 环境变化
const onEnvChange = e => {
  router.replace(`/sa/serviceSource/list/${e.newVal}`);
  selectedEnv.value = e.newVal;
  listServiceSource();
};
</script>
<style scoped>
</style>