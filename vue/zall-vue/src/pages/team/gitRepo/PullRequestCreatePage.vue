<template>
  <div style="padding:10px">
    <div class="header">
      <div class="title">{{t('pullRequest.compareChange')}}</div>
      <div class="desc">{{t('pullRequest.compareChangeDesc')}}</div>
    </div>
    <div class="merge-select">
      <div class="merge-select-left">
        <BranchTagSelect
          style="margin-right:8px"
          :branches="branches"
          :tags="tags"
          @select="headSelect"
        />
        <arrow-left-outlined />
        <BranchTagSelect
          style="margin-left:8px"
          :branches="branches"
          :tags="tags"
          @select="targetSelect"
        />
        <div style="margin: 0 8px" v-show="showLoading">
          <loading-outlined />
        </div>
        <div class="merge-warn" v-show="!canMerge">
          <close-outlined style="color:red" />
          <span style="padding-left:4px">{{t('pullRequest.cantMerge')}}</span>
        </div>
        <div class="merge-warn" v-show="canMerge">
          <check-outlined style="color:green" />
          <span style="padding-left:4px">{{t('pullRequest.canMerge')}}</span>
        </div>
      </div>
      <a-button
        type="primary"
        :disabled="!canMerge"
        v-show="!submitFormVisible"
        @click="showSubmitForm"
        :icon="h(PlusOutlined)"
      >{{t('pullRequest.submitPr')}}</a-button>
    </div>
    <div class="submit-form" v-if="submitFormVisible">
      <div class="form-item">
        <div class="left">{{t('pullRequest.title')}}</div>
        <div class="right">
          <a-input v-model:value="submitInput.title" />
        </div>
      </div>
      <div class="form-item">
        <div class="left">{{t('pullRequest.comment')}}</div>
        <div class="right">
          <a-textarea
            :auto-size="{ minRows: 5, maxRows: 10 }"
            v-model:value="submitInput.comment"
            show-count
            :maxlength="1024"
          />
        </div>
      </div>
      <div style="margin-top:10px;text-align:right">
        <a-button
          type="primary"
          @click="submitPr"
          :icon="h(PlusOutlined)"
        >{{t('pullRequest.submitPr')}}</a-button>
      </div>
    </div>
    <ConflictFiles v-if="conflictFiles.length > 0" :conflictFiles="conflictFiles" />
    <CommitList :commits="commits" :diffNumsStats="diffNumsStats" />
    <div>
      <FileDiffDetail
        v-for="item in fileDetails"
        :stat="item"
        :repoId="repoId"
        v-bind:key="item.filePath"
        :target="headTargetCommitId.target"
        :head="headTargetCommitId.head"
      />
    </div>
  </div>
