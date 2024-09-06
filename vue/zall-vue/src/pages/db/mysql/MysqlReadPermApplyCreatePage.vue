<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span>{{t('mysqlReadPermApply.title')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('mysqlReadPermApply.selectHost')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.dbId" :options="dbList" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('mysqlReadPermApply.fillAccessBaseName')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.accessBase" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('mysqlReadPermApply.fillAccessTablesName')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.accessTables" />
          <div class="input-desc">{{t('mysqlReadPermApply.accessTablesDesc')}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('mysqlReadPermApply.expireDay')}}</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.expireDay">
            <a-radio :value="1">{{t('mysqlReadPermApply.oneDay')}}</a-radio>
            <a-radio :value="30">{{t('mysqlReadPermApply.oneMonth')}}</a-radio>
            <a-radio :value="90">{{t('mysqlReadPermApply.threeMonth')}}</a-radio>
            <a-radio :value="180">{{t('mysqlReadPermApply.sixMonth')}}</a-radio>
            <a-radio :value="365">{{t('mysqlReadPermApply.oneYear')}}</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('mysqlReadPermApply.fillApplyReason')}}</div>
        <div class="section-body">
          <a-textarea
            style="width:100%"
            v-model:value="formState.applyReason"
            :auto-size="{ minRows: 3, maxRows: 3 }"
            :maxlength="255"
          />
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="applyReadPerm">{{t('mysqlReadPermApply.apply')}}</a-button>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
// 数据库列表
const dbList = ref([]);
// 表单数据
const formState = reactive({
  dbId: 0,
  expireDay: 1,
  accessBase: "",
  accessTables: "",
  applyReason: ""
});
// 提交申请
const applyReadPerm = () => {
  if (formState.dbId <= 0) {
    message.warn(t("mysqlReadPermApply.pleaseSelectDatabse"));
    return;
  }
  if (!dbAccessBaseRegexp.test(formState.accessBase)) {
    message.warn(t("mysqlReadPermApply.accessBaseFormatErr"));
    return;
  }
  if (!dbAccessTablesRegexp.test(formState.accessTables)) {
    message.warn(t("mysqlReadPermApply.accessTablesFormatErr"));
    return;
  }
  if (!dbApplyReasonRegexp.test(formState.applyReason)) {
    message.warn(t("mysqlReadPermApply.applyReasonFormatErr"));
    return;
  }
  applyReadPermRequest({
    dbId: formState.dbId,
    accessBase: formState.accessBase,
    accessTables: formState.accessTables,
    expireDay: formState.expireDay,
    applyReason: formState.applyReason
  }).then(() => {
    message.success(t("operationSuccess"));
    router.push(`/db/mysqlReadPermApply/list`);
  });
};
// 获取所有的数据库
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