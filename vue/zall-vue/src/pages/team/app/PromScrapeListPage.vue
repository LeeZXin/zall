<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <div>
        <a-input
          v-model:value="searchEndpointKey"
          placeholder="搜索endpoint"
          style="width:240px;margin-right:6px"
          @pressEnter="searchPromScrape"
        >
          <template #suffix>
            <SearchOutlined />
          </template>
        </a-input>
        <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建抓取任务</a-button>
      </div>

      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex === 'targetType'">{{targetTypeMap[dataItem[dataIndex]]}}</span>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deletePromScrape(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Source</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <edit-outlined />
                  <span style="margin-left:4px">编辑抓取任务</span>
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
      @change="()=>listPromScrape()"
    />
  </div>
</template>
<script setup>
import {
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  SearchOutlined
} from "@ant-design/icons-vue";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, h, createVNode, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listPromScrapeByTeamRequest,
  deletePromScrapeByTeamRequest
} from "@/api/app/promApi";
import { usePromScrapeStore } from "@/pinia/promScrapeStore";
import { Modal, message } from "ant-design-vue";
const promScrapeStore = usePromScrapeStore();
// 分页数据
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
const searchEndpointKey = ref("");
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
const targetTypeMap = {
  1: "服务发现类型",
  2: "直连类型"
};
const columns = [
  {
    title: "endpoint",
    dataIndex: "endpoint",
    key: "endpoint"
  },
  {
    title: "目标",
    dataIndex: "target",
    key: "target"
  },
  {
    title: "目标类型",
    dataIndex: "targetType",
    key: "targetType"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];

const deletePromScrape = item => {
  Modal.confirm({
    title: `你确定要删除${item.endpoint}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deletePromScrapeByTeamRequest(item.id).then(() => {
        message.success("删除成功");
        searchPromScrape();
      });
    }
  });
};

const listPromScrape = () => {
  listPromScrapeByTeamRequest({
    endpoint: searchEndpointKey.value,
    appId: route.params.appId,
    env: selectedEnv.value,
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
// 搜索抓取任务
const searchPromScrape = () => {
  dataPage.current = 1;
  listPromScrape();
};

const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/promScrape/create?env=${selectedEnv.value}`
  );
};

const gotoUpdatePage = item => {
  promScrapeStore.id = item.id;
  promScrapeStore.endpoint = item.endpoint;
  promScrapeStore.env = item.env;
  promScrapeStore.targetType = item.targetType;
  promScrapeStore.target = item.target;
  promScrapeStore.appId = item.appId;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/promScrape/${item.id}/update`
  );
};

const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/promScrape/list/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  searchPromScrape();
};
</script>
<style scoped>
</style>