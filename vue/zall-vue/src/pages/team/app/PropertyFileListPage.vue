<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <div>
        <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建配置文件</a-button>
        <a-button
          type="primary"
          :icon="h(SettingOutlined)"
          @click="showBindPropertySourceModal"
          style="margin-left:6px"
          v-if="appStore.perm?.canManagePropertySource"
        >管理配置来源绑定</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteFile(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete File</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoHistoryListPage(dataItem)">
                  <FileTextOutlined />
                  <span style="margin-left:4px">版本列表</span>
                </li>
                <li @click="showSearchModal(dataItem)">
                  <SearchOutlined />
                  <span style="margin-left:4px">配置查询</span>
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
  </div>
  <a-modal v-model:open="bindModal.open" title="绑定配置来源" @ok="handleBindModalOk">
    <div>
      <div style="font-size:12px;margin-bottom:3px">已选环境</div>
      <div>{{selectedEnv}}</div>
    </div>
    <div style="margin-top: 10px">
      <div style="font-size:12px;margin-bottom:3px">配置来源</div>
      <a-select
        style="width: 100%"
        placeholder="请选择"
        v-model:value="bindModal.selectedIdList"
        :options="bindModal.sourceList"
        mode="multiple"
      />
    </div>
  </a-modal>
  <a-modal v-model:open="searchModal.open" title="配置查询" :width="600" :footer="null">
    <ul class="search-ul">
      <li>
        <div class="left">已选文件:</div>
        <div class="right">{{searchModal.fileName}}</div>
      </li>
      <li>
        <div class="left">已选环境:</div>
        <div class="right">{{selectedEnv}}</div>
      </li>
      <li>
        <div class="left">配置来源:</div>
        <div class="right">
          <a-select
            style="width: 100%"
            placeholder="请选择"
            v-model:value="searchModal.sourceId"
            :options="searchModal.sourceList"
            @change="searchFromSource"
          />
        </div>
      </li>
    </ul>
    <div v-if="searchModal.loading">
      <div style="text-align:center;padding:10px 0;font-size: 14px">
        <LoadingOutlined />
        <span style="margin-left:8px">查询中</span>
      </div>
    </div>
    <div style="margin-top: 10px" v-else>
      <div v-if="searchModal.exist">
        <ul class="search-ul">
          <li>
            <div class="left">版本:</div>
            <div class="right">{{searchModal.version}}</div>
          </li>
        </ul>
        <div class="no-wrap" style="margin-top:10px;font-size:14px;margin-bottom:6px">内容</div>
        <Codemirror
          v-model="searchModal.content"
          style="height:280px;width:100%"
          :extensions="extensions"
          :disabled="true"
        />
      </div>
      <ZNoData v-else />
    </div>
  </a-modal>
</template>
<script setup>
import ZNoData from "@/components/common/ZNoData";
import {
  DeleteOutlined,
  FileTextOutlined,
  SearchOutlined,
  PlusOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  SettingOutlined,
  LoadingOutlined
} from "@ant-design/icons-vue";
import { Codemirror } from "vue-codemirror";
import { oneDark } from "@codemirror/theme-one-dark";
import ZTable from "@/components/common/ZTable";
import { ref, h, createVNode, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listPropertyFileRequest,
  deletePropertyFileRequest,
  listAllPropertySourceRequest,
  listBindPropertySourceRequest,
  bindAppAndPropertySourceRequest,
  searchFromSourceRequest
} from "@/api/app/propertyApi";
import { usePropertyFileStore } from "@/pinia/propertyFileStore";
import EnvSelector from "@/components/app/EnvSelector";
import { Modal, message } from "ant-design-vue";
import { useAppStore } from "@/pinia/appStore";
const extensions = [oneDark];
const appStore = useAppStore();
const propertyFileStore = usePropertyFileStore();
// 来源绑定modal
const bindModal = reactive({
  open: false,
  selectedIdList: [],
  sourceList: []
});
// 配置查询modal
const searchModal = reactive({
  open: false,
  sourceId: null,
  sourceList: [],
  fileId: 0,
  fileName: "",
  exist: false,
  version: "",
  content: "",
  loading: false
});
// 当前环境
const selectedEnv = ref("");
const route = useRoute();
const router = useRouter();
// 表格数据
const dataSource = ref([]);
// 数据项
const columns = [
  {
    title: "配置文件",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];

// 获取配置文件列表
const listPropertyFile = () => {
  listPropertyFileRequest({
    appId: route.params.appId,
    env: selectedEnv.value
  }).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.name,
        ...item
      };
    });
  });
};
// 跳转创建配置文件页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/create?env=${selectedEnv.value}`
  );
};
// 跳转历史版本页面
const gotoHistoryListPage = item => {
  propertyFileStore.id = item.id;
  propertyFileStore.name = item.name;
  propertyFileStore.env = item.env;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${item.id}/history/list`
  );
};
// 删除配置文件
const deleteFile = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deletePropertyFileRequest(item.id).then(() => {
        message.success("删除成功");
        listPropertyFile();
      });
    }
  });
};
// 选择环境变动
const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/list/${e.newVal}`
  );
  listPropertyFile();
};
// 展示绑定配置来源modal
const showBindPropertySourceModal = () => {
  if (!selectedEnv.value) {
    return;
  }
  listAllPropertySourceRequest(selectedEnv.value).then(res => {
    bindModal.sourceList = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
    listBindPropertySourceRequest({
      appId: route.params.appId,
      env: selectedEnv.value
    }).then(res => {
      bindModal.selectedIdList = res.data.map(item => item.id);
      bindModal.open = true;
    });
  });
};
// 展示配置查询modal
const showSearchModal = item => {
  searchModal.sourceId = null;
  searchModal.fileId = item.id;
  searchModal.fileName = item.name;
  searchModal.exist = false;
  listBindPropertySourceRequest({
    appId: route.params.appId,
    env: selectedEnv.value
  }).then(res => {
    searchModal.sourceList = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
    searchModal.open = true;
  });
};
const searchFromSource = () => {
  if (!searchModal.sourceId) {
    message.warn("请选择配置来源");
    return;
  }
  searchModal.loading = true;
  searchFromSourceRequest({
    fileId: searchModal.fileId,
    sourceId: searchModal.sourceId
  })
    .then(res => {
      searchModal.exist = res.data.exist;
      searchModal.version = res.data.version;
      searchModal.content = res.data.content;
      searchModal.loading = false;
    })
    .catch(() => {
      searchModal.loading = false;
    });
};
// 绑定modal点击“确定”按钮
const handleBindModalOk = () => {
  bindAppAndPropertySourceRequest({
    appId: route.params.appId,
    sourceIdList: bindModal.selectedIdList,
    env: selectedEnv.value
  }).then(() => {
    message.success("操作成功");
    bindModal.open = false;
  });
};
</script>
<style scoped>
.search-ul > li {
  width: 100%;
  font-size: 14px;
  display: flex;
  align-items: center;
}
.search-ul > li + li {
  margin-top: 10px;
}
.search-ul > li > .left {
  width: 80px;
  white-space: nowrap;
  overflow: hidden;
  word-break: break-all;
  text-overflow: ellipsis;
}
.search-ul > li > .right {
  width: calc(100% - 80px);
  white-space: nowrap;
  overflow: hidden;
  word-break: break-all;
}
</style>