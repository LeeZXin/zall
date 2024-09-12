<template>
  <div style="padding:10px">
    <template v-if="branches.length > 0">
      <div style="margin-bottom: 10px">
        <BranchTagSelect
          :disableTags="true"
          :branches="branches"
          :defaultBranch="route.params.ref"
          @select="onBranchSelect"
        />
      </div>
      <ul class="commit-list" v-if="commits.length > 0">
        <li v-for="(item, index) in commits" v-bind:key="index">
          <div style="width:70%">
            <div class="commit-msg no-wrap">
              <span class="commit-msg-text" @click="treeCommit(item)">{{item.commitMsg}}</span>
            </div>
            <div class="commit-desc no-wrap flex-center">
              <ZAvatar :url="item.committer?.avatarUrl" :name="item.committer?.name" :showName="true"/>
              <span>{{t('commitList.committedAt')}}</span>
              <span>{{readableTimeComparingNow(item.committedTime)}}</span>
            </div>
          </div>
          <div>
            <a-popover v-if="item.verified" placement="bottomRight">
              <template #content>
                <div style="width: 300px;font-size:14px;padding:6px">
                  <div style="margin-bottom: 12px;" class="flex-center no-wrap">
                    <CheckCircleFilled style="color:green;margin-right:10px" />
                    <span>{{t('commitList.thisCommitIsVerified')}}</span>
                  </div>
                  <div class="flex-center" style="margin-bottom: 12px;">
                    <ZAvatar :url="item.signer?.avatarUrl" :name="item.signer?.name" size="medium" />
                    <div style="margin-left:8px">
                      <div style="margin-bottom: 3px" class="no-wrap">{{item.signer?.account}}</div>
                      <div class="no-wrap">{{item.signer?.name}}</div>
                    </div>
                  </div>
                  <div class="no-wrap">{{item.signer?.type}} KEY</div>
                  <div style="color:gray;word-break:break-all">{{item.signer?.key}}</div>
                </div>
              </template>
              <span style="cursor:pointer">
                <a-tag color="green">{{t('commitList.verified')}}</a-tag>
              </span>
            </a-popover>
            <CommitSha>{{item.shortId}}</CommitSha>
          </div>
        </li>
        <li v-if="lastLoadCount >= 10">
          <div class="more-btn" @click="getCommits()">{{t('commitList.loadMore')}}...</div>
        </li>
      </ul>
    </template>
    <ZNoData v-else />
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import ZNoData from "@/components/common/ZNoData";
import CommitSha from "@/components/git/CommitSha";
import BranchTagSelect from "@/components/git/BranchTagSelect";
import { allBranchesRequest, historyCommitsRequest } from "@/api/git/repoApi";
import { ref, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import { readableTimeComparingNow } from "@/utils/time";
import { CheckCircleFilled } from "@ant-design/icons-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const branches = ref([]);
const commits = ref([]);
const lastLoadCount = ref(0);
const selectedBranch = ref("");
// 获取分支
const getBranches = () => {
  allBranchesRequest(route.params.repoId).then(res => {
    branches.value = res.data;
  });
};
// 获取列表
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
// 分支选择
const onBranchSelect = ({ value }) => {
  router.replace(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/commit/list/${value}`
  );
  selectedBranch.value = value;
  commits.value = [];
  nextTick(() => {
    getCommits();
  });
};
// 跳转提交详情页
const treeCommit = item => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/commit/diff/${item.commitId}`
  );
};
getBranches();
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
  font-size: 14px;
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
.more-btn {
  width: 100%;
  text-align: center;
  cursor: pointer;
}
.more-btn:hover {
  color: #1677ff;
}
</style>