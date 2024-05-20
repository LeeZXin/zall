<template>
  <a-button type="link" @click="toPrDetail">#{{props.prId}}</a-button>
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