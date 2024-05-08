<template>
  <div style="padding:14px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">添加保护分支</span>
        <span v-else-if="mode === 'update'">更新保护分支</span>
        <span v-else-if="mode === 'view'">查看保护分支</span>
      </div>
      <div class="section">
        <div class="section-title">
          <span>保护分支名称模式</span>
          <span style="color:darkred">*</span>
        </div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.pattern" :readonly="mode==='view'" />
          <div class="input-desc">保护分支默认不允许删除,不允许强推, 模式例如: dev_*</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>受保护的文件模式(使用分号 ';' 分隔)</span>
        </div>
        <div class="section-body">
          <a-input
            style="width:100%"
            v-model:value="formState.protectedFilePatterns"
            :readonly="mode==='view'"
          />
          <div
            class="input-desc"
          >即使用户有权添加、编辑或删除此分支中的文件，也不允许直接更改受保护的文件。 可以使用分号 (';') 分隔多个模式。 见github.com/gobwas/glob文档了解模式语法。例如： .drone.yml, /docs/**/*.txt</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>推送</span>
        </div>
        <div class="section-body">
          <div>
            <a-radio-group v-model:value="formState.pushOption" :disabled="mode==='view'">
              <a-radio :style="radioStyle" :value="0">允许推送</a-radio>
              <div class="push-option-desc">任何拥有写访问权限的人将被允许推送到此分支(但不能强行推送)。</div>
              <a-radio :style="radioStyle" :value="1">禁止推送</a-radio>
              <div class="push-option-desc">此分支不允许推送。</div>
              <a-radio :style="radioStyle" :value="2">白名单推送</a-radio>
              <div class="push-option-desc">只有列入白名单的用户或团队才能被允许推送到此分支(但不能强行推送)。</div>
            </a-radio-group>
          </div>
          <div style="margin-top:6px" v-if="formState.pushOption === 2">
            <a-select
              :disabled="mode==='view'"
              v-model:value="formState.pushWhiteList"
              style="width:100%"
              :options="userList.map(item=>({ value: item }))"
              show-search
              mode="multiple"
              :filter-option="filterUserListOption"
            />
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>合并请求</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <div class="input-title">足够的人审核才能合并请求</div>
            <a-input-number
              style="width: 100%"
              :min="0"
              :readonly="mode==='view'"
              v-model:value="formState.reviewCountWhenCreatePr"
            />
          </div>
          <div class="input-item">
            <div class="input-title">
              <div>审批人白名单</div>
              <div
                style="font-size:13px;color:gray;padding-top:4px"
              >只有白名单用户或团队的审核才能计数。 没有批准的白名单，来自任何有写访问权限的人的审核都将计数。</div>
            </div>
            <a-select
              :disabled="mode==='view'"
              v-model:value="formState.reviewerList"
              style="width:100%"
              :options="userList.map(item=>({ value: item }))"
              show-search
              mode="multiple"
              :filter-option="filterUserListOption"
            />
          </div>
          <div style="margin: 14px 0">
            <a-checkbox
              :disabled="mode==='view'"
              v-model:checked="formState.cancelOldReviewApprovalWhenNewCommit"
              style="font-size:14px"
            >当新的提交更改合并请求内容被推送到分支时，旧的批准将被撤销。</a-checkbox>
          </div>
        </div>
      </div>
      <div style="width:100%;border-top:1px solid #d9d9d9;margin: 10px 0" v-if="mode !== 'view'"></div>
      <div style="margin-bottom:20px" v-if="mode !== 'view'">
        <a-button type="primary" @click="createProtectedBranch">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { listAccountRequest } from "@/api/team/teamApi";
