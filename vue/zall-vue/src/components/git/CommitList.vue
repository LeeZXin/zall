<template>
  <div class="commit-list" :style="props.style">
    <div class="title no-wrap">
      <span style="font-weight:bold;padding-left:4px">{{props.commits.length}}</span>
      <span style="color:orange;padding-left:4px">{{t('commitList.commits')}}</span>
      <span style="font-weight:bold;padding-left:12px">{{props.diffNumsStats.fileChangeNums}}</span>
      <span style="color:green;padding-left:4px">{{t('commitList.fileChangeNums')}}</span>
      <span style="font-weight:bold;padding-left:12px">{{props.diffNumsStats.insertNums}}</span>
      <span style="color:green;padding-left:4px">{{t('commitList.insertNums')}}</span>
      <span style="font-weight:bold;padding-left:12px">{{props.diffNumsStats.deleteNums}}</span>
      <span style="color:red;padding-left:4px">{{t('commitList.deleteNums')}}</span>
    </div>
    <ul class="commit-item-list">
      <li v-for="item in props.commits" v-bind:key="item.commitId">
        <div class="commit-item">
          <div class="title">{{item.commitMsg}}</div>
          <div class="desc flex-center">
            <ZAvatar :url="item.committer?.avatarUrl" :name="item.committer?.name" :showName="true"/>
            <span>{{t('commitList.committedAt')}}</span>
            <span>{{readableTimeComparingNow(item.committedTime)}}</span>
          </div>
        </div>
        <CommitSha>{{item.shortId}}</CommitSha>
      </li>
    </ul>
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import CommitSha from "@/components/git/CommitSha";
import { defineProps } from "vue";
import { readableTimeComparingNow } from "@/utils/time";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
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
  padding: 16px;
}
.commit-item-list > li {
  padding: 16px;
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
  font-size: 14px;
  color: gray;
}
.commit-item > .desc > span + span {
  padding-left: 6px;
}
</style>