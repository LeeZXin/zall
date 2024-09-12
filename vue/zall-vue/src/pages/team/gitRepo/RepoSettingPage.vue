<template>
  <div style="padding:10px">
    <div class="container">
      <div class="section">
        <div class="section-title">{{t('repoSetting.repoSize')}}</div>
        <div class="section-body">
          <div
            class="input-item"
          >{{t('repoSetting.gitSize')}}: {{readableVolumeSize(repoInfo.gitSize)}}</div>
          <div
            class="input-item"
          >{{t('repoSetting.lfsSize')}}: {{readableVolumeSize(repoInfo.lfsSize)}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('repoSetting.repoOptimization')}}</div>
        <div class="section-body">
          <div class="input-item">
            <a-button type="primary" @click="triggerGc">{{t('repoSetting.triggerGc')}}</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('repoSetting.setting')}}</div>
        <div class="section-body">
          <div class="input-item">
            <div class="input-title">{{t('repoSetting.repoDesc')}}</div>
            <a-input v-model:value="repoInfo.repoDesc" />
          </div>
          <div class="input-item">
            <a-checkbox v-model:checked="repoInfo.disableLfs">{{t('repoSetting.disallowLfs')}}</a-checkbox>
            <div class="checkbox-option-desc">{{t('repoSetting.disallowLfsDesc')}}</div>
          </div>
          <div class="input-item">
            <div class="input-title">{{t('repoSetting.limitGitSize')}}</div>
            <a-input-number
              style="width:100%"
              v-model:value="repoInfo.gitLimitSize"
              :min="0"
              :precision="0"
            >
              <template #addonAfter>
                <a-select style="width: 100px" v-model:value="repoInfo.gitSizeUnit">
                  <a-select-option
                    v-for="item in volumeList"
                    v-bind:key="item"
                    :value="item"
                  >{{item}}</a-select-option>
                </a-select>
              </template>
            </a-input-number>
            <div class="input-desc">{{t('repoSetting.limitGitSizeDesc')}}</div>
          </div>
          <div class="input-item">
            <div class="input-title">{{t('repoSetting.limitLfsSize')}}</div>
            <a-input-number
              style="width:100%"
              v-model:value="repoInfo.lfsLimitSize"
              :min="0"
              :precision="0"
            >
              <template #addonAfter>
                <a-select style="width: 100px" v-model:value="repoInfo.lfsSizeUnit">
                  <a-select-option
                    v-for="item in volumeList"
                    v-bind:key="item"
                    :value="item"
                  >{{item}}</a-select-option>
                </a-select>
              </template>
            </a-input-number>
            <div class="input-desc">{{t('repoSetting.limitLfsSizeDesc')}}</div>
          </div>
          <div class="input-item">
            <a-button type="primary" @click="updateRepo">{{t('repoSetting.save')}}</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('repoSetting.dangerousAction')}}</div>
        <div class="section-body">
          <div class="input-item">
            <a-button type="primary" danger @click="deleteRepo">{{t('repoSetting.deleteRepo')}}</a-button>
            <div class="input-desc">{{t('repoSetting.deleteRepoDesc')}}</div>
          </div>
          <div class="input-item" v-if="repoInfo.isArchived">
            <a-button
              type="primary"
              danger
              @click="setArchivedStatus(false)"
            >{{t('repoSetting.unArchiveRepo')}}</a-button>
            <div class="input-desc">{{t('repoSetting.unArchiveRepoDesc')}}</div>
          </div>
          <div class="input-item" v-else>
            <a-button
              type="primary"
              danger
              @click="setArchivedStatus(true)"
            >{{t('repoSetting.archiveRepo')}}</a-button>
            <div class="input-desc">{{t('repoSetting.archiveRepoDesc')}}</div>
          </div>
          <div class="input-item" v-if="userStore.isAdmin">
            <a-button
              type="primary"
              danger
              @click="showTransferModal"
            >{{t('repoSetting.transferTeam')}}</a-button>
            <div class="input-desc">{{t('repoSetting.transferTeamDesc')}}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <a-modal
    v-model:open="transferModal.open"
    :title="t('repoSetting.transferTeam')"
    @ok="handleTransferModalOk"
  >
    <a-select
      v-model:value="transferModal.teamId"
      style="width:100%"
      :options="teamList"
      show-search
      :filter-option="filterTeamListOption"
    />
  </a-modal>
