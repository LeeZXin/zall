<template>
  <div class="no-data" :style="props.style">
    <div style="font-size:30px;text-align:center;margin-bottom:20px">
      <InboxOutlined />
    </div>
    <SlotRender v-if="slots.desc" />
    <div v-else style="text-align:center;font-size:14px;font-weight:bold">
      <span>{{t("noData")}}</span>
    </div>
  </div>
</template>
<script setup>
import { defineProps, useSlots, h } from "vue";
import { InboxOutlined } from "@ant-design/icons-vue";
import { useI18n } from "vue-i18n"; 
/*
  无数据展示
*/
const { t } = useI18n();
const slots = useSlots();
const props = defineProps(["style"]);
const SlotRender = () => {
  return slots.desc ? slots.desc() : h("div");
};
</script>
<style scoped>
.no-data {
  width: 100%;
  padding: 40px 0;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
</style>