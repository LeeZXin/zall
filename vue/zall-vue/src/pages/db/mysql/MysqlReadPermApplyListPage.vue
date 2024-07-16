<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">申请Mysql读权限</a-button>
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
                <li @click="cancelApply(dataItem)">
                  <CloseOutlined />
                  <span style="margin-left:4px">撤销申请</span>
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
  </div>
</template>
<script setup>
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
  cancelReadPermRequest
} from "@/api/db/mysqlApi";
import { ref, h, createVNode } from "vue";
import { useRouter } from "vue-router";
import { Modal, message } from "ant-design-vue";
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

const gotoCreatePage = () => {
  router.push(`/db/mysqlReadPermApply/apply`);
};

const listApply = () => {
  listReadPermApplyByOperatorRequest({
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
      cancelReadPermRequest(item.id).then(() => {
        message.success("撤销成功");
        currPage.value = 1;
        listApply();
      });
    },
    onCancel() {}
  });
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