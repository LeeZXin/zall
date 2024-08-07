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
      <a-select
        style="width: 300px;margin-left:6px"
        v-model:value="selectedUser"
        :options="userList"
        show-search
        :filter-option="filterUserListOption"
        @change="selectUserChange"
      />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteReadPerm(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete</span>
              </template>
              <DeleteOutlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="getApply(dataItem)">
                  <EyeOutlined />
                  <span style="margin-left:4px">查看审批单</span>
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
      @change="()=>listPerm()"
    />
    <a-modal v-model:open="apply.open" title="审批单" :footer="null">
      <ul class="apply-ul">
        <li>
          <div class="item-name">数据库名称</div>
          <div class="item-value">{{apply.dbName}}</div>
        </li>
        <li>
          <div class="item-name">申请库</div>
          <div class="item-value">{{apply.accessBase}}</div>
        </li>
        <li>
          <div class="item-name">申请表</div>
          <div class="item-value">{{apply.accessTables}}</div>
        </li>
        <li>
          <div class="item-name">申请人</div>
          <div class="item-value">{{apply.account}}</div>
        </li>
        <li>
          <div class="item-name">审批人</div>
          <div class="item-value">{{apply.auditor}}</div>
        </li>
        <li>
          <div class="item-name">申请原因</div>
          <div class="item-value">{{apply.applyReason}}</div>
        </li>
        <li>
          <div class="item-name">申请时间</div>
          <div class="item-value">{{apply.created}}</div>
        </li>
        <li>
          <div class="item-name">审批时间</div>
          <div class="item-value">{{apply.updated}}</div>
        </li>
      </ul>
    </a-modal>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import {
  listReadPermByDbaRequest,
  getAllMysqlDbRequest,
  getReadPermApplyRequest,
  deleteReadPermRequest
} from "@/api/db/mysqlApi";
import { listAllUserRequest } from "@/api/user/userApi";
import {
  EyeOutlined,
  EllipsisOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import { ref, reactive, createVNode } from "vue";
import { Modal, message } from "ant-design-vue";
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const selectedDbId = ref(0);
const selectedUser = ref("");
const dbList = ref([
  {
    value: 0,
    label: "所有数据库"
  }
]);
const userList = ref([
  {
    value: "",
    label: "所有人"
  }
]);
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
  updated: "",
  account: ""
});
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
    dataIndex: "accessTable",
    key: "accessTable"
  },
  {
    title: "申请人",
    dataIndex: "account",
    key: "account"
  },
  {
    title: "生效时间",
    dataIndex: "created",
    key: "created"
  },
  {
    title: "过期时间",
    dataIndex: "expired",
    key: "expired"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);

const listPerm = () => {
  listReadPermByDbaRequest({
    dbId: selectedDbId.value,
    pageNum: currPage.value,
    account: selectedUser.value
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

const filterDbListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};

const filterUserListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
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

const listAllUser = () => {
  listAllUserRequest().then(res => {
    userList.value = userList.value.concat(
      res.data.map(item => {
        return {
          value: item.account,
          label: `${item.account}(${item.name})`
        };
      })
    );
  });
};

const selectDbIdChange = () => {
  currPage.value = 1;
  listPerm();
};

const selectUserChange = () => {
  currPage.value = 1;
  listPerm();
};

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
    apply.account = data.account;
    apply.open = true;
  });
};

const deleteReadPerm = item => {
  Modal.confirm({
    title: `你确定要删除${item.account}的权限吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteReadPermRequest(item.id).then(() => {
        message.success("操作成功");
        currPage.value = 1;
        listPerm();
      });
    },
    onCancel() {}
  });
};

getAllDb();
listAllUser();
listPerm();
</script>
<style scoped>
.check-btn {
  font-size: 14px;
}
.check-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
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