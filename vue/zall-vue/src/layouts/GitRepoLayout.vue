<template>
  <a-watermark :content="`${userStore.name}${userStore.account}`" :gap="[200,200]">
    <a-layout>
      <a-layout-header style="font-size:22px;color:white">
        <span>{{repo.name}}</span>
        <span class="switch-repo-text" @click="switchRepo">{{t("gitRepo.switchRepo")}}</span>
        <AvatarName style="float:right;" />
        <I18nSelect style="float:right;margin-right: 20px" />
      </a-layout-header>
      <a-layout>
        <a-layout-sider v-model:collapsed="collapsed" collapsible>
          <a-menu theme="dark" mode="inline" @click="selectKey" v-model:selectedKeys="selectedKeys">
            <a-menu-item key="/index">
              <FileOutlined />
              <span>{{t("gitRepoMenu.index")}}</span>
            </a-menu-item>
            <a-menu-item key="/pullRequest/list">
              <PullRequestOutlined />
              <span>{{t("gitRepoMenu.pullRequest")}}</span>
            </a-menu-item>
            <a-menu-item key="/branch/list">
              <BranchesOutlined />
              <span>{{t("gitRepoMenu.branch")}}</span>
            </a-menu-item>
            <a-menu-item key="/tag/list">
              <TagOutlined />
              <span>{{t("gitRepoMenu.tag")}}</span>
            </a-menu-item>
            <a-menu-item key="/commit/list">
              <CloudUploadOutlined />
              <span>{{t("gitRepoMenu.commitHistory")}}</span>
            </a-menu-item>
            <a-menu-item key="/workflow/list">
              <NodeIndexOutlined />
              <span>{{t("gitRepoMenu.workflow")}}</span>
            </a-menu-item>
            <a-menu-item key="/protectedBranch/list" v-if="repo.perm?.canManageProtectedBranch">
              <BranchesOutlined />
              <span>{{t("gitRepoMenu.protectedBranch")}}</span>
            </a-menu-item>
            <a-menu-item key="/webhook/list" v-if="repo.perm?.canManageWebhook">
              <ApiOutlined />
              <span>{{t("gitRepoMenu.webhook")}}</span>
            </a-menu-item>
            <a-menu-item key="/setting" v-if="teamStore.isAdmin">
              <SettingOutlined />
              <span>{{t("gitRepoMenu.setting")}}</span>
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
  </a-watermark>
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
  CloudUploadOutlined,
  ApiOutlined,
  NodeIndexOutlined
} from "@ant-design/icons-vue";
import { getRepoRequest } from "@/api/git/repoApi";
import { getTeamRequest } from "@/api/team/teamApi";
import { useRepoStore } from "@/pinia/repoStore";
import { useTeamStore } from "@/pinia/teamStore";
import { useUserStore } from "@/pinia/userStore";
const userStore = useUserStore();
const teamStore = useTeamStore();
const { t } = useI18n();
// 导航栏是否合上
const collapsed = ref(false);
const router = useRouter();
const route = useRoute();
const repo = useRepoStore();
// 导航栏选择key
const selectedKeys = ref([]);
// 路由前缀
const routeKey = `/team/${route.params.teamId}/gitRepo/${route.params.repoId}`;
// 是否展示routerView
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
  "/setting": "/setting"
};
// 切换仓库
const switchRepo = () => {
  router.push(`/team/${route.params.teamId}/gitRepo/list`);
};
// 导航栏点击
const selectKey = event => {
  router.push({
    path: routeKey + event.key,
    force: true
  });
};
// 获取仓库信息和权限
const getRepo = () => {
  getRepoRequest(route.params.repoId).then(res => {
    repo.repoId = res.data.repoId;
    repo.name = res.data.name;
    repo.teamId = res.data.teamId;
    repo.perm = res.data.perm;
    routerActive.value = true;
  });
};
// 获取团队信息和权限
const getTeam = () => {
  getTeamRequest(route.params.teamId).then(res => {
    teamStore.teamId = res.data.teamId;
    teamStore.name = res.data.name;
    teamStore.isAdmin = res.data.isAdmin;
    teamStore.perm = res.data.perm;
  });
};
// 导航栏key变化触发
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
// 重载页面
provide("gitRepoLayoutReload", () => {
  routerActive.value = false;
  nextTick(() => {
    routerActive.value = true;
  });
});
// 滚动到底部
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
// 滚动到指定位置
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