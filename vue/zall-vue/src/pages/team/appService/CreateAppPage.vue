<template>
  <div style="padding:14px">
    <ZNaviBack url="/team/app/list" name="应用服务列表" />
    <div class="container">
      <div class="body">
        <div class="title">创建应用服务</div>
        <div class="form-item">
          <div class="label">
            <span>所属团队</span>
          </div>
          <div class="form-item-text">{{teamName}}</div>
        </div>
        <div class="form-item">
          <div class="label">
            <span>AppId</span>
          </div>
          <div>
            <a-input type="input" v-model:value="formState.appId" />
          </div>
        </div>
        <div class="form-item">
          <div class="label">
            <span>应用名称</span>
          </div>
          <div>
            <a-input type="input" v-model:value="formState.appName" />
          </div>
        </div>
        <div class="form-item">
          <a-button type="primary" style="margin-top:20px;">立即创建</a-button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import { useTeamStore } from "@/pinia/teamStore";
import ZNaviBack from "@/components/common/ZNaviBack";
const team = useTeamStore();
const teamName = team.teamName;
const formState = reactive({
  appId: "",
  appName: ""
});
</script>
<style scoped>
.body > .title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
}
</style>