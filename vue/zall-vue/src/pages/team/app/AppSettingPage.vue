<template>
  <div style="padding:10px">
    <div class="container">
      <div class="section">
        <div class="section-title">
          <span>{{t('appSetting.appName')}}</span>
        </div>
        <div class="section-body">
          <div class="input-item" style="margin-bottom:10px">
            <a-input v-model:value="appName" />
          </div>
          <a-button type="primary" @click="updateApp">{{t('appSetting.saveAppName')}}</a-button>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>{{t('appSetting.dangerousAction')}}</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-button type="primary" danger @click="deleteApp">{{t('appSetting.deleteApp')}}</a-button>
            <div class="input-desc">{{t('appSetting.deleteAppDesc')}}</div>
          </div>
          <div class="input-item" v-if="userStore.isAdmin">
            <a-button
              type="primary"
              danger
              @click="showTransferModal"
            >{{t('appSetting.transferAppToOtherTeam')}}</a-button>
          </div>
        </div>
      </div>
    </div>
  </div>
  <a-modal
    v-model:open="transferModal.open"
    :title="t('appSetting.transferAppToOtherTeam')"
    @ok="handleTransferModalOk"
  >
    <div style="font-size:12px;margin-bottom:6px">{{t('appSetting.selectTeam')}}</div>
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
import { ref, createVNode, reactive } from "vue";
import {
  getAppRequest,
  updateAppRequest,
  deleteAppRequest,
  transferAppRequest
} from "@/api/app/appApi";
import { useRoute, useRouter } from "vue-router";
import { appNameRegexp } from "@/utils/regexp";
import { message, Modal } from "ant-design-vue";
import { ExclamationCircleOutlined } from "@ant-design/icons-vue";
import { useUserStore } from "@/pinia/userStore";
import { listAllByAdminRequest } from "@/api/team/teamApi";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const userStore = useUserStore();
const appName = ref("");
const route = useRoute();
const router = useRouter();
// 迁移modal
const transferModal = reactive({
  open: false,
  teamId: null
});
const teamList = ref([]);
// 获取应用服务名称
const getApp = () => {
  getAppRequest(route.params.appId).then(res => {
    appName.value = res.data.name;
  });
};
// 编辑应用服务名称
const updateApp = () => {
  if (!appNameRegexp.test(appName.value)) {
    message.warn(t("appSetting.appNameFormatErr"));
    return;
  }
  updateAppRequest({
    appId: route.params.appId,
    name: appName.value
  }).then(() => {
    message.success(t("operationSuccess"));
  });
};
// 删除应用服务
const deleteApp = () => {
  Modal.confirm({
    title: `${t("appSetting.confirmDeleteApp")}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteAppRequest(route.params.appId).then(() => {
        message.success(t("operationSuccess"));
        router.push(`/team/${route.params.teamId}/app/list`);
      });
    },
    onCancel() {}
  });
};
// 展示迁移modal
const showTransferModal = () => {
  if (teamList.value.length === 0) {
    listAllTeam();
  }
  transferModal.open = true;
};
// 获取所有团队
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
// 下拉框过滤
const filterTeamListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 迁移团队点击“确定”
const handleTransferModalOk = () => {
  if (!transferModal.teamId) {
    message.warn(t("appSetting.pleaseSelectTeam"));
    return;
  }
  transferAppRequest({
    appId: route.params.appId,
    teamId: transferModal.teamId
  }).then(() => {
    message.success(t("operationSuccess"));
    router.push(`/team/${transferModal.teamId}/app/list`);
  });
};

getApp();
</script>
<style>
</style>