import {
  createProtectedBranchRequest,
  updateProtecteddBranchRequest
} from "@/api/git/branchApi";
import { useRoute, useRouter } from "vue-router";
import { protectedBranchPatternRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useProtectedBranchStore } from "@/pinia/protectedBranchStore";
import { useRepoStore } from "@/pinia/repoStore";
const route = useRoute();
const repoStore = useRepoStore();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const protectedBranchStore = useProtectedBranchStore();
const userList = ref([]);
const router = useRouter();
const mode = getMode();
const formState = reactive({
  pattern: "",
  protectedFilePatterns: "",
  pushOption: 0,
  pushWhiteList: [],
  reviewCountWhenCreatePr: 1,
  reviewerList: [],
  cancelOldReviewApprovalWhenNewCommit: true
});
const radioStyle = reactive({
  display: "flex",
  alignItems: "flex-start"
});
const filterUserListOption = (input, option) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
const createProtectedBranch = () => {
  if (!protectedBranchPatternRegexp.test(formState.pattern)) {
    message.warn("分支模式错误");
    return;
  }
  if (formState.reviewerList.length > 0 && formState.reviewCountWhenCreatePr > formState.reviewerList.length) {
    message.warn("当限制了白名单, 审批人数量不得大于白名单数量");
    return;
  }
  if (mode === "create") {
    createProtectedBranchRequest({
      pattern: formState.pattern,
      repoId: parseInt(route.params.repoId),
      cfg: {
        protectedFilePatterns: formState.protectedFilePatterns.split(";"),
        pushOption: formState.pushOption,
        pushWhiteList: formState.pushWhiteList,
        reviewCountWhenCreatePr: formState.reviewCountWhenCreatePr,
        reviewerList: formState.reviewerList,
        cancelOldReviewApprovalWhenNewCommit:
          formState.cancelOldReviewApprovalWhenNewCommit
      }
    }).then(() => {
      message.success("添加成功");
      router.push(`/gitRepo/${route.params.repoId}/protectedBranch/list`);
    });
  } else if (mode === "update") {
    updateProtecteddBranchRequest({
      protectedBranchId: protectedBranchStore.id,
      pattern: formState.pattern,
      repoId: parseInt(route.params.repoId),
      cfg: {
        protectedFilePatterns: formState.protectedFilePatterns.split(";"),
        pushOption: formState.pushOption,
        pushWhiteList: formState.pushWhiteList,
        reviewCountWhenCreatePr: formState.reviewCountWhenCreatePr,
        reviewerList: formState.reviewerList,
        cancelOldReviewApprovalWhenNewCommit:
          formState.cancelOldReviewApprovalWhenNewCommit
      }
    }).then(() => {
      message.success("更新成功");
      router.push(`/gitRepo/${route.params.repoId}/protectedBranch/list`);
    });
  }
};
if (mode !== "create") {
  if (
    protectedBranchStore.id === 0 ||
    parseInt(route.params.protectedBranchId) !== protectedBranchStore.id
  ) {
    router.push(`/gitRepo/${route.params.repoId}/protectedBranch/list`);
  } else {
    if (mode !== "create") {
      formState.pattern = protectedBranchStore.pattern;
      formState.protectedFilePatterns = protectedBranchStore.cfg?.protectedFilePatterns?.join(
        ";"
      );
      formState.pushOption = protectedBranchStore.cfg?.pushOption;
      formState.pushWhiteList = protectedBranchStore.cfg?.pushWhiteList;
      formState.reviewCountWhenCreatePr =
        protectedBranchStore.cfg?.reviewCountWhenCreatePr;
      formState.reviewerList = protectedBranchStore.cfg?.reviewerList;
      formState.cancelOldReviewApprovalWhenNewCommit =
        protectedBranchStore.cfg?.cancelOldReviewApprovalWhenNewCommit;
    }
    if (mode !== "view") {
      listAccountRequest(repoStore.teamId).then(res => {
        userList.value = res.data;
      });
    }
  }
} else {
  listAccountRequest(repoStore.teamId).then(res => {
    userList.value = res.data;
  });
}
</script>
<style scoped>
.header {
  font-size: 18px;
  margin-bottom: 10px;
  font-weight: bold;
}
.input-desc {
  font-size: 13px;
  color: gray;
  margin: 10px 0;
  line-height: 20px;
  word-break: break-all;
}
.push-option {
  font-size: 14px;
}
.push-option-desc {
  font-size: 13px;
  padding-left: 24px;
  color: gray;
  margin-bottom: 6px;
}
.input-title {
  font-size: 14px;
  margin-bottom: 6px;
}
.input-item + .input-item {
  margin-top: 16px;
}
</style>