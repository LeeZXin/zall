<template>
  <div style="padding:10px">
    <ul class="info-list">
      <li>
        <div class="info-name">{{t('gitWorkflow.triggerType')}}</div>
        <div
          class="info-value"
        >{{t(taskStore.triggerType === 1 ? 'workflow.hookTriggerType': 'workflow.manualTriggerType')}}</div>
      </li>
      <li>
        <div class="info-name">{{t('gitWorkflow.operator')}}</div>
        <div class="info-value flex-center">
          <ZAvatar
            :url="taskStore.operator?.avatarUrl"
            :name="taskStore.operator?.name"
            :account="taskStore.operator?.account"
            :showName="true"
          />
        </div>
      </li>
      <li>
        <div class="info-name">{{t('gitWorkflow.createTime')}}</div>
        <div class="info-value">{{readableTimeComparingNow(taskStore.created)}}</div>
      </li>
      <li>
        <div class="info-name">{{t('gitWorkflow.branch')}}</div>
        <div class="info-value">{{taskStore.branch}}</div>
      </li>
      <li>
        <div class="info-name">{{t('gitWorkflow.pullRequest')}}</div>
        <div class="info-value">
          <PrIndexTag
            :prIndex="taskStore.prIndex"
            :repoId="route.params.repoId"
            :teamId="route.params.teamId"
            v-if="taskStore.prIndex"
          />
        </div>
      </li>
    </ul>
    <ul class="info-list">
      <li>
        <div class="info-name">{{t('gitWorkflow.taskStatus')}}</div>
        <div class="info-value">
          <RunStatus :status="taskInfo.status" />
        </div>
      </li>
      <li>
        <div class="info-name">{{t('gitWorkflow.duration')}}</div>
        <div class="info-value">{{readableDuration(taskInfo.duration)}}</div>
      </li>
      <li>
        <div class="info-name">{{t('gitWorkflow.workflowCfg')}}</div>
        <div class="info-value check-yaml-btn" @click="showYamlModal">{{t('gitWorkflow.viewCfg')}}</div>
      </li>
    </ul>
    <div class="flow">
      <VueFlow
        :nodes="nodes"
        :edges="edges"
        :style="{ background: 'transparent' }"
        @nodes-initialized="layoutGraph"
        @node-click="clickFlowNode"
      >
        <template #node-custom="nodeProps">
          <WorkflowNode v-bind="nodeProps" />
        </template>
      </VueFlow>
    </div>
    <div
      style="margin-top:10px;"
      v-if="(taskInfo.status === 'running' || taskInfo.status === 'queue') && repoStore.perm?.canTriggerWorkflow"
    >
      <a-button
        type="primary"
        danger
        ghost
        :icon="h(PauseOutlined)"
        @click="killTask"
      >{{t('gitWorkflow.killWorkflow')}}</a-button>
    </div>
    <div class="step-body">
      <div class="left">
        <ul class="job-list">
          <li
            v-for="job in jobList"
            :class="{
            'job-item-selected': job === selectedJob
        }"
            @click="selectJob(job)"
            v-bind:key="job"
          >{{job}}</li>
        </ul>
      </div>
      <div class="right" v-if="jobInfo.status.length > 0 && jobInfo.status !== 'unknown'">
        <div class="run-status">
          <RunStatus :status="jobInfo.status" />
          <span style="margin-left:8px">{{readableDurationWrap(jobInfo.duration)}}</span>
        </div>
        <ul class="step-list">
          <li v-for="(item, index) in stepsList" v-bind:key="'step_' + index">
            <div class="flex-center step-item" @click="showLogs(item, index)">
              <span class="left-down-icon">
                <down-outlined v-if="item.openLog" />
                <right-outlined v-else />
              </span>
              <check-outlined style="color:green" v-if="item.status === 'success'" />
              <close-outlined style="color:red" v-else-if="item.status === 'fail'" />
              <minus-outlined v-else />
              <span class="step-name">{{item.stepName}}</span>
              <span class="step-duration">{{readableDurationWrap(item.duration)}}</span>
            </div>
            <div class="step-log" v-if="item.openLog">
              <table>
                <colgroup>
                  <col width="44" />
                  <col />
                  <col width="44" />
                  <col />
                </colgroup>
                <tr v-for="(line, lineIndex) in item.logs" v-bind:key="'log_' + lineIndex">
                  <td>{{lineIndex + 1}}</td>
                  <td>{{line}}</td>
                </tr>
              </table>
            </div>
          </li>
        </ul>
      </div>
    </div>
    <a-modal v-model:open="yamlModalOpen" :title="t('gitWorkflow.workflowCfg')" :footer="null">
      <Codemirror
        v-model="taskStore.yamlContent"
        :style="codemirrorStyle"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import RunStatus from "@/components/git/WorkflowRunStatus";
