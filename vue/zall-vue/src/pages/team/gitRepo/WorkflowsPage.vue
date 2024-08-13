<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px" v-if="repoStore.perm?.canManageWorkflow">
      <a-button
        type="primary"
        @click="gotoVarsPage"
        :icon="h(KeyOutlined)"
        style="margin-right:8px"
      >管理变量</a-button>
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">创建工作流</a-button>
    </div>
    <ul class="workflow-list" v-if="workflowList.length > 0">
      <li v-for="item in workflowList" v-bind:key="item.id">
        <div class="workflow-header">
          <a-tooltip placement="top">
            <template #title>{{item.name}}</template>
            <div class="name">{{item.name}}</div>
          </a-tooltip>
          <span>
            <template v-if="repoStore.perm?.canTriggerWorkflow">
              <a-tooltip
                placement="top"
                v-if="!item.lastTask || (item.lastTask.taskStatus !== 1 && item.lastTask.taskStatus !== 0)"
              >
                <template #title>手动执行</template>
                <span class="op-icon" @click="showBranchModal(item)">
                  <PlayCircleFilled />
                </span>
              </a-tooltip>
              <a-tooltip placement="top" v-else>
                <template #title>停止执行</template>
                <span class="op-icon" @click="killTask(item)">
                  <PauseOutlined />
                </span>
              </a-tooltip>
            </template>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <template v-if="repoStore.perm?.canManageWorkflow">
                    <li @click="deleteWorkflow(item)">
                      <DeleteOutlined />
                      <span style="margin-left:4px">删除工作流</span>
                    </li>
                    <li @click="gotoDetailPage(item.id)">
                      <EditOutlined />
                      <span style="margin-left:4px">编辑工作流</span>
                    </li>
                  </template>
                  <li @click="gotoTasksPage(item)">
                    <EyeOutlined />
                    <span style="margin-left:4px">查看任务</span>
                  </li>
                </ul>
              </template>
              <div class="op-icon">
                <EllipsisOutlined />
              </div>
            </a-popover>
          </span>
        </div>
        <div class="workflow-status">
          <WorkflowTaskStatusIconText :status="item.lastTask?.taskStatus" />
          <div
            class="no-wrap"
            style="margin-top:10px;"
            v-if="item.lastTask"
          >{{item.lastTask.operator}}推送{{item.lastTask.branch}}触发 {{readableTimeComparingNow(item.lastTask.created)}}</div>
          <div style="margin-top:6px;height:14px" v-else></div>
        </div>
      </li>
    </ul>
    <ZNoData v-else />
    <a-modal
      v-model:open="branchModalOpen"
      :title="branchModalTitle"
      @ok="triggerWorkflow"
      okText="立即触发"
      cancelText="取消"
    >
      <div class="flex-center">
        <div style="line-height:32px;font-size:14px;width:80px">分支:</div>
        <a-input style="width:100%" v-model:value="triggerBranch" />
      </div>
    </a-modal>
  </div>
</template>
<script setup>
import {
  DeleteOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  EditOutlined,
  EyeOutlined,
  PlayCircleFilled,
  PauseOutlined,
  PlusOutlined,
  KeyOutlined
} from "@ant-design/icons-vue";
import WorkflowTaskStatusIconText from "@/components/git/WorkflowTaskStatusIconText";
import ZNoData from "@/components/common/ZNoData";
import { ref, createVNode, onUnmounted, h } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  listWorkflowRequest,
  deleteWorkflowRequest,
  triggerWorkflowRequest,
  killWorkflowTaskRequest
} from "@/api/git/workflowApi";
import { Modal, message } from "ant-design-vue";
import { workflowBranchRegexp } from "@/utils/regexp";
import { readableTimeComparingNow } from "@/utils/time";
import { useWorkflowStore } from "@/pinia/workflowStore";
import { useRepoStore } from "@/pinia/repoStore";
const repoStore = useRepoStore();
const workflowStore = useWorkflowStore();
const router = useRouter();
const route = useRoute();
const workflowList = ref([]);
const branchModalOpen = ref(false);
const branchModalTitle = ref("");
const triggerBranch = ref("");
const triggerWfId = ref(0);
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/create`
  );
};
const listWorkflow = () => {
  listWorkflowRequest(route.params.repoId).then(res => {
    workflowList.value = res.data;
  });
};
const gotoDetailPage = id => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/${id}/update`
  );
};
const gotoVarsPage = () => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/vars`
  );
};
const gotoTasksPage = item => {
  workflowStore.id = item.id;
  workflowStore.name = item.name;
  workflowStore.desc = item.desc;
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/${item.id}/tasks`
  );
};
const showBranchModal = item => {
  branchModalOpen.value = true;
  triggerBranch.value = "";
  branchModalTitle.value = item.name;
  triggerWfId.value = item.id;
};
const deleteWorkflow = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteWorkflowRequest(item.id).then(() => {
        message.success("删除成功");
        listWorkflow();
      });
    },
    onCancel() {}
  });
};
const triggerWorkflow = () => {
  if (!workflowBranchRegexp.test(triggerBranch.value)) {
    message.warn("分支格式错误");
    return;
  }
  triggerWorkflowRequest(triggerWfId.value, triggerBranch.value).then(() => {
    message.success("操作成功");
    listWorkflow();
    branchModalOpen.value = false;
    return;
  });
};
listWorkflow();
const listInterval = setInterval(() => {
  listWorkflow();
}, 5000);
onUnmounted(() => {
  clearInterval(listInterval);
});
const killTask = item => {
  Modal.confirm({
    title: `你确定要停止${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      killWorkflowTaskRequest(item.lastTask.id).then(() => {
        message.success("停止成功");
        listWorkflow();
      });
    },
    onCancel() {}
  });
};
</script>
<style scoped>
.no-data {
  font-size: 14px;
  text-align: center;
}
.workflow-list {
  display: flex;
  flex-wrap: wrap;
}
.workflow-list > li {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  width: calc(50% - 6px);
  padding: 14px 20px;
}
.workflow-list > li:not(:nth-child(2n + 1)) {
  margin-left: 12px;
}
.workflow-list > li:nth-child(n + 3) {
  margin-top: 10px;
}
.workflow-header {
  width: 100%;
  margin-bottom: 14px;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.workflow-header > .name {
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-all;
  max-width: 50%;
}
.workflow-status {
  font-size: 14px;
  color: gray;
  word-break: break-all;
}
.op:hover {
  background-color: #f0f0f0;
}
</style>