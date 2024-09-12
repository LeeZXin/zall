<template>
  <div style="padding:10px">
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
    <ZNoData v-else />
    <a-pagination
      v-model:current="dataPage.current"
      :total="dataPage.totalCount"
      show-less-items
      :pageSize="dataPage.pageSize"
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
import { ref, reactive } from "vue";
import { useRoute } from "vue-router";
import { listTagCommitsRequest } from "@/api/git/repoApi";
const route = useRoute();
const dataSource = ref([]);
// 分页
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
// 获取列表
const listTag = () => {
  listTagCommitsRequest({
    repoId: route.params.repoId,
    pageNum: dataPage.current
  }).then(res => {
    dataPage.totalCount = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.name,
        name: item.name,
        taggerAccount: item.commit.tagger.account,
        taggerTime: item.commit.taggerTime,
        commitId: item.commit.shortId,
        longCommitId: item.commit.commitId,
        tagCommitMsg: item.commit.tagCommitMsg,
        verified: item.commit.verified,
        signer: item.commit.signer
      };
    });
  });
};
//
const onDelete = () => {
  dataPage.current = 1;
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