<template>
  <div style="padding:10px">
    <div class="header">
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('teamHook.createHook')}}</a-button>
    </div>
    <ul class="team-hook-list" v-if="teamHookList.length > 0">
      <li v-for="item in teamHookList" v-bind:key="item.id">
        <div class="team-hook-name no-wrap">{{item.name}}</div>
        <ul class="op-btns">
          <li class="update-btn" @click="updateTeamHook(item)">{{t('teamHook.update')}}</li>
          <li class="del-btn" @click="deleteTeamHook(item)">{{t('teamHook.delete')}}</li>
        </ul>
      </li>
    </ul>
    <ZNoData v-else />
  </div>
</template>
<script setup>
import { ref, createVNode, h } from "vue";
import { useRouter, useRoute } from "vue-router";
import ZNoData from "@/components/common/ZNoData";
import {
  listTeamHookRequest,
  deleteTeamHookRequest
} from "@/api/team/teamHookApi";
import { ExclamationCircleOutlined, PlusOutlined } from "@ant-design/icons-vue";
import { message, Modal } from "ant-design-vue";
import { useTeamHookStore } from "@/pinia/teamHookStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const teamHookList = ref([]);
const teamHookStore = useTeamHookStore();
// 跳转新增team hook页面
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/teamHook/create`);
};
// 删除webhook
const deleteTeamHook = item => {
  Modal.confirm({
    title: `${t('teamHook.confirmDelete')} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteTeamHookRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listTeamHook();
      });
    },
    onCancel() {}
  });
};
// 获取team hook列表
const listTeamHook = () => {
  listTeamHookRequest(route.params.teamId).then(res => {
    teamHookList.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};
// 编辑team hook
const updateTeamHook = item => {
  teamHookStore.id = item.id;
  teamHookStore.name = item.name;
  teamHookStore.hookType = item.hookType;
  teamHookStore.events = item.events;
  teamHookStore.hookCfg = item.hookCfg;
  teamHookStore.teamId = item.teamId;
  router.push(`/team/${route.params.teamId}/teamHook/${item.id}/update`);
};
listTeamHook();
</script>
<style scoped>
.team-hook-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.team-hook-list > li {
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.team-hook-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.team-hook-name {
  font-size: 14px;
  line-height: 32px;
  width: 60%;
  padding-left: 10px;
}
.op-btns {
  display: flex;
  align-items: center;
}
.op-btns > li {
  line-height: 32px;
  font-size: 14px;
  padding: 0 10px;
  cursor: pointer;
}
.op-btns > li:first-child {
  border-top: 1px solid #d9d9d9;
  border-left: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
}
.op-btns > li:not(:first-child, :last-child) {
  border-top: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
}
.op-btns > li:last-child {
  border-top: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
  border-right: 1px solid #d9d9d9;
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
}
.op-btns > li + li {
  border-left: 1px solid #d9d9d9;
}
.header {
  margin-bottom: 10px;
}
.header > span {
  font-size: 18px;
  font-weight: bold;
  line-height: 32px;
  padding-left: 8px;
}
.del-btn {
  color: darkred;
}
.del-btn:hover {
  color: white;
  background-color: darkred;
}
.update-btn:hover {
  background-color: #f0f0f0;
}
.no-data-text {
  font-size: 14px;
  line-height: 18px;
  padding: 10px;
  text-align: center;
  word-break: break-all;
}
</style>