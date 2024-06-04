<template>
  <div style="padding:14px">
    <div style="margin-bottom:10px;font-size:14px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建定时任务</a-button>
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
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'isEnabled'">
          <a-switch :checked="dataItem[dataIndex]" @click="enableOrDisableTask(dataItem)" />
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon" @click="deleteTimerTask(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>删除定时任务</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <edit-outlined />
                  <span style="margin-left:4px">编辑定时任务</span>
                </li>
                <li @click="triggerTimerTask(dataItem)">
                  <play-circle-outlined />
                  <span style="margin-left:4px">手动触发任务</span>
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
      v-model:current="currPage"
      :total="totalCount"
      show-less-items
      :pageSize="pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listTimerTask()"
    />
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { ref, createVNode, h, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  DeleteOutlined,
  EditOutlined,
  ExclamationCircleOutlined,
  EllipsisOutlined,
  PlusOutlined,
  PlayCircleOutlined
} from "@ant-design/icons-vue";
import {
  listTimerTaskRequest,
  enableTimerTaskRequest,
  disableTimerTaskRequest,
  deleteTimerTaskRequest,
  triggerTimerTaskRequest
} from "@/api/team/timerApi";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import { Modal, message } from "ant-design-vue";
import { useTimerTaskStore } from "@/pinia/timerTaskStore";
const timerTaskStore = useTimerTaskStore();
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const selectedEnv = ref();
const envList = ref([]);
const columns = ref([
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "cron表达式",
    dataIndex: "cronExp",
    key: "cronExp"
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
]);
const listTimerTask = () => {
  listTimerTaskRequest(
    {
      teamId: parseInt(route.params.teamId),
      pageNum: currPage.value,
      env: selectedEnv.value
    },
    selectedEnv.value
  ).then(res => {
    totalCount.value = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};
const getEnvList = () => {
  getEnvCfgRequest().then(res => {
    envList.value = res.data.map(item => {
      return {
        value: item,
        label: item
      };
    });
    if (res.data.length > 0) {
      selectedEnv.value = res.data[0];
    }
  });
};
const deleteTimerTask = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteTimerTaskRequest(item.id, item.env).then(() => {
        message.success("删除成功");
        listTimerTask();
      });
    },
    onCancel() {}
  });
};
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/timerTask/create`);
};
const triggerTimerTask = item => {
  Modal.confirm({
    title: `你确定要触发${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      triggerTimerTaskRequest(item.id, item.env).then(() => {
        message.success("触发成功");
      });
    },
    onCancel() {}
  });
};
const gotoUpdatePage = item => {
  timerTaskStore.id = item.id;
  timerTaskStore.name = item.name;
  timerTaskStore.cronExp = item.cronExp;
  timerTaskStore.task = item.task;
  timerTaskStore.teamId = item.teamId;
  timerTaskStore.isEnabled = item.isEnabled;
  timerTaskStore.env = item.env;
  router.push(`/team/${route.params.teamId}/timerTask/${item.id}/update`);
};
const enableOrDisableTask = item => {
  if (item.isEnabled) {
    disableTimerTaskRequest(item.id, item.env).then(() => {
      message.success("关闭成功");
      item.isEnabled = false;
    });
  } else {
    enableTimerTaskRequest(item.id, item.env).then(() => {
      message.success("启动成功");
      item.isEnabled = true;
    });
  }
};
watch(
  () => selectedEnv.value,
  () => {
    listTimerTask();
  }
);
getEnvList();
</script>
<style scoped>
</style>