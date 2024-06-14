<template>
  <div style="padding:14px">
    <ul class="tag-list" v-if="dataSource.length > 0">
      <li v-for="item in dataSource" v-bind:key="item.name">
        <TagItem
          :data="item"
          :repoId="route.params.repoId"
          :teamId="route.params.teamId"
          @delete="onDelete"
        />
      </li>
    </ul>
    <ZNoData v-else>
      <template #desc>
        <div style="font-size:14px;text-align:center">
          <span>无分支数据, 尝试去git tag -a xx -m xx</span>
        </div>
      </template>
    </ZNoData>
    <a-pagination
      v-model:current="currPage"
      :total="totalCount"
      show-less-items
      :pageSize="pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listTag()"
    />
  </div>
</template>
<script setup>
import TagItem from "@/components/git/TagItem";
import ZNoData from "@/components/common/ZNoData";
import { ref } from "vue";
import { useRoute } from "vue-router";
import { pageTagCommitsRequest } from "@/api/git/repoApi";
const route = useRoute();
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const listTag = () => {
  pageTagCommitsRequest({
    repoId: route.params.repoId,
    pageNum: currPage.value
  }).then(res => {
    totalCount.value = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.name,
        name: item.name,
        taggerAccount: item.commit.tagger.account,
        taggerTime: item.commit.taggerTime,
        commitId: item.commit.shortId,
        longCommitId: item.commit.commitId,
        tagCommitMsg: item.commit.tagCommitMsg,
        verified: item.commit.verified
      };
    });
  });
};
const onDelete = () => {
  if (totalCount.value - 1 <= (currPage.value - 1) * pageSize) {
    currPage.value -= 1;
  }
  listTag();
};
listTag();
</script>
<style scoped>
.tag-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.tag-list > li + li {
  border-top: 1px solid #d9d9d9;
}
</style>