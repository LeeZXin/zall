<template>
  <a-tag :color="colorMap[props.status]">{{convertStatus(props.status)}}</a-tag>
</template>
<script setup>
/*
  数据库权限申请状态tag
*/
import { defineProps } from "vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const props = defineProps(["status"]);
const colorMap = {
  1: "orange",
  2: "green",
  3: "red",
  4: "purple"
};
// 转化请求状态
const convertStatus = status => {
  switch (status) {
    case 1:
      return t("mysqlReadPermApply.pendingStatus");
    case 2:
      return t("mysqlReadPermApply.agreeStatus");
    case 3:
      return t("mysqlReadPermApply.disagreeStatus");
    case 4:
      return t("mysqlReadPermApply.canceledStatus");
    default:
      return t("mysqlReadPermApply.unknownStatus");
  }
};
</script>