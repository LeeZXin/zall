<template>
  <div style="padding:14px">
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
      >{{t("gitRepo.createRepoText")}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="repoList" v-if="wholeRepoList.length > 0">
      <template #bodyCell="{dataIndex, dataItem}">
        <div v-if="dataIndex === 'operation'">
          <div class="op-icon" @click="checkRepo">
            <a-tooltip placement="top">
              <template #title>
                <span>查看</span>
              </template>
              <eye-outlined />
            </a-tooltip>
          </div>
        </div>
        <span v-else>{{dataItem[dataIndex]}}</span>
      </template>
    </ZTable>
    <div v-else class="no-data">
      <span v-if="canCreateRepo">暂无仓库数据, 你可点击上方"创建仓库"</span>
      <span v-else>暂无仓库数据, 管理员已禁用“创建仓库”权限</span>
    </div>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { EyeOutlined } from "@ant-design/icons-vue";
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter, useRoute } from "vue-router";
import { getTeamPermRequest } from "@/api/team/teamApi";
import { getRepoListRequest } from "@/api/git/gitApi";
import { useTeamStore } from "@/pinia/teamStore";
const team = useTeamStore();
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
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/gitRepo/create`);
};

const columns = ref([
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
    title: "最近更新时间",
    dataIndex: "updated",
    key: "updated"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);
getTeamPermRequest({
  teamId: team.teamId
}).then(res => {
  canCreateRepo.value = res.data.canCreateRepo;
});
getRepoListRequest({
  teamId: team.teamId
}).then(res => {
  const ret = res.data.map(item => {
    return {
      key: item.id,
      ...item
    };
  });
  wholeRepoList.value = ret;
  repoList.value = ret;
});
</script>
<style scoped>
.no-data {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  line-height: 80px;
  font-size: 16px;
  text-align: center;
  color: gray;
}
</style>