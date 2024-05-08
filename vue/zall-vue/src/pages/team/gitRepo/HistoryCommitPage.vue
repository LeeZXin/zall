<template>
  <div style="padding:14px">
    <template v-if="branches.length > 0">
      <div style="margin-bottom: 10px">
        <BranchTagSelect
          :disableTags="true"
          :branches="branches"
          :defaultBranch="route.params.ref"
          @select="onBranchSelect"
        />
      </div>
      <ul class="commit-list">
        <li v-for="(item, index) in commits" v-bind:key="index">
          <div style="width:70%">
            <div class="commit-msg no-wrap">
              <span class="commit-msg-text" @click="treeCommit(item)">{{item.commitMsg}}</span>
            </div>
            <div class="commit-desc no-wrap">
              <span>{{item.committer.account}}</span>
              <span>提交于</span>
              <span>{{readableTimeComparingNow(item.committedTime)}}</span>
            </div>
          </div>
          <CommitSha>{{item.shortId}}</CommitSha>
        </li>
        <li v-if="lastLoadCount >= 10">
          <div style="width:100%;text-align:center;cursor:pointer" @click="getCommits()">加载更多...</div>
        </li>
      </ul>
    </template>
    <ZNoData v-else>
      <template #desc>
        <div style="text-align:center;font-size:14px">
          <span>无提交数据, 尝试去</span>
          <span class="suggest-text" @click="gotoIndex">提交代码</span>
        </div>
      </template>
    </ZNoData>
  </div>
</template>
<script setup>
import ZNoData from "@/components/common/ZNoData";
import CommitSha from "@/components/git/CommitSha";
import BranchTagSelect from "@/components/git/BranchTagSelect";
import { allBranchesRequest, historyCommitsRequest } from "@/api/git/repoApi";
import { ref, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import { readableTimeComparingNow } from "@/utils/time";
const router = useRouter();
const route = useRoute();
const branches = ref([]);
const commits = ref([]);
const lastLoadCount = ref(0);
const selectedBranch = ref("");
allBranchesRequest(route.params.repoId).then(res => {
  branches.value = res.data;
});
const getCommits = () => {
  historyCommitsRequest({
    repoId: route.params.repoId,
    cursor: commits.value.length,
    ref: selectedBranch.value
  }).then(res => {
    commits.value = [...commits.value, ...res.data];
    lastLoadCount.value = res.data.length;
  });
};
const onBranchSelect = ({ value }) => {
  history.replaceState(
    {},
    "",
    `/gitRepo/${route.params.repoId}/commit/list/${value}`
  );
  selectedBranch.value = value;
  commits.value = [];
  nextTick(() => {
    getCommits();
  });
};
const treeCommit = item => {
  router.push(`/gitRepo/${route.params.repoId}/commit/diff/${item.commitId}`);
};
const gotoIndex = () => {
  router.push(`/gitRepo/${route.params.repoId}/index`);
};
</script>
<style scoped>
.commit-list {
  width: 100%;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.commit-list > li {
  width: 100%;
  padding: 14px;
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}
.commit-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.commit-msg {
  font-size: 16px;
  font-weight: bold;
  padding-bottom: 10px;
}
.commit-desc {
  font-size: 13px;
}
.commit-desc > span + span {
  margin-left: 4px;
}
.load-more {
  width: 100%;
  line-height: 32px;
  text-align: center;
}
.commit-msg-text:hover {
  cursor: pointer;
  color: #1677ff;
}
</style>