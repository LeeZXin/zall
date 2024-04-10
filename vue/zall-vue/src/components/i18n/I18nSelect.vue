<template>
  <a-popover v-model:open="visible" trigger="click" placement="bottomRight">
    <template #content>
      <div
        v-for="(val, key) in localeMap"
        @click="selectLang(key)"
        v-bind:key="key"
        class="select-item"
      >{{val}}</div>
    </template>
    <span class="text-btn" :style="props.style">{{localeText}}</span>
  </a-popover>
</template>
<script setup>
import { ref, defineProps } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";
const currentRoute = useRoute();
const { locale } = useI18n();
const visible = ref(false);
const props = defineProps(["style", "locale"]);
const localeMap = {
  zh: "中文",
  en: "English"
};
const selectLang = key => {
  localeText.value = localeMap[key];
  visible.value = false;
  locale.value = key;
};
if (currentRoute.query.locale) {
  locale.value = currentRoute.query.locale;
}
const localeText = ref(localeMap[locale.value]);
</script>
<style scoped>
.text-btn {
  color: white;
  cursor: pointer;
  font-size: 14px;
}
.select-item {
  line-height: 28px;
  width: 80px;
  text-align: center;
  cursor: pointer;
}
.select-item:hover {
  background-color: #f0f0f0;
}
</style>