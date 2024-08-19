<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px;font-size:14px" class="flex-between">
      <div class="flex-center">
        <a-input
          placeholder="搜索名称"
          style="width: 240px;margin-right:10px"
          v-model:value="searchName"
          @pressEnter="()=>listTimerTask()"
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
        <a-button
          type="primary"
          :icon="h(SettingOutlined)"
          @click="showBindModal"
          v-if="teamStore.perm?.canManageNotifyTpl"
        >任务失败通知模板</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
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
  <a-modal v-model:open="bindModal.open" title="任务失败通知模板" @ok="handleBindModalOk">
    <div>
      <div style="font-size:12px;margin-bottom:3px">已选环境</div>
      <div>{{selectedEnv}}</div>
    </div>
    <div style="margin-top: 10px">
      <div style="font-size:12px;margin-bottom:3px">任务失败通知模板</div>
      <a-select
        style="width: 100%"
        placeholder="请选择"
        v-model:value="bindModal.selectTpl"
        :options="bindModal.tplList"
        show-search
        :filter-option="filterSourceListOption"
      />
    </div>
  </a-modal>
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
  SearchOutlined,
  SettingOutlined
} from "@ant-design/icons-vue";
import {
  listTimerTaskRequest,
  enableTimerTaskRequest,
  disableTimerTaskRequest,
  deleteTimerTaskRequest,
  triggerTimerTaskRequest,
  getFailedTaskNotifyTplRequest,
  bindFailedTaskNotifyTplRequest
} from "@/api/team/timerApi";
import { listAllTplByTeamIdRequest } from "@/api/notify/notifyApi";
import { Modal, message } from "ant-design-vue";
import { useTimerTaskStore } from "@/pinia/timerTaskStore";
import EnvSelector from "@/components/app/EnvSelector";
import { useTeamStore } from "@/pinia/teamStore";
const teamStore = useTeamStore();
const timerTaskStore = useTimerTaskStore();
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const selectedEnv = ref();
const searchName = ref("");
const bindModal = reactive({
  open: false,
  selectTpl: undefined,
  tplList: [
    {
      value: 0,
      label: "不通知"
    }
  ]
});
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
const listTimerTask = () => {
  listTimerTaskRequest({
    teamId: parseInt(route.params.teamId),
    pageNum: currPage.value,
    name: searchName.value,
    env: selectedEnv.value
  }).then(res => {
    totalCount.value = res.totalCount;
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
      deleteTimerTaskRequest(item.id).then(() => {
        message.success("删除成功");
        listTimerTask();
      });
    },
    onCancel() {}
  });
};
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/timerTask/create?env=${selectedEnv.value}`
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
    disableTimerTaskRequest(item.id).then(() => {
      message.success("关闭成功");
      item.isEnabled = false;
    });
  } else {
    enableTimerTaskRequest(item.id).then(() => {
      message.success("启动成功");
      item.isEnabled = true;
    });
  }
};
// 下拉框搜索过滤
const filterSourceListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
const viewLog = item => {
  timerTaskStore.id = item.id;
  timerTaskStore.name = item.name;
  timerTaskStore.cronExp = item.cronExp;
  timerTaskStore.task = item.task;
  timerTaskStore.teamId = item.teamId;
  timerTaskStore.isEnabled = item.isEnabled;
  timerTaskStore.env = item.env;
  router.push(`/team/${route.params.teamId}/timerTask/${item.id}/log`);
};
const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(`/team/${route.params.teamId}/timerTask/list/${e.newVal}`);
  listTimerTask();
};
// 展示绑定modal
const showBindModal = () => {
  bindModal.selectTpl = null;
  if (bindModal.tplList.length === 1) {
    listAllTplByTeamIdRequest(route.params.teamId).then(res => {
      bindModal.tplList = bindModal.tplList.concat(
        res.data.map(item => {
          return {
            value: item.id,
            label: item.name
          };
        })
      );
      getFailedTaskNotifyTplRequest({
        teamId: parseInt(route.params.teamId),
        env: selectedEnv.value
      }).then(res => {
        bindModal.selectTpl = res.data;
        bindModal.open = true;
      });
    });
  } else {
    getFailedTaskNotifyTplRequest({
      teamId: parseInt(route.params.teamId),
      env: selectedEnv.value
    }).then(res => {
      bindModal.selectTpl = res.data;
      bindModal.open = true;
    });
  }
};
// 绑定失败通知模板
const handleBindModalOk = () => {
  bindFailedTaskNotifyTplRequest({
    teamId: parseInt(route.params.teamId),
    env: selectedEnv.value,
    tplId: bindModal.selectTpl
  }).then(() => {
    message.success("操作成功");
    bindModal.open = false;
  });
};
</script>
<style scoped>
</style>