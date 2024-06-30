<template>
  <a-tag
    :style="props.style"
    :color="colorMap[props.status]"
  >{{t("workflow.status." + props.status)}}</a-tag>
</template>
<script setup>
import { defineProps } from "vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const props = defineProps(["status", "style"]);
const colorMap = {
  queue: "orange",
  running: "purple",
  success: "green",
  fail: "red",
  cancel: "red",
  timeout: "red",
  unExecuted: "red",
  unknown: "red"
};
</script>