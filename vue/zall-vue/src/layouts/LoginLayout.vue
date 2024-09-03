<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{t("app")}}</span>
      <I18nSelect style="float:right;" />
    </a-layout-header>
    <a-layout-content style="height: calc(100vh - 128px)">
      <router-view />
    </a-layout-content>
    <a-layout-footer style="text-align: center">{{t("footer")}}</a-layout-footer>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
</script>