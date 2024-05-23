<template>
  <a-tag :style="props.style" :color="colorMap[props.status]">{{convertStatus(props.status)}}</a-tag>
</template>
<script setup>
import { defineProps } from "vue";
const props = defineProps(["status", "style"]);
const colorMap = {
  0: "orange",
  1: "purple",
  2: "green",
  3: "red",
  4: "red",
  5: "red"
};
// 转化请求状态
const convertStatus = status => {
  switch (status) {
    case 0:
      return "排队中";
    case 1:
      return "执行中";
    case 2:
      return "成功";
    case 3:
      return "失败";
    case 4:
      return "中止";
    case 5:
      return "超时";
    default:
      return "未知";
  }
};
</script>