<template>
  <div style="padding:10px">
    <div class="container">
      <div class="section">
        <div class="section-title">
          <span>{{t('teamSetting.teamName')}}</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-input v-model:value="teamName" />
          </div>
          <div class="input-item">
            <a-button type="primary" @click="updateTeam">{{t('teamSetting.saveTeamName')}}</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>{{t('teamSetting.dangerousAction')}}</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-button type="primary" danger @click="deleteTeam">{{t('teamSetting.deleteTeam')}}</a-button>
            <div class="input-desc">{{t('teamSetting.deleteTeamDesc')}}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, createVNode } from "vue";
import { updateTeamRequest, deleteTeamRequest } from "@/api/team/teamApi";
import { teamNameRegexp } from "@/utils/regexp";
import { message, Modal } from "ant-design-vue";
import { ExclamationCircleOutlined } from "@ant-design/icons-vue";
import { useTeamStore } from "@/pinia/teamStore";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const teamStore = useTeamStore();
const teamName = ref(teamStore.name);
const router = useRouter();
// 编辑团队
const updateTeam = () => {
  if (!teamNameRegexp.test(teamName.value)) {
    message.warn(t("teamSetting.teamNameFormatErr"));
    return;
  }
  updateTeamRequest({
    teamId: teamStore.teamId,
    name: teamName.value
  }).then(() => {
    teamStore.name = teamName.value;
    message.success(t("operationSuccess"));
  });
};
// 删除团队
const deleteTeam = () => {
  Modal.confirm({
    title: `${t("teamSetting.confirmDeleteTeam")}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteTeamRequest(teamStore.teamId).then(() => {
        message.success(t("operationSuccess"));
        router.push("/");
      });
    },
    onCancel() {}
  });
};
</script>
<style>
</style>