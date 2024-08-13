<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-radio-group v-model:value="prStatus" @change="selectPrStatus">
        <a-radio-button value="0">
          <span>所有</span>
          <span>({{stats.totalCount}})</span>
        </a-radio-button>
        <a-radio-button value="1">
          <span>已打开</span>
          <span>({{stats.openCount}})</span>
        </a-radio-button>
        <a-radio-button value="3">
          <span>已合并</span>
          <span>({{stats.mergedCount}})</span>
        </a-radio-button>
        <a-radio-button value="2">
          <span>已关闭</span>
          <span>({{stats.closedCount}})</span>
        </a-radio-button>
      </a-radio-group>
    </div>
    <div class="header">
      <a-input placeholder="搜索合并请求" v-model:value="searchKey" @pressEnter="()=>listPullRequest()">
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
      >创建合并请求</a-button>
    </div>
    <ul class="pr-list" v-show="prList.length > 0">
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
        <div class="pr-desc">
          <PrStatusTag :status="item.prStatus" />
          <span>#{{item.id}}</span>
          <span>created by {{item.createBy}}</span>
          <span>{{readableTimeComparingNow(item.created)}}</span>
        </div>
      </li>
    </ul>
    <ZNoData v-show="prList.length === 0" />
    <a-pagination
      v-model:current="currPage"
      :total="totalCount"
      show-less-items
      :pageSize="pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listPullRequest()"
    />
  </div>
</template>
<script setup>
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
import {
  listPullRequestRequest,
  statsPullRequestRequest
} from "@/api/git/prApi";
import { useRoute, useRouter } from "vue-router";
import { readableTimeComparingNow } from "@/utils/time";
import { useRepoStore } from "@/pinia/repoStore";
const repoStore = useRepoStore();
const searchKey = ref("");
const route = useRoute();
const router = useRouter();
const repoId = parseInt(route.params.repoId);
const prStatus = ref("0");
const totalCount = ref(0);
const currPage = ref(1);
const prList = ref([]);
const pageSize = 10;
const stats = reactive({
  totalCount: 0,
  openCount: 0,
  closedCount: 0,
  mergedCount: 0
});
// 搜索
const listPullRequest = () => {
  listPullRequestRequest({
    repoId,
    status: parseInt(prStatus.value),
    pageNum: currPage.value,
    searchKey: searchKey.value,
    pageSize: pageSize
  }).then(res => {
    prList.value = res.data;
    totalCount.value = res.totalCount;
  });
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
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/pullRequest/${item.id}/detail`
  );
};
listPullRequest();
statsPullRequestRequest(repoId).then(res => {
  stats.totalCount = res.data.totalCount;
  stats.openCount = res.data.openCount;
  stats.closedCount = res.data.closedCount;
  stats.mergedCount = res.data.mergedCount;
});
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
  padding: 12px;
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
  margin-top: 6px;
  color: gray;
  font-size: 12px;
}
.pr-desc > span + span {
  margin-left: 4px;
}
.pr-title-right > .message-num {
  padding-left: 4px;
}
</style>