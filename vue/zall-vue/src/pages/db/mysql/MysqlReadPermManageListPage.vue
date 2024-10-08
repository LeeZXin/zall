<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-select style="width: 180px" v-model:value="selectedDbId" @change="searchPerm">
        <a-select-option :value="0">{{t("mysqlReadPermApply.allDatabases")}}</a-select-option>
        <a-select-option
          :value="item.value"
          v-for="item in dbList"
          v-bind:key="item.value"
        >{{item.label}}</a-select-option>
      </a-select>
      <a-select
        style="width: 180px;margin-left:6px"
        v-model:value="selectedUser"
        :options="userList"
        show-search
        :filter-option="filterUserListOption"
        @change="searchPerm"
      />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <div v-if="dataIndex === 'account'" class="flex-center">
          <ZAvatar
            :url="dataItem.account?.avatarUrl"
            :name="dataItem.account?.name"
            :account="dataItem.account?.account"
            :showName="true"
          />
        </div>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteReadPerm(dataItem)">
            <DeleteOutlined />
          </div>
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
    <a-modal
      v-model:open="applyModal.open"
      :title="t('mysqlReadPermApply.approvalForm')"
      :footer="null"
    >
      <ul class="apply-ul">
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.dbName')}}</div>
          <div class="item-value">{{applyModal.dbName}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.accessBase')}}</div>
          <div class="item-value">{{applyModal.accessBase}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.accessTables')}}</div>
          <div class="item-value">{{applyModal.accessTables}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.account')}}</div>
          <div class="item-value">
            <ZAvatar
              :url="applyModal.account?.avatarUrl"
              :name="applyModal.account?.name"
              :account="applyModal.account?.account"
              :showName="true"
            />
          </div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.auditor')}}</div>
          <div class="item-value">
            <ZAvatar
              :url="applyModal.auditor?.avatarUrl"
              :name="applyModal.auditor?.name"
              :account="applyModal.auditor?.account"
              :showName="true"
            />
          </div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.applyReason')}}</div>
          <div class="item-value">{{applyModal.applyReason}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.applyTime')}}</div>
          <div class="item-value">{{applyModal.created}}</div>
        </li>
        <li>
          <div class="item-name">{{t('mysqlReadPermApply.auditTime')}}</div>
          <div class="item-value">{{applyModal.updated}}</div>
        </li>
      </ul>
    </a-modal>
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const dataSource = ref([]);
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
const selectedDbId = ref(0);
const selectedUser = ref("");
// 数据库列表
const dbList = ref([]);
// 用户列表
const userList = ref([
  {
    value: "",
    label: "ALL"
  }
]);
// 审批单
const applyModal = reactive({
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
    i18nTitle: "mysqlReadPermApply.account",
    dataIndex: "account",
    key: "account"
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
    key: "operation",
    width: 130,
    fixed: "right"
  }
]);
// 权限列表
const listPerm = () => {
  listReadPermByDbaRequest({
    dbId: selectedDbId.value,
    pageNum: dataPage.current,
    account: selectedUser.value
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
// 用户下拉框搜索
const filterUserListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
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
// 所有用户
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
// 搜索权限
const searchPerm = () => {
  dataPage.current = 1;
  listPerm();
};
// 获取审批单
const getApply = item => {
  getReadPermApplyRequest(item.applyId).then(res => {
    let data = res.data;
    applyModal.dbName = data.dbName;
    applyModal.accessBase = data.accessBase;
    applyModal.accessTables = data.accessTables;
    applyModal.expireDay = data.expireDay;
    applyModal.applyStatus = data.applyStatus;
    applyModal.auditor = data.auditor;
    applyModal.applyReason = data.applyReason;
    applyModal.created = data.created;
    applyModal.updated = data.updated;
    applyModal.account = data.account;
    applyModal.open = true;
  });
};
// 删除权限
const deleteReadPerm = item => {
  Modal.confirm({
    title: `${t("mysqlReadPermApply.confirmDelete")} ${item.account} ${
      item.dbName
    } ${item.accessBase} ${item.accessTable}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteReadPermRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchPerm();
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
  margin-top: 16px;
}
.item-name {
  font-size: 14px;
  margin-bottom: 4px;
}
.item-value {
  font-size: 14px;
  line-height: 18px;
  padding-left: 20px;
  min-height: 18px;
}
</style>