<template>
  <div :style="props.style" :class="{
    'no-data-bordered': !props.unbordered
    }">
    <a-empty :image="Empty.PRESENTED_IMAGE_SIMPLE" />
  </div>
</template>
<script setup>
const props = defineProps(["unbordered", "style"]);
import { Empty } from "ant-design-vue";
import { defineProps } from "vue";
/*
  无数据展示
*/
</script>
<style scoped>
.no-data {
  width: 100%;
  padding: 40px;
  color: #555555;
}
.no-data-bordered {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
</style>