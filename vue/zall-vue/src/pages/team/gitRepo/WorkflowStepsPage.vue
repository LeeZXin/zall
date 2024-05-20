<template>
  <div style="padding:14px">
    <ZNaviBack
      :url="`/gitRepo/${route.params.repoId}/workflow/${route.params.workflowId}/tasks`"
      name="任务列表"
    />
    <div class="flow">
      <VueFlow
        :nodes="nodes"
        :edges="edges"
        :style="{ background: 'transparent' }"
        @nodes-initialized="layoutGraph"
      >
        <template #node-custom="nodeProps">
          <WorkflowNode v-bind="nodeProps" />
        </template>
      </VueFlow>
    </div>
    <div style="margin:10px 0;">
      <a-button type="primary" danger ghost :icon="h(PauseOutlined)">停止任务</a-button>
    </div>
    <ActionSteps />
  </div>
</template>
<script setup>
import ZNaviBack from "@/components/common/ZNaviBack";
import ActionSteps from "@/components/action/ActionSteps";
import WorkflowNode from "@/components/vueflow/WorkflowNode";
import { PauseOutlined } from "@ant-design/icons-vue";
import { ref, h } from "vue";
import { useRoute } from "vue-router";
import { VueFlow, MarkerType } from "@vue-flow/core";
import "@vue-flow/core/dist/style.css";
import "@vue-flow/core/dist/theme-default.css";
import { useLayout } from "@/utils/dagreLayout";
import jsyaml from "js-yaml";
import { getTaskDetailRequest } from "@/api/git/workflowApi";
const { layout } = useLayout();
// const actionConfig = `
// jobs:
//     job1:
//         needs: []
//         steps:
//             - name: job1 step 1
//               with:
//                 a: "1"
//                 b: "2"
//               script: |
//                 echo $a
//                 echo $b
//                 echo "job1 step 1"
//             - name: job1 step 2
//               with:
//                 c: "3"
//                 d: "4"
//               script: |
//                 echo $c
//                 echo $d
//                 echo "job1 step 2"
//         timeout: 0
//     job2:
//         needs: []
//         steps:
//             - name: job2 step 1
//               with:
//                 f: "5"
//                 g: "6"
//               script: |
//                 echo $f
//                 echo $g
//                 echo "job2 step 1"
//         timeout: 0
//     job3:
//         needs:
//             - job1
//             - job2
//         steps:
//             - name: job3e step 1
//               with:
//                 f: "5"
//                 g: "6"
//               script: |
//                 echo $f
//                 echo $g
//                 echo "job3e step 1"
//             - name: job3e step 2
//               with:
//                 f: "5"
//                 g: "6"
//               script: |
//                 echo $f
//                 echo $g
//                 echo "job3e step 2"
//         timeout: 0
//     job4:
//         needs: []
//         steps:
//             - name: job4 step 1
//               with:
//                 f: "5"
//                 g: "6"
//               script: |
//                 echo $f
//                 echo $g
//                 echo "job4 step 1"
//         timeout: 0
//     job5:
//         needs:
//             - job3
//             - job4
//         steps:
//             - name: job5 step 1
//               with:
//                 f: "5"
//                 g: "6"
//               script: |
//                 echo $f
//                 echo $g
//                 echo "job5 step 1"
//         timeout: 0
//     job6:
//         needs:
//             - job3
//         steps:
//             - name: job6 step 1
//               with:
//                 f: "5"
//                 g: "6"
//               script: |-
//                 echo $f
//                 echo $g
//                 echo "job6 step 1"
//         timeout: 0
// `;

const actionConfig2FlowElements = config => {
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
        hasLeftHandle: jobCfg.needs.length > 0,
        hasRightHandle: false
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

const layoutGraph = () => {
  nodes.value = layout(nodes.value, edges.value);
};
const nodes = ref([]);
const edges = ref([]);
const route = useRoute();
const getTaskDetail = () => {
  getTaskDetailRequest(route.params.taskId).then(res => {
    let data = actionConfig2FlowElements(res.data.yamlContent);
    nodes.value = data.nodes;
    edges.value = data.edges;
  });
};
getTaskDetail();
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
</style>