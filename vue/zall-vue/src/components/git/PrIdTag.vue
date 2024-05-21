<template>
  <span class="link" @click="toPrDetail">#{{props.prId}}</span>
</template>
<script setup>
import { defineProps } from "vue";
import { useRouter } from "vue-router";
const router = useRouter();
const props = defineProps(["repoId", "prId"]);
const toPrDetail = () => {
  router.push(`/gitRepo/${props.repoId}/pullRequest/${props.prId}/detail`);
};
</script>
<style scoped>
.link {
    font-size: 14px;
    color: #1677ff;
    cursor: pointer;
}
</style>
