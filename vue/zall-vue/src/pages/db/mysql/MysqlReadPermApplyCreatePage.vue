<template>
  <div style="padding:14px">
    <div class="container">
      <div class="title">
        <span>申请数据库读权限</span>
      </div>
      <div class="section">
        <div class="section-title">选择数据库</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.dbId" :options="dbList" />
          <div class="input-desc">选择一个数据库</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">申请库</div>
        <div class="section-body">
          <a-input v-model:value="formState.accessBase" placeholder="请填写" />
          <div class="input-desc">填写申请库</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">申请表</div>
        <div class="section-body">
          <a-input v-model:value="formState.accessTables" placeholder="请填写" />
          <div class="input-desc">填写申请的表, 可用*代替全部, 多个表用;隔开</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">时效</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.expireDay">
            <a-radio :value="1">一天</a-radio>
            <a-radio :value="30">一个月</a-radio>
            <a-radio :value="90">三个月</a-radio>
            <a-radio :value="180">半年</a-radio>
            <a-radio :value="365">一年</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section">
        <div class="section-title">申请原因</div>
        <div class="section-body">
          <a-textarea
            style="width:100%"
            v-model:value="formState.applyReason"
            :auto-size="{ minRows: 3, maxRows: 3 }"
            :maxlength="255"
            placeholder="请填写"
          />
          <div class="input-desc">填写申请原因</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="applyReadPerm">立即申请</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import {
  dbAccessBaseRegexp,
  dbAccessTablesRegexp,
  dbApplyReasonRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getAllMysqlDbRequest, applyReadPermRequest } from "@/api/db/mysqlApi";
import { useRouter } from "vue-router";
const router = useRouter();
const dbList = ref([]);
const formState = reactive({
  dbId: 0,
  expireDay: 1,
  accessBase: "",
  accessTables: "",
  applyReason: ""
});

const applyReadPerm = () => {
  if (!dbAccessBaseRegexp.test(formState.accessBase)) {
    message.warn("申请库格式错误");
    return;
  }
  if (!dbAccessTablesRegexp.test(formState.accessTables)) {
    message.warn("申请表格式错误");
    return;
  }
  if (!dbApplyReasonRegexp.test(formState.applyReason)) {
    message.warn("申请原因格式错误");
    return;
  }
  applyReadPermRequest({
    dbId: formState.dbId,
    accessBase: formState.accessBase,
    accessTables: formState.accessTables,
    expireDay: formState.expireDay,
    applyReason: formState.applyReason
  }).then(() => {
    message.success("申请成功");
    router.push(`/db/mysqlReadPermApply/list`);
  });
};
const getAllDb = () => {
  getAllMysqlDbRequest().then(res => {
    dbList.value = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
    if (res.data.length > 0) {
      formState.dbId = res.data[0].id;
    }
  });
};
getAllDb();
</script>
<style scoped>
</style>