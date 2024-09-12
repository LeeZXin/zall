<template>
  <div class="body" :style="props.style">
    <div class="copy-icon" @click="copy()">
      <CopyOutlined />
    </div>
    <div class="sha">
      <span>{{children}}</span>
    </div>
  </div>
</template>
<script setup>
import { CopyOutlined } from "@ant-design/icons-vue";
import { defineProps, useSlots } from "vue";
import { message } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const props = defineProps(["style"]);
const slots = useSlots();
const children = slots.default?.()[0].children;
const copy = () => {
  message.success(t("copySuccess"));
  window.navigator.clipboard.writeText(children);
};
</script>
<style scoped>
.body {
  display: inline-block;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  background-color: white;
}
.copy-icon,
.sha {
  height: 24px;
  line-height: 24px;
  font-size: 14px;
  display: inline-block;
  padding: 0px 8px;
}
.sha {
  display: inline-block;
  border-left: 1px solid #d9d9d9;
}
.copy-icon:hover {
  background-color: #f0f0f0;
  cursor: pointer;
}
</style>