<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>ZALL</span>
      <I18nSelect style="float:right;" />
    </a-layout-header>
    <a-layout-content style="height: calc(100vh - 128px)">
      <router-view/>
    </a-layout-content>
    <a-layout-footer style="text-align: center">Powered by zexin</a-layout-footer>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
</script>
<style scoped>
.layout {
  width: calc(100vw);
  height: calc(100vh);
  background-color: #1e90ff;
  overflow: hidden;
}
</style>