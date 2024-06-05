<template>
  <div style="padding:14px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchApp"
        placeholder="搜索应用服务"
        style="width:240px;margin-right:10px"
        @change="searchChange"
      />
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">创建应用服务</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <span class="check-btn" @click="gotoAppPage(dataItem)">查看</span>
        </div>
      </template>
    </ZTable>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { PlusOutlined } from "@ant-design/icons-vue";
import { listAppRequest } from "@/api/app/appApi";
import { ref, h } from "vue";
import { useRouter, useRoute } from "vue-router";
const router = useRouter();
const route = useRoute();
const searchApp = ref("");
const dataSource = ref([]);
const allAppList = ref([]);
const columns = ref([
  {
    title: "AppId",
    dataIndex: "appId",
    key: "appId"
  },
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);

const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/app/create`);
};

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

const gotoAppPage = item => {
  router.push(`/team/${route.params.teamId}/app/${item.appId}/property/list`);
}
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