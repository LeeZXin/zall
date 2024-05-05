<template>
  <div class="navi-back" :style="props.style">
    <span class="header" @click="goTo">
      <arrow-left-outlined />
      <span style="margin-left:8px">{{props.name}}</span>
    </span>
  </div>
</template>
<script setup>
import { ArrowLeftOutlined } from "@ant-design/icons-vue";
import { defineProps } from "vue";
import { useRouter } from "vue-router";
const router = useRouter();
const props = defineProps(["style", "url", "name"]);
const goTo = () => {
  router.push(props.url);
};
</script>
<style scoped>
.navi-back {
  margin-bottom: 10px;
}
.navi-back > .header {
  font-size: 14px;
  cursor: pointer;
  margin-bottom: 10px;
}
.navi-back > .header:hover {
  color: #1677ff;
}
</style>