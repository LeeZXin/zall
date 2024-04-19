<template>
  <div style="padding:14px">
    <div style="margin-bottom:20px">
      <span class="header" @click="backToAppList">
        <arrow-left-outlined />
        <span style="margin-left:8px">应用服务列表</span>
      </span>
    </div>
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
import { ArrowLeftOutlined } from "@ant-design/icons-vue";
import { useTeamStore } from "@/pinia/TeamStore";
import { useRouter } from "vue-router";
const team = useTeamStore();
const teamName = team.teamName;
const router = useRouter();
const formState = reactive({
  appId: "",
  appName: ""
});
const backToAppList = () => {
  router.push("/team/app/list");
};
</script>
<style scoped>
.header {
  font-size: 14px;
  cursor: pointer;
  margin-bottom: 10px;
}
.header:hover {
  color: #1677ff;
}
.body > .title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
}
</style>