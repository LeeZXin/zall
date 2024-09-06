<template>
  <div class="flex-center" :style="props.style">
    <a-dropdown>
      <template #overlay>
        <a-menu @click="handleMenuClick">
          <a-menu-item v-for="item in envList" v-bind:key="item.value">
            <span :class="{'item-selected': item.value === selectedEnv}">{{item.label}}</span>
          </a-menu-item>
        </a-menu>
      </template>
      <a-button :icon="h(SwapOutlined)">
        <span>{{t('switchEnv')}}</span>
      </a-button>
    </a-dropdown>
  </div>
</template>
<script setup>
import { ref, defineProps, watch, defineEmits } from "vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import { SwapOutlined } from "@ant-design/icons-vue";
import { h } from "vue";
import { useI18n } from "vue-i18n";
/*
  环境选择下拉框 
*/
const { t } = useI18n();
// 默认环境值
const props = defineProps(["defaultEnv"]);
// 监听@change
const emit = defineEmits(["change"]);
// 当前选择的环境
const selectedEnv = ref("");
// 环境列表
const envList = ref([]);
// 获取环境列表
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
const handleMenuClick = event => {
  selectedEnv.value = event.key;
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
.item-selected {
  color: #1677ff;
}
</style>