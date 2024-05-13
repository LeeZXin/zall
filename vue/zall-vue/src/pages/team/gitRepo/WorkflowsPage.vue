<template>
  <div style="padding:14px">
    <div style="margin-bottom:10px">
      <a-button type="primary" @click="gotoCreatePage">创建工作流</a-button>
    </div>
    <ul class="workflow-list" v-if="workflowList.length > 0">
      <li v-for="item in workflowList" v-bind:key="item.id">
        <div class="workflow-header">
          <a-tooltip placement="top">
            <template #title>{{item.name}}</template>
            <div class="name">{{item.name}}</div>
          </a-tooltip>
          <span>
            <a-tooltip placement="top" v-if="!item.lastTask || item.lastTask.taskStatus !== 1">
              <template #title>手动执行</template>
              <span class="op">
                <PlayCircleFilled />
              </span>
            </a-tooltip>
            <a-tooltip placement="top" v-if="item.lastTask?.taskStatus === 1">
              <template #title>暂停执行</template>
              <span class="op">
                <PauseOutlined />
              </span>
            </a-tooltip>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <li @click="deleteWorkflow(item)">
                    <DeleteOutlined />
                    <span style="margin-left:4px">删除工作流</span>
                  </li>
                  <li @click="deleteWorkflow(item)">
                    <EditOutlined />
                    <span style="margin-left:4px">编辑工作流</span>
                  </li>
                  <li @click="gotoDetailPage(item.id)">
                    <EyeOutlined />
                    <span style="margin-left:4px">查看任务</span>
                  </li>
                </ul>
              </template>
              <div class="op">
                <EllipsisOutlined />
              </div>
            </a-popover>
          </span>
        </div>
        <div class="workflow-status">
          <div class="no-wrap" v-if="item.lastTask?.taskStatus === 2">
            <CheckCircleFilled style="color:green" />
            <span style="margin-left: 6px">执行成功</span>
          </div>
          <div class="no-wrap" v-else-if="item.lastTask?.taskStatus === 3">
            <CloseCircleFilled style="color:darkred" />
            <span style="margin-left: 6px">执行失败</span>
          </div>
          <div class="no-wrap" v-else-if="item.lastTask?.taskStatus === 1">
            <LoadingOutlined />
            <span style="margin-left: 6px">执行中</span>
          </div>
          <div class="no-wrap" v-else>
            <ClockCircleOutlined />
            <span style="margin-left: 6px">无执行任务</span>
          </div>
          <div
            class="no-wrap"
            style="margin-top:10px;"
            v-if="item.lastTask"
          >{{item.lastTask.operator}}推送{{item.lastTask.branch}}触发</div>
          <div style="margin-top:6px;height:14px" v-else></div>
        </div>
      </li>
    </ul>
    <ZNoData v-else>
      <template #desc>
        <div class="no-data">
          <span>没有配置工作流, 可点击上方“创建工作流”创建新的工作流</span>
        </div>
      </template>
    </ZNoData>
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
  CheckCircleFilled,
  CloseCircleFilled,
  LoadingOutlined,
  ClockCircleOutlined,
  PauseOutlined
} from "@ant-design/icons-vue";
import ZNoData from "@/components/common/ZNoData";
import { ref, createVNode } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  listWorkflowRequest,
  DeleteWorkflowRequest
} from "@/api/git/workflowApi";
import { Modal, message } from "ant-design-vue";
const router = useRouter();
const route = useRoute();
const workflowList = ref([]);
const gotoCreatePage = () => {
  router.push(`/gitRepo/${route.params.repoId}/workflow/create`);
};
const listWorkflow = () => {
  listWorkflowRequest(route.params.repoId).then(res => {
    workflowList.value = res.data;
  });
};
const gotoDetailPage = id => {
  router.push(`/gitRepo/${route.params.repoId}/workflow/${id}/tasks`);
};
const deleteWorkflow = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      DeleteWorkflowRequest(item.id).then(() => {
        message.success("删除成功");
        listWorkflow();
      });
    },
    onCancel() {}
  });
};
listWorkflow();
</script>
<style scoped>
.op {
  display: inline-block;
  width: 24px;
  height: 24px;
  line-height: 24px;
  font-size: 16px;
  text-align: center;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
}
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
  width: calc(25% - 8px);
  padding: 14px 20px;
}
.workflow-list > li:not(:nth-child(4n + 1)) {
  margin-left: 8px;
}
.workflow-list > li:nth-child(n + 5) {
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