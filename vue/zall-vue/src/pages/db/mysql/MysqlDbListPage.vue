<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchName"
        placeholder="搜索Mysql数据源"
        style="width:240px;margin-right:10px"
        @pressEnter="()=>searchDb()"
      >
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">创建Mysql数据源</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteDb(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Source</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <edit-outlined />
                  <span style="margin-left:4px">编辑数据源</span>
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
      @change="()=>listDb()"
    />
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import {
  PlusOutlined,
  DeleteOutlined,
  EllipsisOutlined,
  EditOutlined,
  ExclamationCircleOutlined,
  SearchOutlined
} from "@ant-design/icons-vue";
import { listMysqlDbRequest, deleteMysqlDbRequest } from "@/api/db/mysqlApi";
import { ref, h, createVNode } from "vue";
import { useRouter } from "vue-router";
import { Modal, message } from "ant-design-vue";
import { useMysqldbStore } from "@/pinia/mysqldbStore";
const dbStore = useMysqldbStore();
const router = useRouter();
const searchName = ref("");
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const columns = [
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "写节点",
    dataIndex: "writeHost",
    key: "writeHost"
  },
  {
    title: "读节点",
    dataIndex: "readHost",
    key: "readHost"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];

const gotoCreatePage = () => {
  router.push(`/db/mysqlDb/create`);
};

const gotoUpdatePage = item => {
  dbStore.id = item.id;
  dbStore.name = item.name;
  dbStore.writeHost = item.config.writeNode.host;
  dbStore.writeUsername = item.config.writeNode.username;
  dbStore.writePassword = item.config.writeNode.password;
  dbStore.readHost = item.config.readNode.host;
  dbStore.readUsername = item.config.readNode.username;
  dbStore.readPassword = item.config.readNode.password;
  router.push(`/db/mysqlDb/${item.id}/update`);
};

const searchDb = () => {
  currPage.value = 1;
  listDb();
};

const listDb = () => {
  listMysqlDbRequest({
    pageNum: currPage.value,
    name: searchName.value
  }).then(res => {
    totalCount.value = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        id: item.id,
        name: item.name,
        writeHost: item.config.writeNode.host,
        readHost: item.config.readNode.host,
        config: item.config
      };
    });
  });
};

const deleteDb = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteMysqlDbRequest(item.id).then(() => {
        message.success("删除成功");
        searchDb();
      });
    },
    onCancel() {}
  });
};

listDb();
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