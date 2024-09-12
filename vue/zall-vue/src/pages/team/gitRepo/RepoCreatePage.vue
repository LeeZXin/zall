<template>
  <div style="padding:10px">
    <div class="create-form container">
      <div class="header">{{t("gitRepo.createRepo")}}</div>
      <div class="section">
        <div class="section-title">{{t("gitRepo.name")}}</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t("gitRepo.repoDesc")}}</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.desc" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">.gitignore</div>
        <div class="section-body">
          <a-select
            v-model:value="formState.gitignore"
            style="width:100%"
            :options="allGitIgnoreTemplateList.map(item=>({ value: item }))"
            show-search
            :filter-option="filterGitIgnoreTemplateListOption"
          />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t("gitRepo.defaultBranch")}}</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.defaultBranch" />
        </div>
      </div>
      <div class="section-item">
        <a-checkbox v-model:checked="formState.addReadme">
          <div>{{t("gitRepo.addReadme")}}</div>
        </a-checkbox>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="create">{{t("gitRepo.save")}}</a-button>
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
    message.error(t("gitRepo.nameFormatErr"));
    return;
  }
  if (!repoDescRegexp.test(formState.desc)) {
    message.error(t("gitRepo.repoDescFormatErr"));
    return;
  }
  if (!defaultBranchRegexp.test(formState.defaultBranch)) {
    message.error(t("gitRepo.defaultBranchFormatErr"));
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
    message.success(t("operationSuccess"));
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