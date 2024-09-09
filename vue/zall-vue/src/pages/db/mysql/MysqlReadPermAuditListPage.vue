<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-select
        v-model:value="applyStatus"
        @change="selectApplyStatus"
        style="margin-right:6px;width:180px"
      >
        <a-select-option :value="1">{{t("mysqlReadPermApply.pendingStatus")}}</a-select-option>
        <a-select-option :value="2">{{t("mysqlReadPermApply.agreeStatus")}}</a-select-option>
        <a-select-option :value="3">{{t("mysqlReadPermApply.disagreeStatus")}}</a-select-option>
        <a-select-option :value="4">{{t("mysqlReadPermApply.canceledStatus")}}</a-select-option>
      </a-select>
      <a-select style="width: 180px" v-model:value="selectedDbId" @change="searchApply">
        <a-select-option :value="0">{{t("mysqlReadPermApply.allDatabases")}}</a-select-option>
        <a-select-option
          :value="item.value"
          v-for="item in dbList"
          v-bind:key="item.value"
        >{{item.label}}</a-select-option>
      </a-select>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <StatusTag v-if="dataIndex === 'applyStatus'" :status="dataItem[dataIndex]" />
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="agreeApply(dataItem)">
                  <CheckOutlined />
                  <span style="margin-left:4px">{{t('mysqlReadPermApply.agree')}}</span>
                </li>
                <li @click="showDisagreeModal(dataItem)">
                  <CloseOutlined />
                  <span style="margin-left:4px">{{t('mysqlReadPermApply.disagree')}}</span>
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
      @change="()=>listApply()"
    />
    <a-modal
      v-model:open="disagreeModal.open"
      :title="t('mysqlReadPermApply.fillDisagreeReason')"
      @ok="disagreeApply"
    >
      <a-textarea
        style="width:100%"
        v-model:value="disagreeModal.reason"
        :auto-size="{ minRows: 3, maxRows: 3 }"
        :maxlength="255"
      />
    </a-modal>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import StatusTag from "@/components/db/MysqlReadPermApplyStatutsTag";
import {
  EllipsisOutlined,
  ExclamationCircleOutlined,
  CloseOutlined,
  CheckOutlined
} from "@ant-design/icons-vue";
import {
  listReadPermApplyByDbaRequest,
  getAllMysqlDbRequest,
  agreeReadPermApplyRequest,
  disagreeReadPermApplyRequest
} from "@/api/db/mysqlApi";
import { ref, createVNode, reactive } from "vue";
import { Modal, message } from "ant-design-vue";
import { dbApplyReasonRegexp } from "@/utils/regexp";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const dataSource = ref([]);
// 状态
const applyStatus = ref(1);
// 分页
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
const selectedDbId = ref(0);
// 不同意modal
const disagreeModal = reactive({
  id: 0,
  open: false,
  reason: ""
});
const dbList = ref([]);
const columns = ref([
  {
    i18nTitle: "mysqlReadPermApply.dbName",
    dataIndex: "dbName",
    key: "dbName",
    fix: "left"
  },
  {
    i18nTitle: "mysqlReadPermApply.accessBase",
    dataIndex: "accessBase",
    key: "accessBase"
  },
  {
    i18nTitle: "mysqlReadPermApply.accessTables",
    dataIndex: "accessTables",
    key: "accessTables"
  },
  {
    i18nTitle: "mysqlReadPermApply.expireDay",
    dataIndex: "expireDay",
    key: "expireDay"
  },
  {
    i18nTitle: "mysqlReadPermApply.applyStatus",
    dataIndex: "applyStatus",
    key: "applyStatus"
  },
  {
    i18nTitle: "mysqlReadPermApply.account",
    dataIndex: "account",
    key: "account"
  },
  {
    i18nTitle: "mysqlReadPermApply.applyReason",
    dataIndex: "applyReason",
    key: "applyReason",
    width: 160
  },
  {
    i18nTitle: "mysqlReadPermApply.applyTime",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "mysqlReadPermApply.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
]);

const searchApply = () => {
  dataPage.current = 1;
  listApply();
};

const selectApplyStatus = () => {
  switch (applyStatus.value) {
    case 1:
      columns.value = [
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
          dataIndex: "accessTables",
          key: "accessTables"
        },
        {
          i18nTitle: "mysqlReadPermApply.expireDay",
          dataIndex: "expireDay",
          key: "expireDay"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlReadPermApply.account",
          dataIndex: "account",
          key: "account"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlReadPermApply.applyTime",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlReadPermApply.operation",
          dataIndex: "operation",
          key: "operation",
          width: 130,
          fixed: "right"
        }
      ];
      break;
    case 2:
      columns.value = [
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
          dataIndex: "accessTables",
          key: "accessTables"
        },
        {
          i18nTitle: "mysqlReadPermApply.expireDay",
          dataIndex: "expireDay",
          key: "expireDay"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlReadPermApply.account",
          dataIndex: "account",
          key: "account"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlReadPermApply.auditor",
          dataIndex: "auditor",
          key: "auditor"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyTime",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlReadPermApply.auditTime",
          dataIndex: "updated",
          key: "updated"
        }
      ];
      break;
    case 3:
      columns.value = [
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
          dataIndex: "accessTables",
          key: "accessTables"
        },
        {
          i18nTitle: "mysqlReadPermApply.expireDay",
          dataIndex: "expireDay",
          key: "expireDay"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlReadPermApply.account",
          dataIndex: "account",
          key: "account"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlReadPermApply.disagreeReason",
          dataIndex: "disagreeReason",
          key: "disagreeReason",
          width: 160
        },
        {
          i18nTitle: "mysqlReadPermApply.auditor",
          dataIndex: "auditor",
          key: "auditor"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyTime",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlReadPermApply.auditTime",
          dataIndex: "updated",
          key: "updated"
        }
      ];
      break;
    case 4:
      columns.value = [
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
          dataIndex: "accessTables",
          key: "accessTables"
        },
        {
          i18nTitle: "mysqlReadPermApply.expireDay",
          dataIndex: "expireDay",
          key: "expireDay"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlReadPermApply.account",
          dataIndex: "account",
          key: "account"
        },
        {
          i18nTitle: "mysqlReadPermApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlReadPermApply.applyTime",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlReadPermApply.cancelTime",
          dataIndex: "updated",
          key: "updated"
        }
      ];
      break;
  }
  searchApply();
};

const listApply = () => {
  listReadPermApplyByDbaRequest({
    dbId: selectedDbId.value,
    pageNum: dataPage.current,
    applyStatus: applyStatus.value
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

const agreeApply = item => {
  Modal.confirm({
    title: `${t("mysqlReadPermApply.confirmAgree")} ${item.account}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      agreeReadPermApplyRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchApply();
      });
    },
    onCancel() {}
  });
};

const showDisagreeModal = item => {
  disagreeModal.id = item.id;
  disagreeModal.open = true;
  disagreeModal.reason = "";
};

const disagreeApply = () => {
  if (!dbApplyReasonRegexp.test(disagreeModal.reason)) {
    message.warn(t("mysqlReadPermApply.disagreeReasonFormatErr"));
    return;
  }
  disagreeReadPermApplyRequest({
    applyId: disagreeModal.id,
    disagreeReason: disagreeModal.reason
  }).then(() => {
    message.success(t("operationSuccess"));
    listApply();
    disagreeModal.open = false;
  });
};

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
getAllDb();
listApply();
</script>
<style scoped>
</style>