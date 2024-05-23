<template>
  <div class="status" :style="props.style">
    <loading-outlined v-if="props.status === 'running'" />
    <check-circle-filled style="color:green" v-else-if="props.status === 'success'" />
    <close-circle-filled
      style="color:red"
      v-else-if="['fail', 'timeout', 'cancel'].includes(props.status)"
    />
    <play-circle-filled style="color:gray" v-else />
    <span style="padding-left:6px">{{props.hideText?'':t(getStatus(props.status))}}</span>
  </div>
</template>
<script setup>
import { useI18n } from "vue-i18n";
import { defineProps } from "vue";
import {
  LoadingOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  PlayCircleFilled
} from "@ant-design/icons-vue";
const { t } = useI18n();
const props = defineProps(["style", "status", "hideText"]);
const statusMap = {
  success: "workflow.status.success",
  fail: "workflow.status.fail",
  timeout: "workflow.status.timeout",
  cancel: "workflow.status.cancel",
  queue: "workflow.status.queue",
  running: "workflow.status.running",
  unknown: "workflow.status.unknown",
  unExecuted: "workflow.status.unExecuted"
};
const getStatus = status => {
  let ret = statusMap[status];
  if (!ret) {
    return statusMap["unknown"];
  }
  return ret;
};
</script>
<style scoped>
.status {
  display: inline;
  font-size: 14px;
}
</style>