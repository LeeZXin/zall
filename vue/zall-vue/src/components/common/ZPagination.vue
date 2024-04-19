<template>
  <div class="z-pagination" :style="props.style">
    <div
      class="z-pagination-flex"
      :class="{
      'z-pagination-align-right': props.placement === 'right',
      'z-pagination-align-center': props.placement === 'center'
  }"
    >
      <div
        @click="onclick('last')"
        class="z-pagination-btn"
        :class="{
          'z-pagination-btn-active': !props.disableLastPage,
          'z-pagination-btn-inactive': props.disableLastPage
      }"
      >
        <left-outlined />
        <span style="padding-left:4px">{{t("lastPage")}}</span>
      </div>
      <div
        @click="onclick('next')"
        class="z-pagination-btn"
        :class="{
          'z-pagination-btn-active': !props.disableNextPage,
          'z-pagination-btn-inactive': props.disableNextPage
      }"
      >
        <span style="padding-right:4px">{{t("nextPage")}}</span>
        <right-outlined />
      </div>
    </div>
  </div>
</template>
<script setup>
import { LeftOutlined, RightOutlined } from "@ant-design/icons-vue";
import { useI18n } from "vue-i18n";
import { defineProps, defineEmits } from "vue";
const props = defineProps([
  "style",
  "disableLastPage",
  "disableNextPage",
  "placement"
]);
const emit = defineEmits(["change"]);
const { t } = useI18n();
const onclick = key => {
  if (
    (key === "last" && !props.disableLastPage) ||
    (key === "next" && !props.disableNextPage)
  ) {
    emit("change", key);
  }
};
</script>
<style scoped>
.z-pagination {
  margin-top: 10px;
  height: 34px;
  width: 100%;
}
.z-pagination-flex {
  background-color: white;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  display: flex;
  align-items: center;
  width: 180px;
}
.z-pagination-btn,
.z-pagination-btn {
  width: 50%;
  overflow: hidden;
  white-space: nowrap;
  height: 32px;
  line-height: 32px;
  text-align: center;
  font-size: 14px;
}
.z-pagination-btn-active {
  cursor: pointer;
  color: black;
}
.z-pagination-btn-inactive {
  cursor: not-allowed;
  color: gray;
}
.z-pagination-btn-active:hover {
  color: #1677ff;
}
.z-pagination-btn + .z-pagination-btn {
  border-left: 1px solid #d9d9d9;
}
.z-pagination-align-right {
  float: right;
}
.z-pagination-align-center {
  margin-left: 50%;
  transform: translateX(-50%);
}
</style>