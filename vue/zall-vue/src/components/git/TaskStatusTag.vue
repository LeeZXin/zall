<template>
  <a-tag :color="colorMap[props.status]">{{convertStatus(props.status)}}</a-tag>
</template>
<script setup>
import { defineProps } from "vue";
const props = defineProps(["status"]);
const colorMap = {
  1: "purple",
  2: "green",
  3: "red",
  4: "red"
};
// 转化请求状态
const convertStatus = status => {
  switch (status) {
    case 1:
      return "执行中";
    case 2:
      return "成功";
    case 3:
      return "失败";
    case 4:
        return "中止";
    default:
      return "未知";
  }
};
</script>