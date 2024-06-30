<template>
  <div class="no-wrap" v-if="props.status">
    <LoadingOutlined v-if="props.status === 'queue'" />
    <LoadingOutlined v-else-if="props.status === 'running'" />
    <CheckCircleFilled style="color:green" v-else-if="props.status === 'success'" />
    <CloseCircleFilled style="color:darkred" v-else-if="props.status === 'fail'" />
    <CloseCircleFilled style="color:darkred" v-else-if="props.status === 'timeout'" />
    <CloseCircleFilled style="color:darkred" v-else-if="props.status === 'cancel'" />
    <ClockCircleOutlined v-else-if="props.status === 'unExecuted'"/>
    <ClockCircleOutlined v-else-if="props.status === 'unknown'"/>
    <span style="margin-left: 6px">{{t("workflow.status." + props.status)}}</span>
  </div>
  <div class="no-wrap" v-else>
    <ClockCircleOutlined />
    <span style="margin-left: 6px">{{t("workflow.status.unExecuted")}}</span>
  </div>
</template>

<script setup>
import {
  CheckCircleFilled,
  CloseCircleFilled,
  LoadingOutlined,
  ClockCircleOutlined
} from "@ant-design/icons-vue";
import { defineProps } from "vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const props = defineProps(["status"]);
</script>

<style scoped>
</style>