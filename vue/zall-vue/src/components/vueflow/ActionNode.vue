<template>
  <div>
    <Handle type="target" :position="Position.Left" />
    <div class="node" :class="{
        'node-selected': props.selected
    }">
      <div class="label">
        <span style="font-weight:bold;">{{ props.label }}</span>
        <span class="label-duration">1m12s</span>
      </div>
      <div class="status" v-if="props.data === 1">
        <rocket-outlined />
        <span style="padding-left:6px">运行中</span>
      </div>
      <div class="status" v-if="props.data === 2">
        <play-circle-outlined />
        <span style="padding-left:6px">未执行</span>
      </div>
      <div class="status" v-if="props.data === 3">
        <check-circle-filled style="color:green" />
        <span style="padding-left:6px">成功</span>
      </div>
      <div class="status">
        <close-circle-filled style="color:red" />
        <span style="padding-left:6px">失败</span>
      </div>
    </div>
    <Handle type="source" :position="Position.Right" />
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
const props = defineProps(["label", "selected"]);
</script>
<style scoped>
.node {
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  padding: 6px;
  width: 160px;
  background-color: white;
}
.node-selected {
  border-color: black;
}
.node > .label {
  font-size: 12px;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.node > .status {
  font-size: 10px;
}
.label-duration {
  float: right;
  color: gray;
  font-size: 10px;
}
</style>