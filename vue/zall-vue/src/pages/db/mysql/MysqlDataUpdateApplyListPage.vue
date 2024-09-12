<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-select
        v-model:value="applyStatus"
        @change="selectApplyStatus"
        style="margin-right:6px;width:180px"
      >
        <a-select-option :value="1">{{t("mysqlDataUpdateApply.pendingStatus")}}</a-select-option>
        <a-select-option :value="2">{{t("mysqlDataUpdateApply.agreeStatus")}}</a-select-option>
        <a-select-option :value="3">{{t("mysqlDataUpdateApply.disagreeStatus")}}</a-select-option>
        <a-select-option :value="4">{{t("mysqlDataUpdateApply.canceledStatus")}}</a-select-option>
        <a-select-option :value="5">{{t("mysqlDataUpdateApply.askToExecuteStatus")}}</a-select-option>
        <a-select-option :value="6">{{t("mysqlDataUpdateApply.executedStatus")}}</a-select-option>
      </a-select>
      <a-button
        type="primary"
        @click="gotoApplyPage"
        :icon="h(PlusOutlined)"
      >{{t("mysqlDataUpdateApply.title")}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" :scroll="{x:1800}">
      <template #bodyCell="{dataIndex, dataItem}">
        <div v-if="dataIndex === 'account'" class="flex-center">
          <ZAvatar
            :url="dataItem.account?.avatarUrl"
            :name="dataItem.account?.name"
            :showName="true"
          />
        </div>
        <div v-else-if="dataIndex === 'auditor'" class="flex-center">
          <ZAvatar
            :url="dataItem.auditor?.avatarUrl"
            :name="dataItem.auditor?.name"
            :showName="true"
          />
        </div>
        <div v-else-if="dataIndex === 'executor'" class="flex-center">
          <ZAvatar
            :url="dataItem.executor?.avatarUrl"
            :name="dataItem.executor?.name"
            :showName="true"
          />
        </div>
        <StatusTag v-else-if="dataIndex === 'applyStatus'" :status="dataItem[dataIndex]" />
        <span
          v-else-if="dataIndex === 'executeImmediatelyAfterApproval'"
        >{{dataItem[dataIndex]?t("mysqlDataUpdateApply.yes"): t("mysqlDataUpdateApply.no")}}</span>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li
                  @click="cancelApply(dataItem)"
                  v-if="applyStatus === 1 || applyStatus === 2 || applyStatus === 5"
                >
                  <CloseOutlined />
                  <span style="margin-left:4px">{{t("mysqlDataUpdateApply.cancel")}}</span>
                </li>
                <li @click="viewExplain(dataItem)" v-if="dataItem.isUnExecuted">
                  <EyeOutlined />
                  <span style="margin-left:4px">{{t("mysqlDataUpdateApply.viewExplain")}}</span>
                </li>
                <template v-if="applyStatus === 2">
                  <li @click="askToExecuteApply(dataItem)">
                    <CloudUploadOutlined />
                    <span style="margin-left:4px">{{t("mysqlDataUpdateApply.askToExecute")}}</span>
                  </li>
                </template>
                <li @click="viewSql(dataItem)">
                  <EyeOutlined />
                  <span style="margin-left:4px">{{t("mysqlDataUpdateApply.viewSql")}}</span>
                </li>
                <li @click="viewLog(dataItem)" v-if="applyStatus === 6">
                  <EyeOutlined />
                  <span style="margin-left:4px">{{t("mysqlDataUpdateApply.viewLog")}}</span>
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
    <a-modal v-model:open="sqlModal.open" :title="t('mysqlDataUpdateApply.viewSql')" :footer="null">
      <Codemirror
        v-model="sqlModal.sql"
        style="height:280px;width:100%"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
    <a-modal
      v-model:open="explainModal.open"
      :title="t('mysqlDataUpdateApply.viewExplain')"
      :footer="null"
      :width="800"
    >
      <Codemirror
        v-model="explainModal.content"
        style="height:280px;width:100%"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
    <a-modal v-model:open="logModal.open" :title="t('mysqlDataUpdateApply.viewLog')" :footer="null">
      <Codemirror
        v-model="logModal.content"
        style="height:280px;width:100%"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
  </div>
</template>
<script setup>
/*
  数据更新申请列表页
*/
import ZAvatar from "@/components/user/ZAvatar";
import ZTable from "@/components/common/ZTable";
import StatusTag from "@/components/db/MysqlDataUpdateApplyStatutsTag";
import {
  PlusOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  CloseOutlined,
  EyeOutlined,
  CloudUploadOutlined
} from "@ant-design/icons-vue";
import {
  listDataUpdateApplyByOperatorRequest,
  cancelDataUpdateApplyRequest,
  explainDataUpdateApplyRequest,
  askToExecuteDataUpdateApplyRequest
} from "@/api/db/mysqlApi";
import { ref, h, createVNode, reactive } from "vue";
import { useRouter } from "vue-router";
import { Modal, message } from "ant-design-vue";
import { Codemirror } from "vue-codemirror";
import { sql } from "@codemirror/lang-sql";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const extensions = [sql()];
// sql modal
const sqlModal = reactive({
  open: false,
  sql: ""
});
// 执行计划modal
const explainModal = reactive({
  open: false,
  content: ""
});
// 日志modal
const logModal = reactive({
  open: false,
  content: ""
});
const router = useRouter();
const dataSource = ref([]);
// 审批状态
const applyStatus = ref(1);
// 分页数据
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
const columns = ref([
  {
    i18nTitle: "mysqlDataUpdateApply.dbName",
    dataIndex: "dbName",
    key: "dbName"
  },
  {
    i18nTitle: "mysqlDataUpdateApply.accessBase",
    dataIndex: "accessBase",
    key: "accessBase"
  },
  {
    i18nTitle: "mysqlDataUpdateApply.applyStatus",
    dataIndex: "applyStatus",
    key: "applyStatus"
  },
  {
    i18nTitle: "mysqlDataUpdateApply.applyReason",
    dataIndex: "applyReason",
    key: "applyReason",
    width: 160
  },
  {
    i18nTitle: "mysqlDataUpdateApply.executeImmediatelyAfterApprovalCol",
    dataIndex: "executeImmediatelyAfterApproval",
    key: "executeImmediatelyAfterApproval"
  },
  {
    i18nTitle: "mysqlDataUpdateApply.created",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "mysqlDataUpdateApply.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
]);
// 选择状态
const selectApplyStatus = () => {
  switch (applyStatus.value) {
    case 1:
      columns.value = [
        {
          i18nTitle: "mysqlDataUpdateApply.dbName",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.accessBase",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.executeImmediatelyAfterApprovalCol",
          dataIndex: "executeImmediatelyAfterApproval",
          key: "executeImmediatelyAfterApproval"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.created",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.operation",
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
          i18nTitle: "mysqlDataUpdateApply.dbName",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.accessBase",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.executeImmediatelyAfterApprovalCol",
          dataIndex: "executeImmediatelyAfterApproval",
          key: "executeImmediatelyAfterApproval"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.auditor",
          dataIndex: "auditor",
          key: "auditor",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.created",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.auditTime",
          dataIndex: "updated",
          key: "updated"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.operation",
          dataIndex: "operation",
          key: "operation",
          width: 130,
          fixed: "right"
        }
      ];
      break;
    case 3:
      columns.value = [
        {
          i18nTitle: "mysqlDataUpdateApply.dbName",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.accessBase",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.executeImmediatelyAfterApprovalCol",
          dataIndex: "executeImmediatelyAfterApproval",
          key: "executeImmediatelyAfterApproval"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.disagreeReason",
          dataIndex: "disagreeReason",
          key: "disagreeReason",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.auditor",
          dataIndex: "auditor",
          key: "auditor",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.created",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.auditTime",
          dataIndex: "updated",
          key: "updated"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.operation",
          dataIndex: "operation",
          key: "operation",
          width: 130,
          fixed: "right"
        }
      ];
      break;
    case 4:
      columns.value = [
        {
          i18nTitle: "mysqlDataUpdateApply.dbName",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.accessBase",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.executeImmediatelyAfterApprovalCol",
          dataIndex: "executeImmediatelyAfterApproval",
          key: "executeImmediatelyAfterApproval"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.created",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.cancelTime",
          dataIndex: "updated",
          key: "updated"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.operation",
          dataIndex: "operation",
          key: "operation",
          width: 130,
          fixed: "right"
        }
      ];
      break;
    case 5:
      columns.value = [
        {
          i18nTitle: "mysqlDataUpdateApply.dbName",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.accessBase",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.executeImmediatelyAfterApprovalCol",
          dataIndex: "executeImmediatelyAfterApproval",
          key: "executeImmediatelyAfterApproval"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.auditor",
          dataIndex: "auditor",
          key: "auditor",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.created",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyTime",
          dataIndex: "updated",
          key: "updated"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.operation",
          dataIndex: "operation",
          key: "operation",
          width: 130,
          fixed: "right"
        }
      ];
      break;
    case 6:
      columns.value = [
        {
          i18nTitle: "mysqlDataUpdateApply.dbName",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.accessBase",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyStatus",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.executeImmediatelyAfterApprovalCol",
          dataIndex: "executeImmediatelyAfterApproval",
          key: "executeImmediatelyAfterApproval"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.auditor",
          dataIndex: "auditor",
          key: "auditor",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.executor",
          dataIndex: "executor",
          key: "executor",
          width: 160
        },
        {
          i18nTitle: "mysqlDataUpdateApply.created",
          dataIndex: "created",
          key: "created"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.executeTime",
          dataIndex: "updated",
          key: "updated"
        },
        {
          i18nTitle: "mysqlDataUpdateApply.operation",
          dataIndex: "operation",
          key: "operation",
          width: 130,
          fixed: "right"
        }
      ];
      break;
  }
  searchApply();
};
// 跳转申请页面
const gotoApplyPage = () => {
  router.push(`/db/mysqlDataUpdateApply/apply`);
};
// 申请列表
const listApply = () => {
  listDataUpdateApplyByOperatorRequest({
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
// 搜索申请
const searchApply = () => {
  dataPage.current = 1;
  listApply();
};
// 取消申请
const cancelApply = item => {
  Modal.confirm({
    title: `${t("mysqlDataUpdateApply.confirmCancel")} ${item.dbName}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      cancelDataUpdateApplyRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchApply();
      });
    },
    onCancel() {}
  });
};
// 请求执行申请
const askToExecuteApply = item => {
  Modal.confirm({
    title: `${t("mysqlDataUpdateApply.confirmAskToExecute")}  ${item.dbName}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      askToExecuteDataUpdateApplyRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchApply();
      });
    },
    onCancel() {}
  });
};
// 查看sql
const viewSql = item => {
  sqlModal.open = true;
  sqlModal.sql = item.updateCmd;
};
// 查看执行计划
const viewExplain = item => {
  explainDataUpdateApplyRequest(item.id).then(res => {
    explainModal.open = true;
    explainModal.content = res.data;
  });
};
// 查看日志
const viewLog = item => {
  logModal.open = true;
  logModal.content = item.executeLog;
};

listApply();
</script>
<style scoped>
</style>