<template>
  <div style="padding:10px">
    <CommitList :commits="commits" :diffNumsStats="diffNumsStats" style="margin-top:0px" />
    <div>
      <FileDiffDetail
        v-for="item in fileDetails"
        :stat="item"
        :repoId="repoId"
        v-bind:key="item.filePath"
        :target="route.params.commitId"
        :head="headCommitId"
      />
    </div>
  </div>
</template>
<script setup>
import FileDiffDetail from "@/components/git/FileDiffDetail";
import CommitList from "@/components/git/CommitList";
import { diffCommitsRequest } from "@/api/git/repoApi";
import { ref, reactive } from "vue";
import { useRoute } from "vue-router";
const route = useRoute();
// 提交
const commits = ref([]);
const headCommitId = ref("");
// 差异信息
const diffNumsStats = reactive({
  deleteNums: 0,
  insertNums: 0,
  fileChangeNums: 0
});
const repoId = route.params.repoId;
// 差异内容
const fileDetails = ref([]);
//
const diffCommits = () => {
  diffCommitsRequest({
    repoId,
    commitId: route.params.commitId
  }).then(res => {
    commits.value = [res.data.commit];
    if (res.data.diffNumsStats) {
      diffNumsStats.insertNums = res.data.diffNumsStats.insertNums;
      diffNumsStats.deleteNums = res.data.diffNumsStats.deleteNums;
      diffNumsStats.fileChangeNums = res.data.diffNumsStats.fileChangeNums;
      if (res.data.commit.parent && res.data.commit.parent.length > 0) {
        headCommitId.value = res.data.commit.parent[0];
      } 
      fileDetails.value = res.data.diffNumsStats.stats;
    } else {
      diffNumsStats.insertNums = 0;
      diffNumsStats.deleteNums = 0;
      diffNumsStats.fileChangeNums = 0;
      fileDetails.value = [];
    }
  });
};
diffCommits();
</script>
<style scoped>
</style>