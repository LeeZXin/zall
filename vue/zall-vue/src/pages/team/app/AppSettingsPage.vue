<template>
  <div style="padding:10px">
    <div class="container">
      <div class="section">
        <div class="section-title">
          <span>应用名称</span>
        </div>
        <div class="section-body">
          <div class="input-item" style="margin-bottom:10px">
            <a-input v-model:value="appName" />
          </div>
          <a-button type="primary" @click="updateApp">保存名称</a-button>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>危险操作</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-button type="primary" danger @click="deleteApp">删除应用</a-button>
            <div class="input-desc">删除应用后, 将删除跟应用相关的配置、部署流水线、相关记录等信息, 不可逆</div>
          </div>
          <div class="input-item" v-if="userStore.isAdmin">
            <a-button type="primary" danger @click="showTransferModal">迁移至其他团队</a-button>
            <div class="input-desc">将应用迁移至其他团队, 该团队成员无法再看到此应用</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <a-modal v-model:open="transferModal.open" title="迁移团队" @ok="handleTransferModalOk">
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
    message.warn("名称格式错误");
    return;
  }
  updateAppRequest({
    appId: route.params.appId,
    name: appName.value
  }).then(() => {
    message.success("编辑成功");
  });
};
// 删除应用服务
const deleteApp = () => {
  Modal.confirm({
    title: `你确定要删除该应用吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteAppRequest(route.params.appId).then(() => {
        message.success("删除成功");
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
    message.warn("请选择团队");
    return;
  }
  transferAppRequest({
    appId: route.params.appId,
    teamId: transferModal.teamId
  }).then(() => {
    message.success("迁移成功");
    router.push(`/team/${transferModal.teamId}/app/list`);
  });
};

getApp();
</script>
<style>
</style>