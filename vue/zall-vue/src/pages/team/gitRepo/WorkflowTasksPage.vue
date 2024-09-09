<template>
  <div style="padding:10px">
    <div class="workflow-name">{{workflowStore.name}}</div>
    <div class="workflow-desc">{{workflowStore.desc}}</div>
    <ZTable
      :columns="columns"
      :dataSource="taskList"
      label="任务列表"
      style="margin-top:10px"
      :scroll="{x:1300}"
    >
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex === 'created'">{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex === 'triggerType'">{{triggerTypeMap[dataItem[dataIndex]]}}</span>
        <template v-else-if="dataIndex === 'prId'">
          <PrIdTag
            v-if="dataItem[dataIndex]"
            :repoId="route.params.repoId"
            :prId="dataItem[dataIndex]"
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
                  <eye-outlined />
                  <span style="margin-left:4px">查看详情</span>
                </li>
                <li v-if="dataItem.taskStatus === 'running'" @click="killTask(dataItem.id)">
                  <close-outlined />
                  <span style="margin-left:4px">停止任务</span>
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
      v-model:current="currPage"
      :total="totalCount"
      show-less-items
      :pageSize="pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>paginationChange()"
    />
  </div>
</template>
<script setup>
import PrIdTag from "@/components/git/PrIdTag";
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
import { ref, createVNode, onUnmounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { readableTimeComparingNow } from "@/utils/time";
import { Modal, message } from "ant-design-vue";
import { useWorkflowStore } from "@/pinia/workflowStore";
import { useWorkflowTaskStore } from "@/pinia/workflowTaskStore";
const workflowStore = useWorkflowStore();
const taskStore = useWorkflowTaskStore();
const totalCount = ref(0);
const pageSize = 10;
const currPage = ref(1);
const route = useRoute();
const router = useRouter();
const taskList = ref([]);
const triggerTypeMap = {
  1: "自动触发",
  2: "手动触发"
};
const listInterval = ref(null);
const columns = [
  {
    title: "触发方式",
    dataIndex: "triggerType",
    key: "triggerType"
  },
  {
    title: "任务状态",
    dataIndex: "taskStatus",
    key: "taskStatus"
  },
  {
    title: "操作人",
    dataIndex: "operator",
    key: "operator"
  },
  {
    title: "创建时间",
    dataIndex: "created",
    key: "created"
  },
  {
    title: "触发分支",
    dataIndex: "branch",
    key: "branch"
  },
  {
    title: "关联合并请求",
    dataIndex: "prId",
    key: "prId"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
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
const listTask = () => {
  listTaskRequest(route.params.workflowId, {
    pageNum: currPage.value,
    pageSize
  }).then(res => {
    totalCount.value = res.totalCount;
    let runningTask = res.data.find(item => {
      return item.taskStatus === 1 || item.taskStatus === 0;
    });
    if (runningTask) {
      if (!listInterval.value) {
        listInterval.value = setInterval(listTask, 5000);
      }
    } else {
      clearListInterval();
    }
    taskList.value = res.data.map(item => {
      return {
        ...item,
        key: item.id
      };
    });
  });
};
const killTask = id => {
  Modal.confirm({
    title: `你确定要停止${workflowStore.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      killWorkflowTaskRequest(id).then(() => {
        message.success("停止成功");
        listTask();
      });
    },
    onCancel() {}
  });
};
const clearListInterval = () => {
  if (listInterval.value) {
    clearInterval(listInterval.value);
    listInterval.value = null;
  }
};
const paginationChange = () => {
  clearListInterval();
  listTask();
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
}
.workflow-desc {
  font-size: 12px;
  color: gray;
  margin-bottom: 20px;
}
</style>