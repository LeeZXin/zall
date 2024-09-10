<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <div>
        <a-input
          v-model:value="searchEndpointKey"
          style="width:240px;margin-right:6px"
          @pressEnter="searchPromScrape"
        >
          <template #suffix>
            <SearchOutlined />
          </template>
        </a-input>
        <a-select
          style="width: 200px;margin-right:10px"
          v-model:value="searchAppIdKey"
          :options="appList"
          show-search
          :filter-option="filterAppListOption"
          @change="()=>searchPromScrape()"
        />
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
        <span v-else-if="dataIndex === 'app'">{{dataItem['appName']}}({{dataItem['appId']}})</span>
        <span v-else-if="dataIndex === 'team'">{{dataItem['teamName']}}</span>
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
  listPromScrapeBySaRequest,
  deletePromScrapeBySaRequest
} from "@/api/app/promApi";
import { listAllAppBySaRequest } from "@/api/app/appApi";
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
const searchEndpointKey = ref("");
const searchAppIdKey = ref("");
const appList = ref([
  {
    value: "",
    label: "ALL"
  }
]);
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
const targetTypeMap = {
  1: "promScrape.discoveryType",
  2: "promScrape.hostType"
};
const columns = [
  {
    i18nTitle: "promScrape.endpoint",
    dataIndex: "endpoint",
    key: "endpoint"
  },
  {
    i18nTitle: "promScrape.team",
    dataIndex: "team",
    key: "team"
  },
  {
    i18nTitle: "promScrape.app",
    dataIndex: "app",
    key: "app"
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
    key: "operation",
    width: 130,
    fixed: "right"
  }
];

const deletePromScrape = item => {
  Modal.confirm({
    title: `${t("promScrape.confirmDelete")} ${item.endpoint}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deletePromScrapeBySaRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        searchPromScrape();
      });
    }
  });
};

const listPromScrape = () => {
  listPromScrapeBySaRequest({
    endpoint: searchEndpointKey.value,
    appId: searchAppIdKey.value,
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
// 下拉框搜索过滤
const filterAppListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};

const gotoCreatePage = () => {
  router.push(`/sa/promScrape/create?env=${selectedEnv.value}`);
};

const gotoUpdatePage = item => {
  promScrapeStore.id = item.id;
  promScrapeStore.endpoint = item.endpoint;
  promScrapeStore.env = item.env;
  promScrapeStore.targetType = item.targetType;
  promScrapeStore.target = item.target;
  promScrapeStore.appId = item.appId;
  router.push(`/sa/promScrape/${item.id}/update`);
};

const onEnvChange = e => {
  router.replace(`/sa/promScrape/list/${e.newVal}`);
  selectedEnv.value = e.newVal;
  searchPromScrape();
};
// 获取所有的应用服务
const listAllApp = () => {
  listAllAppBySaRequest().then(res => {
    appList.value = appList.value.concat(
      res.data.map(item => {
        return {
          value: item.appId,
          label: `${item.name}(${item.appId})`
        };
      })
    );
  });
};
listAllApp();
</script>
<style scoped>
</style>