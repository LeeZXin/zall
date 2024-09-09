<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchName"
        :placeholder="t('mysqlSource.searchPlaceholder')"
        style="width:240px;margin-right:10px"
        @pressEnter="()=>searchDb()"
      >
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('mysqlSource.createSource')}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteDb(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">{{t('mysqlSource.updateSource')}}</span>
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
import { ref, h, createVNode, reactive } from "vue";
import { useRouter } from "vue-router";
import { Modal, message } from "ant-design-vue";
import { useMysqldbStore } from "@/pinia/mysqldbStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const dbStore = useMysqldbStore();
const router = useRouter();
// 搜索key
const searchName = ref("");
// 数据列表
const dataSource = ref([]);
// 分页数据
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
// 数据项
const columns = [
  {
    i18nTitle: "mysqlSource.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "mysqlSource.writeHost",
    dataIndex: "writeHost",
    key: "writeHost"
  },
  {
    i18nTitle: "mysqlSource.readHost",
    dataIndex: "readHost",
    key: "readHost"
  },
  {
    i18nTitle: "mysqlSource.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(`/db/mysqlDb/create`);
};
// 跳转编辑页面
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
// 搜索数据库
const searchDb = () => {
  dataPage.current = 1;
  listDb();
};
// 获取数据库列表
const listDb = () => {
  listMysqlDbRequest({
    pageNum: dataPage.current,
    name: searchName.value
  }).then(res => {
    dataPage.totalCount = res.totalCount;
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
// 删除数据库
const deleteDb = item => {
  Modal.confirm({
    title: `${t("mysqlSource.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteMysqlDbRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchDb();
      });
    },
    onCancel() {}
  });
};

listDb();
</script>
<style scoped>
</style>