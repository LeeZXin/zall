<template>
  <span v-if="props.status === 1">
    <LoadingOutlined />
    <span style="margin-left:8px">{{t('deployPlan.stagePendingStatus')}}</span>
  </span>
  <span v-else-if="props.status === 2">
    <LoadingOutlined />
    <span style="margin-left:8px">{{t('deployPlan.stageRunningStatus')}}</span>
  </span>
  <span v-else-if="props.status === 3">
    <CheckCircleFilled style="color:#52c41a" />
    <span style="margin-left:8px">{{t('deployPlan.stageSuccessfulStatus')}}</span>
  </span>
  <span v-else-if="props.status === 4">
    <CloseCircleFilled style="color:#ff4d4f" />
    <span style="margin-left:8px">{{t('deployPlan.stageFailedStatus')}}</span>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const props = defineProps(["status"]);
</script>

<style scoped>
</style>