</template>
<script setup>
import {
  ArrowLeftOutlined,
  CloseOutlined,
  CheckOutlined,
  LoadingOutlined,
  PlusOutlined
} from "@ant-design/icons-vue";
import FileDiffDetail from "@/components/git/FileDiffDetail";
import BranchTagSelect from "@/components/git/BranchTagSelect";
import { ref, reactive, nextTick, h } from "vue";
import { getBaseInfoRequest, diffRefsRequest } from "@/api/git/repoApi";
import { submitPullRequestRequest } from "@/api/git/prApi";
import { useRoute, useRouter } from "vue-router";
import { prTitleRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import CommitList from "@/components/git/CommitList";
import ConflictFiles from "@/components/git/ConflictFiles";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const repoId = parseInt(route.params.repoId);
const canMerge = ref(false);
const showLoading = ref(false);
const submitFormVisible = ref(false);
const showSubmitForm = () => {
  submitFormVisible.value = true;
};
// 冲突文件
const conflictFiles = ref([]);
// 差异内容
const fileDetails = ref([]);
// 提交
const commits = ref([]);
// 分支
const branches = ref([]);
// 标签
const tags = ref([]);
// 源
const target = reactive({
  refType: "",
  ref: ""
});
// 目标
const head = reactive({
  refType: "",
  ref: ""
});
// 标题、评论
const submitInput = reactive({
  title: "",
  comment: ""
});
// 差异信息
const diffNumsStats = reactive({
  deleteNums: 0,
  insertNums: 0,
  fileChangeNums: 0
});
// 目标/源的commitId
const headTargetCommitId = reactive({
  head: "",
  target: ""
});
// 选择源
const targetSelect = event => {
  target.refType = event.key;
  target.ref = event.value;
  diffRefsOnSelect();
};
// 选择目标
const headSelect = event => {
  head.refType = event.key;
  head.ref = event.value;
  diffRefsOnSelect();
};
const diffRefsOnSelect = () => {
  submitFormVisible.value = false;
  conflictFiles.value = [];
  if (!target.refType || !target.ref || !head.refType || !head.ref) {
    return;
  }
  if (target.refType === head.refType && target.ref === head.ref) {
    diffNumsStats.insertNums = 0;
    diffNumsStats.deleteNums = 0;
    diffNumsStats.fileChangeNums = 0;
    commits.value = [];
    canMerge.value = false;
    fileDetails.value = [];
    return;
  }
  showLoading.value = true;
  diffRefs()
    .then(res => {
      showLoading.value = false;
      canMerge.value = res.data.canMerge;
      commits.value = res.data.commits;
      fileDetails.value = [];
      // 强制重新渲染
      nextTick(() => {
        conflictFiles.value = res.data.conflictFiles;
        if (res.data.diffNumsStats) {
          diffNumsStats.insertNums = res.data.diffNumsStats.insertNums;
          diffNumsStats.deleteNums = res.data.diffNumsStats.deleteNums;
          diffNumsStats.fileChangeNums = res.data.diffNumsStats.fileChangeNums;
          fileDetails.value = res.data.diffNumsStats.stats;
        } else {
          diffNumsStats.insertNums = 0;
          diffNumsStats.deleteNums = 0;
          diffNumsStats.fileChangeNums = 0;
          fileDetails.value = [];
        }
        if (res.data.headCommit) {
          headTargetCommitId.head = res.data.headCommit.commitId;
        }
        if (res.data.targetCommit) {
          headTargetCommitId.target = res.data.targetCommit.commitId;
        }
      });
    })
    .catch(() => {
      showLoading.value = false;
    });
};
// 触发请求
const diffRefs = () => {
  return diffRefsRequest({
    repoId,
    targetType: target.refType,
    target: target.ref,
    headType: head.refType,
    head: head.ref
  });
};
// 创建合并请求
const submitPr = () => {
  // 校验标题
  if (!prTitleRegexp.test(submitInput.title)) {
    message.warn(t("pullRequest.titleFormatErr"));
    return;
  }
  submitPullRequestRequest({
    repoId,
    targetType: target.refType,
    target: target.ref,
    headType: head.refType,
    head: head.ref,
    title: submitInput.title,
    comment: submitInput.comment
  }).then(() => {
    message.success(t("operationSuccess"));
    setTimeout(() => {
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/pullRequest/list`
      );
    }, 1000);
  });
};
// 获取分支或tag列表
getBaseInfoRequest(repoId).then(res => {
  branches.value = res.data.branches;
  tags.value = res.data.tags;
});
</script>
<style scoped>
.header {
  margin-bottom: 16px;
}
.header > .title {
  font-size: 18px;
  font-weight: bold;
  padding-bottom: 8px;
}
.header > .desc {
  font-size: 14px;
  color: gray;
}
.merge-select {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  display: flex;
  align-items: center;
  padding: 10px;
  justify-content: space-between;
}
.merge-warn {
  font-size: 14px;
  margin-left: 12px;
}
.merge-select-left {
  display: flex;
  align-items: center;
}
.submit-form {
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  margin-top: 10px;
  padding: 10px;
  width: 50%;
}
.submit-form > .form-item {
  width: 100%;
  display: flex;
}
.submit-form > .form-item > .left {
  line-height: 32px;
  width: 80px;
}
.submit-form > .form-item > .right {
  width: calc(100% - 80px);
}
.submit-form > .form-item + .form-item {
  margin-top: 12px;
}
</style>