<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode==='create'">{{t('teamRole.createRole')}}</span>
        <span v-else-if="mode==='update'">{{t('teamRole.updateRole')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('teamRole.name')}}</div>
        <div class="section-body">
          <div class="input-item">
            <a-input v-model:value="formState.name" />
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('teamRole.teamPerm')}}</div>
        <div class="section-body">
          <ul class="perm-list">
            <li v-for="item in teamPermList" v-bind:key="item.key">
              <a-checkbox v-model:checked="formState.teamPerm[item.key]">{{t(item.checkbox)}}</a-checkbox>
              <div class="checkbox-desc">{{t(item.desc)}}</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('teamRole.repoPermOption')}}</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.useDefaultRepoPerm">
            <a-radio :value="true">{{t('teamRole.useDefaultRepoPerm')}}</a-radio>
            <div class="radio-option-desc">{{t('teamRole.useDefaultRepoPermDesc')}}</div>
            <a-radio :value="false">{{t('teamRole.useSpecificRepoPerm')}}</a-radio>
            <div class="radio-option-desc">{{t('teamRole.useSpecificRepoPermDesc')}}</div>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.useDefaultRepoPerm">
        <div class="section-title">{{t('teamRole.defaultRepoPerm')}}</div>
        <div class="section-body">
          <ul class="perm-ul">
            <li v-for="item in repoPermList" v-bind:key="item.title">
              <div class="perm-title">{{t(item.title)}}</div>
              <ul class="perm-list">
                <li v-for="(perm, index) in item.perms" v-bind:key="index">
                  <a-checkbox
                    v-model:checked="formState.defaultRepoPerm[perm.key]"
                  >{{t(perm.checkbox)}}</a-checkbox>
                  <div class="checkbox-desc">{{t(perm.desc)}}</div>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-else>
        <div class="section-title">{{t('teamRole.specificRepoPerm')}}</div>
        <ul class="repo-ul" v-if="formState.addRepoPermList.length > 0">
          <li
            v-for="item in formState.addRepoPermList"
            v-bind:key="item.repoId"
            @click="showUpdateRepoPermModal(item)"
          >
            <span>{{item.name}}</span>
            <span style="color:green">+{{item.permCount}} {{t('teamRole.permItem')}}</span>
          </li>
        </ul>
        <div
          class="add-repo-btn"
          @click="showAddRepoPermModal"
          v-if="repoModalState.remainList.length > 0"
        >
          <PlusOutlined />
          <span style="margin-left:4px">{{t('teamRole.addRepoPerm')}}</span>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('teamRole.appPermOption')}}</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.useDefaultAppPerm">
            <a-radio :value="true">{{t('teamRole.useDefaultAppPerm')}}</a-radio>
            <div class="radio-option-desc">{{t('teamRole.useDefaultAppPermDesc')}}</div>
            <a-radio :value="false">{{t('teamRole.useSpecificAppPerm')}}</a-radio>
            <div class="radio-option-desc">{{t('teamRole.useSpecificAppPermDesc')}}</div>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.useDefaultAppPerm">
        <div class="section-title">{{t('teamRole.defaultAppPerm')}}</div>
        <div class="section-body">
          <ul class="perm-list">
            <li v-for="item in appPermList" v-bind:key="item.key">
              <a-checkbox v-model:checked="formState.defaultAppPerm[item.key]">{{t(item.checkbox)}}</a-checkbox>
              <div class="checkbox-desc">{{t(item.desc)}}</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-else>
        <div class="section-title">{{t('teamRole.specificAppPerm')}}</div>
        <ul class="repo-ul" v-if="formState.addAppPermList.length > 0">
          <li
            v-for="item in formState.addAppPermList"
            v-bind:key="item.appId"
            @click="showUpdateAppPermModal(item)"
          >
            <span>{{item.name}}</span>
            <span style="color:green">+{{item.permCount}} {{t('teamRole.permItem')}}</span>
          </li>
        </ul>
        <div
          class="add-repo-btn"
          @click="showAddAppPermModal"
          v-if="appModalState.remainList.length > 0"
        >
          <PlusOutlined />
          <span style="margin-left:4px">{{t('teamRole.addAppPerm')}}</span>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createOrUpdateRole">{{t('teamRole.save')}}</a-button>
      </div>
    </div>
    <a-modal
      v-model:open="repoModalState.addModalOpen"
      :title="t('teamRole.addRepoPerm')"
      @ok="handleAddRepoModalOk"
    >
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">{{t('teamRole.selectMultipleRepo')}}</div>
        <a-select
          style="width:100%"
          v-model:value="repoModalState.repoSelect"
          :options="repoModalState.remainList.map(item=>({ value: item.repoId, label: item.name }))"
          show-search
          mode="multiple"
          :filter-option="filterSelectOption"
        />
      </div>
      <ul class="perm-ul">
        <li v-for="item in repoPermList" v-bind:key="item.title">
          <div class="perm-title">{{t(item.title)}}</div>
          <ul class="perm-list">
            <li v-for="(perm, index) in item.perms" v-bind:key="index">
              <a-checkbox
                v-model:checked="repoModalState.addModalCheckboxs[perm.key]"
              >{{t(perm.checkbox)}}</a-checkbox>
              <div class="checkbox-desc">{{t(perm.desc)}}</div>
            </li>
          </ul>
        </li>
      </ul>
    </a-modal>
    <a-modal v-model:open="repoModalState.updateModalOpen" :title="t('teamRole.updateRepoPerm')">
      <template #footer>
        <a-button type="primary" @click="deleteAddRepoPerm" danger>{{t('teamRole.delete')}}</a-button>
        <a-button type="primary" @click="handleUpdateRepoModalOk">{{t('teamRole.ok')}}</a-button>
      </template>
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">{{t('teamRole.repo')}}</div>
        <div style="font-size:16px">{{repoModalState.updatePerm.name}}</div>
      </div>
      <ul class="perm-ul">
        <li v-for="item in repoPermList" v-bind:key="item.title">
          <div class="perm-title">{{t(item.title)}}</div>
          <ul class="perm-list">
            <li v-for="(perm, index) in item.perms" v-bind:key="index">
              <a-checkbox
                v-model:checked="repoModalState.updatePerm.perms[perm.key]"
              >{{t(perm.checkbox)}}</a-checkbox>
              <div class="checkbox-desc">{{t(perm.desc)}}</div>
            </li>
          </ul>
        </li>
      </ul>
    </a-modal>
    <a-modal
      v-model:open="appModalState.addModalOpen"
      :title="t('teamRole.addAppPerm')"
      @ok="handleAddAppModalOk"
    >
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">{{t('teamRole.selectMultipleApp')}}</div>
        <a-select
          style="width:100%"
          v-model:value="appModalState.appSelect"
          :options="appModalState.remainList.map(item=>({ value: item.appId, label: item.name }))"
          show-search
          mode="multiple"
          :filter-option="filterSelectOption"
        />
      </div>
      <ul class="perm-list">
        <li v-for="item in appPermList" v-bind:key="item.key">
          <a-checkbox v-model:checked="appModalState.addModalCheckboxs[item.key]">{{t(item.checkbox)}}</a-checkbox>
          <div class="checkbox-desc">{{t(item.desc)}}</div>
        </li>
      </ul>
    </a-modal>
    <a-modal v-model:open="appModalState.updateModalOpen" :title="t('teamRole.updateAppPerm')">
      <template #footer>
        <a-button type="primary" @click="deleteAddAppPerm" danger>{{t('teamRole.delete')}}</a-button>
        <a-button type="primary" @click="handleUpdateAppModalOk">{{t('teamRole.ok')}}</a-button>
      </template>
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">{{t('teamRole.app')}}</div>
        <div style="font-size:16px">{{appModalState.updatePerm.name}}</div>
      </div>
      <ul class="perm-list">
        <li v-for="item in appPermList" v-bind:key="item.key">
          <a-checkbox v-model:checked="appModalState.updatePerm.perms[item.key]">{{t(item.checkbox)}}</a-checkbox>
          <div class="checkbox-desc">{{t(item.desc)}}</div>
        </li>
      </ul>
    </a-modal>
  </div>
