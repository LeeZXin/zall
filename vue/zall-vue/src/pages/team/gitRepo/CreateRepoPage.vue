<template>
  <div class="body">
    <div>
      <a-button type="primary" @click="backToRepoList">{{t("createGitRepo.backToRepoList")}}</a-button>
    </div>
    <div class="create-form">
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
          <a-input type="input" v-model:value="formState.username" />
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>{{t("createGitRepo.storeNode")}}</span>
        </div>
        <div>
          <a-select v-model:value="formState.nodeId" style="width:100%">
            <a-select-option value="jack">Jack</a-select-option>
            <a-select-option value="lucy">Lucy</a-select-option>
            <a-select-option value="disabled" disabled>Disabled</a-select-option>
            <a-select-option value="Yiminghe">yiminghe</a-select-option>
          </a-select>
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>{{t("createGitRepo.repoDesc")}}</span>
        </div>
        <div>
          <a-input type="input" v-model:value="formState.username" />
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>{{t("createGitRepo.gitignore")}}</span>
        </div>
        <div>
          <a-select v-model:value="formState.gitignore" style="width:100%">
            <a-select-option value="jack">Jack</a-select-option>
            <a-select-option value="lucy">Lucy</a-select-option>
            <a-select-option value="disabled" disabled>Disabled</a-select-option>
            <a-select-option value="Yiminghe">yiminghe</a-select-option>
          </a-select>
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
        <a-button type="primary" style="margin-top:20px;">{{t("createGitRepo.createBtn")}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import { useUserStore } from "../../../pinia/UserStore";
import { useTeamStore } from "../../../pinia/TeamStore";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
const { t } = useI18n();
const teamName = useTeamStore().teamName;
const username = useUserStore().name;
const router = useRouter();
const formState = reactive({
  username: "",
  password: "",
  gitignore: "",
  nodeId: "",
  addReadme: true,
  defaultBranch: "master"
});
const backToRepoList = () => {
  router.push("/team/gitRepo/list");
};
</script>
<style scoped>
.body {
  padding: 20px;
}
.create-form {
  padding: 30px 260px;
}
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