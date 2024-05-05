<template>
  <a-tag :color="colorMap[props.status]">{{convertPrStatus(props.status)}}</a-tag>
</template>
<script setup>
import { defineProps } from "vue";
const props = defineProps(["status"]);
const colorMap = {
  1: "green",
  2: "red",
  3: "purple"
};
// 转化请求状态
const convertPrStatus = status => {
  switch (status) {
    case 1:
      return "已打开";
    case 2:
      return "已关闭";
    case 3:
      return "已合并";
    default:
      return "未知";
  }
};
</script>