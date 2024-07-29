<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">申请数据修改单</a-button>
    </div>
    <div>
      <a-radio-group v-model:value="applyStatus" @change="selectApplyStatus">
        <a-radio-button :value="1">
          <span>等待审批</span>
        </a-radio-button>
        <a-radio-button :value="2">
          <span>同意</span>
        </a-radio-button>
        <a-radio-button :value="3">
          <span>不同意</span>
        </a-radio-button>
        <a-radio-button :value="4">
          <span>已取消</span>
        </a-radio-button>
        <a-radio-button :value="5">
          <span>请求执行</span>
        </a-radio-button>
        <a-radio-button :value="6">
          <span>已执行</span>
        </a-radio-button>v
      </a-radio-group>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <StatusTag v-if="dataIndex === 'applyStatus'" :status="dataItem[dataIndex]" />
        <span v-else-if="dataIndex === 'executeWhenApply'">{{dataItem[dataIndex]?"是": "否"}}</span>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="cancelApply(dataItem)" v-if="applyStatus === 1">
                  <CloseOutlined />
                  <span style="margin-left:4px">撤销申请</span>
                </li>
                <li @click="checkExplain(dataItem)" v-if="dataItem.isUnExecuted">
                  <EyeOutlined />
                  <span style="margin-left:4px">执行计划</span>
                </li>
                <template v-if="applyStatus === 2">
                  <li @click="askToExecuteApply(dataItem)">
                    <CloudUploadOutlined />
                    <span style="margin-left:4px">请求执行</span>
                  </li>
                </template>
                <li @click="checkSql(dataItem)">
                  <EyeOutlined />
                  <span style="margin-left:4px">查看sql</span>
                </li>
                <li @click="checkLog(dataItem)" v-if="applyStatus === 6">
                  <EyeOutlined />
                  <span style="margin-left:4px">查看执行日志</span>
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
      v-model:current="currPage"
      :total="totalCount"
      show-less-items
      :pageSize="pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listApply()"
    />
    <a-modal v-model:open="sqlModal.open" title="sql" :footer="null">
      <Codemirror
        v-model="sqlModal.sql"
        style="height:280px;width:100%"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
    <a-modal v-model:open="explainModal.open" title="执行计划" :footer="null">
      <Codemirror
        v-model="explainModal.content"
        style="height:280px;width:100%"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
    <a-modal v-model:open="logModal.open" title="执行日志" :footer="null">
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
const extensions = [sql()];
const sqlModal = reactive({
  open: false,
  sql: ""
});
const explainModal = reactive({
  open: false,
  content: ""
});
const logModal = reactive({
  open: false,
  content: ""
});
const router = useRouter();
const dataSource = ref([]);
const applyStatus = ref(1);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const columns = ref([
  {
    title: "数据库名称",
    dataIndex: "dbName",
    key: "dbName"
  },
  {
    title: "申请库",
    dataIndex: "accessBase",
    key: "accessBase"
  },
  {
    title: "状态",
    dataIndex: "applyStatus",
    key: "applyStatus"
  },
  {
    title: "申请原因",
    dataIndex: "applyReason",
    key: "applyReason"
  },
  {
    title: "是否立即执行",
    dataIndex: "executeWhenApply",
    key: "executeWhenApply"
  },
  {
    title: "申请时间",
    dataIndex: "created",
    key: "created"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);

const selectApplyStatus = () => {
  switch (applyStatus.value) {
    case 1:
      columns.value = [
        {
          title: "数据库名称",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          title: "申请库",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "申请原因",
          dataIndex: "applyReason",
          key: "applyReason"
        },
        {
          title: "是否立即执行",
          dataIndex: "executeWhenApply",
          key: "executeWhenApply"
        },
        {
          title: "申请时间",
          dataIndex: "created",
          key: "created"
        },
        {
          title: "操作",
          dataIndex: "operation",
          key: "operation"
        }
      ];
      break;
    case 2:
      columns.value = [
        {
          title: "数据库名称",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          title: "申请库",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "申请原因",
          dataIndex: "applyReason",
          key: "applyReason"
        },
        {
          title: "是否立即执行",
          dataIndex: "executeWhenApply",
          key: "executeWhenApply"
        },
        {
          title: "审批人",
          dataIndex: "auditor",
          key: "auditor"
        },
        {
          title: "申请时间",
          dataIndex: "created",
          key: "created"
        },
        {
          title: "审批时间",
          dataIndex: "updated",
          key: "updated"
        },
        {
          title: "操作",
          dataIndex: "operation",
          key: "operation"
        }
      ];
      break;
    case 3:
      columns.value = [
        {
          title: "数据库名称",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          title: "申请库",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "是否立即执行",
          dataIndex: "executeWhenApply",
          key: "executeWhenApply"
        },
        {
          title: "申请原因",
          dataIndex: "applyReason",
          key: "applyReason"
        },
        {
          title: "不同意原因",
          dataIndex: "disagreeReason",
          key: "disagreeReason"
        },
        {
          title: "审批人",
          dataIndex: "auditor",
          key: "auditor"
        },
        {
          title: "申请时间",
          dataIndex: "created",
          key: "created"
        },
        {
          title: "审批时间",
          dataIndex: "updated",
          key: "updated"
        },
        {
          title: "操作",
          dataIndex: "operation",
          key: "operation"
        }
      ];
      break;
    case 4:
      columns.value = [
        {
          title: "数据库名称",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          title: "申请库",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "申请原因",
          dataIndex: "applyReason",
          key: "applyReason"
        },
        {
          title: "是否立即执行",
          dataIndex: "executeWhenApply",
          key: "executeWhenApply"
        },
        {
          title: "申请时间",
          dataIndex: "created",
          key: "created"
        },
        {
          title: "取消时间",
          dataIndex: "updated",
          key: "updated"
        },
        {
          title: "操作",
          dataIndex: "operation",
          key: "operation"
        }
      ];
      break;
    case 5:
      columns.value = [
        {
          title: "数据库名称",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          title: "申请库",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "申请原因",
          dataIndex: "applyReason",
          key: "applyReason"
        },
        {
          title: "是否立即执行",
          dataIndex: "executeWhenApply",
          key: "executeWhenApply"
        },
        {
          title: "审批人",
          dataIndex: "auditor",
          key: "auditor"
        },
        {
          title: "申请时间",
          dataIndex: "created",
          key: "created"
        },
        {
          title: "请求时间",
          dataIndex: "updated",
          key: "updated"
        },
        {
          title: "操作",
          dataIndex: "operation",
          key: "operation"
        }
      ];
      break;
    case 6:
      columns.value = [
        {
          title: "数据库名称",
          dataIndex: "dbName",
          key: "dbName"
        },
        {
          title: "申请库",
          dataIndex: "accessBase",
          key: "accessBase"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "申请原因",
          dataIndex: "applyReason",
          key: "applyReason"
        },
        {
          title: "是否立即执行",
          dataIndex: "executeWhenApply",
          key: "executeWhenApply"
        },
        {
          title: "审批人",
          dataIndex: "auditor",
          key: "auditor"
        },
        {
          title: "申请时间",
          dataIndex: "created",
          key: "created"
        },
        {
          title: "执行时间",
          dataIndex: "updated",
          key: "updated"
        },
        {
          title: "操作",
          dataIndex: "operation",
          key: "operation"
        }
      ];
      break;
  }
  listApply();
};

const gotoCreatePage = () => {
  router.push(`/db/mysqlDataUpdateApply/apply`);
};

const listApply = () => {
  listDataUpdateApplyByOperatorRequest({
    pageNum: currPage.value,
    applyStatus: applyStatus.value
  }).then(res => {
    totalCount.value = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};

const cancelApply = item => {
  Modal.confirm({
    title: `你确定要撤销${item.dbName}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      cancelDataUpdateApplyRequest(item.id).then(() => {
        message.success("撤销成功");
        currPage.value = 1;
        listApply();
      });
    },
    onCancel() {}
  });
};

const askToExecuteApply = item => {
  Modal.confirm({
    title: `你确定要请求执行${item.dbName}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      askToExecuteDataUpdateApplyRequest(item.id).then(() => {
        message.success("操作成功");
        currPage.value = 1;
        listApply();
      });
    },
    onCancel() {}
  });
};

const checkSql = item => {
  sqlModal.open = true;
  sqlModal.sql = item.updateCmd;
};

const checkExplain = item => {
  explainDataUpdateApplyRequest(item.id).then(res => {
    explainModal.open = true;
    explainModal.content = res.data;
  });
};

const checkLog = item => {
  logModal.open = true;
  logModal.content = item.executeLog;
};

listApply();
</script>
<style scoped>
.check-btn {
  font-size: 14px;
}
.check-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>