</template>
<script setup>
import { readableVolumeSize, calcUnit, Unit } from "@/utils/size";
import { reactive, createVNode, ref } from "vue";
import {
  getDetailInfoRequest,
  gcRequest,
  updateRepoRequest,
  setArchivedRequest,
  setUnArchivedRequest,
  deleteRepoRequest,
  transferRepoRequest
} from "@/api/git/repoApi";
import { ExclamationCircleOutlined } from "@ant-design/icons-vue";
import { useRoute, useRouter } from "vue-router";
import { message, Modal } from "ant-design-vue";
import { useRepoStore } from "@/pinia/repoStore";
import { useUserStore } from "@/pinia/userStore";
import { listAllByAdminRequest } from "@/api/team/teamApi";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const userStore = useUserStore();
const route = useRoute();
const router = useRouter();
const repoId = parseInt(route.params.repoId);
const volumeList = ["KB", "MB", "GB", "TB"];
// 仓库信息
const repoInfo = reactive({
  gitSize: 0,
  lfsSize: 0,
  disableLfs: false,
  lfsLimitSize: 0,
  gitLimitSize: 0,
  gitSizeUnit: "KB",
  lfsSizeUnit: "KB",
  isArchived: false,
  repoDesc: "",
  loaded: false
});
// 迁移团队modal
const transferModal = reactive({
  open: false,
  teamId: null
});
// 团队列表
const teamList = ref([]);
// 展示迁移团队modal
const showTransferModal = () => {
  if (teamList.value.length === 0) {
    listAllTeam();
  }
  transferModal.open = true;
};
// 获取团队列表
const listAllTeam = () => {
  listAllByAdminRequest().then(res => {
    let t = res.data.map(item => {
      return {
        value: item.teamId,
        label: item.name
      };
    });
    let teamId = parseInt(route.params.teamId);
    t = t.filter(item => item.value !== teamId);
    teamList.value = t;
    if (t.length > 0) {
      transferModal.teamId = t[0].value;
    }
  });
};
// 团队下拉框过滤
const filterTeamListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 确认迁移团队
const handleTransferModalOk = () => {
  if (!transferModal.teamId) {
    message.warn(t("repoSetting.pleaseSelectTeam"));
    return;
  }
  transferRepoRequest({
    repoId: repoId,
    teamId: transferModal.teamId
  }).then(() => {
    message.success(t("operationSuccess"));
    router.push(`/team/${transferModal.teamId}/gitRepo/list`);
  });
};
// 获取仓库信息
const getRepo = () => {
  getDetailInfoRequest(repoId).then(res => {
    repoInfo.gitSize = res.data.gitSize;
    repoInfo.lfsSize = res.data.lfsSize;
    if (!repoInfo.loaded) {
      repoInfo.repoDesc = res.data.repoDesc;
      repoInfo.loaded = true;
      repoInfo.disableLfs = res.data.disableLfs;
      let lfs = calcUnit(res.data.lfsLimitSize);
      repoInfo.lfsLimitSize = lfs.size;
      repoInfo.lfsSizeUnit = lfs.unit.unit;
      let git = calcUnit(res.data.gitLimitSize);
      repoInfo.gitLimitSize = git.size;
      repoInfo.gitSizeUnit = git.unit.unit;
    }
    repoInfo.isArchived = res.data.isArchived;
  });
};
// 编辑仓库
const updateRepo = () => {
  updateRepoRequest({
    repoId,
    desc: repoInfo.repoDesc,
    disableLfs: repoInfo.disableLfs,
    gitLimitSize: new Unit(repoInfo.gitSizeUnit).toNumber(
      repoInfo.gitLimitSize
    ),
    lfsLimitSize: new Unit(repoInfo.lfsSizeUnit).toNumber(repoInfo.lfsLimitSize)
  }).then(() => {
    message.success(t("operationSuccess"));
  });
};
const triggerGc = () => {
  gcRequest(repoId).then(() => {
    message.success(t("operationSuccess"));
    getRepo();
  });
};
const setArchivedStatus = isArchived => {
  let warnMsg;
  let request;
  if (isArchived) {
    warnMsg = `${t("repoSetting.confirmArchiveRepo")}?`;
    request = setArchivedRequest;
  } else {
    warnMsg = `${t("repoSetting.confirmUnArchiveRepo")}?`;
    request = setUnArchivedRequest;
  }
  Modal.confirm({
    title: warnMsg,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      request(repoId).then(() => {
        message.success(t("operationSuccess"));
        getRepo();
      });
    },
    onCancel() {}
  });
};
const deleteRepo = () => {
  Modal.confirm({
    title: `${t("repoSetting.confirmDeleteRepo")}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteRepoRequest(repoId).then(() => {
        message.success(t("operationSuccess"));
        router.push(`/team/${useRepoStore().teamId}/gitRepo/list`);
      });
    },
    onCancel() {}
  });
};
getRepo();
</script>
<style scoped>
</style>