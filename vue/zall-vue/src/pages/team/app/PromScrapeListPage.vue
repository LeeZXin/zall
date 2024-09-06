<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <div>
        <a-input
          v-model:value="searchEndpointKey"
          :placeholder="t('promScrape.searchEndpoint')"
          style="width:240px;margin-right:6px"
          @pressEnter="searchPromScrape"
        >
          <template #suffix>
            <SearchOutlined />
          </template>
        </a-input>
        <a-button
          type="primary"
          :icon="h(PlusOutlined)"
          @click="gotoCreatePage"
        >{{t('promScrape.createScrape')}}</a-button>
      </div>

      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex === 'targetType'">{{t(targetTypeMap[dataItem[dataIndex]])}}</span>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deletePromScrape(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">{{t('promScrape.updateScrape')}}</span>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const promScrapeStore = usePromScrapeStore();
// 分页数据
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
// 搜索关键词
const searchEndpointKey = ref("");
const route = useRoute();
// 选择的环境
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
// 服务发现类型map
const targetTypeMap = {
  1: "promScrape.discoveryType",
  2: "promScrape.hostType"
};
// 数据项
const columns = [
  {
    i18nTitle: "promScrape.endpoint",
    dataIndex: "endpoint",
    key: "endpoint"
  },
  {
    i18nTitle: "promScrape.target",
    dataIndex: "target",
    key: "target"
  },
  {
    i18nTitle: "promScrape.targetType",
    dataIndex: "targetType",
    key: "targetType"
  },
  {
    i18nTitle: "promScrape.operation",
    dataIndex: "operation",
    key: "operation"
  }
];
// 删除任务
const deletePromScrape = item => {
  Modal.confirm({
    title: `${t("promScrape.confirmDelete")} ${item.endpoint}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deletePromScrapeByTeamRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchPromScrape();
      });
    }
  });
};
// 任务列表
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
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/promScrape/create?env=${selectedEnv.value}`
  );
};
// 跳转编辑页面
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
// 环境变化
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