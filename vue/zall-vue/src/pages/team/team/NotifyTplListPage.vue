<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchNameKey"
        style="width:240px;margin-right:6px"
        :placeholder="t('notifyTpl.searchName')"
        @pressEnter="searchTpl"
      >
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('notifyTpl.createTpl')}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'avatarUrl'">
          <a-image :width="20" :height="20" :src="dataItem[dataIndex]" :fallback="fallbackAvatar" />
        </template>
        <template v-else-if="dataIndex === 'notifyType'">
          <span v-if="dataItem[dataIndex] === 'wework'">{{t('notifyTpl.wework')}}</span>
          <span v-else-if="dataItem[dataIndex] === 'feishu'">{{t('notifyTpl.feishu')}}</span>
        </template>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteTpl(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="changeApiKey(dataItem)">
                  <KeyOutlined />
                  <span style="margin-left:4px">{{t('notifyTpl.changeApiKey')}}</span>
                </li>
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">{{t('notifyTpl.updateTpl')}}</span>
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
      @change="()=>listTpl()"
    />
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  KeyOutlined,
  DeleteOutlined
} from "@ant-design/icons-vue";
import {
  deleteNotifyTplRequest,
  listNotifyTplRequest,
  changeNotifyTplApiKeyRequest
} from "@/api/team/notifyApi";
import { ref, h, reactive, createVNode } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useNotifyTplStore } from "@/pinia/notifyTplStore";
import { Modal, message } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const notifyTplStore = useNotifyTplStore();
// 分页数据
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
// 搜索关键词
const searchNameKey = ref("");
const route = useRoute();
const router = useRouter();
// 数据
const dataSource = ref([]);
// 数据项
const columns = [
  {
    i18nTitle: "notifyTpl.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "notifyTpl.notifyType",
    dataIndex: "notifyType",
    key: "notifyType"
  },
  {
    i18nTitle: "notifyTpl.apiKey",
    dataIndex: "apiKey",
    key: "apiKey"
  },
  {
    i18nTitle: "notifyTpl.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 跳转创建模板界面
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/notifyTpl/create`);
};
// 跳转编辑模板界面
const gotoUpdatePage = item => {
  notifyTplStore.id = item.origin.id;
  notifyTplStore.name = item.origin.name;
  notifyTplStore.url = item.origin.notifyCfg.url;
  notifyTplStore.notifyType = item.origin.notifyCfg.notifyType;
  notifyTplStore.template = item.origin.notifyCfg.template;
  notifyTplStore.feishuSignKey = item.origin.notifyCfg.feishuSignKey;
  router.push(`/team/${route.params.teamId}/notifyTpl/${item.id}/update`);
};
// 获取模板列表
const listTpl = () => {
  listNotifyTplRequest({
    teamId: route.params.teamId,
    pageNum: dataPage.current,
    name: searchNameKey.value
  }).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        id: item.id,
        name: item.name,
        notifyType: item.notifyCfg?.notifyType,
        apiKey: item.apiKey,
        // 原始数据
        origin: item
      };
    });
    dataPage.totalCount = res.totalCount;
  });
};
// 变更api key
const changeApiKey = item => {
  Modal.confirm({
    title: `${t("notifyTpl.confirmChangeApiKey")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      changeNotifyTplApiKeyRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listTpl();
      });
    }
  });
};
// 删除模板
const deleteTpl = item => {
  Modal.confirm({
    title: `${t("notifyTpl.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteNotifyTplRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchTpl();
      });
    }
  });
};
// 搜索模板
const searchTpl = () => {
  dataPage.current = 1;
  listTpl();
};
listTpl();
</script>
<style scoped>
</style>