</template>
<script setup>
import { PlusOutlined } from "@ant-design/icons-vue";
import { reactive, nextTick } from "vue";
import { message } from "ant-design-vue";
import { getRepoListByAdminRequest } from "@/api/git/repoApi";
import { createRoleRequest, updateRoleRequest } from "@/api/team/teamApi";
import { listAllAppByAdminRequest } from "@/api/app/appApi";
import { useRoute, useRouter } from "vue-router";
import { teamRoleNameRegexp } from "@/utils/regexp";
import { useTeamRoleStore } from "@/pinia/teamRoleStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const teamRoleStore = useTeamRoleStore();
const mode = getMode();
const router = useRouter();
// 表单数据
const formState = reactive({
  name: "",
  teamPerm: {},
  // 默认仓库权限
  defaultRepoPerm: {
    canAccessRepo: true,
    canPushRepo: true,
    canSubmitPullRequest: true,
    canAddCommentInPullRequest: true,
    canTriggerWorkflow: true
  },
  // 默认应用服务权限
  defaultAppPerm: {
    canDevelop: true
  },
  useDefaultRepoPerm: true,
  useDefaultAppPerm: true,
  addRepoPermList: [],
  addAppPermList: []
});
// 仓库modal
const repoModalState = reactive({
  allList: [],
  remainList: [],
  repoSelect: [],
  addModalCheckboxs: {},
  addModalOpen: false,
  updateModalOpen: false,
  updatePerm: {
    repoId: 0,
    name: "",
    perms: {},
    // 用于回显
    target: null
  }
});
// 应用服务modal
const appModalState = reactive({
  allList: [],
  remainList: [],
  appSelect: [],
  addModalCheckboxs: {},
  addModalOpen: false,
  updateModalOpen: false,
  updatePerm: {
    appId: "",
    name: "",
    perms: {},
    // 用于回显
    target: null
  }
});
// 应用服务权限列表
const appPermList = [
  {
    checkbox: "teamRole.canDevelopApp",
    key: "canDevelop",
    desc: "teamRole.canDevelopAppDesc"
  },
  {
    checkbox: "teamRole.canCreateDeployPlan",
    key: "canCreateDeployPlan",
    desc: "teamRole.canCreateDeployPlanDesc"
  },
  {
    checkbox: "teamRole.canManagePipeline",
    key: "canManagePipeline",
    desc: "teamRole.canManagePipelineDesc"
  },
  {
    checkbox: "teamRole.canManageDiscoverySource",
    key: "canManageDiscoverySource",
    desc: "teamRole.canManageDiscoverySourceDesc"
  },
  {
    checkbox: "teamRole.canManagePropertySource",
    key: "canManagePropertySource",
    desc: "teamRole.canManagePropertySourceDesc"
  },
  {
    checkbox: "teamRole.canDeployProperty",
    key: "canDeployProperty",
    desc: "teamRole.canDeployPropertyDesc"
  },
  {
    checkbox: "teamRole.canManageServiceSource",
    key: "canManageServiceSource",
    desc: "teamRole.canManageServiceSourceDesc"
  }
];
const teamPermList = [
  {
    checkbox: "teamRole.canManageTimer",
    key: "canManageTimer",
    desc: "teamRole.canManageTimerDesc"
  },
  {
    checkbox: "teamRole.canManageNotifyTpl",
    key: "canManageNotifyTpl",
    desc: "teamRole.canManageNotifyTplDesc"
  },
  {
    checkbox: "teamRole.canManageTeamHook",
    key: "canManageTeamHook",
    desc: "teamRole.canManageTeamHookDesc"
  },
  {
    checkbox: "teamRole.canManagePromScrape",
    key: "canManagePromScrape",
    desc: "teamRole.canManagePromScrapeDesc"
  },
  {
    checkbox: "teamRole.canManageWeworkAccessToken",
    key: "canManageWeworkAccessToken",
    desc: "teamRole.canManageWeworkAccessTokenDesc"
  },
  {
    checkbox: "teamRole.canManageFeishuAccessToken",
    key: "canManageFeishuAccessToken",
    desc: "teamRole.canManageFeishuAccessTokenDesc"
  }
];
// 仓库权限key
const repoPermKeys = [
  "canAccessRepo",
  "canPushRepo",
  "canSubmitPullRequest",
  "canAddCommentInPullRequest",
  "canManageWorkflow",
  "canTriggerWorkflow",
  "canManageWebhook",
  "canManageProtectedBranch"
];
// 应用服务权限key
const appPermKeys = [
  "canDevelop",
  "canCreateDeployPlan",
  "canManagePipeline",
  "canManagePropertySource",
  "canManageServiceSource",
  "canManageDiscoverySource",
  "canDeployProperty"
];
// 仓库权限列表
const repoPermList = [
  {
    title: "teamRole.codePerm",
    perms: [
      {
        checkbox: "teamRole.canAccessRepo",
        key: "canAccessRepo",
        desc: "teamRole.canAccessRepoDesc"
      },
      {
        checkbox: "teamRole.canPushRepo",
        key: "canPushRepo",
        desc: "teamRole.canPushRepoDesc"
      }
    ]
  },
  {
    title: "teamRole.pullRequestPerm",
    perms: [
      {
        checkbox: "teamRole.canSubmitPullRequest",
        key: "canSubmitPullRequest",
        desc: "teamRole.canSubmitPullRequestDesc"
      },
      {
        checkbox: "teamRole.canAddCommentInPullRequest",
        key: "canAddCommentInPullRequest",
        desc: "teamRole.canAddCommentInPullRequestDesc"
      }
    ]
  },
  {
    title: "teamRole.workflowPerm",
    perms: [
      {
        checkbox: "teamRole.canManageWorkflow",
        key: "canManageWorkflow",
        desc: "teamRole.canManageWorkflowDesc"
      },
      {
        checkbox: "teamRole.canTriggerWorkflow",
        key: "canTriggerWorkflow",
        desc: "teamRole.canTriggerWorkflowDesc"
      }
    ]
  },
  {
    title: "teamRole.webhookPerm",
    perms: [
      {
        checkbox: "teamRole.canManageWebhook",
        key: "canManageWebhook",
        desc: "teamRole.canManageWebhookDesc"
      }
    ]
  },
  {
    title: "teamRole.protectedBranchPerm",
    perms: [
      {
        checkbox: "teamRole.canManageProtectedBranch",
        key: "canManageProtectedBranch",
        desc: "teamRole.canManageProtectedBranchDesc"
      }
    ]
  }
];
// 展示添加仓库权限modal
const showAddRepoPermModal = () => {
  // 默认权限
  repoModalState.addModalCheckboxs = {
    canAccessRepo: true,
    canPushRepo: true,
    canSubmitPullRequest: true,
    canAddCommentInPullRequest: true,
    canTriggerWorkflow: true
  };
  // 展示
  repoModalState.addModalOpen = true;
  // 清空选定仓库
  repoModalState.repoSelect = [];
};
// 展示添加应用服务权限modal
const showAddAppPermModal = () => {
  // 默认权限
  appModalState.addModalCheckboxs = {
    canDevelop: true
  };
  // 展示
  appModalState.addModalOpen = true;
  // 清空选定应用服务
  appModalState.appSelect = [];
};
// 展示编辑仓库权限modal
const showUpdateRepoPermModal = item => {
  repoModalState.updateModalOpen = true;
  repoModalState.updatePerm.repoId = item.repoId;
  repoModalState.updatePerm.name = item.name;
  repoModalState.updatePerm.perms = { ...item.perms };
  repoModalState.updatePerm.target = item;
};
// 展示编辑应用服务权限modal
const showUpdateAppPermModal = item => {
  appModalState.updateModalOpen = true;
  appModalState.updatePerm.appId = item.appId;
  appModalState.updatePerm.name = item.name;
  appModalState.updatePerm.perms = { ...item.perms };
  appModalState.updatePerm.target = item;
};
// 点击添加指定仓库权限“确定”按钮
const handleAddRepoModalOk = () => {
  if (repoModalState.repoSelect.length === 0) {
    message.warn("请选择仓库");
    return;
  }
  let permCount = 0;
  // 计算打勾权限数量
  for (let key in repoModalState.addModalCheckboxs) {
    if (repoModalState.addModalCheckboxs[key]) {
      permCount += 1;
    }
  }
  // 存储数据
  repoModalState.repoSelect.forEach(item => {
    let repo = repoModalState.remainList.find(repo => repo.repoId === item);
    formState.addRepoPermList.push({
      ...repo,
      perms: { ...repoModalState.addModalCheckboxs },
      permCount: permCount
    });
  });
  // 刷新剩余仓库数量
  refreshRemainRepoList();
  repoModalState.addModalOpen = false;
};
// 点击添加指定应用服务权限“确定”按钮
const handleAddAppModalOk = () => {
  if (appModalState.appSelect.length === 0) {
    message.warn("请选择应用服务");
    return;
  }
  let permCount = 0;
  // 计算打勾权限数量
  for (let key in appModalState.addModalCheckboxs) {
    if (appModalState.addModalCheckboxs[key]) {
      permCount += 1;
    }
  }
  // 存储数据
  appModalState.appSelect.forEach(item => {
    let app = appModalState.remainList.find(app => app.appId === item);
    formState.addAppPermList.push({
      ...app,
      perms: { ...appModalState.addModalCheckboxs },
      permCount: permCount
    });
  });
  // 刷新剩余应用服务数量
  refreshRemainAppList();
  appModalState.addModalOpen = false;
};
// 刷新剩余仓库数据
const refreshRemainRepoList = () => {
  repoModalState.remainList = repoModalState.allList.filter(repo => {
    return !formState.addRepoPermList.find(item => item.repoId === repo.repoId);
  });
};
// 刷新剩余app数量
const refreshRemainAppList = () => {
  appModalState.remainList = appModalState.allList.filter(app => {
    return !formState.addAppPermList.find(item => item.appId === app.appId);
  });
};
// 点击编辑指定仓库权限“确定”按钮
const handleUpdateRepoModalOk = () => {
  let permCount = 0;
  for (let key in repoModalState.updatePerm.perms) {
    if (repoModalState.updatePerm.perms[key]) {
      permCount += 1;
    }
  }
  // 回显数据
  repoModalState.updatePerm.target.perms = {
    ...repoModalState.updatePerm.perms
  };
  repoModalState.updatePerm.target.permCount = permCount;
  repoModalState.updateModalOpen = false;
};
// 点击编辑指定应用服务权限“确定”按钮
const handleUpdateAppModalOk = () => {
  let permCount = 0;
  for (let key in appModalState.updatePerm.perms) {
    if (appModalState.updatePerm.perms[key]) {
      permCount += 1;
    }
  }
  // 回显数据
  appModalState.updatePerm.target.perms = {
    ...appModalState.updatePerm.perms
  };
  appModalState.updatePerm.target.permCount = permCount;
  appModalState.updateModalOpen = false;
};
// 下拉框搜索选择
const filterSelectOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 删除仓库权限
const deleteAddRepoPerm = () => {
  formState.addRepoPermList = formState.addRepoPermList.filter(item => {
    return item.repoId !== repoModalState.updatePerm.repoId;
  });
  // 刷新剩余数量
  refreshRemainRepoList();
  repoModalState.updateModalOpen = false;
};
// 删除应用权限
const deleteAddAppPerm = () => {
  formState.addAppPermList = formState.addAppPermList.filter(item => {
    return item.appId !== appModalState.updatePerm.appId;
  });
  // 刷新剩余数量
  refreshRemainAppList();
  appModalState.updateModalOpen = false;
};
/* 
  获取所有仓库
  callback 为获取数据后的回调
 */
