<template>
  <div>
    <div v-if="props.label" class="z-label">{{props.label}}</div>
    <div class="z-table" :style="props.style">
      <table>
        <tr>
          <td v-for="item in props.columns" v-bind:key="item.key">{{item.title}}</td>
        </tr>
        <tr v-for="dataItem in props.dataSource" v-bind:key="dataItem.key">
          <td v-for="(columnItem, index) in columns" v-bind:key="dataItem.key + index">
            <CellRender :dataIndex="columnItem.dataIndex" :dataItem="dataItem" />
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>
<script setup>
import { defineProps, h, useSlots } from "vue";
const props = defineProps(["columns", "dataSource", "style", "label"]);
const slots = useSlots();
const CellRender = event => {
  return slots.bodyCell
    ? slots.bodyCell({
        dataIndex: event.dataIndex,
        dataItem: event.dataItem
      })
    : h("span", null, event.dataItem[event.dataIndex]);
};
</script>
<style scoped>
.z-table {
  margin-top: 10px;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
}
.z-table,
.z-table > table {
  width: 100%;
}
.z-table > table > tr > td {
  background-color: white;
}
.z-table > table > tr:first-child > td {
  background-color: #f0f0f0;
}

.z-table > table > tr:first-child > td:last-child {
  border-top-right-radius: 4px;
}

.z-table > table > tr:first-child > td:first-child {
  border-top-left-radius: 4px;
}

.z-table > table > tr:last-child > td:last-child {
  border-bottom-right-radius: 4px;
}

.z-table > table > tr:last-child > td:first-child {
  border-bottom-left-radius: 4px;
}

.z-table > table > tr > td {
  text-align: center;
  padding: 0 8px;
}
.z-table > table > tr {
  line-height: 32px;
  font-size: 14px;
}
.z-table > table > tr + tr {
  border-top: 1px solid #d9d9d9;
}
.z-table > table > tr:first-child {
  line-height: 32px;
  font-size: 14px;
  font-weight: bold;
}
.z-label {
  font-size: 16px;
  font-weight: bold;
  line-height: 32px;
}
</style>