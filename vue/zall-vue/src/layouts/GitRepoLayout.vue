<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{repo.name}}</span>
      <span class="switch-repo-text" @click="switchRepo">{{t("gitRepo.switchRepo")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout>
      <a-layout-sider v-model:collapsed="collapsed" collapsible>
        <a-menu theme="dark" mode="inline" @select="onselect">
          <a-menu-item key="/team/gitRepo/list">
            <file-outlined />
            <span>代码文件</span>
          </a-menu-item>
          <a-menu-item key="/team/gitRepo/pullRequestsList">
            <pull-request-outlined />
            <span>合并请求</span>
          </a-menu-item>
          <a-menu-item key="/team/gitRepo/branches">
            <branches-outlined />
            <span>分支</span>
          </a-menu-item>
          <a-menu-item key="/team/gitRepo/tags">
            <tag-outlined />
            <span>标签</span>
          </a-menu-item>
          <a-menu-item key="/team/gitRepo/opLogs">
            <calendar-outlined />
            <span>操作日志</span>
          </a-menu-item>
          <a-menu-item key="/team/gitRepo/accessToken">
            <key-outlined />
            <span>访问令牌</span>
          </a-menu-item>
          <a-menu-item key="/team/gitRepo/settings">
            <setting-outlined />
            <span>设置</span>
          </a-menu-item>
        </a-menu>
      </a-layout-sider>
      <a-layout-content style="height: calc(100vh - 64px); overflow: scroll;background-color:white">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
import AvatarName from "../components/user/AvatarName";
import { useI18n } from "vue-i18n";
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  BranchesOutlined,
  FileOutlined,
  PullRequestOutlined,
  TagOutlined,
  SettingOutlined,
  CalendarOutlined,
  KeyOutlined
} from "@ant-design/icons-vue";
import { getRepoRequest } from "@/api/git/gitApi";
import { useRepoStore } from "@/pinia/repoStore";
const { t } = useI18n();
const collapsed = ref(false);
const router = useRouter();
const route = useRoute();
const repo = useRepoStore();
const switchRepo = () => {
  router.push(`/team/${repo.teamId}/gitRepo/list`);
};
const onselect = event => {
  router.push(event.key);
};
if (repo.repoId === 0) {
  getRepoRequest({
    repoId: parseInt(route.params.repoId)
  }).then(res => {
    repo.repoId = res.data.repoId;
    repo.name = res.data.name;
    repo.teamId = res.data.teamId;
  });
}
</script>
<style scoped>
.switch-repo-text {
  color: white;
  margin-left: 12px;
  font-size: 12px;
  cursor: pointer;
}
.switch-repo-text:hover {
  color: #1677ff;
}
</style>