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
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('mysqlReadPermApply.applyReadPerm')}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" :scroll="{x:1800}">
      <template #bodyCell="{dataIndex, dataItem}">
        <div v-if="dataIndex === 'auditor'" class="flex-center">
          <ZAvatar
            :url="dataItem.auditor?.avatarUrl"
            :name="dataItem.auditor?.name"
            :account="dataItem.auditor?.account"
            :showName="true"
          />
        </div>
        <StatusTag v-else-if="dataIndex === 'applyStatus'" :status="dataItem[dataIndex]" />
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="cancelApply(dataItem)">
                  <CloseOutlined />
                  <span style="margin-left:4px">{{t('mysqlReadPermApply.cancelApply')}}</span>
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
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import ZTable from "@/components/common/ZTable";
import StatusTag from "@/components/db/MysqlReadPermApplyStatutsTag";
import {
  PlusOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  CloseOutlined
} from "@ant-design/icons-vue";
import {
  listReadPermApplyByOperatorRequest,
  cancelReadPermApplyRequest
} from "@/api/db/mysqlApi";
import { ref, h, createVNode, reactive } from "vue";
import { useRouter } from "vue-router";
import { Modal, message } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const dataSource = ref([]);
// 审批状态
const applyStatus = ref(1);
// 分页
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
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
    dataIndex: "accessTables",
    key: "accessTables",
    width: 130
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
// 选择审批状态
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
          key: "accessTables",
          width: 130
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
          key: "accessTables",
          width: 130
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
          i18nTitle: "mysqlReadPermApply.applyReason",
          dataIndex: "applyReason",
          key: "applyReason",
          width: 160
        },
        {
          i18nTitle: "mysqlReadPermApply.auditor",
          dataIndex: "auditor",
          key: "auditor",
          width: 160
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
          key: "accessTables",
          width: 130
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
          key: "auditor",
          width: 160
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
          key: "accessTables",
          width: 130
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
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(`/db/mysqlReadPermApply/apply`);
};
// 获取列表
const listApply = () => {
  listReadPermApplyByOperatorRequest({
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
// 撤销申请
const cancelApply = item => {
  Modal.confirm({
    title: `${t("mysqlReadPermApply.confirmCancel")} ${item.dbName}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      cancelReadPermApplyRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchApply();
      });
    },
    onCancel() {}
  });
};
// 搜索列表
const searchApply = () => {
  dataPage.current = 1;
  listApply();
};

listApply();
</script>
<style scoped>
</style>