import PrIndexTag from "@/components/git/PrIndexTag";
import WorkflowNode from "@/components/vueflow/WorkflowNode";
import { Codemirror } from "vue-codemirror";
import { yaml } from "@codemirror/lang-yaml";
import { oneDark } from "@codemirror/theme-one-dark";
import {
  PauseOutlined,
  CheckOutlined,
  CloseOutlined,
  MinusOutlined,
  RightOutlined,
  DownOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import { ref, h, onUnmounted, reactive, createVNode } from "vue";
import { useRoute } from "vue-router";
import { VueFlow, MarkerType, useVueFlow } from "@vue-flow/core";
import "@vue-flow/core/dist/style.css";
import "@vue-flow/core/dist/theme-default.css";
import { useLayout } from "@/utils/dagreLayout";
import jsyaml from "js-yaml";
import {
  getTaskStatusRequest,
  killWorkflowTaskRequest,
  getLogContentRequest,
  getTaskDetailRequest
} from "@/api/git/workflowApi";
import { readableDuration, readableTimeComparingNow } from "@/utils/time";
import { useI18n } from "vue-i18n";
import { message, Modal } from "ant-design-vue";
import { useRepoStore } from "@/pinia/repoStore";
const repoStore = useRepoStore();
const extensions = [yaml(), oneDark];
const codemirrorStyle = { height: "380px", width: "100%" };
const { findNode } = useVueFlow();
const { t } = useI18n();
const taskStore = ref({});
const yamlModalOpen = ref(false);
// 任务状态
const taskInfo = reactive({
  status: "unknown",
  duration: 0
});
// 单个job状态
const jobInfo = reactive({
  status: "unknown",
  duration: 0
});
const { layout } = useLayout();
const stepsList = ref([]);
const readableDurationWrap = duration => {
  if (duration > 0) {
    return readableDuration(duration);
  }
  return "";
};
const actionConfig2FlowElements = jobMap => {
  let config = taskStore.value.yamlContent;
  let action = jsyaml.load(config);
  let nodes = [];
  let edges = [];
  let nodesMap = {};
  let hasRightHandleNodesMap = {};
  for (let jobName in action.jobs) {
    let jobCfg = action.jobs[jobName];
    nodesMap[jobName] = {
      id: jobName,
      label: jobName,
      type: "custom",
      position: { x: 0, y: 0 },
      data: {
        hasLeftHandle: jobCfg.needs?.length > 0,
        hasRightHandle: false,
        duration: readableDurationWrap(
          jobMap[jobName] ? jobMap[jobName].duration : 0
        ),
        status: jobMap[jobName] ? jobMap[jobName].status : ""
      }
    };
    if (jobCfg.needs) {
      jobCfg.needs.forEach(item => {
        edges.push({
          id: "edge_" + jobName + "_" + item,
          source: item,
          target: jobName,
          markerEnd: MarkerType.ArrowClosed
        });
        hasRightHandleNodesMap[item] = true;
      });
    }
  }
  for (let jobName in hasRightHandleNodesMap) {
    let node = nodesMap[jobName];
    if (node) {
      node.data.hasRightHandle = true;
    }
  }
  for (let jobName in nodesMap) {
    nodes.push(nodesMap[jobName]);
  }
  return { nodes, edges };
};
const getJobAndSteps = jobName => {
  let job = allJobs.value[jobName];
  if (job) {
    jobInfo.status = job.status;
    jobInfo.duration = job.duration;
    let steps = job.steps ? job.steps : [];
    steps = steps.map(item => {
      return {
        ...item,
        logs: [],
        openLog: false
      };
    });
    stepsList.value = steps;
  }
};
const layoutGraph = () => {
  nodes.value = layout(nodes.value, edges.value);
};
const jobList = ref([]);
const selectedJob = ref("");
const selectJob = job => {
  if (selectedJob.value === job) {
    return;
  }
  nodes.value.forEach(item => {
    const node = findNode(item.id);
    if (item.id === job) {
      if (node) {
        node.selected = true;
      }
    } else {
      node.selected = false;
    }
  });
  selectedJob.value = job;
  getJobAndSteps(job);
};
const nodes = ref([]);
const edges = ref([]);
const route = useRoute();
const statusInterval = ref(null);
const allJobs = ref({});
const getTaskStatus = () => {
  getTaskStatusRequest(route.params.taskId).then(res => {
    taskInfo.status = res.data.status;
    taskInfo.duration = res.data.duration;
    let jobs = res.data.jobStatus ? res.data.jobStatus : [];
    let jobMap = {};
    jobs.forEach(item => {
      jobMap[item.jobName] = item;
    });
    allJobs.value = jobMap;
    if (nodes.value.length === 0) {
      let data = actionConfig2FlowElements(jobMap);
      nodes.value = data.nodes;
      edges.value = data.edges;
      jobList.value = data.nodes.map(item => {
        return item.id;
      });
    } else {
      nodes.value.forEach(item => {
        let job = jobMap[item.id];
        if (job) {
          if (job.status !== item.data.status) {
            item.data.status = job.status;
          }
          item.data.duration = readableDurationWrap(job.duration);
        }
      });
    }
    if (res.data.status === "queue" || res.data.status === "running") {
      if (!statusInterval.value) {
        statusInterval.value = setInterval(() => {
          getTaskStatus();
        }, 5000);
      }
    } else if (statusInterval.value !== null) {
      clearInterval(statusInterval.value);
      statusInterval.value = null;
    }
  });
};
// 获取任务详情
const getTaskDetail = () => {
  getTaskDetailRequest(route.params.taskId).then(res => {
    taskStore.value = res.data;
    getTaskStatus();
  });
};
// 点击节点
const clickFlowNode = e => {
  selectJob(e.node.id);
};
// 展示配置
const showYamlModal = () => {
  yamlModalOpen.value = true;
};
// 停止任务
const killTask = () => {
  Modal.confirm({
    title: `${t("gitWorkflow.confirmKill")}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      killWorkflowTaskRequest(taskStore.value.id).then(() => {
        message.success(t("operationSuccess"));
      });
    },
    onCancel() {}
  });
};
// 查看日志
const showLogs = (item, index) => {
  if (item.openLog) {
    item.openLog = false;
  } else if (item.logs && item.logs.length > 0) {
    item.openLog = true;
  } else {
    getLogContentRequest(taskStore.value.id, selectedJob.value, index).then(
      res => {
        item.logs = res.data;
        item.openLog = true;
      }
    );
  }
};
getTaskDetail();
onUnmounted(() => {
  if (statusInterval.value) {
    clearInterval(statusInterval.value);
    statusInterval.value = null;
  }
});
</script>
<style scoped>
.header {
  font-size: 14px;
  cursor: pointer;
  margin-bottom: 10px;
}
.header:hover {
  color: #1677ff;
}
.task-name {
  font-size: 16px;
  line-height: 32px;
  margin-right: 8px;
  font-weight: bold;
}
.flow {
  width: 100%;
  height: 400px;
  border-radius: 4px;
  background-color: #f0f0f0;
}
.step-body {
  margin-top: 10px;
  display: flex;
  background-color: #001529;
  border-radius: 4px;
  color: white;
}
.step-body > .left {
  width: 20%;
  padding: 10px;
  height: 500px;
  overflow: scroll;
}
.step-body > .right {
  width: 80%;
  padding: 10px;
  height: 500px;
  overflow: scroll;
}
.job-list > li {
  font-size: 14px;
  white-space: pre-wrap;
  padding: 12px 20px;
  cursor: pointer;
  width: 100%;
  border-radius: 4px;
  position: relative;
}
.job-list > li:hover {
  background-color: #b1bac41f;
}
.job-item-selected {
  background-color: #b1bac41f;
}
.job-item-selected::before {
  content: " ";
  width: 4px;
  height: 38px;
  background-color: #1677ff;
  border-radius: 2px;
  position: absolute;
  top: 0px;
  left: 0px;
}
.run-status {
  margin-bottom: 10px;
  font-size: 12px;
}
.step-list > li {
  font-size: 14px;
}
.step-name {
  margin-left: 20px;
  font-weight: bold;
}
.step-duration {
  margin-left: 12px;
  color: gray;
}
.left-down-icon {
  margin-right: 20px;
  font-size: 12px;
}
.step-log {
  background-color: #b1bac41f;
  border-bottom-right-radius: 4px;
  border-bottom-left-radius: 4px;
  color: white;
}
.step-log > table > tr > td {
  line-height: 22px;
  font-size: 14px;
  word-break: break-all;
}
.step-log > table {
  width: 100%;
}
.step-log > table > tr > td:first-child {
  text-align: center;
  align-content: baseline;
}
.step-item {
  line-height: 38px;
  padding: 0 10px;
  cursor: pointer;
}
.check-yaml-btn:hover {
  cursor: pointer;
  color: #1677ff;
}
</style>