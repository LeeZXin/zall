<template>
  <a-dropdown>
    <template #overlay>
      <a-menu @click="selectLang">
        <a-menu-item key="zh">中文</a-menu-item>
        <a-menu-item key="en">English</a-menu-item>
      </a-menu>
    </template>
    <div class="lang no-wrap" :style="props.style">{{localeText}}</div>
  </a-dropdown>
</template>
<script setup>
import { ref, defineProps } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";
const props = defineProps(["style"])
const route = useRoute();
const { locale } = useI18n();
const localeMap = {
  zh: "中文",
  en: "English"
};
// 点击选择语言
const selectLang = event => {
  localeText.value = localeMap[event.key];
  locale.value = event.key;
};
if (route.query.locale) {
  locale.value = route.query.locale;
}
const localeText = ref(localeMap[locale.value]);
</script>
<style scoped>
.lang {
  font-size: 14px;
  color: white;
  line-height: 64px;
}
</style>