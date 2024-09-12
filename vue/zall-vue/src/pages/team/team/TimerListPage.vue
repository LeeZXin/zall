<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px;font-size:14px" class="flex-between">
      <div class="flex-center">
        <a-input
          :placeholder="t('timerTask.searchName')"
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
        >{{t('timerTask.createTimer')}}</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <div v-if="dataIndex === 'creator'" class="flex-center">
          <ZAvatar
            :url="dataItem.creator?.avatarUrl"
            :name="dataItem.creator?.name"
            :showName="true"
          />
        </div>
        <template v-else-if="dataIndex === 'isEnabled'">
          <a-switch :checked="dataItem[dataIndex]" @click="enableOrDisableTimer(dataItem)" />
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon" @click="deleteTimerTask(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <edit-outlined />
                  <span style="margin-left:4px">{{t('timerTask.updateTimer')}}</span>
                </li>
                <li @click="triggerTimerTask(dataItem)">
                  <play-circle-outlined />
                  <span style="margin-left:4px">{{t('timerTask.triggerTimer')}}</span>
                </li>
                <li @click="viewLog(dataItem)">
                  <eye-outlined />
                  <span style="margin-left:4px">{{t('timerTask.viewLog')}}</span>
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
import ZAvatar from "@/components/user/ZAvatar";
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const teamStore = useTeamStore();
const timerTaskStore = useTimerTaskStore();
const router = useRouter();
const route = useRoute();
// 数据
const dataSource = ref([]);
// 分页
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
// 当前环境
const selectedEnv = ref();
// 搜索名称
const searchName = ref("");
// 数据项
const columns = [
  {
    i18nTitle: "timerTask.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "timerTask.cronExp",
    dataIndex: "cronExp",
    key: "cronExp"
  },
  {
    i18nTitle: "timerTask.creator",
    dataIndex: "creator",
    key: "creator"
  },
  {
    i18nTitle: "timerTask.isEnabled",
    dataIndex: "isEnabled",
    key: "isEnabled"
  },
  {
    i18nTitle: "timerTask.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 定时任务列表
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
// 删除定时任务
const deleteTimerTask = item => {
  Modal.confirm({
    title: `${t("timerTask.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteTimerRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchTimer();
      });
    },
    onCancel() {}
  });
};
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/timer/create?env=${selectedEnv.value}`
  );
};
// 手动触发任务
const triggerTimerTask = item => {
  Modal.confirm({
    title: `${t("timerTask.confirmTrigger")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      triggerTimerTaskRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
      });
    }
  });
};
// 跳转编辑页面
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
// 启动或停用定时任务
const enableOrDisableTimer = item => {
  if (item.isEnabled) {
    disableTimerRequest(item.id).then(() => {
      message.success(t("operationSuccess"));
      item.isEnabled = false;
    });
  } else {
    enableTimerRequest(item.id).then(() => {
      message.success(t("operationSuccess"));
      item.isEnabled = true;
    });
  }
};
// 跳转日志页面
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
// 环境变化
const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(`/team/${route.params.teamId}/timer/list/${e.newVal}`);
  searchTimer();
};
// 搜索定时任务
const searchTimer = () => {
  dataPage.current = 1;
  listTimer();
};
</script>
<style scoped>
</style>