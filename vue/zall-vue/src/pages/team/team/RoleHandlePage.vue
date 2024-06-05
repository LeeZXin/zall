<template>
  <div style="padding:14px">
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
        <ul class="repo-ul" v-if="addRepoPermList.length > 0">
          <li
            v-for="item in addRepoPermList"
            v-bind:key="item.name"
            @click="showUpdateRepoPermModal(item)"
          >
            <span>{{item.name}}</span>
            <span style="color:green">+{{item.permCount}}项权限</span>
          </li>
        </ul>
        <div class="add-repo-btn" @click="showAddRepoPermModal" v-if="repoList.length > 0">
          <PlusOutlined />
          <span>添加仓库权限</span>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>可开发应用服务列表</span>
          <span class="select-all-btn" @click="selectAllApp">全选</span>
        </div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.appList"
            :options="allAppList"
            mode="multiple"
          />
          <div class="input-desc">可对应用服务配置中心、发布部署等有权限操作</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createOrUpdateRole">立即保存</a-button>
      </div>
    </div>
    <a-modal v-model:open="addPermModalOpen" title="添加仓库权限" @ok="handleAddModalOk">
      <template #footer>
        <a-button type="primary" @click="handleAddModalOk">确定</a-button>
      </template>
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">多选仓库</div>
        <a-select
          style="width:100%"
          v-model:value="repoSelect"
          :options="repoList.map(item=>({ value: item.repoId, label: item.name }))"
          show-search
          mode="multiple"
          :filter-option="filterRepoListOption"
        />
      </div>
      <ul class="perm-ul">
        <li v-for="item in repoPermList" v-bind:key="item.title">
          <div class="perm-title">{{item.title}}</div>
          <ul class="perm-list">
            <li v-for="(perm, index) in item.perms" v-bind:key="index">
              <a-checkbox v-model:checked="addModalCheckboxs[perm.key]">{{perm.checkbox}}</a-checkbox>
              <div class="checkbox-desc">{{perm.desc}}</div>
            </li>
          </ul>
        </li>
      </ul>
    </a-modal>
    <a-modal v-model:open="updatePermModalOpen" title="编辑仓库权限">
      <template #footer>
        <a-button type="primary" @click="deleteAddRepoPerm" danger>删除</a-button>
        <a-button type="primary" @click="handleUpdateRepoOk">确定</a-button>
      </template>
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">仓库</div>
        <div style="font-size:16px">{{updateRepoPerm.name}}</div>
      </div>
      <ul class="perm-ul">
        <li v-for="item in repoPermList" v-bind:key="item.title">
          <div class="perm-title">{{item.title}}</div>
          <ul class="perm-list">
            <li v-for="(perm, index) in item.perms" v-bind:key="index">
              <a-checkbox v-model:checked="updateRepoPerm.perms[perm.key]">{{perm.checkbox}}</a-checkbox>
              <div class="checkbox-desc">{{perm.desc}}</div>
            </li>
          </ul>
        </li>
      </ul>
    </a-modal>
  </div>
