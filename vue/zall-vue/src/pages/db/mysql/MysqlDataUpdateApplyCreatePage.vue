<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span>{{t("mysqlDataUpdateApply.title")}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t("mysqlDataUpdateApply.selectHost")}}</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.dbId"
            :options="dbList"
            show-search
            :filter-option="filterDbListOption"
          />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t("mysqlDataUpdateApply.fillAccessBaseName")}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.accessBase" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t("mysqlDataUpdateApply.fillApplyReason")}}</div>
        <div class="section-body">
          <a-textarea
            style="width:100%"
            v-model:value="formState.applyReason"
            :auto-size="{ minRows: 3, maxRows: 3 }"
            :maxlength="255"
          />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t("mysqlDataUpdateApply.fillSql")}}</div>
        <div class="section-body">
          <Codemirror
            v-model="formState.cmd"
            style="height:280px;width:100%"
            :extensions="extensions"
          />
        </div>
      </div>
      <div class="section-item">
        <a-checkbox v-model:checked="formState.executeImmediatelyAfterApproval">
          <div>{{t("mysqlDataUpdateApply.executeImmediatelyAfterApproval")}}</div>
        </a-checkbox>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="applyDataUpdate">{{t("mysqlDataUpdateApply.apply")}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
/*
  数据库修改单申请页面
*/
import { reactive, ref } from "vue";
import {
  dbAccessBaseRegexp,
  dbDataUpdateCmdRegexp,
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const extensions = [sql()];
const router = useRouter();
// 数据库列表
const dbList = ref([]);
// 表单数据
const formState = reactive({
  dbId: 0,
  accessBase: "",
  cmd: "",
  applyReason: "",
  executeImmediatelyAfterApproval: false
});
// 提交申请
const applyDataUpdate = () => {
  if (formState.dbId <= 0) {
    message.warn(t("mysqlDataUpdateApply.pleaseSelectDatabse"));
    return;
  }
  if (!dbAccessBaseRegexp.test(formState.accessBase)) {
    message.warn(t("mysqlDataUpdateApply.accessBaseFormatErr"));
    return;
  }
  if (!dbDataUpdateCmdRegexp.test(formState.cmd)) {
    message.warn(t("mysqlDataUpdateApply.sqlFormatErr"));
    return;
  }
  if (!dbApplyReasonRegexp.test(formState.applyReason)) {
    message.warn(t("mysqlDataUpdateApply.applyReasonFormatErr"));
    return;
  }
  applyDataUpdateRequest({
    dbId: formState.dbId,
    accessBase: formState.accessBase,
    cmd: formState.cmd,
    applyReason: formState.applyReason,
    executeImmediatelyAfterApproval: formState.executeImmediatelyAfterApproval
  }).then(res => {
    if (res.data.allPass) {
      message.success(t("operationSuccess"));
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
// 获取所有数据库列表
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
// 下拉框过滤
const filterDbListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
getAllDb();
</script>
<style scoped>
</style>