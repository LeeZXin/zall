<template>
  <div style="padding:10px">
    <div class="create-form container">
      <div class="header">{{t("createGitRepo.createText")}}</div>
      <div class="section">
        <div class="section-title">{{t("createGitRepo.repoName")}}</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.name" />
          <div class="input-desc">不包含特殊字符,长度不得超过32</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t("createGitRepo.repoDesc")}}</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.desc" />
          <div class="input-desc">为仓库添加一段简短的描述,长度不得超过255</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t("createGitRepo.gitignore")}}</div>
        <div class="section-body">
          <a-select
            v-model:value="formState.gitignore"
            style="width:100%"
            :options="allGitIgnoreTemplateList.map(item=>({ value: item }))"
            show-search
            :filter-option="filterGitIgnoreTemplateListOption"
          />
          <div class="input-desc">为仓库添加一段简短的描述,长度不得超过255</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t("createGitRepo.defaultBranch")}}</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.defaultBranch" />
          <div class="input-desc">长度不得超过32</div>
        </div>
      </div>
      <div class="form-item">
        <a-checkbox v-model:checked="formState.addReadme">
          <div>{{t("createGitRepo.addReadme")}}</div>
        </a-checkbox>
      </div>
      <div class="form-item">
        <a-button
          type="primary"
          style="margin-top:20px;"
          @click="create"
        >{{t("createGitRepo.createBtn")}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { useTeamStore } from "@/pinia/teamStore";
import { useI18n } from "vue-i18n";
import {
  allGitIgnoreTemplateListRequest,
  createRepoRequest
} from "@/api/git/repoApi";
import {
  repoNameRegexp,
  defaultBranchRegexp,
  repoDescRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useRouter, useRoute } from "vue-router";
const team = useTeamStore();
const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const formState = reactive({
  name: "",
  desc: "",
  gitignore: "",
  addReadme: true,
  defaultBranch: "main"
});
const allGitIgnoreTemplateList = ref([]);
allGitIgnoreTemplateListRequest().then(res => {
  allGitIgnoreTemplateList.value = res.data;
});
const filterGitIgnoreTemplateListOption = (input, option) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
const create = () => {
  if (!repoNameRegexp.test(formState.name)) {
    message.error("仓库名称不正确");
    return;
  }
  if (!repoDescRegexp.test(formState.desc)) {
    message.error("仓库描述不正确");
    return;
  }
  if (!defaultBranchRegexp.test(formState.defaultBranch)) {
    message.error("默认分支不正确");
    return;
  }
  createRepoRequest({
    name: formState.name,
    desc: formState.desc,
    addReadme: formState.addReadme,
    teamId: team.teamId,
    gitIgnoreName: formState.gitignore,
    defaultBranch: formState.defaultBranch
  }).then(() => {
    message.success("创建成功");
    setTimeout(() => {
      router.push(`/team/${route.params.teamId}/gitRepo/list`);
    }, 1000);
  });
};
</script>
<style scoped>
.star-text {
  font-size: 14px;
  line-height: 24px;
}
</style>