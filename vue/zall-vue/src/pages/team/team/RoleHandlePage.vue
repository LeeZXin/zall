<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode==='create'">创建角色</span>
        <span v-else-if="mode==='update'">编辑角色</span>
      </div>
      <div class="section">
        <div class="section-title">角色名称</div>
        <div class="section-body">
          <div class="input-item">
            <a-input v-model:value="formState.name" />
            <div class="input-desc">简短的角色名称</div>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">团队权限</div>
        <div class="section-body">
          <ul class="perm-list">
            <li v-for="item in teamPermList" v-bind:key="item.key">
              <a-checkbox v-model:checked="formState.teamPerm[item.key]">{{item.checkbox}}</a-checkbox>
              <div class="checkbox-desc">{{item.desc}}</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="section">
        <div class="section-title">仓库权限选项</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.useDefaultRepoPerm">
            <a-radio :value="true">使用默认仓库权限</a-radio>
            <div class="radio-option-desc">所有仓库将应用一套权限配置</div>
            <a-radio :value="false">指定仓库权限</a-radio>
            <div class="radio-option-desc">针对不同仓库有不同的权限配置</div>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.useDefaultRepoPerm">
        <div class="section-title">默认仓库权限</div>
        <div class="section-body">
          <ul class="perm-ul">
            <li v-for="item in repoPermList" v-bind:key="item.title">
              <div class="perm-title">{{item.title}}</div>
              <ul class="perm-list">
                <li v-for="(perm, index) in item.perms" v-bind:key="index">
                  <a-checkbox
                    v-model:checked="formState.defaultRepoPerm[perm.key]"
                  >{{perm.checkbox}}</a-checkbox>
                  <div class="checkbox-desc">{{perm.desc}}</div>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-else>
        <div class="section-title">指定仓库权限</div>
        <ul class="repo-ul" v-if="formState.addRepoPermList.length > 0">
          <li
            v-for="item in formState.addRepoPermList"
            v-bind:key="item.repoId"
            @click="showUpdateRepoPermModal(item)"
          >
            <span>{{item.name}}</span>
            <span style="color:green">+{{item.permCount}}项权限</span>
          </li>
        </ul>
        <div
          class="add-repo-btn"
          @click="showAddRepoPermModal"
          v-if="repoModalState.remainList.length > 0"
        >
          <PlusOutlined />
          <span>添加仓库权限</span>
        </div>
      </div>
      <div class="section">
        <div class="section-title">应用服务权限选项</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.useDefaultAppPerm">
            <a-radio :value="true">使用默认应用服务权限</a-radio>
            <div class="radio-option-desc">所有应用服务将应用一套权限配置</div>
            <a-radio :value="false">指定应用服务权限</a-radio>
            <div class="radio-option-desc">针对不同应用服务有不同的权限配置</div>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.useDefaultAppPerm">
        <div class="section-title">默认应用服务权限</div>
        <div class="section-body">
          <ul class="perm-list">
            <li v-for="item in appPermList" v-bind:key="item.key">
              <a-checkbox v-model:checked="formState.defaultAppPerm[item.key]">{{item.checkbox}}</a-checkbox>
              <div class="checkbox-desc">{{item.desc}}</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-else>
        <div class="section-title">指定应用服务权限</div>
        <ul class="repo-ul" v-if="formState.addAppPermList.length > 0">
          <li
            v-for="item in formState.addAppPermList"
            v-bind:key="item.appId"
            @click="showUpdateAppPermModal(item)"
          >
            <span>{{item.name}}</span>
            <span style="color:green">+{{item.permCount}}项权限</span>
          </li>
        </ul>
        <div
          class="add-repo-btn"
          @click="showAddAppPermModal"
          v-if="appModalState.remainList.length > 0"
        >
          <PlusOutlined />
          <span>添加应用服务权限</span>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createOrUpdateRole">立即保存</a-button>
      </div>
    </div>
    <a-modal v-model:open="repoModalState.addModalOpen" title="添加仓库权限" @ok="handleAddRepoModalOk">
      <template #footer>
        <a-button type="primary" @click="handleAddRepoModalOk">确定</a-button>
      </template>
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">多选仓库</div>
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
          <div class="perm-title">{{item.title}}</div>
          <ul class="perm-list">
            <li v-for="(perm, index) in item.perms" v-bind:key="index">
              <a-checkbox
                v-model:checked="repoModalState.addModalCheckboxs[perm.key]"
              >{{perm.checkbox}}</a-checkbox>
              <div class="checkbox-desc">{{perm.desc}}</div>
            </li>
          </ul>
        </li>
      </ul>
    </a-modal>
    <a-modal v-model:open="repoModalState.updateModalOpen" title="编辑仓库权限">
      <template #footer>
        <a-button type="primary" @click="deleteAddRepoPerm" danger>删除</a-button>
        <a-button type="primary" @click="handleUpdateRepoModalOk">确定</a-button>
      </template>
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">仓库</div>
        <div style="font-size:16px">{{repoModalState.updatePerm.name}}</div>
      </div>
      <ul class="perm-ul">
        <li v-for="item in repoPermList" v-bind:key="item.title">
          <div class="perm-title">{{item.title}}</div>
          <ul class="perm-list">
            <li v-for="(perm, index) in item.perms" v-bind:key="index">
              <a-checkbox
                v-model:checked="repoModalState.updatePerm.perms[perm.key]"
              >{{perm.checkbox}}</a-checkbox>
              <div class="checkbox-desc">{{perm.desc}}</div>
            </li>
          </ul>
        </li>
      </ul>
    </a-modal>
    <a-modal v-model:open="appModalState.addModalOpen" title="添加应用服务权限" @ok="handleAddAppModalOk">
      <template #footer>
        <a-button type="primary" @click="handleAddAppModalOk">确定</a-button>
      </template>
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">多选应用服务</div>
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
          <a-checkbox v-model:checked="appModalState.addModalCheckboxs[item.key]">{{item.checkbox}}</a-checkbox>
          <div class="checkbox-desc">{{item.desc}}</div>
        </li>
      </ul>
    </a-modal>
    <a-modal v-model:open="appModalState.updateModalOpen" title="编辑应用服务权限">
      <template #footer>
        <a-button type="primary" @click="deleteAddAppPerm" danger>删除</a-button>
        <a-button type="primary" @click="handleUpdateAppModalOk">确定</a-button>
      </template>
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">应用服务</div>
        <div style="font-size:16px">{{appModalState.updatePerm.name}}</div>
      </div>
      <ul class="perm-list">
        <li v-for="item in appPermList" v-bind:key="item.key">
          <a-checkbox v-model:checked="appModalState.updatePerm.perms[item.key]">{{item.checkbox}}</a-checkbox>
          <div class="checkbox-desc">{{item.desc}}</div>
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
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const teamRoleStore = useTeamRoleStore();
const mode = getMode();
const router = useRouter();
const formState = reactive({
  name: "",
  teamPerm: {},
  defaultRepoPerm: {
    canAccessRepo: true,
    canPushRepo: true,
    canSubmitPullRequest: true,
    canAddCommentInPullRequest: true,
    canTriggerWorkflow: true
  },
  defaultAppPerm: {
    canDevelop: true
  },
  useDefaultRepoPerm: true,
  useDefaultAppPerm: true,
  addRepoPermList: [],
  addAppPermList: []
});
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
const appPermList = [
  {
    checkbox: "开发应用服务",
    key: "canDevelop",
    desc: "拥有该权限,则可以查看应用服务、发布配置、查看服务状态等"
  },
  {
    checkbox: "创建发布计划",
    key: "canCreateDeployPlan",
    desc: "拥有该权限,则可以创建发布计划来部署服务"
  },
  {
    checkbox: "管理部署流水线",
    key: "canManagePipeline",
    desc: "拥有该权限,则可以管理应用部署流水线"
  },
  {
    checkbox: "管理注册中心来源",
    key: "canManageDiscoverySource",
    desc: "拥有该权限,则可以管理注册中心来源与应用服务的绑定"
  },
  {
    checkbox: "管理配置中心来源",
    key: "canManagePropertySource",
    desc: "拥有该权限,则可以管理配置中心来源与应用服务的绑定"
  },
  {
    checkbox: "发布配置",
    key: "canDeployProperty",
    desc: "拥有该权限,则可以自由发布配置"
  },
  {
    checkbox: "管理服务状态来源",
    key: "canManageServiceSource",
    desc: "拥有该权限,则可以管理服务状态来源与应用服务的绑定"
  }
];
const teamPermList = [
  {
    checkbox: "管理定时任务",
    key: "canManageTimer",
    desc: "拥有该权限,则可以对定时任务新增、查看、编辑、触发、删除"
  },
  {
    checkbox: "管理外部通知",
    key: "canManageNotifyTpl",
    desc: "拥有该权限,则可以对外部通知新增、查看、编辑、删除"
  },
  {
    checkbox: "管理Team Hook",
    key: "canManageTeamHook",
    desc: "拥有该权限,则可以对Team Hook新增、查看、编辑、删除"
  }
];
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
const appPermKeys = [
  "canDevelop",
  "canCreateDeployPlan",
  "canManagePipeline",
  "canManagePropertySource",
  "canManageServiceSource",
  "canManageDiscoverySource",
  "canDeployProperty"
];
const repoPermList = [
  {
    title: "代码权限",
    perms: [
      {
        checkbox: "访问、拉取代码",
        key: "canAccessRepo",
        desc: "对代码仓库有读的权限"
      },
      {
        checkbox: "推送代码",
        key: "canPushRepo",
        desc: "对代码仓库有写的权限"
      }
    ]
  },
  {
    title: "合并请求权限",
    perms: [
      {
        checkbox: "发起合并请求",
        key: "canSubmitPullRequest",
        desc: "可发起合并请求"
      },
      {
        checkbox: "发表评论",
        key: "canAddCommentInPullRequest",
        desc: "可在合并请求里发表评论"
      }
    ]
  },
  {
    title: "工作流权限",
    perms: [
      {
        checkbox: "管理工作流",
        key: "canManageWorkflow",
        desc: "新增、编辑、删除工作流"
      },
      {
        checkbox: "触发工作流",
        key: "canTriggerWorkflow",
        desc: "可手动触发工作流"
      }
    ]
  },
  {
    title: "webhook权限",
    perms: [
      {
        checkbox: "管理webhook",
        key: "canManageWebhook",
        desc: "查看、新增、编辑、删除webhook"
      }
    ]
  },
  {
    title: "保护分支",
    perms: [
      {
        checkbox: "管理保护分支",
        key: "canManageProtectedBranch",
        desc: "查看、新增、编辑、删除保护分支"
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
    message.warn("名称不正确");
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
      message.success("添加成功");
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
      message.success("编辑成功");
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