<template>
  <div class="flex-center" :style="props.style">
    <span style="margin-right:6px">环境:</span>
    <a-select
      style="width: 200px"
      placeholder="选择环境"
      v-model:value="selectedEnv"
      :options="envList"
    />
  </div>
</template>

<script setup>
import { ref, defineProps, watch, defineEmits } from "vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
const props = defineProps(["defaultEnv"]);
const emit = defineEmits(["change"]);
const selectedEnv = ref("");
const envList = ref([]);

const getEnvList = () => {
  getEnvCfgRequest().then(res => {
    envList.value = res.data.map(item => {
      return {
        value: item,
        label: item
      };
    });
    if (props.defaultEnv && res.data?.includes(props.defaultEnv)) {
      selectedEnv.value = props.defaultEnv;
    } else if (res.data.length > 0) {
      selectedEnv.value = res.data[0];
    }
  });
};

watch(
  () => selectedEnv.value,
  (newVal, oldVal) => {
    emit("change", {
      newVal,
      oldVal
    });
  }
);

getEnvList();
</script>

<style scoped>
</style>