<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px" class="flex-center">
      <a-input
        v-model:value="searchKey"
        placeholder="搜索名称或appId"
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
      >刷新列表</a-button>
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">创建AccessToken任务</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon" @click="deleteAccessToken(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>删除任务</span>
              </template>
              <DeleteOutlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">编辑任务</span>
                </li>
                <li @click="changeApiKey(dataItem)">
                  <ReloadOutlined />
                  <span style="margin-left:4px">变更api key</span>
                </li>
                <li @click="refreshAccessToken(dataItem)">
                  <ReloadOutlined />
                  <span style="margin-left:4px">重刷token</span>
                </li>
                <li @click="viewToken(dataItem)">
                  <EyeOutlined />
                  <span style="margin-left:4px">查看token</span>
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
          <span>复制</span>
        </div>
      </div>
      <div>tenant token</div>
      <div style="padding:10px">
        <span style="word-break:break-all;">{{tokenModal.tenantToken}}</span>
        <div class="copy-btn" @click="copyTenantToken">
          <CopyOutlined />
          <span>复制</span>
        </div>
      </div>
      <div>过期时间</div>
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
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "appId",
    dataIndex: "appId",
    key: "appId"
  },
  {
    title: "secret",
    dataIndex: "secret",
    key: "secret"
  },
  {
    title: "Api密钥",
    dataIndex: "apiKey",
    key: "apiKey"
  },
  {
    title: "创建人",
    dataIndex: "creator",
    key: "creator"
  },
  {
    title: "操作",
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
  message.success("复制成功");
  window.navigator.clipboard.writeText(tokenModal.token);
};
// 复制tenant token
const copyTenantToken = () => {
  message.success("复制成功");
  window.navigator.clipboard.writeText(tokenModal.tenantToken);
};
// 删除token
const deleteAccessToken = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteAccessTokenRequest(item.id).then(() => {
        message.success("删除成功");
        searchAccessToken();
      });
    },
    onCancel() {}
  });
};
// 刷新token
const refreshAccessToken = item => {
  Modal.confirm({
    title: `你确定要刷新${item.name}的access token吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      refreshAccessTokenRequest(item.id).then(() => {
        message.success("刷新成功");
        listAccessToken();
      });
    },
    onCancel() {}
  });
};
// 变更api key
const changeApiKey = item => {
  Modal.confirm({
    title: `你确定要刷新${item.name}的api密钥吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      changeAccessTokenApiKeyRequest(item.id).then(() => {
        message.success("变更成功");
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