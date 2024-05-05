<template>
  <div style="padding:14px">
    <ZNaviBack :name="t('createGitRepo.backToRepoList')" :url="`/team/${route.params.teamId}/gitRepo/list`" />
    <div class="create-form container">
      <div class="title">{{t("createGitRepo.createText")}}</div>
      <!--div class="migrate-text">仓库包含所有项目文件，包括修订历史。已经在别处有了吗？</div-->
      <div class="star-text">{{t("createGitRepo.starText")}}</div>
      <div style="display:flex;align-items:center">
        <div class="form-item">
          <div class="label">
            <span>{{t("createGitRepo.owner")}}</span>
          </div>
          <div class="form-item-text">{{username}}</div>
        </div>
        <div class="form-item">
          <div class="label">
            <span>{{t("createGitRepo.team")}}</span>
          </div>
          <div class="form-item-text">{{teamName}}</div>
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>{{t("createGitRepo.repoName")}}</span>
        </div>
        <div>
          <a-input type="input" v-model:value="formState.name" />
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>{{t("createGitRepo.repoDesc")}}</span>
        </div>
        <div>
          <a-input type="input" v-model:value="formState.desc" />
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>{{t("createGitRepo.gitignore")}}</span>
        </div>
        <div>
          <a-select
            v-model:value="formState.gitignore"
            style="width:100%"
            :options="allGitIgnoreTemplateList.map(item=>({ value: item }))"
            show-search
            :filter-option="filterGitIgnoreTemplateListOption"
          />
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>{{t("createGitRepo.defaultBranch")}}</span>
        </div>
        <div>
          <a-input type="input" v-model:value="formState.defaultBranch" />
        </div>
      </div>
      <div class="form-item">
        <a-checkbox v-model:checked="formState.addReadme" style="margin-top:6px">
          <div class="add-readme-text">{{t("createGitRepo.addReadme")}}</div>
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
import ZNaviBack from "@/components/common/ZNaviBack";
import { reactive, ref } from "vue";
import { useUserStore } from "@/pinia/userStore";
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
const teamName = useTeamStore().name;
const username = useUserStore().name;
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
.create-form > .title {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 10px;
  line-height: 32px;
}
.migrate-text {
  font-size: 14px;
  color: gray;
  line-height: 24px;
}
.star-text {
  font-size: 14px;
  line-height: 24px;
}
.add-readme-text {
  font-size: 14px;
}
</style>