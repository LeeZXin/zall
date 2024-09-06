<template>
  <div style="padding: 10px">
    <div class="container">
      <div class="header">{{t('team.createTeam')}}</div>
      <div class="section">
        <div class="section-title">
          <span>{{t("team.name")}}</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-input v-model:value="teamName" />
          </div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createTeam">{{t("team.save")}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { useI18n } from "vue-i18n";
import { ref } from "vue";
import { teamNameRegexp } from "@/utils/regexp";
import { createTeamRequest } from "@/api/team/teamApi";
import { message } from "ant-design-vue";
import { useRouter } from "vue-router";
const router = useRouter();
const { t } = useI18n();
const teamName = ref("");
const createTeam = () => {
  if (!teamNameRegexp.test(teamName.value)) {
    message.error(t("team.nameFormatErr"));
    return;
  }
  createTeamRequest({
    name: teamName.value
  }).then(() => {
    message.success(t("operationSuccess"));
    setTimeout(() => {
      router.push("/index/team/list");
    }, 1000);
  });
};
</script>
<style scoped>
</style>