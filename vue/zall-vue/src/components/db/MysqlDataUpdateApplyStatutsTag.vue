<template>
  <a-tag :color="colorMap[props.status]">{{convertStatus(props.status)}}</a-tag>
</template>
<script setup>
import { defineProps } from "vue";
import { useI18n } from "vue-i18n";
/*
  数据修改单状态tag
*/
const { t } = useI18n();
const props = defineProps(["status"]);
const colorMap = {
  1: "orange",
  2: "green",
  3: "red",
  4: "purple",
  5: "blue",
  6: "green"
};
// 转化请求状态
const convertStatus = status => {
  switch (status) {
    case 1:
      return t("mysqlDataUpdateApply.pendingStatus");
    case 2:
      return t("mysqlDataUpdateApply.agreeStatus");
    case 3:
      return t("mysqlDataUpdateApply.disagreeStatus");
    case 4:
      return t("mysqlDataUpdateApply.canceledStatus");
    case 5:
      return t("mysqlDataUpdateApply.askToExecuteStatus");
    case 6:
      return t("mysqlDataUpdateApply.executedStatus");
    default:
      return t("mysqlDataUpdateApply.unknownStatus");
  }
};
</script>