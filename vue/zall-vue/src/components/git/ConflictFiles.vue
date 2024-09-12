<template>
  <div class="conflict-files">
    <div style="color:darkred;padding-bottom:6px">
      <warning-outlined />
      <span style="padding-left:6px">{{t('existConflictFiles')}}</span>
    </div>
    <div class="conflict-items">
      <span v-for="item in props.conflictFiles" v-bind:key="item">{{item}}</span>
    </div>
  </div>
</template>
<script setup>
import { defineProps } from "vue";
import { WarningOutlined } from "@ant-design/icons-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const props = defineProps(["conflictFiles"]);
</script>
<style scoped>
.conflict-files {
  padding: 16px;
  font-size: 14px;
  word-break: break-all;
  width: 100%;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  margin-top: 10px;
}
.conflict-items > span {
  padding-left: 4px;
  line-height: 18px;
  font-size: 12px;
}
</style>