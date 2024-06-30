<template>
  <a-tag :color="colorMap[props.status]">{{convertStatus(props.status)}}</a-tag>
</template>
<script setup>
import { defineProps } from "vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const props = defineProps(["status"]);
const colorMap = {
  1: "orange",
  2: "green",
  3: "purple",
  4: "red"
};
// 转化请求状态
const convertStatus = status => {
  switch (status) {
    case 1:
      return t("deployPlan.pendingStatus");
    case 2:
      return t("deployPlan.runningStatus");
    case 3:
      return t("deployPlan.successfulStatus");
    case 4:
      return t("deployPlan.closedStatus");
    default:
      return t("deployPlan.unknownStatus");
  }
};
</script>