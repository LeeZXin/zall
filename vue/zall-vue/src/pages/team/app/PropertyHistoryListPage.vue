<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <div style="font-size:20px;">
        <span style="margin-right:8px;font-weight:bold">{{propertyFileStore.name}}</span>
        <a-tag color="orange">{{propertyFileStore.env}}</a-tag>
      </div>
    </div>
    <div class="body">
      <ZTable :columns="columns" :dataSource="dataSource">
        <template #bodyCell="{dataIndex, dataItem}">
          <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <div v-else>
            <div class="op-icon" @click="gotoNewPage(dataItem)">
              <a-tooltip placement="top">
                <template #title>
                  <span>{{t('propertyFile.newVersion')}}</span>
                </template>
                <PlusOutlined />
              </a-tooltip>
            </div>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <li @click="showDiffModal(dataItem)">
                    <CodeOutlined />
                    <span style="margin-left:4px">{{t('propertyFile.compareLastVersion')}}</span>
                  </li>
                  <li @click="gotoPubPage(dataItem)" v-if="appStore.perm?.canDeployProperty">
                    <CloudUploadOutlined />
                    <span style="margin-left:4px">{{t('propertyFile.deployVersion')}}</span>
                  </li>
                  <li @click="showDeployModal(dataItem)">
                    <EyeOutlined />
                    <span style="margin-left:4px">{{t('propertyFile.deployRecord')}}</span>
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
        @change="()=>listHistory()"
      />
      <a-modal
        :title="t('propertyFile.compareLastVersion')"
        :footer="null"
        v-model:open="diffModalOpen"
        :width="800"
      >
        <code-diff
          :old-string="diffState.oldContent"
          :new-string="diffState.newContent"
          :context="10"
          outputFormat="side-by-side"
          :hideStat="true"
          :filename="diffState.oldVersion"
          :newFilename="diffState.newVersion"
          style="max-height:400px;overflow:scroll"
        />
      </a-modal>
      <a-modal
        :title="deployModal.title"
        :footer="null"
        v-model:open="deployModal.open"
        :width="800"
      >
        <div style="max-height:600px;overflow:scroll">
          <ZTable :columns="deployColumns" :dataSource="deployDataSource" />
        </div>
      </a-modal>
    </div>
  </div>
</template>
<script setup>
import {
  CodeOutlined,
  EllipsisOutlined,
  PlusOutlined,
  CloudUploadOutlined,
  EyeOutlined
} from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref, reactive } from "vue";
import { usePropertyFileStore } from "@/pinia/propertyFileStore";
import { useAppStore } from "@/pinia/appStore";
import { usePropertyHistoryStore } from "@/pinia/propertyHistoryStore";
import { useRouter, useRoute } from "vue-router";
import {
  listHistoryRequest,
  getHistoryByVersionRequest,
  listDeployRequest
} from "@/api/app/propertyApi";
import { message } from "ant-design-vue";
import { CodeDiff } from "v-code-diff";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const appStore = useAppStore();
// 发布记录modal数据
const deployModal = reactive({
  open: false,
  title: ""
});
// 发布记录表项
const deployColumns = [
  {
    i18nTitle: "propertyFile.nodeName",
    dataIndex: "nodeName",
    key: "nodeName"
  },
  {
    title: "endpoints",
    dataIndex: "endpoints",
    key: "endpoints"
  },
  {
    i18nTitle: "propertyFile.deployTime",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "propertyFile.deployer",
    dataIndex: "creator",
    key: "creator"
  }
];
// 发布记录列表
const deployDataSource = ref([]);
// 对比modal是否展示
const diffModalOpen = ref(false);
const route = useRoute();
const router = useRouter();
const propertyFileStore = usePropertyFileStore();
const propertyHistoryStore = usePropertyHistoryStore();
// 版本数据
const dataSource = ref([]);
// 分页数据
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
// 对比数据
const diffState = reactive({
  newContent: "",
  newVersion: "",
  oldCotent: "",
  oldVersion: ""
});
// 版本数据表项
const columns = [
  {
    i18nTitle: "propertyFile.version",
    dataIndex: "version",
    key: "version"
  },
  {
    i18nTitle: "propertyFile.lastVersion",
    dataIndex: "lastVersion",
    key: "lastVersion"
  },
  {
    i18nTitle: "propertyFile.created",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "propertyFile.creator",
    dataIndex: "creator",
    key: "creator"
  },
  {
    i18nTitle: "propertyFile.operation",
    dataIndex: "operation",
    key: "operation"
  }
];
// 获取版本历史
const listHistory = () => {
  listHistoryRequest({
    fileId: propertyFileStore.id,
    pageNum: dataPage.current
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
// 展示对比modal
const showDiffModal = item => {
  if (item.lastVersion) {
    getHistoryByVersionRequest({
      fileId: item.fileId,
      version: item.lastVersion
    }).then(res => {
      if (res.data.exist) {
        diffModalOpen.value = true;
        diffState.oldContent = res.data.value.content;
        diffState.oldVersion = item.lastVersion;
        diffState.newContent = item.content;
        diffState.newVersion = item.version;
      } else {
        message.warn("跟随版本数据不存在");
      }
    });
    return;
  }
  diffModalOpen.value = true;
  diffState.oldContent = "";
  diffState.oldVersion = "";
  diffState.newContent = item.content;
  diffState.newVersion = item.version;
};

if (propertyFileStore.id === 0) {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/list`
  );
} else {
  listHistory();
}
// 跳转新增版本页面
const gotoNewPage = item => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${propertyFileStore.id}/new?from=${item.version}`
  );
};
// 跳转发布页面
const gotoPubPage = item => {
  propertyHistoryStore.id = item.id;
  propertyHistoryStore.fileName = item.fileName;
  propertyHistoryStore.fileId = item.fileId;
  propertyHistoryStore.content = item.content;
  propertyHistoryStore.version = item.version;
  propertyHistoryStore.created = item.created;
  propertyHistoryStore.creator = item.creator;
  propertyHistoryStore.lastVersion = item.lastVersion;
  propertyHistoryStore.env = item.env;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${propertyFileStore.id}/deploy/${item.version}`
  );
};
// 展示发布记录modal
const showDeployModal = item => {
  deployModal.title = `${item.fileName} ${item.version}`;
  listDeployRequest(item.id).then(res => {
    deployDataSource.value = res.data.map((item, index) => {
      return {
        key: index,
        ...item
      };
    });
    deployModal.open = true;
  });
};
</script>
<style scoped>
.body {
  width: 100%;
  height: calc(100% - 42px);
}
</style>