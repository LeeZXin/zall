<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('protectedBranch.createBranch')}}</span>
        <span v-else-if="mode === 'update'">{{t('protectedBranch.updateBranch')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('protectedBranch.pattern')}}</div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.pattern" />
          <div class="input-desc">{{t('protectedBranch.patternDesc')}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('protectedBranch.push')}}</div>
        <div class="section-body">
          <div>
            <a-radio-group v-model:value="formState.pushOption">
              <a-radio :style="radioStyle" :value="0">{{t('protectedBranch.allowPush')}}</a-radio>
              <div class="radio-option-desc">{{t('protectedBranch.allowPushDesc')}}</div>
              <a-radio :style="radioStyle" :value="1">{{t('protectedBranch.disallowPush')}}</a-radio>
              <div class="radio-option-desc">{{t('protectedBranch.disallowPushDesc')}}</div>
              <a-radio :style="radioStyle" :value="2">{{t('protectedBranch.whiteListPush')}}</a-radio>
              <div class="radio-option-desc">{{t('protectedBranch.whiteListPushDesc')}}</div>
            </a-radio-group>
          </div>
          <div style="margin-top:6px" v-if="formState.pushOption === 2">
            <a-select
              v-model:value="formState.pushWhiteList"
              style="width:100%"
              :options="userList"
              show-search
              mode="multiple"
              :filter-option="filterUserListOption"
            />
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('protectedBranch.pullRequest')}}</div>
        <div class="section-body">
          <div class="input-item">
            <div class="input-title">{{t('protectedBranch.mergePullRequestDesc')}}</div>
            <a-input-number
              style="width: 100%"
              :min="0"
              v-model:value="formState.reviewCountWhenCreatePr"
            />
          </div>
          <div class="input-item">
            <div class="input-title">{{t('protectedBranch.reviewPullRequestWhiteList')}}</div>
            <a-select
              v-model:value="formState.reviewerList"
              style="width:100%"
              :options="userList"
              show-search
              mode="multiple"
              :filter-option="filterUserListOption"
            />
            <div class="input-desc">{{t('protectedBranch.reviewPullRequestWhiteListDesc')}}</div>
          </div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button
          type="primary"
          @click="createOrUpdateProtectedBranch"
        >{{t('protectedBranch.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { listUserByTeamIdRequest } from "@/api/team/teamApi";
import {
  createProtectedBranchRequest,
  updateProtectedBranchRequest
} from "@/api/git/branchApi";
import { useRoute, useRouter } from "vue-router";
import { protectedBranchPatternRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useProtectedBranchStore } from "@/pinia/protectedBranchStore";
import { useRepoStore } from "@/pinia/repoStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const route = useRoute();
const repoStore = useRepoStore();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const protectedBranchStore = useProtectedBranchStore();
// 团队成员列表
const userList = ref([]);
const router = useRouter();
const mode = getMode();
// 表单数据
const formState = reactive({
  pattern: "",
  pushOption: 0,
  pushWhiteList: [],
  reviewCountWhenCreatePr: 1,
  reviewerList: []
});
const radioStyle = reactive({
  display: "flex",
  alignItems: "flex-start"
});
// 成员列表下拉框过滤
const filterUserListOption = (input, option) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 获取团队成员
const listUser = () => {
  listUserByTeamIdRequest(repoStore.teamId).then(res => {
    userList.value = res.data.map(item => {
      return {
        value: item.account,
        label: `${item.account}(${item.name})`
      };
    });
  });
};
// 新增或编辑保护分支
const createOrUpdateProtectedBranch = () => {
  if (!protectedBranchPatternRegexp.test(formState.pattern)) {
    message.warn(t("protectedBranch.patternFormatErr"));
    return;
  }
  if (
    formState.reviewerList.length > 0 &&
    formState.reviewCountWhenCreatePr > formState.reviewerList.length
  ) {
    message.warn(t("protectedBranch.assignedReviewerCountErr"));
    return;
  }
  if (mode === "create") {
    createProtectedBranchRequest({
      pattern: formState.pattern,
      repoId: parseInt(route.params.repoId),
      cfg: {
        pushOption: formState.pushOption,
        pushWhiteList: formState.pushWhiteList,
        reviewCountWhenCreatePr: formState.reviewCountWhenCreatePr,
        reviewerList: formState.reviewerList
      }
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/protectedBranch/list`
      );
    });
  } else if (mode === "update") {
    updateProtectedBranchRequest({
      protectedBranchId: protectedBranchStore.id,
      pattern: formState.pattern,
      repoId: parseInt(route.params.repoId),
      cfg: {
        pushOption: formState.pushOption,
        pushWhiteList: formState.pushWhiteList,
        reviewCountWhenCreatePr: formState.reviewCountWhenCreatePr,
        reviewerList: formState.reviewerList
      }
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/protectedBranch/list`
      );
    });
  }
};
if (mode !== "create") {
  if (
    protectedBranchStore.id === 0 ||
    parseInt(route.params.protectedBranchId) !== protectedBranchStore.id
  ) {
    router.push(
      `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/protectedBranch/list`
    );
  } else {
    if (mode !== "create") {
      formState.pattern = protectedBranchStore.pattern;
      formState.pushOption = protectedBranchStore.cfg?.pushOption;
      formState.pushWhiteList = protectedBranchStore.cfg?.pushWhiteList;
      formState.reviewCountWhenCreatePr =
        protectedBranchStore.cfg?.reviewCountWhenCreatePr;
      formState.reviewerList = protectedBranchStore.cfg?.reviewerList;
    }
    listUser();
  }
} else {
  listUser();
}
</script>
<style scoped>
</style>