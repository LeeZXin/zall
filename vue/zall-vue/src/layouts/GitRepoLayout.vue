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
        <a-menu theme="dark" mode="inline" @click="clickPage" v-model:selectedKeys="selectedKeys">
          <a-menu-item key="/index">
            <file-outlined />
            <span>代码文件</span>
          </a-menu-item>
          <a-menu-item key="/pullRequest/list">
            <pull-request-outlined />
            <span>合并请求</span>
          </a-menu-item>
          <a-menu-item key="/branch/list">
            <branches-outlined />
            <span>分支列表</span>
          </a-menu-item>
          <a-menu-item key="/tag/list">
            <tag-outlined />
            <span>标签列表</span>
          </a-menu-item>
          <a-menu-item key="/commit/list">
            <cloud-upload-outlined />
            <span>提交历史</span>
          </a-menu-item>
          <a-menu-item
            key="/workflow/list"
            v-if="repo.perm?.canManageWorkflow ||  repo.perm?.canTriggerWorkflow"
          >
            <node-index-outlined />
            <span>工作流</span>
          </a-menu-item>
          <a-menu-item key="/protectedBranch/list" v-if="teamStore.isAdmin">
            <branches-outlined />
            <span>保护分支</span>
          </a-menu-item>
          <a-menu-item key="/webhook/list" v-if="repo.perm?.canManageWebhook">
            <api-outlined />
            <span>Webhook</span>
          </a-menu-item>
          <a-menu-item key="/opLogs" v-if="teamStore.isAdmin">
            <calendar-outlined />
            <span>操作日志</span>
          </a-menu-item>
          <a-menu-item key="/config" v-if="teamStore.isAdmin">
            <setting-outlined />
            <span>设置</span>
          </a-menu-item>
        </a-menu>
      </a-layout-sider>
      <a-layout-content>
        <div
          style="height: calc(100vh - 64px); overflow: scroll;background-color:white; width: 100%"
          ref="container"
        >
          <router-view v-if="routerActive" />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
import AvatarName from "../components/user/AvatarName";
import { useI18n } from "vue-i18n";
import { ref, provide, nextTick, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  BranchesOutlined,
  FileOutlined,
  PullRequestOutlined,
  TagOutlined,
  SettingOutlined,
  CalendarOutlined,
  CloudUploadOutlined,
  ApiOutlined,
  NodeIndexOutlined
} from "@ant-design/icons-vue";
import { getRepoRequest } from "@/api/git/repoApi";
import { getTeamRequest } from "@/api/team/teamApi";
import { useRepoStore } from "@/pinia/repoStore";
import { useTeamStore } from "@/pinia/teamStore";
const teamStore = useTeamStore();
const { t } = useI18n();
const collapsed = ref(false);
const router = useRouter();
const route = useRoute();
const repo = useRepoStore();
const selectedKeys = ref([]);
const routeKey = `/team/${route.params.teamId}/gitRepo/${route.params.repoId}`;
const routerActive = ref(false);
const container = ref(null);
// 为了子页面能体现在导航栏
const pagesMap = {
  "/index": "/index",
  "/tree": "/index",
  "/pullRequest": "/pullRequest/list",
  "/branch": "/branch/list",
  "/commit": "/commit/list",
  "/tag": "/tag/list",
  "/protectedBranch": "/protectedBranch/list",
  "/webhook": "/webhook/list",
  "/workflow": "/workflow/list",
  "/config": "/config",
  "/opLogs": "/opLogs"
};
const switchRepo = () => {
  router.push(`/team/${route.params.teamId}/gitRepo/list`);
};
const clickPage = event => {
  router.push({
    path: routeKey + event.key,
    force: true
  });
};
const getRepo = () => {
  getRepoRequest(route.params.repoId).then(res => {
    repo.repoId = res.data.repoId;
    repo.name = res.data.name;
    repo.teamId = res.data.teamId;
    repo.perm = res.data.perm;
    routerActive.value = true;
  });
};
const getTeam = () => {
  getTeamRequest(route.params.teamId).then(res => {
    teamStore.teamId = res.data.teamId;
    teamStore.name = res.data.name;
    teamStore.isAdmin = res.data.isAdmin;
    teamStore.perm = res.data.perm;
  });
};
const changeSelectedKey = path => {
  const routeSuffix = path.replace(new RegExp(`^${routeKey}`), "");
  for (let key in pagesMap) {
    let value = pagesMap[key];
    if (routeSuffix.startsWith(key)) {
      selectedKeys.value = [value];
      break;
    }
  }
};
if (teamStore.teamId === 0) {
  getTeam();
}
getRepo();
changeSelectedKey(route.path);
provide("gitRepoLayoutReload", () => {
  routerActive.value = false;
  nextTick(() => {
    routerActive.value = true;
  });
});
provide("gitRepoLayoutScrollToBottom", () => {
  if (container.value) {
    nextTick(() => {
      container.value.scrollTo({
        top: container.value.scrollHeight,
        behavior: "smooth"
      });
    });
  }
});
provide("gitRepoLayoutScrollToElem", id => {
  let c = container.value;
  let doc = document.getElementById(id);
  if (c && doc) {
    let bounding = doc.getBoundingClientRect();
    nextTick(() => {
      container.value.scrollTo({
        top: c.scrollTop + bounding.top - bounding.height,
        behavior: "smooth"
      });
    });
  }
});
watch(
  () => router.currentRoute.value.path,
  newPath => changeSelectedKey(newPath)
);
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