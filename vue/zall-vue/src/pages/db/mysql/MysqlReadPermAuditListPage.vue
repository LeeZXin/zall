<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-select
        style="width: 300px"
        v-model:value="selectedDbId"
        :options="dbList"
        show-search
        :filter-option="filterDbListOption"
        @change="selectDbIdChange"
      />
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
      </a-radio-group>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <StatusTag v-if="dataIndex === 'applyStatus'" :status="dataItem[dataIndex]" />
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="agreeApply(dataItem)">
                  <CheckOutlined />
                  <span style="margin-left:4px">同意</span>
                </li>
                <li @click="showDisagreeModal(dataItem)">
                  <CloseOutlined />
                  <span style="margin-left:4px">不同意</span>
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
    <a-modal v-model:open="disagreeObj.modalOpen" title="填写不同意原因" @ok="disagreeApply">
      <a-textarea
        style="width:100%"
        v-model:value="disagreeObj.reason"
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
import { dbApplyReasonRegexp } from "../../../utils/regexp";
const dataSource = ref([]);
const applyStatus = ref(1);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const selectedDbId = ref(0);
const disagreeObj = reactive({
  id: 0,
  modalOpen: false,
  reason: ""
});
const dbList = ref([
  {
    value: 0,
    label: "所有数据库"
  }
]);
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
    title: "申请表",
    dataIndex: "accessTables",
    key: "accessTables"
  },
  {
    title: "时效(天)",
    dataIndex: "expireDay",
    key: "expireDay"
  },
  {
    title: "状态",
    dataIndex: "applyStatus",
    key: "applyStatus"
  },
  {
    title: "申请人",
    dataIndex: "account",
    key: "account"
  },
  {
    title: "申请原因",
    dataIndex: "applyReason",
    key: "applyReason"
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

const selectDbIdChange = () => {
  currPage.value = 1;
  listApply();
};

const filterDbListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};

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
          title: "申请表",
          dataIndex: "accessTables",
          key: "accessTables"
        },
        {
          title: "时效(天)",
          dataIndex: "expireDay",
          key: "expireDay"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "申请人",
          dataIndex: "account",
          key: "account"
        },
        {
          title: "申请原因",
          dataIndex: "applyReason",
          key: "applyReason"
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
          title: "申请表",
          dataIndex: "accessTables",
          key: "accessTables"
        },
        {
          title: "时效(天)",
          dataIndex: "expireDay",
          key: "expireDay"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "申请人",
          dataIndex: "account",
          key: "account"
        },
        {
          title: "申请原因",
          dataIndex: "applyReason",
          key: "applyReason"
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
          title: "申请表",
          dataIndex: "accessTables",
          key: "accessTables"
        },
        {
          title: "时效(天)",
          dataIndex: "expireDay",
          key: "expireDay"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "申请人",
          dataIndex: "account",
          key: "account"
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
          title: "申请表",
          dataIndex: "accessTables",
          key: "accessTables"
        },
        {
          title: "时效(天)",
          dataIndex: "expireDay",
          key: "expireDay"
        },
        {
          title: "状态",
          dataIndex: "applyStatus",
          key: "applyStatus"
        },
        {
          title: "申请人",
          dataIndex: "account",
          key: "account"
        },
        {
          title: "申请原因",
          dataIndex: "applyReason",
          key: "applyReason"
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
        }
      ];
      break;
  }
  listApply();
};

const listApply = () => {
  listReadPermApplyByDbaRequest({
    dbId: selectedDbId.value,
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

const agreeApply = item => {
  Modal.confirm({
    title: `你确定要同意${item.account}的申请吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      agreeReadPermApplyRequest(item.id).then(() => {
        message.success("操作成功");
        currPage.value = 1;
        listApply();
      });
    },
    onCancel() {}
  });
};

const showDisagreeModal = item => {
  disagreeObj.id = item.id;
  disagreeObj.modalOpen = true;
  disagreeObj.reason = "";
};

const disagreeApply = () => {
  if (!dbApplyReasonRegexp.test(disagreeObj.reason)) {
    message.warn("原因格式错误");
    return;
  }
  disagreeReadPermApplyRequest({
    applyId: disagreeObj.id,
    disagreeReason: disagreeObj.reason
  }).then(() => {
    message.success("操作成功");
    currPage.value = 1;
    listApply();
    disagreeObj.modalOpen = false;
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
.check-btn {
  font-size: 14px;
}
.check-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>