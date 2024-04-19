<template>
  <div class="body">
    <div style="margin-bottom:20px">
      <a-input
        v-model:value="searchRepo"
        :placeholder="t('gitRepo.searchText')"
        style="width:240px;margin-right:10px"
        @change="searchChange"
      />
      <a-button type="primary" @click="gotoCreatePage">{{t("gitRepo.createRepoText")}}</a-button>
    </div>
    <ul class="repo-list">
      <li v-for="item in repoList" v-bind:key="item.id">
        <div class="repo-name">{{item.name}}</div>
        <div class="repo-desc">
          <span>{{item.repoDesc}}</span>
          <span class="last-updated-text">{{item.lastUpdated}}</span>
        </div>
      </li>
      <li v-show="repoList.length === 0">
          <div style="font-size:32px;text-align:center;color:gray"><inbox-outlined/></div>
          <div class="no-data">无数据</div>
      </li>
    </ul>
  </div>
</template>
<script setup>
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import {
  InboxOutlined
} from "@ant-design/icons-vue";
const { t } = useI18n();
const router = useRouter();
const searchRepo = ref("");
const wholeRepoList = ref([
  {
    id: 1,
    name: "zall",
    repoDesc: "一体化devops",
    lastUpdated: "updated 18 hours ago"
  },
  {
    id: 2,
    name: "zsf",
    repoDesc: "一体化devops",
    lastUpdated: "updated 18 hours ago"
  }
]);
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
  router.push("/team/gitRepo/create");
};
</script>
<style scoped>
.body {
  padding: 20px;
}
.repo-name {
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
  margin-bottom: 10px;
}
.repo-name:hover {
  color: #1677ff;
}
.last-updated-text {
  margin-left: 8px;
}
.repo-list {
  border-radius: 4px;
  border: 1px solid #d9d9d9;
}
.repo-list > li {
  padding:14px;
}
.repo-desc {
  line-height: 28px;
  font-size: 14px;
  color: gray;
}
.repo-list > li+li {
    border-top:1px solid #d9d9d9;
}
.no-data {
    line-height: 48px;
    font-size: 16px;
    text-align: center;
    color:gray;
}
</style>