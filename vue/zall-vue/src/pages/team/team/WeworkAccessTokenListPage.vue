<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px" class="flex-center">
      <a-input
        v-model:value="searchKey"
        style="width:240px;margin-right:6px"
        @pressEnter="searchAccessToken"
        :placeholder="t('weworkAccessToken.searchNameOrCorpId')"
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
      >{{t('weworkAccessToken.refreshList')}}</a-button>
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('weworkAccessToken.createTask')}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <div v-if="dataIndex === 'creator'" class="flex-center">
          <ZAvatar
            :url="dataItem.creator?.avatarUrl"
            :name="dataItem.creator?.name"
            :showName="true"
          />
        </div>
        <template v-else-if="dataIndex !== 'operation'">
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
                  <span style="margin-left:4px">{{t('weworkAccessToken.updateTask')}}</span>
                </li>
                <li @click="changeApiKey(dataItem)">
                  <ReloadOutlined />
                  <span style="margin-left:4px">{{t('weworkAccessToken.changeApiKey')}}</span>
                </li>
                <li @click="refreshAccessToken(dataItem)">
                  <ReloadOutlined />
                  <span style="margin-left:4px">{{t('weworkAccessToken.refreshToken')}}</span>
                </li>
                <li @click="viewToken(dataItem)">
                  <EyeOutlined />
                  <span style="margin-left:4px">{{t('weworkAccessToken.viewToken')}}</span>
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
          <span>{{t('weworkAccessToken.copy')}}</span>
        </div>
      </div>
      <div>{{t('weworkAccessToken.expiredTime')}}</div>
      <div style="word-break:break-all;padding:10px">{{tokenModal.expired}}</div>
    </div>
  </a-modal>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
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
} from "@/api/team/weworkApi";
import { useRoute, useRouter } from "vue-router";
import { message, Modal } from "ant-design-vue";
import { useWeworkAccessTokenStore } from "@/pinia/weworkAccessTokenStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const wwatStore = useWeworkAccessTokenStore();
const router = useRouter();
const route = useRoute();
// 搜索关键词
const searchKey = ref("");
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/weworkAccessToken/create`);
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
    i18nTitle: "weworkAccessToken.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "weworkAccessToken.corpId",
    dataIndex: "corpId",
    key: "corpId"
  },
  {
    i18nTitle: "weworkAccessToken.secret",
    dataIndex: "secret",
    key: "secret"
  },
  {
    i18nTitle: "weworkAccessToken.apiKey",
    dataIndex: "apiKey",
    key: "apiKey"
  },
  {
    i18nTitle: "weworkAccessToken.creator",
    dataIndex: "creator",
    key: "creator"
  },
  {
    i18nTitle: "weworkAccessToken.operation",
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
  tokenModal.expired = item.expired;
};
// 复制token
const copyToken = () => {
  message.success(t("copySuccess"));
  window.navigator.clipboard.writeText(tokenModal.token);
};
// 删除token
const deleteAccessToken = item => {
  Modal.confirm({
    title: `${t("weworkAccessToken.confirmDelete")} ${item.name}?`,
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
    title: `${t("weworkAccessToken.confirmRefresh")} ${
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
    title: `${t("weworkAccessToken.confirmChange")} ${item.name} ${t(
      "weworkAccessToken.apiKey"
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
  wwatStore.id = item.id;
  wwatStore.name = item.name;
  wwatStore.corpId = item.corpId;
  wwatStore.secret = item.secret;
  router.push(
    `/team/${route.params.teamId}/weworkAccessToken/${item.id}/update`
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