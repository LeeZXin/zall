<template>
  <a-tag :color="colorMap[props.status]">{{t(convertStatus(props.status))}}</a-tag>
</template>
<script setup>
import { defineProps } from "vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const props = defineProps(["status"]);
const colorMap = {
  1: "green",
  2: "red",
  3: "purple"
};
// 转化请求状态
const convertStatus = status => {
  switch (status) {
    case 1:
      return "pullRequest.openStatus";
    case 2:
      return "pullRequest.closedStatus";
    case 3:
      return "pullRequest.mergedStatus";
  }
  return "";
};
</script>