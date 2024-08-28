<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px;font-size:14px" class="flex-between">
      <div class="flex-center">
        <a-input
          placeholder="搜索名称"
          style="width: 240px;margin-right:10px"
          v-model:value="searchName"
          @pressEnter="()=>searchTimer()"
        >
          <template #suffix>
            <SearchOutlined />
          </template>
        </a-input>
        <a-button
          type="primary"
          :icon="h(PlusOutlined)"
          @click="gotoCreatePage"
          style="margin-right:10px"
          v-if="teamStore.perm?.canManageTimer"
        >创建定时任务</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'isEnabled'">
          <a-switch :checked="dataItem[dataIndex]" @click="enableOrDisableTimer(dataItem)" />
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
                <li @click="viewLog(dataItem)">
                  <eye-outlined />
                  <span style="margin-left:4px">查看日志</span>
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
      @change="()=>listTimer()"
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
  PlusOutlined,
  PlayCircleOutlined,
  EyeOutlined,
  SearchOutlined
} from "@ant-design/icons-vue";
import {
  listTimerRequest,
  enableTimerRequest,
  disableTimerRequest,
  deleteTimerRequest,
  triggerTimerTaskRequest
} from "@/api/team/timerApi";
import { Modal, message } from "ant-design-vue";
import { useTimerTaskStore } from "@/pinia/timerTaskStore";
import EnvSelector from "@/components/app/EnvSelector";
import { useTeamStore } from "@/pinia/teamStore";
const teamStore = useTeamStore();
const timerTaskStore = useTimerTaskStore();
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
const selectedEnv = ref();
const searchName = ref("");
const columns = [
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
const listTimer = () => {
  listTimerRequest({
    teamId: parseInt(route.params.teamId),
    pageNum: dataPage.current,
    name: searchName.value,
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
const deleteTimerTask = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteTimerRequest(item.id).then(() => {
        message.success("删除成功");
        searchTimer();
      });
    },
    onCancel() {}
  });
};
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/timer/create?env=${selectedEnv.value}`
  );
};
const triggerTimerTask = item => {
  Modal.confirm({
    title: `你确定要触发${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      triggerTimerTaskRequest(item.id).then(() => {
        message.success("触发成功");
      });
    }
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
  router.push(`/team/${route.params.teamId}/timer/${item.id}/update`);
};
const enableOrDisableTimer = item => {
  if (item.isEnabled) {
    disableTimerRequest(item.id).then(() => {
      message.success("关闭成功");
      item.isEnabled = false;
    });
  } else {
    enableTimerRequest(item.id).then(() => {
      message.success("启动成功");
      item.isEnabled = true;
    });
  }
};
const viewLog = item => {
  timerTaskStore.id = item.id;
  timerTaskStore.name = item.name;
  timerTaskStore.cronExp = item.cronExp;
  timerTaskStore.task = item.task;
  timerTaskStore.teamId = item.teamId;
  timerTaskStore.isEnabled = item.isEnabled;
  timerTaskStore.env = item.env;
  router.push(`/team/${route.params.teamId}/timer/${item.id}/log`);
};
const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(`/team/${route.params.teamId}/timer/list/${e.newVal}`);
  searchTimer();
};
const searchTimer = () => {
  dataPage.current = 1;
  listTimer();
}
</script>
<style scoped>
</style>