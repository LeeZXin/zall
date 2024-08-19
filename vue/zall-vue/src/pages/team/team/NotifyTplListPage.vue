<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchNameKey"
        placeholder="搜索名称"
        style="width:240px;margin-right:6px"
        @pressEnter="searchTpl"
      >
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">创建通知模板</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'avatarUrl'">
          <a-image :width="20" :height="20" :src="dataItem[dataIndex]" :fallback="fallbackAvatar" />
        </template>
        <template v-else-if="dataIndex === 'notifyType'">
          <span v-if="dataItem[dataIndex] === 'wework'">企业微信</span>
          <span v-else-if="dataItem[dataIndex] === 'feishu'">飞书</span>
        </template>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteTpl(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>删除模板</span>
              </template>
              <DeleteOutlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="hideApiKey(dataItem)" v-if="dataItem['showApiKey']">
                  <EyeInvisibleOutlined />
                  <span style="margin-left:4px">隐藏api密钥</span>
                </li>
                <li @click="showApiKey(dataItem)" v-else>
                  <EyeOutlined />
                  <span style="margin-left:4px">查看api密钥</span>
                </li>
                <li @click="changeApiKey(dataItem)">
                  <KeyOutlined />
                  <span style="margin-left:4px">变更api密钥</span>
                </li>
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">编辑模板</span>
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
  DeleteOutlined,
  EyeOutlined,
  EyeInvisibleOutlined
} from "@ant-design/icons-vue";
import {
  deleteNotifyTplRequest,
  listNotifyTplRequest,
  changeNotifyTplApiKeyRequest
} from "@/api/notify/notifyApi";
import { ref, h, reactive, createVNode } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useNotifyTplStore } from "@/pinia/notifyTplStore";
import { Modal, message } from "ant-design-vue";
const notifyTplStore = useNotifyTplStore();
// 密钥展示****
const sensitiveStr = "********************************";
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
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "类型",
    dataIndex: "notifyType",
    key: "notifyType"
  },
  {
    title: "Api密钥",
    dataIndex: "sensitiveApiKey",
    key: "sensitiveApiKey"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
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
        sensitiveApiKey: sensitiveStr,
        // 是否展示api密钥
        showApiKey: false,
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
    title: `是否要变更${item.name}的api密钥吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      changeNotifyTplApiKeyRequest(item.id).then(() => {
        message.success("操作成功");
        listTpl();
      });
    }
  });
};
// 查看api密钥
const showApiKey = item => {
  item.showApiKey = true;
  item.sensitiveApiKey = item.apiKey;
};
// 隐藏密钥
const hideApiKey = item => {
  item.showApiKey = false;
  item.sensitiveApiKey = sensitiveStr;
};
// 删除模板
const deleteTpl = item => {
  Modal.confirm({
    title: `是否要删除${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteNotifyTplRequest(item.id).then(() => {
        message.success("操作成功");
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