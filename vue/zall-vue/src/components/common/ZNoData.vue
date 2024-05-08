<template>
  <div class="no-data" :style="props.style">
    <div style="font-size:30px;text-align:center;margin-bottom:20px">
      <InboxOutlined />
    </div>
    <SlotRender v-if="slots.desc" />
    <div v-else style="text-align:center;font-size:14px;">
      <span>暂无数据</span>
    </div>
  </div>
</template>
<script setup>
import { defineProps, useSlots, h } from "vue";
import { InboxOutlined } from "@ant-design/icons-vue";
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