const listAllRepo = callback => {
  getRepoListByAdminRequest(route.params.teamId).then(res => {
    repoModalState.allList = res.data;
    refreshRemainRepoList();
    if (callback) {
      nextTick(() => {
        callback(res.data);
      });
    }
  });
};
// 获取所有的应用服务
const listAllApp = callback => {
  listAllAppByAdminRequest(route.params.teamId).then(res => {
    appModalState.allList = res.data;
    refreshRemainAppList();
    if (callback) {
      nextTick(() => {
        callback(res.data);
      });
    }
  });
};
// 点击“立即保存”
const createOrUpdateRole = () => {
  if (!teamRoleNameRegexp.test(formState.name)) {
    message.warn(t('teamRole.roleNameFormatErr'));
    return;
  }
  let repoPermList;
  if (formState.useDefaultRepoPerm) {
    formState.addRepoPermList.value = [];
    repoPermList = [];
  } else {
    formState.defaultRepoPerm = {};
    repoPermList = formState.addRepoPermList.map(item => {
      return {
        repoId: item.repoId,
        ...item.perms
      };
    });
  }
  let appPermList;
  if (formState.useDefaultAppPerm) {
    formState.addAppPermList.value = [];
    appPermList = [];
  } else {
    formState.defaultAppPerm = {};
    appPermList = formState.addAppPermList.map(item => {
      return {
        appId: item.appId,
        ...item.perms
      };
    });
  }
  if (mode === "create") {
    createRoleRequest({
      teamId: parseInt(route.params.teamId),
      name: formState.name,
      perm: {
        teamPerm: formState.teamPerm,
        defaultRepoPerm: formState.defaultRepoPerm,
        defaultAppPerm: formState.defaultAppPerm,
        repoPermList: repoPermList,
        appPermList: appPermList
      }
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/team/${route.params.teamId}/role/list`);
    });
  } else {
    updateRoleRequest({
      roleId: parseInt(route.params.roleId),
      name: formState.name,
      perm: {
        teamPerm: formState.teamPerm,
        defaultRepoPerm: formState.defaultRepoPerm,
        defaultAppPerm: formState.defaultAppPerm,
        repoPermList: repoPermList,
        appPermList: appPermList
      }
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/team/${route.params.teamId}/role/list`);
    });
  }
};
if (mode === "update") {
  if (teamRoleStore.roleId === 0) {
    router.push(`/team/${route.params.teamId}/role/list`);
  } else {
    listAllRepo(allRepo => {
      formState.name = teamRoleStore.name;
      formState.teamPerm = teamRoleStore.teamPerm;
      formState.defaultRepoPerm = teamRoleStore.defaultRepoPerm;
      if (teamRoleStore.repoPermList?.length > 0) {
        formState.addRepoPermList = teamRoleStore.repoPermList.map(item => {
          let findRepo = allRepo.find(repo => {
            return repo.repoId === item.repoId;
          });
          let perms = {};
          let permCount = 0;
          repoPermKeys.forEach(key => {
            perms[key] = item[key];
            if (item[key]) {
              permCount += 1;
            }
          });
          return {
            repoId: item.repoId,
            name: findRepo.name,
            permCount: permCount,
            perms: perms
          };
        });
        refreshRemainRepoList();
        formState.useDefaultRepoPerm = false;
      } else {
        formState.useDefaultRepoPerm = true;
      }
    });
    listAllApp(allApp => {
      formState.defaultAppPerm = teamRoleStore.defaultAppPerm;
      if (teamRoleStore.appPermList?.length > 0) {
        formState.addAppPermList = teamRoleStore.appPermList.map(item => {
          let findApp = allApp.find(app => {
            return app.appId === item.appId;
          });
          let perms = {};
          let permCount = 0;
          appPermKeys.forEach(key => {
            perms[key] = item[key];
            if (item[key]) {
              permCount += 1;
            }
          });
          return {
            appId: item.appId,
            name: findApp.name,
            permCount: permCount,
            perms: perms
          };
        });
        refreshRemainAppList();
        formState.useDefaultAppPerm = false;
      } else {
        formState.useDefaultAppPerm = true;
      }
    });
  }
} else {
  listAllRepo();
  listAllApp();
}
</script>
<style scoped>
.perm-title {
  font-size: 12px;
  padding-bottom: 8px;
}
.perm-list {
  font-size: 14px;
  display: flex;
  flex-wrap: wrap;
}
.perm-list > li {
  padding-right: 10px;
  width: 50%;
  margin-bottom: 16px;
}
.perm-ul > li {
  width: 100%;
}
.perm-ul > li + li {
  margin-top: 6px;
}
.add-repo-btn {
  border-top: 1px solid #d9d9d9;
  width: 100%;
  text-align: center;
  font-size: 14px;
  line-height: 42px;
}
.add-repo-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
.repo-ul {
  border-top: 1px solid #d9d9d9;
}
.repo-ul > li {
  width: 100%;
  display: flex;
  padding: 10px;
  font-size: 14px;
  justify-content: space-between;
  cursor: pointer;
}
.repo-ul > li:hover {
  background-color: #f0f0f0;
}
</style>