</template>
<script setup>
import { PlusOutlined } from "@ant-design/icons-vue";
import { ref, reactive, nextTick } from "vue";
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
const addPermModalOpen = ref(false);
const updatePermModalOpen = ref(false);
const addModalCheckboxs = ref({});
const allAppList = ref([]);
const updateRepoPerm = reactive({
  repoId: 0,
  name: "",
  perms: {},
  target: null
});
const formState = reactive({
  name: "",
  teamPerm: {
    canCreateRepo: true,
    canManageTimer: true
  },
  defaultRepoPerm: {
    canAccessRepo: true,
    canPushRepo: true,
    canSubmitPullRequest: true,
    canReviewPullRequest: true,
    canAddCommentInPullRequest: true,
    canTriggerWorkflow: true
  },
  useDefaultRepoPerm: true,
  appList: []
});
const addRepoPermList = ref([]);
const repoSelect = ref([]);
const allRepoList = ref([]);
const repoList = ref([]);
const teamPermList = [
  {
    checkbox: "创建仓库",
    key: "canCreateRepo",
    desc: "拥有该权限,则可以创建仓库"
  },
  {
    checkbox: "管理定时任务",
    key: "canManageTimer",
    desc: "拥有该权限,则可以对定时任务新增、查看、编辑、触发、删除"
  },
  {
    checkbox: "创建发布计划",
    key: "canCreateDeployPlan",
    desc: "拥有该权限,则可以创建发布计划来部署服务"
  },
  {
    checkbox: "管理服务部署配置",
    key: "canManageDeployConfig",
    desc: "拥有该权限,则可以对服务部署新增、查看、编辑、触发、删除"
  }
];
const p = [
  "canAccessRepo",
  "canPushRepo",
  "canSubmitPullRequest",
  "canReviewPullRequest",
  "canAddCommentInPullRequest",
  "canManageWorkflow",
  "canTriggerWorkflow",
  "canManageWebhook"
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
        checkbox: "评审合并请求",
        key: "canReviewPullRequest",
        desc: "可评审并同意合并请求"
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
        desc: "新增、编辑、删除webhook"
      }
    ]
  }
];
const showAddRepoPermModal = () => {
  addModalCheckboxs.value = {
    canAccessRepo: true,
    canPushRepo: true,
    canSubmitPullRequest: true,
    canReviewPullRequest: true,
    canAddCommentInPullRequest: true,
    canTriggerWorkflow: true
  };
  addPermModalOpen.value = true;
  repoSelect.value = [];
};
const showUpdateRepoPermModal = item => {
  updatePermModalOpen.value = true;
  updateRepoPerm.repoId = item.repoId;
  updateRepoPerm.name = item.name;
  updateRepoPerm.perms = { ...item.perms };
  updateRepoPerm.target = item;
};
const handleAddModalOk = () => {
  if (repoSelect.value.length === 0) {
    message.warn("请选择仓库");
    return;
  }
  let permCount = 0;
  for (let key in addModalCheckboxs.value) {
    if (addModalCheckboxs.value[key]) {
      permCount += 1;
    }
  }
  repoSelect.value.forEach(item => {
    let repo = repoList.value.find(repo => repo.repoId === item);
    addRepoPermList.value.push({
      ...repo,
      perms: addModalCheckboxs.value,
      permCount: permCount
    });
  });
  repoList.value = allRepoList.value.filter(repo => {
    return !addRepoPermList.value.find(item => item.repoId === repo.repoId);
  });
  addPermModalOpen.value = false;
};
const handleUpdateRepoOk = () => {
  let permCount = 0;
  for (let key in updateRepoPerm.perms) {
    if (updateRepoPerm.perms[key]) {
      permCount += 1;
    }
  }
  updateRepoPerm.target.perms = { ...updateRepoPerm.perms };
  updateRepoPerm.target.permCount = permCount;
  updatePermModalOpen.value = false;
};
const filterRepoListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
const deleteAddRepoPerm = () => {
  addRepoPermList.value = addRepoPermList.value.filter(item => {
    return item.repoId !== updateRepoPerm.repoId;
  });
  repoList.value = allRepoList.value.filter(repo => {
    return !addRepoPermList.value.find(item => item.repoId === repo.repoId);
  });
  updatePermModalOpen.value = false;
};
const getAllRepoList = callback => {
  getRepoListByAdminRequest(route.params.teamId).then(res => {
    allRepoList.value = res.data;
    repoList.value = [...res.data];
    nextTick(() => {
      if (callback) {
        callback(res.data);
      }
    });
  });
};
const listAllApp = () => {
  listAllAppByAdminRequest(route.params.teamId).then(res => {
    allAppList.value = res.data.map(item => {
      return {
        value: item.appId,
        label: item.name
      };
    });
  });
};
const selectAllApp = () => {
  formState.appList = allAppList.value.map(item => item.value);
};
const createOrUpdateRole = () => {
  if (!teamRoleNameRegexp.test(formState.name)) {
    message.warn("名称不正确");
    return;
  }
  let repoPermList = [];
  if (formState.useDefaultRepoPerm) {
    addRepoPermList.value = [];
  } else {
    formState.defaultRepoPerm = {};
    repoPermList = addRepoPermList.value.map(item => {
      return {
        repoId: item.repoId,
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
        repoPermList: repoPermList,
        developAppList: formState.appList
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
        repoPermList: repoPermList,
        developAppList: formState.appList
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
    getAllRepoList(data => {
      formState.name = teamRoleStore.name;
      formState.teamPerm = teamRoleStore.teamPerm;
      formState.defaultRepoPerm = teamRoleStore.defaultRepoPerm;
      formState.appList = teamRoleStore.developAppList;
      if (teamRoleStore.repoPermList && teamRoleStore.repoPermList.length > 0) {
        addRepoPermList.value = teamRoleStore.repoPermList.map(item => {
          let r = data.find(repo => {
            return repo.repoId === item.repoId;
          });
          let perms = {};
          let permCount = 0;
          p.forEach(key => {
            perms[key] = item[key];
            if (item[key]) {
              permCount += 1;
            }
          });
          return {
            repoId: item.repoId,
            name: r.name,
            permCount: permCount,
            perms: perms
          };
        });
        repoList.value = allRepoList.value.filter(ar => {
          return !addRepoPermList.value.find(ap => ap.repoId === ar.repoId);
        });
        formState.useDefaultRepoPerm = false;
      } else {
        formState.useDefaultRepoPerm = true;
      }
    });
  }
} else {
  getAllRepoList();
}
listAllApp();
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
.select-all-btn {
  font-size: 14px;
  float: right;
}
.select-all-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>