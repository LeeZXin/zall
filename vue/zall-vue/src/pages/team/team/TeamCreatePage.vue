<template>
  <div style="padding: 10px">
    <div class="container">
      <div class="title">创建团队</div>
      <div class="section">
        <div class="section-title">
          <span>{{t("createTeam.teamName")}}</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-input v-model:value="teamName" />
          </div>
          <div class="input-desc">不包含特殊字符,长度不得超过32</div>
        </div>
      </div>
      <div class="form-item">
        <a-button type="primary" @click="createTeam">立即创建</a-button>
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
    message.error("团队名称长度在1-32之间");
    return;
  }
  createTeamRequest({
    name: teamName.value
  }).then(() => {
    message.success("创建成功");
    setTimeout(() => {
      router.push("/index/team/list");
    }, 1000);
  });
};
</script>
<style scoped>
</style>