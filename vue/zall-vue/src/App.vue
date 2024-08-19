<template>
  <a-config-provider :locale="antdLocale">
    <router-view />
  </a-config-provider>
</template>
<script setup>
import { useI18n } from "vue-i18n";
import enUS from "ant-design-vue/es/locale/en_US";
import zhCN from "ant-design-vue/es/locale/zh_CN";
import { ref, watch } from "vue";
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
const { locale } = useI18n();
const antdLocale = ref();
// 国际化配置
const changeLocale = val => {
  if (val === "zh") {
    antdLocale.value = zhCN;
    dayjs.locale(zhCN.locale);
  } else if (val === "en") {
    antdLocale.value = enUS;
    dayjs.locale(enUS.locale);
  }
};
// 监听配置
watch(locale, newVal => {
  changeLocale(newVal);
});
changeLocale(locale.value);
</script>
<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  height: calc(100vh);
  width: calc(100vw);
  overflow: scroll;
}
</style>