<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{t("app")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout-content style="height: calc(100vh - 128px)">
      <div style="padding:10px;">
        <div>
          <a-button @click="backToIndex" type="primary">选择团队</a-button>
        </div>
      </div>
      <div class="container">
        <div class="header">创建团队</div>
        <div class="form-item">
          <div class="label">
            <span>{{t("createTeam.teamName")}}</span>
          </div>
          <div>
            <a-input type="input" v-model:value="formState.teamName" />
          </div>
        </div>
        <div class="form-item">
          <a-button type="primary" style="margin-top:20px;">{{t("createGitRepo.createBtn")}}</a-button>
        </div>
      </div>
    </a-layout-content>
    <a-layout-footer style="text-align: center">{{t("footer")}}</a-layout-footer>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
import AvatarName from "../components/user/AvatarName";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import { reactive } from "vue";
const { t } = useI18n();
const router = useRouter();
const formState = reactive({
  teamName: ""
});
const backToIndex = () => {
  router.push("/index");
};
</script>
<style scoped>
.container {
  width: 50%;
  margin-left: 25%;
}
.container > .header {
  font-size: 18px;
  font-weight: bold;
  line-height: 42px;
}
</style>