<template>
  <div class="z-table">
    <a-table
      :columns="props.columns"
      :dataSource="props.dataSource"
      :pagination="false"
      :style="props.style"
      :scroll="props.scroll"
      bordered
      size="middle"
    >
      <template v-if="props.label" #title>{{props.label}}</template>
      <template #headerCell="{column}">
        <template v-for="(item, index) in props.columns" :key="index">
          <span
            v-if="column.dataIndex === item.dataIndex"
          >{{item.i18nTitle?t(item.i18nTitle):item.title}}</span>
        </template>
      </template>
      <template #bodyCell="{column, record}">
        <CellRender :dataIndex="column.dataIndex" :dataItem="record" />
      </template>
    </a-table>
  </div>
</template>
<script setup>
import { defineProps, h, useSlots } from "vue";
import { useI18n } from "vue-i18n";
/*
  table表格 支持i18n转换
*/
const { t } = useI18n();
const props = defineProps([
  "columns",
  "dataSource",
  "style",
  "label",
  "scroll"
]);
const slots = useSlots();
const CellRender = params => {
  return slots.bodyCell
    ? slots.bodyCell({
        dataIndex: params.dataIndex,
        dataItem: params.dataItem
      })
    : h("span", null, params.dataItem[params.dataIndex]);
};
</script>
<style scoped>
</style>