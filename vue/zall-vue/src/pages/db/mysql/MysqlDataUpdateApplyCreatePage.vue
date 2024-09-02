<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span>申请数据库修改单</span>
      </div>
      <div class="section">
        <div class="section-title">选择数据库</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.dbId"
            :options="dbList"
            show-search
            :filter-option="filterDbListOption"
          />
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
      <div class="section">
        <div class="section-title">执行sql</div>
        <div class="section-body">
          <Codemirror
            v-model="formState.cmd"
            style="height:280px;width:100%"
            :extensions="extensions"
          />
        </div>
      </div>
      <div class="form-item">
        <a-checkbox v-model:checked="formState.executeWhenApply">
          <div>是否立即执行</div>
        </a-checkbox>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="applyDataUpdate">立即申请</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import {
  dbAccessBaseRegexp,
  dbDatUpdateCmdRegexp,
  dbApplyReasonRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import {
  getAllMysqlDbRequest,
  applyDataUpdateRequest
} from "@/api/db/mysqlApi";
import { useRouter } from "vue-router";
import { Codemirror } from "vue-codemirror";
import { sql } from "@codemirror/lang-sql";
const extensions = [sql()];
const router = useRouter();
const dbList = ref([]);
const formState = reactive({
  dbId: 0,
  accessBase: "",
  cmd: "",
  applyReason: "",
  executeWhenApply: false
});

const applyDataUpdate = () => {
  if (!dbAccessBaseRegexp.test(formState.accessBase)) {
    message.warn("申请库格式错误");
    return;
  }
  if (!dbDatUpdateCmdRegexp.test(formState.accessTables)) {
    message.warn("申请表格式错误");
    return;
  }
  if (!dbApplyReasonRegexp.test(formState.applyReason)) {
    message.warn("申请原因格式错误");
    return;
  }
  applyDataUpdateRequest({
    dbId: formState.dbId,
    accessBase: formState.accessBase,
    cmd: formState.cmd,
    applyReason: formState.applyReason,
    executeWhenApply: formState.executeWhenApply
  }).then(res => {
    if (res.data.allPass) {
      message.success("申请成功");
      router.push(`/db/mysqlDataUpdateApply/list`);
      return;
    }
    if (res.data.results?.length > 0) {
      let item = res.data.results.find(item => item.pass === false);
      if (item) {
        message.warn(item.errMsg);
        return;
      }
    }
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
const filterDbListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
getAllDb();
</script>
<style scoped>
</style>