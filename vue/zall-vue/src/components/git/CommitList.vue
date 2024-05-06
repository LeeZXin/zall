<template>
  <div class="commit-list" :style="props.style">
    <div class="title no-wrap">
      <span style="font-weight:bold">{{props.commits.length}}</span>
      <span style="color:orange">个提交</span>
      <span style="font-weight:bold">{{props.diffNumsStats.fileChangeNums}}</span>
      <span style="color:green">个文件修改</span>
      <span style="font-weight:bold">{{props.diffNumsStats.insertNums}}</span>
      <span style="color:green">次新增</span>
      <span style="font-weight:bold">{{props.diffNumsStats.deleteNums}}</span>
      <span style="color:red">次删除</span>
    </div>
    <ul class="commit-item-list">
      <li v-for="item in props.commits" v-bind:key="item.commitId">
        <div class="commit-item">
          <div class="title">{{item.commitMsg}}</div>
          <div class="desc">
            <span>{{item.committer.account}}</span>
            <span>提交于</span>
            <span>{{readableTimeComparingNow(item.committedTime)}}</span>
          </div>
        </div>
        <CommitSha>{{item.shortId}}</CommitSha>
      </li>
    </ul>
  </div>
</template>
<script setup>
import CommitSha from "@/components/git/CommitSha";
import { defineProps } from "vue";
import { readableTimeComparingNow } from "@/utils/time";
const props = defineProps(["commits", "diffNumsStats", "style"]);
</script>
<style scoped>
.commit-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  margin-top: 10px;
}
.commit-list > .title {
  width: 100%;
  font-size: 14px;
  padding: 10px;
}
.commit-list > .title > span+span {
  padding-left: 6px;
}
.commit-item-list > li {
  padding: 10px;
  border-top: 1px solid #d9d9d9;
  display: flex;
  width: 100%;
  justify-content: space-between;
  align-items: baseline;
}
.commit-item {
  width: 50%;
}
.commit-item > .title {
  font-size: 16px;
  padding-bottom: 8px;
  font-weight: bold;
  width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.commit-item > .desc {
  font-size: 12px;
  color: gray;
}
.commit-item > .desc > span+span {
  margin-left: 4px;
}
</style>