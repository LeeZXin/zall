<template>
  <div>
    <Handle type="target" :position="Position.Left" v-if="props.data.hasLeftHandle" />
    <div class="node" :class="{
        'node-selected': props.selected
    }">
      <div class="label">
        <span style="font-weight:bold;">{{ props.label }}</span>
        <span class="label-duration">{{props.data.duration}}</span>
      </div>
      <div class="status" v-if="props.data.status === 'running'">
        <rocket-outlined />
        <span style="padding-left:6px">运行中</span>
      </div>
      <div class="status" v-else-if="props.data.status === 'success'">
        <check-circle-filled style="color:green" />
        <span style="padding-left:6px">成功</span>
      </div>
      <div class="status" v-else-if="props.data.status === 'fail'">
        <close-circle-filled style="color:red" />
        <span style="padding-left:6px">失败</span>
      </div>
      <div class="status" v-else-if="props.data.status === 'timeout'">
        <close-circle-filled style="color:red" />
        <span style="padding-left:6px">超时</span>
      </div>
      <div class="status" v-else>
        <play-circle-outlined />
        <span style="padding-left:6px">未执行</span>
      </div>
    </div>
    <Handle type="source" :position="Position.Right" v-if="props.data.hasRightHandle" />
  </div>
</template>
<script setup>
import { Position, Handle } from "@vue-flow/core";
import { defineProps } from "vue";
import {
  RocketOutlined,
  PlayCircleOutlined,
  CheckCircleFilled,
  CloseCircleFilled
} from "@ant-design/icons-vue";
// props were passed from the slot using `v-bind="customNodeProps"`
const props = defineProps(["label", "selected", "data"]);
</script>
<style scoped>
.node {
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  padding: 10px;
  width: 240px;
  background-color: white;
}
.node-selected {
  border-color: black;
}
.node > .label {
  font-size: 14px;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.node > .status {
  font-size: 14px;
}
.label-duration {
  float: right;
  color: gray;
  font-size: 14px;
}
</style>