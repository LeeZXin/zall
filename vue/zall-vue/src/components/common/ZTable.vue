<template>
  <div class="z-table-section">
    <div v-if="props.label" class="z-label">{{props.label}}</div>
    <div class="z-table" :style="props.style">
      <table>
        <tr>
          <td
            v-for="item in props.columns"
            v-bind:key="item.key"
          >{{item.i18nTitle?t(item.i18nTitle):item.title}}</td>
        </tr>
        <tr v-for="dataItem in props.dataSource" v-bind:key="dataItem.key">
          <td v-for="(columnItem, index) in columns" v-bind:key="dataItem.key + index">
            <CellRender :dataIndex="columnItem.dataIndex" :dataItem="dataItem" />
          </td>
        </tr>
        <tr v-if="props?.dataSource?.length === 0">
          <td :colspan="props.columns.length">
            <div class="z-no-data">
              <div style="font-size:30px;text-align:center;margin-bottom:14px">
                <InboxOutlined />
              </div>
              <NoDataRender v-if="slots.noData" />
              <div v-else style="text-align:center;font-size:14px;">
                <span>{{t("ztable.noDataText")}}</span>
              </div>
            </div>
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>
<script setup>
import { defineProps, h, useSlots } from "vue";
import { InboxOutlined } from "@ant-design/icons-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const props = defineProps(["columns", "dataSource", "style", "label"]);
const slots = useSlots();
const CellRender = params => {
  return slots.bodyCell
    ? slots.bodyCell({
        dataIndex: params.dataIndex,
        dataItem: params.dataItem
      })
    : h("span", null, params.dataItem[params.dataIndex]);
};
const NoDataRender = () => {
  return slots.noData ? slots.noData() : h("div");
};
</script>
<style scoped>
.z-table-section + .z-table-section {
  margin-top: 10px;
}
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
  margin-bottom: 10px;
}
.z-no-data {
  width: 100%;
  padding: 20px 0;
  border-radius: 4px;
}
</style>