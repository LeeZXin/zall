<template>
  <div class="no-data no-wrap" :style="props.style">暂无数据</div>
</template>
<script setup>
import { defineProps } from "vue";

const props = defineProps(["style"]);
</script>
<style scoped>
.no-data {
  width: 100%;
  text-align: center;
  padding: 40px 0;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
</style>