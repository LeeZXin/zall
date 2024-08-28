<template>
  <span v-if="props.status === 1">
    <LoadingOutlined />
    <span style="margin-left:8px">等待执行</span>
  </span>
  <span v-else-if="props.status === 2">
    <LoadingOutlined />
    <span style="margin-left:8px">执行中</span>
  </span>
  <span v-else-if="props.status === 3">
    <CheckCircleFilled style="color:#52c41a" />
    <span style="margin-left:8px">成功</span>
  </span>
  <span v-else-if="props.status === 4">
    <CloseCircleFilled style="color:#ff4d4f" />
    <span style="margin-left:8px">失败</span>
  </span>
</template>
<script setup>
import { defineProps } from "vue";
import {
  CheckCircleFilled,
  CloseCircleFilled,
  LoadingOutlined
} from "@ant-design/icons-vue";
/*
  发布计划阶段状态值
*/
const props = defineProps(["status"]);
</script>

<style scoped>
</style>