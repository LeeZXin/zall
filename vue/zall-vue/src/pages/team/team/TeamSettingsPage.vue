<template>
  <div style="padding:10px">
    <div class="container">
      <div class="section">
        <div class="section-title">
          <span>团队名称</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-input v-model:value="teamName" />
          </div>
          <div class="input-item">
            <a-button type="primary" @click="updateTeam">保存名称</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>危险操作</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-button type="primary" danger @click="deleteApp">删除团队</a-button>
            <div
              class="input-desc"
            >删除团队前, 会判断与团队关联的git仓库、应用服务、定时任务等数据, 若存在关联未迁移的数据, 则会阻止删除, 直到相关仓库、服务等数据被删除或被迁移到其他团队</div>
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
const teamStore = useTeamStore();
const teamName = ref(teamStore.name);
const router = useRouter();
const updateTeam = () => {
  if (!teamNameRegexp.test(teamName.value)) {
    message.warn("名称格式错误");
    return;
  }
  updateTeamRequest({
    teamId: teamStore.teamId,
    name: teamName.value
  }).then(() => {
    teamStore.name = teamName.value;
    message.success("编辑成功");
  });
};

const deleteApp = () => {
  Modal.confirm({
    title: `你确定要删除该团队吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteTeamRequest(teamStore.teamId).then(() => {
        message.success("删除成功");
        router.push("/");
      });
    },
    onCancel() {}
  });
};
</script>
<style>
</style>