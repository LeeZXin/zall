<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchRepo"
        :placeholder="t('gitRepo.searchText')"
        style="width:240px;margin-right:10px"
        @change="searchChange"
      />
      <a-button
        type="primary"
        @click="gotoCreatePage"
        v-if="canCreateRepo"
        :icon="h(PlusOutlined)"
      >{{t("gitRepo.createRepoText")}}</a-button>
      <a-button
        type="primary"
        @click="gotoRecyclePage"
        :icon="h(DeleteOutlined)"
        danger
        style="float:right"
      >仓库回收站</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="repoList" v-if="wholeRepoList.length > 0">
      <template #bodyCell="{dataIndex, dataItem}">
        <span @click="checkRepo(dataItem)" class="check-btn" v-if="dataIndex === 'operation'">查看</span>
        <span v-else-if="dataIndex === 'gitSize'">{{readableVolumeSize(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex === 'lfsSize'">{{readableVolumeSize(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex === 'lastOperated'">{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex === 'name'">
          <span>{{dataItem[dataIndex]}}</span>
          <a-tag v-if="dataItem['isArchived']" color="red" style="margin-left:4px">已归档</a-tag>
        </span>
        <span v-else>{{dataItem[dataIndex]}}</span>
      </template>
    </ZTable>
    <ZNoData v-else>
      <template #desc>
        <div class="no-data">
          <span v-if="canCreateRepo">暂无仓库数据, 你可点击上方"创建仓库"</span>
          <span v-else>暂无仓库数据, 管理员已禁用“创建仓库”权限</span>
        </div>
      </template>
    </ZNoData>
  </div>
</template>
<script setup>
import ZNoData from "@/components/common/ZNoData";
import ZTable from "@/components/common/ZTable";
import { PlusOutlined, DeleteOutlined } from "@ant-design/icons-vue";
import { ref, h } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter, useRoute } from "vue-router";
import { getTeamPermRequest } from "@/api/team/teamApi";
import { getRepoListRequest } from "@/api/git/repoApi";
import { useRepoStore } from "@/pinia/repoStore";
import { readableVolumeSize } from "@/utils/size";
import { readableTimeComparingNow } from "@/utils/time";
const { t } = useI18n();
const route = useRoute();
const router = useRouter();
// 是否可以创建仓库
const canCreateRepo = ref(false);
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
    key: "operation"
  }
];
// 获取团队权限 判断是否可以创建仓库
getTeamPermRequest(route.params.teamId).then(res => {
  canCreateRepo.value = res.data.canCreateRepo;
});
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
  const repo = useRepoStore();
  repo.repoId = item.repoId;
  repo.name = item.name;
  repo.teamId = item.teamId;
  router.push(`/team/${route.params.teamId}/gitRepo/${item.repoId}/index`);
};
</script>
<style scoped>
.no-data {
  font-size: 16px;
  text-align: center;
}
.check-btn {
  font-size: 14px;
}
.check-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>