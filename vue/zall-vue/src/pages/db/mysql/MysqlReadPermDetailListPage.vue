<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-select style="width: 180px" v-model:value="selectedDbId" @change="selectDbIdChange">
        <a-select-option :value="0">{{t("mysqlReadPermApply.allDatabases")}}</a-select-option>
        <a-select-option
          :value="item.value"
          v-for="item in dbList"
          v-bind:key="item.value"
        >{{item.label}}</a-select-option>
      </a-select>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="getApply(dataItem)">
                  <EyeOutlined />
                  <span style="margin-left:4px">{{t('mysqlReadPermApply.viewApprovalForm')}}</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </div>
      </template>
    </ZTable>
    <a-pagination
      v-model:current="dataPage.current"
      :total="dataPage.totalCount"
      show-less-items
      :pageSize="dataPage.pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listPerm()"
    />
    <a-modal v-model:open="apply.open" :title="t('mysqlReadPermApply.approvalForm')" :footer="null">
      <ul class="apply-ul">
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.dbName')}}</div>
          <div class="item-value">{{apply.dbName}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.accessBase')}}</div>
          <div class="item-value">{{apply.accessBase}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.accessTables')}}</div>
          <div class="item-value">{{apply.accessTables}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.auditor')}}</div>
          <div class="item-value">{{apply.auditor}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.applyReaason')}}</div>
          <div class="item-value">{{apply.applyReason}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.applyTime')}}</div>
          <div class="item-value">{{apply.created}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.auditTime')}}</div>
          <div class="item-value">{{apply.updated}}</div>
        </li>
      </ul>
    </a-modal>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import {
  listReadPermByOperatorRequest,
  getAllMysqlDbRequest,
  getReadPermApplyRequest
} from "@/api/db/mysqlApi";
import { EyeOutlined, EllipsisOutlined } from "@ant-design/icons-vue";
import { ref, reactive } from "vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const dataSource = ref([]);
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
const selectedDbId = ref(0);
// 数据库列表
const dbList = ref([]);
// 审批单
const apply = reactive({
  open: false,
  dbName: "",
  accessBase: "",
  accessTables: "",
  expireDay: 0,
  applyStatus: 0,
  auditor: "",
  applyReason: "",
  created: "",
  updated: ""
});
// 数据项
const columns = ref([
  {
    i18nTitle: "mysqlReadPermApply.dbName",
    dataIndex: "dbName",
    key: "dbName"
  },
  {
    i18nTitle: "mysqlReadPermApply.accessBase",
    dataIndex: "accessBase",
    key: "accessBase"
  },
  {
    i18nTitle: "mysqlReadPermApply.accessTables",
    dataIndex: "accessTable",
    key: "accessTable"
  },
  {
    i18nTitle: "mysqlReadPermApply.effectiveTime",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "mysqlReadPermApply.expireTime",
    dataIndex: "expired",
    key: "expired"
  },
  {
    i18nTitle: "mysqlReadPermApply.operation",
    dataIndex: "operation",
    key: "operation"
  }
]);
// 权限列表
const listPerm = () => {
  listReadPermByOperatorRequest({
    dbId: selectedDbId.value,
    pageNum: dataPage.current
  }).then(res => {
    dataPage.totalCount = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};
// 所有数据库
const getAllDb = () => {
  getAllMysqlDbRequest().then(res => {
    dbList.value = dbList.value.concat(
      res.data.map(item => {
        return {
          value: item.id,
          label: item.name
        };
      })
    );
  });
};
// 选择数据库
const selectDbIdChange = () => {
  dataPage.current = 1;
  listPerm();
};
// 获取审批单
const getApply = item => {
  getReadPermApplyRequest(item.applyId).then(res => {
    let data = res.data;
    apply.dbName = data.dbName;
    apply.accessBase = data.accessBase;
    apply.accessTables = data.accessTables;
    apply.expireDay = data.expireDay;
    apply.applyStatus = data.applyStatus;
    apply.auditor = data.auditor;
    apply.applyReason = data.applyReason;
    apply.created = data.created;
    apply.updated = data.updated;
    apply.open = true;
  });
};
getAllDb();
listPerm();
</script>
<style scoped>
.apply-ul {
  width: 100%;
  padding-bottom: 20px;
}
.apply-ul > li {
  width: 100%;
}
.apply-ul > li + li {
  margin-top: 12px;
}
.item-name {
  font-size: 12px;
  margin-bottom: 4px;
}
.item-value {
  font-size: 14px;
  line-height: 18px;
  padding-left: 20px;
  min-height: 18px;
}
</style>