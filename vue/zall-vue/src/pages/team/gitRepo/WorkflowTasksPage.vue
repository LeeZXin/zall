<template>
  <div style="padding:10px">
    <div class="workflow-name">{{workflowStore.name}}</div>
    <div class="workflow-desc">{{workflowStore.desc}}</div>
    <ZTable :columns="columns" :dataSource="taskList" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <div v-if="dataIndex === 'operator'" class="flex-center">
          <ZAvatar
            :url="dataItem.operator?.avatarUrl"
            :name="dataItem.operator?.name"
            :account="dataItem.operator?.account"
            :showName="true"
          />
        </div>
        <span v-else-if="dataIndex === 'created'">{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex === 'triggerType'">{{t(triggerTypeMap[dataItem[dataIndex]])}}</span>
        <template v-else-if="dataIndex === 'prIndex'">
          <PrIndexTag
            v-if="dataItem[dataIndex]"
            :repoId="route.params.repoId"
            :prIndex="dataItem[dataIndex]"
            :teamId="route.params.teamId"
          />
        </template>
        <TaskStatusTag :status="dataItem[dataIndex]" v-else-if="dataIndex === 'taskStatus'" />
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoTaskDetail(dataItem)">
                  <EyeOutlined />
                  <span style="margin-left:4px">{{t('gitWorkflow.viewDetail')}}</span>
                </li>
                <li v-if="dataItem.taskStatus === 'running'" @click="killTask(dataItem.id)">
                  <CloseOutlined />
                  <span style="margin-left:4px">{{t('gitWorkflow.killWorkflow')}}</span>
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
    <a-pagination
      v-model:current="dataPage.current"
      :total="dataPage.totalCount"
      show-less-items
      :pageSize="dataPage.pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listTask()"
    />
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import PrIndexTag from "@/components/git/PrIndexTag";
import TaskStatusTag from "@/components/git/WorkflowTaskStatusTag";
import ZTable from "@/components/common/ZTable";
import {
  listTaskRequest,
  killWorkflowTaskRequest
} from "@/api/git/workflowApi";
import {
  EyeOutlined,
  CloseOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import { ref, createVNode, onUnmounted, reactive } from "vue";
import { useRouter, useRoute } from "vue-router";
import { readableTimeComparingNow } from "@/utils/time";
import { Modal, message } from "ant-design-vue";
import { useWorkflowStore } from "@/pinia/workflowStore";
import { useWorkflowTaskStore } from "@/pinia/workflowTaskStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const workflowStore = useWorkflowStore();
const taskStore = useWorkflowTaskStore();
// 分页数据
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
const route = useRoute();
const router = useRouter();
// 任务列表
const taskList = ref([]);
const triggerTypeMap = {
  1: "gitWorkflow.automatic",
  2: "gitWorkflow.manual"
};
// 数据项
const columns = [
  {
    i18nTitle: "gitWorkflow.triggerType",
    dataIndex: "triggerType",
    key: "triggerType"
  },
  {
    i18nTitle: "gitWorkflow.taskStatus",
    dataIndex: "taskStatus",
    key: "taskStatus"
  },
  {
    i18nTitle: "gitWorkflow.operator",
    dataIndex: "operator",
    key: "operator"
  },
  {
    i18nTitle: "gitWorkflow.createTime",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "gitWorkflow.branch",
    dataIndex: "branch",
    key: "branch"
  },
  {
    i18nTitle: "gitWorkflow.pullRequest",
    dataIndex: "prIndex",
    key: "prIndex"
  },
  {
    i18nTitle: "gitWorkflow.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 跳转任务详情页
const gotoTaskDetail = item => {
  taskStore.id = item.id;
  taskStore.triggerType = item.triggerType;
  taskStore.operator = item.operator;
  taskStore.created = item.created;
  taskStore.branch = item.branch;
  taskStore.prId = item.prId;
  taskStore.yamlContent = item.yamlContent;
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/${route.params.workflowId}/${item.id}/steps`
  );
};
// 获取任务列表
const listTask = () => {
  listTaskRequest({
    workflowId: route.params.workflowId,
    pageNum: dataPage.current
  }).then(res => {
    dataPage.totalCount = res.totalCount;
    taskList.value = res.data.map(item => {
      return {
        ...item,
        key: item.id
      };
    });
  });
};
// 停止任务
const killTask = id => {
  Modal.confirm({
    title: `${t("gitWorkflow.confirmKill")} ${workflowStore.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      killWorkflowTaskRequest(id).then(() => {
        message.success(t("operationSuccess"));
        listTask();
      });
    },
    onCancel() {}
  });
};
// 定时刷新任务列表数据
const listInterval = setInterval(listTask, 5000);
const clearListInterval = () => {
  if (listInterval) {
    clearInterval(listInterval);
  }
};
onUnmounted(() => {
  clearListInterval();
});
if (workflowStore.id === 0) {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/list`
  );
} else {
  listTask();
}
</script>
<style scoped>
.workflow-name {
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: bold;
}
.workflow-desc {
  font-size: 14px;
  color: gray;
  margin-bottom: 10px;
}
</style>