<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px;font-size:14px" class="flex-between">
      <div class="flex-center">
        <a-button
          type="primary"
          :icon="h(PlusOutlined)"
          @click="gotoCreatePage"
          style="margin-right:10px"
        >创建监控告警</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'isEnabled'">
          <a-switch :checked="dataItem[dataIndex]" @click="enableOrDisableAlertConfig(dataItem)" />
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon" @click="deleteAlertConfigTask(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>删除监控告警</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">编辑监控告警</span>
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
const alertConfigStore = useAlertConfigStore();
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
const selectedEnv = ref();
const columns = [
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "创建人",
    dataIndex: "creator",
    key: "creator"
  },
  {
    title: "是否启动",
    dataIndex: "isEnabled",
    key: "isEnabled"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];
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
const deleteAlertConfigTask = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteAlertConfigRequest(item.id).then(() => {
        message.success("删除成功");
        searchAlertConfig();
      });
    },
    onCancel() {}
  });
};
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/alertConfig/create?env=${selectedEnv.value}`
  );
};
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
const enableOrDisableAlertConfig = item => {
  if (item.isEnabled) {
    disableAlertConfigRequest(item.id).then(() => {
      message.success("关闭成功");
      item.isEnabled = false;
    });
  } else {
    enableAlertConfigRequest(item.id).then(() => {
      message.success("启动成功");
      item.isEnabled = true;
    });
  }
};
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