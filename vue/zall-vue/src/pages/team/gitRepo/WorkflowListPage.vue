<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px" v-if="repoStore.perm?.canManageWorkflow">
      <a-button
        type="primary"
        @click="gotoVarsPage"
        :icon="h(KeyOutlined)"
        style="margin-right:8px"
      >{{t('gitWorkflow.manageVariables')}}</a-button>
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('gitWorkflow.createWorkflow')}}</a-button>
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
                v-if="!item.lastTask || (item.lastTask.taskStatus !== 'running' && item.lastTask.taskStatus !== 'queue')"
              >
                <template #title>{{t('gitWorkflow.triggerWorkflow')}}</template>
                <span class="op-icon" @click="showBranchModal(item)">
                  <PlayCircleFilled />
                </span>
              </a-tooltip>
              <a-tooltip placement="top" v-else>
                <template #title>{{t('gitWorkflow.killWorkflow')}}</template>
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
                      <span style="margin-left:4px">{{t('gitWorkflow.deleteWorkflow')}}</span>
                    </li>
                    <li @click="gotoDetailPage(item.id)">
                      <EditOutlined />
                      <span style="margin-left:4px">{{t('gitWorkflow.updateWorkflow')}}</span>
                    </li>
                  </template>
                  <li @click="gotoTasksPage(item)">
                    <EyeOutlined />
                    <span style="margin-left:4px">{{t('gitWorkflow.viewTasks')}}</span>
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
          <div v-if="item.lastTask" class="flex-center no-wrap" style="margin-top:10px;">
            <ZAvatar
              :url="item.lastTask?.operator.avatarUrl"
              :name="item.lastTask?.operator.name"
              :account="item.lastTask?.operator.account"
              :showName="true"
            />
            <span style="padding-left:4px">{{t('gitWorkflow.push')}}</span>
            <span style="padding-left:4px">{{item.lastTask.branch}}</span>
            <span style="padding-left:4px">{{readableTimeComparingNow(item.lastTask.created)}}</span>
          </div>
          <div v-else style="height:20px;margin-top:10px"></div>
        </div>
      </li>
    </ul>
    <ZNoData v-else />
    <a-modal v-model:open="branchModal.open" :title="branchModal.title" @ok="triggerWorkflow">
      <div class="flex-center">
        <div style="line-height:32px;font-size:14px;width:80px">{{t('gitWorkflow.branch')}}:</div>
        <a-select
          style="width:100%"
          v-model:value="branchModal.branch"
          :options="branchModal.branchList"
        />
      </div>
    </a-modal>
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
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
import { ref, createVNode, onUnmounted, h, reactive } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  listWorkflowRequest,
  deleteWorkflowRequest,
  triggerWorkflowRequest,
  killWorkflowTaskRequest
} from "@/api/git/workflowApi";
import { allBranchesRequest } from "@/api/git/repoApi";
import { Modal, message } from "ant-design-vue";
import { workflowBranchRegexp } from "@/utils/regexp";
import { readableTimeComparingNow } from "@/utils/time";
import { useWorkflowStore } from "@/pinia/workflowStore";
import { useRepoStore } from "@/pinia/repoStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const repoStore = useRepoStore();
const workflowStore = useWorkflowStore();
const router = useRouter();
const route = useRoute();
// 工作流列表
const workflowList = ref([]);
// 手动触发modal
const branchModal = reactive({
  open: false,
  title: "",
  branch: "",
  wfId: 0,
  branchList: []
});
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/create`
  );
};
// 获取工作流列表
const listWorkflow = () => {
  listWorkflowRequest(route.params.repoId).then(res => {
    workflowList.value = res.data;
  });
};
// 跳转详情页面
const gotoDetailPage = id => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/${id}/update`
  );
};
// 跳转工作流变量页面
const gotoVarsPage = () => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/vars`
  );
};
// 跳转任务列表页面
const gotoTasksPage = item => {
  workflowStore.id = item.id;
  workflowStore.name = item.name;
  workflowStore.desc = item.desc;
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/${item.id}/tasks`
  );
};
// 展示分支modal
const showBranchModal = item => {
  if (branchModal.branchList.length === 0) {
    allBranchesRequest(route.params.repoId).then(res => {
      branchModal.branchList = res.data.map(item => {
        return {
          value: item,
          label: item
        };
      });
      if (res.data.length > 0) {
        branchModal.branch = res.data[0];
      }
    });
  } else {
    branchModal.branch = branchModal.branchList[0].value;
  }
  branchModal.open = true;
  branchModal.title = item.name;
  branchModal.wfId = item.id;
};
// 删除工作流
const deleteWorkflow = item => {
  Modal.confirm({
    title: `${t("gitWorkflow.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteWorkflowRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listWorkflow();
      });
    },
    onCancel() {}
  });
};
// 手动触发工作流
const triggerWorkflow = () => {
  if (!workflowBranchRegexp.test(branchModal.branch)) {
    message.warn("分支格式错误");
    return;
  }
  triggerWorkflowRequest(branchModal.wfId, branchModal.branch).then(() => {
    message.success(t("operationSuccess"));
    listWorkflow();
    branchModal.open = false;
    return;
  });
};
listWorkflow();
// 每5秒获取工作流状态
const listInterval = setInterval(() => {
  listWorkflow();
}, 5000);
// unmounted取消定时器
onUnmounted(() => {
  clearInterval(listInterval);
});
// 停止任务
const killTask = item => {
  Modal.confirm({
    title: `${t("gitWorkflow.confirmKill")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      killWorkflowTaskRequest(item.lastTask.id).then(() => {
        message.success(t("operationSuccess"));
        listWorkflow();
      });
    },
    onCancel() {}
  });
};
</script>
<style scoped>
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
</style>