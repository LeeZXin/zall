<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px" class="flex-center">
      <a-input
        v-model:value="searchKey"
        style="width:240px;margin-right:6px"
        @pressEnter="searchAccessToken"
      >
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
      <a-button
        type="primary"
        @click="listAccessToken"
        :icon="h(ReloadOutlined)"
        style="margin-right:6px"
      >{{t('feishuAccessToken.refreshList')}}</a-button>
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('feishuAccessToken.createTask')}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon" @click="deleteAccessToken(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">{{t('feishuAccessToken.updateTask')}}</span>
                </li>
                <li @click="changeApiKey(dataItem)">
                  <ReloadOutlined />
                  <span style="margin-left:4px">{{t('feishuAccessToken.changeApiKey')}}</span>
                </li>
                <li @click="refreshAccessToken(dataItem)">
                  <ReloadOutlined />
                  <span style="margin-left:4px">{{t('feishuAccessToken.refreshToken')}}</span>
                </li>
                <li @click="viewToken(dataItem)">
                  <EyeOutlined />
                  <span style="margin-left:4px">{{t('feishuAccessToken.viewToken')}}</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </template>
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
      @change="()=>listAccessToken()"
    />
  </div>
  <a-modal v-model:open="tokenModal.open" :title="tokenModal.title" :footer="null">
    <div style="font-size:14px">
      <div>access token</div>
      <div style="padding:10px">
        <span style="word-break:break-all;">{{tokenModal.token}}</span>
        <div class="copy-btn" @click="copyToken">
          <CopyOutlined />
          <span>{{t('feishuAccessToken.copy')}}</span>
        </div>
      </div>
      <div>tenant token</div>
      <div style="padding:10px">
        <span style="word-break:break-all;">{{tokenModal.tenantToken}}</span>
        <div class="copy-btn" @click="copyTenantToken">
          <CopyOutlined />
          <span>{{t('feishuAccessToken.copy')}}</span>
        </div>
      </div>
      <div>{{t('feishuAccessToken.expiredTime')}}</div>
      <div style="word-break:break-all;padding:10px">{{tokenModal.expired}}</div>
    </div>
  </a-modal>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { ref, h, reactive, createVNode } from "vue";
import {
  SearchOutlined,
  PlusOutlined,
  DeleteOutlined,
  EllipsisOutlined,
  ReloadOutlined,
  EyeOutlined,
  CopyOutlined,
  ExclamationCircleOutlined,
  EditOutlined
} from "@ant-design/icons-vue";
import {
  listAccessTokenRequest,
  deleteAccessTokenRequest,
  refreshAccessTokenRequest,
  changeAccessTokenApiKeyRequest
} from "@/api/team/feishuApi";
import { useRoute, useRouter } from "vue-router";
import { message, Modal } from "ant-design-vue";
import { useFeishuAccessTokenStore } from "@/pinia/feishuAccessTokenStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const fsatStore = useFeishuAccessTokenStore();
const router = useRouter();
const route = useRoute();
// 搜索关键词
const searchKey = ref("");
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/feishuAccessToken/create`);
};
// 页面跳转
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
// 数据项
const columns = [
  {
    i18nTitle: "feishuAccessToken.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "feishuAccessToken.appId",
    dataIndex: "appId",
    key: "appId"
  },
  {
    i18nTitle: "feishuAccessToken.secret",
    dataIndex: "secret",
    key: "secret"
  },
  {
    i18nTitle: "feishuAccessToken.apiKey",
    dataIndex: "apiKey",
    key: "apiKey"
  },
  {
    i18nTitle: "feishuAccessToken.creator",
    dataIndex: "creator",
    key: "creator"
  },
  {
    i18nTitle: "feishuAccessToken.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
const tokenModal = reactive({
  open: false,
  title: "",
  token: "",
  tenantToken: "",
  expired: ""
});
// 数据列表
const dataSource = ref([]);
// 搜索列表
const listAccessToken = () => {
  listAccessTokenRequest({
    pageNum: dataPage.current,
    key: searchKey.value,
    teamId: route.params.teamId
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
// 搜索token
const searchAccessToken = () => {
  dataPage.current = 1;
  listAccessToken();
};
// 展示token
const viewToken = item => {
  tokenModal.open = true;
  tokenModal.title = item.name;
  tokenModal.token = item.token;
  tokenModal.tenantToken = item.tenantToken;
  tokenModal.expired = item.expired;
};
// 复制token
const copyToken = () => {
  message.success(t("operationSuccess"));
  window.navigator.clipboard.writeText(tokenModal.token);
};
// 复制tenant token
const copyTenantToken = () => {
  message.success(t("operationSuccess"));
  window.navigator.clipboard.writeText(tokenModal.tenantToken);
};
// 删除token
const deleteAccessToken = item => {
  Modal.confirm({
    title: `${t("feishuAccessToken.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteAccessTokenRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchAccessToken();
      });
    },
    onCancel() {}
  });
};
// 刷新token
const refreshAccessToken = item => {
  Modal.confirm({
    title: `${t("feishuAccessToken.confirmRefresh")} ${
      item.name
    } access token?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      refreshAccessTokenRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listAccessToken();
      });
    },
    onCancel() {}
  });
};
// 变更api key
const changeApiKey = item => {
  Modal.confirm({
    title: `${t("feishuAccessToken.confirmChange")} ${item.name} ${t(
      "feishuAccessToken.apiKey"
    )}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      changeAccessTokenApiKeyRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listAccessToken();
      });
    },
    onCancel() {}
  });
};
// 跳转变更页面
const gotoUpdatePage = item => {
  fsatStore.id = item.id;
  fsatStore.name = item.name;
  fsatStore.appId = item.appId;
  fsatStore.secret = item.secret;
  router.push(
    `/team/${route.params.teamId}/feishuAccessToken/${item.id}/update`
  );
};
listAccessToken();
</script>
<style scoped>
.copy-btn {
  display: inline-block;
}
.copy-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>