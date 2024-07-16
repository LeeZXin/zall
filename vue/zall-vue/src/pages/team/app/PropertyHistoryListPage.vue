<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <div style="font-size:20px;font-weight:bold">
        <span style="margin-right:8px">{{propertyFileStore.name}}</span>
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
                  <span>在此版本号新增版本</span>
                </template>
                <plus-outlined />
              </a-tooltip>
            </div>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <li @click="showDiffModal(dataItem)">
                    <code-outlined />
                    <span style="margin-left:4px">对比跟随版本</span>
                  </li>
                  <li @click="gotoPubPage(dataItem)">
                    <cloud-upload-outlined />
                    <span style="margin-left:4px">发布此版本</span>
                  </li>
                  <li @click="showDeployModal(dataItem)">
                    <eye-outlined />
                    <span style="margin-left:4px">发布记录</span>
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
        @change="()=>listHistory()"
      />
      <a-modal title="跟随版本对比" :footer="null" v-model:open="diffModalOpen" :width="800">
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
      <a-modal :title="deployModalTitle" :footer="null" v-model:open="deployModalOpen" :width="800">
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
import { usePropertyHistoryStore } from "@/pinia/propertyHistoryStore";
import { useRouter, useRoute } from "vue-router";
import {
  listHistoryRequest,
  getHistoryByVersionRequest,
  listDeployRequest
} from "@/api/app/propertyApi";
import { message } from "ant-design-vue";
import { CodeDiff } from "v-code-diff";
const deployModalOpen = ref(false);
const deployModalTitle = ref("");
const deployColumns = [
  {
    title: "发布节点",
    dataIndex: "nodeName",
    key: "nodeName"
  },
  {
    title: "endpoints",
    dataIndex: "endpoints",
    key: "endpoints"
  },
  {
    title: "发布时间",
    dataIndex: "created",
    key: "created"
  },
  {
    title: "发布人",
    dataIndex: "creator",
    key: "creator"
  }
];
const deployDataSource = ref([]);
const diffModalOpen = ref(false);
const route = useRoute();
const router = useRouter();
const propertyFileStore = usePropertyFileStore();
const propertyHistoryStore = usePropertyHistoryStore();
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const diffState = reactive({
  newContent: "",
  newVersion: "",
  oldCotent: "",
  oldVersion: ""
});
const columns = [
  {
    title: "跟随版本号",
    dataIndex: "lastVersion",
    key: "lastVersion"
  },
  {
    title: "版本号",
    dataIndex: "version",
    key: "version"
  },
  {
    title: "创建时间",
    dataIndex: "created",
    key: "created"
  },
  {
    title: "创建人",
    dataIndex: "creator",
    key: "creator"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];

const listHistory = () => {
  listHistoryRequest({
    fileId: propertyFileStore.id,
    pageNum: currPage.value
  }).then(res => {
    totalCount.value = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};

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

const gotoNewPage = item => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${propertyFileStore.id}/new?from=${item.version}`
  );
};

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
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${propertyFileStore.id}/publish/${item.version}`
  );
};

const showDeployModal = item => {
  deployModalTitle.value = `${item.fileName}_${item.version}`;
  listDeployRequest(item.id).then(res => {
    deployDataSource.value = res.data.map((item, index) => {
      return {
        key: index,
        ...item
      };
    });
    deployModalOpen.value = true;
  });
};
</script>
<style scoped>
.body {
  width: 100%;
  height: calc(100% - 42px);
}
</style>