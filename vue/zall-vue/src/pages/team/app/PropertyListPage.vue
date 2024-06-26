<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建配置文件</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <div class="body">
      <ZTable :columns="columns" :dataSource="dataSource">
        <template #bodyCell="{dataIndex, dataItem}">
          <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <div v-else>
            <div class="op-icon" @click="deleteApp">
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
                    <file-text-outlined />
                    <span style="margin-left:4px">版本列表</span>
                  </li>
                  <li>
                    <eye-outlined />
                    <span style="margin-left:4px">发布历史</span>
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
  </div>
</template>
<script setup>
import {
  DeleteOutlined,
  FileTextOutlined,
  EyeOutlined,
  PlusOutlined,
  EllipsisOutlined
} from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref, h } from "vue";
import { useRoute, useRouter } from "vue-router";
import { listPropertyFileRequest } from "@/api/app/propertyApi";
import { usePropertyFileStore } from "@/pinia/propertyFileStore";
import EnvSelector from "@/components/app/EnvSelector";
const propertyFileStore = usePropertyFileStore();
const selectedEnv = ref("");
const route = useRoute();
const router = useRouter();
const dataSource = ref([]);

const columns = ref([
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
]);

const listPropertyFile = () => {
  listPropertyFileRequest(
    {
      appId: route.params.appId,
      env: selectedEnv.value
    },
    selectedEnv.value
  ).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.name,
        ...item
      };
    });
  });
};

const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/property/create?env=${selectedEnv.value}`
  );
};

const gotoHistoryListPage = item => {
  propertyFileStore.id = item.id;
  propertyFileStore.name = item.name;
  propertyFileStore.env = item.env;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/property/${item.id}/history/list`
  );
};

const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/property/list/${e.newVal}`
  );
  listPropertyFile();
};
</script>
<style scoped>
.body {
  width: 100%;
  height: calc(100% - 42px);
}
</style>