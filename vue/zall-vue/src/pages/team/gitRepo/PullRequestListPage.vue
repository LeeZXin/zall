<template>
  <div style="padding:10px">
    <div class="header">
      <a-select v-model:value="prStatus" @change="selectPrStatus" style="width:180px">
        <a-select-option :value="0">{{t('pullRequest.allStatus')}}</a-select-option>
        <a-select-option :value="1">{{t('pullRequest.openStatus')}}</a-select-option>
        <a-select-option :value="3">{{t('pullRequest.mergedStatus')}}</a-select-option>
        <a-select-option :value="2">{{t('pullRequest.closedStatus')}}</a-select-option>
      </a-select>
      <a-input
        :placeholder="t('pullRequest.searchPr')"
        v-model:value="searchKey"
        @pressEnter="searchPullRequest"
        style="width:240px;margin-left:10px"
      >
        <template #prefix>
          <search-outlined />
        </template>
      </a-input>
      <a-button
        type="primary"
        style="margin-left:10px"
        @click="toCreatePage"
        :icon="h(PlusOutlined)"
        v-if="repoStore.perm?.canSubmitPullRequest"
      >{{t('pullRequest.createPr')}}</a-button>
    </div>
    <ul class="pr-list" v-if="prList.length > 0">
      <li v-for="item in prList" v-bind:key="item.id" @click="toDetail(item)">
        <div class="pr-title">
          <plus-circle-outlined v-if="item.prStatus === 1" />
          <close-circle-outlined v-else-if="item.prStatus === 2" />
          <check-circle-outlined v-else-if="item.prStatus === 3" />
          <span class="pr-title-span">{{item.prTitle}}</span>
          <span class="pr-title-right">
            <message-outlined />
            <span class="message-num">{{item.commentCount}}</span>
          </span>
        </div>
        <div class="pr-desc flex-center">
          <PrStatusTag :status="item.prStatus" />
          <ZAvatar
            :url="item.createBy?.avatarUrl"
            :name="item.createBy?.name"
            :account="item.createBy?.account"
            :showName="true"
          />
          <span>#{{item.prIndex}}</span>
          <span style="margin-left:4px">{{t('pullRequest.createdAt')}}</span>
          <span>{{readableTimeComparingNow(item.created)}}</span>
        </div>
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
      @change="()=>listPullRequest()"
    />
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import {
  SearchOutlined,
  CloseCircleOutlined,
  CheckCircleOutlined,
  MessageOutlined,
  PlusCircleOutlined,
  PlusOutlined
} from "@ant-design/icons-vue";
import ZNoData from "@/components/common/ZNoData";
import PrStatusTag from "@/components/git/PrStatusTag";
import { ref, reactive, h } from "vue";
import { listPullRequestRequest } from "@/api/git/prApi";
import { useRoute, useRouter } from "vue-router";
import { readableTimeComparingNow } from "@/utils/time";
import { useRepoStore } from "@/pinia/repoStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const repoStore = useRepoStore();
const searchKey = ref("");
const route = useRoute();
const router = useRouter();
const repoId = parseInt(route.params.repoId);
const prStatus = ref(0);
const prList = ref([]);
// 分页
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
// 搜索
const listPullRequest = () => {
  listPullRequestRequest({
    repoId,
    status: parseInt(prStatus.value),
    pageNum: dataPage.current,
    searchKey: searchKey.value
  }).then(res => {
    prList.value = res.data;
    dataPage.totalCount = res.totalCount;
  });
};
const searchPullRequest = () => {
  dataPage.current = 1;
  listPullRequest();
};
// 选择状态并搜索
const selectPrStatus = () => {
  searchKey.value = "";
  listPullRequest();
};
// 跳转创建请求页面
const toCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/pullRequest/create`
  );
};
// 详情页
const toDetail = item => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/pullRequest/${item.prIndex}/detail`
  );
};
listPullRequest();
</script>
<style scoped>
.header {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
}
.pr-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.pr-list > li {
  padding: 18px;
}
.pr-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.pr-list > li:hover {
  cursor: pointer;
  background-color: #f0f0f0;
}
.pr-title {
  font-size: 16px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.pr-title-span {
  margin-left: 6px;
  font-weight: bold;
}
.pr-title-right {
  float: right;
  font-size: 14px;
}
.pr-desc {
  margin-top: 12px;
  color: gray;
  font-size: 12px;
}
.pr-title-right > .message-num {
  padding-left: 4px;
}
</style>