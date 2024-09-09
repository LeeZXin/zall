<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px;font-size:14px" class="flex-between">
      <div class="flex-center">
        <a-button
          type="primary"
          :icon="h(PlusOutlined)"
          @click="gotoCreatePage"
          style="margin-right:10px"
        >{{t('alertConfig.createConfig')}}</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'isEnabled'">
          <a-switch :checked="dataItem[dataIndex]" @click="enableOrDisableAlertConfig(dataItem)" />
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon" @click="deleteAlertConfig(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">{{t('alertConfig.updateConfig')}}</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </template>
      </template>
    </ZTable>
    <a-pagination
      v-model:current="dataPage.current"
      :total="dataPage.totalCount"
      show-less-items
      :pageSize="dataPage.pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listAlertConfig()"
    />
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { ref, createVNode, h, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  DeleteOutlined,
  EditOutlined,
  ExclamationCircleOutlined,
  EllipsisOutlined,
  PlusOutlined
} from "@ant-design/icons-vue";
import {
  listAlertConfigRequest,
  enableAlertConfigRequest,
  disableAlertConfigRequest,
  deleteAlertConfigRequest
} from "@/api/app/alertApi";
import { Modal, message } from "ant-design-vue";
import { useAlertConfigStore } from "@/pinia/alertConfigStore";
import EnvSelector from "@/components/app/EnvSelector";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const alertConfigStore = useAlertConfigStore();
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
// 分页
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
// 选择的环境
const selectedEnv = ref();
// 数据项
const columns = [
  {
    i18nTitle: "alertConfig.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "alertConfig.creator",
    dataIndex: "creator",
    key: "creator"
  },
  {
    i18nTitle: "alertConfig.isEnabled",
    dataIndex: "isEnabled",
    key: "isEnabled"
  },
  {
    i18nTitle: "alertConfig.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 获取配置列表
const listAlertConfig = () => {
  listAlertConfigRequest({
    appId: route.params.appId,
    pageNum: dataPage.current,
    env: selectedEnv.value
  }).then(res => {
    dataPage.totalCount = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};
// 删除配置
const deleteAlertConfig = item => {
  Modal.confirm({
    title: `${t('alertConfig.confirmDelete')} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteAlertConfigRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchAlertConfig();
      });
    },
    onCancel() {}
  });
};
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/alertConfig/create?env=${selectedEnv.value}`
  );
};
// 跳转编辑页面
const gotoUpdatePage = item => {
  alertConfigStore.id = item.id;
  alertConfigStore.name = item.name;
  alertConfigStore.content = item.content;
  alertConfigStore.intervalSec = item.intervalSec;
  alertConfigStore.isEnabled = item.isEnabled;
  alertConfigStore.env = item.env;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/alertConfig/${item.id}/update`
  );
};
// 启动或停用
const enableOrDisableAlertConfig = item => {
  if (item.isEnabled) {
    disableAlertConfigRequest(item.id).then(() => {
      message.success(t("operationSuccess"));
      item.isEnabled = false;
    });
  } else {
    enableAlertConfigRequest(item.id).then(() => {
      message.success(t("operationSuccess"));
      item.isEnabled = true;
    });
  }
};
// 环境变化
const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/alertConfig/list/${e.newVal}`
  );
  dataPage.current = 1;
  listAlertConfig();
};
const searchAlertConfig = () => {
  dataPage.current = 1;
  listAlertConfig();
};
</script>
<style scoped>
</style>