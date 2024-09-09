<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchApp"
        :placeholder="t('appService.searchApp')"
        style="width:240px;margin-right:10px"
        @change="searchChange"
      >
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
        v-if="teamStore.isAdmin"
      >{{t('appService.createApp')}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <span class="check-btn" @click="gotoAppPage(dataItem)">{{t('appService.viewApp')}}</span>
        </div>
      </template>
    </ZTable>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { PlusOutlined, SearchOutlined } from "@ant-design/icons-vue";
import { listAppRequest } from "@/api/app/appApi";
import { ref, h } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useTeamStore } from "@/pinia/teamStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const teamStore = useTeamStore();
const router = useRouter();
const route = useRoute();
const searchApp = ref("");
const dataSource = ref([]);
// 所有应用服务列表
const allAppList = ref([]);
// 数据项
const columns = [
  {
    i18nTitle: "appService.appId",
    dataIndex: "appId",
    key: "appId"
  },
  {
    i18nTitle: "appService.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "appService.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/app/create`);
};
// app应用
const listApp = () => {
  listAppRequest(route.params.teamId).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.appId,
        ...item
      };
    });
    allAppList.value = dataSource.value;
  });
};
// 搜索框搜索
const searchChange = () => {
  if (!searchApp.value || searchApp.value === "") {
    dataSource.value = allAppList.value;
  } else {
    dataSource.value = allAppList.value.filter(item => {
      return (
        item.appId.indexOf(searchApp.value) >= 0 ||
        item.name.indexOf(searchApp.value) >= 0
      );
    });
  }
};
// 进入应用服务layout
const gotoAppPage = item => {
  router.push(
    `/team/${route.params.teamId}/app/${item.appId}/propertyFile/list`
  );
};
listApp();
</script>
<style scoped>
.check-btn {
  font-size: 14px;
}
.check-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>