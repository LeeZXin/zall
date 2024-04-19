<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{t("app")}}</span>
      <I18nSelect style="float:right;" />
    </a-layout-header>
    <a-layout-content style="height: calc(100vh - 128px)">
      <div class="not-found">:( {{t("notFound")}}</div>
      <div class="back-to-index" @click="backToIndex">{{t("backToIndex")}}</div>
    </a-layout-content>
    <a-layout-footer style="text-align: center">{{t("footer")}}</a-layout-footer>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
const { t } = useI18n();
const router = useRouter();
const backToIndex = () => {
    router.push("/index");
};
</script>
<style scoped>
.not-found {
  font-size: 32px;
  line-height: 64px;
  text-align: center;
  margin-top: 40px;
}
.back-to-index {
  line-height: 32px;
  text-align: center;
  font-size: 18px;
  cursor: pointer;
  margin-top: 20px;
}
.back-to-index:hover {
  color: #1677ff;
}
</style>