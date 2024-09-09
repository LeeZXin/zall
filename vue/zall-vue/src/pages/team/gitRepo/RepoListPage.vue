<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchRepo"
        :placeholder="t('gitRepo.searchText')"
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
        v-if="teamStore.isAdmin"
        :icon="h(PlusOutlined)"
      >{{t("gitRepo.createRepoText")}}</a-button>
      <a-button
        type="primary"
        @click="gotoRecyclePage"
        :icon="h(DeleteOutlined)"
        danger
        style="float:right"
        v-if="teamStore.isAdmin"
      >仓库回收站</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="repoList" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <span @click="checkRepo(dataItem)" class="check-btn" v-if="dataIndex === 'operation'">查看</span>
        <span v-else-if="dataIndex === 'gitSize'">{{readableVolumeSize(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex === 'lfsSize'">{{readableVolumeSize(dataItem[dataIndex])}}</span>
        <span
          v-else-if="dataIndex === 'lastOperated'"
        >{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex === 'name'">
          <span>{{dataItem[dataIndex]}}</span>
          <a-tag v-if="dataItem['isArchived']" color="red" style="margin-left:4px">已归档</a-tag>
        </span>
        <span v-else>{{dataItem[dataIndex]}}</span>
      </template>
    </ZTable>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import {
  PlusOutlined,
  DeleteOutlined,
  SearchOutlined
} from "@ant-design/icons-vue";
import { ref, h } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter, useRoute } from "vue-router";
import { getRepoListRequest } from "@/api/git/repoApi";
import { readableVolumeSize } from "@/utils/size";
import { readableTimeComparingNow } from "@/utils/time";
import { useTeamStore } from "@/pinia/teamStore";
const teamStore = useTeamStore();
const { t } = useI18n();
const route = useRoute();
const router = useRouter();
// 搜索框key
const searchRepo = ref("");
// 所有仓库列表
const wholeRepoList = ref([]);
// 搜索框检索后仓库列表
const repoList = ref(wholeRepoList.value);
// searchChange 搜索框触发搜索
const searchChange = () => {
  let searchKey = searchRepo.value;
  if (!searchKey || searchKey === "") {
    repoList.value = wholeRepoList.value;
    return;
  }
  repoList.value = wholeRepoList.value.filter(item => {
    return item.name.indexOf(searchKey) >= 0;
  });
};
// 跳转创建仓库页面
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/gitRepo/create`);
};
// 跳转仓库回收站页面
const gotoRecyclePage = () => {
  router.push(`/team/${route.params.teamId}/gitRepo/recycle`);
};
const columns = [
  {
    title: "仓库名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "描述",
    dataIndex: "repoDesc",
    key: "repoDesc"
  },
  {
    title: "仓库大小",
    dataIndex: "gitSize",
    key: "gitSize"
  },
  {
    title: "lfs大小",
    dataIndex: "lfsSize",
    key: "lfsSize"
  },
  {
    title: "最近更新时间",
    dataIndex: "lastOperated",
    key: "lastOperated"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 获取仓库列表
getRepoListRequest(route.params.teamId).then(res => {
  const ret = res.data.map(item => {
    return {
      key: item.repoId,
      ...item
    };
  });
  wholeRepoList.value = ret;
  repoList.value = ret;
});
// 跳转仓库代码首页
const checkRepo = item => {
  router.push(`/team/${route.params.teamId}/gitRepo/${item.repoId}/index`);
};
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