<template>
  <div style="padding:14px">
    <div style="margin-bottom:10px">
      <span class="header" @click="backToTaskList">
        <arrow-left-outlined />
        <span style="margin-left:8px">任务列表</span>
      </span>
    </div>
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
    <div class="table-top">
      <div class="task-name">
        <span>xxxxx</span>
        <span>-</span>
        <span>任务详情</span>
      </div>
      <a-button type="primary" danger size="small">停止任务</a-button>
    </div>
    <ActionSteps />
  </div>
</template>
<script setup>
import ActionSteps from "@/components/action/ActionSteps";
import WorkflowNode from "@/components/vueflow/WorkflowNode";
import { ArrowLeftOutlined } from "@ant-design/icons-vue";
import { ref, nextTick } from "vue";
import { useRouter } from "vue-router";
import { VueFlow, MarkerType, useVueFlow } from "@vue-flow/core";
import "@vue-flow/core/dist/style.css";
import "@vue-flow/core/dist/theme-default.css";
import { useLayout } from "@/utils/dagreLayout";
import jsyaml from "js-yaml";
const { layout } = useLayout();
const { fitView } = useVueFlow();
const actionConfig = `
jobs:
    job1:
        needs: []
        steps:
            - name: job1 step 1
              with:
                a: "1"
                b: "2"
              script: |
                echo $a
                echo $b
                echo "job1 step 1"
            - name: job1 step 2
              with:
                c: "3"
                d: "4"
              script: |
                echo $c
                echo $d
                echo "job1 step 2"
        timeout: 0
    job2:
        needs: []
        steps:
            - name: job2 step 1
              with:
                f: "5"
                g: "6"
              script: |
                echo $f
                echo $g
                echo "job2 step 1"
        timeout: 0
    job3:
        needs:
            - job1
            - job2
        steps:
            - name: job3e step 1
              with:
                f: "5"
                g: "6"
              script: |
                echo $f
                echo $g
                echo "job3e step 1"
            - name: job3e step 2
              with:
                f: "5"
                g: "6"
              script: |
                echo $f
                echo $g
                echo "job3e step 2"
        timeout: 0
    job4:
        needs: []
        steps:
            - name: job4 step 1
              with:
                f: "5"
                g: "6"
              script: |
                echo $f
                echo $g
                echo "job4 step 1"
        timeout: 0
    job5:
        needs:
            - job3
            - job4
        steps:
            - name: job5 step 1
              with:
                f: "5"
                g: "6"
              script: |
                echo $f
                echo $g
                echo "job5 step 1"
        timeout: 0
    job6:
        needs:
            - job3
        steps:
            - name: job6 step 1
              with:
                f: "5"
                g: "6"
              script: |-
                echo $f
                echo $g
                echo "job6 step 1"
        timeout: 0
`;

const actionConfig2FlowElements = config => {
  let action = jsyaml.load(config);
  let nodes = [];
  let edges = [];
  for (let jobName in action.jobs) {
    nodes.push({
      id: jobName,
      label: jobName,
      type: "custom",
      position: { x: 0, y: 0 }
    });
    let jobCfg = action.jobs[jobName];
    if (jobCfg.needs) {
      jobCfg.needs.forEach(item => {
        edges.push({
          id: "edge_" + jobName + "_" + item,
          source: jobName,
          target: item,
          markerEnd: MarkerType.ArrowClosed
        });
      });
    }
  }
  return { nodes, edges };
};

const layoutGraph = () => {
  nodes.value = layout(nodes.value, edges.value);
  nextTick(() => {
    fitView();
  });
};

const data = actionConfig2FlowElements(actionConfig)
const nodes = ref(data.nodes);
const edges = ref(data.edges);
const router = useRouter();

const backToTaskList = () => {
  router.push("/team/action/list");
};
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
.table-top {
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  width: 100%;
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